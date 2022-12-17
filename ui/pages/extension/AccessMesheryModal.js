import { makeStyles } from '@material-ui/core/styles';
import { Typography, Button, Dialog, DialogActions, DialogContent, DialogTitle, IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';

const styles = makeStyles((theme) => ({
  root : {
    margin : 0,
    padding : theme.spacing(2),
  },
  closeButton : {
    position : "absolute",
    right : theme.spacing(1),
    top : theme.spacing(1),
    color : theme.palette.common.white,
    transform : "rotate(-90deg)",
    "&:hover" : {
      transform : "rotate(90deg)",
      transition : "all .3s ease-in"
    },
  },
  dialogTitle : {
    textAlign : "center",
    backgroundColor : "#252E31",
    color : "#fff"
  },
  imgWrapper : {
    display : "flex",
    justifyContent : "center",
    padding : "1rem"
  }
}));


export default function PlaygroundMeshDeploy(props) {
  const classes = styles();

  const handlePage = (e) => {
    window.open("https://meshery.io/#getting-started", "_blank")
    e.stopPropagation();
  };

  return (
    <div>
      <Dialog aria-labelledby="customized-dialog-title" open={props.isOpen} >
        <DialogTitle id="customized-dialog-title" className={classes.dialogTitle}>
          <Typography variant="h6">The Cloud Native Playground</Typography>
          <IconButton aria-label="close" className={classes.closeButton} onClick={props.closeForm}>
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent>
          <div className={classes.imgWrapper}>
            <img width="60%" height="60%" src="/static/img/meshery-logo-light-text-side.png" />
          </div>
          <Typography gutterBottom>
            Meshery Playground gives you hands-on experience using not just Kubernetes, but every CNCF project and surrounding ecosystem.
            Choose your <a href="https://layer5.io/learn/learning-paths">Learning Path</a> and follow their labs in your browser as you visually and collaboratively learn by doing without having to install a single thing.
          </Typography>
          <Typography gutterBottom>
            To ensure that Meshery Playground remains a clean sandbox of all to use, many of Meshery&apos;s features are disabled.
            aaFor full access to all of Meshery&apos;s features, deploy your own instance of Meshery.
          </Typography>
        </DialogContent>
        <DialogActions style={{ diasplay : "flex", justifyContent : "center" }}>
          <Button size="medium" variant="outlined" color="primary" onClick={(e) => handlePage(e)}>Get Started</Button>
        </DialogActions>
      </Dialog>
    </div>
  )
}