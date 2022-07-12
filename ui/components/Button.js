import{ Button, Box} from '@mui/material/';
import { styled } from "@mui/material/styles";
import { deepSpaceSparkle } from '../styles/colors';

export const MetricsButton = styled(Button)(({theme})=>({
    backgroundColor: deepSpaceSparkle,
}))


 const RightAlignButtonContainer = styled(Box)(({theme})=>({
      display : 'flex',
    justifyContent : 'flex-end',
    width : '100%',
    marginTop : theme.spacing(3),
}))    

export const RightAlignButton = ({title}) => {
    return (
        <RightAlignButtonContainer>
            <Button
                type="submit" 
                variant="contained"
                color="primary"
                size="large">
                {title}
            </Button>
        </RightAlignButtonContainer>
    )
} 

