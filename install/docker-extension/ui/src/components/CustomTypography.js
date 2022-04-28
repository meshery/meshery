import Typography from "@mui/material/Typography";
import { styled } from "@mui/material/styles";

const CustomTypography = styled(Typography)(({theme})=>({  
  root: {
  color: theme.palette.common.white
}})) 
  
export default CustomTypography;