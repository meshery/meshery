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

// The RTK Query ping hook shape both telemetry modules expose. Kept loose so
// either module's generated hook satisfies it.
type PingHook = (
  arg: { connectionID: string },
  opts: { skip?: boolean; refetchOnMountOrArgChange?: boolean },
) => { isFetching: boolean; isError: boolean; error?: unknown };

interface PingStatusProps {
  connectionID: string;
  // The module's ping query hook (e.g. usePingGrafanaConnectionQuery).
  usePingQuery: PingHook;
  // Product name for the tooltip text, e.g. "Grafana" / "Prometheus".
  productName?: string;
}

/**
 * PingStatus shows whether the selected telemetry connection is currently
 * reachable and its credential accepted. The ping hook is injected so this
 * component is reusable across telemetry modules.
 */
const PingStatus: React.FC<PingStatusProps> = ({
  connectionID,
  usePingQuery,
  productName = 'The service',
}) => {
  const theme = useTheme();
  const { isFetching, isError, error } = usePingQuery(
    { connectionID },
    { skip: !connectionID, refetchOnMountOrArgChange: true },
  );

  let icon: React.ReactNode;
  let label: string;
  let tip: string;
  if (isFetching) {
    icon = <CircularProgress size={14} />;
    label = 'Checking…';
    tip = `Checking ${productName} reachability`;
  } else if (isError) {
    icon = <ErrorIcon style={{ color: theme.palette.error.default, fontSize: 16 }} />;
    label = 'Unreachable';
    tip =
      (error as any)?.data?.error || `${productName} is unreachable or the credential was rejected`;
  } else {
    icon = <CheckCircleIcon style={{ color: theme.palette.success.main, fontSize: 16 }} />;
    label = 'Connected';
    tip = `${productName} is reachable`;
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
