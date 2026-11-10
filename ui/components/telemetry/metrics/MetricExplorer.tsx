import React, { useEffect, useMemo, useState } from 'react';
import {
  Box,
  Button,
  Chip,
  CircularProgress,
  Divider,
  FormControl,
  InputLabel,
  List,
  ListItemButton,
  MenuItem,
  Select,
  TextField,
  Typography,
  styled,
} from '@sistent/sistent';
import {
  useLazyQueryPrometheusRangeQuery,
  useListPrometheusMetricsQuery,
  usePrometheusMetricMetadataQuery,
} from '@/rtk-query/telemetryPrometheus';
import TimeSeriesChart from '../common/TimeSeriesChart';
import { parsePromMatrix } from '../common/time';
import type { ChartSeries, TimeWindow } from '../common/types';
import { PANEL_TYPES, type MetricMetadata, type MetricPanel } from './types';

interface MetricExplorerProps {
  connectionID: string;
  timeWindow: TimeWindow;
  onAdd: (draft: Omit<MetricPanel, 'id' | 'gridPos'>) => void;
}

const Section = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(1),
}));

const MetricList = styled(List)(({ theme }) => ({
  maxHeight: 220,
  overflow: 'auto',
  border: `1px solid ${theme.palette.border.default}`,
  borderRadius: theme.shape.borderRadius,
  padding: 0,
}));

const PreviewBox = styled(Box)(({ theme }) => ({
  height: 180,
  border: `1px solid ${theme.palette.border.default}`,
  borderRadius: theme.shape.borderRadius,
  padding: theme.spacing(1),
}));

// Pull the first metadata entry for a metric out of Prometheus' /metadata shape:
// { data: { metricName: [ { type, help, unit } ] } }.
const firstMetadata = (resp: any, metric: string): MetricMetadata | undefined =>
  resp?.data?.[metric]?.[0];

/**
 * MetricExplorer lets a user search Prometheus metric names, inspect a metric's
 * metadata, compose a PromQL expression, preview it live, and save it as a panel.
 */
