package models

// This is the output of this board: Kubernetes / Compute Resources / Cluster, which comes as part of
// prometheus operator install &
// $datasource template variable replaced with "prometheus" &
// $cluster template variable replaced with ""
const staticBoardCluster = `
{
	"annotations": {
	  "list": [
		{
		  "builtIn": 1,
		  "datasource": "-- Grafana --",
		  "enable": true,
		  "hide": true,
		  "iconColor": "rgba(0, 211, 255, 1)",
		  "name": "Annotations & Alerts",
		  "type": "dashboard"
		}
	  ]
	},
	"editable": true,
	"gnetId": null,
	"graphTooltip": 0,
	"id": 2,
	"iteration": 1568661178372,
	"links": [],
	"panels": [
	  {
		"collapsed": false,
		"gridPos": {
		  "h": 1,
		  "w": 24,
		  "x": 0,
		  "y": 0
		},
		"id": 11,
		"panels": [],
		"repeat": null,
		"title": "Headlines",
		"type": "row"
	  },
	  {
		"aliasColors": {},
		"bars": false,
		"cacheTimeout": null,
		"colorBackground": false,
		"colorValue": false,
		"colors": [
		  "#299c46",
		  "rgba(237, 129, 40, 0.89)",
		  "#d44a3a"
		],
		"dashLength": 10,
		"dashes": false,
		"datasource": "prometheus",
		"fill": 1,
		"format": "percentunit",
		"gauge": {
		  "maxValue": 100,
		  "minValue": 0,
		  "show": false,
		  "thresholdLabels": false,
		  "thresholdMarkers": true
		},
		"gridPos": {
		  "h": 3,
		  "w": 4,
		  "x": 0,
		  "y": 1
		},
		"id": 1,
		"interval": null,
		"legend": {
		  "avg": false,
		  "current": false,
		  "max": false,
		  "min": false,
		  "show": true,
		  "total": false,
		  "values": false
		},
		"lines": true,
		"linewidth": 1,
		"links": [],
		"mappingType": 1,
		"mappingTypes": [
		  {
			"name": "value to text",
			"value": 1
		  },
		  {
			"name": "range to text",
			"value": 2
		  }
		],
		"maxDataPoints": 100,
		"nullPointMode": "null as zero",
		"nullText": null,
		"percentage": false,
		"pointradius": 5,
		"points": false,
		"postfix": "",
		"postfixFontSize": "50%",
		"prefix": "",
		"prefixFontSize": "50%",
		"rangeMaps": [
		  {
			"from": "null",
			"text": "N/A",
			"to": "null"
		  }
		],
		"renderer": "flot",
		"seriesOverrides": [],
		"spaceLength": 10,
		"sparkline": {
		  "fillColor": "rgba(31, 118, 189, 0.18)",
		  "full": false,
		  "lineColor": "rgb(31, 120, 193)",
		  "show": false
		},
		"stack": false,
		"steppedLine": false,
		"tableColumn": "",
		"targets": [
		  {
			"expr": "1 - avg(rate(node_cpu_seconds_total{mode=\"idle\", cluster=\"\"}[1m]))",
			"format": "time_series",
			"instant": true,
			"intervalFactor": 2,
			"refId": "A"
		  }
		],
		"thresholds": "70,80",
		"timeFrom": null,
		"timeShift": null,
		"title": "CPU Utilization",
		"tooltip": {
		  "shared": false,
		  "sort": 0,
		  "value_type": "individual"
		},
		"type": "singlestat",
		"valueFontSize": "80%",
		"valueMaps": [
		  {
			"op": "=",
			"text": "N/A",
			"value": "null"
		  }
		],
		"valueName": "avg",
		"xaxis": {
		  "buckets": null,
		  "mode": "time",
		  "name": null,
		  "show": true,
		  "values": []
		},
		"yaxes": [
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": 0,
			"show": true
		  },
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": false
		  }
		]
	  },
	  {
		"aliasColors": {},
		"bars": false,
		"cacheTimeout": null,
		"colorBackground": false,
		"colorValue": false,
		"colors": [
		  "#299c46",
		  "rgba(237, 129, 40, 0.89)",
		  "#d44a3a"
		],
		"dashLength": 10,
		"dashes": false,
		"datasource": "prometheus",
		"fill": 1,
		"format": "percentunit",
		"gauge": {
		  "maxValue": 100,
		  "minValue": 0,
		  "show": false,
		  "thresholdLabels": false,
		  "thresholdMarkers": true
		},
		"gridPos": {
		  "h": 3,
		  "w": 4,
		  "x": 4,
		  "y": 1
		},
		"id": 2,
		"interval": null,
		"legend": {
		  "avg": false,
		  "current": false,
		  "max": false,
		  "min": false,
		  "show": true,
		  "total": false,
		  "values": false
		},
		"lines": true,
		"linewidth": 1,
		"links": [],
		"mappingType": 1,
		"mappingTypes": [
		  {
			"name": "value to text",
			"value": 1
		  },
		  {
			"name": "range to text",
			"value": 2
		  }
		],
		"maxDataPoints": 100,
		"nullPointMode": "null as zero",
		"nullText": null,
		"percentage": false,
		"pointradius": 5,
		"points": false,
		"postfix": "",
		"postfixFontSize": "50%",
		"prefix": "",
		"prefixFontSize": "50%",
		"rangeMaps": [
		  {
			"from": "null",
			"text": "N/A",
			"to": "null"
		  }
		],
		"renderer": "flot",
		"seriesOverrides": [],
		"spaceLength": 10,
		"sparkline": {
		  "fillColor": "rgba(31, 118, 189, 0.18)",
		  "full": false,
		  "lineColor": "rgb(31, 120, 193)",
		  "show": false
		},
		"stack": false,
		"steppedLine": false,
		"tableColumn": "",
		"targets": [
		  {
			"expr": "sum(kube_pod_container_resource_requests_cpu_cores{cluster=\"\"}) / sum(node:node_num_cpu:sum{cluster=\"\"})",
			"format": "time_series",
			"instant": true,
			"intervalFactor": 2,
			"refId": "A"
		  }
		],
		"thresholds": "70,80",
		"timeFrom": null,
		"timeShift": null,
		"title": "CPU Requests Commitment",
		"tooltip": {
		  "shared": false,
		  "sort": 0,
		  "value_type": "individual"
		},
		"type": "singlestat",
		"valueFontSize": "80%",
		"valueMaps": [
		  {
			"op": "=",
			"text": "N/A",
			"value": "null"
		  }
		],
		"valueName": "avg",
		"xaxis": {
		  "buckets": null,
		  "mode": "time",
		  "name": null,
		  "show": true,
		  "values": []
		},
		"yaxes": [
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": 0,
			"show": true
		  },
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": false
		  }
		]
	  },
	  {
		"aliasColors": {},
		"bars": false,
		"cacheTimeout": null,
		"colorBackground": false,
		"colorValue": false,
		"colors": [
		  "#299c46",
		  "rgba(237, 129, 40, 0.89)",
		  "#d44a3a"
		],
		"dashLength": 10,
		"dashes": false,
		"datasource": "prometheus",
		"fill": 1,
		"format": "percentunit",
		"gauge": {
		  "maxValue": 100,
		  "minValue": 0,
		  "show": false,
		  "thresholdLabels": false,
		  "thresholdMarkers": true
		},
		"gridPos": {
		  "h": 3,
		  "w": 4,
		  "x": 8,
		  "y": 1
		},
		"id": 3,
		"interval": null,
		"legend": {
		  "avg": false,
		  "current": false,
		  "max": false,
		  "min": false,
		  "show": true,
		  "total": false,
		  "values": false
		},
		"lines": true,
		"linewidth": 1,
		"links": [],
		"mappingType": 1,
		"mappingTypes": [
		  {
			"name": "value to text",
			"value": 1
		  },
		  {
			"name": "range to text",
			"value": 2
		  }
		],
		"maxDataPoints": 100,
		"nullPointMode": "null as zero",
		"nullText": null,
		"percentage": false,
		"pointradius": 5,
		"points": false,
		"postfix": "",
		"postfixFontSize": "50%",
		"prefix": "",
		"prefixFontSize": "50%",
		"rangeMaps": [
		  {
			"from": "null",
			"text": "N/A",
			"to": "null"
		  }
		],
		"renderer": "flot",
		"seriesOverrides": [],
		"spaceLength": 10,
		"sparkline": {
		  "fillColor": "rgba(31, 118, 189, 0.18)",
		  "full": false,
		  "lineColor": "rgb(31, 120, 193)",
		  "show": false
		},
		"stack": false,
		"steppedLine": false,
		"tableColumn": "",
		"targets": [
		  {
			"expr": "sum(kube_pod_container_resource_limits_cpu_cores{cluster=\"\"}) / sum(node:node_num_cpu:sum{cluster=\"\"})",
			"format": "time_series",
			"instant": true,
			"intervalFactor": 2,
			"refId": "A"
		  }
		],
		"thresholds": "70,80",
		"timeFrom": null,
		"timeShift": null,
		"title": "CPU Limits Commitment",
		"tooltip": {
		  "shared": false,
		  "sort": 0,
		  "value_type": "individual"
		},
		"type": "singlestat",
		"valueFontSize": "80%",
		"valueMaps": [
		  {
			"op": "=",
			"text": "N/A",
			"value": "null"
		  }
		],
		"valueName": "avg",
		"xaxis": {
		  "buckets": null,
		  "mode": "time",
		  "name": null,
		  "show": true,
		  "values": []
		},
		"yaxes": [
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": 0,
			"show": true
		  },
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": false
		  }
		]
	  },
	  {
		"aliasColors": {},
		"bars": false,
		"cacheTimeout": null,
		"colorBackground": false,
		"colorValue": false,
		"colors": [
		  "#299c46",
		  "rgba(237, 129, 40, 0.89)",
		  "#d44a3a"
		],
		"dashLength": 10,
		"dashes": false,
		"datasource": "prometheus",
		"fill": 1,
		"format": "percentunit",
		"gauge": {
		  "maxValue": 100,
		  "minValue": 0,
		  "show": false,
		  "thresholdLabels": false,
		  "thresholdMarkers": true
		},
		"gridPos": {
		  "h": 3,
		  "w": 4,
		  "x": 12,
		  "y": 1
		},
		"id": 4,
		"interval": null,
		"legend": {
		  "avg": false,
		  "current": false,
		  "max": false,
		  "min": false,
		  "show": true,
		  "total": false,
		  "values": false
		},
		"lines": true,
		"linewidth": 1,
		"links": [],
		"mappingType": 1,
		"mappingTypes": [
		  {
			"name": "value to text",
			"value": 1
		  },
		  {
			"name": "range to text",
			"value": 2
		  }
		],
		"maxDataPoints": 100,
		"nullPointMode": "null as zero",
		"nullText": null,
		"percentage": false,
		"pointradius": 5,
		"points": false,
		"postfix": "",
		"postfixFontSize": "50%",
		"prefix": "",
		"prefixFontSize": "50%",
		"rangeMaps": [
		  {
			"from": "null",
			"text": "N/A",
			"to": "null"
		  }
		],
		"renderer": "flot",
		"seriesOverrides": [],
		"spaceLength": 10,
		"sparkline": {
		  "fillColor": "rgba(31, 118, 189, 0.18)",
		  "full": false,
		  "lineColor": "rgb(31, 120, 193)",
		  "show": false
		},
		"stack": false,
		"steppedLine": false,
		"tableColumn": "",
		"targets": [
		  {
			"expr": "1 - sum(:node_memory_MemFreeCachedBuffers_bytes:sum{cluster=\"\"}) / sum(:node_memory_MemTotal_bytes:sum{cluster=\"\"})",
			"format": "time_series",
			"instant": true,
			"intervalFactor": 2,
			"refId": "A"
		  }
		],
		"thresholds": "70,80",
		"timeFrom": null,
		"timeShift": null,
		"title": "Memory Utilization",
		"tooltip": {
		  "shared": false,
		  "sort": 0,
		  "value_type": "individual"
		},
		"type": "singlestat",
		"valueFontSize": "80%",
		"valueMaps": [
		  {
			"op": "=",
			"text": "N/A",
			"value": "null"
		  }
		],
		"valueName": "avg",
		"xaxis": {
		  "buckets": null,
		  "mode": "time",
		  "name": null,
		  "show": true,
		  "values": []
		},
		"yaxes": [
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": 0,
			"show": true
		  },
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": false
		  }
		]
	  },
	  {
		"aliasColors": {},
		"bars": false,
		"cacheTimeout": null,
		"colorBackground": false,
		"colorValue": false,
		"colors": [
		  "#299c46",
		  "rgba(237, 129, 40, 0.89)",
		  "#d44a3a"
		],
		"dashLength": 10,
		"dashes": false,
		"datasource": "prometheus",
		"fill": 1,
		"format": "percentunit",
		"gauge": {
		  "maxValue": 100,
		  "minValue": 0,
		  "show": false,
		  "thresholdLabels": false,
		  "thresholdMarkers": true
		},
		"gridPos": {
		  "h": 3,
		  "w": 4,
		  "x": 16,
		  "y": 1
		},
		"id": 5,
		"interval": null,
		"legend": {
		  "avg": false,
		  "current": false,
		  "max": false,
		  "min": false,
		  "show": true,
		  "total": false,
		  "values": false
		},
		"lines": true,
		"linewidth": 1,
		"links": [],
		"mappingType": 1,
		"mappingTypes": [
		  {
			"name": "value to text",
			"value": 1
		  },
		  {
			"name": "range to text",
			"value": 2
		  }
		],
		"maxDataPoints": 100,
		"nullPointMode": "null as zero",
		"nullText": null,
		"percentage": false,
		"pointradius": 5,
		"points": false,
		"postfix": "",
		"postfixFontSize": "50%",
		"prefix": "",
		"prefixFontSize": "50%",
		"rangeMaps": [
		  {
			"from": "null",
			"text": "N/A",
			"to": "null"
		  }
		],
		"renderer": "flot",
		"seriesOverrides": [],
		"spaceLength": 10,
		"sparkline": {
		  "fillColor": "rgba(31, 118, 189, 0.18)",
		  "full": false,
		  "lineColor": "rgb(31, 120, 193)",
		  "show": false
		},
		"stack": false,
		"steppedLine": false,
		"tableColumn": "",
		"targets": [
		  {
			"expr": "sum(kube_pod_container_resource_requests_memory_bytes{cluster=\"\"}) / sum(:node_memory_MemTotal_bytes:sum{cluster=\"\"})",
			"format": "time_series",
			"instant": true,
			"intervalFactor": 2,
			"refId": "A"
		  }
		],
		"thresholds": "70,80",
		"timeFrom": null,
		"timeShift": null,
		"title": "Memory Requests Commitment",
		"tooltip": {
		  "shared": false,
		  "sort": 0,
		  "value_type": "individual"
		},
		"type": "singlestat",
		"valueFontSize": "80%",
		"valueMaps": [
		  {
			"op": "=",
			"text": "N/A",
			"value": "null"
		  }
		],
		"valueName": "avg",
		"xaxis": {
		  "buckets": null,
		  "mode": "time",
		  "name": null,
		  "show": true,
		  "values": []
		},
		"yaxes": [
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": 0,
			"show": true
		  },
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": false
		  }
		]
	  },
	  {
		"aliasColors": {},
		"bars": false,
		"cacheTimeout": null,
		"colorBackground": false,
		"colorValue": false,
		"colors": [
		  "#299c46",
		  "rgba(237, 129, 40, 0.89)",
		  "#d44a3a"
		],
		"dashLength": 10,
		"dashes": false,
		"datasource": "prometheus",
		"fill": 1,
		"format": "percentunit",
		"gauge": {
		  "maxValue": 100,
		  "minValue": 0,
		  "show": false,
		  "thresholdLabels": false,
		  "thresholdMarkers": true
		},
		"gridPos": {
		  "h": 3,
		  "w": 4,
		  "x": 20,
		  "y": 1
		},
		"id": 6,
		"interval": null,
		"legend": {
		  "avg": false,
		  "current": false,
		  "max": false,
		  "min": false,
		  "show": true,
		  "total": false,
		  "values": false
		},
		"lines": true,
		"linewidth": 1,
		"links": [],
		"mappingType": 1,
		"mappingTypes": [
		  {
			"name": "value to text",
			"value": 1
		  },
		  {
			"name": "range to text",
			"value": 2
		  }
		],
		"maxDataPoints": 100,
		"nullPointMode": "null as zero",
		"nullText": null,
		"percentage": false,
		"pointradius": 5,
		"points": false,
		"postfix": "",
		"postfixFontSize": "50%",
		"prefix": "",
		"prefixFontSize": "50%",
		"rangeMaps": [
		  {
			"from": "null",
			"text": "N/A",
			"to": "null"
		  }
		],
		"renderer": "flot",
		"seriesOverrides": [],
		"spaceLength": 10,
		"sparkline": {
		  "fillColor": "rgba(31, 118, 189, 0.18)",
		  "full": false,
		  "lineColor": "rgb(31, 120, 193)",
		  "show": false
		},
		"stack": false,
		"steppedLine": false,
		"tableColumn": "",
		"targets": [
		  {
			"expr": "sum(kube_pod_container_resource_limits_memory_bytes{cluster=\"\"}) / sum(:node_memory_MemTotal_bytes:sum{cluster=\"\"})",
			"format": "time_series",
			"instant": true,
			"intervalFactor": 2,
			"refId": "A"
		  }
		],
		"thresholds": "70,80",
		"timeFrom": null,
		"timeShift": null,
		"title": "Memory Limits Commitment",
		"tooltip": {
		  "shared": false,
		  "sort": 0,
		  "value_type": "individual"
		},
		"type": "singlestat",
		"valueFontSize": "80%",
		"valueMaps": [
		  {
			"op": "=",
			"text": "N/A",
			"value": "null"
		  }
		],
		"valueName": "avg",
		"xaxis": {
		  "buckets": null,
		  "mode": "time",
		  "name": null,
		  "show": true,
		  "values": []
		},
		"yaxes": [
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": 0,
			"show": true
		  },
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": false
		  }
		]
	  },
	  {
		"collapsed": false,
		"gridPos": {
		  "h": 1,
		  "w": 24,
		  "x": 0,
		  "y": 4
		},
		"id": 12,
		"panels": [],
		"repeat": null,
		"title": "CPU",
		"type": "row"
	  },
	  {
		"aliasColors": {},
		"bars": false,
		"dashLength": 10,
		"dashes": false,
		"datasource": "prometheus",
		"fill": 10,
		"gridPos": {
		  "h": 7,
		  "w": 24,
		  "x": 0,
		  "y": 5
		},
		"id": 7,
		"legend": {
		  "avg": false,
		  "current": false,
		  "max": false,
		  "min": false,
		  "show": true,
		  "total": false,
		  "values": false
		},
		"lines": true,
		"linewidth": 0,
		"links": [],
		"nullPointMode": "null as zero",
		"paceLength": 10,
		"percentage": false,
		"pointradius": 5,
		"points": false,
		"renderer": "flot",
		"seriesOverrides": [],
		"spaceLength": 10,
		"stack": true,
		"steppedLine": false,
		"targets": [
		  {
			"expr": "sum(namespace_pod_name_container_name:container_cpu_usage_seconds_total:sum_rate{cluster=\"\"}) by (namespace)",
			"format": "time_series",
			"intervalFactor": 2,
			"legendFormat": "{{namespace}}",
			"legendLink": null,
			"refId": "A",
			"step": 10
		  }
		],
		"thresholds": [],
		"timeFrom": null,
		"timeRegions": [],
		"timeShift": null,
		"title": "CPU Usage",
		"tooltip": {
		  "shared": false,
		  "sort": 0,
		  "value_type": "individual"
		},
		"type": "graph",
		"xaxis": {
		  "buckets": null,
		  "mode": "time",
		  "name": null,
		  "show": true,
		  "values": []
		},
		"yaxes": [
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": 0,
			"show": true
		  },
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": false
		  }
		],
		"yaxis": {
		  "align": false,
		  "alignLevel": null
		}
	  },
	  {
		"collapsed": false,
		"gridPos": {
		  "h": 1,
		  "w": 24,
		  "x": 0,
		  "y": 12
		},
		"id": 13,
		"panels": [],
		"repeat": null,
		"title": "CPU Quota",
		"type": "row"
	  },
	  {
		"aliasColors": {},
		"bars": false,
		"columns": [],
		"dashLength": 10,
		"dashes": false,
		"datasource": "prometheus",
		"fill": 1,
		"fontSize": "100%",
		"gridPos": {
		  "h": 7,
		  "w": 24,
		  "x": 0,
		  "y": 13
		},
		"id": 8,
		"legend": {
		  "avg": false,
		  "current": false,
		  "max": false,
		  "min": false,
		  "show": true,
		  "total": false,
		  "values": false
		},
		"lines": true,
		"linewidth": 1,
		"links": [],
		"nullPointMode": "null as zero",
		"pageSize": null,
		"percentage": false,
		"pointradius": 5,
		"points": false,
		"renderer": "flot",
		"scroll": true,
		"seriesOverrides": [],
		"showHeader": true,
		"sort": {
		  "col": 0,
		  "desc": true
		},
		"spaceLength": 10,
		"stack": false,
		"steppedLine": false,
		"styles": [
		  {
			"alias": "Time",
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"pattern": "Time",
			"type": "hidden"
		  },
		  {
			"alias": "Pods",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 0,
			"link": true,
			"linkTooltip": "Drill down to pods",
			"linkUrl": "/d/85a562078cdf77779eaa1add43ccec1e/k8s-resources-namespace?var-datasource=prometheus&var-cluster=&var-namespace=$__cell_1",
			"pattern": "Value #A",
			"thresholds": [],
			"type": "number",
			"unit": "short"
		  },
		  {
			"alias": "Workloads",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 0,
			"link": true,
			"linkTooltip": "Drill down to workloads",
			"linkUrl": "/d/a87fb0d919ec0ea5f6543124e16c42a5/k8s-resources-workloads-namespace?var-datasource=prometheus&var-cluster=&var-namespace=$__cell_1",
			"pattern": "Value #B",
			"thresholds": [],
			"type": "number",
			"unit": "short"
		  },
		  {
			"alias": "CPU Usage",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #C",
			"thresholds": [],
			"type": "number",
			"unit": "short"
		  },
		  {
			"alias": "CPU Requests",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #D",
			"thresholds": [],
			"type": "number",
			"unit": "short"
		  },
		  {
			"alias": "CPU Requests %",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #E",
			"thresholds": [],
			"type": "number",
			"unit": "percentunit"
		  },
		  {
			"alias": "CPU Limits",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #F",
			"thresholds": [],
			"type": "number",
			"unit": "short"
		  },
		  {
			"alias": "CPU Limits %",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #G",
			"thresholds": [],
			"type": "number",
			"unit": "percentunit"
		  },
		  {
			"alias": "Namespace",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": true,
			"linkTooltip": "Drill down to pods",
			"linkUrl": "/d/85a562078cdf77779eaa1add43ccec1e/k8s-resources-namespace?var-datasource=prometheus&var-cluster=&var-namespace=$__cell",
			"pattern": "namespace",
			"thresholds": [],
			"type": "number",
			"unit": "short"
		  },
		  {
			"alias": "",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"pattern": "/.*/",
			"thresholds": [],
			"type": "string",
			"unit": "short"
		  }
		],
		"targets": [
		  {
			"expr": "count(mixin_pod_workload{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "A",
			"step": 10
		  },
		  {
			"expr": "count(avg(mixin_pod_workload{cluster=\"\"}) by (workload, namespace)) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "B",
			"step": 10
		  },
		  {
			"expr": "sum(namespace_pod_name_container_name:container_cpu_usage_seconds_total:sum_rate{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "C",
			"step": 10
		  },
		  {
			"expr": "sum(kube_pod_container_resource_requests_cpu_cores{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "D",
			"step": 10
		  },
		  {
			"expr": "sum(namespace_pod_name_container_name:container_cpu_usage_seconds_total:sum_rate{cluster=\"\"}) by (namespace) / sum(kube_pod_container_resource_requests_cpu_cores{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "E",
			"step": 10
		  },
		  {
			"expr": "sum(kube_pod_container_resource_limits_cpu_cores{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "F",
			"step": 10
		  },
		  {
			"expr": "sum(namespace_pod_name_container_name:container_cpu_usage_seconds_total:sum_rate{cluster=\"\"}) by (namespace) / sum(kube_pod_container_resource_limits_cpu_cores{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "G",
			"step": 10
		  }
		],
		"thresholds": [],
		"timeFrom": null,
		"timeShift": null,
		"title": "CPU Quota",
		"tooltip": {
		  "shared": false,
		  "sort": 0,
		  "value_type": "individual"
		},
		"transform": "table",
		"type": "table",
		"xaxis": {
		  "buckets": null,
		  "mode": "time",
		  "name": null,
		  "show": true,
		  "values": []
		},
		"yaxes": [
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": 0,
			"show": true
		  },
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": false
		  }
		]
	  },
	  {
		"collapsed": false,
		"gridPos": {
		  "h": 1,
		  "w": 24,
		  "x": 0,
		  "y": 20
		},
		"id": 14,
		"panels": [],
		"repeat": null,
		"title": "Memory",
		"type": "row"
	  },
	  {
		"aliasColors": {},
		"bars": false,
		"dashLength": 10,
		"dashes": false,
		"datasource": "prometheus",
		"fill": 10,
		"gridPos": {
		  "h": 7,
		  "w": 24,
		  "x": 0,
		  "y": 21
		},
		"id": 9,
		"legend": {
		  "avg": false,
		  "current": false,
		  "max": false,
		  "min": false,
		  "show": true,
		  "total": false,
		  "values": false
		},
		"lines": true,
		"linewidth": 0,
		"links": [],
		"nullPointMode": "null as zero",
		"paceLength": 10,
		"percentage": false,
		"pointradius": 5,
		"points": false,
		"renderer": "flot",
		"seriesOverrides": [],
		"spaceLength": 10,
		"stack": true,
		"steppedLine": false,
		"targets": [
		  {
			"expr": "sum(container_memory_rss{cluster=\"\", container_name!=\"\"}) by (namespace)",
			"format": "time_series",
			"intervalFactor": 2,
			"legendFormat": "{{namespace}}",
			"legendLink": null,
			"refId": "A",
			"step": 10
		  }
		],
		"thresholds": [],
		"timeFrom": null,
		"timeRegions": [],
		"timeShift": null,
		"title": "Memory Usage (w/o cache)",
		"tooltip": {
		  "shared": false,
		  "sort": 0,
		  "value_type": "individual"
		},
		"type": "graph",
		"xaxis": {
		  "buckets": null,
		  "mode": "time",
		  "name": null,
		  "show": true,
		  "values": []
		},
		"yaxes": [
		  {
			"format": "bytes",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": 0,
			"show": true
		  },
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": false
		  }
		],
		"yaxis": {
		  "align": false,
		  "alignLevel": null
		}
	  },
	  {
		"collapsed": false,
		"gridPos": {
		  "h": 1,
		  "w": 24,
		  "x": 0,
		  "y": 28
		},
		"id": 15,
		"panels": [],
		"repeat": null,
		"title": "Memory Requests",
		"type": "row"
	  },
	  {
		"aliasColors": {},
		"bars": false,
		"columns": [],
		"dashLength": 10,
		"dashes": false,
		"datasource": "prometheus",
		"fill": 1,
		"fontSize": "100%",
		"gridPos": {
		  "h": 7,
		  "w": 24,
		  "x": 0,
		  "y": 29
		},
		"id": 10,
		"legend": {
		  "avg": false,
		  "current": false,
		  "max": false,
		  "min": false,
		  "show": true,
		  "total": false,
		  "values": false
		},
		"lines": true,
		"linewidth": 1,
		"links": [],
		"nullPointMode": "null as zero",
		"pageSize": null,
		"percentage": false,
		"pointradius": 5,
		"points": false,
		"renderer": "flot",
		"scroll": true,
		"seriesOverrides": [],
		"showHeader": true,
		"sort": {
		  "col": 0,
		  "desc": true
		},
		"spaceLength": 10,
		"stack": false,
		"steppedLine": false,
		"styles": [
		  {
			"alias": "Time",
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"pattern": "Time",
			"type": "hidden"
		  },
		  {
			"alias": "Pods",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 0,
			"link": true,
			"linkTooltip": "Drill down to pods",
			"linkUrl": "/d/85a562078cdf77779eaa1add43ccec1e/k8s-resources-namespace?var-datasource=prometheus&var-cluster=&var-namespace=$__cell_1",
			"pattern": "Value #A",
			"thresholds": [],
			"type": "number",
			"unit": "short"
		  },
		  {
			"alias": "Workloads",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 0,
			"link": true,
			"linkTooltip": "Drill down to workloads",
			"linkUrl": "/d/a87fb0d919ec0ea5f6543124e16c42a5/k8s-resources-workloads-namespace?var-datasource=prometheus&var-cluster=&var-namespace=$__cell_1",
			"pattern": "Value #B",
			"thresholds": [],
			"type": "number",
			"unit": "short"
		  },
		  {
			"alias": "CPU Usage",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #C",
			"thresholds": [],
			"type": "number",
			"unit": "short"
		  },
		  {
			"alias": "Memory Usage",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #D",
			"thresholds": [],
			"type": "number",
			"unit": "bytes"
		  },
		  {
			"alias": "Memory Requests",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #E",
			"thresholds": [],
			"type": "number",
			"unit": "bytes"
		  },
		  {
			"alias": "Memory Requests %",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #F",
			"thresholds": [],
			"type": "number",
			"unit": "percentunit"
		  },
		  {
			"alias": "Memory Limits",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #G",
			"thresholds": [],
			"type": "number",
			"unit": "bytes"
		  },
		  {
			"alias": "Memory Limits %",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": false,
			"linkTooltip": "Drill down",
			"linkUrl": "",
			"pattern": "Value #H",
			"thresholds": [],
			"type": "number",
			"unit": "percentunit"
		  },
		  {
			"alias": "Namespace",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"link": true,
			"linkTooltip": "Drill down to pods",
			"linkUrl": "/d/85a562078cdf77779eaa1add43ccec1e/k8s-resources-namespace?var-datasource=prometheus&var-cluster=&var-namespace=$__cell",
			"pattern": "namespace",
			"thresholds": [],
			"type": "number",
			"unit": "short"
		  },
		  {
			"alias": "",
			"colorMode": null,
			"colors": [],
			"dateFormat": "YYYY-MM-DD HH:mm:ss",
			"decimals": 2,
			"pattern": "/.*/",
			"thresholds": [],
			"type": "string",
			"unit": "short"
		  }
		],
		"targets": [
		  {
			"expr": "count(mixin_pod_workload{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "A",
			"step": 10
		  },
		  {
			"expr": "count(avg(mixin_pod_workload{cluster=\"\"}) by (workload, namespace)) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "B",
			"step": 10
		  },
		  {
			"expr": "sum(container_memory_rss{cluster=\"\", container_name!=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "C",
			"step": 10
		  },
		  {
			"expr": "sum(kube_pod_container_resource_requests_memory_bytes{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "D",
			"step": 10
		  },
		  {
			"expr": "sum(container_memory_rss{cluster=\"\", container_name!=\"\"}) by (namespace) / sum(kube_pod_container_resource_requests_memory_bytes{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "E",
			"step": 10
		  },
		  {
			"expr": "sum(kube_pod_container_resource_limits_memory_bytes{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "F",
			"step": 10
		  },
		  {
			"expr": "sum(container_memory_rss{cluster=\"\", container_name!=\"\"}) by (namespace) / sum(kube_pod_container_resource_limits_memory_bytes{cluster=\"\"}) by (namespace)",
			"format": "table",
			"instant": true,
			"intervalFactor": 2,
			"legendFormat": "",
			"refId": "G",
			"step": 10
		  }
		],
		"thresholds": [],
		"timeFrom": null,
		"timeShift": null,
		"title": "Requests by Namespace",
		"tooltip": {
		  "shared": false,
		  "sort": 0,
		  "value_type": "individual"
		},
		"transform": "table",
		"type": "table",
		"xaxis": {
		  "buckets": null,
		  "mode": "time",
		  "name": null,
		  "show": true,
		  "values": []
		},
		"yaxes": [
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": 0,
			"show": true
		  },
		  {
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": false
		  }
		]
	  }
	],
	"refresh": "10s",
	"schemaVersion": 18,
	"style": "dark",
	"tags": [
	  "kubernetes-mixin"
	],
	"templating": {
	  "list": [
		{
		  "current": {
			"text": "prometheus",
			"value": "prometheus"
		  },
		  "hide": 0,
		  "label": null,
		  "name": "datasource",
		  "options": [],
		  "query": "prometheus",
		  "refresh": 1,
		  "regex": "",
		  "skipUrlSync": false,
		  "type": "datasource"
		},
		{
		  "allValue": null,
		  "current": {
			"isNone": true,
			"text": "None",
			"value": ""
		  },
		  "datasource": "prometheus",
		  "definition": "",
		  "hide": 2,
		  "includeAll": false,
		  "label": "cluster",
		  "multi": false,
		  "name": "cluster",
		  "options": [],
		  "query": "label_values(node_cpu_seconds_total, cluster)",
		  "refresh": 1,
		  "regex": "",
		  "skipUrlSync": false,
		  "sort": 2,
		  "tagValuesQuery": "",
		  "tags": [],
		  "tagsQuery": "",
		  "type": "query",
		  "useTags": false
		}
	  ]
	},
	"time": {
	  "from": "now-1h",
	  "to": "now"
	},
	"timepicker": {
	  "refresh_intervals": [
		"5s",
		"10s",
		"30s",
		"1m",
		"5m",
		"15m",
		"30m",
		"1h",
		"2h",
		"1d"
	  ],
	  "time_options": [
		"5m",
		"15m",
		"1h",
		"6h",
		"12h",
		"24h",
		"2d",
		"7d",
		"30d"
	  ]
	},
	"timezone": "",
	"title": "Kubernetes / Compute Resources / Cluster",
	"uid": "efa86fd1d0c121a26444b636a3f509a8",
	"version": 1
  }
`
