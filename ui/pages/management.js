import React, {useState} from "react";
import { Grid,TextField, MenuItem,  Chip} from "@mui/material";
import { PaperWithoutTitle, PaperWithTitle } from "@/components/Paper";
import { styled } from "@mui/material/styles";
const {img} = "../../docs/assets/img/adapters/linkerd/linkerd-adapter.png" 
import MesheryMetrics from "@/components/MesheryMetrics" 
import ReactSelectWrapper from "@/components/ReactSelectWrapper"

export default function Lifestyle() {
 
  const CustomChip = styled(Chip)(({ theme }) => ({
    height : "50px",
    fontSize : "15px",
    position : "relative",
    top : theme.spacing(0.5),
    [theme.breakpoints.down("md")] : { fontSize : "12px", },
  }));


  function SelectServiceMesh () {
   
   return (
   <PaperWithoutTitle>

    <TextField
    select
    id="adapter_id"
    name="adapter_name"
    label="Select Service Mesh Type"
    fullWidth
    margin="normal"
    variant="outlined"
  > 
  <MenuItem>
  meshery-linkerd:10001
  </MenuItem>
    </TextField>
    </PaperWithoutTitle>
  )
  }
  
  

  function ManageServiceMesh () {
    return(
      <PaperWithTitle title="Manage Service Mesh" titleVariant="h6">
        <Grid container spacing={4}>
    <Grid container item xs={12} alignItems="flex-start" justify="space-between" style={{ gap : '2rem', }} >
     <CustomChip label="meshery-linkerd:10001" icon={<img src={img} />}  variant="outlined" />
     <div style={{width: "80%"}}>
     <ReactSelectWrapper  label="Namespace"  />
     </div>
     </Grid>
      </Grid>
      </PaperWithTitle>
    
    )
  }
  return (
    <>
    <Grid container spacing={3} >
        <Grid item xs={12}>
          <SelectServiceMesh />
        </Grid>
    <Grid item xs={12}>
    <ManageServiceMesh />
    </Grid>
    <Grid item xs={12}>
    <MesheryMetrics />
    </Grid>
    </Grid>
    </>
  );
}
