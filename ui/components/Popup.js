import { Button, Grid, IconButton, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import CloseIcon from '@material-ui/icons/Close';
import { useRouter } from 'next/router';
import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import Cookies from 'universal-cookie';
import { mesheryExtensionRoute } from '../pages/_app';

const styles = makeStyles((theme) => ({
  paper: {
    position: 'fixed',
    width: 450,
    backgroundColor: theme.palette.background.paper,
    border: '0px solid #000',
    boxShadow: theme.shadows[5],
    padding: theme.spacing(1, 2, 3, 4),
    right: 0,
    bottom: 0,
    borderRadius: 10,
    ['@media (max-width: 455px)']: {
      width: '100%',
    },
    zIndex: 1201,
  },
  grid: {
    width: '100%',
  },
  designerImg: {
    height: 300,
    margin: 'auto',
  },
  header: {
    paddingBottom: '0.5rem',
    paddingTop: '0.6rem',
    fontWeight: 'bold',
    ['@media (max-width: 455px)']: {
      fontSize: '1rem',
    },
  },
  closeButtonContainer: {
    display: 'flex',
    justifyContent: 'flex-end',
    whiteSpace: 'nowrap',
    alignItems: 'center',
  },
  caption: {
    lineHeight: '1.2',
    paddingBottom: '15px',
    fontSize: '.75rem',
    textAlign: 'center',
  },
  imgWrapper: {
    padding: '15px 10px 15px 0',
    display: 'flex',
  },
  headerWrapper: {
    marginBottom: 12,
    display: 'flex',
    justifyContent: 'space-between',
  },
}));

const isMeshMapRegisteredUser = (capabilitiesRegistry) => {
  if (!capabilitiesRegistry) {
    return false;
  }

  return (
    capabilitiesRegistry.extensions?.navigator?.length > 0 &&
    capabilitiesRegistry.extensions.navigator.find((ext) => ext.title === 'MeshMap')
  );
};

function MeshMapEarlyAccessCardPopup({ capabilitiesRegistry }) {
  const [isOpen, setIsOpen] = useState(false);
  const cookies = new Cookies('registered');

  const closeCallback = () => {
    cookies.set('registered', 'true', { path: '/' });
  };

  useEffect(() => {
    // cookies return string and not boolean thus truthy,falsy doesnt work as intended
    const isAlreadyRegistered = cookies.get('registered') && cookies.get('registered') === 'true';

    if (isAlreadyRegistered) {
      return;
    }

    const timer = setTimeout(() => {
      setIsOpen(true);
    }, 10000); // 10sec waiting time

    return () => clearTimeout(timer);
  }, []);

  if (isOpen) {
    return (
      <MeshMapEarlyAccessCard
        closeForm={() => {
          setIsOpen(false);
          closeCallback();
        }}
        capabilitiesRegistry={capabilitiesRegistry}
      />
    );
  } else {
    return <></>;
  }
}

export function MeshMapEarlyAccessCard({
  rootStyle = {},
  closeForm = () => {},
  capabilitiesRegistry,
}) {
  const signUpText = 'Sign up';
  const signupHeader = 'Get early access to MeshMap!';
  const classes = styles();
  const [buttonText, setButtonText] = useState(signUpText);
  const [title, setTitle] = useState(signupHeader);
  const { push } = useRouter();

  const handleButtonClick = (e) => {
    if (buttonText === signUpText) {
      window.open('https://layer5.io/meshmap', '_blank');
    } else {
      push(mesheryExtensionRoute);
    }
    e.stopPropagation();
  };

  useState(() => {
    const isMeshMapUser = isMeshMapRegisteredUser(capabilitiesRegistry);
    if (isMeshMapUser) {
      setTitle('Your access to collaborative cloud native management is enabled!');
      setButtonText('Open MeshMap');
    } else {
      setTitle(signupHeader);
      setButtonText(signUpText);
    }
  }, [capabilitiesRegistry]);

  return (
    <div className={classes.paper} style={rootStyle}>
      <div className={classes.headerWrapper}>
        <Typography className={classes.header} variant="h6">
          {title}
        </Typography>

        <div className={classes.closeButtonContainer}>
          <IconButton
            key="close"
            aria-label="Close"
            color="inherit"
            onClick={closeForm}
            style={{ height: '2.5rem' }}
          >
            <CloseIcon />
          </IconButton>
        </div>
      </div>

      <div className={classes.imgWrapper}>
        <img className={classes.designerImg} src="/static/img/designer.png" />
      </div>
      <Typography className={classes.caption} variant="subtitle1">
        <i>
          Friends don&apos;t let friends GitOps alone. Visually design and collaborate in real-time
          with other MeshMap users.
        </i>
      </Typography>
      <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
        <Grid>
          <Button
            fullWidth
            variant="contained"
            color="primary"
            onClick={(e) => handleButtonClick(e)}
          >
            {buttonText}
          </Button>
        </Grid>
      </div>
    </div>
  );
}

const mapStateToProps = (state) => ({
  capabilitiesRegistry: state.get('capabilitiesRegistry'),
});

export default connect(mapStateToProps)(MeshMapEarlyAccessCardPopup);
