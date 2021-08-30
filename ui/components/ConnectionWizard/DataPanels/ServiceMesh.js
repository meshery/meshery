/* eslint-disable react/display-name */
/* eslint-disable no-unused-vars */
import CloseIcon from "@material-ui/icons/Close";
import {
  withStyles,
  Typography,
  Grid,
  Chip,
  IconButton,
} from "@material-ui/core/";
import { withSnackbar } from "notistack";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { updateProgress, updateK8SConfig, updateAdaptersInfo } from "../../../lib/store";
import { pingAdapterWithNotification, handleDeleteAdapter } from "../helpers/serviceMeshes"
import AdapterChip from "./AdapterChip"

const styles = theme => ({

  infoContainer : {
    width : "20rem",
    height : "15rem",
    padding : "1rem 1rem",
    boxShadow : "0px 1px 6px 1px rgba(0,0,0,0.20)",
    borderRadius : '1rem',
  },
  infoTitle : {
    color : "#647881",
    width : "3rem",
    background : "#F1F3F4",
    padding : ".5rem 5rem .75rem 1.5rem",
    borderRadius : "0.25rem",
    fontSize : ".8rem",
  },
  infoLabel : {
    fontSize : ".9rem",
    color : theme.palette.text.primary,
    marginRight : "1rem",
    fontWeight : 500
  },
  infoData : { fontSize : ".9rem",
    color : theme.palette.text.secondary, },


})

const chipStyles = (theme) => ({ chipIcon : { width : theme.spacing(2.5) },
  chip : { marginRight : theme.spacing(1),
    marginBottom : theme.spacing(1), }, })


const MeshAdapterChip = ({
  handleAdapterClick, isActive, adapter, handleAdapterDelete
}) => {
  let image = isActive ? "/static/img/" + adapter.name?.toLowerCase() + ".svg" : "/static/img/meshery-logo.png";
  return (
    <AdapterChip
      label={adapter?.label.split(":")[0]}
      handleClick={() => handleAdapterClick(adapter?.value)}
      image={image}
      onDelete={isActive
        ? () =>  handleAdapterDelete(adapter?.value)
        : null}
      isActive={isActive}
    />
  )
}


const AdapterPingSnackbarAction = (closeSnackbar) => (key) => (
  <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
    <CloseIcon />
  </IconButton>
)



const ServiceMeshDataPlane= ({
  classes, updateProgress, enqueueSnackbar, closeSnackbar, adapterInfo, isActive, updateAdaptersInfo
}) => {

  const handleAdapterDelete =  handleDeleteAdapter(

    (result) => {
      updateProgress({ showProgress : false });
      if (typeof result !== "undefined") {
        enqueueSnackbar("Adapter was successfully removed!", { variant : "success",
          autoHideDuration : 2000,
          action : AdapterPingSnackbarAction(closeSnackbar) });
        updateAdaptersInfo({ meshAdapters : result });
      }
    },

    (err) => {
      updateProgress({ showProgress : false });
      enqueueSnackbar("Adapter was not removed! "+err, { variant : "error",
        autoHideDuration : 2000,
        action : AdapterPingSnackbarAction(closeSnackbar) });
    }
  )

  return (
    <Grid container className={classes.infoContainer} xs={10}>

      <Grid item xs={12}>
        <MeshAdapterChip adapter={adapterInfo} handleAdapterDelete={handleAdapterDelete} isActive={isActive} handleAdapterClick={(location) => pingAdapterWithNotification(
          updateProgress,
          enqueueSnackbar,
          AdapterPingSnackbarAction(closeSnackbar),
          location)
        }/>
      </Grid>

      <Grid item xs={12}>
        <Typography className={classes.infoTitle}>Details</Typography>
      </Grid>

      <Grid item xs={12} container>
        <Grid item xs={12}>
          <Typography className={classes.infoLabel}>Adapter Server Location: </Typography>
          <Typography className={classes.infoData}>{adapterInfo?.value}</Typography>
        </Grid>
        {adapterInfo?.version ?
          <Grid item xs={12}>
            <Typography className={classes.infoLabel}>Adapter Version</Typography>
            <Typography className={classes.infoData}>
              {adapterInfo.version}
            </Typography>
          </Grid>
          :null}
      </Grid>
    </Grid>

  )
}


const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch),
  updateAdaptersInfo : bindActionCreators(updateAdaptersInfo, dispatch), });

export default withStyles(styles)(connect(null, mapDispatchToProps)(withSnackbar(ServiceMeshDataPlane)))
