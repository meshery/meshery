import { styled } from '@mui/material/styles'
import MenuItem from '@mui/material/MenuItem'
import Dialog from "@mui/material/Dialog";

export const CustomDiv = styled('div')(({ theme }) => ({
  width : '60%',
  marginLeft : 'auto',
  marginRight : 'auto',
  marginTop : theme.spacing(5)
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
}));