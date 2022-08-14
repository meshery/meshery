import React from 'react'
import {Grid} from "@mui/material"
import {
  fortioResultToJsChartData, makeChart, makeOverlayChart, makeMultiChart,
}  from "@/lib/chartjs-formatter"

function NonRecursiveConstructDisplayCells(data) {
  // console.log(data)
  return Object.keys(data).map(el => {
    if (typeof data[el].display?.value === "string" && !data[el].display?.hide) {
      return (
        <>
          <b>{data[el].display?.key}</b>: {data[el].display?.value}
        </>
      )
    }
  })
}

function MesheryChart({rawdata, data}) {
  let chartData;
  
 if (typeof data !== 'undefined') {
    const results = data;
    console.log(results.length)
    if (results.length == 1) {
      chartData = makeOverlayChart(fortioResultToJsChartData(rawdata,results[0]));
      console.log(chartData?.options?.metadata)
    }
    if (results.length == 2) {
      chartData = makeOverlayChart(fortioResultToJsChartData(rawdata,results[0]), fortioResultToJsChartData(rawdata,results[1]));
    } else if (results.length > 2) {
      chartData = makeMultiChart(rawdata,results);
    }
  }

  
  return (
    <div>
      
      <Grid container spacing={1} style={{ margin : "1rem" }} justifyContent="center">
            {NonRecursiveConstructDisplayCells(chartData?.options?.metadata )?.map((el, i) => {
              return <Grid item xs={4} key={`nri-${i}`}>{el}  </Grid>
            })}
          </Grid>
    </div>
  )
}

export default MesheryChart