package academy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
)

func TestAcademyCreate(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	utils.Log = utils.SetupMeshkitLogger("mesheryctl", false, os.Stdout)
	mesheryctlflags.InitValidators(AcademyCmd)

	resetFlags := func(t *testing.T) {
		t.Helper()
		defaults := map[string]string{
			"type": "", "title": "", "description": "", "into": "",
			"org": "", "level": "", "category": "", "tags": "",
			"id": "", "force": "false", "banner": "", "draft": "false",
		}
		for name, value := range defaults {
			if err := createCmd.Flags().Set(name, value); err != nil {
				t.Fatalf("failed to reset flag %q: %v", name, err)
			}
		}
	}

	tests := []struct {
		name         string
		args         []string
		expectErr    bool
		expectedCode string
	}{
		{
			name:      "missing type flag",
			args:      []string{"create", "--title", "Test Title", "--description", "Desc"},
			expectErr: true,
		},
		{
			name:      "missing description flag",
			args:      []string{"create", "--type", "learning-path", "--title", "Test", "--into", tempDir},
			expectErr: true,
		},
		{
			name:         "invalid type flag",
			args:         []string{"create", "--type", "invalid-type", "--title", "Test", "--description", "Desc", "--into", tempDir},
			expectErr:    true,
			expectedCode: ErrTaxonomyTypeCode,
		},
		{
			name:         "fresh scaffold learning path without org",
			args:         []string{"create", "--type", "learning-path", "--title", "My Path", "--description", "Desc"},
			expectErr:    true,
			expectedCode: ErrMissingOrgIDCode,
		},
		{
			name:      "fresh scaffold learning path with org",
			args:      []string{"create", "--type", "learning-path", "--title", "My Path", "--description", "Desc", "--org", "test-org-uuid"},
			expectErr: false,
		},
		{
			name:      "fresh scaffold learning path into explicit dir",
			args:      []string{"create", "--type", "learning-path", "--title", "My Explicit Path", "--description", "Desc", "--into", tempDir, "--category", "cloud", "--tags", "kubernetes, infra", "--org", "explicit-org", "--level", "advanced"},
			expectErr: false,
		},
		{
			name:      "setup for collision test",
			args:      []string{"create", "--type", "learning-path", "--title", "Collision Path", "--description", "Desc", "--into", tempDir, "--org", "collision-org"},
			expectErr: false,
		},
		{
			name:         "existing file collision without force",
			args:         []string{"create", "--type", "learning-path", "--title", "Collision Path", "--description", "Desc", "--into", tempDir, "--org", "collision-org"},
			expectErr:    true,
			expectedCode: ErrScaffoldExistsCode,
		},
		{
			name:      "existing file collision with force",
			args:      []string{"create", "--type", "learning-path", "--title", "Collision Path", "--description", "Desc", "--into", tempDir, "--force", "--org", "collision-org"},
			expectErr: false,
		},
		{
			name:         "invalid nesting error",
			args:         []string{"create", "module", "Bad Module", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path")},
			expectErr:    true,
			expectedCode: ErrInvalidNestingCode,
		},

		{
			name:      "scaffold single node (course) with correct nesting",
			args:      []string{"create", "course", "New Course", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path")},
			expectErr: false,
		},
		{
			name:      "scaffold single node (module) with correct nesting",
			args:      []string{"create", "module", "New Module", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "new-course")},
			expectErr: false,
		},
		{
			name:         "scaffold single node (course) without into should fail",
			args:         []string{"create", "course", "Missing Into", "--description", "Desc"},
			expectErr:    true,
			expectedCode: ErrMissingIntoCode,
		},
		{
			name:      "scaffold flat test under course",
			args:      []string{"create", "test", "Course Test", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "new-course")},
			expectErr: false,
		},
		{
			name:      "scaffold exam under course (no collision with test)",
			args:      []string{"create", "exam", "Course Exam", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "new-course")},
			expectErr: false,
		},
		{
			name:      "scaffold lab under module",
			args:      []string{"create", "lab", "Lab 1", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "course-1", "module-1")},
			expectErr: false,
		},
		{
			name:      "scaffold flat test under module",
			args:      []string{"create", "test", "Module Test", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "course-1", "module-1")},
			expectErr: false,
		},
		{
			name:      "scaffold certification",
			args:      []string{"create", "--type", "certification", "--title", "Cert 1", "--description", "Desc", "--into", tempDir, "--org", "cert-org"},
			expectErr: false,
		},
		{
			name:      "scaffold flat test under certification 1",
			args:      []string{"create", "test", "Cert Test 1", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "certifications", "cert-org", "cert-1")},
			expectErr: false,
		},
		{
			name:      "scaffold flat test under certification 2",
			args:      []string{"create", "test", "Cert Test 2", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "certifications", "cert-org", "cert-1")},
			expectErr: false,
		},
		{
			name:      "scaffold exam under certification",
			args:      []string{"create", "exam", "Final Exam", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "certifications", "cert-org", "cert-1")},
			expectErr: false,
		},
		{
			name:         "fresh scaffold challenge without org",
			args:         []string{"create", "--type", "challenge", "--title", "My Challenge", "--description", "Desc"},
			expectErr:    true,
			expectedCode: ErrMissingOrgIDCode,
		},
		{
			name:      "fresh scaffold challenge with org",
			args:      []string{"create", "--type", "challenge", "--title", "My Challenge", "--description", "Desc", "--org", "test-org-uuid"},
			expectErr: false,
		},
		{
			name:         "re-run challenge scaffold without force",
			args:         []string{"create", "--type", "challenge", "--title", "My Challenge", "--description", "Desc", "--org", "test-org-uuid"},
			expectErr:    true,
			expectedCode: ErrScaffoldExistsCode,
		},
		{
			name:      "re-run challenge scaffold with force",
			args:      []string{"create", "--type", "challenge", "--title", "My Challenge", "--description", "Desc", "--org", "test-org-uuid", "--force"},
			expectErr: false,
		},
		{
			name:      "reject orgID with path traversal",
			args:      []string{"create", "--type", "learning-path", "--title", "My Path", "--description", "Desc", "--org", "../invalid-org"},
			expectErr: true,
		},
		{
			name:         "scaffold page under lab should fail",
			args:         []string{"create", "page", "Bad Page", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "course-1", "module-1", "lab-1")},
			expectErr:    true,
			expectedCode: ErrInvalidNestingCode,
		},
		{
			name:         "scaffold test under test should fail",
			args:         []string{"create", "test", "Bad Test", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "certifications", "cert-org", "cert-1", "test-1")},
			expectErr:    true,
			expectedCode: ErrInvalidNestingCode,
		},
		{
			name:      "inherit level from parent",
			args:      []string{"create", "page", "Inherited Level Page", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "explicit-org", "my-explicit-path", "course-1", "module-1")},
			expectErr: false,
		},
		{
			name:      "reject invalid title slug",
			args:      []string{"create", "course", "../../", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path")},
			expectErr: true,
		},
		{
			name:      "reject reserved Windows device name title",
			args:      []string{"create", "course", "CON", "--description", "Desc", "--into", filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path")},
			expectErr: true,
		},
		{
			name:      "escape special characters in title and ID",
			args:      []string{"create", "--type", "learning-path", "--title", `Title with "quotes" and \backslashes`, "--description", `Desc`, "--id", `"my-id"`, "--org", "escape-org", "--into", tempDir, "--level", "beginner"},
			expectErr: false,
		},
		{
			name:      "scaffold learning path with banner and draft",
			args:      []string{"create", "--type", "learning-path", "--title", "Banner Path", "--description", "Desc", "--org", "banner-org", "--into", tempDir, "--banner", "hero.svg", "--draft"},
			expectErr: false,
		},
		{
			name:         "reject invalid level",
			args:         []string{"create", "--type", "learning-path", "--title", "Invalid Level", "--description", "Desc", "--org", "level-org", "--level", "invalid-level"},
			expectErr:    true,
			expectedCode: ErrInvalidLevelCode,
		},
		{
			name:         "structural node into empty dir (no parent _index.md) should fail",
			args:         []string{"create", "course", "Orphan Course", "--description", "Desc", "--into", tempDir},
			expectErr:    true,
			expectedCode: ErrInvalidParentMetadataCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags(t)

			AcademyCmd.SetArgs(tt.args)
			err := AcademyCmd.Execute()
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if tt.expectErr && tt.expectedCode != "" && err != nil {
				code := errors.GetCode(err)
				if code != tt.expectedCode {
					t.Fatalf("expected error code %q, got %q", tt.expectedCode, code)
				}
			}
		})
	}

	// Additional verification for frontmatter contents
	// 1. Check path construction for fresh scaffold
	lpPath := filepath.Join(tempDir, "content", "learning-paths", "test-org-uuid", "my-path", "_index.md")
	content, err := os.ReadFile(lpPath)
	if err != nil {
		t.Fatalf("Failed to read scaffolded learning-path _index.md: %v", err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "id: \"REPLACE_WITH_INSTRUCTOR_CONSOLE_ID\"") {
		t.Errorf("learning-path should contain id placeholder")
	}
	if !strings.Contains(contentStr, "type: \"learning-path\"") { // singular type
		t.Errorf("learning-path should have singular type")
	}
	if !strings.Contains(contentStr, "level: \"beginner\"") {
		t.Errorf("learning-path should default to beginner level")
	}

	// 2. Check path construction for --into explicit scaffold
	lp2Path := filepath.Join(tempDir, "content", "learning-paths", "explicit-org", "my-explicit-path", "_index.md")
	content2, err := os.ReadFile(lp2Path)
	if err != nil {
		t.Fatalf("Failed to read scaffolded learning-path explicit _index.md: %v", err)
	}
	contentStr2 := string(content2)

	if !strings.Contains(contentStr2, "categories: \"cloud\"") {
		t.Errorf("learning-path explicit should contain singular category")
	}
	if !strings.Contains(contentStr2, `- "kubernetes"`) || !strings.Contains(contentStr2, `- "infra"`) {
		t.Errorf("learning-path explicit should contain tags as list with quotes, got: %v", contentStr2)
	}

	// 3. Check course
	coursePath := filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "new-course", "_index.md")
	content3, err := os.ReadFile(coursePath)
	if err != nil {
		t.Fatalf("Failed to read scaffolded course _index.md: %v", err)
	}
	contentStr3 := string(content3)

	if strings.Contains(contentStr3, "id: \"REPLACE_WITH_INSTRUCTOR_CONSOLE_ID\"") {
		t.Errorf("course should NOT contain id placeholder")
	}
	if !strings.Contains(contentStr3, "tags: []") {
		t.Errorf("course without tags should have empty list")
	}
	if strings.Contains(contentStr3, "categories:") {
		t.Errorf("course without category should omit categories line")
	}
	if strings.Contains(contentStr3, "orgId:") {
		t.Errorf("course should NOT contain orgId")
	}
	// Course/module scaffolded via --into must be single nodes
	for _, stub := range []string{"module-1", "page-1"} {
		if _, err := os.Stat(filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "new-course", stub)); !os.IsNotExist(err) {
			t.Errorf("course scaffolded via --into should not auto-generate a %s stub beneath it", stub)
		}
	}

	standaloneModulePath := filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "new-course", "new-module", "_index.md")
	if _, err := os.Stat(standaloneModulePath); err != nil {
		t.Fatalf("Failed to verify module scaffolded via --into: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "new-course", "new-module", "page-1")); !os.IsNotExist(err) {
		t.Errorf("module scaffolded via --into should not auto-generate a page-1 stub beneath it")
	}

	// 4. Check flat test and exam under course
	courseTestPath := filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "new-course", "test.md")
	_, err = os.Stat(courseTestPath)
	if err != nil {
		t.Fatalf("Failed to verify flat test under course: %v", err)
	}

	courseExamPath := filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "new-course", "course-exam.md")
	_, err = os.Stat(courseExamPath)
	if err != nil {
		t.Fatalf("Failed to verify exam under course: %v", err)
	}

	// 5. Check flat test under module (test.md)
	moduleTestPath := filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "course-1", "module-1", "test.md")
	_, err = os.Stat(moduleTestPath)
	if err != nil {
		t.Fatalf("Failed to verify flat test under module: %v", err)
	}

	// 6. Check lab under module
	labPath := filepath.Join(tempDir, "content", "learning-paths", "collision-org", "collision-path", "course-1", "module-1", "lab-1", "_index.md")
	_, err = os.Stat(labPath)
	if err != nil {
		t.Fatalf("Failed to verify lab under module: %v", err)
	}

	// 7. Check standard directory tests under certification
	certTest1Path := filepath.Join(tempDir, "content", "certifications", "cert-org", "cert-1", "test-1", "_index.md")
	_, err = os.Stat(certTest1Path)
	if err != nil {
		t.Fatalf("Failed to verify standard directory test-1 under certification: %v", err)
	}

	certTest2Path := filepath.Join(tempDir, "content", "certifications", "cert-org", "cert-1", "test-2", "_index.md")
	_, err = os.Stat(certTest2Path)
	if err != nil {
		t.Fatalf("Failed to verify standard directory test-2 under certification: %v", err)
	}

	// 8. Check exam under certification
	examPath := filepath.Join(tempDir, "content", "certifications", "cert-org", "cert-1", "final-exam", "_index.md")
	_, err = os.Stat(examPath)
	if err != nil {
		t.Fatalf("Failed to verify exam under certification: %v", err)
	}

	// 9. Check escaping of special characters
	escapedPath := filepath.Join(tempDir, "content", "learning-paths", "escape-org", "title-with-quotes-and-backslashes", "_index.md")
	escapedContent, err := os.ReadFile(escapedPath)
	if err != nil {
		t.Fatalf("Failed to read scaffolded escaped learning-path _index.md: %v", err)
	}
	escapedStr := string(escapedContent)

	if !strings.Contains(escapedStr, `title: "Title with \"quotes\" and \\backslashes"`) {
		t.Errorf("title was not properly escaped: %v", escapedStr)
	}
	if !strings.Contains(escapedStr, `id: "\"my-id\""`) {
		t.Errorf("id was not properly escaped: %v", escapedStr)
	}

	// 10. Check level inheritance
	inheritedPagePath := filepath.Join(tempDir, "content", "learning-paths", "explicit-org", "my-explicit-path", "course-1", "module-1", "inherited-level-page", "_index.md")
	inheritedContent, err := os.ReadFile(inheritedPagePath)
	if err != nil {
		t.Fatalf("Failed to read scaffolded inherited level page _index.md: %v", err)
	}
	if !strings.Contains(string(inheritedContent), `level: "advanced"`) {
		t.Errorf("expected inherited level 'advanced' in page frontmatter, got: %v", string(inheritedContent))
	}

	// 10b. Check banner and draft frontmatter
	bannerPath := filepath.Join(tempDir, "content", "learning-paths", "banner-org", "banner-path", "_index.md")
	bannerContent, err := os.ReadFile(bannerPath)
	if err != nil {
		t.Fatalf("Failed to read scaffolded banner learning-path _index.md: %v", err)
	}
	bannerStr := string(bannerContent)
	if !strings.Contains(bannerStr, `banner: "hero.svg"`) {
		t.Errorf("learning-path should contain banner frontmatter, got: %v", bannerStr)
	}
	if !strings.Contains(bannerStr, "draft: true") {
		t.Errorf("learning-path should contain draft: true, got: %v", bannerStr)
	}
	if !strings.Contains(bannerStr, `id: "REPLACE_WITH_INSTRUCTOR_CONSOLE_ID"`) {
		t.Errorf("learning-path (root type) should contain id placeholder, got: %v", bannerStr)
	}

	// 11. Check challenge structure
	challengePath := filepath.Join(tempDir, "content", "challenges", "test-org-uuid", "my-challenge", "_index.md")
	challengeContent, err := os.ReadFile(challengePath)
	if err != nil {
		t.Fatalf("Failed to read scaffolded challenge _index.md: %v", err)
	}
	challengeStr := string(challengeContent)
	if !strings.Contains(challengeStr, `type: "challenge"`) {
		t.Errorf("challenge should have type: 'challenge'")
	}
	if !strings.Contains(challengeStr, `id: "REPLACE_WITH_INSTRUCTOR_CONSOLE_ID"`) {
		t.Errorf("challenge should contain id placeholder")
	}

	challengeLabPath := filepath.Join(tempDir, "content", "challenges", "test-org-uuid", "my-challenge", "lab", "_index.md")
	challengeLabContent, err := os.ReadFile(challengeLabPath)
	if err != nil {
		t.Fatalf("Failed to read challenge lab _index.md: %v", err)
	}
	if !strings.Contains(string(challengeLabContent), `type: "lab"`) {
		t.Errorf("challenge lab should have type: 'lab'")
	}

	challengeExamPath := filepath.Join(tempDir, "content", "challenges", "test-org-uuid", "my-challenge", "exam", "_index.md")
	challengeExamContent, err := os.ReadFile(challengeExamPath)
	if err != nil {
		t.Fatalf("Failed to read challenge exam _index.md: %v", err)
	}
	if !strings.Contains(string(challengeExamContent), `type: "exam"`) {
		t.Errorf("challenge exam should have type: 'exam'")
	}

	expectedDirs := []string{
		filepath.Join(tempDir, "content", "challenges", "test-org-uuid", "my-challenge", "content", "description"),
		filepath.Join(tempDir, "content", "challenges", "test-org-uuid", "my-challenge", "content", "getting-started"),
		filepath.Join(tempDir, "content", "challenges", "test-org-uuid", "my-challenge", "content", "faq"),
	}
	for _, dir := range expectedDirs {
		if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
			t.Errorf("expected directory %s to exist", dir)
		}
	}
}
