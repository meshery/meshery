import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import { NoSsr, Typography, Grid, Link } from "@material-ui/core";
import FavoriteIcon from '@material-ui/icons/Favorite';
import MesheryLogo from "../img/meshery-logo/meshery-logo.svg";
import { makeStyles } from "@material-ui/core/styles";
import {
  Container,
  Button, Paper
} from "@material-ui/core";
import TwitterIcon from '@material-ui/icons/Twitter';
import GitHubIcon from '@material-ui/icons/GitHub';
import LinkedInIcon from '@material-ui/icons/LinkedIn';
import GoogleIcon from "../img/SVGs/googleIcon";
import MesheryIcon from "../img/meshery-logo/CustomMesheryLogo";


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
    borderRadius: "5px ",
    padding: "1rem"
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
  button: {
    padding: "0.5rem"
  },

  Icon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
}));

const ExtensionsComponent = props => {

  const classes = useStyles();

  return (
    <div className={classes.root}>

      <div className={classes.main}>
        <MesheryIcon />

        <div className={classes.OAuth}>
          <Grid justify="center">
            <Grid >
              <div className={classes.button}>
                <Button style={{ backgroundColor: "#55acee" }}
                  className={classes.paper}>
                  <TwitterIcon className={classes.Icon} />
                  SIGN IN WITH TWITTER
                </Button>
              </div>
            </Grid>
            <Grid >
              <div className={classes.button}>
                <Button style={{ backgroundColor: "#444" }} className={classes.paper}>
                  <GitHubIcon className={classes.Icon} />
                  SIGN IN WITH GITHUB
                </Button>
              </div>
            </Grid>
            <Grid item xs>
              <div className={classes.button}>
                <Button style={{ backgroundColor: "#DD4B39" }} className={classes.paper}>
                  <div className={classes.Icon}> <GoogleIcon width={18} height={22} /></div>
                  SIGN IN WITH GOOGLE
                </Button>
              </div>
            </Grid>
            <Grid item xs>
              <div className={classes.button}>
                <Button style={{ backgroundColor: "#007bb6" }} className={classes.paper}>
                  <LinkedInIcon className={classes.Icon} />
                  SIGN IN WITH LINKEDIN
                </Button>
              </div>
            </Grid>
          </Grid>
        </div>
      </div>
    </div>
  );
}

export default ExtensionsComponent;
