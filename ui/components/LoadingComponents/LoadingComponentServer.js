import React from 'react';
import AnimatedMeshery from './Animations/AnimatedMesheryCSS';
import PropTypes from 'prop-types';

function LoadingScreen(props) {
  const { message, children, isLoading, ...other } = props;

  if (isLoading) {
    return (
      <div
        {...other}
        style={{
          display: 'grid',
          placeItems: 'center',
          minHeight: '100vh',
          textAlign: 'center',
        }}
      >
        <div>
          <AnimatedMeshery style={{ height: '100px', margin: '4px 0px 8px' }} />
          <h1 style={{ fontFamily: 'sans-serif', fontSize: '18px', fontWeight: 'normal' }}>
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
