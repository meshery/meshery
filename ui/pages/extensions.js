import { Grid, Typography, Button, Switch, IconButton } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import React, { useEffect, useState } from "react";
import CloseIcon from "@material-ui/icons/Close";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { toggleCatalogContent } from "../lib/store";
import Head from 'next/head';
import { withSnackbar } from "notistack";
import dataFetch from "../lib/data-fetch";
import { EXTENSIONS } from "../utils/Enum";
import { extensionStyles as styles } from "../css/icons.styles";
import { Adapters } from "../components/extensions"
import { LARGE_6_MED_12_GRID_STYLE } from "../css/grid.styles";

const MeshMapSignUpcard = ({ classes, hasAccessToMeshMap = false }) => {

  const handleSignUp = (e) => {
    window.open("https://layer5.io/meshmap", "_blank")
    e.stopPropagation();
  };

  return (
    <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
      <div className={classes.card}>
        <Typography className={classes.frontContent} variant="h5" component="div">
          MeshMap
        </Typography>

        <Typography className={classes.frontSideDescription} variant="body">
          <img className={classes.img} src="/static/img/meshmap.svg" />
          Collaboratively design and manage your Kubernetes clusters, service mesh deployments, and cloud native apps.
          MeshMap is now in private beta. {!hasAccessToMeshMap && "Sign-up today to for early access!"}
        </Typography>
        {<div style={{ textAlign : "right" }}>
          <Button
            variant="contained"
            color="primary"
            disabled={hasAccessToMeshMap}
            className={classes.button}
            onClick={(e) => handleSignUp(e)}>
            {hasAccessToMeshMap ? "Enabled" : "Sign Up"}
          </Button>
        </div>}
      </div>
    </Grid>
  )
}

export const WrappedMeshMapSignupCard = withStyles(styles)(MeshMapSignUpcard);

const Extensions = ({ classes, toggleCatalogContent, enqueueSnackbar, closeSnackbar, capabilitiesRegistry }) => {
  const [catalogContent, setCatalogContent] = useState(true);
  const [extensionPreferences, setExtensionPreferences] = useState({})
  const [hasAccessToMeshMap, setHasAccessToMeshMap] = useState(false)

  const handleToggle = () => {
    toggleCatalogContent({ catalogVisibility : !catalogContent });
    setCatalogContent(!catalogContent);
    handleCatalogPreference(!catalogContent);
  }

  useEffect(() => {
    dataFetch(
      "/api/user/prefs",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        if (result) {
          setExtensionPreferences(result?.usersExtensionPreferences)
          setCatalogContent(result?.usersExtensionPreferences?.catalogContent)
        }
      },
      err => console.error(err)
    )
  }, []);

  const handleCatalogPreference = (catalogPref) => {
    let body = Object.assign({}, extensionPreferences)
    body["catalogContent"] = catalogPref

    dataFetch(
      "/api/user/prefs",
      {
        method : "POST",
        credentials : "include",
        body : JSON.stringify({ usersExtensionPreferences : body })
      },
      () => {
        enqueueSnackbar(`Catalog Content was ${catalogPref ? "enab" : "disab"}led`,
          {
            variant : 'success',
            autoHideDuration : 4000,
            action : (key) => (
              <IconButton
                key="close"
                aria-label="Close"
                color="inherit"
                onClick={() => closeSnackbar(key)}
              >
                <CloseIcon />
              </IconButton>
            ),
          });
      },
      err => console.error(err),
    )
  }

  useEffect(() => {
    const meshMapExtensionExists = capabilitiesRegistry?.extensions?.navigator?.filter(
      (val) => val.title.toLowerCase() === EXTENSIONS.MESHMAP
    );
    if (typeof meshMapExtensionExists === "object" && meshMapExtensionExists.length)
      setHasAccessToMeshMap(true);
  }, [])


  return (
    <React.Fragment>
      <Head>
        <title>Extensions | Meshery</title>
      </Head>
      <Grid container spacing={1} >
        <WrappedMeshMapSignupCard hasAccessToMeshMap={hasAccessToMeshMap} />
        <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
          <div className={classes.card} >
            <Typography className={classes.frontContent} variant="h5" component="div">
              {"Meshery Catalog"}
            </Typography>

            <Typography className={classes.frontSideDescription} variant="body">
              <img className={classes.img} src="/static/img/meshery_catalog.svg" />
              <div style={{
                display : "inline", position : "relative",
              }}>
                Enable access to the cloud native catalog, supporting <a href="https://service-mesh-patterns.github.io/service-mesh-patterns" className={classes.link
                }>Service Mesh Patterns</a>, WebAssembly filters, eBPF programs (<span style={{ fontStyle : "italic" }}>soon</span>), and OPA policies (<span style={{ fontStyle : "italic" }}>soon</span>). Import any catalog item and customize.
              </div>
            </Typography>

            <Grid container spacing={2} className={classes.grid} direction="row" justifyContent="space-between" alignItems="baseline" style={{ position : "absolute", paddingRight : "3rem", paddingLeft : ".5rem", bottom : "1.5rem", }}>
              <Typography variant="subtitle2" style={{ fontStyle : "italic" }}>
                Explore the <a href="https://meshery.io/catalog" target="_blank" rel="noreferrer" className={classes.link}>Meshery Catalog</a>
              </Typography>

              <div style={{ textAlign : "right" }}>
                <Switch
                  checked={catalogContent}
                  onChange={handleToggle}
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
        <Adapters />
      </Grid>
    </React.Fragment>
  )
}

const mapStateToProps = state => ({
  catalogVisibility : state.get('catalogVisibility'),
  capabilitiesRegistry : state.get('capabilitiesRegistry')
})

const mapDispatchToProps = dispatch => ({
  toggleCatalogContent : bindActionCreators(toggleCatalogContent, dispatch)
})

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(Extensions)));