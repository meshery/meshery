import React from 'react';
// import AnimatedMeshery from './Animations/AnimatedMesheryCSS';
import PropTypes from 'prop-types';
import uiConfig from '../../ui.config';
import { Typography, styled } from '@layer5/sistent';
import { UsesSistent } from '../SistentWrapper';

const StyledAnimatedLogoDark = styled(uiConfig.AnimatedLogoDark)(({ theme }) => ({
  height: '100px',
  margin: '4px 0px 8px',
  fill: theme.palette.mode === 'dark' ? '#fff' : '',
}));

function LoadingScreen(props) {
  const { message, children, isLoading, ...other } = props;

  if (isLoading) {
    return (
      <UsesSistent>
        <div
          {...other}
          style={{
            display: 'grid',
            placeItems: 'center',
            minHeight: '100vh',
            textAlign: 'center',
            backgroundColor: '#263238', // color of navigation menu
            color: '#dedede', // soften the subtitle / message
          }}
        >
          <div>
            <StyledAnimatedLogoDark />
            <Typography
              sx={{
                '@font-face': {
                  fontFamily: 'Qanelas Soft, Roboto',
                  src: 'url("/static/fonts/qanelas-soft/QanelasSoftBlack.otf") format("opentype")',
                  fontWeight: 'normal',
                  fontSize: '.9rem',
                 },
                marginTop: '1rem',
              }}
            >
              {message}
            </Typography>
          </div>
        </div>
      </UsesSistent>
    );
  }

  return <>{children}</>;
}

LoadingScreen.propTypes = {
  message: PropTypes.string.isRequired,
  className: PropTypes.string,
  isLoading: PropTypes.bool.isRequired,
};

export default LoadingScreen;
