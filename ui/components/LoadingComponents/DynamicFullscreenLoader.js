import { useEffect } from 'react';

export const DynamicFullScrrenLoader = ({ children, isLoading }) => {
  useEffect(() => {
    const loader = window.Loader;
    if (!loader) {
      return;
    }
    if (isLoading) {
      loader.show();
    }

    if (loader && !isLoading) {
      loader.hide();
    }
  }, [isLoading]);

  if (!isLoading) return children;

  return null;

  // if (document.body) {
  //   // return createPortal(
  //   //   <LoadingScreen
  //   //     {...props}
  //   //     message={'Dynamic: ' + message}
  //   //     style={{
  //   //       position: 'absolute',
  //   //       top: 0,
  //   //       left: 0,
  //   //       zIndex: 9999999,
  //   //       height: '100vh',
  //   //       width: '100vw',
  //   //     }}
  //   //   />,
  //   //   document.body,
  //   // );
  // }

  // return <LoadingScreen {...props}> {children} </LoadingScreen>;
};
