import React from 'react';
import uiConfig from '../../ui.config';

type LoadingScreenProps = {
  message: string;
  children?: React.ReactNode;
  isLoading: boolean;
  id?: string;
  className?: string;
  style?: React.CSSProperties;
  [key: string]: any;
};

function LoadingScreen(props: LoadingScreenProps) {
  const { message, children, isLoading, id, ...other } = props;

  if (isLoading) {
    return <PureHtmlLoadingScreen message={message} id={id} {...other} />;
  }

  return <>{children}</>;
}

type PureHtmlLoadingScreenProps = {
  message: string;
  id?: string | undefined;
  style?: React.CSSProperties;
  [key: string]: any;
};

export const PureHtmlLoadingScreen = (props: PureHtmlLoadingScreenProps) => {
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
            fontSize: '.9rem',
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

export default LoadingScreen;