const MetricExplorer: React.FC<MetricExplorerProps> = ({ connectionID, timeWindow, onAdd }) => {
  const [search, setSearch] = useState('');
  const [debounced, setDebounced] = useState('');
  const [metric, setMetric] = useState('');
  const [expr, setExpr] = useState('');
  const [title, setTitle] = useState('');
  const [type, setType] = useState<string>('timeseries');
  const [unit, setUnit] = useState('');
  const [preview, setPreview] = useState<ChartSeries[] | null>(null);

  // Debounce the metric-name search so we don't query on every keystroke.
  useEffect(() => {
    const id = setTimeout(() => setDebounced(search.trim()), 300);
    return () => clearTimeout(id);
  }, [search]);

  const { data: metricsData, isFetching: metricsLoading } = useListPrometheusMetricsQuery(
    { connectionID, search: debounced, limit: 50 },
    { skip: !connectionID },
  );
  const metrics = (metricsData as string[] | undefined) ?? [];

  const { data: metadataResp } = usePrometheusMetricMetadataQuery(
    { connectionID, metric },
    { skip: !connectionID || !metric },
  );
  const metadata = useMemo(
    () => (metric ? firstMetadata(metadataResp, metric) : undefined),
    [metadataResp, metric],
  );

  const [triggerPreview, { isFetching: previewLoading, error: previewError }] =
    useLazyQueryPrometheusRangeQuery();

  const pickMetric = (name: string) => {
    setMetric(name);
    if (!expr) setExpr(name);
    if (!title) setTitle(name);
  };

  const runPreview = async () => {
    if (!expr.trim()) return;
    try {
      const resp = await triggerPreview({
        connectionID,
        query: expr,
        start: String(timeWindow.start),
        end: String(timeWindow.end),
        step: String(timeWindow.step),
      }).unwrap();
      setPreview(parsePromMatrix(resp));
    } catch {
      setPreview([]);
    }
  };

  const canAdd = expr.trim() !== '' && title.trim() !== '';

  const handleAdd = () => {
    if (!canAdd) return;
    onAdd({ title: title.trim(), expr: expr.trim(), type, unit: unit.trim() || undefined });
    // Reset the query-specific fields, keep the metric list state.
    setExpr('');
    setTitle('');
    setUnit('');
    setPreview(null);
  };

  return (
    <Section>
      <Typography variant="subtitle2">1. Find a metric</Typography>
      <TextField
        size="small"
        placeholder="Search metrics (e.g. container_memory)"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        fullWidth
      />
      <MetricList>
        {metricsLoading ? (
          <Box sx={{ p: 2, display: 'flex', justifyContent: 'center' }}>
            <CircularProgress size={18} />
          </Box>
        ) : metrics.length === 0 ? (
          <Box sx={{ p: 2 }}>
            <Typography variant="caption" color="textSecondary">
              {debounced ? 'No matching metrics' : 'Type to search metric names'}
            </Typography>
          </Box>
        ) : (
          metrics.map((name) => (
            <ListItemButton
              key={name}
              selected={name === metric}
              onClick={() => pickMetric(name)}
              dense
            >
              <Typography variant="body2" noWrap>
                {name}
              </Typography>
            </ListItemButton>
          ))
        )}
      </MetricList>
      {metadata && (
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, alignItems: 'center' }}>
          {metadata.type && <Chip size="small" label={metadata.type} variant="outlined" />}
          {metadata.unit && <Chip size="small" label={metadata.unit} variant="outlined" />}
          {metadata.help && (
            <Typography variant="caption" color="textSecondary">
              {metadata.help}
            </Typography>
          )}
        </Box>
      )}

      <Divider sx={{ my: 1 }} />

      <Typography variant="subtitle2">2. Write PromQL</Typography>
      <TextField
        size="small"
        placeholder="PromQL expression"
        value={expr}
        onChange={(e) => setExpr(e.target.value)}
        multiline
        minRows={2}
        fullWidth
      />
      <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
        <Button
          variant="outlined"
          size="small"
          onClick={runPreview}
          disabled={!expr.trim() || previewLoading}
        >
          {previewLoading ? 'Previewing…' : 'Preview'}
        </Button>
        {previewError && (
          <Typography variant="caption" color="error">
            {(previewError as any)?.data?.error || 'Query failed'}
          </Typography>
        )}
      </Box>
      {preview !== null && (
        <PreviewBox>
          {preview.length === 0 ? (
            <Box
              sx={{
                height: '100%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <Typography variant="caption" color="textSecondary">
                No data for this query/range
              </Typography>
            </Box>
          ) : (
            <TimeSeriesChart series={preview} unit={unit} filled />
          )}
        </PreviewBox>
      )}

      <Divider sx={{ my: 1 }} />

      <Typography variant="subtitle2">3. Save as panel</Typography>
      <TextField
        size="small"
        label="Panel title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        fullWidth
      />
      <Box sx={{ display: 'flex', gap: 1 }}>
        <FormControl size="small" sx={{ minWidth: 140 }}>
          <InputLabel id="metric-panel-type">Visualization</InputLabel>
          <Select
            labelId="metric-panel-type"
            label="Visualization"
            value={type}
            onChange={(e) => setType(String(e.target.value))}
          >
            {PANEL_TYPES.map((t) => (
              <MenuItem key={t} value={t}>
                {t}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <TextField
          size="small"
          label="Unit (optional)"
          placeholder="bytes, percent, s…"
          value={unit}
          onChange={(e) => setUnit(e.target.value)}
          sx={{ flex: 1 }}
        />
      </Box>
      <Button variant="contained" onClick={handleAdd} disabled={!canAdd}>
        Add panel
      </Button>
    </Section>
  );
};

export default MetricExplorer;
