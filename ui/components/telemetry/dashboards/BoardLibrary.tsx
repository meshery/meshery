import React, { useEffect, useMemo, useState } from 'react';
import {
  AddIcon,
  Box,
  Chip,
  CircularProgress,
  DoneIcon,
  IconButton,
  InputAdornment,
  LaunchIcon,
  SearchIcon,
  TextField,
  Tooltip,
  Typography,
  WarningIcon,
  styled,
  useTheme,
} from '@sistent/sistent';
import { useSearchGrafanaBoardsQuery } from '@/rtk-query/telemetryGrafana';
import type { BoardSummary, PinnedBoard } from './types';

interface BoardLibraryProps {
  connectionID: string;
  baseURL?: string;
  pinned: PinnedBoard[];
  onTogglePin: (board: PinnedBoard, shouldPin: boolean) => void;
}

const Row = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1),
  padding: theme.spacing(1, 1.5),
  borderRadius: theme.shape.borderRadius,
  border: `1px solid ${theme.palette.border.default}`,
  '&:hover': { background: theme.palette.background.hover ?? theme.palette.action?.hover },
}));

const List = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(1),
  marginTop: theme.spacing(2),
  overflowY: 'auto',
}));

/**
 * BoardLibrary lets the user search a Grafana instance's dashboards and add or
 * remove them from the connection's pinned set.
 */
const BoardLibrary: React.FC<BoardLibraryProps> = ({
  connectionID,
  baseURL,
  pinned,
  onTogglePin,
}) => {
  const theme = useTheme();
  const [input, setInput] = useState('');
  const [search, setSearch] = useState('');

  // Debounce the search input so we don't query Grafana on every keystroke.
  useEffect(() => {
    const t = setTimeout(() => setSearch(input.trim()), 350);
    return () => clearTimeout(t);
  }, [input]);

  const { data, isFetching, isError } = useSearchGrafanaBoardsQuery(
    { connectionID, search },
    { skip: !connectionID },
  );
  const boards = (data as BoardSummary[] | undefined) ?? [];
  const pinnedUIDs = useMemo(() => new Set(pinned.map((b) => b.uid)), [pinned]);

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <TextField
        autoFocus
        fullWidth
        size="small"
        placeholder="Search dashboards…"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon />
            </InputAdornment>
          ),
          endAdornment: isFetching ? <CircularProgress size={16} /> : null,
        }}
      />

      {isError ? (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 3 }}>
          <WarningIcon style={{ color: theme.palette.warning.main }} />
          <Typography color="textSecondary">
            Could not reach Grafana. Check the connection and its credential.
          </Typography>
        </Box>
      ) : !isFetching && boards.length === 0 ? (
        <Typography color="textSecondary" sx={{ mt: 3 }}>
          {search ? `No dashboards match “${search}”.` : 'No dashboards found on this Grafana.'}
        </Typography>
      ) : (
        <List>
          {boards.map((board) => {
            const isPinned = pinnedUIDs.has(board.uid);
            return (
              <Row key={board.uid}>
                <Box sx={{ flex: 1, minWidth: 0 }}>
                  <Typography variant="body2" noWrap sx={{ fontWeight: 500 }}>
                    {board.title}
                  </Typography>
                  <Box sx={{ display: 'flex', gap: 0.5, mt: 0.5, flexWrap: 'wrap' }}>
                    {(board.tags ?? []).slice(0, 5).map((tag) => (
                      <Chip key={tag} label={tag} size="small" variant="outlined" />
                    ))}
                  </Box>
                </Box>
                {baseURL && board.url && (
                  <Tooltip title="Open in Grafana">
                    <IconButton
                      size="small"
                      aria-label="Open dashboard in Grafana"
                      component="a"
                      href={`${baseURL.replace(/\/$/, '')}${board.url}`}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <LaunchIcon fontSize="small" />
                    </IconButton>
                  </Tooltip>
                )}
                <Tooltip title={isPinned ? 'Remove from telemetry' : 'Add to telemetry'}>
                  <IconButton
                    size="small"
                    color={isPinned ? 'primary' : 'default'}
                    onClick={() => onTogglePin({ uid: board.uid, title: board.title }, !isPinned)}
                  >
                    {isPinned ? <DoneIcon /> : <AddIcon />}
                  </IconButton>
                </Tooltip>
              </Row>
            );
          })}
        </List>
      )}
    </Box>
  );
};

export default BoardLibrary;
