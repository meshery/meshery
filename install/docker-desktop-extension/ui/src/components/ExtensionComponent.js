import React, {useState} from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import { NoSsr, Typography, Grid, Link } from "@material-ui/core";
import FavoriteIcon from '@material-ui/icons/Favorite';
import MesheryLogo from "../img/meshery-logo/meshery-logo.svg";
import { makeStyles } from "@material-ui/core/styles";
import {
  Container,
  Button, Paper, Switch
} from "@material-ui/core";
import ConsulIcon from "../img/SVGs/consulIcon";
import IstioIcon from "../img/SVGs/IstioIcon";
import KumaIcon from "../img/SVGs/kumaIcon";
import LinkerdIcon from "../img/SVGs/linkerdIcon";
import NginxIcon from "../img/SVGs/nginxIcon";
import MesheryIcon from "../img/meshery-logo/CustomMesheryLogo";
import { CustomTypography } from "./CustomTypography";

const useStyles = makeStyles((theme) => ({
  root: {
    textAlign: "center",
    backgroundColor: "#222C32",
    padding: "5rem",
    minHeight: "100vh"
  },
  main: {
    margin: theme.spacing(5),
    backgroundColor: "#393F49",
    borderRadius: "20px ",
    padding: "1rem",
    height: "300px",
    [theme.breakpoints.down("xs")]: {
      height: "400px",
    },
  },
  paper: {
    padding: theme.spacing(1.5),
    textAlign: "center",
    color: "#ffffff",
    width: "240px",
    height: "45px"
  },
  OAuth: {
    padding: "2rem",

  },
  serviceMeshAdapters: {
    width: "50%",
    float: "right",
    [theme.breakpoints.down("xs")]: {
      width: "100%"
    }
  },
  account: {
    width: "50%",
    float: "left",
  },
  sm: {
    width: "16%",
    float: "left",
    flexDirection: "row",
    padding: "0.3rem"
  },
  subText: {
    color: "#AAAAAA",
    padding: "0.7rem"
  },
  headText: {
    maxWidth: "60%",
    margin: "auto",
    padding: "1rem"
  },
  button: {
    padding: "0.5rem"
  },
  inactiveAdapter: {
    filter: "grayscale(1) invert(0.35)"
  },
  mesheryConfig: {
    backgroundColor: "#7794AB",
    color: "#FFFFFF",
  },
  Icon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
}));

const ExtensionsComponent = props => {
  const [consulChecked, isConsulChecked] = useState(true);
  const [istioChecked, isIstioChecked] = useState(false);
  const [linkerdChecked, isLinkerdChecked] = useState(false);
  const [nginxChecked, isNginxChecked] = useState(false);
  const [kumaChecked, isKumaChecked] = useState(false);
  const classes = useStyles();

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
    <div className={classes.root}>
    <MesheryIcon />
    <CustomTypography className={classes.headText}>Design and operate your cloud native deployments with the extensible management plane, Meshery.</CustomTypography>
      <div className={classes.main}>
    
      <CustomTypography variant="h6" className={classes.subText}>
                        CONFIGURE YOUR MESHERY DEPLOYMENT
                    </CustomTypography>
        <div className={classes.OAuth}>
          <div className={classes.account}>
          <CustomTypography Style="margin-bottom:2rem">Account</CustomTypography>
          <div><Button className={classes.mesheryConfig} variant="contained">Open Meshery</Button></div>
          </div>
          <Grid justify="center">
            <div className={classes.serviceMeshAdapters}> 
            <CustomTypography Style="margin-bottom:2rem">Deploy a Service Mesh</CustomTypography>
            <div className={classes.sm}>
            <div className={consulChecked ?   null : classes.inactiveAdapter }>
              
        <ConsulIcon width={40} height={40}  /> </div>
          <Switch onChange={handleConsul} color="primary" defaultChecked></Switch>
          </div>
          <div className={classes.sm}>
          <div className={istioChecked ? null : classes.inactiveAdapter }>
            <IstioIcon width={40} height={40} /></div>
          <Switch onChange={handleIstio} color="primary"></Switch></div>
           
            <div className={classes.sm}>
            <div className={linkerdChecked ? null : classes.inactiveAdapter }><LinkerdIcon width={40} height={40} /></div>
              <Switch onChange={handleLinkerd} color="primary"></Switch></div>
              <div className={classes.sm}>
            <div className={nginxChecked ? null : classes.inactiveAdapter }><NginxIcon width={38} height={40} /></div><Switch onChange={handleNginx} color="primary"></Switch></div>
            <div className={classes.sm}>
            <div className={kumaChecked ? null : classes.inactiveAdapter }><KumaIcon width={40} height={40} /></div><Switch onChange={handleKuma} color="primary"></Switch></div>
             </div> 
          </Grid>
        </div>
      </div>
    </div>
  );
}

export default ExtensionsComponent;
