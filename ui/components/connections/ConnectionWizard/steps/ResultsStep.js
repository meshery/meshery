import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Grid2,
  Chip,
  Alert,
  List,
  ListItem,
  ListItemText,
  CircularProgress,
  Button,
  Collapse,
} from '@sistent/sistent';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';
import WarningIcon from '@mui/icons-material/Warning';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import { StepWrapper } from '../common/StepWrapper';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { useAddKubernetesConfigMutation } from '@/rtk-query/connection';

const ConnectionResult = ({ context, selection, result, expanded, onToggle }) => {
  console.log('selection', selection);
  console.log('result', result);
  console.log('expanded', expanded);

  if (!selection || selection.action !== 'manage') {
    return (
      <Card sx={{ mb: 2, bgcolor: 'grey.50' }}>
        <CardContent>
          <Box display="flex" alignItems="center" justifyContent="space-between">
            <Box>
              <Typography variant="h6">{context.name}</Typography>
              <Typography variant="body2" color="textSecondary">
                {context.server}
              </Typography>
            </Box>
            <Chip icon={<WarningIcon />} label="Skipped" color="default" size="small" />
          </Box>
        </CardContent>
      </Card>
    );
  }

  const getStatusChip = () => {
    if (!result) {
      return <Chip icon={<CircularProgress size={16} />} label="Processing..." color="primary" />;
    }
    if (result.error) {
      return <Chip icon={<ErrorIcon />} label="Failed" color="error" />;
    }
    if (result.warnings?.length > 0) {
      return <Chip icon={<WarningIcon />} label="Success with warnings" color="warning" />;
    }
    return <Chip icon={<CheckCircleIcon />} label="Success" color="success" />;
  };

  const hasDetails = result && (result.details || result.warnings || result.error);

  return (
    <Card
      sx={{
        mb: 2,
        border: result?.error
          ? '1px solid #f44336'
          : result?.warnings
            ? '1px solid #ff9800'
            : result?.success
              ? '1px solid #4caf50'
              : '1px solid #ddd',
      }}
    >
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="flex-start">
          <Box flex={1}>
            <Typography variant="h6" gutterBottom>
              {context.name}
            </Typography>
            <Typography variant="body2" color="textSecondary" gutterBottom>
              {context.server}
            </Typography>

            <Box display="flex" gap={1} alignItems="center" mb={1}>
              {getStatusChip()}
              {selection?.deploymentMode && (
                <Chip
                  label={
                    selection.deploymentMode === 'operator' ? 'Operator Mode' : 'Embedded Mode'
                  }
                  size="small"
                  variant="outlined"
                />
              )}
            </Box>

            {!result && (
              <Box display="flex" alignItems="center" mt={1}>
                <CircularProgress size={16} sx={{ mr: 1 }} />
                <Typography variant="body2">
                  {selection.deploymentMode === 'operator'
                    ? 'Installing operator...'
                    : 'Establishing connection...'}
                </Typography>
              </Box>
            )}
          </Box>

          {hasDetails && (
            <Box
              onClick={onToggle}
              sx={{ cursor: 'pointer', display: 'flex', alignItems: 'center' }}
            >
              {expanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
            </Box>
          )}
        </Box>

        <Collapse in={expanded && hasDetails}>
          {result?.success && result.details && (
            <Alert severity="info" sx={{ mt: 2, mb: 1 }}>
              <Typography variant="subtitle2">Connection Details:</Typography>
              <List dense>
                <ListItem sx={{ py: 0 }}>
                  <ListItemText primary={`Context: ${result.details.contextName}`} />
                </ListItem>
                <ListItem sx={{ py: 0 }}>
                  <ListItemText primary={`Server: ${result.details.server}`} />
                </ListItem>
                <ListItem sx={{ py: 0 }}>
                  <ListItemText primary={`Namespace: ${result.details.namespace}`} />
                </ListItem>
                <ListItem sx={{ py: 0 }}>
                  <ListItemText primary={`Connection ID: ${result.details.connectionId}`} />
                </ListItem>
                {result.details.version && (
                  <ListItem sx={{ py: 0 }}>
                    <ListItemText primary={`Version: ${result.details.version}`} />
                  </ListItem>
                )}
                <ListItem sx={{ py: 0 }}>
                  <ListItemText primary={`Deployment Mode: ${result.deploymentMode}`} />
                </ListItem>
              </List>
            </Alert>
          )}

          {result?.warnings?.length > 0 && (
            <Alert severity="warning" sx={{ mt: 2, mb: 1 }}>
              <Typography variant="subtitle2">Warnings:</Typography>
              <List dense>
                {result.warnings.map((warning, index) => (
                  <ListItem key={index} sx={{ py: 0 }}>
                    <ListItemText primary={warning} />
                  </ListItem>
                ))}
              </List>
            </Alert>
          )}

          {result?.error && (
            <Alert severity="error" sx={{ mt: 2 }}>
              <Typography variant="subtitle2">Error:</Typography>
              <Typography variant="body2">{result.error}</Typography>
            </Alert>
          )}
        </Collapse>
      </CardContent>
    </Card>
  );
};

export const ResultsStep = ({
  //   connectionType = 'kubernetes',
  wizardData,
  setWizardData,
  onComplete,
  onCancel,
}) => {
  const [executionResults, setExecutionResults] = useState({});
  // const [isExecuting, setIsExecuting] = useState(false);
  const [expandedCards, setExpandedCards] = useState({});
  // const [executionProgress, setExecutionProgress] = useState(0);
  const { notify } = useNotification();
  const [addK8sConfig] = useAddKubernetesConfigMutation();

  const discoveredConnections = wizardData?.discoveredConnections || [];

  const selectedConnections = wizardData?.selectedConnections || {};
  console.log('discoveredConnections', discoveredConnections);
  console.log('selectedConnections', selectedConnections);

  const connectionsToManage = discoveredConnections.filter(
    (ctx) => selectedConnections[ctx.id]?.action === 'manage',
  );
  console.log('connectionsToManage', connectionsToManage);

  useEffect(() => {
    if (connectionsToManage.length > 0 && Object.keys(executionResults).length === 0) {
      executeConnections();
    }
  }, []);

  const establishKubernetesConnection = async (context, deploymentMode, uploadedFile) => {
    try {
      console.log('I am called');
     
      const formData = new FormData();
      formData.append('k8sfile', uploadedFile);
      formData.append('meshsync_deployment_mode', deploymentMode);

      return await addK8sConfig({ body: formData }).unwrap();
    } catch (error) {
      throw new Error(`Failed to establish connection: ${error.message}`);
    }
  };

  const executeConnections = async () => {
    const results = {};
    // let completed = 0;

    // Execute connections sequentially with progress updates
    for (const context of connectionsToManage) {
      const selection = selectedConnections[context.id];

      try {
        const result = await establishKubernetesConnection(
          context,
          selection.deploymentMode,
          wizardData.uploadedFile,
        );
        results[context.id] = result;
        console.log('results', results);

        notify({
          message: `Successfully connected to ${context.name}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (error) {
        results[context.id] = { error: error.message };

        notify({
          message: `Failed to connect to ${context.name}: ${error.message}`,
          event_type: EVENT_TYPES.ERROR,
        });
      }

      setExecutionResults({ ...results });
    }

    setWizardData((prev) => ({
      ...prev,
      executionResults: results,
    }));
    console.log('wizardData', wizardData);
    // setIsExecuting(false);
  };

  const handleToggleCard = (contextId) => {
    setExpandedCards((prev) => ({
      ...prev,
      [contextId]: !prev[contextId],
    }));
  };

  // const getExecutionSummary = () => {
  //   const total = connectionsToManage.length;
  //   const completed = Object.keys(executionResults).length;
  //   const successful = Object.values(executionResults).filter((r) => r.success).length;
  //   const withWarnings = Object.values(executionResults).filter(
  //     (r) => r.success && r.warnings?.length > 0,
  //   ).length;
  //   const failed = Object.values(executionResults).filter((r) => r.error).length;
  //   const skipped = discoveredConnections.length - total;

  //   return { total, completed, successful, withWarnings, failed, skipped };
  // };

  // const summary = getExecutionSummary();
  // const isComplete = !isExecuting && summary.completed === summary.total;

  return (
    <StepWrapper
      title="Connection Results"
      description="Review the results of your connection setup"
      showBackButton={false}
      showNextButton={false}
      showCancelButton={false}
      customFooter={
        <Box
          display="flex"
          justifyContent="space-between"
          width="100%"
          mt={3}
          pt={2}
          sx={{ borderTop: '1px solid', borderColor: 'divider' }}
        >
          <Box>
            <Button variant="outlined" onClick={onCancel}>
              Cancel
            </Button>
          </Box>
          <Box>
            <Button variant="contained" onClick={() => onComplete(wizardData)}>
              Finish
            </Button>
          </Box>
        </Box>
      }
    >
      <Box>
        <Grid2 container spacing={2}>
          {discoveredConnections.map((context) => (
            <Grid2 key={context.id} size={{ xs: 12, lg: 6 }}>
              <ConnectionResult
                context={context}
                selection={selectedConnections[context.id]}
                result={executionResults[context.id]}
                expanded={expandedCards[context.id]}
                onToggle={() => handleToggleCard(context.id)}
              />
            </Grid2>
          ))}
        </Grid2>
      </Box>
    </StepWrapper>
  );
};
