import React, { useState, useEffect } from "react";
import {
  Grid,
  Button, Switch
} from "@mui/material";
import { styled } from "@mui/material/styles";
import ConsulIcon from "../../img/SVGs/consulIcon";
import IstioIcon from "../../img/SVGs/IstioIcon";
import KumaIcon from "../../img/SVGs/kumaIcon";
import LinkerdIcon from "../../img/SVGs/linkerdIcon";
import NginxIcon from "../../img/SVGs/nginxIcon";
import MesheryIcon from "../../img/meshery-logo/CustomMesheryLogo";
import CustomTypography from "../CustomTypography"
import { DockerMuiThemeProvider } from '@docker/docker-mui-theme';
import CssBaseline from '@mui/material/CssBaseline';
import dataFetch from "../../lib/data-fetch";
import { StyledDiv, AccountDiv, ServiceMeshAdapters, ExtensionWrapper, AdapterDiv, ComponentWrapper, SectionWrapper } from "./styledComponents";


const baseURL = "http://localhost:9081"

const ExtensionsComponent = () => {
  const [consulChecked, setConsulChecked] = useState(false);
  const [istioChecked, isIstioChecked] = useState(false);
  const [linkerdChecked, isLinkerdChecked] = useState(false);
  const [nginxChecked, isNginxChecked] = useState(false);
  const [kumaChecked, isKumaChecked] = useState(false);


  const theme = {
    textAlign: "center",
    backgroundColor: "#222C32",
    padding: "5rem",
    maxHeight: "100vh"
  };

  useEffect(() => {
    // window.ddClient.extension.vm.service.get("/ping").then(console.log);
  }, [])

  const submitConfig = (adapterLocation) => {

    const data = { meshLocationURL: adapterLocation };

    const params = Object.keys(data)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)

    console.log("asdf")
    window.ddClient.extension.vm.service.post("/api/system/adapter/manage", {
      credentials: "same-origin",
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8", },
      body: params,
    }).then(console.log).catch(console.log);
  }


  // Wrote separate functions since we need these functions to provision the adapters as well
  const handleConsul = () => {
    window.ddClient.desktopUI.toast.success(`Request received. ${consulChecked ? "Deprovisioning" : "Provisioning"} Consul Service Mesh...`);
    setTimeout(() => {
      window.ddClient.desktopUI.toast.success(`Consul Service Mesh ${consulChecked ? "Deprovisioned" : "Provisioned"} successfully`);
    }, 3000)
    setConsulChecked(prev => !prev)
  }
  const handleIstio = () => {
    window.ddClient.desktopUI.toast.success(`Request received. ${istioChecked ? "Deprovisioning" : "Provisioning"} Istio Service Mesh...`);
    setTimeout(() => {
      window.ddClient.desktopUI.toast.success(`Istio Service Mesh ${istioChecked ? "Deprovisioned" : "Provisioned"} successfully`);
    }, 3000)
    isIstioChecked(prev => !prev);
  }
  const handleLinkerd = () => {
    window.ddClient.desktopUI.toast.success(`Request received. ${linkerdChecked ? "Deprovisioning" : "Provisioning"} Linkerd Service Mesh...`);
    setTimeout(() => {
      window.ddClient.desktopUI.toast.success(`Linkerd Service Mesh ${linkerdChecked ? "Deprovisioned" : "Provisioned"} successfully`);
    }, 3000)
    isLinkerdChecked(prev => !prev);
  }
  const handleNginx = () => {
    window.ddClient.desktopUI.toast.success(`Request received. ${nginxChecked ? "Deprovisioning" : "Provisioning"} Nginx Service Mesh...`);
    setTimeout(() => {
      window.ddClient.desktopUI.toast.success(`Nginx Service Mesh ${nginxChecked ? "Deprovisioned" : "Provisioned"} successfully`);
    }, 3000)
    isNginxChecked(prev => !prev);
  }
  const handleKuma = () => {
    window.ddClient.desktopUI.toast.success(`Request received. ${kumaChecked ? "Deprovisioning" : "Provisioning"} Kuma Service Mesh...`);
    setTimeout(() => {
      window.ddClient.desktopUI.toast.success(`Kuma Service Mesh ${kumaChecked ? "Deprovisioned" : "Provisioned"} successfully`);
    }, 3000)
    isKumaChecked(prev => !prev);
  }


  return (
    <DockerMuiThemeProvider theme={theme}>
      <ComponentWrapper>
        <CssBaseline />
        <MesheryIcon />
        <CustomTypography sx={{ margin: "auto", paddingTop: "1rem" }}>Design and operate your cloud native deployments with the extensible management plane, Meshery.</CustomTypography>

        <SectionWrapper>
          <ExtensionWrapper>

            <AccountDiv>
              <CustomTypography sx={{ marginBottom: "2rem" }}>
                Account
              </CustomTypography>
              <div style={{ paddingBottom: "2rem" }}>
                <a style={{ textDecoration: "none" }} href="http://localhost:9081">
                  <Button sx={{ color: "#FFFFFF", whiteSpace: "nowrap" }} variant="contained">
                    Open Meshery
                  </Button>
                </a>
              </div>
            </AccountDiv>

          </ExtensionWrapper>

          <ExtensionWrapper>
            <AccountDiv>
              <CustomTypography sx={{ marginBottom: "2rem", whiteSpace: " nowrap" }}>Import Compose App</CustomTypography>
              <div style={{ paddingBottom: "2rem" }}>
                <label htmlFor="upload-button" >
                  <Button sx={{ backgroundColor: "#7794AB" }} variant="contained" color="primary" aria-label="Upload Button" component="span" >
                    <input id="upload-button" type="file" accept=".yaml, .yml" hidden name="upload-button" />
                    Browse...
                  </Button>
                </label>
              </div>
            </AccountDiv>
          </ExtensionWrapper>
          <ExtensionWrapper>
            <div>
              <CustomTypography sx={{ marginBottom: "2rem" }}>Deploy a Service Mesh</CustomTypography>
              <Grid style={{ display: "flex", justifyContent: 'center', alignItems: 'center' }}>
                <ServiceMeshAdapters>
                  <StyledDiv>
                    <AdapterDiv inactiveAdapter={!consulChecked}>
                      <ConsulIcon width={40} height={40} /> </AdapterDiv>
                    <Switch checked={consulChecked} onChange={handleConsul} color="primary" ></Switch>
                  </StyledDiv>
                  <StyledDiv>
                    <AdapterDiv inactiveAdapter={!istioChecked}>
                      <IstioIcon width={40} height={40} /></AdapterDiv>
                    <Switch checked={istioChecked} onChange={handleIstio} color="primary"></Switch> </StyledDiv>

                  <StyledDiv>
                    <AdapterDiv inactiveAdapter={!linkerdChecked}><LinkerdIcon width={40} height={40} /></AdapterDiv>
                    <Switch checked={linkerdChecked} onChange={handleLinkerd} color="primary"></Switch> </StyledDiv>
                  <StyledDiv>
                    <AdapterDiv inactiveAdapter={!nginxChecked}><NginxIcon width={38} height={40} /></AdapterDiv><Switch checked={nginxChecked} onChange={handleNginx} color="primary"></Switch> </StyledDiv>
                  <StyledDiv>
                    <AdapterDiv inactiveAdapter={!kumaChecked}><KumaIcon width={40} height={40} /></AdapterDiv><Switch checked={kumaChecked} onChange={handleKuma} color="primary"></Switch> </StyledDiv>
                </ServiceMeshAdapters>
              </Grid>
            </div>
          </ExtensionWrapper>
        </SectionWrapper>
      </ComponentWrapper>
    </DockerMuiThemeProvider >
  );
}

export default ExtensionsComponent;
