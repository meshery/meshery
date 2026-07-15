import React, { useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/router';
import {
  AddIcon,
  Box,
  Button,
  CircularProgress,
  Drawer,
  InsertChartIcon,
  Typography,
  styled,
  useTheme,
} from '@sistent/sistent';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
import {
  useGetPrometheusPanelsQuery,
  usePingPrometheusConnectionQuery,
  useUpdatePrometheusPanelsMutation,
} from '@/rtk-query/telemetryPrometheus';
import ConnectionPicker, { TelemetryConnection } from '../common/ConnectionPicker';
import PingStatus from '../common/PingStatus';
import TimeRangePicker from '../common/TimeRangePicker';
import RefreshControl from '../common/RefreshControl';
import { DEFAULT_RANGE, resolveWindow } from '../common/time';
import MetricsGrid from './MetricsGrid';
import MetricExplorer from './MetricExplorer';
import type { MetricPanel } from './types';

const Toolbar = styled(Box)(({ theme }) => ({
  position: 'sticky',
  top: 0,
  zIndex: 2,
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1.5),
  flexWrap: 'wrap',
  padding: theme.spacing(1.5, 2),
  background: theme.palette.background.default,
  borderBottom: `1px solid ${theme.palette.border.default}`,
}));

const Content = styled(Box)(({ theme }) => ({
  padding: theme.spacing(2),
}));

const Centered = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  gap: theme.spacing(2),
  textAlign: 'center',
  padding: theme.spacing(8, 2),
  color: theme.palette.text.secondary,
}));

const DrawerBody = styled(Box)(({ theme }) => ({
  width: 460,
  maxWidth: '90vw',
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  padding: theme.spacing(3),
  gap: theme.spacing(1),
}));

// Assigns a default 12-wide, 8-tall slot, flowing two panels per row.
const nextGridPos = (count: number) => ({
  x: (count % 2) * 12,
  y: Math.floor(count / 2) * 8,
  w: 12,
  h: 8,
});

const newId = () =>
  typeof crypto !== 'undefined' && 'randomUUID' in crypto
    ? crypto.randomUUID()
    : `${Date.now()}-${Math.floor(Math.random() * 1e6)}`;

/**
 * TelemetryMetrics is the root of the Prometheus-native telemetry experience. A
 * user picks a Prometheus connection, explores its metrics, writes PromQL, and
 * saves panels that render live on a shared time range — all backed by a single
 * batched query per refresh.
 */
