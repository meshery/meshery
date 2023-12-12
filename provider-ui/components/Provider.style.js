import { styled } from '@mui/material/styles'
import MenuItem from '@mui/material/MenuItem'
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";

export const CustomDiv = styled('div')(({ theme }) => ({
  width : '60%',
  marginLeft : 'auto',
  marginRight : 'auto',
  marginTop : theme.spacing(3)
}))

export const MesheryLogo = styled('img')(({ theme }) => ({
  width : theme.spacing(50),
  maxWidth : '100%',
  height : 'auto'
}))

export const MenuProviderDisabled = styled(MenuItem)(() => ({
  display : 'flex',
  justifyContent : 'space-between'
}))

export const CustomDialog = styled(Dialog)(({ theme }) => ({
  '& .MuiDialogContent-root' : {
    padding : theme.spacing(2),
  },
  '& .MuiDialogActions-root' : {
    padding : theme.spacing(1),
  },
  '& .MuiDialogContentText-root > a' : {
    color : "#222",
  },
}));
export const CustomDialogActions = styled(DialogActions)(({ theme }) => ({
  display : "flex", justifyContent : "space-between",
  background : "#eee",
  padding : theme.spacing(2),
  '& div > a' : {
    color : "#222",
  },
}));

export const LearnMore = styled('div')(() => ({
  width : '60%',
  marginLeft : 'auto',
  marginRight : 'auto',
  marginTop : '3rem',
}))