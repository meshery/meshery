import React, { useCallback, useEffect, useMemo, useState } from 'react';
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
import { useLazyQueryGrafanaRangeQuery } from '@/rtk-query/telemetryGrafana';
import TimeSeriesChart from './TimeSeriesChart';
import { formatValue, parsePromMatrix } from './time';
import { buildVarValues, resolveDatasourceUid, resolveExpr } from './resolve';
import type { ChartSeries, Datasource, Panel as PanelType, TemplateVar, TimeWindow } from './types';

interface PanelProps {
  connectionID: string;
  panel: PanelType;
  timeWindow: TimeWindow;
  templateVars?: TemplateVar[];
  datasources?: Datasource[];
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
 * Panel fetches the data for a single dashboard panel (one query per target)
 * and renders it as a time-series chart or a single-stat value, with explicit
 * loading / error / empty states.
 */
const Panel: React.FC<PanelProps> = ({
  connectionID,
  panel,
  timeWindow,
  templateVars = [],
  datasources = [],
}) => {
  const theme = useTheme();
  const [trigger] = useLazyQueryGrafanaRangeQuery();
  const [series, setSeries] = useState<ChartSeries[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const varValues = useMemo(() => buildVarValues(templateVars), [templateVars]);

  const fetchData = useCallback(async () => {
    const targets = (panel.targets ?? []).filter((t) => t.expr && t.expr.trim() !== '');
    if (targets.length === 0) {
      setSeries([]);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const collected: ChartSeries[] = [];
      for (const target of targets) {
        // Resolve Grafana template variables and macros into a concrete
        // datasource + executable PromQL before proxying the query.
        const ds = resolveDatasourceUid(target.datasourceUid, templateVars, datasources);
        const query = resolveExpr(target.expr, varValues, timeWindow);
        const resp = await trigger({
          connectionID,
          ds,
          query,
          start: String(timeWindow.start),
          end: String(timeWindow.end),
          step: String(timeWindow.step),
        }).unwrap();
        collected.push(...parsePromMatrix(resp, target.legendFormat));
      }
      setSeries(collected);
    } catch (e: any) {
      setError(e?.data?.error || e?.error || e?.message || 'Failed to load panel data');
      setSeries([]);
    } finally {
      setLoading(false);
    }
  }, [connectionID, panel, timeWindow, trigger, templateVars, datasources, varValues]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

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

export default Panel;
