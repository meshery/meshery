import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Autocomplete,
  Box,
  Typography,
  CircularProgress,
  Alert,
  Tabs,
  Tab,
} from '@mui/material';
import TerminalIcon from '@mui/icons-material/Terminal';
import DescriptionIcon from '@mui/icons-material/Description';
import { graphql, requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../lib/relayEnvironment';
import K8sTerminal from './Terminal';
import K8sLogsViewer from './LogsViewer';

/**
 * PodSelector - Dialog for selecting pod and opening terminal/logs
 * 
 * @param {boolean} open - Whether dialog is open
 * @param {function} onClose - Callback when dialog is closed
 * @param {string} contextId - Kubernetes context ID
 */
const PodSelector = ({ open, onClose, contextId }) => {
  const [namespaces, setNamespaces] = useState([]);
  const [pods, setPods] = useState([]);
  const [containers, setContainers] = useState([]);
  const [selectedNamespace, setSelectedNamespace] = useState(null);
  const [selectedPod, setSelectedPod] = useState(null);
  const [selectedContainer, setSelectedContainer] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [tabValue, setTabValue] = useState(0); // 0 for terminal, 1 for logs
  const [showTerminal, setShowTerminal] = useState(false);
  const [showLogs, setShowLogs] = useState(false);

  // Fetch namespaces when dialog opens
  useEffect(() => {
    if (open && contextId) {
      fetchNamespaces();
    }
  }, [open, contextId]);

  // Fetch pods when namespace is selected
  useEffect(() => {
    if (selectedNamespace) {
      fetchPods(selectedNamespace.name);
    }
  }, [selectedNamespace]);

  // Extract containers when pod is selected
  useEffect(() => {
    if (selectedPod) {
      extractContainers(selectedPod);
    }
  }, [selectedPod]);

  const fetchNamespaces = async () => {
    setLoading(true);
    setError(null);
    try {
      // TODO: Replace with actual GraphQL query to fetch namespaces
      // This is a placeholder - you'll need to implement the actual query
      const response = await fetch(
        `/api/system/kubernetes/namespaces?context=${contextId}`,
      );
      if (!response.ok) throw new Error('Failed to fetch namespaces');
      const data = await response.json();
      setNamespaces(data.namespaces || []);
    } catch (err) {
      setError(err.message);
      console.error('Error fetching namespaces:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchPods = async (namespace) => {
    setLoading(true);
    setError(null);
    try {
      // TODO: Replace with actual GraphQL query to fetch pods
      const response = await fetch(
        `/api/system/kubernetes/pods?context=${contextId}&namespace=${namespace}`,
      );
      if (!response.ok) throw new Error('Failed to fetch pods');
      const data = await response.json();
      setPods(data.pods || []);
    } catch (err) {
      setError(err.message);
      console.error('Error fetching pods:', err);
    } finally {
      setLoading(false);
    }
  };

  const extractContainers = (pod) => {
    // Extract containers from pod spec
    const containerList = [];
    if (pod.spec && pod.spec.containers) {
      containerList.push(...pod.spec.containers.map((c) => ({ name: c.name, type: 'container' })));
    }
    if (pod.spec && pod.spec.initContainers) {
      containerList.push(
        ...pod.spec.initContainers.map((c) => ({ name: c.name, type: 'initContainer' })),
      );
    }
    setContainers(containerList);
    // Auto-select first container if only one
    if (containerList.length === 1) {
      setSelectedContainer(containerList[0]);
    }
  };

  const handleOpenTerminal = () => {
    if (!selectedNamespace || !selectedPod) {
      setError('Please select namespace and pod');
      return;
    }
    setShowTerminal(true);
  };

  const handleOpenLogs = () => {
    if (!selectedNamespace || !selectedPod) {
      setError('Please select namespace and pod');
      return;
    }
    setShowLogs(true);
  };

  const handleClose = () => {
    setShowTerminal(false);
    setShowLogs(false);
    setSelectedNamespace(null);
    setSelectedPod(null);
    setSelectedContainer(null);
    setContainers([]);
    setPods([]);
    setError(null);
    onClose();
  };

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  if (showTerminal) {
    return (
      <Dialog open={true} onClose={() => setShowTerminal(false)} maxWidth="xl" fullWidth>
        <K8sTerminal
          namespace={selectedNamespace.name}
          podName={selectedPod.metadata.name}
          containerName={selectedContainer?.name || ''}
          contextId={contextId}
          onClose={() => setShowTerminal(false)}
        />
      </Dialog>
    );
  }

  if (showLogs) {
    return (
      <Dialog open={true} onClose={() => setShowLogs(false)} maxWidth="xl" fullWidth>
        <K8sLogsViewer
          namespace={selectedNamespace.name}
          podName={selectedPod.metadata.name}
          containerName={selectedContainer?.name || ''}
          contextId={contextId}
          onClose={() => setShowLogs(false)}
        />
      </Dialog>
    );
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <TerminalIcon />
          <Typography variant="h6">Pod Terminal & Logs</Typography>
        </Box>
      </DialogTitle>

      <DialogContent>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
          {error && (
            <Alert severity="error" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          <Autocomplete
            options={namespaces}
            getOptionLabel={(option) => option.name || ''}
            value={selectedNamespace}
            onChange={(event, newValue) => {
              setSelectedNamespace(newValue);
              setSelectedPod(null);
              setSelectedContainer(null);
              setPods([]);
              setContainers([]);
            }}
            renderInput={(params) => (
              <TextField {...params} label="Namespace" placeholder="Select namespace" required />
            )}
            loading={loading && namespaces.length === 0}
          />

          <Autocomplete
            options={pods}
            getOptionLabel={(option) => option.metadata?.name || ''}
            value={selectedPod}
            onChange={(event, newValue) => {
              setSelectedPod(newValue);
              setSelectedContainer(null);
            }}
            disabled={!selectedNamespace}
            renderInput={(params) => (
              <TextField {...params} label="Pod" placeholder="Select pod" required />
            )}
            loading={loading && pods.length === 0}
            renderOption={(props, option) => (
              <li {...props}>
                <Box>
                  <Typography variant="body2">{option.metadata?.name}</Typography>
                  <Typography variant="caption" color="text.secondary">
                    Status: {option.status?.phase || 'Unknown'}
                  </Typography>
                </Box>
              </li>
            )}
          />

          {containers.length > 1 && (
            <Autocomplete
              options={containers}
              getOptionLabel={(option) => `${option.name} (${option.type})`}
              value={selectedContainer}
              onChange={(event, newValue) => setSelectedContainer(newValue)}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Container (Optional)"
                  placeholder="Select container"
                />
              )}
            />
          )}

          {containers.length === 1 && (
            <Alert severity="info">
              Container: {containers[0].name} (auto-selected)
            </Alert>
          )}

          <Tabs value={tabValue} onChange={handleTabChange} centered>
            <Tab icon={<TerminalIcon />} label="Terminal" />
            <Tab icon={<DescriptionIcon />} label="Logs" />
          </Tabs>
        </Box>
      </DialogContent>

      <DialogActions>
        <Button onClick={handleClose}>Cancel</Button>
        <Button
          onClick={tabValue === 0 ? handleOpenTerminal : handleOpenLogs}
          variant="contained"
          disabled={!selectedNamespace || !selectedPod || loading}
          startIcon={tabValue === 0 ? <TerminalIcon /> : <DescriptionIcon />}
        >
          {tabValue === 0 ? 'Open Terminal' : 'View Logs'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default PodSelector;

