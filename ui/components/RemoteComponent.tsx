import React from 'react';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import LoadingScreen from './LoadingComponents/LoadingComponent';

const requires = createRequires(getDependencies);

const useRemoteComponent = createUseRemoteComponent({ requires });

interface RemoteComponentUrl {
  url: string;
}

interface RemoteComponentProps {
  url: RemoteComponentUrl;
  loaderType?: 'circular' | 'default';
  [key: string]: unknown;
}

const RemoteComponent: React.FC<RemoteComponentProps> = ({ url, loaderType, ...props }) => {
  const [loading, err, RemoteComponentInstance] = useRemoteComponent(url.url);

  if (loading) {
    if (loaderType === 'circular') {
      return null;
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
      <RemoteComponentInstance {...props} />
    </div>
  );
};

export default RemoteComponent;
