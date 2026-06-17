import React from 'react';
import {
  Box,
  CheckCircleIcon,
  CircularProgress,
  ErrorIcon,
  Tooltip,
  Typography,
  useTheme,
} from '@sistent/sistent';
import { usePingGrafanaConnectionQuery } from '@/rtk-query/telemetryGrafana';

interface PingStatusProps {
  connectionID: string;
}

/**
 * PingStatus shows whether the selected Grafana connection is currently
 * reachable and its credential accepted.
 */
const PingStatus: React.FC<PingStatusProps> = ({ connectionID }) => {
  const theme = useTheme();
  const { isFetching, isError, error } = usePingGrafanaConnectionQuery(
    { connectionID },
    { skip: !connectionID, refetchOnMountOrArgChange: true },
  );

  let icon: React.ReactNode;
  let label: string;
  let tip: string;
  if (isFetching) {
    icon = <CircularProgress size={14} />;
    label = 'Checking…';
    tip = 'Checking Grafana reachability';
  } else if (isError) {
    icon = <ErrorIcon style={{ color: theme.palette.error.default, fontSize: 16 }} />;
    label = 'Unreachable';
    tip = (error as any)?.data?.error || 'Grafana is unreachable or the credential was rejected';
  } else {
    icon = <CheckCircleIcon style={{ color: theme.palette.success.main, fontSize: 16 }} />;
    label = 'Connected';
    tip = 'Grafana is reachable';
  }

  return (
    <Tooltip title={tip}>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
        {icon}
        <Typography variant="caption" color="textSecondary">
          {label}
        </Typography>
      </Box>
    </Tooltip>
  );
};

export default PingStatus;
