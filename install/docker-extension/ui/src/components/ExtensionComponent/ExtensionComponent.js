import React, { useState, useEffect } from "react";
import {
  Grid,
  Button, Switch, Typography
} from "@mui/material";
import ConsulIcon from "../../img/SVGs/consulIcon";
import IstioIcon from "../../img/SVGs/IstioIcon";
import KumaIcon from "../../img/SVGs/kumaIcon";
import Joyride from 'react-joyride';
import Tour from "../Walkthrough/Tour";
import { createTheme } from '@mui/material/styles';
import LinkerdIcon from "../../img/SVGs/linkerdIcon";
import NginxIcon from "../../img/SVGs/nginxIcon";
import OsmIcon from "../../img/SVGs/osmIcon";
import Meshery from "../../img/SVGs/meshery";
import MesheryIcon from "../../img/meshery-logo/CustomMesheryLogo";
import { DockerMuiThemeProvider } from '@docker/docker-mui-theme';
import CssBaseline from '@mui/material/CssBaseline';
import { StyledDiv, AccountDiv, ServiceMeshAdapters, ExtensionWrapper, AdapterDiv, ComponentWrapper, SectionWrapper } from "./styledComponents";
import { MesheryAnimation } from "../MesheryAnimation/MesheryAnimation";
import axios from "axios";


const baseURL = "http://localhost:9081"

export function trueRandom() {
  return crypto.getRandomValues(new Uint32Array(1))[0] / 2 ** 32;
}

export function randomPatternNameGenerator() {
  return "meshery_" + Math.floor(trueRandom() * 100)
}

const useThemeDetector = () => {
  const getCurrentTheme = () => window.matchMedia("(prefers-color-scheme: dark)").matches;
  const [isDarkTheme, setIsDarkTheme] = useState(getCurrentTheme());
  const mqListener = (e => {
    setIsDarkTheme(e.matches);
  });

  useEffect(() => {
    const darkThemeMq = window.matchMedia("(prefers-color-scheme: dark)");
    darkThemeMq.addListener(mqListener);
    return () => darkThemeMq.removeListener(mqListener);
  }, []);
  return isDarkTheme;
}

