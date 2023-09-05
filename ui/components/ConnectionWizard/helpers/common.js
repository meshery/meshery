/* eslint-disable react/display-name */
import { IconButton } from "@material-ui/core"
import CloseIcon from "@material-ui/icons/Close";
import { EVENT_TYPES } from "../../../lib/event-types";

export const closeButtonForSnackbarAction = (closeSnackbar) => (key) => (
  <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
    <CloseIcon />
  </IconButton>
)

export const successHandlerGenerator = (notify, msg, cb) => (res) => {
  if (res !== undefined) {
    if (cb !== undefined) cb(res)
    notify({ message: `${msg}: ${res}`, details: res, event_type: EVENT_TYPES.SUCCESS })
  }
}

export const errorHandlerGenerator = (notify, msg, cb) => (err) => {
  if (cb !== undefined) cb(err)
  err = typeof err !== "string" ? err.toString() : err
  notify({ message: `${msg}`, details: err, event_type: EVENT_TYPES.ERROR })
}