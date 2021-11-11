import Button from '@mui/material/Button';
import { styled } from "@mui/material/styles";
import { seaGreenColor } from '../styles/colors';

export const MetricsButton = styled(Button)(({theme})=>({
    backgroundColor: seaGreenColor,
    width: "240px",
}))
