[1mdiff --git a/cmd/main.go b/cmd/main.go[m
[1mindex aa5e1f60..dd07e694 100644[m
[1m--- a/cmd/main.go[m
[1m+++ b/cmd/main.go[m
[36m@@ -185,10 +185,7 @@[m [mfunc main() {[m
 		GenericPersister:                dbHandler,[m
 	}[m
 	lProv.Initialize()[m
[31m-	seededUUIDs, err := lProv.SeedContent(log)[m
[31m-	if err != nil {[m
[31m-		logrus.Error(err)[m
[31m-	}[m
[32m+[m	[32mseededUUIDs := lProv.SeedContent(log)[m
 	provs[lProv.Name()] = lProv[m
 [m
 	cPreferencePersister, err := models.NewBitCaskPreferencePersister(viper.GetString("USER_DATA_FOLDER"))[m
[1mdiff --git a/models/default_local_provider.go b/models/default_local_provider.go[m
[1mindex db866ab6..3e36db35 100644[m
[1m--- a/models/default_local_provider.go[m
[1m+++ b/models/default_local_provider.go[m
[36m@@ -841,49 +841,74 @@[m [mfunc (l *DefaultLocalProvider) GetKubeClient() *mesherykube.Client {[m
 	return l.KubeClient[m
 }[m
 [m
[31m-func (l *DefaultLocalProvider) SeedContent(log logger.Handler) ([]uuid.UUID, error) {[m
[32m+[m[32mfunc (l *DefaultLocalProvider) SeedContent(log logger.Handler) []uuid.UUID {[m
 	var seededUUIDs []uuid.UUID[m
[31m-	log.Info("Starting to seed patterns")[m
 	names, content, err := getSeededComponents("Pattern", log)[m
 	if err != nil {[m
[31m-		return nil, ErrGettingSeededComponents(err, "Patterns")[m
[31m-	}[m
[31m-[m
[31m-	for i, name := range names {[m
[31m-		id, _ := uuid.NewV4()[m
[31m-		var pattern = &MesheryPattern{[m
[31m-			PatternFile: content[i],[m
[31m-			Name:        name,[m
[31m-			ID:          &id,[m
[31m-		}[m
[31m-		log.Info("[SEEDING] ", "Saving pattern- ", name)[m
[31m-		_, err := l.MesheryPatternPersister.SaveMesheryPattern(pattern)[m
[31m-		if err != nil {[m
[31m-			return nil, ErrSavingSeededComponents(err, "Patterns")[m
[32m+[m		[32mlogrus.Error(ErrGettingSeededComponents(err, "Patterns"))[m
[32m+[m	[32m} else {[m
[32m+[m		[32mlog.Info("Starting to seed patterns")[m
[32m+[m		[32mfor i, name := range names {[m
[32m+[m			[32mid, _ := uuid.NewV4()[m
[32m+[m			[32mvar pattern = &MesheryPattern{[m
[32m+[m				[32mPatternFile: content[i],[m
[32m+[m				[32mName:        name,[m
[32m+[m				[32mID:          &id,[m
[32m+[m			[32m}[m
[32m+[m			[32mlog.Info("[SEEDING] ", "Saving pattern- ", name)[m
[32m+[m			[32m_, err := l.MesheryPatternPersister.SaveMesheryPattern(pattern)[m
[32m+[m			[32mif err != nil {[m
[32m+[m				[32mlogrus.Error(ErrGettingSeededComponents(err, "Patterns"))[m
[32m+[m			[32m}[m
[32m+[m			[32mseededUUIDs = append(seededUUIDs, id)[m
 		}[m
[31m-		seededUUIDs = append(seededUUIDs, id)[m
 	}[m
[31m-	log.Info("Starting to seed filters")[m
[32m+[m
 	names, content, err = getSeededComponents("Filter", log)[m
 	if err != nil {[m
[31m-		return nil, ErrGettingSeededComponents(err, "Filters")[m
[32m+[m		[32mlogrus.Error(ErrGettingSeededComponents(err, "Filters"))[m
[32m+[m	[32m} else {[m
[32m+[m		[32mlog.Info("Starting to seed filters")[m
[32m+[m		[32mfor i, name := range names {[m
[32m+[m			[32mid, _ := uuid.NewV4()[m
[32m+[m			[32mvar filter = &MesheryFilter{[m
[32m+[m				[32mFilterFile: content[i],[m
[32m+[m				[32mName:       name,[m
[32m+[m				[32mID:         &id,[m
[32m+[m			[32m}[m
[32m+[m			[32mlog.Info("[SEEDING] ", "Saving filter- ", name)[m
[32m+[m			[32m_, err := l.MesheryFilterPersister.SaveMesheryFilter(filter)[m
[32m+[m			[32mif err != nil {[m
[32m+[m				[32mlogrus.Error(ErrGettingSeededComponents(err, "Filters"))[m
[32m+[m			[32m}[m
[32m+[m			[32mseededUUIDs = append(seededUUIDs, id)[m
[32m+[m		[32m}[m
[32m+[m
 	}[m
 [m
[31m-	for i, name := range names {[m
[31m-		id, _ := uuid.NewV4()[m
[31m-		var filter = &MesheryFilter{[m
[31m-			FilterFile: content[i],[m
[31m-			Name:       name,[m
[31m-			ID:         &id,[m
[31m-		}[m
[31m-		log.Info("[SEEDING] ", "Saving filter- ", name)[m
[31m-		_, err := l.MesheryFilterPersister.SaveMesheryFilter(filter)[m
[31m-		if err != nil {[m
[31m-			return nil, ErrSavingSeededComponents(err, "Filters")[m
[32m+[m	[32mnames, content, err = getSeededComponents("Application", log)[m
[32m+[m	[32mif err != nil {[m
[32m+[m		[32mlogrus.Error(ErrGettingSeededComponents(err, "Applications"))[m
[32m+[m	[32m} else {[m
[32m+[m		[32mlog.Info("Starting to seed applications")[m
[32m+[m
[32m+[m		[32mfor i, name := range names {[m
[32m+[m			[32mid, _ := uuid.NewV4()[m
[32m+[m			[32mvar app = &MesheryApplication{[m
[32m+[m				[32mApplicationFile: content[i],[m
[32m+[m				[32mName:            name,[m
[32m+[m				[32mID:              &id,[m
[32m+[m			[32m}[m
[32m+[m			[32mlog.Info("[SEEDING] ", "Saving application- ", name)[m
[32m+[m			[32m_, err := l.MesheryApplicationPersister.SaveMesheryApplication(app)[m
[32m+[m			[32mif err != nil {[m
[32m+[m				[32mlogrus.Error(ErrGettingSeededComponents(err, "Applications"))[m
[32m+[m			[32m}[m
[32m+[m			[32mseededUUIDs = append(seededUUIDs, id)[m
 		}[m
[31m-		seededUUIDs = append(seededUUIDs, id)[m
 	}[m
[31m-	return seededUUIDs, nil[m
[32m+[m
[32m+[m	[32mreturn seededUUIDs[m
 }[m
 func (l *DefaultLocalProvider) CleanupSeeded(seededUUIDs []uuid.UUID) {[m
 	for _, id := range seededUUIDs {[m
[36m@@ -1167,7 +1192,10 @@[m [mfunc getSeededComponents(comp string, log logger.Handler) ([]string, []string, e[m
 		wd = filepath.Join(wd, ".meshery", "seed_content", "patterns")[m
 	case "Filter":[m
 		wd = filepath.Join(wd, ".meshery", "seed_content", "filters", "binaries")[m
[32m+[m	[32mcase "Application":[m
[32m+[m		[32mwd = filepath.Join(wd, ".meshery", "seed_content", "applications")[m
 	}[m
[32m+[m
 	log.Info("[SEEDING] ", "Extracting "+comp+"s from ", wd)[m
 	var names []string[m
 	var contents []string[m
