import Link from 'next/link';
import { Box, CheckCircleIcon, Chip, ErrorIcon, Typography, WarningIcon } from '@sistent/sistent';
import { useSelector } from 'react-redux';
import { alpha, styled } from '@/theme';
import {
  CONNECTIONS_PATH,
  kubernetesImportPresentation,
  type KubernetesImportIssue,
  type KubernetesImportIssueKind,
  type KubernetesImportPresentation,
} from '../ConnectionWizard.helpers';
import FormatConnectionMetadata from '../metadata';
import { ConnectionStateChip } from '../ConnectionChip';
import { StepHeader } from '../ConnectionWizardStepContent';
import type { GenericRecord, WizardContext, WizardStep } from './types';

const KUBERNETES_CONNECTION_DOCS_URL =
  'https://docs.meshery.io/guides/infrastructure-management/kubernetes-connection-lifecycle';

const OutcomeBadge = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'tone',
})<{ tone: KubernetesImportPresentation }>(({ theme, tone }) => {
  const color =
    tone === 'partial'
      ? theme.palette.warning.main
      : tone === 'failure'
        ? theme.palette.error.main
        : theme.palette.background.brand?.default;
  return {
    width: 88,
    height: 88,
    borderRadius: '50%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: alpha(color, 0.12),
  };
});

const SuccessIcon = styled(CheckCircleIcon)(({ theme }) => ({
  width: 52,
  height: 52,
  fill: theme.palette.background.brand?.default,
}));

const PartialIcon = styled(WarningIcon)(({ theme }) => ({
  width: 52,
  height: 52,
  fill: theme.palette.warning.main,
}));

const FailureIcon = styled(ErrorIcon)(({ theme }) => ({
  width: 52,
  height: 52,
  fill: theme.palette.error.main,
}));

const IMPORT_ISSUE_LABEL: Record<KubernetesImportIssueKind, string> = {
  errored: 'Errored',
  ignored: 'Ignored',
  'failed-connect': 'Failed to connect',
};

const ConnectionsLink = styled(Link)(({ theme }) => ({
  color: theme.palette.background.brand?.default || theme.palette.primary.main,
  fontWeight: 600,
  textDecoration: 'underline',
  textUnderlineOffset: 2,
}));

// ---------------------------------------------------------------------------
// 5. Receipt
// ---------------------------------------------------------------------------

