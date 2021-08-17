/* eslint-disable no-unused-vars */
/* eslint-disable react/display-name */
import {
  withStyles,
  Grid,
  Button,
  IconButton,
} from "@material-ui/core/";
import { configureAdapterWithNotification } from "../helpers/serviceMeshes";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { updateAdaptersInfo, updateProgress } from "../../../lib/store";


const styles = () => ({})

const ServiceMeshConfig = ({
  classes, adapterLoc, updateProgress, enqueueSnackbar, closeSnackbar, updateAdaptersInfo
}) => {

  const handleAdapterConfigure = () => {
    configureAdapterWithNotification(enqueueSnackbar, updateProgress, ServiceMeshConfigNotificationAction(closeSnackbar), adapterLoc, updateAdaptersInfo)
  }

  const ServiceMeshConfigNotificationAction = (closeSnackbar) => (key) => (
    <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
      <CloseIcon />
    </IconButton>
  )
  return (
    <>
      <Grid item xs={12} style={{ padding : "1rem" }}>
        <Button
          id="service-mesh-config-setup"
          name="serviceMeshAdapterConfigureButton"
          variant="contained"
          fullWidth
          color="primary"
          onClick={handleAdapterConfigure}
        >
          Connect
        </Button>
      </Grid>
    </>
  )
}


const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch),
  updateAdaptersInfo : bindActionCreators(updateAdaptersInfo, dispatch), });

export default connect(null, mapDispatchToProps)(withStyles(styles)(withSnackbar(ServiceMeshConfig)))
