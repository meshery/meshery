import { Grid, Typography, Button, Switch } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import { useState } from "react";
import FlipCard from "../components/FlipCard";

const styles = (theme) => ({
  card : {
    height : '100%',
    width : "30%",
    backgroundColor : "transparent",
    perspective : theme.spacing(15),
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
    display : 'inline-flex',
    flexDirection : 'row',
    justifyContent : 'flex-start',
    alignItems : 'flex-start',

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
    margin : "1rem",
    display : 'inline-flex',
    flexDirection : 'row',
    justifyContent : 'flex-start',
    alignItems : 'flex-start',
  },
  a : {
    textDecoration : "none",
    marginTop : ".5rem",
  },
  img : {
    marginLeft : "0.5rem",
    marginTop : "1rem",
    height : "auto",
    width : "auto",
    maxWidth : "100px",
    maxHeight : "100px",
    display : "inline-flex",
    flexDirection : 'row',
    justifyContent : 'flex-start',
    alignItems : 'flex-start',
    textAlign : "left"

  },
  frontSideDescription : {
    textAlign : "left",
    display : "inline-flex",
    flexDirection : 'row',
    justifyContent : 'flex-start',
    alignItems : 'flex-start',
    textAlign : "left"
  },
  backContent : {
    marginBottom : ".75rem",
  },
  backSideDescription : {
    textAlign : "left",
    padding : "1rem",
    borderRadius : ".25rem",
    backgroundColor : "#eeeeee",
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
    window.open("https://layer5.io/meshmap", "_blank")
    e.stopPropagation();
  };


  return (
    <Grid container display="inline-flex" >
      <Grid item {...INITIAL_GRID_SIZE} display="inline-flex" className={classes.grid} >
        <FlipCard >
          <div>
            <Typography className={classes.frontContent} variant="h5" component="div">
              {"MeshMap"}
            </Typography>

            <img className={classes.img} src="/static/img/meshmap.svg" />
            <Typography className={classes.frontSideDescription} variant="body">
              {"Collaboratively design and manage your Kubernetes clusters, service mesh deployments, and cloud native apps."}
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

          {/* Back part */}
          <div>
            <Typography variant="h5" component="div">
              {"MeshMap"}
            </Typography>
            <Typography className={classes.backSideDescription} variant="body" component="div">
              MeshMap is now in private beta. Sign-up today to for early access!
            </Typography>

          </div>
        </FlipCard>
      </Grid>
      <Grid item {...INITIAL_GRID_SIZE} display="inline-flex" className={classes.grid}>
        <FlipCard >
          <div>
            <Typography className={classes.frontContent} variant="h5" component="div">
              {"Meshery Catalog"}
            </Typography>
            <img className={classes.img} src="/static/img/meshery_catalog.svg" />
            <Typography variant="body">
              {"Collaboratively design and manage your Kubernetes clusters, service mesh deployments, and cloud native apps."}
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

          {/* Back part */}
          <div className={classes.backContent}>
            <Typography className={classes.backContent} variant="h5" component="div">
              {"Meshery Catalog"}
            </Typography>
            <Typography variant="body" component="div" className={classes.backSideDescription}>
              {"Connect to enable access to the cloud native catalog. Import any catalog item and customize."}
            </Typography>
            <Typography variant="subtitle2" style={{ fontStyle : "italic", marginTop : ".75rem", marginBottom : "0rem" }}>
              Explore at: <a style={{ textDecoration : "none" }} href="https://meshery.io/catalog" target="_blank" rel="noreferrer"> Meshery Catalog </a>
            </Typography>
          </div>
        </FlipCard>
      </Grid>
    </Grid>
  )
}

export default withStyles(styles)(Extensions);