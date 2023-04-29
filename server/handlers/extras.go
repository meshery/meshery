package handlers

// swagger:route GET /api/system/meshsync/grafana SystemAPI idMeshSyncGrafana
// Handle GET request for mesh-sync grafana
//
// Fetches Prometheus and Grafana
// responses:
// 	200: v1ServicesMapResponseWrapper

// ScanPromGrafanaHandler - fetches  Prometheus and Grafana
// func (h *Handler) ScanPromGrafanaHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
// 	errs := []string{}
// 	var wg sync.WaitGroup
// 	customK8scontexts, ok := req.Context().Value(models.KubeClustersKey).([]models.K8sContext)
// 	if ok && len(customK8scontexts) > 0 {
// 		for _, mk8scontext := range customK8scontexts {
// 			wg.Add(1)
// 			go func(mk8scontext models.K8sContext) {
// 				defer wg.Done()
// 				k8sconfig, err := mk8scontext.GenerateKubeConfig()
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				availablePromGrafana, err := helpers.ScanPromGrafana(k8sconfig, mk8scontext.Name)
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				if err = json.NewEncoder(w).Encode(availablePromGrafana); err != nil {
// 					obj := "payloads"
// 					h.log.Error(ErrMarshal(err, obj))
// 					errs = append(errs, ErrMarshal(err, obj).Error())
// 					return
// 				}
// 			}(mk8scontext)
// 		}
// 	}
// 	if len(errs) != 0 {
// 		http.Error(w, mergeMsgs(errs), http.StatusInternalServerError)
// 	}
// 	wg.Wait()
// }

// swagger:route GET /api/system/meshsync/prometheus SystemAPI idMeshSyncPrometheus
// Handle GET request for fetching prometheus
//
// Fetches Prometheus
// responses:
// 	200: v1ServicesMapResponseWrapper

// ScanPrometheusHandler - fetches  Prometheus
// func (h *Handler) ScanPrometheusHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
// 	errs := []string{}
// 	var wg sync.WaitGroup
// 	customK8scontexts, ok := req.Context().Value(models.KubeClustersKey).([]models.K8sContext)
// 	if ok && len(customK8scontexts) > 0 {
// 		for _, mk8scontext := range customK8scontexts {
// 			wg.Add(1)
// 			go func(mk8scontext models.K8sContext) {
// 				defer wg.Done()
// 				k8sconfig, err := mk8scontext.GenerateKubeConfig()
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				availablePromGrafana, err := helpers.ScanPrometheus(k8sconfig, mk8scontext.Name)
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				if err = json.NewEncoder(w).Encode(availablePromGrafana); err != nil {
// 					obj := "payloads"
// 					h.log.Error(ErrMarshal(err, obj))
// 					errs = append(errs, ErrMarshal(err, obj).Error())
// 					return
// 				}
// 			}(mk8scontext)
// 		}
// 	}
// 	if len(errs) != 0 {
// 		http.Error(w, mergeMsgs(errs), http.StatusInternalServerError)
// 	}
// 	wg.Wait()
// }

// swagger:route GET /api/telemetry/metrics/grafana/scan GrafanaAPI idGetGrafana
// Handle GET request for Grafana
//
// Fetches and returns Grafana
// responses:
// 	200: v1ServicesMapResponseWrapper

// ScanGrafanaHandler - fetches  Grafana
// func (h *Handler) ScanGrafanaHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
// 	errs := []string{}
// 	var wg sync.WaitGroup
// 	customK8scontexts, ok := req.Context().Value(models.KubeClustersKey).([]models.K8sContext)
// 	if ok && len(customK8scontexts) > 0 {
// 		for _, mk8scontext := range customK8scontexts {
// 			wg.Add(1)
// 			go func(mk8scontext models.K8sContext) {
// 				defer wg.Done()
// 				k8sconfig, err := mk8scontext.GenerateKubeConfig()
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				availablePromGrafana, err := helpers.ScanGrafana(k8sconfig, mk8scontext.Name)
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				if err = json.NewEncoder(w).Encode(availablePromGrafana); err != nil {
// 					obj := "payloads"
// 					h.log.Error(ErrMarshal(err, obj))
// 					errs = append(errs, ErrMarshal(err, obj).Error())
// 					return
// 				}
// 			}(mk8scontext)
// 		}
// 	}
// 	if len(errs) != 0 {
// 		http.Error(w, mergeMsgs(errs), http.StatusInternalServerError)
// 	}
// 	wg.Wait()
// }

// swagger:route GET /api/telemetry/metrics/config PrometheusAPI idGetPrometheusConfig
// Handle GET for Prometheus configuration
//
// Used for fetching Prometheus configuration
// responses:
//  200: prometheusConfigResponseWrapper

// swagger:route POST /api/telemetry/metrics/config PrometheusAPI idPostPrometheusConfig
// Handle POST for Prometheus configuration
//
// Used for persisting Prometheus configuration
// responses:
//  200:

///

// func writeDefK8sOnFileSystem(def string, path string) {
// 	err := ioutil.WriteFile(path, []byte(def), 0777)
// 	if err != nil {
// 		fmt.Println("err def: ", err.Error())
// 	}
// }

// func writeSchemaK8sFileSystem(schema string, path string) {
// 	err := ioutil.WriteFile(path, []byte(schema), 0777)
// 	if err != nil {
// 		fmt.Println("err schema: ", err.Error())
// 	}
// }
