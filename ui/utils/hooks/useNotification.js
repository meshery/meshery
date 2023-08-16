import { IconButton } from "@material-ui/core"
import { ToggleButtonGroup } from "@mui/material"
import { useSnackbar } from "notistack"
import { iconMedium } from "../../css/icons.styles"
import CloseIcon from "@material-ui/icons/Close";
import { useDispatch } from "react-redux";
import { pushEvent ,openEventInNotificationCenter, toggleNotificationCenter } from "../../lib/store";
import moment from "moment";
import { v4 } from "uuid";
import BellIcon from '@material-ui/icons/Notifications';

export const useNotification = () => {

  const { enqueueSnackbar,closeSnackbar }= useSnackbar()

  const dispatch = useDispatch()

  const openEvent = (eventId) => {
    dispatch(toggleNotificationCenter())
    dispatch(openEventInNotificationCenter({
      eventId
    }))
  }
  const notify = ({ id=null,message,details=null,event_type,timestamp=null,customEvent=null,showInNotificationCenter=true }) => {

    timestamp = timestamp ?? moment.utc().valueOf()
    id = id  || v4()


    if (showInNotificationCenter) {
      dispatch(pushEvent({

        ...customEvent,
        summary : message ,
        id,
        event_type,
        timestamp,
        details,
      }))
    }

    enqueueSnackbar(
      message,{
        variant : event_type.type,
        action : function Action(key) {
          return (
            <ToggleButtonGroup>
              <IconButton key={`closeevent-${id}`} aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon style={iconMedium} />
              </IconButton>
              {showInNotificationCenter &&
              <IconButton key={`openevent-${id}`} aria-label="Open" color="inherit" onClick={() => openEvent(id)}>
                <BellIcon style={iconMedium} />
              </IconButton> }
            </ToggleButtonGroup>
          );
        },
      })
  }

  return {
    notify
  }
}

// For use in legacy class based components
export function withNotify( Component ){
  return function WrappedWithNotify(props) {
    const { notify } = useNotification()
    return <Component {...props} notify={notify} />
  }
}