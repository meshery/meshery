import React, {
  useEffect,
  useState
} from 'react';
import Cookies from 'universal-cookie';
import { makeStyles } from '@material-ui/core/styles';
import { Typography, Grid, Button, IconButton } from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';

const styles = makeStyles((theme) => ({
  paper : {
    position : 'fixed',
    width : 450,
    backgroundColor : theme.palette.background.paper,
    border : '0px solid #000',
    boxShadow : theme.shadows[5],
    padding : theme.spacing(1, 2, 3, 4),
    right : 0,
    bottom : 0,
    borderRadius : 10,
    ["@media (max-width: 455px)"] : {
      width : "100%"
    },
    zIndex : 5
  },
  grid : {
    width : '100%'
  },
  designerImg : {
    height : 300,
    margin : "auto"
  },
  header : {
    paddingBottom : "0.5rem",
    paddingTop : "0.6rem",
    position : "absolute",
    fontWeight : "bold",
    ["@media (max-width: 455px)"] : {
      fontSize : "1rem"
    },
  },
  caption : {
    lineHeight : "1.2",
    paddingBottom : "15px",
    fontSize : ".75rem",
    textAlign : "center"
  },
  imgWrapper : {
    padding : "15px 10px 15px 0",
    display : "flex"
  },
  headerWrapper : {
    marginBottom : 12
  }
}));

function MeshMapEarlyAccessCardPopup() {
  const [isOpen, setIsOpen] = useState(true);
  const cookies = new Cookies('registered');

  const handleOpen = () => {
    const timer = setTimeout(() => {
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
  }, [])

  if (isOpen) {
    return <MeshMapEarlyAccessCard closeForm={() => setIsOpen(false)} />
  } else {
    return <></>
  }
}

export function MeshMapEarlyAccessCard({ rootStyle = {}, closeForm = () => {} }) {
  const classes = styles();

  const handleSignUp = (e) => {
    window.open("https://layer5.io/meshmap", "_blank")
    e.stopPropagation();
  };

  return (
    <div
      className={classes.paper}
      style={rootStyle}
    >
      <div className={classes.headerWrapper}>
        <Typography className={classes.header} variant="h6">Get early access to MeshMap!
        </Typography>

        <div style={{ display : "flex", justifyContent : "flex-end", whiteSpace : "nowrap", position : "relative" }}>
          <IconButton key="close" aria-label="Close" color="inherit" onClick={closeForm}>
            <CloseIcon />
          </IconButton>
        </div>
      </div>

      <div className={classes.imgWrapper}>
        <img className={classes.designerImg} src="/static/img/designer.png" />
      </div>
      <Typography className={classes.caption} variant="subtitle1"><i>Friends dont let friends GitOps alone.
        Visually design and collaborate in real-time with other MeshMap users.</i></Typography>
      <div style={{ display : "flex", justifyContent : "flex-end" }}>
        <Grid item xs={3}>
          <Button fullWidth variant="contained" color="primary" onClick={(e) => handleSignUp(e)}>Sign up</Button>
        </Grid>
      </div>
    </div>
  )
}

export default MeshMapEarlyAccessCardPopup;
