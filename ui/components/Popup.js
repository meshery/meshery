import React, { useEffect, useState } from 'react';
import Cookies from 'universal-cookie';
import GenericModal from './GenericModal';
import { makeStyles, MuiThemeProvider } from '@material-ui/core/styles';
import { createTheme } from '@material-ui/core/styles';
import { Typography } from '@material-ui/core';


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
    width : 600,
    backgroundColor : theme.palette.background.paper,
    border : '0px solid #000',
    boxShadow : theme.shadows[5],
    padding : theme.spacing(2, 4, 3),
    top : '50%',
    left : '50%',
    transform : `translate(-50%, -50%)`,
    borderRadius : 10,
  },
  grid : {
    width : '100%'
  },


}));

function Popup() {
  const [isOpen, setIsOpen] = useState(true);
  const cookies = new Cookies('registered');
  const classes = styles();


  useEffect(() => {
    if (cookies.get('registered')) {
      setIsOpen(false);
    } else if (!cookies.get('registered')) {
      cookies.set('registered', 'true', {
        path : '/',
      });
      setIsOpen(true);
    }
  },[])

  const handleClose = () => {
    setIsOpen(false);
  }

  return (
    <GenericModal
      open={isOpen}
      handleClose={handleClose}
      Content={
        <MuiThemeProvider theme={getMuiTheme()}>
          <div
            className={classes.paper}
          >
            <Typography>Have you taken a look at MeshMap</Typography>
          </div>
        </MuiThemeProvider>
      }
    />
  )
}

export default Popup