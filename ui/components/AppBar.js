import { AppBar } from '@mui/material';
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";
import { deepSpaceSparkle } from '../styles/colors';

const CustomAppBarWrapper = styled(AppBar)(({theme})=>({
    backgroundColor: deepSpaceSparkle, 
    position: "sticky",
    padding: theme.spacing(1.4),
    zIndex: 1100,
}))

export default CustomAppBarWrapper;

