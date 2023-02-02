import { IconButton } from "@material-ui/core";
import { Cancel } from "@material-ui/icons";
import { useSnackbar } from "notistack";
import { iconMedium } from "../css/icons.styles";

/**
 *  Show Snackbar when error occurs. Can be used in catch blocks
 *  of functions.
 *
 */
function HandleError() {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  /**
    *
    * @param {Object} err
    * @param {string} prefixMessage
    * @param {("error"|"warning")} variant
    */
  const errorH = (err, prefixMessage, variant) => {
    console.error("an error occured with severity: ", variant, { err })
    return enqueueSnackbar(
      `${prefixMessage}: ${err?.message}`,
      {
        variant : variant || "error",
        autoHideDuration : 8000,
        preventDuplicate : true,
        action : (key) => (
          <IconButton
            onClick={() => closeSnackbar(key)}
            color="secondary"
            style={iconMedium}
          >
            <Cancel />
          </IconButton>
        )
      })
  }

  return errorH
}

export default HandleError