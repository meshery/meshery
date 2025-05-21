import React from 'react';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import LoadingScreen from './LoadingComponents/LoadingComponent';

const requires = createRequires(getDependencies);

const useRemoteComponent = createUseRemoteComponent({ requires });

const RemoteComponent = ({ url, loaderType, props = {} }) => {
  const [loading, err, RemoteComponent] = useRemoteComponent(url.url);
  if (loading) {
    if (loaderType === 'circular') {
      return '';
    } else {
      return (
        <LoadingScreen animatedIcon="AnimatedMeshery" message="Establishing Remote Connection" />
      );
    }
  }
  if (err != null) {
    /* Debugging log */
    console.error(`Extension Error: ${err.toString()}`);
    return <></>;
  }

  return (
    <div>
      <RemoteComponent {...props} />
    </div>
  );
};

export default RemoteComponent;