const TelemetryMetrics: React.FC = () => {
  const theme = useTheme();
  const router = useRouter();

  const { data: connectionsData, isLoading: connectionsLoading } = useGetConnectionsQuery({
    kind: JSON.stringify(['prometheus']),
    pagesize: 200,
  });

  const connections: TelemetryConnection[] = useMemo(
    () =>
      ((connectionsData as any)?.connections ?? [])
        .filter((c: any) => c.kind === 'prometheus' && c.status !== 'deleted')
        .map((c: any) => ({ id: c.id, name: c.name, kind: c.kind, metadata: c.metadata })),
    [connectionsData],
  );

  const [connectionID, setConnectionID] = useState('');
  const [durationSec, setDurationSec] = useState(DEFAULT_RANGE.durationSec);
  const [refreshMs, setRefreshMs] = useState(0);
  const [refreshSignal, setRefreshSignal] = useState(0);
  const [explorerOpen, setExplorerOpen] = useState(false);

  useEffect(() => {
    if (connections.length === 0) {
      if (connectionID) setConnectionID('');
      return;
    }
    if (!connections.some((c) => c.id === connectionID)) {
      setConnectionID(connections[0].id);
    }
  }, [connections, connectionID]);

  const selectedConnection = connections.find((c) => c.id === connectionID);

  const { data: panelsData } = useGetPrometheusPanelsQuery(
    { connectionID },
    { skip: !connectionID },
  );
  const panels = (panelsData as MetricPanel[] | undefined) ?? [];
  const [updatePanels] = useUpdatePrometheusPanelsMutation();

  // Fresh absolute window whenever range or refresh signal changes.
  const timeWindow = useMemo(() => resolveWindow(durationSec), [durationSec, refreshSignal]);

  const persist = async (next: MetricPanel[]) => {
    try {
      await updatePanels({ connectionID, panels: next }).unwrap();
    } catch {
      /* invalidation will resync */
    }
  };

  const handleAdd = (draft: Omit<MetricPanel, 'id' | 'gridPos'>) => {
    const panel: MetricPanel = { ...draft, id: newId(), gridPos: nextGridPos(panels.length) };
    persist([...panels, panel]);
    setExplorerOpen(false);
  };

  const handleRemove = (panel: MetricPanel) => persist(panels.filter((p) => p.id !== panel.id));

  if (connectionsLoading) {
    return (
      <Centered>
        <CircularProgress />
      </Centered>
    );
  }

  if (connections.length === 0) {
    return (
      <Centered data-testid="telemetry-prometheus-empty">
        <InsertChartIcon style={{ width: 56, height: 56, fill: theme.palette.icon.secondary }} />
        <Box>
          <Typography variant="h6" gutterBottom>
            No Prometheus connections yet
          </Typography>
          <Typography color="textSecondary" sx={{ maxWidth: 460 }}>
            Add a Prometheus connection to explore its metrics and build panels here. Connections
            are managed from the Connections page.
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => router.push('/management/connections')}
        >
          Add a Prometheus connection
        </Button>
      </Centered>
    );
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <Toolbar data-testid="telemetry-prometheus-toolbar">
        <ConnectionPicker
          connections={connections}
          value={connectionID}
          onChange={setConnectionID}
          label="Prometheus connection"
        />
        {connectionID && (
          <PingStatus
            connectionID={connectionID}
            usePingQuery={usePingPrometheusConnectionQuery}
            productName="Prometheus"
          />
        )}
        <Box sx={{ flex: 1 }} />
        <Button
          variant="outlined"
          startIcon={<AddIcon />}
          onClick={() => setExplorerOpen(true)}
          disabled={!connectionID}
        >
          Add metric
        </Button>
        <TimeRangePicker durationSec={durationSec} onChange={setDurationSec} />
        <RefreshControl
          intervalMs={refreshMs}
          onIntervalChange={setRefreshMs}
          onRefresh={() => setRefreshSignal((s) => s + 1)}
        />
      </Toolbar>

      {panels.length === 0 ? (
        <Centered>
          <InsertChartIcon style={{ width: 56, height: 56, fill: theme.palette.icon.secondary }} />
          <Box>
            <Typography variant="h6" gutterBottom>
              No panels yet
            </Typography>
            <Typography color="textSecondary" sx={{ maxWidth: 460 }}>
              Explore this Prometheus instance&apos;s metrics and save PromQL queries as panels to
              monitor them from Meshery.
            </Typography>
          </Box>
          <Button variant="contained" startIcon={<AddIcon />} onClick={() => setExplorerOpen(true)}>
            Explore metrics
          </Button>
        </Centered>
      ) : (
        <Content>
          <MetricsGrid
            connectionID={connectionID}
            panels={panels}
            timeWindow={timeWindow}
            onRemove={handleRemove}
          />
        </Content>
      )}

      <Drawer anchor="right" open={explorerOpen} onClose={() => setExplorerOpen(false)}>
        <DrawerBody>
          <Typography variant="h6">Explore metrics</Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
            {selectedConnection?.name
              ? `Browsing ${selectedConnection.name}`
              : 'Browse and add metrics'}
          </Typography>
          {connectionID && (
            <MetricExplorer connectionID={connectionID} timeWindow={timeWindow} onAdd={handleAdd} />
          )}
        </DrawerBody>
      </Drawer>
    </Box>
  );
};

export default TelemetryMetrics;
