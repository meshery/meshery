package connections

var PossibleTransitionnsMap = map[string]map[ConnectionStatus][]ConnectionStatus{
	"Kubernetes": {
		DISCOVERED: {},
		REGISTERED: {
			CONNECTED, IGNORED,
		},
		CONNECTED: {
			DISCONNECTED, MAINTENANCE, DELETED, NOTFOUND,
		},
		IGNORED: {
		 	DELETED, NOTFOUND,
		},
		MAINTENANCE: {
			REGISTERED, CONNECTED, IGNORED, NOTFOUND,
		},
		DISCONNECTED: {
			CONNECTED, DELETED,
		},
		DELETED: {},
		NOTFOUND: {
			DISCOVERED,
		},
	},
}