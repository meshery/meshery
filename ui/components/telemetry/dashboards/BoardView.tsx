import React, { useMemo, useState } from 'react';
import {
  Box,
  Chip,
  CircularProgress,
  Collapse,
  DeleteIcon,
  ExpandMoreIcon,
  IconButton,
  Tooltip,
  Typography,
  WarningIcon,
  styled,
  useTheme,
} from '@sistent/sistent';
import {
  useGetGrafanaBoardQuery,
  useGetGrafanaDatasourcesQuery,
  useQueryGrafanaRangeBatchQuery,
} from '@/rtk-query/telemetryGrafana';
import Panel from '../common/Panel';
import { parsePromMatrix } from '../common/time';
import { buildVarValues, resolveDatasourceUid, resolveExpr } from './resolve';
import type { Board, ChartSeries, Datasource, PinnedBoard, TimeWindow } from './types';

interface BoardViewProps {
  connectionID: string;
  board: PinnedBoard;
  timeWindow: TimeWindow;
  onRemove: (board: PinnedBoard) => void;
}

// One resolved query for the batch request, plus the client-side context needed
// to route its result back to a panel (panelId) and label it (legendFormat).
interface ResolvedQuery {
  id: string;
  panelId: number;
  legendFormat?: string;
  ds: string;
  query: string;
}

interface PanelState {
  series: ChartSeries[];
  loading: boolean;
  error: string | null;
}

const Section = styled(Box)(({ theme }) => ({
  border: `1px solid ${theme.palette.border.default}`,
  borderRadius: theme.shape.borderRadius,
  background: theme.palette.background.surfaces,
  overflow: 'hidden',
}));

const Header = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1),
  padding: theme.spacing(1, 2),
  cursor: 'pointer',
  background: theme.palette.background.default,
  borderBottom: `1px solid ${theme.palette.border.default}`,
}));

// A 24-column grid mirrors Grafana's panel layout model.
const PanelGrid = styled(Box)(({ theme }) => ({
  display: 'grid',
  gridTemplateColumns: 'repeat(24, 1fr)',
  gridAutoRows: '30px',
  gap: theme.spacing(1),
  padding: theme.spacing(2),
}));

const ExpandButton = styled(IconButton, {
  shouldForwardProp: (prop) => prop !== 'open',
})<{ open: boolean }>(({ open, theme }) => ({
  transform: open ? 'rotate(0deg)' : 'rotate(-90deg)',
  transition: theme.transitions.create('transform', { duration: 150 }),
}));

/**
 * BoardView loads a Grafana dashboard, resolves every panel's queries, and
 * fetches them all in a SINGLE batched request (the backend fans out to Grafana
 * concurrently). Results are distributed to each presentational Panel, so a
 * board renders with one round trip instead of one request per panel target.
 */
