import MesherySnackbarWrapper from "./MesherySnackbarWrapper"


/* Component for the  event(notification) snackbar */ 
function MesheryEventViewer ({eventVariant, eventSummary, eventDetails}) {

    return(
     <MesherySnackbarWrapper
     key={`event_-_${eventVariant}`}
            message={eventSummary}
            details={eventDetails}
             />
    )
}

export default MesheryEventViewer
