import React from 'react';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  CircularProgress,
  Tooltip,
  Typography,
  WarningIcon,
  styled,
  useTheme,
} from '@sistent/sistent';
import TimeSeriesChart from './TimeSeriesChart';
import { formatValue } from './time';
import type { ChartSeries, PanelMeta } from './types';

interface PanelProps {
  panel: PanelMeta;
  // Data + state are computed once at the board/grid level (one batched request)
  // and passed down, so this component is purely presentational.
  series: ChartSeries[];
  loading: boolean;
  error: string | null;
}

const STAT_TYPES = new Set(['stat', 'singlestat', 'gauge']);
const FILLED_TYPES = new Set(['timeseries', 'graph']);

const PanelCard = styled(Card)(({ theme }) => ({
  height: '100%',
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  border: `1px solid ${theme.palette.border.default}`,
  boxShadow: 'none',
}));

const Body = styled(CardContent)(() => ({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  minHeight: 0,
  paddingTop: 0,
  '&:last-child': { paddingBottom: 8 },
}));

// Centers small states (loading / error / no-data / single stat).
const Center = styled(Box)(() => ({
  flex: 1,
  minHeight: 0,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  textAlign: 'center',
}));

// Lets the chart fill the panel body.
const ChartFill = styled(Box)(() => ({
  flex: 1,
  minHeight: 0,
  width: '100%',
}));

const StatValue = styled(Typography)(() => ({
  fontVariantNumeric: 'tabular-nums',
  fontWeight: 600,
  lineHeight: 1.1,
}));

/**
 * Panel renders a single panel from already-fetched series as either a
 * time-series chart or a single-stat value, with explicit loading / error /
 * empty states. Data fetching is owned by the parent grid (one batched request).
 */
const Panel: React.FC<PanelProps> = ({ panel, series, loading, error }) => {
  const theme = useTheme();

  const isStat = STAT_TYPES.has(panel.type);
  const hasData = series.length > 0 && series.some((s) => s.points.length > 0);

  let body: React.ReactNode;
  if (loading && series.length === 0) {
    body = (
      <Center>
        <CircularProgress size={22} />
      </Center>
    );
  } else if (error) {
    body = (
      <Center>
        <Tooltip title={error}>
          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 1 }}>
            <WarningIcon style={{ color: theme.palette.warning.default }} />
            <Typography variant="caption" color="textSecondary">
              Query failed
            </Typography>
          </Box>
        </Tooltip>
      </Center>
    );
  } else if (!hasData) {
    body = (
      <Center>
        <Typography variant="caption" color="textSecondary">
          No data
        </Typography>
      </Center>
    );
  } else if (isStat) {
    const last = series[0].points[series[0].points.length - 1];
    const value = last ? last[1] : NaN;
    body = (
      <Center>
        <Box>
          <StatValue variant="h3" color="primary">
            {formatValue(value, panel.unit)}
          </StatValue>
          {series[0].name && (
            <Typography variant="caption" color="textSecondary" noWrap display="block">
              {series[0].name}
            </Typography>
          )}
        </Box>
      </Center>
    );
  } else {
    body = (
      <ChartFill>
        <TimeSeriesChart series={series} unit={panel.unit} filled={FILLED_TYPES.has(panel.type)} />
      </ChartFill>
    );
  }

  return (
    <PanelCard>
      <CardHeader
        title={
          <Tooltip title={panel.title || ''}>
            <Typography variant="subtitle2" noWrap sx={{ fontWeight: 600 }}>
              {panel.title || 'Untitled'}
            </Typography>
          </Tooltip>
        }
        action={loading && series.length > 0 ? <CircularProgress size={14} /> : null}
        sx={{ py: 1, px: 2 }}
      />
      <Body>{body}</Body>
    </PanelCard>
  );
};

// Memoized so a panel only re-renders when its own data/state changes. Parents
// pass referentially-stable series objects (computed in a memoized map), so a
// board/grid re-render doesn't cascade into every panel's chart.
export default React.memo(Panel);
