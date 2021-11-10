import { AppBar } from '@mui/material';
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";
import theme from '../styles/theme'
import { appBarColor } from '../styles/colors';

const CustomAppBarWrapper = styled(AppBar)(({theme})=>({
    backgroundColor: appBarColor, 
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

