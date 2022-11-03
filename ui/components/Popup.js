import React, {
  useEffect,
  useState } from 'react';
import Cookies from 'universal-cookie';
// import GenericModal from './GenericModal';
import { makeStyles, MuiThemeProvider } from '@material-ui/core/styles';
import { createTheme } from '@material-ui/core/styles';
import { Typography, Grid, Button,Modal, IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';


const getMuiTheme = () => createTheme({
  palette : {
    primary : {
      main : "#607d8b"
    }
  },
  overrides : {
    MuiGrid : {
      input : {
        color : '#607d8b'
      }
    },
  }
})

const styles = makeStyles((theme) => ({
  paper : {
    position : 'absolute',
    width : 450,
    backgroundColor : theme.palette.background.paper,
    border : '0px solid #000',
    boxShadow : theme.shadows[5],
    padding : theme.spacing(1, 2, 3, 4),
    top : '70%',
    left : '83%',
    transform : `translate(-50%, -50%)`,
    borderRadius : 10,
  },
  grid : {
    width : '100%'
  },
  drawer : {
    height : 300,
    position : "relative"
  },

  canvas : {
    width : 275,
    position : "absolute",
    paddingTop : "1rem",
    right : 15
  },
  header : {
    paddingBottom : "0.5rem",
    paddingTop : "0.6rem",
    position : "absolute"
  }


}));

function Popup() {
  const [isOpen, setIsOpen] = useState(true);
  const cookies = new Cookies('registered');
  const classes = styles();

  const handleOpen = () => {
    const timer =  setTimeout(() => {
      setIsOpen(true);
    }, 10000)
    return () => clearTimeout(timer);
  }

  useEffect(() => {
    if (cookies.get('registered')) {
      setIsOpen(false);
    } else if (!cookies.get('registered')) {
      cookies.set('registered', 'true', {
        path : '/',
      });
      handleOpen();
    }
  },[])

  const handleClose = () => {
    setIsOpen(false);
  }

  const handleSignUp = (e) => {
    window.open("https://layer5.io/meshmap", "_blank")
    e.stopPropagation();
  };

  return (
    <div>

      <Modal hideBackdrop={true}    open={isOpen} onClose={handleClose}>

        <MuiThemeProvider theme={getMuiTheme()}>
          <div
            className={classes.paper}
          >

            <Typography className={classes.header} variant="h6" >Get early access to MeshMap!
            </Typography>

            <div style={{ display : "flex", justifyContent : "flex-end", whiteSpace : "nowrap", position : "relative" }}>
              <IconButton  key="close" aria-label="Close" color="inherit" onClick={() => setIsOpen(false)}>
                <CloseIcon />
              </IconButton>

            </div>
            <img className={classes.drawer} src="/static/img/designer-drawer.png" />
            <img className={classes.canvas} src="/static/img/designer-canvas.png" />
            <Typography variant="subtitle1"><i>Friends dont let friends GitOps alone. Visually design and collaborate in real-time with other MeshMap users.</i></Typography>
            <div style={{ display : "flex", justifyContent : "flex-end" }}>
              <Grid item xs={3}>
                <Button fullWidth variant="contained" color="primary" onClick={(e) => handleSignUp(e)}>Sign up</Button>
              </Grid>
              {/* <Grid item xs={3}>
                <Button fullWidth variant="contained" color="primary" onClick={() => setIsOpen(false)}>close</Button>
              </Grid> */}
            </div>
          </div>
        </MuiThemeProvider>

      </Modal>

      {/* <GenericModal
        open={isOpen}
        handleClose={handleClose}
        Content={

        }
      /> */}
    </div>
  )
}

export default Popup