const KubernetesReceiptBody = ({ ctx }: { ctx: WizardContext }) => {
  const controllerState = useSelector(
    (state: { ui: { controllerState: unknown } }) => state.ui.controllerState,
  );

  // Configure mode operates on a single connection: show the same live detail
  // receipt that the Connections table renders when a row is expanded.
  // That includes status chips, controller versions, and the Diagnostics section.
  const configuredConnection = (ctx.data.registrationResult as GenericRecord) || {};
  if (
    ctx.mode === 'configure' &&
    configuredConnection.kind === 'kubernetes' &&
    configuredConnection.id
  ) {
    return (
      <Box sx={{ display: 'grid', gap: 2 }}>
        <StepHeader
          title="Connection details"
          subtitle="Live status and diagnostics for this Kubernetes connection."
        />
        <FormatConnectionMetadata
          connection={configuredConnection}
          meshsyncControllerState={controllerState}
        />
      </Box>
    );
  }

  const imported =
    (ctx.data.postConfig.importedContexts as { name: string; status: string }[]) || [];
  const created = Number(ctx.data.postConfig.createdCount ?? imported.length);
  const connected = Number(ctx.data.postConfig.connectedCount ?? 0);
  const unreachable = Number(ctx.data.postConfig.unreachableCount ?? 0);
  const isImport = ctx.data.postConfig.createdCount !== undefined;
  const issues = (ctx.data.postConfig.importIssues as KubernetesImportIssue[]) || [];
  const presentation: KubernetesImportPresentation = isImport
    ? kubernetesImportPresentation(created, issues.length)
    : 'success';
  const OutcomeIcon =
    presentation === 'failure'
      ? FailureIcon
      : presentation === 'partial'
        ? PartialIcon
        : SuccessIcon;
  const heading = isImport
    ? presentation === 'failure'
      ? 'Kubernetes import failed'
      : presentation === 'partial'
        ? 'Kubernetes import completed with issues'
        : 'Kubernetes import complete'
    : 'Configuration saved';

  const summaryParts = [
    `${created} imported`,
    connected > 0 ? `${connected} connected` : null,
    unreachable > 0 ? `${unreachable} not found` : null,
  ].filter(Boolean);
  const summary =
    presentation === 'failure' ? 'Nothing was imported.' : `${summaryParts.join(' · ')}.`;

  return (
    <Box
      sx={{ display: 'grid', gap: 2, justifyItems: 'center', textAlign: 'center', py: 4, px: 2 }}
    >
      <OutcomeBadge
        tone={presentation}
        data-testid="kubernetes-import-outcome"
        data-outcome={isImport ? presentation : 'saved'}
      >
        <OutcomeIcon />
      </OutcomeBadge>
      <Typography
        variant="h6"
        sx={{
          fontWeight: 600,
          ...(presentation === 'failure' || presentation === 'partial'
            ? {
                color: (theme) =>
                  presentation === 'failure'
                    ? theme.palette.error.main
                    : theme.palette.warning.main,
              }
            : {}),
        }}
      >
        {heading}
      </Typography>
      {isImport && (
        <Typography variant="body2" color="text.secondary" sx={{ maxWidth: 440 }}>
          {summary} Manage these anytime from{' '}
          <ConnectionsLink href={CONNECTIONS_PATH}>connections</ConnectionsLink>.
        </Typography>
      )}
      {issues.length > 0 && (
        <Box sx={{ display: 'grid', gap: 1, width: '100%', maxWidth: 460, textAlign: 'left' }}>
          {issues.map((issue, index) => (
            <Box
              key={`${issue.kind}-${issue.name}-${index}`}
              data-testid="kubernetes-import-issue"
              sx={{ display: 'grid', gap: 0.25 }}
            >
              <Typography variant="body2" sx={{ fontWeight: 600 }}>
                {issue.name || 'Unnamed context'}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                {IMPORT_ISSUE_LABEL[issue.kind]}
                {issue.reason ? `: ${issue.reason}` : ''}
              </Typography>
            </Box>
          ))}
        </Box>
      )}
      {imported.length > 0 && (
        <Box sx={{ display: 'grid', gap: 1, width: '100%', maxWidth: 460, textAlign: 'left' }}>
          {imported.map((context, index) => (
            <Box
              key={`${context.name}-${index}`}
              sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                gap: 1.5,
                minWidth: 0,
              }}
            >
              <Chip
                label={context.name || 'Unnamed context'}
                size="small"
                title={context.name}
                sx={{
                  minWidth: 0,
                  flex: '1 1 auto',
                  '& .MuiChip-label': {
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                  },
                }}
              />
              <ConnectionStateChip status={context.status} />
            </Box>
          ))}
        </Box>
      )}
    </Box>
  );
};

export const kubernetesReceiptStep: WizardStep = {
  id: 'kubernetes-receipt',
  label: 'Done',
  Component: KubernetesReceiptBody,
  nextLabel: () => 'Finish',
  // Create mode has already imported. Back would return to Review, where
  // Import runs again. Configure mode can still step back to settings.
  hideBack: (ctx) => ctx.mode === 'create',
  helpText: (ctx) => {
    const created = Number(ctx.data.postConfig.createdCount ?? 0);
    const issues = (ctx.data.postConfig.importIssues as KubernetesImportIssue[]) || [];
    const presentation =
      ctx.data.postConfig.createdCount !== undefined
        ? kubernetesImportPresentation(created, issues.length)
        : 'success';
    const lead =
      presentation === 'failure'
        ? 'Nothing was imported.'
        : presentation === 'partial'
          ? 'Import completed with issues.'
          : 'Import complete.';
    return `${lead} Manage connections anytime from the Connections page. [Learn more about Kubernetes connection lifecycle](${KUBERNETES_CONNECTION_DOCS_URL}).`;
  },
};
