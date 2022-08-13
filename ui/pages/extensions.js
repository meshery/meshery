import { Grid, Typography, Button, Switch } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import { useState } from "react";


const styles = (theme) => ({
  button : {
    marginRight : "0.5rem",
    borderRadius : 5,
    minWidth : 100,
    color : "#fff",
    "&:hover" : {
      boxShadow : "0px 2px 4px -1px rgb(0 0 0 / 20%), 0px 4px 5px 0px rgb(0 0 0 / 14%), 0px 1px 10px 0px rgb(0 0 0 / 12%)"
    },
  },
  card : {
    padding : theme.spacing(2),
    borderRadius : theme.spacing(1),
    transformStyle : "preserve-3d",
    boxShadow : "0 4px 8px 0 rgba(0,0,0,0.2)",
    backgroundColor : "#fff",
    minHeight : "100%",
    position : "relative"
  },
  a : {
    textDecoration : "none",
    padding : "0.5rem"
  },
  img : {
    paddingRight : "1rem",
    height : "auto",
    width : "auto",
    maxWidth : "120px",
    maxHeight : "120px",
  },
  frontSideDescription : {
    padding : "1rem",
    textAlign : "left",
    display : "flex",
    flexDirection : 'row',
    justifyContent : 'flex-start',
    alignItems : 'flex-start',
    textAlign : "left"
  },
  link : {
    textDecoration : "none",
    display : "flex",
    justifyContent : "flex-end",
    padding : "0.5rem"
  },
  comingSoon : {
    transform : "rotate(30deg) translateX(20%) translateY(-50%)",
    position : "absolute",
    right : "0",
  }
});

const INITIAL_GRID_SIZE = { lg : 6,  md : 12, xs : 12 };

const Extensions = ({ classes }) => {
  const [catalog, setCatalog] = useState(false);

  const handleToggle = (e) => {
    e.stopPropagation();
    console.log(`Meshery Catalog: ${catalog ? "enabled" : "disabled"}`)
    setCatalog(catalog => !catalog)
  }

  const handleSignUp = (e) => {
    window.open("https://layer5.io/meshmap", "_blank")
    e.stopPropagation();
  };


  return (
    <Grid container spacing={1} >
      <Grid item {...INITIAL_GRID_SIZE} className={classes.grid} >
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
      <Grid item {...INITIAL_GRID_SIZE} className={classes.grid}>

        <div className={classes.card}>
          <Typography className={classes.comingSoon} variant="h5" component="div">Coming Soon</Typography>
          <Typography className={classes.frontContent} variant="h5" component="div">
            {"Meshery Catalog"}
          </Typography>

          <Typography className={classes.frontSideDescription} variant="body">
            <img className={classes.img} src="/static/img/meshery_catalog.svg" />
              Connect to enable access to the cloud native catalog. Import any catalog item and customize.


          </Typography>

          <Typography className={classes.link} variant="subtitle2" style={{ fontStyle : "italic" }}>
              Explore at:<a href="https://meshery.io/catalog" target="_blank" rel="noreferrer" style={{ textDecoration : "none", paddingLeft : "0.3rem ", color : "#00b39F" }}> Meshery Catalog </a>
          </Typography>
          <div style={{ textAlign : "right" }}>
            <Switch
              disabled
              checked={catalog}
              onClick={(e) => handleToggle(e)}
              name="OperatorSwitch"
              color="primary"
            />
          </div>
        </div>
      </Grid>
    </Grid>
  )
}

export default withStyles(styles)(Extensions);