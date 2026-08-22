// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/pkg/encryption"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	encKeyFlag     string
	encKeyFileFlag string
	dbPathFlag     string
	dryRunFlag     bool
	decryptFlag    bool
)

var linkDocEncryptDatastore = map[string]string{
	"link":    "https://docs.meshery.io/reference/references/mesheryctl/system/encrypt-datastore",
	"caption": "Usage of mesheryctl system encrypt-datastore",
}

// encryptDatastoreCmd represents the command to encrypt or decrypt sensitive data at rest
var encryptDatastoreCmd = &cobra.Command{
	Use:   "encrypt-datastore",
	Short: "Encrypt or decrypt sensitive credential and kubeconfig data at rest in the Meshery datastore",
	Long: `Encrypt or decrypt sensitive columns (credential secrets and kubeconfig auth/cluster sections)
in the Meshery SQLite datastore in-place using AES-256-GCM.

Before running this command, ensure the Meshery server is STOPPED to prevent concurrent database access.

To encrypt:
  mesheryctl system encrypt-datastore --key <64-char-hex>
  mesheryctl system encrypt-datastore --key-file /path/to/key.txt

To decrypt (restore plaintext):
  mesheryctl system encrypt-datastore --decrypt --key <64-char-hex>

To preview changes without modifying the database:
  mesheryctl system encrypt-datastore --dry-run --key <64-char-hex>`,
	Example: `
// Encrypt datastore with an explicit hex key
mesheryctl system encrypt-datastore --key 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef

// Encrypt datastore using a key stored in a file
mesheryctl system encrypt-datastore --key-file ~/.meshery/key.txt

// Dry-run preview
mesheryctl system encrypt-datastore --dry-run --key-file ~/.meshery/key.txt

// Decrypt datastore back to plaintext
mesheryctl system encrypt-datastore --decrypt --key-file ~/.meshery/key.txt
`,
	Annotations: linkDocEncryptDatastore,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEncryptDatastore()
	},
}

func init() {
	encryptDatastoreCmd.Flags().StringVar(&encKeyFlag, "key", "", "32-byte AES-256 key in hex (64 chars) or base64 format")
	encryptDatastoreCmd.Flags().StringVar(&encKeyFileFlag, "key-file", "", "Path to a file containing the 32-byte key")
	encryptDatastoreCmd.Flags().StringVar(&dbPathFlag, "db-path", "", "Path to the mesherydb.sql file (defaults to ~/.meshery/config/mesherydb.sql)")
	encryptDatastoreCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Preview changes without modifying the datastore")
	encryptDatastoreCmd.Flags().BoolVar(&decryptFlag, "decrypt", false, "Decrypt existing encrypted rows back to plaintext")
}

func runEncryptDatastore() error {
	// 1. Resolve key
	if encKeyFlag != "" {
		_ = os.Setenv(encryption.EncryptionKeyEnv, encKeyFlag)
	}
	if encKeyFileFlag != "" {
		_ = os.Setenv(encryption.EncryptionKeyFileEnv, encKeyFileFlag)
	}

	svc, err := encryption.NewFromEnv()
	if err != nil {
		return errors.Wrap(err, "invalid encryption key")
	}
	if svc == nil {
		return fmt.Errorf("no encryption key specified. Provide --key, --key-file, or set %s", encryption.EncryptionKeyEnv)
	}

	// 2. Resolve database path
	dbPath := dbPathFlag
	if dbPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return errors.Wrap(err, "unable to find user home directory")
		}
		dbPath = filepath.Join(home, ".meshery", "config", "mesherydb.sql")
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database file does not exist at %s", dbPath)
	}

	action := "Encrypting"
	if decryptFlag {
		action = "Decrypting"
	}
	if dryRunFlag {
		utils.Log.Infof("[DRY-RUN] %s datastore at %s...", action, dbPath)
	} else {
		if !utils.SilentFlag {
			prompt := fmt.Sprintf("%s datastore at %s. Please ensure Meshery server is stopped. Continue?", action, dbPath)
			if !utils.AskForConfirmation(prompt) {
				utils.Log.Info("Operation cancelled.")
				return nil
			}
		}
		utils.Log.Infof("%s datastore at %s...", action, dbPath)
	}

	// 3. Connect to database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to open database at %s", dbPath)
	}

	// 4. Migrate credentials table
	credCount, err := migrateCredentialsTable(db, svc, decryptFlag, dryRunFlag)
	if err != nil {
		return errors.Wrap(err, "error processing credentials table")
	}

	// 5. Migrate k8s_contexts table
	k8sCount, err := migrateK8sContextsTable(db, svc, decryptFlag, dryRunFlag)
	if err != nil {
		return errors.Wrap(err, "error processing k8s_contexts table")
	}

	verb := "Encrypted"
	if decryptFlag {
		verb = "Decrypted"
	}
	if dryRunFlag {
		verb = "Would have " + verb
	}

	utils.Log.Infof("Migration complete: %s %d credential(s) and %d Kubernetes context(s).", verb, credCount, k8sCount)
	return nil
}

