import { styled } from '@mui/material/styles'
import MenuItem from '@mui/material/MenuItem'
import DialogTitle from '@mui/material/DialogTitle'
import IconButton from '@mui/material/IconButton'

export const Div = styled('div')(({ theme }) => ({
  width: '60%',
  marginLeft: 'auto',
  marginRight: 'auto',
  marginTop: theme.spacing(5)
}))

export const MesheryLogo = styled('img')(({ theme }) => ({
  width: theme.spacing(50),
  maxWidth: '100%',
  height: 'auto'
}))

export const MenuProviderDisabled = styled(MenuItem)(() => ({
  display: 'flex',
  justifyContent: 'space-between'
}))

export const CustomDialogTitle = styled(DialogTitle)(({ theme }) => ({
  margin: 0,
  padding: theme.spacing(2),
}))

export const CustomIconButton = styled(IconButton)(({ theme }) => ({
  position: "absolute",
  right: theme.spacing(1),
  top: theme.spacing(1),
  color: theme.palette.grey[500],
}))