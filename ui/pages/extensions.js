import { Grid, Typography, Button, Switch } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import { useState } from "react";
import FlipCard from "../components/FlipCard";

const styles = (theme) => ({
  card : {
    height : '100%',
    width : "40%",
    backgroundColor : "transparent",
    perspective : theme.spacing(125),
    margin : theme.spacing(5),
    overFlow : "hidden",
  },
  innerCard : {
    padding : theme.spacing(2),
    borderRadius : theme.spacing(1),
    transformStyle : "preserve-3d",
    boxShadow : "0 4px 8px 0 rgba(0,0,0,0.2)",
    backgroundColor : "#fff",
    cursor : "pointer",
    transition : `transform 1ms`,
    transformOrigin : "50% 50% 10%",

  },
  button : {
    marginRight : "0.5rem",
    borderRadius : 5,
    minWidth : 100,
    color : "#fff",
    "&:hover" : {
      boxShadow : "0px 2px 4px -1px rgb(0 0 0 / 20%), 0px 4px 5px 0px rgb(0 0 0 / 14%), 0px 1px 10px 0px rgb(0 0 0 / 12%)"
    },
  },
  grid : {
    margin : "1rem"
  },
  text : {
    padding : "1rem"
  },
  img : {
    marginLeft : "0.5rem",
    height : 'auto',
    width : "auto",
    maxWidth : "100px",
    maxHeight : "100px"
  },
  frontContent : {
    textAlign : "right"
  },
  backContent : {
    textAlign : "left"
  }
});

const INITIAL_GRID_SIZE = { xl : 4, md : 8, xs : 12 };

const Extensions = ({ classes }) => {
  const [catalog, setCatalog] = useState(false);

  const handleToggle = (e) => {
    e.stopPropagation();
    console.log(`Meshery Catalog: ${catalog ? "enabled" : "disabled"}`)
    setCatalog(catalog => !catalog)
  }

  const handleSignUp = (e) => {
    window.open("https://layer5.io/cloud-native-management/meshmap", "_blank")
    e.stopPropagation();
  };


  return (
    <Grid container justifyContent="center" alignItems="center">
      <Grid item {...INITIAL_GRID_SIZE} className={classes.grid} >
        <FlipCard >
          <div>
            <Typography className={classes.frontContent} variant="h5" component="div">
              {"MeshMap"}
            </Typography>

            <img className={classes.img} src="/static/img/meshmap.svg" />
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

          {/* Back part */}
          <div>
            <Typography className={classes.backContent} variant="h5" component="div">
              {"MeshMap"}
            </Typography>
            <Typography variant="subtitle1" component="div" className={classes.text}>
              {"Collaboratively design and manage your Kubernetes clusters, service mesh deployments, and cloud native apps."}
            </Typography>
          </div>
        </FlipCard>
      </Grid>
      <Grid item {...INITIAL_GRID_SIZE} className={classes.grid}>
        <FlipCard >
          <div>
            <Typography  className={classes.frontContent} variant="h5" component="div">
              {"Meshery Catalog"}
            </Typography>
            <img className={classes.img} src="/static/img/meshery_catalog.svg" />
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

          {/* Back part */}
          <div>
            <Typography className={classes.backContent} variant="h5" component="div">
              {"Meshery Catalog"}
            </Typography>
            <Typography variant="subtitle1" component="div" className={classes.text}>
              {"Connect to enable access to the cloud native catalog. Import any catalog item and customize."}
            </Typography>
            <Typography variant="subtitle2" style={{ fontStyle : "italic" }}>
                Explore at: <a href="https://meshery.io/catalog" target="_blank" rel="noreferrer"> Meshery Catalog </a>
            </Typography>
          </div>
        </FlipCard>
      </Grid>
    </Grid>
  )
}

export default withStyles(styles)(Extensions);