// rawCredential holds raw row data from the credentials table.
type rawCredential struct {
	ID     string `gorm:"column:id;primaryKey"`
	Secret string `gorm:"column:secret"`
}

func (rawCredential) TableName() string {
	return "credentials"
}

func migrateCredentialsTable(db *gorm.DB, svc *encryption.Service, decrypt bool, dryRun bool) (int, error) {
	if !db.Migrator().HasTable("credentials") {
		return 0, nil
	}

	var rows []rawCredential
	if err := db.Find(&rows).Error; err != nil {
		return 0, err
	}

	count := 0
	for _, row := range rows {
		if row.Secret == "" {
			continue
		}

		var secretMap map[string]interface{}
		if err := json.Unmarshal([]byte(row.Secret), &secretMap); err != nil {
			return count, fmt.Errorf("credentials row %s: column secret is not a valid JSON map: %w", row.ID, err)
		}

		if decrypt {
			if !svc.IsEncrypted(secretMap) {
				continue // Already plaintext
			}
			decrypted, err := svc.DecryptMap(secretMap)
			if err != nil {
				return count, fmt.Errorf("failed to decrypt credential %s: %w", row.ID, err)
			}
			if !dryRun {
				bytes, _ := json.Marshal(decrypted)
				if err := db.Model(&rawCredential{}).Where("id = ?", row.ID).Update("secret", string(bytes)).Error; err != nil {
					return count, err
				}
			}
			count++
		} else {
			if svc.IsEncrypted(secretMap) {
				continue // Already encrypted
			}
			encrypted, err := svc.EncryptMap(secretMap)
			if err != nil {
				return count, fmt.Errorf("failed to encrypt credential %s: %w", row.ID, err)
			}
			if !dryRun {
				bytes, _ := json.Marshal(encrypted)
				if err := db.Model(&rawCredential{}).Where("id = ?", row.ID).Update("secret", string(bytes)).Error; err != nil {
					return count, err
				}
			}
			count++
		}
	}

	return count, nil
}

// rawK8sContext holds raw row data from the k8s_contexts table.
type rawK8sContext struct {
	ID      string `gorm:"column:id;primaryKey"`
	Auth    string `gorm:"column:auth"`
	Cluster string `gorm:"column:cluster"`
}

func (rawK8sContext) TableName() string {
	return "k8s_contexts"
}

func migrateK8sContextsTable(db *gorm.DB, svc *encryption.Service, decrypt bool, dryRun bool) (int, error) {
	if !db.Migrator().HasTable("k8s_contexts") {
		return 0, nil
	}

	var rows []rawK8sContext
	if err := db.Find(&rows).Error; err != nil {
		return 0, err
	}

	count := 0
	for _, row := range rows {
		modified := false
		newAuth := row.Auth
		newCluster := row.Cluster

		// Handle Auth column
		if row.Auth != "" {
			var authMap map[string]interface{}
			if err := json.Unmarshal([]byte(row.Auth), &authMap); err != nil {
				return count, fmt.Errorf("k8s_contexts row %s: column auth is not a valid JSON map: %w", row.ID, err)
			}
			if decrypt && svc.IsEncrypted(authMap) {
				decrypted, err := svc.DecryptMap(authMap)
				if err != nil {
					return count, fmt.Errorf("failed to decrypt auth for context %s: %w", row.ID, err)
				}
				b, _ := json.Marshal(decrypted)
				newAuth = string(b)
				modified = true
			} else if !decrypt && !svc.IsEncrypted(authMap) {
				encrypted, err := svc.EncryptMap(authMap)
				if err != nil {
					return count, fmt.Errorf("failed to encrypt auth for context %s: %w", row.ID, err)
				}
				b, _ := json.Marshal(encrypted)
				newAuth = string(b)
				modified = true
			}
		}

		// Handle Cluster column
		if row.Cluster != "" {
			var clusterMap map[string]interface{}
			if err := json.Unmarshal([]byte(row.Cluster), &clusterMap); err != nil {
				return count, fmt.Errorf("k8s_contexts row %s: column cluster is not a valid JSON map: %w", row.ID, err)
			}
			if decrypt && svc.IsEncrypted(clusterMap) {
				decrypted, err := svc.DecryptMap(clusterMap)
				if err != nil {
					return count, fmt.Errorf("failed to decrypt cluster for context %s: %w", row.ID, err)
				}
				b, _ := json.Marshal(decrypted)
				newCluster = string(b)
				modified = true
			} else if !decrypt && !svc.IsEncrypted(clusterMap) {
				encrypted, err := svc.EncryptMap(clusterMap)
				if err != nil {
					return count, fmt.Errorf("failed to encrypt cluster for context %s: %w", row.ID, err)
				}
				b, _ := json.Marshal(encrypted)
				newCluster = string(b)
				modified = true
			}
		}

		if modified {
			if !dryRun {
				updates := map[string]interface{}{
					"auth":    newAuth,
					"cluster": newCluster,
				}
				if err := db.Model(&rawK8sContext{}).Where("id = ?", row.ID).Updates(updates).Error; err != nil {
					return count, err
				}
			}
			count++
		}
	}

	return count, nil
}
