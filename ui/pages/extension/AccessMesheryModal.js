import React, {
} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Typography, Grid, Button, IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';

const styles = makeStyles((theme) => ({
  paper : {
    position : 'absolute',
    width : "30rem",
    height : "20rem",
    backgroundColor : theme.palette.background.paper,
    boxShadow : theme.shadows[5],
    borderRadius : "4px",
    top : "0",
    bottom : "0",
    left : "0",
    right : "0",
    margin : "auto",
    ["@media (max-width: 455px)"] : {
      width : "100%"
    },
    zIndex : 5
  },
  header : {
    display : "flex",
    justifyContent : "center",
    position : "absolute",
    width : "100%",
    top : "0.7rem",
    fontWeight : "bold",
    ["@media (max-width: 455px)"] : {
      fontSize : "1rem"
    },
  },
  imgWrapper : {
    display : "flex",
    justifyContent : "center"
  },
  headerWrapper : {
    marginBottom : 12,
    borderRadius : "4px 4px 0 0",
    textAlign : 'center',
    width : "100%",
    padding : theme.spacing(1),
    color : '#fff',
    backgroundColor : "#252E31",
    fontSize : "1rem",
  },
  closing : {
    transform : "rotate(-90deg)",
    "&:hover" : {
      transform : "rotate(90deg)",
      transition : "all .3s ease-in",
      cursor : "pointer",
    }
  }

}));

function PlaygroundMeshDeploy(props) {

  if (props.isOpen) {
    return <PlayCard closeForm={props.closeForm} rootStyle={props.rootStyle} />
  }
  return <></>;
}

export function PlayCard({ rootStyle, closeForm }) {
  const classes = styles();

  const handlePage = (e) => {
    window.open("https://meshery.io/#getting-started", "_blank")
    e.stopPropagation();
  };

  return (
    <div
      className={classes.paper}
      style={rootStyle}
    >
      <div className={classes.headerWrapper}>
        <Typography className={classes.header} variant="h6">Access all features</Typography>
        <div style={{ display : "flex", justifyContent : "flex-end", whiteSpace : "nowrap", position : "relative" }}>
          <IconButton key="close" aria-label="Close" color="inherit" onClick={closeForm}>
            <CloseIcon className={classes.closing} />
          </IconButton>
        </div>
      </div>
      <div className={classes.imgWrapper}>
        <img width="60%" height="60%" src="/static/img/meshery-logo-light-text-side.png" />
      </div>
      <div style={{ display : "flex", justifyContent : "center", padding : "1.5rem", width : "100%" }} >
        <Typography>Access is limited in Meshery <span style={{ color : "#F0A303" }}>Playground.</span> <br />
                    To get full access, Deploy Meshery locally.</Typography>
      </div>
      <Typography style={{ textAlign : "center", marginBottom : "1rem" }}>&#x26A1; Deploy Meshery in seconds</Typography>
      <div style={{ display : "flex", justifyContent : "center" }}>
        <Grid item xs={3}>
          <Button fullWidth variant="outlined" color="primary" onClick={(e) => handlePage(e)}>Get Started</Button>
        </Grid>
      </div>
    </div>
  )
}

export default PlaygroundMeshDeploy;
