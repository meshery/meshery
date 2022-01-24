import CloseIcon from "@material-ui/icons/Close";
import IconButton from "@material-ui/core/IconButton";

function handleError(msg, enqueueSnackbar, closeSnackbar, updateProgress) {

  return function (error) {
    (updateProgress({ showProgress : false }))

    enqueueSnackbar(`${msg}: ${error}`, { variant : "error",
      action : function Action(key) {
        return (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        );
      },
      autoHideDuration : 8000, });
  };
}

export default handleError;