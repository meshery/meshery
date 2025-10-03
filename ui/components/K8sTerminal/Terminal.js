import React, { useEffect, useRef, useState } from 'react';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { Box, IconButton, Tooltip, Paper, Typography } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import RefreshIcon from '@mui/icons-material/Refresh';
import 'xterm/css/xterm.css';

/**
 * K8sTerminal - Interactive terminal component for Kubernetes pod exec
 * 
 * @param {string} namespace - Kubernetes namespace
 * @param {string} podName - Pod name
 * @param {string} containerName - Container name (optional)
 * @param {string} contextId - Kubernetes context ID
 * @param {string} shell - Shell to use (default: /bin/sh)
 * @param {function} onClose - Callback when terminal is closed
 */
const K8sTerminal = ({
  namespace,
  podName,
  containerName = '',
  contextId,
  shell = '/bin/sh',
  onClose,
}) => {
  const terminalRef = useRef(null);
  const terminalInstance = useRef(null);
  const socketRef = useRef(null);
  const fitAddonRef = useRef(null);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isConnected, setIsConnected] = useState(false);

  const connectTerminal = () => {
    if (terminalInstance.current) {
      terminalInstance.current.dispose();
    }

    // Create terminal instance
    const term = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      theme: {
        background: '#1e1e1e',
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
      allowProposedApi: true,
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.loadAddon(new WebLinksAddon());

    term.open(terminalRef.current);
    fitAddon.fit();

    terminalInstance.current = term;
    fitAddonRef.current = fitAddon;

    // Construct WebSocket URL
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    const params = new URLSearchParams({
      namespace,
      pod: podName,
      context: contextId,
      shell,
    });

    if (containerName) {
      params.append('container', containerName);
    }

    const wsUrl = `${protocol}//${host}/api/system/kubernetes/exec?${params.toString()}`;

    // Create WebSocket connection
    const ws = new WebSocket(wsUrl);
    socketRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
      term.writeln('\x1b[1;32mConnecting to pod...\x1b[0m');

      // Send initial terminal size
      ws.send(
        JSON.stringify({
          op: 'resize',
          cols: term.cols,
          rows: term.rows,
        }),
      );
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.op === 'stdout') {
          term.write(msg.data);
        }
      } catch (e) {
        console.error('Error parsing WebSocket message:', e);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      term.writeln('\x1b[1;31m\r\nWebSocket connection error\x1b[0m');
      setIsConnected(false);
    };

    ws.onclose = () => {
      term.writeln('\x1b[1;33m\r\nConnection closed\x1b[0m');
      setIsConnected(false);
    };

    // Send terminal input to WebSocket
    term.onData((data) => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(
          JSON.stringify({
            op: 'stdin',
            data: data,
          }),
        );
      }
    });

    // Handle terminal resize
    const handleResize = () => {
      if (fitAddon && terminalRef.current) {
        fitAddon.fit();
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(
            JSON.stringify({
              op: 'resize',
              cols: term.cols,
              rows: term.rows,
            }),
          );
        }
      }
    };

    // Resize observer
    const resizeObserver = new ResizeObserver(handleResize);
    if (terminalRef.current) {
      resizeObserver.observe(terminalRef.current);
    }

    // Handle window resize
    window.addEventListener('resize', handleResize);

    // Cleanup function
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
    const cleanup = connectTerminal();
    return cleanup;
  }, [namespace, podName, containerName, contextId, shell]);

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
      connectTerminal();
    }, 100);
  };

  const handleClose = () => {
    if (socketRef.current) {
      socketRef.current.close();
    }
    if (onClose) {
      onClose();
    }
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
        backgroundColor: '#1e1e1e',
      }}
    >
      {/* Terminal Header */}
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: '8px 16px',
          backgroundColor: '#2d2d2d',
          borderBottom: '1px solid #3e3e3e',
        }}
      >
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

        <Box sx={{ display: 'flex', gap: 1 }}>
          <Tooltip title="Refresh Connection">
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

          <Tooltip title="Close Terminal">
            <IconButton size="small" onClick={handleClose} sx={{ color: '#e5e5e5' }}>
              <CloseIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>

      {/* Terminal Container */}
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

export default K8sTerminal;

