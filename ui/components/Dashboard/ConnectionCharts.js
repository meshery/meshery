import React , { useEffect,useState } from "react";
import { Typography } from "@material-ui/core";
import { donut } from "billboard.js";
import BBChart from "../BBChart";
import { dataToColors , isValidColumnName } from "../../utils/charts"
import { getConnectionStatusSummary } from "../../api/connections";


export default function ConnectionStatsChart({ classes }) {

  const [chartData,setChartData] = useState([])

  useEffect(() => {
    getConnectionStatusSummary().then(json => {
      setChartData(json.connections_status
        .filter(data => isValidColumnName(data.status))
        .map(data => [data.status,data.count]))
    })
  },[])

  const chartOptions = {
    data : {
      columns : chartData,
      type : donut(),
      colors : dataToColors(chartData),
    },
    arc : {
      cornerRadius : {
        ratio : 0.05,
      },
    },
    donut : {
      title : "Connection Stats",
      padAngle : 0.03,
      label : {
        format : function (value) {
          return value;
        },
      },
    },
    tooltip : {
      format : {
        value : function (v) {
          return v;
        },
      },
    },
  };

  return (
    <div className={classes.dashboardSection}>
      <Typography variant="h6" gutterBottom className={classes.chartTitle}>
        Connections
      </Typography>

      <div>
        <BBChart options={chartOptions} />
      </div>
    </div>
  );
}