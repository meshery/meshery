import React from 'react';
import {
  NoSsr,
  ErrorBoundary,
  PermissionSessionContext,
  type PermissionKeySpec,
} from '@sistent/sistent';
import CustomErrorFallback from '../../shared/ErrorBoundary/ErrorBoundary';

interface CurrentSessionInfoProps {
  permissionKey?: PermissionKeySpec;
}

const CurrentSessionInfo: React.FC<CurrentSessionInfoProps> = ({ permissionKey }) => {
  return <PermissionSessionContext variant="card" permissionKey={permissionKey} />;
};

const CurrentSessionInfoWithErrorBoundary: React.FC<CurrentSessionInfoProps> = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary customFallback={CustomErrorFallback}>
        <CurrentSessionInfo {...props} />
      </ErrorBoundary>
    </NoSsr>
  );
};

export default CurrentSessionInfoWithErrorBoundary;
