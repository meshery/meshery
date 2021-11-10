import CircularProgress from '@mui/material/CircularProgress';
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";

const CustomLoadingSpinnerWrapper = styled(CircularProgress)(({theme})=>({
    backgroundColor: '#396679', 
    size: 30,
    thickness: 2.7,
}))

export const CustomLoadingSpinner = () =>  { const theme = useTheme();
    return (
      <CustomLoadingSpinnerWrapper/>
    );
  }

export default CustomLoadingSpinnerWrapper;