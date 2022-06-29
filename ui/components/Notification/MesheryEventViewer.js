import MesherySnackbarWrapper from "./MesherySnackbarWrapper"



function MesheryEventViewer () {

    return(
     <MesherySnackbarWrapper
     key={`event_-_${eventVariant}`}
            variant={eventTypes[eventVariant]
              ? eventTypes[eventVariant].type
              : eventTypes[0].type}
            message={eventSummary}
            details={eventDetails}
            // onClose={handleSnackbarClose}
             />
    )
}

export default MesheryEventViewer
