import CircularProgress from '@mui/material/CircularProgress';
import { styled } from "@mui/material/styles";
import { appBarColor } from '../styles/colors';

const CustomLoadingSpinnerWrapper = styled(CircularProgress)(({theme})=>({
    backgroundColor: appBarColor, 
    size: 30,
    thickness: 2.7,
}))

export default CustomLoadingSpinnerWrapper;