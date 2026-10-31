import React from 'react';
import { Box, IconButton, Typography } from '@sistent/sistent';
import ClearIcon from '@/assets/icons/ClearIcon';
import CrossCircleIcon from '@/assets/icons/CrossCircleIcon';
import { useTheme } from '@/theme';
import { ClusterInfo, ResourceRow, ResourceValidationDetail } from './compatibilityAnalysisUtils';
import {
  CodeBox,
  DetailCard,
  DetailSection,
  InsightBox,
} from './CrossClusterCompatibilityAnalysis.styles';

export interface ValidationDetailsPanelProps {
  selectedCell: { resource: ResourceRow; clusterId: string } | null;
  selectedCluster: ClusterInfo | null;
  selectedResult: ResourceValidationDetail | null;
  onClearSelection: () => void;
}

export const ValidationDetailsPanel: React.FC<ValidationDetailsPanelProps> = ({
  selectedCell,
  selectedCluster,
  selectedResult,
  onClearSelection,
}) => {
  const theme = useTheme();

  return (
    <DetailCard data-testid="validation-details-panel">
      {selectedCell && selectedCluster && selectedResult ? (
        <>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
            <Box>
              <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
                {selectedCell.resource.resourceName} — Validation Details
              </Typography>
              <Typography variant="caption" color="text.secondary">
                Cluster: <strong>{selectedCluster.name}</strong> ({selectedCluster.version}) | Kind:{' '}
                <strong>{selectedCell.resource.kind}</strong> | Status:{' '}
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
            <IconButton size="small" onClick={onClearSelection}>
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
            Click any cell in the Compatibility Matrix table to review validation errors, suggested
            fixes, and dependency impacts.
          </Typography>
        </Box>
      )}
    </DetailCard>
  );
};

export default ValidationDetailsPanel;
