import { Button, Grid, IconButton, Typography, styled } from '@layer5/sistent';
import { UsesSistent } from './SistentWrapper';
import { useTheme } from '@material-ui/core';
import CloseIcon from '@mui/icons-material/Close';
import { useRouter } from 'next/router';
import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import Cookies from 'universal-cookie';
import { mesheryExtensionRoute } from '../pages/_app';
import { Colors } from '@/themes/app';

const StyledPaper = styled('div')(({ theme }) => ({
  position: 'fixed',
  width: 450,
  backgroundColor:
    theme.palette.mode === 'dark'
      ? theme.palette.background.hover
      : theme.palette.background.neutral?.hover,
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
}));

const StyledHeaderWrapper = styled('div')(() => ({
  marginBottom: 12,
  display: 'flex',
  justifyContent: 'space-between',
}));

const StyledCloseButtonContainer = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'flex-end',
  whiteSpace: 'nowrap',
  alignItems: 'center',
  color: '#F6F8F8',
}));

const StyledImgWrapper = styled('div')(() => ({
  padding: '0px 10px 15px 0',
  display: 'flex',
}));

const StyledDesignerImg = styled('img')(() => ({
  height: '205px',
  width: 'auto',
  margin: 'auto',
  boxShadow:
    '1px 2px 2px hsl(173deg, 100%, 35% , 0.133), \
      2px 4px 4px hsl(173deg, 100%, 35% , 0.133),  \
      3px 6px 6px hsl(173deg, 100%, 35% , 0.133)',
}));

const isMesheryExtensionRegisteredUser = (capabilitiesRegistry) => {
  if (!capabilitiesRegistry) {
    return false;
  }

  return (
    capabilitiesRegistry.extensions?.navigator?.length > 0 &&
    capabilitiesRegistry.extensions.navigator.find((ext) => ext.title === 'Kanvas')
  );
};

export function MesheryExtensionEarlyAccessCardPopup({ capabilitiesRegistry }) {
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
      <MesheryExtensionEarlyAccessCard
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

export function MesheryExtensionEarlyAccessCard({
  rootStyle = {},
  closeForm = () => {},
  capabilitiesRegistry,
}) {
  const signUpText = 'Sign up';
  const signupHeader = 'Get early access to Kanvas!';
  const [buttonText, setButtonText] = useState(signUpText);
  const [title, setTitle] = useState(signupHeader);
  const { push } = useRouter();
  const theme = useTheme();
  const popupImageSrc =
    theme.palette.type === 'dark' ? '/static/img/aws.svg' : '/static/img/aws-light.svg';

  const handleButtonClick = (e) => {
    if (buttonText === signUpText) {
      window.open('https://layer5.io/cloud-native-management/kanvas', '_blank');
    } else {
      push(mesheryExtensionRoute);
    }
    e.stopPropagation();
  };

  useState(() => {
    const isMesheryExtensionUser = isMesheryExtensionRegisteredUser(capabilitiesRegistry);
    if (isMesheryExtensionUser) {
      setTitle('Collaborative management enabled');
      setButtonText('Open Kanvas');
    } else {
      setTitle(signupHeader);
      setButtonText(signUpText);
    }
  }, [capabilitiesRegistry]);

  return (
    <UsesSistent>
      <StyledPaper style={rootStyle}>
        <StyledHeaderWrapper>
          <Typography
            sx={{
              paddingBottom: '0.5rem',
              paddingTop: '0.6rem',
              fontWeight: 'bold',
              color: '#F6F8F8',
              ['@media (max-width: 455px)']: {
                fontSize: '1rem',
              },
            }}
            variant="h6"
          >
            {title}
          </Typography>

          <StyledCloseButtonContainer>
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={closeForm}
              style={{ height: '2.5rem' }}
            >
              <CloseIcon />
            </IconButton>
          </StyledCloseButtonContainer>
        </StyledHeaderWrapper>

        <StyledImgWrapper>
          <StyledDesignerImg src={popupImageSrc} alt="pop-up" />
        </StyledImgWrapper>
        <Typography
          sx={{
            lineHeight: '1.2',
            paddingBottom: '15px',
            fontSize: '.85rem',
            textAlign: 'center',
            color: '#F6F8F8',
          }}
          variant="subtitle1"
        >
          <i>
            Friends don&apos;t let friends GitOps alone. Visually design and collaborate in
            real-time with other Meshery users.
          </i>
        </Typography>
        <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
          <Grid>
            <Button
              variant="contained"
              sx={{
                backgroundColor: Colors.keppelGreen,
                borderRadius: 2,
                height: 35,
                fontSize: '0.8rem',
              }}
              onClick={(e) => handleButtonClick(e)}
            >
              {buttonText}
            </Button>
          </Grid>
        </div>
      </StyledPaper>
    </UsesSistent>
  );
}

const mapStateToProps = (state) => ({
  capabilitiesRegistry: state.get('capabilitiesRegistry'),
});

export default connect(mapStateToProps)(MesheryExtensionEarlyAccessCardPopup);
