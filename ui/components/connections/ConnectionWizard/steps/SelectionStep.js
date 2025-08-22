import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  FormControlLabel,
  RadioGroup,
  Radio,
  Grid2,
  Chip,
  Alert,
  Divider,
  FormControl,
  FormLabel,
  Button,
} from '@sistent/sistent';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import WarningIcon from '@mui/icons-material/Warning';
import ErrorIcon from '@mui/icons-material/Error';
import { StepWrapper } from '../common/StepWrapper';

const DEPLOYMENT_MODES = {
  operator: {
    label: 'Operator Mode',
    description: 'Deploy MeshSync as a separate operator in the cluster',
    requirements: ['Can install CRDs', 'Cluster admin permissions'],
  },
  embedded: {
    label: 'Embedded Mode',
    description: 'Run MeshSync embedded within Meshery Server',
    requirements: ['Network connectivity', 'Basic read permissions'],
  },
};

// const CONNECTION_ACTIONS = {
//   manage: 'Manage with Meshery',
//   ignore: 'Do not manage',
// };

const SelectionCard = ({ context, selection, onSelectionChange, globalDeploymentMode }) => {
  const { status, warnings = [], errors = [] } = context;

  // Debug logging to understand the status values
  console.log('SelectionCard context:', context.name, 'status:', status);

  const isReachable = status.connectivity === 'success' && status.authentication === 'success';
  const canManage =
    isReachable && (status.authorization === 'success' || status.authorization === 'warning');
  const canUseOperator = canManage && status.canInstallCRD;

  console.log('SelectionCard computed:', {
    name: context.name,
    isReachable,
    canManage,
    canUseOperator,
  });

  const getRecommendedMode = () => {
    if (!canManage) return null;
    return canUseOperator ? 'operator' : 'embedded';
  };

  const recommendedMode = getRecommendedMode();

  const handleActionChange = (action) => {
    onSelectionChange(context.id, {
      action,
      deploymentMode: action === 'manage' ? recommendedMode || globalDeploymentMode : null,
    });
  };

  const handleModeChange = (mode) => {
    onSelectionChange(context.id, {
      ...selection,
      deploymentMode: mode,
    });
  };

  return (
    <Card
      sx={{
        mb: 2,
        opacity: canManage ? 1 : 0.7,
        border: !canManage
          ? '1px solid #ccc'
          : selection?.action === 'manage'
            ? '2px solid #00B39F'
            : '1px solid #ddd',
      }}
    >
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="flex-start" mb={2}>
          <Box flex={1}>
            <Typography variant="h6" component="div">
              {context.name}
            </Typography>
            <Typography variant="body2" color="textSecondary" gutterBottom>
              {context.server}
            </Typography>

            <Box display="flex" gap={1} flexWrap="wrap" mb={2}>
              {isReachable ? (
                <Box display="flex" alignItems="center" gap={0.5}>
                  <CheckCircleIcon color="success" fontSize="small" />
                  <Typography variant="body2">Reachable</Typography>
                </Box>
              ) : (
                <Box display="flex" alignItems="center" gap={0.5}>
                  <ErrorIcon fontSize="small" />
                  <Typography variant="body2">Unreachable</Typography>
                </Box>
              )}

              {canUseOperator && (
                <Box display="flex" alignItems="center" gap={0.5}>
                  <CheckCircleIcon color="success" fontSize="small" />
                  <Typography variant="body2">Operator Ready</Typography>
                </Box>
              )}

              {warnings.length > 0 && (
                <Chip
                  icon={<WarningIcon />}
                  label={`${warnings.length} Warning${warnings.length > 1 ? 's' : ''}`}
                  color="warning"
                  size="small"
                />
              )}

              {errors.length > 0 && (
                <Chip
                  icon={<ErrorIcon />}
                  label={`${errors.length} Error${errors.length > 1 ? 's' : ''}`}
                  color="error"
                  size="small"
                />
              )}
            </Box>
          </Box>
        </Box>

        <Divider sx={{ my: 2 }} />

        <FormControl component="fieldset" fullWidth>
          <FormLabel component="legend">Action</FormLabel>
          <RadioGroup
            value={selection?.action || 'ignore'}
            onChange={(e) => handleActionChange(e.target.value)}
          >
            <FormControlLabel
              value="manage"
              control={<Radio />}
              label="Connect"
              disabled={!canManage}
            />
            <FormControlLabel value="ignore" control={<Radio />} label="Ignore" />
          </RadioGroup>
        </FormControl>

        {selection?.action === 'manage' && (
          <Box mt={2}>
            <FormControl component="fieldset" fullWidth>
              <FormLabel component="legend">Deployment Mode</FormLabel>
              <RadioGroup
                value={selection.deploymentMode || recommendedMode || 'embedded'}
                onChange={(e) => handleModeChange(e.target.value)}
              >
                <FormControlLabel
                  value="operator"
                  control={<Radio />}
                  label={
                    <Box>
                      <Typography variant="body2">
                        <strong>Operator Mode</strong>
                      </Typography>
                      <Typography variant="caption" color="textSecondary">
                        {DEPLOYMENT_MODES.operator.description}
                      </Typography>
                    </Box>
                  }
                  disabled={!canUseOperator}
                />
                <FormControlLabel
                  value="embedded"
                  control={<Radio />}
                  label={
                    <Box>
                      <Typography variant="body2">
                        <strong>Embedded Mode</strong>
                        {recommendedMode === 'embedded' && (
                          <Chip label="Recommended" color="primary" size="small" sx={{ ml: 1 }} />
                        )}
                      </Typography>
                      <Typography variant="caption" color="textSecondary">
                        {DEPLOYMENT_MODES.embedded.description}
                      </Typography>
                    </Box>
                  }
                />
              </RadioGroup>
            </FormControl>

            {!canUseOperator && selection.deploymentMode === 'operator' && (
              <Alert severity="warning" sx={{ mt: 1 }}>
                Operator mode is not available for this cluster due to insufficient permissions.
                Embedded mode will be used instead.
              </Alert>
            )}
          </Box>
        )}

        {!canManage && (
          <Alert severity="error" sx={{ mt: 2 }}>
            This cluster cannot be managed due to connectivity or authentication issues.
          </Alert>
        )}
      </CardContent>
    </Card>
  );
};