const BoardView: React.FC<BoardViewProps> = ({ connectionID, board, timeWindow, onRemove }) => {
  const theme = useTheme();
  const [open, setOpen] = useState(true);
  const { data, isFetching, isError } = useGetGrafanaBoardQuery(
    { connectionID, uid: board.uid },
    { skip: !connectionID || !board.uid },
  );
  const fullBoard = data as Board | undefined;

  // Datasources resolve `$datasource` variables and datasource names to concrete
  // UIDs; cached per connection by RTK Query. Wait for them before querying so
  // references resolve correctly on the first request.
  const { data: datasourcesData, isLoading: datasourcesLoading } = useGetGrafanaDatasourcesQuery(
    { connectionID },
    { skip: !connectionID },
  );
  const datasources = (datasourcesData as Datasource[] | undefined) ?? [];

  // Resolve every panel target into an executable query once per board/window.
  const resolved = useMemo<ResolvedQuery[]>(() => {
    if (!fullBoard) return [];
    const varValues = buildVarValues(fullBoard.templateVars);
    const out: ResolvedQuery[] = [];
    for (const panel of fullBoard.panels) {
      (panel.targets ?? []).forEach((t, idx) => {
        if (!t.expr || t.expr.trim() === '') return;
        out.push({
          id: `${panel.id}:${t.refId || idx}`,
          panelId: panel.id,
          legendFormat: t.legendFormat,
          ds: resolveDatasourceUid(t.datasourceUid, fullBoard.templateVars, datasources),
          query: resolveExpr(t.expr, varValues, timeWindow),
        });
      });
    }
    return out;
  }, [fullBoard, datasources, timeWindow]);

  const batchArg = useMemo(
    () => ({
      connectionID,
      start: String(timeWindow.start),
      end: String(timeWindow.end),
      step: String(timeWindow.step),
      queries: resolved.map((q) => ({ id: q.id, ds: q.ds, query: q.query })),
    }),
    [connectionID, timeWindow, resolved],
  );

  const {
    data: batchData,
    isFetching: batchFetching,
    isError: batchError,
  } = useQueryGrafanaRangeBatchQuery(batchArg, {
    skip: datasourcesLoading || resolved.length === 0,
  });

  // Index batch results by query id for distribution to panels.
  const resultsById = useMemo(() => {
    const map = new Map<string, { response?: any; error?: string }>();
    for (const r of (batchData as any)?.results ?? []) map.set(r.id, r);
    return map;
  }, [batchData]);

  // Compute each panel's series / loading / error from the shared batch result.
  const panelStates = useMemo(() => {
    const states = new Map<number, PanelState>();
    if (!fullBoard) return states;
    for (const panel of fullBoard.panels) {
      const targets = resolved.filter((q) => q.panelId === panel.id);
      if (targets.length === 0) {
        states.set(panel.id, { series: [], loading: false, error: null });
        continue;
      }
      const series: ChartSeries[] = [];
      let error: string | null = batchError ? 'Failed to load panel data' : null;
      for (const t of targets) {
        const r = resultsById.get(t.id);
        if (!r) continue;
        if (r.error) {
          error = error ?? r.error;
          continue;
        }
        series.push(...parsePromMatrix(r.response, t.legendFormat));
      }
      states.set(panel.id, { series, loading: batchFetching, error });
    }
    return states;
  }, [fullBoard, resolved, resultsById, batchFetching, batchError]);

  return (
    <Section>
      <Header onClick={() => setOpen((o) => !o)}>
        <ExpandButton open={open} size="small" aria-label={open ? 'collapse' : 'expand'}>
          <ExpandMoreIcon />
        </ExpandButton>
        <Typography variant="subtitle1" sx={{ fontWeight: 600, flex: 1 }} noWrap>
          {fullBoard?.title || board.title || board.uid}
        </Typography>
        {(fullBoard?.tags ?? []).slice(0, 4).map((tag) => (
          <Chip key={tag} label={tag} size="small" variant="outlined" />
        ))}
        {(isFetching || batchFetching) && <CircularProgress size={16} />}
        {isError && (
          <Tooltip title="Failed to load this dashboard">
            <WarningIcon style={{ color: theme.palette.warning.main }} />
          </Tooltip>
        )}
        <Tooltip title="Remove from telemetry">
          <IconButton
            size="small"
            onClick={(e) => {
              e.stopPropagation();
              onRemove(board);
            }}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </Header>

      <Collapse in={open} unmountOnExit>
        {isError ? (
          <Box sx={{ p: 3 }}>
            <Typography color="textSecondary">
              Could not load this dashboard from Grafana. It may have been deleted or the connection
              may be unreachable.
            </Typography>
          </Box>
        ) : !fullBoard ? (
          <Box sx={{ p: 4, display: 'flex', justifyContent: 'center' }}>
            <CircularProgress size={26} />
          </Box>
        ) : fullBoard.panels.length === 0 ? (
          <Box sx={{ p: 3 }}>
            <Typography color="textSecondary">This dashboard has no renderable panels.</Typography>
          </Box>
        ) : (
          <PanelGrid>
            {fullBoard.panels.map((panel) => {
              const state = panelStates.get(panel.id) ?? {
                series: [],
                loading: batchFetching,
                error: null,
              };
              return (
                <Box
                  key={panel.id}
                  sx={{
                    gridColumn: `${(panel.gridPos?.x ?? 0) + 1} / span ${panel.gridPos?.w || 24}`,
                    gridRow: `${(panel.gridPos?.y ?? 0) + 1} / span ${panel.gridPos?.h || 8}`,
                    minWidth: 0,
                  }}
                >
                  <Panel
                    panel={panel}
                    series={state.series}
                    loading={state.loading}
                    error={state.error}
                  />
                </Box>
              );
            })}
          </PanelGrid>
        )}
      </Collapse>
    </Section>
  );
};

export default BoardView;
