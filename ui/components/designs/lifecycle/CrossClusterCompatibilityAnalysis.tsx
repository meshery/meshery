import React, { useMemo, useState } from 'react';
import {
  Box,
  Button,
  Card,
  Chip,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@sistent/sistent';
import KubernetesIcon from '@/assets/icons/KubernetesIcon';
import CheckIcon from '@/assets/icons/CheckIcon';
import AlertIcon from '@/assets/icons/AlertIcon';
import CrossCircleIcon from '@/assets/icons/CrossCircleIcon';
import ClearIcon from '@/assets/icons/ClearIcon';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import { alpha, useTheme } from '@/theme';
import {
  CompatibilityAnalysisData,
  CompatibilityStatus,
  ResourceRow,
  buildCompatibilityAnalysis,
} from './compatibilityAnalysisUtils';
import {
  BottomGrid,
  CodeBox,
  DetailCard,
  DetailSection,
  HeaderContainer,
  InsightBox,
  MatrixCard,
  MetricCard,
  MetricsGrid,
  RootContainer,
  StatusChipStyled,
  TwoColumnGrid,
} from './CrossClusterCompatibilityAnalysis.styles';

export interface CrossClusterCompatibilityAnalysisProps {
  dryRunResponse?: any;
  design?: any;
  selectedK8sConnections?: any[];
  selectedK8sContexts?: string[];
  onReRunAnalysis?: () => void;
  isReRunning?: boolean;
}

/**
 * CrossClusterCompatibilityAnalysis renders a consolidated multi-cluster matrix
 * showing resource-level dry run results across environments.
 */
export const CrossClusterCompatibilityAnalysis: React.FC<
  CrossClusterCompatibilityAnalysisProps
> = ({
  dryRunResponse,
  design,
  selectedK8sConnections = [],
  selectedK8sContexts = [],
  onReRunAnalysis,
  isReRunning = false,
}) => {
  const theme = useTheme();

  const data: CompatibilityAnalysisData = useMemo(() => {
    return buildCompatibilityAnalysis({
      dryRunResponse,
      design,
      selectedK8sConnections,
      selectedK8sContexts,
    });
  }, [dryRunResponse, design, selectedK8sConnections, selectedK8sContexts]);

  // Selected cell state
  const [selectedCell, setSelectedCell] = useState<{
    resource: ResourceRow;
    clusterId: string;
  } | null>(() => {
    // Default select first failed or warning resource if available
    const failedRes = data.resources.find((r) =>
      Object.values(r.clusterResults).some((res) => res.status === 'failed'),
    );
    if (failedRes) {
      const clusterId =
        Object.keys(failedRes.clusterResults).find(
          (cid) => failedRes.clusterResults[cid].status === 'failed',
        ) || data.clusters[0]?.id;
      return { resource: failedRes, clusterId };
    }
    return null;
  });

  const selectedCluster = useMemo(() => {
    if (!selectedCell) return null;
    return data.clusters.find((c) => c.id === selectedCell.clusterId) || null;
  }, [selectedCell, data.clusters]);

  const selectedResult = useMemo(() => {
    if (!selectedCell || !selectedCluster) return null;
    return selectedCell.resource.clusterResults[selectedCluster.id] || null;
  }, [selectedCell, selectedCluster]);

  const renderStatusBadge = (
    status: CompatibilityStatus,
    isSelected: boolean,
    onClick: () => void,
  ) => {
    let icon = <CheckIcon width={14} height={14} fill={theme.palette.success.main} />;
    let label = 'Success';

    if (status === 'failed') {
      icon = <CrossCircleIcon width={14} height={14} fill={theme.palette.error.main} />;
      label = 'Failed';
    } else if (status === 'warning') {
      icon = <AlertIcon width={14} height={14} fill={theme.palette.warning.main} />;
      label = 'Warning';
    } else if (status === 'incomplete') {
      label = 'Incomplete';
    }

    return (
      <StatusChipStyled
        size="small"
        status={status}
        isSelected={isSelected}
        icon={icon}
        label={label}
        onClick={onClick}
        data-testid={`status-badge-${status}`}
      />
    );
  };

  return (
    <RootContainer data-testid="cross-cluster-compatibility-analysis">
      {/* 1. Header */}
      <HeaderContainer>
        <Box>
          <Typography
            variant="h5"
            sx={{ fontWeight: 600, display: 'flex', alignItems: 'center', gap: 1 }}
          >
            <DryRunIcon width={24} height={24} />
            Cross-Cluster Compatibility Analysis
            <Chip
              label="Beta"
              size="small"
              sx={{
                fontSize: '0.75rem',
                height: '20px',
                backgroundColor: alpha(theme.palette.primary.main, 0.2),
                color: theme.palette.primary.main,
                fontWeight: 600,
              }}
            />
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Compare the dry-run validation results of this design across multiple Kubernetes
            clusters.
          </Typography>
        </Box>
        {onReRunAnalysis && (
          <Button
            variant="contained"
            color="primary"
            startIcon={<DryRunIcon width={18} height={18} />}
            onClick={onReRunAnalysis}
            disabled={isReRunning}
            data-testid="rerun-analysis-btn"
          >
            {isReRunning ? 'Analyzing...' : 'Re-run Analysis'}
          </Button>
        )}
      </HeaderContainer>

      {/* 2. Metric Cards Grid */}
      <MetricsGrid>
        <MetricCard>
          <KubernetesIcon width={32} height={32} />
          <Box>
            <Typography variant="h4" sx={{ fontWeight: 700 }}>
              {data.metrics.totalClusters}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              Clusters Analyzed
            </Typography>
          </Box>
        </MetricCard>

        <MetricCard>
          <Box
            sx={{
              width: 32,
              height: 32,
              borderRadius: '50%',
              backgroundColor: alpha(theme.palette.primary.main, 0.15),
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontWeight: 700,
              color: theme.palette.primary.main,
            }}
          >
            #
          </Box>
          <Box>
            <Typography variant="h4" sx={{ fontWeight: 700 }}>
              {data.metrics.totalResources}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              Resources Checked
            </Typography>
          </Box>
        </MetricCard>

        <MetricCard accentColor={theme.palette.success.main}>
          <Box
            sx={{
              width: 32,
              height: 32,
              borderRadius: '50%',
              backgroundColor: alpha(theme.palette.success.main, 0.15),
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <CheckIcon width={18} height={18} fill={theme.palette.success.main} />
          </Box>
          <Box>
            <Typography variant="h4" sx={{ fontWeight: 700, color: theme.palette.success.main }}>
              {data.metrics.succeededCount}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              Succeeded
            </Typography>
          </Box>
        </MetricCard>

        <MetricCard accentColor={theme.palette.warning.main}>
          <Box
            sx={{
              width: 32,
              height: 32,
              borderRadius: '50%',
              backgroundColor: alpha(theme.palette.warning.main, 0.15),
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <AlertIcon width={18} height={18} fill={theme.palette.warning.main} />
          </Box>
          <Box>
            <Typography variant="h4" sx={{ fontWeight: 700, color: theme.palette.warning.main }}>
              {data.metrics.warningCount}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              Warnings
            </Typography>
          </Box>
        </MetricCard>

        <MetricCard accentColor={theme.palette.error.main}>
          <Box
            sx={{
              width: 32,
              height: 32,
              borderRadius: '50%',
              backgroundColor: alpha(theme.palette.error.main, 0.15),
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <CrossCircleIcon width={18} height={18} fill={theme.palette.error.main} />
          </Box>
          <Box>
            <Typography variant="h4" sx={{ fontWeight: 700, color: theme.palette.error.main }}>
              {data.metrics.failedCount}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              Failed
            </Typography>
          </Box>
        </MetricCard>
      </MetricsGrid>

      {/* 3. Main Analysis Grid: Matrix + Details Inspector */}
      <TwoColumnGrid>
        {/* Left: Compatibility Matrix Table */}
        <MatrixCard>
          <Box sx={{ p: 2, borderBottom: `1px solid ${theme.palette.divider}` }}>
            <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
              Compatibility Matrix
            </Typography>
            <Typography variant="caption" color="text.secondary">
              Click any resource cell to inspect error details and dependency impact.
            </Typography>
          </Box>
          <TableContainer sx={{ maxHeight: 440 }}>
            <Table stickyHeader size="small">
              <TableHead>
                <TableRow>
                  <TableCell sx={{ fontWeight: 600, width: '25%' }}>Resource</TableCell>
                  <TableCell sx={{ fontWeight: 600, width: '20%' }}>Kind</TableCell>
                  {data.clusters.map((cluster) => (
                    <TableCell key={cluster.id} align="center" sx={{ fontWeight: 600 }}>
                      <Box
                        sx={{
                          display: 'flex',
                          flexDirection: 'column',
                          alignItems: 'center',
                          gap: 0.25,
                        }}
                      >
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                          <KubernetesIcon width={16} height={16} />
                          <Typography variant="body2" sx={{ fontWeight: 600 }}>
                            {cluster.name}
                          </Typography>
                        </Box>
                        <Chip
                          label={cluster.version}
                          size="small"
                          sx={{ fontSize: '0.65rem', height: 18 }}
                        />
                      </Box>
                    </TableCell>
                  ))}
                </TableRow>
              </TableHead>
              <TableBody>
                {data.resources.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={2 + data.clusters.length} align="center" sx={{ py: 4 }}>
                      <Typography variant="body2" color="text.secondary">
                        No resources found in current design validation.
                      </Typography>
                    </TableCell>
                  </TableRow>
                ) : (
                  data.resources.map((res) => (
                    <TableRow
                      key={res.resourceName}
                      hover
                      selected={selectedCell?.resource.resourceName === res.resourceName}
                      sx={{ cursor: 'pointer' }}
                    >
                      <TableCell sx={{ fontWeight: 500 }}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Typography variant="body2">{res.resourceName}</Typography>
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={res.kind}
                          size="small"
                          variant="outlined"
                          sx={{ fontSize: '0.75rem', height: 22 }}
                        />
                      </TableCell>
                      {data.clusters.map((cluster) => {
                        const cellRes = res.clusterResults[cluster.id] || { status: 'incomplete' };
                        const isSelected =
                          selectedCell?.resource.resourceName === res.resourceName &&
                          selectedCell?.clusterId === cluster.id;
                        return (
                          <TableCell key={cluster.id} align="center">
                            {renderStatusBadge(cellRes.status, isSelected, () => {
                              setSelectedCell({ resource: res, clusterId: cluster.id });
                            })}
                          </TableCell>
                        );
                      })}
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </MatrixCard>

        {/* Right: Validation Details Inspector */}
        <DetailCard data-testid="validation-details-panel">
          {selectedCell && selectedCluster && selectedResult ? (
            <>
              <Box
                sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}
              >
                <Box>
                  <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
                    {selectedCell.resource.resourceName} — Validation Details
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    Cluster: <strong>{selectedCluster.name}</strong> ({selectedCluster.version}) |
                    Kind: <strong>{selectedCell.resource.kind}</strong> | Status:{' '}
                    <Box
                      component="span"
                      sx={{
                        color:
                          selectedResult.status === 'failed'
                            ? theme.palette.error.main
                            : selectedResult.status === 'warning'
                              ? theme.palette.warning.main
                              : selectedResult.status === 'incomplete'
                                ? theme.palette.text.secondary
                                : theme.palette.success.main,
                        fontWeight: 600,
                        textTransform: 'capitalize',
                      }}
                    >
                      {selectedResult.status}
                    </Box>
                  </Typography>
                </Box>
                <IconButton size="small" onClick={() => setSelectedCell(null)}>
                  <ClearIcon width={16} height={16} />
                </IconButton>
              </Box>

              {/* Validation Error Section */}
              <DetailSection>
                <Typography
                  variant="body2"
                  sx={{ fontWeight: 600, display: 'flex', alignItems: 'center', gap: 0.5 }}
                >
                  <CrossCircleIcon width={16} height={16} fill={theme.palette.error.main} />
                  Validation Error
                </Typography>
                {selectedResult.errorMessage ? (
                  <CodeBox>{selectedResult.errorMessage}</CodeBox>
                ) : (
                  <Typography variant="body2" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                    No error detected. Component passed validation on {selectedCluster.name}.
                  </Typography>
                )}
              </DetailSection>

              {/* Possible Fix Section */}
              {selectedResult.possibleFix && (
                <DetailSection>
                  <Typography
                    variant="body2"
                    sx={{
                      fontWeight: 600,
                      display: 'flex',
                      alignItems: 'center',
                      gap: 0.5,
                      color: theme.palette.primary.main,
                    }}
                  >
                    Possible Fix
                  </Typography>
                  <InsightBox tone="fix">{selectedResult.possibleFix}</InsightBox>
                </DetailSection>
              )}

              {/* Dependency Impact Section */}
              {selectedResult.dependencyImpact && (
                <DetailSection>
                  <Typography
                    variant="body2"
                    sx={{
                      fontWeight: 600,
                      display: 'flex',
                      alignItems: 'center',
                      gap: 0.5,
                      color: theme.palette.warning.main,
                    }}
                  >
                    Dependency Impact
                  </Typography>
                  <InsightBox tone="dependency">{selectedResult.dependencyImpact}</InsightBox>
                </DetailSection>
              )}
            </>
          ) : (
            <Box sx={{ textAlign: 'center', py: 6 }}>
              <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 1 }}>
                Select a Resource to Inspect Details
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Click any cell in the Compatibility Matrix table to review validation errors,
                suggested fixes, and dependency impacts.
              </Typography>
            </Box>
          )}
        </DetailCard>
      </TwoColumnGrid>

      {/* 4. Bottom Row: Validation Summary + Next Steps */}
      <BottomGrid>
        {/* Left: Validation Summary */}
        <Card
          sx={{
            p: 2.5,
            backgroundColor: theme.palette.background.paper,
            border: `1px solid ${theme.palette.divider}`,
          }}
        >
          <Typography
            variant="subtitle1"
            sx={{ fontWeight: 600, mb: 1.5, display: 'flex', alignItems: 'center', gap: 1 }}
          >
            <KubernetesIcon width={18} height={18} />
            Validation Summary
          </Typography>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
            {data.clusterSummaries.map((summary) => {
              const isAllPass = summary.failed === 0 && summary.warnings === 0;
              const hasErrors = summary.failed > 0;
              return (
                <Box
                  key={summary.clusterId}
                  sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    p: 1.25,
                    borderRadius: 1,
                    backgroundColor: alpha(theme.palette.background.default, 0.6),
                  }}
                >
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Box
                      sx={{
                        width: 10,
                        height: 10,
                        borderRadius: '50%',
                        backgroundColor: hasErrors
                          ? theme.palette.error.main
                          : summary.warnings > 0
                            ? theme.palette.warning.main
                            : theme.palette.success.main,
                      }}
                    />
                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                      {summary.clusterName}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      ({summary.version})
                    </Typography>
                  </Box>

                  <Typography
                    variant="body2"
                    sx={{
                      fontWeight: 600,
                      color: hasErrors
                        ? theme.palette.error.main
                        : summary.warnings > 0
                          ? theme.palette.warning.main
                          : theme.palette.success.main,
                    }}
                  >
                    {isAllPass
                      ? `${summary.passed} / ${summary.total} Passed`
                      : `${summary.passed} Passed${summary.warnings > 0 ? `, ${summary.warnings} Warning` : ''}${
                          summary.failed > 0 ? `, ${summary.failed} Failed` : ''
                        }`}
                  </Typography>
                </Box>
              );
            })}
          </Box>
        </Card>

        {/* Right: Next Steps */}
        <Card
          sx={{
            p: 2.5,
            backgroundColor: theme.palette.background.paper,
            border: `1px solid ${theme.palette.divider}`,
          }}
        >
          <Typography
            variant="subtitle1"
            sx={{ fontWeight: 600, mb: 1.5, display: 'flex', alignItems: 'center', gap: 1 }}
          >
            <DryRunIcon width={18} height={18} />
            Next Steps
          </Typography>
          <Box
            component="ul"
            sx={{ pl: 2.5, m: 0, display: 'flex', flexDirection: 'column', gap: 1 }}
          >
            {data.nextSteps.map((step, idx) => (
              <Typography key={idx} component="li" variant="body2" color="text.secondary">
                {step}
              </Typography>
            ))}
          </Box>
          {onReRunAnalysis && (
            <Box sx={{ mt: 2.5, display: 'flex', justifyContent: 'flex-end' }}>
              <Button
                variant="outlined"
                size="small"
                startIcon={<DryRunIcon width={16} height={16} />}
                onClick={onReRunAnalysis}
                disabled={isReRunning}
              >
                Re-run Analysis
              </Button>
            </Box>
          )}
        </Card>
      </BottomGrid>
    </RootContainer>
  );
};

export default CrossClusterCompatibilityAnalysis;
