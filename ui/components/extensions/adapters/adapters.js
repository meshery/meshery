import { isNil, isUndefined } from "lodash";
import { useEffect, useState } from "react";
import { withRouter } from "next/router";
import { extensionStyles as styles } from "../../../css/icons.styles";
import { Grid, Typography, Switch, Tooltip } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { updateProgress } from "../../../lib/store";
import { adaptersList } from "./constants";
import changeAdapterState from '../../graphql/mutations/AdapterStatusMutation';
import { LARGE_6_MED_12_GRID_STYLE } from "../../../css/grid.style";
import { promisifiedDataFetch } from "../../../lib/data-fetch";
import { useNotification } from "../../../utils/hooks/useNotification";
import { EVENT_TYPES } from "../../../lib/event-types";
import CircleIcon from '@mui/icons-material/Circle';

const Adapters = ({ updateProgress, classes }) => {

  // States.
  const [availableAdapters, setAvailableAdapters] = useState(adaptersList);
  const [adapterStates, setAdapterStates] = useState({});

  // Hooks.
  const { notify } = useNotification();

  // useEffects.
  useEffect(() => {
    handleAdapterSync();
    Object.keys(adaptersList).forEach( adapterId => adapterStates[adapterId] = adaptersList[adapterId].enabled? "Deployed" : "Undeployed" );
  }, [])

  // Handlers.
  const handleAdapterSync = async (showLoader = true) => {
    showLoader && updateProgress({ showProgress : true });

    promisifiedDataFetch(
      "/api/system/sync",
      {
        method : "GET",
        credentials : "include",
      },
    ).then((result) => {
      showLoader && updateProgress({ showProgress : false });

      if (!isUndefined(result)) {
        // Deep copying to avoid mutability.
        // Ref: https://developer.mozilla.org/en-US/docs/Web/API/structuredClone
        let currentAdaptersList = structuredClone(adaptersList);

        result.meshAdapters.forEach(element => {
          const adapterId = element.name;
          if (adapterId && currentAdaptersList[adapterId]) {
            currentAdaptersList[adapterId].enabled = true;
            currentAdaptersList[adapterId].url = element.adapter_location;
          }
        })
        setAvailableAdapters(currentAdaptersList);
      }
    }).catch(() =>
      handleError("Unable to fetch list of adapters.")
    );
  }

  const handleAdapterDeployment = (payload, msg, selectedAdapter, adapterId) => {
    updateProgress({ showProgress : true });

    changeAdapterState((response, errors) => {
      updateProgress({ showProgress : false });

      if (!isNil(errors)) {
        // Toggle the switch to it's previous state if the request fails.
        setAvailableAdapters({ ...availableAdapters, [adapterId] : { ...selectedAdapter, enabled : !selectedAdapter.enabled } });
        handleError(msg);
      } else {
        notify({ message : `Adapter ${response.adapterStatus.toLowerCase()}`, event_type : EVENT_TYPES.SUCCESS });
        setAdapterStates({ ...adapterStates, [adapterId] : payload.status === "ENABLED" ? "Deployed" : "Undeployed" });
      }
    }, payload);



  };

  const handleError = (msg) => (error) => {
    updateProgress({ showProgress : false });
    notify({ message : msg, event_type : EVENT_TYPES.ERROR, details : error.toString() })
  };

  const handleToggle = (selectedAdapter, adapterId) => {
    setAvailableAdapters({ ...availableAdapters, [adapterId] : { ...selectedAdapter, enabled : !selectedAdapter.enabled } });
    let payload = {}, msg = "";
    if (!selectedAdapter.enabled) {
      setAdapterStates({ ...adapterStates, [adapterId] : "Deploying" });
      payload = {
        status : "ENABLED",
        adapter : selectedAdapter.label,
        targetPort : selectedAdapter.defaultPort
      };
      msg = "Unable to Deploy adapter";
    } else {
      setAdapterStates({ ...adapterStates, [adapterId] : "Undeploying" });
      payload = {
        status : "DISABLED",
        adapter : selectedAdapter.label,
        targetPort : selectedAdapter.defaultPort
      };
      msg = "Unable to Undeploy adapter";
    }
    handleAdapterDeployment(payload, msg, selectedAdapter, adapterId);
  }

  const getColorForState = (adapterId) => {
    if ( adapterStates[adapterId] === "Deploying" ) return "#ffcc00";
    if ( adapterStates[adapterId] === "Undeploying" ) return "#ff9400";
    if ( availableAdapters[adapterId].enabled ) return "#008000";
    if ( !availableAdapters[adapterId].enabled ) return "#808080";
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
                <Grid item xs={7} md={8}>
                  <Typography variant="subtitle2" style={{ fontStyle : "italic" }}>
                    <a href="https://docs.meshery.io/concepts/architecture/adapters" target="_blank" rel="noreferrer" className={classes.link}>Open Adapter docs</a>
                  </Typography>
                </Grid>
                <Grid item xs={2} md={1} style={{ margin : "auto" }}>
                  <Tooltip title={ adapterStates[adapterId] }>
                    <CircleIcon style={{ color : getColorForState(adapterId) }} />
                  </Tooltip>
                </Grid>
                <Grid item xs={3}>
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
              </Grid>
            </div>
          </Grid>
        )
      }
    </>
  )
}


const mapDispatchToProps = (dispatch) => ({
  updateProgress : bindActionCreators(updateProgress, dispatch),
});



export default withStyles(styles)(
  connect(() => { }, mapDispatchToProps)(withRouter(Adapters))
);