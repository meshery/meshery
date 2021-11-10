import { AppBar } from '@mui/material';
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";
import theme from '../styles/theme'

const CustomAppBarWrapper = styled(AppBar)(({theme})=>({
    backgroundColor: '#396679', 
    position: "sticky",
    elevation: 0,
    padding: theme.spacing(1.4),
    zIndex: 1100,
}))

export const CustomAppBar = () =>  { const theme = useTheme();
    return (
      <CustomAppBarWrapper>
        
      </CustomAppBarWrapper>
    );
  }

export default CustomAppBarWrapper;

