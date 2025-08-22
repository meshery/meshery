import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Grid2,
  Alert,
  CircularProgress,
  List,
  ListItem,
  ListItemText,
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

// Real API call for Kubernetes discovery using actual Meshery APIs
const performKubernetesDiscovery = async (uploadedFormData, addK8sConfig) => {
  if (!uploadedFormData) {
    throw new Error('No kubeconfig file uploaded');
  }

  try {
    const result = await addK8sConfig({ body: uploadedFormData }).unwrap();
    console.log('result', result);
    const discoveredContexts = [];

    // Process registered contexts
    if (result.registered_contexts) {
      result.registered_contexts.forEach((context) => {
        discoveredContexts.push({
          id: context.connection_id,
          name: context.name,
          server: context.server,
          namespace: context.namespace || 'default',
          status: {
            connectivity: 'success',
            authentication: 'success',
            authorization: 'success',
            canInstallCRD: true,
          },
          details: {
            version: context.version || 'Unknown',
          },
          warnings: [],
          errors: [],
          connectionState: 'registered',
        });
      });
    }

    // Process connected contexts
    if (result.connected_contexts) {
      result.connected_contexts.forEach((context) => {
        discoveredContexts.push({
          id: context.connection_id,
          name: context.name,
          server: context.server,
          namespace: context.namespace || 'default',
          status: {
            connectivity: 'success', // Match SelectionStep expectations
            authentication: 'success', // Match SelectionStep expectations
            authorization: 'success', // Match SelectionStep expectations
            canInstallCRD: true,
          },
          details: {
            version: context.version || 'Unknown',
            nodes: context.nodes || 0,
          },
          warnings: [],
          errors: [],
          connectionState: 'connected',
        });
      });
    }

    // Process ignored contexts
    if (result.ignored_contexts) {
      result.ignored_contexts.forEach((context) => {
        discoveredContexts.push({
          id: context.connection_id,
          name: context.name,
          server: context.server,
          namespace: context.namespace || 'default',
          status: {
            connectivity: 'failed', // -> add reason from backend (nikita)
            authentication: 'failed', // -> add reason from backend (nikita)
            authorization: 'failed', // -> add reason from backend (nikita)
            canInstallCRD: false,
          },
          details: {},
          warnings: [''],
          errors: [],
          connectionState: 'ignored',
        });
      });
    }

    // Process errored contexts
    if (result.errored_contexts) {
      result.errored_contexts.forEach((context) => {
        discoveredContexts.push({
          id: context.connection_id,
          name: context.name,
          server: context.server || 'Unknown Server',
          namespace: context.namespace || 'default',
          status: {
            connectivity: 'failed', // -> add reason from backend (nikita)
            authentication: 'failed', // -> add reason from backend (nikita)
            authorization: 'failed', // -> add reason from backend (nikita)
            canInstallCRD: false,
          },
          details: {},
          warnings: [],
          errors: [context.error || 'Connection failed'],
          connectionState: 'errored',
        });
      });
    }

    return {
      success: true,
      contexts: discoveredContexts,
      connected:
        (result.registered_contexts?.length || 0) + (result.connected_contexts?.length || 0),
      failed: result.errored_contexts?.length || 0,
    };
  } catch (error) {
    throw new Error(`Discovery failed: ${error.message}`);
  }
};

const getStatusColor = (status) => {
  switch (status) {
    case 'success':
      return 'success';
    case 'warning':
      return 'warning';
    case 'failed':
      return 'error';
    default:
      return 'default';
  }
};

const getStatusIcon = (status) => {
  switch (status) {
    case 'success':
      return <CheckCircleIcon />;
    case 'warning':
      return <WarningIcon />;
    case 'failed':
      return <ErrorIcon />;
    default:
      return null;
  }
};

