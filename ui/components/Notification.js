import { IconButton } from '@mui/material';
import HighlightOffIcon from '@mui/icons-material/HighlightOff';

export const closeButtonForSnackbarAction = (closeSnackbar) => (key) => (
  <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
    <HighlightOffIcon />
  </IconButton>
)

export const successHandlerGenerator = (snackbar, action, msg, cb) => (res) => {
  if (res !== undefined) {
    hideProgress()
    if (cb !== undefined) cb(res)
    snackbar(msg, { variant : "success",
      action,
      autoHideDuration : 7000, })
  }
}

export const errorHandlerGenerator = (snackbar, action,msg, cb) => (err) => {
  hideProgress()
  if (cb !== undefined) cb(err)
  snackbar(`${msg}: ${err}`, { variant : "error",
    action,
    autoHideDuration : 7000, })
}