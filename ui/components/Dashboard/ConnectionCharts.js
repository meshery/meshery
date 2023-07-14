import React  from "react"
import {
  Typography
} from "@material-ui/core";
import { donut } from "billboard.js";
import BBChart from "../BBChart"

export default function ConnectionStatsChart({ classes }) {
  const data = [
    ["Connected", 10],
    ["Discovered", 24],
    ["Ignored",23],
    ["Registered",15],
  ]
  const chartOptions = {
    data : {
      columns : data,
      type : donut(),
      colors : {
        Connected : "#006B56",   // Light Green
        Discovered : "#EBC017",  // Saffron
        Ignored : "#51636B",     // Lighter Dark
        Registered : "#477E96",  // Teal Blue
      },
    },
    arc : {
      cornerRadius : {
        ratio : 0.05
      }
    },
    donut : {
      title : "Connection Stats",
      padAngle : 0.03,
      label : {
        format : function(value) {
          return value
        }
      }
    },
    tooltip : {
      format : {
        value : function (v) {
          return v
        }
      }
    },
  }

  return (
    <div className={classes.dashboardSection}>
      <Typography variant="h6" gutterBottom className={classes.chartTitle}>
        Connections
      </Typography>

      <div>
        <BBChart options={chartOptions} />
      </div>
    </div>
  )

}