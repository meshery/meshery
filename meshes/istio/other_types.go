package istio

type supportedOperation struct {
	// a unique identifier
	key string
	// a friendly name
	name string
	// the template file name
	templateName string

	resetOp bool
}

var supportedOps = []supportedOperation{
	{
		key:          "istio_1",
		name:         "Shift All traffic to version V1 of all the services",
		templateName: "shift_all_traffic_to_v1_of_all_services.tmpl",
	},
	{
		key:          "istio_2",
		name:         "Shift logged in user traffic to version V2 of reviews service",
		templateName: "shift_user_traffic_to_v2_of_reviews.tmpl",
	},
	{
		key:          "istio_3",
		name:         "Inject a HTTP abort for ratings service for the logged in user",
		templateName: "inject_abort_for_ratings_service_for_user.tmpl",
	},
	{
		key:          "istio_4",
		name:         "Inject a HTTP delay for ratings service for the logged in user",
		templateName: "inject_delay_for_ratings_service_for_user.tmpl",
	},
	{
		key:          "istio_5",
		name:         "Shift 50 percent of traffic to version V3 of reviews service",
		templateName: "shift_50_percent_of_traffic_to_v3_of_reviews.tmpl",
	},
	{
		key:          "istio_6",
		name:         "Shift 100 percent of traffic to version V3 of reviews service",
		templateName: "shift_all_traffic_to_v3_of_reviews.tmpl",
	},
	{
		key:     "istio_7",
		name:    "Reset all applied rules",
		resetOp: true,
	},
}
