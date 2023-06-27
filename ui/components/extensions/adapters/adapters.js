import { useEffect, useState } from "react";
import { withRouter } from "next/router";
import { iconMedium, extensionStyles as styles } from "../../../css/icons.styles";
import { Grid, Typography, IconButton, Switch } from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import { withStyles } from "@material-ui/core/styles";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { useSnackbar } from "notistack";
import { updateProgress } from "../../../lib/store";
import { adaptersList } from "./constants";
import changeAdapterState from '../../graphql/mutations/AdapterStatusMutation';
import { LARGE_6_MED_12_GRID_STYLE } from "../../../css/grid.styles";

const Adapters = ({ updateProgress, classes, meshAdapters, meshAdaptersts }) => {

  // States.
  const [availableAdapters, setAvailableAdapters] = useState(adaptersList);
  const [mts, setMts] = useState(new Date());

  // Hooks.
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  // useEffects.
  useEffect(() => {
    updateProgress({ showProgress : true })
  }, [])

  useEffect(() => {
    if (meshAdaptersts > mts) {
      let currentAdaptersList = { ...adaptersList };
      meshAdapters.forEach(element => {
        const adapterId = element.name;
        if (adapterId && currentAdaptersList[adapterId]) {
          currentAdaptersList[adapterId].enabled = true;
          currentAdaptersList[adapterId].url = element.adapter_location;
        }
      })
      setAvailableAdapters(currentAdaptersList);
      setMts(meshAdaptersts);
      updateProgress({ showProgress : false })
    }

  }, [meshAdapters, meshAdaptersts])


  // Handlers.
  const handleAdapterDeployment = (payload, msg, selectedAdapter, adapterId) => {
    updateProgress({ showProgress : true });

    changeAdapterState((response, errors) => {
      updateProgress({ showProgress : false });

      if (errors !== undefined) {
        // Toggle the switch to it's previous state if the request fails.
        setAvailableAdapters({ ...availableAdapters, [adapterId] : { ...selectedAdapter, enabled : !selectedAdapter.enabled } });
        handleError(msg);
      }

      enqueueSnackbar(`Adapter ${response.adapterStatus.toLowerCase()}`, {
        variant : "info",
        autoHideDuration : 2000,
        action : (key) => (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon style={iconMedium} />
          </IconButton>
        ),
      }
      );
    }, payload);

  };

  const handleError = (msg) => (error) => {
    updateProgress({ showProgress : false });
    enqueueSnackbar(`${msg}: ${error}`, {
      variant : "error", preventDuplicate : true,
      action : (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
          <CloseIcon style={iconMedium} />
        </IconButton>
      ),
      autoHideDuration : 7000,
    });
  };

  const handleToggle = (selectedAdapter, adapterId) => {
    setAvailableAdapters({ ...availableAdapters, [adapterId] : { ...selectedAdapter, enabled : !selectedAdapter.enabled } });
    let payload = {}, msg = "";
    if (!selectedAdapter.enabled) {
      payload = {
        status : "ENABLED",
        adapter : selectedAdapter.label,
        targetPort : selectedAdapter.defaultPort
      };
      msg = "Unable to Deploy adapter";
    } else {
      payload = {
        status : "DISABLED",
        adapter : selectedAdapter.label,
        targetPort : selectedAdapter.defaultPort
      };
      msg = "Unable to Undeploy adapter";
    }
    handleAdapterDeployment(payload, msg, selectedAdapter, adapterId);
  }

  // Render.
  return (
    <>
      {
        Object.entries(availableAdapters).map(([adapterId, adapter]) =>
          <Grid item {...LARGE_6_MED_12_GRID_STYLE} key={adapterId}>
            <div className={classes.card} >

              <Typography className={classes.frontContent} variant="h5" component="div">
                Meshery Adapter for {adapter.name}
              </Typography>

              <Typography className={classes.frontSideDescription} variant="body">
                <img className={classes.img} src={adapter.imageSrc} />
                <div style={{
                  display : "inline", position : "relative",
                }}>
                  {adapter.description}
                </div>
              </Typography>

              <Grid container spacing={2} className={classes.grid} direction="row" justifyContent="space-between" alignItems="baseline" style={{ position : "absolute", paddingRight : "3rem", paddingLeft : ".5rem", bottom : "1.5rem", }}>
                <Typography variant="subtitle2" style={{ fontStyle : "italic" }}>
                  <a href="https://docs.meshery.io/concepts/architecture/adapters" target="_blank" rel="noreferrer" className={classes.link}>Open Adapter docs</a>
                </Typography>

                <div style={{ textAlign : "right" }}>
                  <Switch
                    checked={adapter.enabled}
                    onChange={() => handleToggle(adapter, adapterId)}
                    name="OperatorSwitch"
                    color="primary"
                    classes={{
                      switchBase : classes.switchBase,
                      track : classes.track,
                      checked : classes.checked,
                    }}
                  />
                </div>
              </Grid>
            </div>
          </Grid>
        )
      }
    </>
  )
}

const mapStateToProps = (state) => {
  const meshAdapters = state.get("meshAdapters");
  const meshAdaptersts = state.get("meshAdaptersts");

  return {
    meshAdapters,
    meshAdaptersts
  };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress : bindActionCreators(updateProgress, dispatch),
});



export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(Adapters))
);