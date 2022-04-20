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
import Meshery from "../../img/SVGs/meshery";
import MesheryBlackIcon from "../../img/SVGs/mesheryBlackIcon";
import MesheryIcon from "../../img/meshery-logo/CustomMesheryLogo";
import CustomTypography from "../CustomTypography"
import { DockerMuiThemeProvider } from '@docker/docker-mui-theme';
import CssBaseline from '@mui/material/CssBaseline';
import { StyledDiv, AccountDiv, ServiceMeshAdapters, ExtensionWrapper, AdapterDiv, ComponentWrapper, SectionWrapper } from "./styledComponents";
import { MesheryAnimation } from "../MesheryAnimation/MesheryAnimation";


const ExtensionsComponent = props => {
  const [consulChecked, isConsulChecked] = useState(true);
  const [istioChecked, isIstioChecked] = useState(false);
  const [linkerdChecked, isLinkerdChecked] = useState(false);
  const [nginxChecked, isNginxChecked] = useState(false);
  const [kumaChecked, isKumaChecked] = useState(false);
  const [isHovered, setIsHovered] = useState(false);


  const theme = {
    textAlign: "center",
    backgroundColor: "#222C32",
    padding: "5rem",
    maxHeight: "100vh"
  };

  const onMouseOver = e => {
    let target = e.target.closest("div");
    target.style.transition = "all .5s";
    target.style.transform = "scale(1.11)";
  }
  const onMouseOut = e => {
    setIsHovered(!isHovered);
    let target = e.target.closest("div");
    target.style.transition = "all .8s";
    target.style.transform = "scale(1)";
  }

  const onClick = e => {
    let target = e.target.closest("div");
    target.style.transition = "all .2s";
    target.style.transform = "scale(0.9)";
    target.style.cursor = "grabbing";
    setTimeout(() => {
      target.style.transition = "none";
    }, 100);
  };

  useEffect(() => {
    // window.ddClient.extension.vm.service.get("/ping").then(console.log);
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
      <ComponentWrapper>
        <CssBaseline />
        <MesheryIcon />
        <CustomTypography sx={{ margin: "auto", paddingTop: "1rem" }}>Design and operate your cloud native deployments with the extensible management plane, Meshery.</CustomTypography>

        <SectionWrapper>

        <ExtensionWrapper>
            <AccountDiv>
              <CustomTypography sx={{ marginBottom: "1rem", whiteSpace: "nowrap" }}>
                Launch Meshery
              </CustomTypography>
              <div style={{marginBottom: "0.5rem"}}>
                <a style={{ textDecoration: "none" }} href="http://localhost:9081">
               
                <div
          onMouseEnter={() => setIsHovered(!isHovered)}
          onMouseLeave={onMouseOut}
          onClick={onClick}
          onMouseOver={onMouseOver}
        > 
         {isHovered ?    <MesheryAnimation height={70} width={72}/> : <Meshery height={70} width={72} /> } 
        
        </div>
        {/* <MesheryBlackIcon width={110} /> */}
             {/* <div onMouseEnter={trigger}>   <MeshSyncAnimation style={{ height : "75px" }} /> </div> */}
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
          <Grid style = {{display: "flex", justifyContent: 'center', alignItems: 'center' }}>
            <ServiceMeshAdapters>
             
              <StyledDiv>
                <AdapterDiv inactiveAdapter={!consulChecked}>
                  <ConsulIcon width={40} height={40} /> </AdapterDiv>
                <Switch onChange={handleConsul} color="primary" defaultChecked></Switch>
              </StyledDiv>
              <StyledDiv>
                <AdapterDiv inactiveAdapter={!istioChecked}>
                  <IstioIcon width={40} height={40} /></AdapterDiv>
                <Switch onChange={handleIstio} color="primary"></Switch> </StyledDiv>

              <StyledDiv>
                <AdapterDiv inactiveAdapter={!linkerdChecked}><LinkerdIcon width={40} height={40} /></AdapterDiv>
                <Switch onChange={handleLinkerd} color="primary"></Switch> </StyledDiv>
              <StyledDiv>
                <AdapterDiv inactiveAdapter={!nginxChecked}><NginxIcon width={38} height={40} /></AdapterDiv><Switch onChange={handleNginx} color="primary"></Switch> </StyledDiv>
              <StyledDiv>
                <AdapterDiv inactiveAdapter={!kumaChecked}><KumaIcon width={40} height={40} /></AdapterDiv><Switch onChange={handleKuma} color="primary"></Switch> </StyledDiv>
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
