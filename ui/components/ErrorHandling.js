import { IconButton } from "@material-ui/core";
import { Cancel } from "@material-ui/icons";
import { useSnackbar } from "notistack";

/**
 *  Show Snackbar when error occurs. Can be used in catch blocks
 *  of functions.
 *
 * @param {Object} err
 * @param {string} prefixMessage
 * @param {("normal"|"fatal")} severity
 */
function handleError(err, prefixMessage, severity = "normal") {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar()
  if (severity.length > 0) {
    severity = severity[0].toUpperCase() + severity.slice(1);
  }

  console.error("an error occured with severity: ", severity, { err })

  return enqueueSnackbar(
    `${severity}: ${prefixMessage}: ${err?.message}`,
    {
      variant : "error",
      autoHideDuration : 8000,
      preventDuplicate : true,
      action : (key) => (
        <IconButton
          onClick={() => closeSnackbar(key)}
          color="secondary"
        >
          <Cancel />
        </IconButton>
      )
    })
}

export default handleError