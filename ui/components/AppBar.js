import { AppBar } from '@mui/material';
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";
import { seaGreenColor } from '../styles/colors';

const CustomAppBarWrapper = styled(AppBar)(({theme})=>({
    backgroundColor: seaGreenColor, 
    position: "sticky",
    elevation: 0,
    padding: theme.spacing(1.4),
    zIndex: 1100,
}))

export default CustomAppBarWrapper;

