import React from 'react';
import { randomLoadingMessage } from './loadingMessages';
import LoadingScreen from './LoadingComponentServer';
import { createPortal } from 'react-dom/cjs/react-dom.production.min';
import { NoSsr } from '@material-ui/core';

export const DynamicFullScrrenLoader = ({ children, ...props }) => {
  if (!props.isLoading) return children;

  if (document.body) {
    return (
      <NoSsr>
        {createPortal(
          <LoadingScreen
            {...props}
            message={randomLoadingMessage}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              zIndex: 9999999,
              height: '100vh',
              width: '100vw',
            }}
          />,
          document.body,
        )}
      </NoSsr>
    );
  }

  return <LoadingScreen {...props}> {children} </LoadingScreen>;
};
