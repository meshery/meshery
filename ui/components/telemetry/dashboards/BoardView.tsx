import React, { useState } from 'react';
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
} from '@/rtk-query/telemetryGrafana';
import Panel from './Panel';
import type { Board, Datasource, PinnedBoard, TimeWindow } from './types';

interface BoardViewProps {
  connectionID: string;
  board: PinnedBoard;
  timeWindow: TimeWindow;
  onRemove: (board: PinnedBoard) => void;
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
 * BoardView loads a Grafana dashboard and lays its panels out on a 24-column
 * grid that mirrors Grafana's own layout, with each panel fetching its own data.
 */
const BoardView: React.FC<BoardViewProps> = ({ connectionID, board, timeWindow, onRemove }) => {
  const theme = useTheme();
  const [open, setOpen] = useState(true);
  const { data, isFetching, isError } = useGetGrafanaBoardQuery(
    { connectionID, uid: board.uid },
    { skip: !connectionID || !board.uid },
  );
  const fullBoard = data as Board | undefined;

  // Datasources are needed to resolve `$datasource` template variables to a
  // concrete UID; cached per connection by RTK Query, so this is shared.
  const { data: datasourcesData } = useGetGrafanaDatasourcesQuery(
    { connectionID },
    { skip: !connectionID },
  );
  const datasources = (datasourcesData as Datasource[] | undefined) ?? [];

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
        {isFetching && <CircularProgress size={16} />}
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
            {fullBoard.panels.map((panel) => (
              <Box
                key={panel.id}
                sx={{
                  gridColumn: `${(panel.gridPos?.x ?? 0) + 1} / span ${panel.gridPos?.w || 24}`,
                  gridRow: `${(panel.gridPos?.y ?? 0) + 1} / span ${panel.gridPos?.h || 8}`,
                  minWidth: 0,
                }}
              >
                <Panel
                  connectionID={connectionID}
                  panel={panel}
                  timeWindow={timeWindow}
                  templateVars={fullBoard.templateVars}
                  datasources={datasources}
                />
              </Box>
            ))}
          </PanelGrid>
        )}
      </Collapse>
    </Section>
  );
};

export default BoardView;
