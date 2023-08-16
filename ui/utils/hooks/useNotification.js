import { IconButton } from "@material-ui/core"
import { ToggleButtonGroup } from "@mui/material"
import { useSnackbar } from "notistack"
import { iconMedium } from "../../css/icons.styles"
import CloseIcon from "@material-ui/icons/Close";
import { NorthEast } from "@mui/icons-material";
import { useDispatch } from "react-redux";
import { pushEvent ,openEventInNotificationCenter, toggleNotificationCenter } from "../../lib/store";
import moment from "moment";
import { v4 } from "uuid";

// Meshery Event Interface
// event = {
//   event_type ,
//   summary,
//   deleteEvent,
//   eventCause,
//   eventRemediation,
//   eventErrorCode,
//   componentType,
//   componentName,
// }

export const useNotification = () => {

  const { enqueueSnackbar,closeSnackbar }= useSnackbar()

  const dispatch = useDispatch()

  const openEvent = (eventId) => {
    dispatch(openEventInNotificationCenter({
      eventId
    }))
    dispatch(toggleNotificationCenter())
  }
  const notify = ({ id=null,message,description=null,event_type,timestamp=null,customEvent=null,showInNotificationCenter=true }) => {

    timestamp = timestamp ?? moment.utc().valueOf()
    id = id  || v4()


    if (showInNotificationCenter) {
      dispatch(pushEvent({
        summary : message ,
        id,
        event_type,
        timestamp,
        eventCause : description,
        ...customEvent
      }))
    }

    enqueueSnackbar(message,{
      variant : event_type.type,
      action : function Action(key) {
        return (
          <ToggleButtonGroup>
            <IconButton key={`closeevent-${id}`} aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon style={iconMedium} />
            </IconButton>
            {showInNotificationCenter &&
              <IconButton key={`openevent-${id}`} aria-label="Open" color="inherit" onClick={() => openEvent(id)}>
                <NorthEast style={iconMedium} />
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