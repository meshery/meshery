import { NOTIFICATIONCOLORS } from "../../themes"
import AlertIcon from "../../assets/icons/AlertIcon";
import ErrorIcon from "../../assets/icons/ErrorIcon.js"
import { Colors } from "../../themes/app";
import ReadIcon from "../../assets/icons/ReadIcon";

export const SEVERITY = {
  INFO : "informational",
  ERROR : "error",
  WARNING : "warning",
  // SUCCESS: "success"
}

export const STATUS = {
  READ : "read",
  UNREAD : "unread",
}

export const STATUS_STYLE = {
  [STATUS.READ] : {
    icon : ReadIcon,
    color : Colors.charcoal
  }
}

export const SEVERITY_STYLE = {
  [SEVERITY.INFO] : {
    icon : ErrorIcon,
    color : NOTIFICATIONCOLORS.INFO
  },
  [SEVERITY.ERROR] : {
    icon : ErrorIcon,
    color : NOTIFICATIONCOLORS.ERROR
  },
  [SEVERITY.WARNING] : {
    icon : AlertIcon,
    color : NOTIFICATIONCOLORS.WARNING
  },

}