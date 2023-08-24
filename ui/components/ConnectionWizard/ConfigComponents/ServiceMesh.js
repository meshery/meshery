/* eslint-disable no-unused-vars */
/* eslint-disable react/display-name */
import {
  withStyles,
  Grid,
  Button,
} from "@material-ui/core/";
import { configureAdapterWithNotification } from "../helpers/serviceMeshes";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { updateAdaptersInfo, updateProgress } from "../../../lib/store";
import { useNotification } from "../../../utils/hooks/useNotification";


const styles = () => ({})

const ServiceMeshConfig = ({
  adapterLoc, updateProgress, closeSnackbar, updateAdaptersInfo
}) => {
  const { notify } = useNotification()
  const handleAdapterConfigure = () => {
    configureAdapterWithNotification(notify, updateProgress, adapterLoc, updateAdaptersInfo)
  }

  // const ServiceMeshConfigNotificationAction = (closeSnackbar) => (key) => (
  //   <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
  //     <CloseIcon />
  //   </IconButton>
  // )
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

export default connect(null, mapDispatchToProps)(withStyles(styles)(ServiceMeshConfig))
