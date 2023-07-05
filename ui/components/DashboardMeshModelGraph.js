import Grid from "@material-ui/core/Grid";
import React, { useEffect, useRef, useState } from "react";
import {
  Typography
} from "@material-ui/core";
import bb, { donut,pie } from "billboard.js";
import { getAllComponents,getMeshModels,getRelationshipsDetail,fetchCategories, getModelFromCategoryApi } from "../api/meshmodel"

const BBChart = ({ options }) => {
  const chartRef = useRef(null);

  useEffect(() => {
    const chart = bb.generate({ ...options,bindto : chartRef.current });

    return () => {
      chart.destroy();
    };
  }, [options]);

  return <div ref={chartRef} ></div>;
};

const useFetchTotal = ( fetchr ) => {
  const [total , setTotal] = useState(0)

  useEffect(() => {
    fetchr().then((json) => {
      setTotal(json["total_count"])
    }).catch(e => console.log("Api Error : ",e))
  },[])
  return total
}

const CHART_COLORS = ["#00b39f","#B32700" ,"#396679" , '#EE5351']
const dataToColors = (data) => {
  const columns = data.map((item) => item[0] )
  const colors = { }
  CHART_COLORS.forEach((color,idx) => {
    if (idx < columns.length ) {
      colors[columns[idx]] = color
    }
  })
  return colors
}


function MeshModelContructs({ classes }){
  const totalModels = useFetchTotal(() => getMeshModels(1,1))
  const totalComponents = useFetchTotal(() => getAllComponents(1,1))
  const totalRelations = useFetchTotal(() => getRelationshipsDetail(1 , 1))
  const data = [["Components", totalComponents], ["Models", totalModels], ["Relations",totalRelations]]
  const chartOptions = {
    data : {
      columns : data ,
      type : donut(),
      colors : dataToColors(data)
    },
    tooltip : {
      format : {
        value : function(v){
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

  const [categoryMap , setCategoryMap] = useState({})
  const cleanedData = Object.keys(categoryMap).map((key) => ([key,categoryMap[key]]) )
  const chartOptions = {

    data : {
      columns : cleanedData ,
      colors : dataToColors(cleanedData),
      type : pie(),
    },
    tooltip : {
      format : {
        value : function(v){
          return `${v} Models`
        }
      }
    },
  }
  useEffect(() => {
    fetchCategories().then((categoriesJson) => {
      categoriesJson['categories'].forEach((category) => {
        let categoryName = category.name
        getModelFromCategoryApi(categoryName).then((modelsJson) => {
          setCategoryMap((prevState) => ( { ...prevState,[categoryName] : modelsJson["total_count"] }) )
        })
      })
    })
  },[])

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
    <Grid container spacing={4}>
      <Grid item xs={12} md={6}>
        <MeshModelCategories classes={classes}/>
      </Grid>

      <Grid item xs={12} md={6}>
        <MeshModelContructs classes={classes}/>
      </Grid>

    </Grid>)
}