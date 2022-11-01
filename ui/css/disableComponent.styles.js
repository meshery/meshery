import { makeStyles } from "@material-ui/core"

const opacity = 0.5

export const disabledStyle = {
  pointerEvents : 'none',
  opacity
}

export const cursorNotAllowed = {
  cursor : "not-allowed",
  opacity
}

export const disabledStylesMui = makeStyles(() => ({
  disabled : disabledStyle,
  cursorNotAllowed
}))

export const enabledComponents = {
  navigator : {
    extension : {
      meshmap : {
        designer : [
          "design",
          "application",
          "filter",
          "save"
        ]
      }
    }
  },
  header : {
    settings : "*",
    contextSwitcher : "*",
    profile : "*"
  }
}