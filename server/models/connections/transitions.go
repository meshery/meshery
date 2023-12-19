package connections

var PossibleTransitionnsMap = map[string]map[ConnectionStatus][]ConnectionStatus{
	"kubernetes": {
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
	"meshery": {
		CONNECTED: {
			DELETED,
		},
	},
}
