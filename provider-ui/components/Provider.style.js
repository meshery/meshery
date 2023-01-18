import { styled } from '@mui/material/styles'
import MenuItem from '@mui/material/MenuItem'

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
