import CircularProgress from '@mui/material/CircularProgress';
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";

const CustomLoadingSpinnerWrapper = styled(CircularProgress)(({theme})=>({
    color: "primary", 
    size: 30,
    thickness: 2.7,
}))

const CustomLoadingSpinner = () =>  { const theme = useTheme();
    return (
      <CustomLoadingSpinnerWrapper/>
    );
  }

export default CustomLoadingSpinner;