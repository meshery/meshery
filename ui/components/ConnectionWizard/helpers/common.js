import { IconButton } from "@material-ui/core"
import CloseIcon from "@material-ui/icons/Close";

export const showProgress = (updateProgress) => updateProgress({showProgress: true})
export const hideProgress = (updateProgress) => updateProgress({showProgress: false})


export const closeButtonForSnackbarAction = (closeSnackbar) => (key) => (
  <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
    <CloseIcon />
  </IconButton>
)

export const successHandlerGenerator = (snackbar, action, msg, cb) => (res) => {
    if(res !== undefined)
      {
        if(cb !== undefined) cb(res)
        snackbar(msg, {
          variant: "success",
          action,
          autoHideDuration: 7000,
        })
      }
}

export const errorHandlerGenerator = (snackbar, action,msg, cb) => (err) => {
      if(cb !== undefined) cb(err)      
      snackbar(`${msg}: ${err}`, {
        variant: "error",
        action,
        autoHideDuration: 7000,
      })
}
