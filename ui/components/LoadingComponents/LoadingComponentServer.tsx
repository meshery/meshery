// @ts-nocheck
import React from 'react';
import PropTypes from 'prop-types';
import uiConfig from '../../ui.config';

function LoadingScreen(props) {
  const { message, children, isLoading, id, ...other } = props;

  if (isLoading) {
    return <PureHtmlLoadingScreen message={message} id={id} {...other} />;
  }

  return <>{children}</>;
}

export const PureHtmlLoadingScreen = (props) => {
  const { message, id, ...other } = props;
  const AnimatedLogoDark = uiConfig.AnimatedLogoDark;

  return (
    <div
      {...other}
      style={{
        minHeight: '100vh',
        position: 'absolute',
        inset: '0',
        zIndex: 99999,
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        textAlign: 'center',
        fontFamily: 'QanelasSoftRegular, sans-serif',
        fontWeight: 'normal',
        fontSize: '1.15rem',
        backgroundColor: '#263238', // color of navigation menu
        color: '#dedede', // soften the subtitle / message
        ...(other.style || {}),
      }}
      id={id}
    >
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          gap: '1rem',
        }}
      >
        <AnimatedLogoDark />
        <h1
          id={id + '-text-message'}
          style={{
            fontFamily: 'QanelasSoftRegular, sans-serif', // this is important to have consistent font between prereact render
            fontSize: '1.15rem',
            fontWeight: 'normal',
            marginTop: '0rem',
          }}
        >
          {message}
        </h1>
      </div>
    </div>
  );
};

LoadingScreen.propTypes = {
  message: PropTypes.string.isRequired,
  className: PropTypes.string,
  isLoading: PropTypes.bool.isRequired,
};

export default LoadingScreen;
