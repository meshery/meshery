/**
 * DesignImpactPreviewModal
 *
 * Lightweight modal providing pre-change impact preview for modified or deleted
 * components and relationships within a Meshery Design before committing changes.
 */
import React, { FC } from 'react';
import {
  Box,
  Card,
  Chip,
  DialogContent,
  ModalButtonPrimary,
  ModalButtonSecondary,
  Typography,
} from '@sistent/sistent';
import { Modal } from '@/components/shared/Modal';
import AlertIcon from '@/assets/icons/AlertIcon';
import CheckIcon from '@/assets/icons/CheckIcon';
import { useTheme } from '@/theme';
import {
  DesignImpactAnalysisResult,
  ImpactedComponent,
  ImpactedRelationship,
} from './designImpactAnalysis';

export interface DesignImpactPreviewModalProps {
  isOpen: boolean;
  onClose: () => void;
  onProceed: () => void;
  impactResult: DesignImpactAnalysisResult;
  isProceedDisabled?: boolean;
}

export const DesignImpactPreviewModal: FC<DesignImpactPreviewModalProps> = ({
  isOpen,
  onClose,
  onProceed,
  impactResult,
  isProceedDisabled = false,
}) => {
  const theme = useTheme();

  const {
    targetName,
    targetKind,
    targetType,
    modificationType,
    directlyAffected,
    indirectlyAffected,
    affectedRelationships,
    totalAffectedCount,
    summary,
  } = impactResult;

  const actionVerb = modificationType === 'delete' ? 'Remove' : 'Modify';

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={`Design Impact Preview: ${actionVerb} ${targetName || 'Resource'}`}
      headerIcon={
        totalAffectedCount > 0 || affectedRelationships.length > 0 ? (
          <AlertIcon width={22} height={22} fill={theme.palette.warning.main} />
        ) : (
          <CheckIcon width={22} height={22} fill={theme.palette.success.main} />
        )
      }
      size="md"
      actions={
        <>
          <ModalButtonSecondary onClick={onClose}>Cancel</ModalButtonSecondary>
          <ModalButtonPrimary onClick={onProceed} disabled={isProceedDisabled}>
            Continue
          </ModalButtonPrimary>
        </>
      }
    >
      <DialogContent sx={{ p: 1 }}>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          {/* Target summary banner */}
          <Box
            sx={{
              p: 1.5,
              borderRadius: 1,
              backgroundColor: theme.palette.background.paper,
              border: `1px solid ${theme.palette.divider}`,
              display: 'flex',
              flexDirection: 'column',
              gap: 0.5,
            }}
          >
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
                Target:
              </Typography>
              <Chip
                size="small"
                label={
                  targetType === 'component'
                    ? `${targetKind || 'Component'}: ${targetName}`
                    : `Relationship: ${targetName}`
                }
              />
              <Chip
                size="small"
                color={modificationType === 'delete' ? 'error' : 'primary'}
                label={modificationType.toUpperCase()}
              />
            </Box>
            <Typography variant="body2" color="text.secondary">
              {summary}
            </Typography>
          </Box>

          {/* Directly Affected Components */}
          {directlyAffected.length > 0 && (
            <Card
              sx={{
                p: 1.5,
                border: `1px solid ${theme.palette.divider}`,
                backgroundColor: theme.palette.background.paper,
              }}
            >
              <Typography variant="subtitle2" sx={{ fontWeight: 600, mb: 1 }}>
                Directly Affected Components ({directlyAffected.length})
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                {directlyAffected.map((item: ImpactedComponent) => (
                  <Box
                    key={item.id}
                    sx={{
                      p: 1,
                      borderRadius: 1,
                      backgroundColor: theme.palette.background.default,
                      border: `1px solid ${theme.palette.divider}`,
                    }}
                  >
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                      <Typography variant="body2" sx={{ fontWeight: 600 }}>
                        {item.kind}/{item.name}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        ({item.model})
                      </Typography>
                    </Box>
                    <Typography variant="caption" color="text.secondary">
                      {item.explanation}
                    </Typography>
                  </Box>
                ))}
              </Box>
            </Card>
          )}

          {/* Indirectly Affected Components */}
          {indirectlyAffected.length > 0 && (
            <Card
              sx={{
                p: 1.5,
                border: `1px solid ${theme.palette.divider}`,
                backgroundColor: theme.palette.background.paper,
              }}
            >
              <Typography variant="subtitle2" sx={{ fontWeight: 600, mb: 1 }}>
                Indirectly / Transitively Affected ({indirectlyAffected.length})
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                {indirectlyAffected.map((item: ImpactedComponent) => (
                  <Box
                    key={item.id}
                    sx={{
                      p: 1,
                      borderRadius: 1,
                      backgroundColor: theme.palette.background.default,
                      border: `1px solid ${theme.palette.divider}`,
                    }}
                  >
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                      <Typography variant="body2" sx={{ fontWeight: 600 }}>
                        {item.kind}/{item.name}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        ({item.model})
                      </Typography>
                    </Box>
                    <Typography variant="caption" color="text.secondary">
                      {item.explanation}
                    </Typography>
                  </Box>
                ))}
              </Box>
            </Card>
          )}

          {/* Connected Relationships */}
          {affectedRelationships.length > 0 && (
            <Card
              sx={{
                p: 1.5,
                border: `1px solid ${theme.palette.divider}`,
                backgroundColor: theme.palette.background.paper,
              }}
            >
              <Typography variant="subtitle2" sx={{ fontWeight: 600, mb: 1 }}>
                Involved Relationships ({affectedRelationships.length})
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.75 }}>
                {affectedRelationships.map((rel: ImpactedRelationship, idx: number) => {
                  const relTitle = rel.subType || rel.type || `${rel.kind || 'Edge'} connection`;
                  return (
                    <Box
                      key={rel.id || idx}
                      sx={{
                        p: 0.75,
                        borderRadius: 1,
                        backgroundColor: theme.palette.background.default,
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                      }}
                    >
                      <Typography variant="body2" sx={{ fontWeight: 500 }}>
                        {rel.kind || 'Edge'} &rarr; {relTitle}
                      </Typography>
                      {rel.id && (
                        <Typography variant="caption" color="text.secondary">
                          {rel.id}
                        </Typography>
                      )}
                    </Box>
                  );
                })}
              </Box>
            </Card>
          )}

          {totalAffectedCount === 0 && affectedRelationships.length === 0 && (
            <Box
              sx={{
                p: 2,
                borderRadius: 1,
                backgroundColor: theme.palette.background.paper,
                textAlign: 'center',
              }}
            >
              <Typography variant="body2" color="text.secondary">
                No dependencies or connected relationships will be disrupted by this change.
              </Typography>
            </Box>
          )}
        </Box>
      </DialogContent>
    </Modal>
  );
};

export default DesignImpactPreviewModal;
