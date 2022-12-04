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

const styles = (theme) => ({
  button : {
    borderRadius : 5,
    minWidth : 100,
    color : "#fff",
    "&:hover" : {
      boxShadow : "0px 2px 4px -1px rgb(0 0 0 / 20%), 0px 4px 5px 0px rgb(0 0 0 / 14%), 0px 1px 10px 0px rgb(0 0 0 / 12%)"
    },
  },
  card : {
    padding : theme.spacing(3),
    borderRadius : theme.spacing(1),
    transformStyle : "preserve-3d",
    boxShadow : "0 4px 8px 0 rgba(0,0,0,0.2)",
    backgroundColor : "#fff",
    minHeight : "250px",
    position : "relative",
  },
  a : {
    textDecoration : "none",
  },
  img : {
    paddingRight : "1rem",
    height : "auto",
    width : "auto",
    maxWidth : "120px",
    maxHeight : "120px",
  },
  frontSideDescription : {
    paddingTop : "1rem",
    paddingBottom : "1rem",
    textAlign : "left",
    display : "flex",
    flexDirection : 'row',
  },
  link : {
    textDecoration : "none",
    color : "#00b39F",
  },
});

const INITIAL_GRID_SIZE = { lg : 6, md : 12, xs : 12 };

const Extensions = ({ classes, toggleCatalogContent,  enqueueSnackbar, closeSnackbar }) => {
  const [catalogContent, setCatalogContent] = useState(true);
  const [extensionPreferences, setExtensionPreferences] = useState({})

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
          { variant : 'success',
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

  const handleSignUp = (e) => {
    window.open("https://layer5.io/meshmap", "_blank")
    e.stopPropagation();
  };


  return (
    <React.Fragment>
      <Head>
        <title>Extensions | Meshery</title>
      </Head>
      <Grid container spacing={1} >
        <Grid item {...INITIAL_GRID_SIZE}>
          <div className={classes.card}>
            <Typography className={classes.frontContent} variant="h5" component="div">
              {"MeshMap"}
            </Typography>

            <Typography className={classes.frontSideDescription} variant="body">
              <img className={classes.img} src="/static/img/meshmap.svg" />
              Collaboratively design and manage your Kubernetes clusters, service mesh deployments, and cloud native apps.
              MeshMap is now in private beta. Sign-up today to for early access!
            </Typography>
            <div style={{ textAlign : "right" }}>
              <Button
                variant="contained"
                color="primary"
                className={classes.button}
                onClick={(e) => handleSignUp(e)}>
                Sign Up
              </Button>
            </div>
          </div>
        </Grid>
        <Grid item {...INITIAL_GRID_SIZE}>

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
                />
              </div>
            </Grid>
          </div>
        </Grid>
      </Grid>
    </React.Fragment>
  )
}

const mapStateToProps = state => ({
  catalogVisibility : state.get('catalogVisibility')
})

const mapDispatchToProps = dispatch => ({
  toggleCatalogContent : bindActionCreators(toggleCatalogContent, dispatch)
})

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(Extensions)));