const ExtensionsComponent = (props) => {
  const {steps} = props;
  const [consulChecked, setConsulChecked] = useState(false);
  const [istioChecked, isIstioChecked] = useState(false);
  const [linkerdChecked, isLinkerdChecked] = useState(false);
  const [nginxChecked, isNginxChecked] = useState(false);
  const [kumaChecked, isKumaChecked] = useState(false);
  const [isHovered, setIsHovered] = useState(false);
  const isDarkTheme = useThemeDetector();
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [userName, setUserName] = useState("")
  const [token, setToken] = useState()

  useEffect(() => {
        let socket = new WebSocket("ws://127.0.0.1:7877/ws");
        console.log("Attempting Connection...");

        socket.onopen = () => {
            console.log("Successfully Connected");
            socket.send("Hi From the Client!")
        };

        socket.onclose = event => {
            console.log("Socket Closed Connection: ", event);
            socket.send("Client Closed!")
        };

        socket.onmessage = msg  => {
          console.log("From proxy ws connection: ", msg)
          if(msg.data == "Authenticated")
          setIsLoggedIn(true)
    }

        socket.onerror = error => {
            console.log("Socket Error: ", error);
        };
  }, [])

  useEffect(() => {
    fetch("http://127.0.0.1:7877/token").then(res => res.text()).then(res => {
      setToken(res)
    }).catch(console.log)
  }, [])

  useEffect(() => {
    fetch("http://127.0.0.1:7877/token/store")
      .then((obj) => {
        console.log(obj)
        if (obj.status >= 200 && obj.status < 300) {
          setIsLoggedIn(true)
          fetch("http://127.0.0.1:7877/api/user").then(res => res.text()).then(res => setUserName(JSON.parse(res)?.user_id))
        }
        else setIsLoggedIn(false)
      })
      .catch((obj) => {
        setIsLoggedIn(false)
        console.log(obj)
      })
  }, [])

  const onMouseOver = e => {
    let target = e.target.closest("div");
    target.style.transition = "all .5s";
    target.style.transform = "scale(1)";
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
    target.style.transform = "scale(0.8)";
    // fetch("http://127.0.0.1:7877/token").then(res => res.text()).then(res => {
    //   console.log(res)
    //   window.ddClient.host.openExternal("http://localhost:7877/api/user/token?token=" + res)
    // }).catch(console.log)
  };



  const submitConfig = (adapterLocation) => {
    const data = { meshLocationURL: adapterLocation };

    const params = Object.keys(data)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)

    fetch("http://127.0.0.1:7877/api/system/adapter/manage", {
      credentials: "same-origin",
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8", },
      mode: "no-cors",
      body: params,
    }).then(() => {
      window.ddClient.desktopUI.toast.success("Service Mesh was successfully provisioned.")
    }).catch(() => window.ddClient.desktopUI.toast.error("Some error occured while trying to provision the service mesh."));
  }


  // Wrote separate functions since we need these functions to provision the adapters as well
  const handleConsul = () => {
    // window.ddClient.desktopUI.toast.success(`Request received. ${consulChecked ? "Deprovisioning" : "Provisioning"} Consul Service Mesh...`);
    // setTimeout(() => {
    //   window.ddClient.desktopUI.toast.success(`Consul Service Mesh ${consulChecked ? "Deprovisioned" : "Provisioned"} successfully`);
    // }, 3000)
    submitConfig("localhost:10002")
    setConsulChecked(prev => !prev)
  }
  const handleIstio = () => {
    // window.ddClient.desktopUI.toast.success(`Request received. ${istioChecked ? "Deprovisioning" : "Provisioning"} Istio Service Mesh...`);
    // setTimeout(() => {
    //   window.ddClient.desktopUI.toast.success(`Istio Service Mesh ${istioChecked ? "Deprovisioned" : "Provisioned"} successfully`);
    // }, 3000)
    submitConfig("localhost:10000")
    isIstioChecked(prev => !prev);
  }
  const handleLinkerd = () => {
    // window.ddClient.desktopUI.toast.success(`Request received. ${linkerdChecked ? "Deprovisioning" : "Provisioning"} Linkerd Service Mesh...`);
    // setTimeout(() => {
    //   window.ddClient.desktopUI.toast.success(`Linkerd Service Mesh ${linkerdChecked ? "Deprovisioned" : "Provisioned"} successfully`);
    // }, 3000)
    submitConfig("localhost:10001")
    isLinkerdChecked(prev => !prev);
  }
  const handleNginx = () => {
    // window.ddClient.desktopUI.toast.success(`Request received. ${nginxChecked ? "Deprovisioning" : "Provisioning"} Nginx Service Mesh...`);
    // setTimeout(() => {
    //   window.ddClient.desktopUI.toast.success(`Nginx Service Mesh ${nginxChecked ? "Deprovisioned" : "Provisioned"} successfully`);
    // }, 3000)
    // submitConfig("localhost:10001")
    isNginxChecked(prev => !prev);
  }
  const handleKuma = () => {
    // window.ddClient.desktopUI.toast.success(`Request received. ${kumaChecked ? "Deprovisioning" : "Provisioning"} Kuma Service Mesh...`);
    // setTimeout(() => {
    //   window.ddClient.desktopUI.toast.success(`Kuma Service Mesh ${kumaChecked ? "Deprovisioned" : "Provisioned"} successfully`);
    // }, 3000)
    submitConfig("localhost:10007")
    isKumaChecked(prev => !prev);
  }



  const handleImport = () => {
    const file = document.getElementById("upload-button").files[0];
    // Create a reader
    const reader = new FileReader();
    reader.addEventListener("load", (event) => {

      let body = { save: true }
      let name = randomPatternNameGenerator()
      body = JSON.stringify({
        ...body, application_data: { name, application_file: event.target.result }
      })

      fetch("http://localhost:7877/api/application", {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8", },
        body,
      }).then((res) => {
        console.log(res)
        window.ddClient.desktopUI.toast.success("Compose file has been uploaded with name: " + name)
      }).catch(() => window.ddClient.desktopUI.toast.error("Some error occured while uploading the compose file."));

    });
    reader.readAsText(file);

    // window.ddClient.desktopUI.toast.success(`Importing Compose App...`);
    // setTimeout(() => {
    //   window.ddClient.desktopUI.toast.success(`Compose App imported successfully`);
    // }, 3000)
  }

  return (
    <DockerMuiThemeProvider>
      <CssBaseline />
      <Tour />
      <ComponentWrapper>
        
        <MesheryIcon CustomColor={isDarkTheme ? "white" : "#3C494F"} />
        <Typography sx={{ margin: "auto", paddingTop: "1rem" }}>Design and operate your cloud native deployments with the extensible management plane, Meshery.</Typography>

        <SectionWrapper>

          <ExtensionWrapper className="third-step" sx={{ backgroundColor: isDarkTheme ? "#393F49" : "#D7DADE" }}>
            <AccountDiv>
              <Typography sx={{ marginBottom: "1rem", whiteSpace: "nowrap" }}>
                Launch Meshery
              </Typography>
              <div style={{ marginBottom: "0.5rem" }}>
                <a style={{ textDecoration: "none" }} href={token && "http://localhost:7877/api/user/token?token=" + token} >

                  <div
                    onMouseEnter={() => setIsHovered(!isHovered)}
                    onMouseLeave={onMouseOut}
                    onClick={onClick}
                    onMouseOver={onMouseOver}
                  >
                    {isHovered ? <MesheryAnimation height={70} width={72} /> : <Meshery height={70} width={72} />}

                  </div>
                </a>
              </div>

              {userName && <Typography sx={{ marginBottom: "1rem", whiteSpace: "nowrap" }}>
                User: {userName}
              </Typography>}
            </AccountDiv>
          </ExtensionWrapper>


          {isLoggedIn && <ExtensionWrapper sx={{ backgroundColor: isDarkTheme ? "#393F49" : "#a5b1ba" }}>
            <AccountDiv>
              <Typography sx={{ marginBottom: "2rem", whiteSpace: " nowrap" }}>Import Compose App</Typography>
              <div style={{ paddingBottom: "2rem" }}>
                <label htmlFor="upload-button" >
                  <Button variant="contained" color="primary" disabled={!isLoggedIn} aria-label="Upload Button" component="span" >
                    <input id="upload-button" type="file" accept=".yaml, .yml" hidden name="upload-button" onChange={handleImport} />
                    Browse...
                  </Button>
                </label>
              </div>
            </AccountDiv>
          </ExtensionWrapper>}

          {!isLoggedIn && <ExtensionWrapper sx={{ backgroundColor: isDarkTheme ? "#393F49" : "#a5b1ba" }}>
            <AccountDiv style={{ width: "" }}>
              <Button variant="contained" disabled={isLoggedIn} color="primary" component="span" onClick={() => {
                window.ddClient.host.openExternal("https://meshery.layer5.io?source=aHR0cDovL2xvY2FsaG9zdDo3ODc3L3Rva2VuL3N0b3Jl&provider_version=v0.3.14")
              }}>
                Login
              </Button>
            </AccountDiv>
          </ExtensionWrapper>}

          {!!isLoggedIn && <ExtensionWrapper sx={{ backgroundColor: isDarkTheme ? "#393F49" : "#a5b1ba" }} >
            <div>
              <Typography sx={{ marginBottom: "1rem" }}>Deploy a Service Mesh</Typography>
              <Grid style={{ display: "flex", justifyContent: 'center', alignItems: 'center' }}>
                <ServiceMeshAdapters>
                  <StyledDiv>
                    <AdapterDiv inactiveAdapter={!consulChecked}>
                      <ConsulIcon width={40} height={40} /> </AdapterDiv>
                    <Typography>Consul</Typography>
                    <Switch checked={consulChecked} disabled={!isLoggedIn} onChange={handleConsul} color="primary" ></Switch>
                  </StyledDiv>
                  <StyledDiv>
                    <AdapterDiv inactiveAdapter={!istioChecked}>
                      <IstioIcon width={40} height={40} /></AdapterDiv>
                    <Typography >Istio</Typography>
                    <Switch checked={istioChecked} disabled={!isLoggedIn} onChange={handleIstio} color="primary"></Switch> </StyledDiv>

                  <StyledDiv>
                    <AdapterDiv inactiveAdapter={!linkerdChecked}><LinkerdIcon width={40} height={40} /></AdapterDiv>
                    <Typography>Linkerd</Typography>
                    <Switch checked={linkerdChecked} disabled={!isLoggedIn} onChange={handleLinkerd} color="primary"></Switch> </StyledDiv>
                  <StyledDiv>
                    <AdapterDiv inactiveAdapter={!nginxChecked}><NginxIcon width={38} height={40} /></AdapterDiv>
                    <Typography>Nginx</Typography>
                    <Switch checked={nginxChecked} disabled={!isLoggedIn} onChange={handleNginx} color="primary"></Switch> </StyledDiv>
                  <StyledDiv>
                    <AdapterDiv inactiveAdapter={!kumaChecked}><KumaIcon width={40} height={40} /></AdapterDiv>
                    <Typography>Kuma</Typography>
                    <Switch checked={kumaChecked} disabled={!isLoggedIn} onChange={handleKuma} color="primary"></Switch> </StyledDiv>
                </ServiceMeshAdapters>
              </Grid>
            </div>
          </ExtensionWrapper>}
        </SectionWrapper>
      </ComponentWrapper>
    </DockerMuiThemeProvider>
  );
}

export default ExtensionsComponent;
