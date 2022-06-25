import React from 'react';
import CreatableSelect from 'react-select/creatable';
import { styled } from "@mui/material/styles";
import { MenuItem,TextField,Typography, Chip} from "@mui/material";


const CustomCreatableSelect = styled(CreatableSelect)(({ theme }) => ({
// marginLeft: "4rem",
// width: "70%",
verticalAlign: "center",
margin: "auto 0",
minWidth : '250px',
flex : '1',
}));


function ReactSelectWrapper (){
return(                         
<CustomCreatableSelect
  textField = "xyz"
    isClearable  />

)
}

export default ReactSelectWrapper;