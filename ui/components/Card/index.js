import React from "react";
import { styled } from "@mui/material/styles";
import { Card, CardHeader, CardActions, IconButton,} from "@mui/material";
// import AddIcon from "@material-ui/icons/Add";
// import DeleteIcon from "@material-ui/icons/Delete";
// import PlayIcon from "@material-ui/icons/PlayArrow";


const CustomCard = styled(Card)(({ theme }) => ({
    height : '100%',
    display : 'flex',
    flexDirection : 'column'
  }));


  export function LifestyleCard  ( cat )  {
    
    if (typeof cat === "undefined") {
      cat = 0;
    }

    let content;
    let description;
    switch (cat) {
      case 0:
        content = "Manage Service Mesh Lifecycle";
        description = "Deploy a service mesh or SMI adapter on your cluster.";
        break;

      case 1:
        content = "Manage Sample Application Lifecycle";
        description = "Deploy sample applications on/off the service mesh.";
        break;

      case 2:
        content = "Apply Service Mesh Configuration";
        description = "Configure your service mesh using some pre-defined options.";
        // selectedAdapterOps = selectedAdapterOps.filter((ops) => !ops.value.startsWith("Add-on:"));
        break;

      case 3:
        content = "Validate Service Mesh Configuration";
        description = "Validate your service mesh configuration against best practices.";
        break;

      case 4:
        content = "Apply Custom Configuration";
        description = "Customize the configuration of your service mesh.";
        break;
    }

    return(
        <CustomCard>
        <CardHeader title={content} subheader={description} style={{ flexGrow : 1 }} />

        </CustomCard>
    )

}  