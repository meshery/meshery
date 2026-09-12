import React, { useMemo } from 'react';
import { Box, DeleteIcon, IconButton, Tooltip, styled } from '@sistent/sistent';
import { useQueryPrometheusRangeBatchQuery } from '@/rtk-query/telemetryPrometheus';
import Panel from '../common/Panel';
import { parsePromMatrix } from '../common/time';
import { resolveExpr } from '../common/promql';
import type { ChartSeries, TimeWindow } from '../common/types';
import type { MetricPanel } from './types';

interface MetricsGridProps {
  connectionID: string;
  panels: MetricPanel[];
  timeWindow: TimeWindow;
  onRemove: (panel: MetricPanel) => void;
}

interface PanelState {
  series: ChartSeries[];
  loading: boolean;
  error: string | null;
}

const EMPTY_STATE: PanelState = { series: [], loading: false, error: null };

// A 24-column grid mirrors the Grafana layout model used across telemetry.
const PanelGrid = styled(Box)(({ theme }) => ({
  display: 'grid',
  gridTemplateColumns: 'repeat(24, 1fr)',
  gridAutoRows: '30px',
  gap: theme.spacing(1),
}));

const Cell = styled(Box)(() => ({
  position: 'relative',
  minWidth: 0,
  '&:hover .panel-actions': { opacity: 1 },
}));

const Actions = styled(Box)(({ theme }) => ({
  position: 'absolute',
  top: theme.spacing(0.5),
  right: theme.spacing(0.5),
  opacity: 0,
  transition: 'opacity 120ms',
  zIndex: 1,
}));

/**
 * MetricsGrid lays saved Prometheus panels out on a 24-column grid and fetches
 * all of their queries in a SINGLE batched request (the backend fans out to
 * Prometheus concurrently), distributing results to presentational panels.
 */
const MetricsGrid: React.FC<MetricsGridProps> = ({
  connectionID,
  panels,
  timeWindow,
  onRemove,
}) => {
  // Resolve PromQL macros ($__rate_interval etc.) once per window. No template
  // variables — Prometheus panels are hand-written queries.
  const queries = useMemo(
    () => panels.map((p) => ({ id: p.id, query: resolveExpr(p.expr, {}, timeWindow) })),
    [panels, timeWindow],
  );

  const batchArg = useMemo(
    () => ({
      connectionID,
      start: String(timeWindow.start),
      end: String(timeWindow.end),
      step: String(timeWindow.step),
      queries,
    }),
    [connectionID, timeWindow, queries],
  );

  const {
    data: batchData,
    isFetching,
    isError,
  } = useQueryPrometheusRangeBatchQuery(batchArg, { skip: panels.length === 0 });

  const resultsById = useMemo(() => {
    const map = new Map<string, { response?: any; error?: string }>();
    for (const r of (batchData as any)?.results ?? []) map.set(r.id, r);
    return map;
  }, [batchData]);

  // Compute each panel's state once per batch result. This keeps `series`
  // referentially stable across unrelated re-renders (so the memoized Panel can
  // skip re-rendering) and avoids re-parsing the response on every render.
  const panelStates = useMemo(() => {
    const states = new Map<string, PanelState>();
    for (const panel of panels) {
      const r = resultsById.get(panel.id);
      if (isError) {
        states.set(panel.id, {
          series: [],
          loading: isFetching,
          error: 'Failed to load panel data',
        });
      } else if (!r) {
        states.set(panel.id, { series: [], loading: isFetching, error: null });
      } else if (r.error) {
        states.set(panel.id, { series: [], loading: isFetching, error: r.error });
      } else {
        states.set(panel.id, {
          series: parsePromMatrix(r.response),
          loading: isFetching,
          error: null,
        });
      }
    }
    return states;
  }, [panels, resultsById, isFetching, isError]);

  return (
    <PanelGrid>
      {panels.map((panel) => {
        const state = panelStates.get(panel.id) ?? EMPTY_STATE;
        return (
          <Cell
            key={panel.id}
            sx={{
              gridColumn: `${(panel.gridPos?.x ?? 0) + 1} / span ${panel.gridPos?.w || 12}`,
              gridRow: `${(panel.gridPos?.y ?? 0) + 1} / span ${panel.gridPos?.h || 8}`,
            }}
          >
            <Actions className="panel-actions">
              <Tooltip title="Remove panel">
                <IconButton size="small" onClick={() => onRemove(panel)} aria-label="Remove panel">
                  <DeleteIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </Actions>
            <Panel
              panel={panel}
              series={state.series}
              loading={state.loading}
              error={state.error}
            />
          </Cell>
        );
      })}
    </PanelGrid>
  );
};

export default MetricsGrid;
