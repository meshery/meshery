import React, { useEffect } from 'react';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import { Button, styled } from '@sistent/sistent';
import LoadingScreen from '../shared/LoadingState/LoadingComponent';
import Troubleshoot from '../TroubleshootingComponent';
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

const ErrorDetail = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.background.code ?? theme.palette.background.default,
  color: theme.palette.text.primary,
  padding: '.85rem',
  borderRadius: '.2rem',
}));

const ErrorAction = styled('div')({
  marginTop: '1rem',
});

const RemoteComponent: React.FC<RemoteComponentProps> = ({ url, loaderType, ...props }) => {
  const [loading, err, RemoteComponentInstance] = useRemoteComponent(url.url);
  const { notify } = useNotification();
  const router = useRouter();
  const isFullPage = loaderType !== 'circular';

  useEffect(() => {
    if (err == null || !isFullPage) return;

    notify({
      message: `Failed to load extension: ${err.toString()}`,
      event_type: EVENT_TYPES.ERROR,
    });
  }, [err, notify, isFullPage]);

  if (loading) {
    if (!isFullPage) {
      return null;
    } else {
      return (
        <LoadingScreen animatedIcon="AnimatedMeshery" message="Establishing Remote Connection" />
      );
    }
  }

  if (err != null) {
    if (!isFullPage) {
      return null;
    }

    return (
      <div role="alert">
        <h2>Uh-oh! Please pardon our mesh.</h2>
        <ErrorDetail>
          <code>{err.toString()}</code>
        </ErrorDetail>
        <ErrorAction>
          <Button variant="outlined" onClick={() => router.push('/')}>
            Go to dashboard
          </Button>
        </ErrorAction>
        <ErrorAction>
          <Troubleshoot showDesignerButton={false} />
        </ErrorAction>
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
