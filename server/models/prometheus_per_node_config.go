package models

// This is the output of this board: Kubernetes / Compute Resources / Nodes, which comes as part of
// prometheus operator install &
// $datasource template variable replaced with "prometheus" &
// $cluster template variable replaced with ""
const staticBoardNodes = `
[[ $indexCheck := .indexCheck ]]
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
"editable": false,
"gnetId": null,
"graphTooltip": 0,
"id": 7,
"iteration": 1568396170452,
"links": [],
"panels": [
 [[ range $ind, $instance := .instances ]]
  {
	"aliasColors": {},
	"bars": false,
	"dashLength": 10,
	"dashes": false,
	"datasource": "prometheus",
	"fill": 1,
	"gridPos": {
	  "h": 7,
	  "w": 12,
	  "x": 0,
	  "y": 0
	},
	"id": 2,
	"legend": {
	  "alignAsTable": false,
	  "avg": false,
	  "current": false,
	  "max": false,
	  "min": false,
	  "rightSide": false,
	  "show": true,
	  "total": false,
	  "values": false
	},
	"lines": true,
	"linewidth": 1,
	"links": [],
	"nullPointMode": "null",
	"paceLength": 10,
	"percentage": false,
	"pointradius": 5,
	"points": false,
	"renderer": "flot",
	"repeat": null,
	"seriesOverrides": [],
	"spaceLength": 10,
	"stack": false,
	"steppedLine": false,
	"targets": [
	  {
		"expr": "max(node_load1{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"})",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "load 1m",
		"refId": "A"
	  },
	  {
		"expr": "max(node_load5{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"})",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "load 5m",
		"refId": "B"
	  },
	  {
		"expr": "max(node_load15{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"})",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "load 15m",
		"refId": "C"
	  },
	  {
		"expr": "count(node_cpu_seconds_total{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\", mode=\"user\"})",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "logical cores",
		"refId": "D"
	  }
	],
	"thresholds": [],
	"timeFrom": null,
	"timeRegions": [],
	"timeShift": null,
	"title": "System load - [[ $instance ]]",
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
		"min": null,
		"show": true
	  },
	  {
		"format": "short",
		"label": null,
		"logBase": 1,
		"max": null,
		"min": null,
		"show": true
	  }
	],
	"yaxis": {
	  "align": false,
	  "alignLevel": null
	}
  },
  {
	"aliasColors": {},
	"bars": false,
	"dashLength": 10,
	"dashes": false,
	"datasource": "prometheus",
	"fill": 1,
	"gridPos": {
	  "h": 7,
	  "w": 12,
	  "x": 12,
	  "y": 0
	},
	"id": 3,
	"legend": {
	  "alignAsTable": false,
	  "avg": false,
	  "current": false,
	  "max": false,
	  "min": false,
	  "rightSide": false,
	  "show": true,
	  "total": false,
	  "values": false
	},
	"lines": true,
	"linewidth": 1,
	"links": [],
	"nullPointMode": "null",
	"paceLength": 10,
	"percentage": false,
	"pointradius": 5,
	"points": false,
	"renderer": "flot",
	"repeat": null,
	"seriesOverrides": [],
	"spaceLength": 10,
	"stack": false,
	"steppedLine": false,
	"targets": [
	  {
		"expr": "sum by (cpu) (irate(node_cpu_seconds_total{cluster=\"\", job=\"node-exporter\", mode!=\"idle\", instance=\"[[ $instance ]]\"}[5m]))",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "{{cpu}}",
		"refId": "A"
	  }
	],
	"thresholds": [],
	"timeFrom": null,
	"timeRegions": [],
	"timeShift": null,
	"title": "Usage Per Core - [[ $instance ]]",
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
		"format": "percentunit",
		"label": null,
		"logBase": 1,
		"max": null,
		"min": null,
		"show": true
	  },
	  {
		"format": "percentunit",
		"label": null,
		"logBase": 1,
		"max": null,
		"min": null,
		"show": true
	  }
	],
	"yaxis": {
	  "align": false,
	  "alignLevel": null
	}
  },
  {
	"aliasColors": {},
	"bars": false,
	"dashLength": 10,
	"dashes": false,
	"datasource": "prometheus",
	"fill": 1,
	"gridPos": {
	  "h": 7,
	  "w": 18,
	  "x": 0,
	  "y": 7
	},
	"id": 4,
	"legend": {
	  "alignAsTable": true,
	  "avg": true,
	  "current": true,
	  "max": false,
	  "min": false,
	  "rightSide": true,
	  "show": true,
	  "total": false,
	  "values": true
	},
	"lines": true,
	"linewidth": 1,
	"links": [],
	"nullPointMode": "null",
	"paceLength": 10,
	"percentage": false,
	"pointradius": 5,
	"points": false,
	"renderer": "flot",
	"repeat": null,
	"seriesOverrides": [],
	"spaceLength": 10,
	"stack": false,
	"steppedLine": false,
	"targets": [
	  {
		"expr": "max (sum by (cpu) (irate(node_cpu_seconds_total{cluster=\"\", job=\"node-exporter\", mode!=\"idle\", instance=\"[[ $instance ]]\"}[2m])) ) * 100",
		"format": "time_series",
		"intervalFactor": 10,
		"legendFormat": "{{ cpu }}",
		"refId": "A"
	  }
	],
	"thresholds": [],
	"timeFrom": null,
	"timeRegions": [],
	"timeShift": null,
	"title": "CPU Utilization - [[ $instance ]]",
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
		"format": "percent",
		"label": null,
		"logBase": 1,
		"max": 100,
		"min": 0,
		"show": true
	  },
	  {
		"format": "percent",
		"label": null,
		"logBase": 1,
		"max": 100,
		"min": 0,
		"show": true
	  }
	],
	"yaxis": {
	  "align": false,
	  "alignLevel": null
	}
  },
  {
	"cacheTimeout": null,
	"colorBackground": false,
	"colorValue": false,
	"colors": [
	  "rgba(50, 172, 45, 0.97)",
	  "rgba(237, 129, 40, 0.89)",
	  "rgba(245, 54, 54, 0.9)"
	],
	"datasource": "prometheus",
	"format": "percent",
	"gauge": {
	  "maxValue": 100,
	  "minValue": 0,
	  "show": true,
	  "thresholdLabels": false,
	  "thresholdMarkers": true
	},
	"gridPos": {
	  "h": 7,
	  "w": 6,
	  "x": 18,
	  "y": 7
	},
	"id": 5,
	"interval": null,
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
	"nullPointMode": "connected",
	"nullText": null,
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
	"sparkline": {
	  "fillColor": "rgba(31, 118, 189, 0.18)",
	  "full": false,
	  "lineColor": "rgb(31, 120, 193)",
	  "show": false
	},
	"tableColumn": "",
	"targets": [
	  {
		"expr": "avg(sum by (cpu) (irate(node_cpu_seconds_total{cluster=\"\", job=\"node-exporter\", mode!=\"idle\", instance=\"[[ $instance ]]\"}[2m]))) * 100",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "",
		"refId": "A"
	  }
	],
	"thresholds": "80, 90",
	"title": "CPU Usage - [[ $instance ]]",
	"tooltip": {
	  "shared": false
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
	"valueName": "current"
  },
  {
	"aliasColors": {},
	"bars": false,
	"dashLength": 10,
	"dashes": false,
	"datasource": "prometheus",
	"fill": 1,
	"gridPos": {
	  "h": 7,
	  "w": 18,
	  "x": 0,
	  "y": 14
	},
	"id": 6,
	"legend": {
	  "alignAsTable": false,
	  "avg": false,
	  "current": false,
	  "max": false,
	  "min": false,
	  "rightSide": false,
	  "show": true,
	  "total": false,
	  "values": false
	},
	"lines": true,
	"linewidth": 1,
	"links": [],
	"nullPointMode": "null",
	"paceLength": 10,
	"percentage": false,
	"pointradius": 5,
	"points": false,
	"renderer": "flot",
	"repeat": null,
	"seriesOverrides": [],
	"spaceLength": 10,
	"stack": false,
	"steppedLine": false,
	"targets": [
	  {
		"expr": "max(node_memory_MemTotal_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}  - node_memory_MemFree_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}  - node_memory_Buffers_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}  - node_memory_Cached_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"})",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "memory used",
		"refId": "A"
	  },
	  {
		"expr": "max(node_memory_Buffers_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"})",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "memory buffers",
		"refId": "B"
	  },
	  {
		"expr": "max(node_memory_Cached_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"})",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "memory cached",
		"refId": "C"
	  },
	  {
		"expr": "max(node_memory_MemFree_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"})",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "memory free",
		"refId": "D"
	  }
	],
	"thresholds": [],
	"timeFrom": null,
	"timeRegions": [],
	"timeShift": null,
	"title": "Memory Usage - [[ $instance ]]",
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
		"min": null,
		"show": true
	  },
	  {
		"format": "bytes",
		"label": null,
		"logBase": 1,
		"max": null,
		"min": null,
		"show": true
	  }
	],
	"yaxis": {
	  "align": false,
	  "alignLevel": null
	}
  },
  {
	"cacheTimeout": null,
	"colorBackground": false,
	"colorValue": false,
	"colors": [
	  "rgba(50, 172, 45, 0.97)",
	  "rgba(237, 129, 40, 0.89)",
	  "rgba(245, 54, 54, 0.9)"
	],
	"datasource": "prometheus",
	"format": "percent",
	"gauge": {
	  "maxValue": 100,
	  "minValue": 0,
	  "show": true,
	  "thresholdLabels": false,
	  "thresholdMarkers": true
	},
	"gridPos": {
	  "h": 7,
	  "w": 6,
	  "x": 18,
	  "y": 14
	},
	"id": 7,
	"interval": null,
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
	"nullPointMode": "connected",
	"nullText": null,
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
	"sparkline": {
	  "fillColor": "rgba(31, 118, 189, 0.18)",
	  "full": false,
	  "lineColor": "rgb(31, 120, 193)",
	  "show": false
	},
	"tableColumn": "",
	"targets": [
	  {
		"expr": "max(  (    (      node_memory_MemTotal_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}    - node_memory_MemFree_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}    - node_memory_Buffers_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}    - node_memory_Cached_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}    )    / node_memory_MemTotal_bytes{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}  ) * 100)",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "",
		"refId": "A"
	  }
	],
	"thresholds": "80, 90",
	"title": "Memory Usage - [[ $instance ]]",
	"tooltip": {
	  "shared": false
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
	"valueName": "current"
  },
  {
	"aliasColors": {},
	"bars": false,
	"dashLength": 10,
	"dashes": false,
	"datasource": "prometheus",
	"fill": 1,
	"gridPos": {
	  "h": 7,
	  "w": 12,
	  "x": 0,
	  "y": 21
	},
	"id": 8,
	"legend": {
	  "alignAsTable": false,
	  "avg": false,
	  "current": false,
	  "max": false,
	  "min": false,
	  "rightSide": false,
	  "show": true,
	  "total": false,
	  "values": false
	},
	"lines": true,
	"linewidth": 1,
	"links": [],
	"nullPointMode": "null",
	"paceLength": 10,
	"percentage": false,
	"pointradius": 5,
	"points": false,
	"renderer": "flot",
	"repeat": null,
	"seriesOverrides": [
	  {
		"alias": "read",
		"yaxis": 1
	  },
	  {
		"alias": "io time",
		"yaxis": 2
	  }
	],
	"spaceLength": 10,
	"stack": false,
	"steppedLine": false,
	"targets": [
	  {
		"expr": "max(rate(node_disk_read_bytes_total{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}[2m]))",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "read",
		"refId": "A"
	  },
	  {
		"expr": "max(rate(node_disk_written_bytes_total{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}[2m]))",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "written",
		"refId": "B"
	  },
	  {
		"expr": "max(rate(node_disk_io_time_seconds_total{cluster=\"\", job=\"node-exporter\",  instance=\"[[ $instance ]]\"}[2m]))",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "io time",
		"refId": "C"
	  }
	],
	"thresholds": [],
	"timeFrom": null,
	"timeRegions": [],
	"timeShift": null,
	"title": "Disk I/O - [[ $instance ]]",
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
		"min": null,
		"show": true
	  },
	  {
		"format": "ms",
		"label": null,
		"logBase": 1,
		"max": null,
		"min": null,
		"show": true
	  }
	],
	"yaxis": {
	  "align": false,
	  "alignLevel": null
	}
  },
  {
	"aliasColors": {},
	"bars": false,
	"dashLength": 10,
	"dashes": false,
	"datasource": "prometheus",
	"fill": 1,
	"gridPos": {
	  "h": 7,
	  "w": 12,
	  "x": 12,
	  "y": 21
	},
	"id": 9,
	"legend": {
	  "alignAsTable": false,
	  "avg": false,
	  "current": false,
	  "max": false,
	  "min": false,
	  "rightSide": false,
	  "show": true,
	  "total": false,
	  "values": false
	},
	"lines": true,
	"linewidth": 1,
	"links": [],
	"nullPointMode": "null",
	"paceLength": 10,
	"percentage": false,
	"pointradius": 5,
	"points": false,
	"renderer": "flot",
	"repeat": null,
	"seriesOverrides": [],
	"spaceLength": 10,
	"stack": false,
	"steppedLine": false,
	"targets": [
	  {
		"expr": "node:node_filesystem_usage:{cluster=\"\"}",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "{{device}}",
		"refId": "A"
	  }
	],
	"thresholds": [],
	"timeFrom": null,
	"timeRegions": [],
	"timeShift": null,
	"title": "Disk Space Usage - [[ $instance ]]",
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
		"format": "percentunit",
		"label": null,
		"logBase": 1,
		"max": null,
		"min": null,
		"show": true
	  },
	  {
		"format": "percentunit",
		"label": null,
		"logBase": 1,
		"max": null,
		"min": null,
		"show": true
	  }
	],
	"yaxis": {
	  "align": false,
	  "alignLevel": null
	}
  },
  {
	"aliasColors": {},
	"bars": false,
	"dashLength": 10,
	"dashes": false,
	"datasource": "prometheus",
	"fill": 1,
	"gridPos": {
	  "h": 7,
	  "w": 12,
	  "x": 0,
	  "y": 28
	},
	"id": 10,
	"legend": {
	  "alignAsTable": false,
	  "avg": false,
	  "current": false,
	  "max": false,
	  "min": false,
	  "rightSide": false,
	  "show": true,
	  "total": false,
	  "values": false
	},
	"lines": true,
	"linewidth": 1,
	"links": [],
	"nullPointMode": "null",
	"paceLength": 10,
	"percentage": false,
	"pointradius": 5,
	"points": false,
	"renderer": "flot",
	"repeat": null,
	"seriesOverrides": [],
	"spaceLength": 10,
	"stack": false,
	"steppedLine": false,
	"targets": [
	  {
		"expr": "max(rate(node_network_receive_bytes_total{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\", device!~\"lo\"}[5m]))",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "{{device}}",
		"refId": "A"
	  }
	],
	"thresholds": [],
	"timeFrom": null,
	"timeRegions": [],
	"timeShift": null,
	"title": "Network Received - [[ $instance ]]",
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
		"min": null,
		"show": true
	  },
	  {
		"format": "bytes",
		"label": null,
		"logBase": 1,
		"max": null,
		"min": null,
		"show": true
	  }
	],
	"yaxis": {
	  "align": false,
	  "alignLevel": null
	}
  },
  {
	"aliasColors": {},
	"bars": false,
	"dashLength": 10,
	"dashes": false,
	"datasource": "prometheus",
	"fill": 1,
	"gridPos": {
	  "h": 7,
	  "w": 12,
	  "x": 12,
	  "y": 28
	},
	"id": 11,
	"legend": {
	  "alignAsTable": false,
	  "avg": false,
	  "current": false,
	  "max": false,
	  "min": false,
	  "rightSide": false,
	  "show": true,
	  "total": false,
	  "values": false
	},
	"lines": true,
	"linewidth": 1,
	"links": [],
	"nullPointMode": "null",
	"paceLength": 10,
	"percentage": false,
	"pointradius": 5,
	"points": false,
	"renderer": "flot",
	"repeat": null,
	"seriesOverrides": [],
	"spaceLength": 10,
	"stack": false,
	"steppedLine": false,
	"targets": [
	  {
		"expr": "max(rate(node_network_transmit_bytes_total{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\", device!~\"lo\"}[5m]))",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "{{device}}",
		"refId": "A"
	  }
	],
	"thresholds": [],
	"timeFrom": null,
	"timeRegions": [],
	"timeShift": null,
	"title": "Network Transmitted - [[ $instance ]]",
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
		"min": null,
		"show": true
	  },
	  {
		"format": "bytes",
		"label": null,
		"logBase": 1,
		"max": null,
		"min": null,
		"show": true
	  }
	],
	"yaxis": {
	  "align": false,
	  "alignLevel": null
	}
  },
  {
	"aliasColors": {},
	"bars": false,
	"dashLength": 10,
	"dashes": false,
	"datasource": "prometheus",
	"fill": 1,
	"gridPos": {
	  "h": 7,
	  "w": 18,
	  "x": 0,
	  "y": 35
	},
	"id": 12,
	"legend": {
	  "alignAsTable": false,
	  "avg": false,
	  "current": false,
	  "max": false,
	  "min": false,
	  "rightSide": false,
	  "show": true,
	  "total": false,
	  "values": false
	},
	"lines": true,
	"linewidth": 1,
	"links": [],
	"nullPointMode": "null",
	"paceLength": 10,
	"percentage": false,
	"pointradius": 5,
	"points": false,
	"renderer": "flot",
	"repeat": null,
	"seriesOverrides": [],
	"spaceLength": 10,
	"stack": false,
	"steppedLine": false,
	"targets": [
	  {
		"expr": "max(  node_filesystem_files{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}  - node_filesystem_files_free{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"})",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "inodes used",
		"refId": "A"
	  },
	  {
		"expr": "max(node_filesystem_files_free{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"})",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "inodes free",
		"refId": "B"
	  }
	],
	"thresholds": [],
	"timeFrom": null,
	"timeRegions": [],
	"timeShift": null,
	"title": "Inodes Usage - [[ $instance ]]",
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
		"min": null,
		"show": true
	  },
	  {
		"format": "short",
		"label": null,
		"logBase": 1,
		"max": null,
		"min": null,
		"show": true
	  }
	],
	"yaxis": {
	  "align": false,
	  "alignLevel": null
	}
  },
  {
	"cacheTimeout": null,
	"colorBackground": false,
	"colorValue": false,
	"colors": [
	  "rgba(50, 172, 45, 0.97)",
	  "rgba(237, 129, 40, 0.89)",
	  "rgba(245, 54, 54, 0.9)"
	],
	"datasource": "prometheus",
	"format": "percent",
	"gauge": {
	  "maxValue": 100,
	  "minValue": 0,
	  "show": true,
	  "thresholdLabels": false,
	  "thresholdMarkers": true
	},
	"gridPos": {
	  "h": 7,
	  "w": 6,
	  "x": 18,
	  "y": 35
	},
	"id": 13,
	"interval": null,
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
	"nullPointMode": "connected",
	"nullText": null,
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
	"sparkline": {
	  "fillColor": "rgba(31, 118, 189, 0.18)",
	  "full": false,
	  "lineColor": "rgb(31, 120, 193)",
	  "show": false
	},
	"tableColumn": "",
	"targets": [
	  {
		"expr": "max(  (    (      node_filesystem_files{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}    - node_filesystem_files_free{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}    )    / node_filesystem_files{cluster=\"\", job=\"node-exporter\", instance=\"[[ $instance ]]\"}  ) * 100)",
		"format": "time_series",
		"intervalFactor": 2,
		"legendFormat": "",
		"refId": "A"
	  }
	],
	"thresholds": "80, 90",
	"title": "Inodes Usage - [[ $instance ]]",
	"tooltip": {
	  "shared": false
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
	"valueName": "current"
  }[[if ne $indexCheck $ind ]],
  [[ end ]]
  [[ end ]]
],
"refresh": "",
"schemaVersion": 18,
"style": "dark",
"tags": [
  "kubernetes-mixin"
],
"templating": {
  "list": [
	{
	  "current": {
		"selected": true,
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
		"selected": true,
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
	  "query": "label_values(kube_pod_info, cluster)",
	  "refresh": 2,
	  "regex": "",
	  "skipUrlSync": false,
	  "sort": 0,
	  "tagValuesQuery": "",
	  "tags": [],
	  "tagsQuery": "",
	  "type": "query",
	  "useTags": false
	},
	{
	  "allValue": null,
	  "current": {
		"selected": false,
		"tags": [],
		"text": "10.199.75.57:9100",
		"value": "10.199.75.57:9100"
	  },
	  "datasource": "prometheus",
	  "definition": "",
	  "hide": 0,
	  "includeAll": false,
	  "label": null,
	  "multi": false,
	  "name": "instance",
	  "options": [],
	  "query": "label_values(node_boot_time_seconds{cluster=\"\", job=\"node-exporter\"}, instance)",
	  "refresh": 2,
	  "regex": "",
	  "skipUrlSync": false,
	  "sort": 0,
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
"title": "Kubernetes / Nodes",
"uid": "fa49a4706d07a042595b664c87fb33ea",
"version": 1
}`
