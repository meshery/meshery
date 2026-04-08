import React, { useEffect } from 'react';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import LoadingScreen from './LoadingComponents/LoadingComponent';
import Troubleshoot from './TroubleshootingComponent';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { useRouter } from 'next/router';

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

const DASHBOARD_REDIRECT_DELAY_MS = 5000;

const RemoteComponent: React.FC<RemoteComponentProps> = ({ url, loaderType, ...props }) => {
  const [loading, err, RemoteComponentInstance] = useRemoteComponent(url.url);
  const { notify } = useNotification();
  const router = useRouter();

  useEffect(() => {
    if (err == null) return;

    notify({
      message: `Failed to load extension: ${err.toString()}`,
      event_type: EVENT_TYPES.ERROR,
    });

    const timer = setTimeout(() => {
      router.push('/');
    }, DASHBOARD_REDIRECT_DELAY_MS);

    return () => clearTimeout(timer);
  }, [err]);

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
    return (
      <div role="alert">
        <h2>Uh-oh! Please pardon our mesh.</h2>
        <div
          style={{
            backgroundColor: '#1E2117',
            color: '#FFFFFF',
            padding: '.85rem',
            borderRadius: '.2rem',
          }}
        >
          <code>{err.toString()}</code>
        </div>
        <p style={{ marginTop: '1rem', color: '#666' }}>
          Redirecting to dashboard in a few seconds...
        </p>
        <div style={{ marginTop: '1rem' }}>
          <Troubleshoot showDesignerButton={false} />
        </div>
      </div>
    );
  }

  return (
    <div>
      <RemoteComponentInstance {...props} />
    </div>
  );
};

export default RemoteComponent;
