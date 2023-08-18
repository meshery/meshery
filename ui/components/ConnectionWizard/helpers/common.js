/* eslint-disable react/display-name */
import { IconButton } from "@material-ui/core"
import CloseIcon from "@material-ui/icons/Close";
import { useNotification } from "../../../utils/hooks/useNotification";
import { EVENT_TYPES } from "../../../lib/event-types";

export const closeButtonForSnackbarAction = (closeSnackbar) => (key) => (
  <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
    <CloseIcon />
  </IconButton>
)

export const successHandlerGenerator = (msg, cb) => (res) => {
  const { notify } = useNotification()
  if (res !== undefined) {
    if (cb !== undefined) cb(res)
    notify({ message : `${msg}: ${res}`, event_type : EVENT_TYPES.SUCCESS })
  }
}

export const errorHandlerGenerator = (msg, cb) => (err) => {
  const { notify } = useNotification()
  if (cb !== undefined) cb(err)
  notify({ message : `${msg}: ${err}`, event_type : EVENT_TYPES.ERROR })
}
