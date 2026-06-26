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
  useGetPinnedBoardsQuery,
  usePingGrafanaConnectionQuery,
  useUpdatePinnedBoardsMutation,
} from '@/rtk-query/telemetryGrafana';
import ConnectionPicker, { TelemetryConnection } from '../common/ConnectionPicker';
import PingStatus from '../common/PingStatus';
import TimeRangePicker from '../common/TimeRangePicker';
import RefreshControl from '../common/RefreshControl';
import { DEFAULT_RANGE, resolveWindow } from '../common/time';
import BoardLibrary from './BoardLibrary';
import BoardView from './BoardView';
import type { PinnedBoard } from './types';

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
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(2),
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

/**
 * TelemetryDashboards is the root of the Telemetry section. It lets a user pick
 * a registered Grafana connection, browse and add its dashboards, and render
 * the added dashboards' panels live with a shared time range and refresh.
 */
const TelemetryDashboards: React.FC = () => {
  const theme = useTheme();
  const router = useRouter();

  const { data: connectionsData, isLoading: connectionsLoading } = useGetConnectionsQuery({
    kind: JSON.stringify(['grafana']),
    pagesize: 200,
  });

  const connections: TelemetryConnection[] = useMemo(
    () =>
      ((connectionsData as any)?.connections ?? [])
        .filter((c: any) => c.kind === 'grafana' && c.status !== 'deleted')
        .map((c: any) => ({ id: c.id, name: c.name, kind: c.kind, metadata: c.metadata })),
    [connectionsData],
  );

  const [connectionID, setConnectionID] = useState('');
  const [durationSec, setDurationSec] = useState(DEFAULT_RANGE.durationSec);
  const [refreshMs, setRefreshMs] = useState(0);
  const [refreshSignal, setRefreshSignal] = useState(0);
  const [libraryOpen, setLibraryOpen] = useState(false);

  // Default the selection to the first connection once loaded, and keep it valid.
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
  const baseURL = selectedConnection?.metadata?.url as string | undefined;

  const { data: pinnedData } = useGetPinnedBoardsQuery({ connectionID }, { skip: !connectionID });
  const pinned = (pinnedData as PinnedBoard[] | undefined) ?? [];
  const [updatePinnedBoards] = useUpdatePinnedBoardsMutation();

  // Resolve a fresh absolute timeWindow whenever the range or refresh signal changes,
  // so every panel re-queries against an aligned, current timeWindow.
  const timeWindow = useMemo(
    () => resolveWindow(durationSec),

    [durationSec, refreshSignal],
  );

  const handleTogglePin = async (board: PinnedBoard, shouldPin: boolean) => {
    const next = shouldPin
      ? [...pinned.filter((b) => b.uid !== board.uid), board]
      : pinned.filter((b) => b.uid !== board.uid);
    try {
      await updatePinnedBoards({ connectionID, boards: next }).unwrap();
    } catch {
      /* invalidation will resync; surfaced via the library's own states */
    }
  };

  if (connectionsLoading) {
    return (
      <Centered>
        <CircularProgress />
      </Centered>
    );
  }

  if (connections.length === 0) {
    return (
      <Centered data-testid="telemetry-grafana-empty">
        <InsertChartIcon style={{ width: 56, height: 56, fill: theme.palette.icon.secondary }} />
        <Box>
          <Typography variant="h6" gutterBottom>
            No Grafana connections yet
          </Typography>
          <Typography color="textSecondary" sx={{ maxWidth: 460 }}>
            Add a Grafana connection to browse and render its dashboards here. Connections are
            managed from the Connections page.
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => router.push('/management/connections')}
        >
          Add a Grafana connection
        </Button>
      </Centered>
    );
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <Toolbar data-testid="telemetry-grafana-toolbar">
        <ConnectionPicker
          connections={connections}
          value={connectionID}
          onChange={setConnectionID}
          label="Grafana connection"
        />
        {connectionID && (
          <PingStatus
            connectionID={connectionID}
            usePingQuery={usePingGrafanaConnectionQuery}
            productName="Grafana"
          />
        )}
        <Box sx={{ flex: 1 }} />
        <Button
          variant="outlined"
          startIcon={<AddIcon />}
          onClick={() => setLibraryOpen(true)}
          disabled={!connectionID}
        >
          Add dashboards
        </Button>
        <TimeRangePicker durationSec={durationSec} onChange={setDurationSec} />
        <RefreshControl
          intervalMs={refreshMs}
          onIntervalChange={setRefreshMs}
          onRefresh={() => setRefreshSignal((s) => s + 1)}
        />
      </Toolbar>

      {pinned.length === 0 ? (
        <Centered>
          <InsertChartIcon style={{ width: 56, height: 56, fill: theme.palette.icon.secondary }} />
          <Box>
            <Typography variant="h6" gutterBottom>
              No dashboards added
            </Typography>
            <Typography color="textSecondary" sx={{ maxWidth: 460 }}>
              Browse this Grafana instance and add the dashboards you want to monitor from Meshery.
            </Typography>
          </Box>
          <Button variant="contained" startIcon={<AddIcon />} onClick={() => setLibraryOpen(true)}>
            Browse dashboards
          </Button>
        </Centered>
      ) : (
        <Content>
          {pinned.map((board) => (
            <BoardView
              key={board.uid}
              connectionID={connectionID}
              board={board}
              timeWindow={timeWindow}
              onRemove={(b) => handleTogglePin(b, false)}
            />
          ))}
        </Content>
      )}

      <Drawer anchor="right" open={libraryOpen} onClose={() => setLibraryOpen(false)}>
        <DrawerBody>
          <Typography variant="h6">Add dashboards</Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
            {selectedConnection?.name
              ? `Browsing ${selectedConnection.name}`
              : 'Browse and add dashboards'}
          </Typography>
          {connectionID && (
            <BoardLibrary
              connectionID={connectionID}
              baseURL={baseURL}
              pinned={pinned}
              onTogglePin={handleTogglePin}
            />
          )}
        </DrawerBody>
      </Drawer>
    </Box>
  );
};

export default TelemetryDashboards;
