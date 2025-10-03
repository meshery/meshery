import React, { useEffect, useRef, useState } from 'react';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { SearchAddon } from 'xterm-addon-search';
import {
  Box,
  IconButton,
  Tooltip,
  Paper,
  Typography,
  Switch,
  FormControlLabel,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import RefreshIcon from '@mui/icons-material/Refresh';
import SearchIcon from '@mui/icons-material/Search';
import 'xterm/css/xterm.css';

/**
 * K8sLogsViewer - Log viewer component for Kubernetes pods
 * 
 * @param {string} namespace - Kubernetes namespace
 * @param {string} podName - Pod name
 * @param {string} containerName - Container name (optional)
 * @param {string} contextId - Kubernetes context ID
 * @param {boolean} follow - Follow log output (default: true)
 * @param {boolean} previous - Show logs from previous container (default: false)
 * @param {function} onClose - Callback when viewer is closed
 */
const K8sLogsViewer = ({
  namespace,
  podName,
  containerName = '',
  contextId,
  follow: initialFollow = true,
  previous: initialPrevious = false,
  onClose,
}) => {
  const terminalRef = useRef(null);
  const terminalInstance = useRef(null);
  const socketRef = useRef(null);
  const fitAddonRef = useRef(null);
  const searchAddonRef = useRef(null);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isConnected, setIsConnected] = useState(false);
  const [follow, setFollow] = useState(initialFollow);
  const [showPrevious, setShowPrevious] = useState(initialPrevious);

  const connectLogs = (followLogs, previousLogs) => {
    if (terminalInstance.current) {
      terminalInstance.current.dispose();
    }

    // Create terminal instance (read-only for logs)
    const term = new Terminal({
      cursorBlink: false,
      fontSize: 12,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      theme: {
        background: '#0a0a0a',
        foreground: '#d4d4d4',
        cursor: '#ffffff',
        selection: 'rgba(255, 255, 255, 0.3)',
        black: '#000000',
        red: '#cd3131',
        green: '#0dbc79',
        yellow: '#e5e510',
        blue: '#2472c8',
        magenta: '#bc3fbc',
        cyan: '#11a8cd',
        white: '#e5e5e5',
        brightBlack: '#666666',
        brightRed: '#f14c4c',
        brightGreen: '#23d18b',
        brightYellow: '#f5f543',
        brightBlue: '#3b8eea',
        brightMagenta: '#d670d6',
        brightCyan: '#29b8db',
        brightWhite: '#e5e5e5',
      },
      convertEol: true,
      disableStdin: true,
      allowProposedApi: true,
    });

    const fitAddon = new FitAddon();
    const searchAddon = new SearchAddon();
    term.loadAddon(fitAddon);
    term.loadAddon(new WebLinksAddon());
    term.loadAddon(searchAddon);

    term.open(terminalRef.current);
    fitAddon.fit();

    terminalInstance.current = term;
    fitAddonRef.current = fitAddon;
    searchAddonRef.current = searchAddon;

    // Construct WebSocket URL
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    const params = new URLSearchParams({
      namespace,
      pod: podName,
      context: contextId,
      follow: followLogs.toString(),
      previous: previousLogs.toString(),
    });

    if (containerName) {
      params.append('container', containerName);
    }

    const wsUrl = `${protocol}//${host}/api/system/kubernetes/logs?${params.toString()}`;

    // Create WebSocket connection
    const ws = new WebSocket(wsUrl);
    socketRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
      term.writeln('\x1b[1;36m=== Connecting to log stream ===\x1b[0m');
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.op === 'stdout') {
          term.write(msg.data);
        }
      } catch (e) {
        // If not JSON, just write the raw data
        term.write(event.data);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      term.writeln('\x1b[1;31m\r\nWebSocket connection error\x1b[0m');
      setIsConnected(false);
    };

    ws.onclose = () => {
      term.writeln('\x1b[1;33m\r\n=== Log stream closed ===\x1b[0m');
      setIsConnected(false);
    };

    // Handle terminal resize
    const handleResize = () => {
      if (fitAddon && terminalRef.current) {
        fitAddon.fit();
      }
    };

    const resizeObserver = new ResizeObserver(handleResize);
    if (terminalRef.current) {
      resizeObserver.observe(terminalRef.current);
    }

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      resizeObserver.disconnect();
      if (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) {
        ws.close();
      }
      term.dispose();
    };
  };

  useEffect(() => {
    const cleanup = connectLogs(follow, showPrevious);
    return cleanup;
  }, [namespace, podName, containerName, contextId, follow, showPrevious]);

  const handleFullscreen = () => {
    setIsFullscreen(!isFullscreen);
    setTimeout(() => {
      if (fitAddonRef.current) {
        fitAddonRef.current.fit();
      }
    }, 100);
  };

  const handleRefresh = () => {
    if (socketRef.current) {
      socketRef.current.close();
    }
    setTimeout(() => {
      connectLogs(follow, showPrevious);
    }, 100);
  };

  const handleSearch = () => {
    if (searchAddonRef.current && terminalInstance.current) {
      const searchTerm = prompt('Search logs:');
      if (searchTerm) {
        searchAddonRef.current.findNext(searchTerm);
      }
    }
  };

  const handleClose = () => {
    if (socketRef.current) {
      socketRef.current.close();
    }
    if (onClose) {
      onClose();
    }
  };

  const handleFollowChange = (event) => {
    setFollow(event.target.checked);
  };

  const handlePreviousChange = (event) => {
    setShowPrevious(event.target.checked);
  };

  return (
    <Paper
      elevation={3}
      sx={{
        position: isFullscreen ? 'fixed' : 'relative',
        top: isFullscreen ? 0 : 'auto',
        left: isFullscreen ? 0 : 'auto',
        right: isFullscreen ? 0 : 'auto',
        bottom: isFullscreen ? 0 : 'auto',
        width: isFullscreen ? '100vw' : '100%',
        height: isFullscreen ? '100vh' : '600px',
        zIndex: isFullscreen ? 9999 : 'auto',
        display: 'flex',
        flexDirection: 'column',
        backgroundColor: '#0a0a0a',
      }}
    >
      {/* Logs Header */}
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: '8px 16px',
          backgroundColor: '#1a1a1a',
          borderBottom: '1px solid #2a2a2a',
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Box
              sx={{
                width: 12,
                height: 12,
                borderRadius: '50%',
                backgroundColor: isConnected ? '#0dbc79' : '#cd3131',
              }}
            />
            <Typography variant="body2" sx={{ color: '#e5e5e5' }}>
              {namespace}/{podName}
              {containerName && ` : ${containerName}`}
            </Typography>
          </Box>

          <FormControlLabel
            control={
              <Switch
                size="small"
                checked={follow}
                onChange={handleFollowChange}
                sx={{
                  '& .MuiSwitch-switchBase.Mui-checked': {
                    color: '#0dbc79',
                  },
                }}
              />
            }
            label={<Typography variant="caption" sx={{ color: '#e5e5e5' }}>Follow</Typography>}
          />

          <FormControlLabel
            control={
              <Switch
                size="small"
                checked={showPrevious}
                onChange={handlePreviousChange}
                sx={{
                  '& .MuiSwitch-switchBase.Mui-checked': {
                    color: '#2472c8',
                  },
                }}
              />
            }
            label={<Typography variant="caption" sx={{ color: '#e5e5e5' }}>Previous</Typography>}
          />
        </Box>

        <Box sx={{ display: 'flex', gap: 1 }}>
          <Tooltip title="Search Logs">
            <IconButton size="small" onClick={handleSearch} sx={{ color: '#e5e5e5' }}>
              <SearchIcon fontSize="small" />
            </IconButton>
          </Tooltip>

          <Tooltip title="Refresh Logs">
            <IconButton size="small" onClick={handleRefresh} sx={{ color: '#e5e5e5' }}>
              <RefreshIcon fontSize="small" />
            </IconButton>
          </Tooltip>

          <Tooltip title={isFullscreen ? 'Exit Fullscreen' : 'Fullscreen'}>
            <IconButton size="small" onClick={handleFullscreen} sx={{ color: '#e5e5e5' }}>
              {isFullscreen ? (
                <FullscreenExitIcon fontSize="small" />
              ) : (
                <FullscreenIcon fontSize="small" />
              )}
            </IconButton>
          </Tooltip>

          <Tooltip title="Close Logs">
            <IconButton size="small" onClick={handleClose} sx={{ color: '#e5e5e5' }}>
              <CloseIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>

      {/* Logs Container */}
      <Box
        ref={terminalRef}
        sx={{
          flex: 1,
          overflow: 'hidden',
          padding: '8px',
        }}
      />
    </Paper>
  );
};

export default K8sLogsViewer;

