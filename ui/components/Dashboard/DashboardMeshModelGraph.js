import Grid from "@material-ui/core/Grid";
import React, { useEffect, useState } from "react";
import {
  Typography
} from "@material-ui/core";
import BBChart from "../BBChart"
import { donut, pie } from "billboard.js";
import { getAllComponents, getMeshModels, getRelationshipsDetail, fetchCategories, getModelFromCategoryApi } from "../../api/meshmodel"
import { dataToColors } from "../../utils/chartColors"

const useFetchTotal = (fetchr) => {
  const [total, setTotal] = useState(0)

  useEffect(() => {
    fetchr().then((json) => {
      setTotal(json["total_count"])
    }).catch(e => console.log("Api Error : ", e))
  }, [])
  return total
}

function MeshModelContructs({ classes }) {
  const totalModels = useFetchTotal(() => getMeshModels(1, 1))
  const totalComponents = useFetchTotal(() => getAllComponents(1, 1))
  const totalRelationships = useFetchTotal(() => getRelationshipsDetail(1, 1))
  // TODO: Add Policies
  const data = [
    ["Models", totalModels],
    ["Components", totalComponents],
    ["Relationships", totalRelationships]
    // TODO: Add Policies
  ]
  const chartOptions = {
    data : {
      columns : data,
      type : donut(),
      colors : dataToColors(data),
    },
    arc : {
      cornerRadius : {
        ratio : 0.1
      }
    },
    donut : {
      title : "Capabilities\nRegistry"
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
        MeshModel Constructs
      </Typography>

      <div>
        <BBChart options={chartOptions} />
      </div>
    </div>
  )
}



function MeshModelCategories({ classes }) {

  const [categoryMap, setCategoryMap] = useState({})
  const cleanedData = Object.keys(categoryMap).map((key) => ([key, categoryMap[key]]))
  const chartOptions = {

    data : {
      columns : cleanedData,
      colors : dataToColors(cleanedData),
      type : pie(),
    },
    tooltip : {
      format : {
        value : function (v) {
          return `${v} Models`
        }
      }
    },
    legend : {
      show : false
    }
  }
  useEffect(() => {
    fetchCategories().then((categoriesJson) => {
      categoriesJson['categories'].forEach((category) => {
        let categoryName = category.name
        getModelFromCategoryApi(categoryName).then((modelsJson) => {
          setCategoryMap((prevState) => ({ ...prevState, [categoryName] : modelsJson["total_count"] }))
        })
      })
    })
  }, [])

  return (
    <div className={classes.dashboardSection}>
      <Typography variant="h6" gutterBottom className={classes.chartTitle}>
        MeshModel Categories
      </Typography>

      <div>
        <BBChart options={chartOptions} />
      </div>
    </div>
  )
}

export default function DashboardMeshModelGraph({ classes }) {

  return (
    <Grid container spacing={2}>
      <Grid item xs={12} md={6}>
        <MeshModelCategories classes={classes} />
      </Grid>

      <Grid item xs={12} md={6}>
        <MeshModelContructs classes={classes} />
      </Grid>

    </Grid>)
}