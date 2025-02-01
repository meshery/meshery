import { styled } from "@mui/material/styles";
import MenuItem from "@mui/material/MenuItem";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import Typography from "@mui/material/Typography";
import Popover from "@mui/material/Popover";
export const CustomDiv = styled("div")(({ theme }) => ({
  width : "60%",
  marginLeft : "auto",
  marginRight : "auto",
  marginTop : theme.spacing(3),
}));

export const CustomTypography = styled(Typography)(() => ({
  "@font-face" : {
    fontFamily : "Qanelas Soft",
    src : "url('/provider/static/fonts/qanelas-soft/QanelasSoftRegular.otf') format('opentype')",
  },
  fontWeight : 400,
  fontStyle : "normal",
  fontSize : "1.5rem",
  lineHeight : "2rem",
  letterSpacing : "0.01em",
  color : "#222",
}));

export const MesheryLogo = styled("img")(({ theme }) => ({
  width : theme.spacing(50),
  maxWidth : "100%",
  height : "auto",
}));

export const MenuProviderDisabled = styled(MenuItem)(() => ({
  display : "flex",
  justifyContent : "space-between",
}));

export const CustomDialog = styled(Dialog)(({ theme }) => ({
  "& .MuiDialogContent-root" : {
    padding : theme.spacing(2),
  },
  "& .MuiDialogActions-root" : {
    padding : theme.spacing(1),
  },
  "& .MuiDialogContentText-root > a" : {
    color : "#222",
  },
}));
export const CustomDialogActions = styled(DialogActions)(({ theme }) => ({
  display : "flex",
  justifyContent : "space-between",
  background : "#eee",
  padding : theme.spacing(2),
  "& div > a" : {
    color : "#222",
  },
}));
export const StyledPopover = styled(Popover)(() => ({
  ".MuiPaper-root" : {
    backgroundColor : "#202020", // Background color
    color : "white", // Optional: for readable text
  },
}));
export const LearnMore = styled("div")(() => ({
  width : "60%",
  marginLeft : "auto",
  marginRight : "auto",
  marginTop : "3rem",
}));
