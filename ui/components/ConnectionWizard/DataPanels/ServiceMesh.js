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
import { updateProgress, updateK8SConfig } from "../../../lib/store";
import {pingAdapterWithNotification} from "../helpers/serviceMeshes"

const styles = theme => ({

  infoContainer: {
    width: "20rem",
    height: "15rem",
    padding: "1rem 1rem",
    boxShadow: "0px 1px 6px 1px rgba(0,0,0,0.20)",
    borderRadius: '1rem',
  },
  infoTitle: {
    color: "#647881",
    width: "3rem",
    background: "#F1F3F4",
    padding: ".5rem 5rem .75rem 1.5rem",
    borderRadius: "0.25rem",
    fontSize: ".8rem",
  },
  infoLabel: {
    fontSize: ".9rem",
    color: theme.palette.text.primary,
    marginRight: "1rem",
  },
  infoData: {
    fontSize: ".9rem",
    color: theme.palette.text.secondary,
  },


})

const chipStyles = (theme) => ({
  chipIcon: {
    width: theme.spacing(2.5)
  },
  chip: {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
})


const AdapterChip = withStyles(chipStyles)(({classes, handleAdapterClick, adapter}) => {
  console.log(adapter)
  let image = "/static/img/meshery-logo.png";
  let logoIcon = <img src={image} className={classes.chipIcon} />;
  return(
    <Chip
      label={adapter?.label.split(":")[0]}
      onClick={() => handleAdapterClick(adapter?.value)}
      icon={logoIcon}
      className={classes.chip}
      key={`adapters-${adapter?.label}`}
    />
  )
})


const AdapterPingSnackbarAction = (closeSnackbar) => (key) => (
  <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
    <CloseIcon />
  </IconButton>
)



const ServiceMeshDataPlane= ({classes, updateProgress, enqueueSnackbar, closeSnackbar, adapterInfo}) => {
  console.log("in the main component", adapterInfo)

  return (
    <Grid container className={classes.infoContainer} xs={10}>

      <Grid item xs={12}>
        <AdapterChip adapter={adapterInfo} handleAdapterClick={(location) => pingAdapterWithNotification(
          updateProgress,
          enqueueSnackbar,
          AdapterPingSnackbarAction(closeSnackbar),
          location)
        }/>
      </Grid>

      <Grid item xs={12} style={{marginBottom: "1rem"}}>
        <Typography className={classes.infoTitle}>Status</Typography>
      </Grid>

      <Grid item xs={12} container>
        <Grid item xs={12}>
          <Typography className={classes.infoLabel}>Adapter Server Location: {adapterInfo?.value}</Typography>
        </Grid>
        {adapterInfo?.version ? 
          <Grid item xs={12}>
            <Typography className={classes.infoLabel}>Adapter Version</Typography>
            <Typography className={classes.infoData}>
              {adapterInfo.version}
            </Typography>
          </Grid> :null}
      </Grid>
    </Grid>

  )
}


const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default withStyles(styles)(connect(null, mapDispatchToProps)(withSnackbar(ServiceMeshDataPlane)))