export const SelectionStep = ({
  //   connectionType = 'kubernetes',
  wizardData,
  setWizardData,
  onNext,
  onBack,
}) => {
  const [selections, setSelections] = useState(wizardData?.selectedConnections || {});
  const [globalDeploymentMode, setGlobalDeploymentMode] = useState('operator');

  const discoveredConnections = wizardData?.discoveredConnections || [];

  // Debug logging to understand what we received from DiscoveryStep
  console.log('SelectionStep discoveredConnections:', discoveredConnections);

  useEffect(() => {
    setWizardData((prev) => ({
      ...prev,
      selectedConnections: selections,
    }));
  }, [selections, setWizardData]);

  const handleSelectionChange = (contextId, newSelection) => {
    setSelections((prev) => ({
      ...prev,
      [contextId]: newSelection,
    }));
  };

  const handleApplyGlobalMode = () => {
    const updatedSelections = { ...selections };

    discoveredConnections.forEach((context) => {
      if (selections[context.id]?.action === 'manage') {
        const canUseOperator =
          context.status.canInstallCRD &&
          context.status.connectivity === 'success' &&
          context.status.authentication === 'success';

        updatedSelections[context.id] = {
          ...updatedSelections[context.id],
          deploymentMode:
            globalDeploymentMode === 'operator' && canUseOperator ? 'operator' : 'embedded',
        };
      }
    });

    setSelections(updatedSelections);
  };

  const getSelectionSummary = () => {
    const toManage = Object.values(selections).filter((s) => s.action === 'manage').length;
    const operatorMode = Object.values(selections).filter(
      (s) => s.action === 'manage' && s.deploymentMode === 'operator',
    ).length;
    const embeddedMode = Object.values(selections).filter(
      (s) => s.action === 'manage' && s.deploymentMode === 'embedded',
    ).length;
    const toIgnore = discoveredConnections.length - toManage;

    return { toManage, operatorMode, embeddedMode, toIgnore, total: discoveredConnections.length };
  };

  const summary = getSelectionSummary();
  const hasSelections = summary.toManage > 0;

  return (
    <StepWrapper
      title="Configure Connection Management"
      description="Choose which connections to manage and how to deploy them"
      onNext={onNext}
      onBack={onBack}
      nextDisabled={!hasSelections}
      nextButtonText="Proceed with Configuration"
    >
      <Box>
        {/* Global Settings */}
        <Card sx={{ mb: 3, bgcolor: 'background.paper' }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Quick Configuration
            </Typography>
            <Box display="flex" alignItems="center" gap={2} mb={2}>
              <FormControl>
                <FormLabel>Default Deployment Mode:</FormLabel>
                <RadioGroup
                  row
                  value={globalDeploymentMode}
                  onChange={(e) => setGlobalDeploymentMode(e.target.value)}
                >
                  <FormControlLabel value="operator" control={<Radio />} label="Operator" />
                  <FormControlLabel value="embedded" control={<Radio />} label="Embedded" />
                </RadioGroup>
              </FormControl>
              <Button
                variant="outlined"
                onClick={handleApplyGlobalMode}
                disabled={summary.toManage === 0}
              >
                Apply to All
              </Button>
            </Box>
            <Typography variant="body2" color="textSecondary">
              Set a default deployment mode and apply it to all selected connections. Individual
              connections can still be customized below.
            </Typography>
          </CardContent>
        </Card>

        {/* Selection Summary */}
        {summary.total > 0 && (
          <Alert severity={hasSelections ? 'info' : 'warning'} sx={{ mb: 3 }}>
            <Typography variant="subtitle1" gutterBottom>
              Selection Summary
            </Typography>
            <Typography variant="body2">
              <strong>{summary.toManage}</strong> connections selected for management
              {summary.toManage > 0 && (
                <>
                  {' '}
                  ({summary.operatorMode} operator mode, {summary.embeddedMode} embedded mode)
                </>
              )}
              {summary.toIgnore > 0 && (
                <>
                  , <strong>{summary.toIgnore}</strong> will be ignored
                </>
              )}
            </Typography>
            {!hasSelections && (
              <Typography variant="body2" sx={{ mt: 1 }}>
                Please select at least one connection to continue.
              </Typography>
            )}
          </Alert>
        )}

        {/* Individual Connection Selection */}
        <Grid2 container spacing={2}>
          {discoveredConnections.map((context) => (
            <Grid2 key={context.id} size={{ xs: 12, lg: 6 }}>
              <SelectionCard
                context={context}
                selection={selections[context.id]}
                onSelectionChange={handleSelectionChange}
                globalDeploymentMode={globalDeploymentMode}
              />
            </Grid2>
          ))}
        </Grid2>
      </Box>
    </StepWrapper>
  );
};
