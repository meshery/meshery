import { makeStyles } from "@material-ui/core"

const opacity = 0.5

export const disabledStyle = {
  pointerEvents : 'none',
  opacity
}

export const cursorNotAllowed = {
  cursor : "not-allowed",
}

export const cursorNotAllowedWithLowOpacity = {
  cursor : "not-allowed",
  opacity
}

export const disabledStylesMui = makeStyles(() => ({
  disabled : disabledStyle,
  cursorNotAllowed
}))