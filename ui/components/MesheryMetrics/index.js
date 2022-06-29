import React from "react";
import { styled } from "@mui/material/styles";
import { PaperWithoutTitle } from "@/components/Paper";
import { MetricsButton } from "@/components/Button"
import AddIcon from '@mui/icons-material/Add';
import {Typography} from "@mui/material";

function MesheryMetrics(){

    const CustomPaperWithoutTitle = styled(PaperWithoutTitle)(() => ({
        padding : "2rem",
         display : "flex",
          justifyContent : "center", 
          alignItems : "center",
        flexDirection : "column"
    
    }))
    const CustomTypography = styled(Typography)(() => ({
        fontSize : "1.5rem",
        marginBottom : "2rem"
    
    }))

 return(
   <CustomPaperWithoutTitle>
    <div       style={{
        padding : "2rem",
        display : "flex",
        justifyContent : "center",
        alignItems : "center",
        flexDirection : "column",
        margin: "auto",
      }}>
    <CustomTypography align="center" color="textSecondary"> 
    No Service Mesh Metrics Configurations Found
    </CustomTypography>
     <MetricsButton variant="contained" style={{width: "auto"}} startIcon={<AddIcon />} size="large" >
        Configure Service Mesh Metrics
     </MetricsButton>
     </div>
   </CustomPaperWithoutTitle>
 )
}


export default MesheryMetrics;
