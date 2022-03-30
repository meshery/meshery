import React, { useState, useEffect } from "react";
import { makeStyles } from "@mui/styles";
import {
  Grid,
  Button, Switch
} from "@mui/material";
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";
import ConsulIcon from "../../img/SVGs/consulIcon";
import IstioIcon from "../../img/SVGs/IstioIcon";
import KumaIcon from "../../img/SVGs/kumaIcon";
import LinkerdIcon from "../../img/SVGs/linkerdIcon";
import NginxIcon from "../../img/SVGs/nginxIcon";
import MesheryIcon from "../../img/meshery-logo/CustomMesheryLogo";
import { CustomTypography } from "../CustomTypography";
import {useStyles} from "./ExtensionComponent.styles";
import { createTheme, ThemeProvider } from "@mui/material";
import { DockerMuiThemeProvider } from '@docker/docker-mui-theme';
import CssBaseline from '@mui/material/CssBaseline';
const ExtensionsComponent = props => {
  const [consulChecked, isConsulChecked] = useState(true);
  const [istioChecked, isIstioChecked] = useState(false);
  const [linkerdChecked, isLinkerdChecked] = useState(false);
  const [nginxChecked, isNginxChecked] = useState(false);
  const [kumaChecked, isKumaChecked] = useState(false);
  const classes = useStyles();

  const theme = createTheme();

  useEffect(() => {
    window.ddClient.extension.vm.service.get("/ping").then(console.log);
  }, [])

  // Wrote separate functions since we need these functions to provision the adapters as well
  const handleConsul = () => {
    isConsulChecked(prev => !prev);
  }
  const handleIstio = () => {
    isIstioChecked(prev => !prev);
  }
  const handleLinkerd = () => {
    isLinkerdChecked(prev => !prev);
  }
  const handleNginx = () => {
    isNginxChecked(prev => !prev);
  }
  const handleKuma = () => {
    isKumaChecked(prev => !prev);
  }


  return (
    <DockerMuiThemeProvider theme={theme}>
      <CssBaseline />
    <div className={classes.root}>
      <MesheryIcon />
      <CustomTypography className={classes.headText}>Design and operate your cloud native deployments with the extensible management plane, Meshery.</CustomTypography>
      <div className={classes.main}>
        <CustomTypography variant="h6" className={classes.subText}>
                        CONFIGURE YOUR MESHERY DEPLOYMENT
        </CustomTypography>
        <div className={classes.OAuth}>
          <div className={classes.account}>
            <CustomTypography style="margin-bottom:2rem">Account</CustomTypography>
            <div><a className={classes.link} href="http://localhost:9081"><Button className={classes.mesheryConfig} variant="contained">
              Open Meshery
            </Button></a></div>
          </div>
          <Grid justify="center">
            <div className={classes.serviceMeshAdapters}>
              <CustomTypography style="margin-bottom:2rem">Deploy a Service Mesh</CustomTypography>
              <div className={classes.sm}>
                <div className={consulChecked ? null : classes.inactiveAdapter}>

                  <ConsulIcon width={40} height={40} /> </div>
                <Switch onChange={handleConsul} color="primary" defaultChecked></Switch>
              </div>
              <div className={classes.sm}>
                <div className={istioChecked ? null : classes.inactiveAdapter}>
                  <IstioIcon width={40} height={40} /></div>
                <Switch onChange={handleIstio} color="primary"></Switch></div>

              <div className={classes.sm}>
                <div className={linkerdChecked ? null : classes.inactiveAdapter}><LinkerdIcon width={40} height={40} /></div>
                <Switch onChange={handleLinkerd} color="primary"></Switch></div>
              <div className={classes.sm}>
                <div className={nginxChecked ? null : classes.inactiveAdapter}><NginxIcon width={38} height={40} /></div><Switch onChange={handleNginx} color="primary"></Switch></div>
              <div className={classes.sm}>
                <div className={kumaChecked ? null : classes.inactiveAdapter}><KumaIcon width={40} height={40} /></div><Switch onChange={handleKuma} color="primary"></Switch></div>
            </div>
          </Grid>
        </div>
      </div>
    </div>
    </DockerMuiThemeProvider>
  );
}

export default ExtensionsComponent;
