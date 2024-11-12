import React from 'react';
// import AnimatedMeshery from './Animations/AnimatedMesheryCSS';
import PropTypes from 'prop-types';
import uiConfig from '../../ui.config';
// import { Typography } from '@layer5/sistent';
// import { UsesSistent } from '../SistentWrapper';

// const StyledAnimatedLogoDark = styled(uiConfig.AnimatedLogoDark)(({ theme }) => ({
//   height: '100px',
//   margin: '4px 0px 8px',
//   fill: theme.palette.mode === 'dark' ? '#fff' : '',
// }));

function LoadingScreen(props) {
  const { message, children, isLoading, ...other } = props;
  const AnimatedLogoDark = uiConfig.AnimatedLogoDark;

  if (isLoading) {
    return (
      <div
        {...other}
        style={{
          '@font-face': {
            fontFamily: 'Qanelas Soft',
            src: 'url("/static/fonts/qanelas-soft/QanelasSoftBlack.otf") format("opentype")',
            fontWeight: 'normal',
          },
          fontFamily: 'QanelasSoftRegular, Qanelas Soft Regular, sans-serif',
          display: 'grid',
          placeItems: 'center',
          minHeight: '100vh',
          textAlign: 'center',
          backgroundColor: '#263238', // color of navigation menu
          color: '#dedede', // soften the subtitle / message
        }}
      >
        <div>
          <AnimatedLogoDark />
          <h1
            style={{
              fontFamily: 'QanelasSoftRegular',
              fontSize: '.9rem',
              fontWeight: 'normal',
              marginTop: '1rem',
            }}
          >
            {message}
          </h1>
        </div>
      </div>
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