const ContextCard = ({ context, expanded, onToggle }) => {
  const { status, details, warnings = [], errors = [] } = context;
  const hasIssues = warnings.length > 0 || errors.length > 0;

  const overallStatus =
    status.connectivity === 'failed'
      ? 'failed'
      : status.authorization === 'warning' || warnings.length > 0
        ? 'warning'
        : 'success';

  return (
    <Card
      sx={{
        mb: 2,
        border:
          overallStatus === 'failed'
            ? '1px solid #f44336'
            : overallStatus === 'warning'
              ? '1px solid #ff9800'
              : '1px solid #4caf50',
      }}
    >
      <CardContent sx={{ pb: 1 }}>
        <Box display="flex" justifyContent="space-between" alignItems="flex-start">
          <Box flex={1}>
            <Typography variant="h6" component="div" gutterBottom>
              {context.name}
            </Typography>
            <Typography variant="body2" color="textSecondary" gutterBottom>
              {context.server}
            </Typography>

            <Box display="flex" gap={1} flexWrap="wrap" mb={1}>
              <Box display="flex" alignItems="center" gap={0.5}>
                {getStatusIcon(status.connectivity, 'small')}
                <Typography variant="body2" color={getStatusColor(status.connectivity)}>
                  {`Connectivity: ${status.connectivity}`}
                </Typography>
              </Box>
              <Box display="flex" alignItems="center" gap={0.5}>
                {getStatusIcon(status.authentication, 'small')}
                <Typography variant="body2" color={getStatusColor(status.authentication)}>
                  {`Auth: ${status.authentication}`}
                </Typography>
              </Box>
            </Box>

            {details && (
              <Typography variant="body2" color="textSecondary">
                Version: {details.version}
              </Typography>
            )}
          </Box>

          {hasIssues && (
            <Box
              onClick={onToggle}
              sx={{ cursor: 'pointer', display: 'flex', alignItems: 'center' }}
            >
              {expanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
            </Box>
          )}
        </Box>

        <Collapse in={expanded}>
          {warnings.length > 0 && (
            <Alert severity="warning" sx={{ mt: 2, mb: 1 }}>
              <Typography variant="subtitle2">Warnings:</Typography>
              <List dense>
                {warnings.map((warning, index) => (
                  <ListItem key={index} sx={{ py: 0 }}>
                    <ListItemText primary={warning} />
                  </ListItem>
                ))}
              </List>
            </Alert>
          )}

          {errors.length > 0 && (
            <Alert severity="error" sx={{ mt: 2 }}>
              <Typography variant="subtitle2">Errors:</Typography>
              <List dense>
                {errors.map((error, index) => (
                  <ListItem key={index} sx={{ py: 0 }}>
                    <ListItemText primary={error} />
                  </ListItem>
                ))}
              </List>
            </Alert>
          )}
        </Collapse>
      </CardContent>
    </Card>
  );
};

export const DiscoveryStep = ({
  connectionType = 'kubernetes',
  wizardData,
  setWizardData,
  onNext,
  onBack,
}) => {
  const [isDiscovering, setIsDiscovering] = useState(false);
  const [discoveryComplete, setDiscoveryComplete] = useState(false);
  const [discoveredConnections, setDiscoveredConnections] = useState([]);
  const [expandedCards, setExpandedCards] = useState({});
  const { notify } = useNotification();
  const [addK8sConfig] = useAddKubernetesConfigMutation();
  console.log('wizardData', wizardData);
  useEffect(() => {
    if (!wizardData?.discoveredConnections?.length) {
      performDiscovery();
    } else if (wizardData?.discoveredConnections?.length) {
      setDiscoveredConnections(wizardData.discoveredConnections);
      setDiscoveryComplete(true);
    }
  }, [wizardData?.uploadedFormData]);

  const performDiscovery = async () => {
    setIsDiscovering(true);

    try {
      let results;

      if (connectionType === 'kubernetes') {
        results = await performKubernetesDiscovery(wizardData.uploadedFormData, addK8sConfig);
      }
      setDiscoveredConnections(results.contexts);
      console.log('discoveredConnections', discoveredConnections);
      setWizardData((prev) => ({
        ...prev,
        discoveredConnections: results.contexts,
        discoveryResults: results,
      }));

      setDiscoveryComplete(true);

      notify({
        message: `Discovered ${results.contexts.length} contexts`,
        event_type: EVENT_TYPES.INFO,
      });
    } catch (error) {
      notify({
        message: `Discovery failed: ${error.message}`,
        event_type: EVENT_TYPES.ERROR,
      });
    } finally {
      setIsDiscovering(false);
    }
  };

  const handleToggleCard = (contextId) => {
    setExpandedCards((prev) => ({
      ...prev,
      [contextId]: !prev[contextId],
    }));
  };

  const getStatusSummary = () => {
    const successful = discoveredConnections.filter(
      (ctx) =>
        ctx.status.connectivity === 'success' &&
        ctx.status.authentication === 'success' &&
        ctx.status.authorization === 'success',
    ).length;

    const withWarnings = discoveredConnections.filter(
      (ctx) =>
        ctx.status.connectivity === 'success' &&
        ctx.status.authentication === 'success' &&
        (ctx.status.authorization === 'warning' || ctx.warnings?.length > 0),
    ).length;

    const failed = discoveredConnections.filter(
      (ctx) => ctx.status.connectivity === 'failed' || ctx.status.authentication === 'failed',
    ).length;

    return { successful, withWarnings, failed, total: discoveredConnections.length };
  };

  const summary = discoveryComplete ? getStatusSummary() : null;

  if (isDiscovering) {
    return (
      <StepWrapper
        title="Discovering Connections..."
        description="Analyzing your configuration and checking connectivity"
        showNextButton={false}
        showBackButton={false}
      >
        <Box display="flex" flexDirection="column" alignItems="center" mt={4}>
          <CircularProgress size={60} />
          <Typography variant="h6" sx={{ mt: 2 }}>
            Performing pre-flight checks...
          </Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mt: 1 }}>
            This may take a few moments
          </Typography>
        </Box>
      </StepWrapper>
    );
  }

  return (
    <StepWrapper
      title="Discovery Results"
      description="Review the discovered connections and their status"
      onNext={onNext}
      onBack={onBack}
      nextDisabled={!discoveryComplete}
      nextButtonText="Continue to Selection"
    >
      <Box>
        {summary && (
          <Alert severity="info" sx={{ mb: 3 }}>
            <Typography variant="subtitle1" gutterBottom>
              Discovery Summary
            </Typography>
            <Typography variant="body2">
              Found <strong>{summary.total}</strong> contexts:{' '}
              <strong style={{ color: '#4caf50' }}>{summary.successful}</strong> ready for
              management, <strong style={{ color: '#ff9800' }}>{summary.withWarnings}</strong> with
              limitations, <strong style={{ color: '#f44336' }}>{summary.failed}</strong>{' '}
              unreachable
            </Typography>
          </Alert>
        )}

        <Grid2 container spacing={2}>
          {discoveredConnections.map((context) => (
            <Grid2 key={context.id} size={{ xs: 12, md: 6 }}>
              <ContextCard
                context={context}
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
