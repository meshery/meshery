import React from 'react';
// import AnimatedMeshery from './Animations/AnimatedMesheryCSS';
import PropTypes from 'prop-types';
import uiConfig from '../../ui.config';

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
          display: 'grid',
          placeItems: 'center',
          minHeight: '100vh',
          textAlign: 'center',
          fontFamily: 'QanelasSoftRegular, Qanelas Soft Regular, sans-serif',
          backgroundColor: '#263238', // color of navigation menu
          color: '#dedede', // soften the subtitle / message
        }}
      >
        <div>
          <AnimatedLogoDark style={{ height: '100px', margin: '4px 0px 8px' }} />
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
