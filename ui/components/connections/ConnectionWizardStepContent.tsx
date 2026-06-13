import {
  Box,
  Checkbox,
  CheckCircleIcon,
  CircularProgress,
  CloudUploadIcon,
  MenuItem,
  TextField,
  Typography,
} from '@sistent/sistent';
import { ReactNode } from 'react';
import { alpha, styled } from '@/theme';
import RJSFWrapper from '@/components/meshery-mesh-interface/PatternService/RJSF_wrapper';
import ConnectionIcon from '@/assets/icons/Connection';
import {
  type ConnectionWizardKindConfig,
  type SupportedConnectionWizardKind,
} from './ConnectionWizard.helpers';

type RjsfFormHandle = {
  validateForm: () => boolean;
  state: { formData: Record<string, unknown> };
};

type CredentialOption = {
  id?: string;
  name?: string;
};

const StepLayout = styled(Box)(({ theme }) => ({
  display: 'grid',
  gap: theme.spacing(3),
}));

// ---------------------------------------------------------------------------
// Shared step header: a bold title with a muted one-line subtitle. Replaces the
// loose body paragraphs each step used to lead with, giving every step a
// consistent visual hierarchy.
// ---------------------------------------------------------------------------

const StepHeader = ({ title, subtitle }: { title: string; subtitle?: ReactNode }) => (
  <Box sx={{ display: 'grid', gap: 0.5 }}>
    <Typography variant="h6" sx={{ fontWeight: 600 }}>
      {title}
    </Typography>
    {subtitle && (
      <Typography variant="body2" color="text.secondary">
        {subtitle}
      </Typography>
    )}
  </Box>
);

const KindGrid = styled(Box)(({ theme }) => ({
  display: 'grid',
  gap: theme.spacing(2),
  gridTemplateColumns: 'repeat(auto-fill, minmax(220px, 1fr))',
}));

const KindCard = styled('button', {
  shouldForwardProp: (prop) => prop !== 'selected',
})<{ selected?: boolean }>(({ theme, selected }) => ({
  position: 'relative',
  textAlign: 'left',
  borderRadius: theme.spacing(1.5),
  border: `1.5px solid ${selected ? theme.palette.background.brand.default : theme.palette.divider}`,
  background: selected
    ? alpha(theme.palette.background.brand.default, 0.06)
    : theme.palette.background.card,
  padding: theme.spacing(2.5),
  cursor: 'pointer',
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(1),
  minHeight: 150,
  transition: 'transform 0.15s ease, box-shadow 0.15s ease, border-color 0.15s ease',
  boxShadow: selected ? `0 0 0 1px ${theme.palette.background.brand.default}` : 'none',
  '&:hover:not(:disabled)': {
    borderColor: theme.palette.background.brand.default,
    transform: 'translateY(-2px)',
    boxShadow: theme.shadows[4],
  },
  '&:focus-visible': {
    outline: `2px solid ${theme.palette.background.brand.default}`,
    outlineOffset: 2,
  },
  '&:disabled': {
    opacity: 0.55,
    cursor: 'not-allowed',
  },
}));

const KindIconWrap = styled(Box)(({ theme }) => ({
  width: 48,
  height: 48,
  borderRadius: theme.spacing(1),
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: alpha(theme.palette.background.brand.default, 0.1),
  marginBottom: theme.spacing(0.5),
}));

const KindIcon = styled('img')({
  width: 30,
  height: 30,
  objectFit: 'contain',
});

const SelectedBadge = styled(CheckCircleIcon)(({ theme }) => ({
  position: 'absolute',
  top: theme.spacing(1.5),
  right: theme.spacing(1.5),
  width: 20,
  height: 20,
  fill: theme.palette.background.brand.default,
}));

const Segmented = styled(Box)(({ theme }) => ({
  display: 'inline-flex',
  padding: theme.spacing(0.5),
  gap: theme.spacing(0.5),
  borderRadius: theme.spacing(1),
  border: `1px solid ${theme.palette.divider}`,
  background: theme.palette.background.card,
  width: 'fit-content',
}));

const SegmentButton = styled('button', {
  shouldForwardProp: (prop) => prop !== 'active',
})<{ active?: boolean }>(({ theme, active }) => ({
  border: 'none',
  cursor: 'pointer',
  borderRadius: theme.spacing(0.75),
  padding: theme.spacing(0.75, 2),
  fontSize: '0.8125rem',
  fontWeight: 600,
  transition: 'background 0.15s ease, color 0.15s ease',
  background: active ? theme.palette.background.brand.default : 'transparent',
  color: active ? theme.palette.common.white : theme.palette.text.secondary,
  '&:hover': {
    color: active ? theme.palette.common.white : theme.palette.text.primary,
  },
}));

const SummaryCard = styled(Box)(({ theme }) => ({
  borderRadius: theme.spacing(1.5),
  border: `1px solid ${theme.palette.divider}`,
  background: theme.palette.background.card,
  overflow: 'hidden',
}));

const SummaryRow = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: theme.spacing(2),
  padding: theme.spacing(1.75, 2.5),
  '& + &': {
    borderTop: `1px solid ${theme.palette.divider}`,
  },
}));

const SummaryLabel = styled(Typography)(({ theme }) => ({
  color: theme.palette.text.secondary,
  flexShrink: 0,
}));

const SummaryValue = styled(Typography)({
  textAlign: 'right',
  wordBreak: 'break-word',
});

const InlineNotice = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'flex-start',
  gap: theme.spacing(1),
  borderRadius: theme.spacing(1),
  border: `1px solid ${theme.palette.divider}`,
  background: alpha(theme.palette.info.main, 0.06),
  padding: theme.spacing(1, 1.5),
}));

const UploadDropzone = styled('button')(({ theme }) => ({
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  gap: theme.spacing(1),
  padding: theme.spacing(4),
  borderRadius: theme.spacing(1.5),
  border: `1.5px dashed ${theme.palette.divider}`,
  background: theme.palette.background.card,
  cursor: 'pointer',
  textAlign: 'center',
  transition: 'border-color 0.15s ease, background 0.15s ease',
  '&:hover': {
    borderColor: theme.palette.background.brand.default,
    background: alpha(theme.palette.background.brand.default, 0.04),
  },
}));

const UploadIcon = styled(CloudUploadIcon)(({ theme }) => ({
  width: 36,
  height: 36,
  fill: theme.palette.text.secondary,
}));

const HiddenFileInput = styled('input')({
  display: 'none',
});

const FormContainer = styled(Box)(({ theme }) => ({
  borderRadius: theme.spacing(1.5),
  border: `1px solid ${theme.palette.divider}`,
  background: theme.palette.background.card,
  padding: theme.spacing(2.5),
}));

type ConnectionKindSelectionStepProps = {
  kinds: ConnectionWizardKindConfig[];
  isLoading: boolean;
  selectedKind: SupportedConnectionWizardKind | null;
  connectionIconMap: Record<string, { icon?: string }>;
  onSelectKind: (kind: SupportedConnectionWizardKind) => void;
  canUseKind: (config: ConnectionWizardKindConfig) => boolean;
};

export const ConnectionKindSelectionStep = ({
  kinds,
  isLoading,
  selectedKind,
  connectionIconMap,
  onSelectKind,
  canUseKind,
}: ConnectionKindSelectionStepProps) => (
  <StepLayout>
    <StepHeader
      title="Choose a connection type"
      subtitle="Select the kind of first-class connection you want Meshery to create or register."
    />
    {isLoading ? (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 6 }}>
        <CircularProgress />
      </Box>
    ) : kinds.length === 0 ? (
      <Typography variant="body2" color="text.secondary">
        No connection definitions are registered in the registry yet.
      </Typography>
    ) : (
      <KindGrid>
        {kinds.map((config) => {
          const isPermitted = canUseKind(config);
          const iconSrc = connectionIconMap?.[config.kind]?.icon;
          const isSelected = selectedKind === config.kind;

          return (
            <KindCard
              key={config.kind}
              type="button"
              onClick={() => isPermitted && onSelectKind(config.kind)}
              selected={isSelected}
              disabled={!isPermitted}
            >
              {isSelected && <SelectedBadge />}
              <KindIconWrap>
                {iconSrc ? (
                  <KindIcon src={iconSrc} alt={`${config.label} icon`} />
                ) : (
                  <ConnectionIcon height={30} width={30} />
                )}
              </KindIconWrap>
              <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                {config.label}
              </Typography>
              {config.description && (
                <Typography variant="body2" color="text.secondary">
                  {config.description}
                </Typography>
              )}
              {!isPermitted && (
                <Typography variant="caption" color="error">
                  You don&apos;t have permission to add this connection type.
                </Typography>
              )}
            </KindCard>
          );
        })}
      </KindGrid>
    )}
  </StepLayout>
);

type GenericConnectionDetailsStepProps = {
  label?: string;
  isInitializing: boolean;
  schema: Record<string, unknown> | null;
  formData: Record<string, unknown>;
  formRef: { current: RjsfFormHandle | null };
  onChange: (formData: Record<string, unknown>) => void;
};

export const GenericConnectionDetailsStep = ({
  label,
  isInitializing,
  schema,
  formData,
  formRef,
  onChange,
}: GenericConnectionDetailsStepProps) => (
  <StepLayout>
    <StepHeader
      title={`Configure ${label ?? 'connection'}`}
      subtitle="These fields are rendered from the connection definition's registration schema."
    />
    {isInitializing || !schema ? (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 6 }}>
        <CircularProgress />
      </Box>
    ) : (
      <FormContainer>
        <RJSFWrapper
          key={`${label}-connection-form`}
          jsonSchema={schema}
          formData={formData}
          formRef={formRef}
          liveValidate={false}
          hideTitle
          onChange={onChange}
        />
      </FormContainer>
    )}
  </StepLayout>
);

type CredentialAssociationStepProps = {
  label?: string;
  existingCredentials: CredentialOption[];
  credentialMode: 'existing' | 'new';
  selectedCredentialId: string;
  onCredentialModeChange: (mode: 'existing' | 'new') => void;
  onSelectedCredentialChange: (id: string) => void;
  credentialSchema: Record<string, unknown> | null;
  credentialFormData: Record<string, unknown>;
  formRef: { current: RjsfFormHandle | null };
  onCredentialFormChange: (formData: Record<string, unknown>) => void;
  skipCredentialVerification: boolean;
  onSkipCredentialVerificationChange: (value: boolean) => void;
};

export const CredentialAssociationStep = ({
  label,
  existingCredentials,
  credentialMode,
  selectedCredentialId,
  onCredentialModeChange,
  onSelectedCredentialChange,
  credentialSchema,
  credentialFormData,
  formRef,
  onCredentialFormChange,
  skipCredentialVerification,
  onSkipCredentialVerificationChange,
}: CredentialAssociationStepProps) => (
  <StepLayout>
    <StepHeader
      title="Associate a credential"
      subtitle={`Reuse an existing credential or create a new one for this ${label ?? ''} connection.`}
    />
    {existingCredentials.length > 0 && (
      <>
        <Segmented role="tablist">
          <SegmentButton
            type="button"
            active={credentialMode === 'existing'}
            onClick={() => onCredentialModeChange('existing')}
          >
            Use existing
          </SegmentButton>
          <SegmentButton
            type="button"
            active={credentialMode === 'new'}
            onClick={() => onCredentialModeChange('new')}
          >
            Create new
          </SegmentButton>
        </Segmented>
        {credentialMode === 'existing' && (
          <TextField
            select
            fullWidth
            label="Existing credential"
            value={selectedCredentialId}
            onChange={(event) => onSelectedCredentialChange(event.target.value)}
          >
            <MenuItem value="">Select a credential</MenuItem>
            {existingCredentials.map((credential) => (
              <MenuItem key={credential.id} value={credential.id}>
                {credential.name}
              </MenuItem>
            ))}
          </TextField>
        )}
      </>
    )}
    {(credentialMode === 'new' || existingCredentials.length === 0) && credentialSchema && (
      <FormContainer>
        <RJSFWrapper
          key={`${label}-credential-form`}
          jsonSchema={credentialSchema}
          formData={credentialFormData}
          formRef={formRef}
          liveValidate={false}
          hideTitle
          onChange={onCredentialFormChange}
        />
      </FormContainer>
    )}
    <InlineNotice
      onClick={() => onSkipCredentialVerificationChange(!skipCredentialVerification)}
      sx={{ cursor: 'pointer' }}
    >
      <Checkbox
        checked={skipCredentialVerification}
        onChange={(event) => onSkipCredentialVerificationChange(event.target.checked)}
        sx={{ p: 0, mt: 0.25 }}
        onClick={(event) => event.stopPropagation()}
      />
      <Box>
        <Typography variant="body2" sx={{ fontWeight: 600 }}>
          Bypass credential verification
        </Typography>
        <Typography variant="caption" color="text.secondary">
          Skip the connectivity check and register this connection without verifying the credential.
        </Typography>
      </Box>
    </InlineNotice>
  </StepLayout>
);

type KubernetesImportStepProps = {
  kubeconfigFile: File | null;
  onPickFile: (file: File | null) => void;
};

export const KubernetesImportStep = ({ kubeconfigFile, onPickFile }: KubernetesImportStepProps) => (
  <StepLayout>
    <StepHeader
      title="Import a kubeconfig"
      subtitle="Meshery discovers the contexts inside the file and registers each one as a Kubernetes connection."
    />
    <UploadDropzone
      type="button"
      onClick={() => document.getElementById('connection-wizard-kubeconfig-input')?.click()}
    >
      <UploadIcon />
      {kubeconfigFile ? (
        <Typography variant="body1" sx={{ fontWeight: 600 }}>
          {kubeconfigFile.name}
        </Typography>
      ) : (
        <Typography variant="body1" sx={{ fontWeight: 600 }}>
          Click to choose a kubeconfig file
        </Typography>
      )}
      <Typography variant="caption" color="text.secondary">
        {kubeconfigFile ? 'Click to replace the selected file' : 'Accepts a standard kubeconfig'}
      </Typography>
    </UploadDropzone>
    <HiddenFileInput
      id="connection-wizard-kubeconfig-input"
      type="file"
      onChange={(event) => onPickFile(event.target.files?.[0] || null)}
    />
  </StepLayout>
);

type ConnectionReviewStepProps = {
  isKubernetes: boolean;
  label?: string;
  kubeconfigFile: File | null;
  connectionName: string;
  credentialName: string;
  skipCredentialVerification: boolean;
};

export const ConnectionReviewStep = ({
  isKubernetes,
  label,
  kubeconfigFile,
  connectionName,
  credentialName,
  skipCredentialVerification,
}: ConnectionReviewStepProps) => (
  <StepLayout>
    <StepHeader
      title={isKubernetes ? 'Review Kubernetes import' : 'Review connection'}
      subtitle="Confirm the details below before Meshery registers the connection."
    />
    <SummaryCard>
      <SummaryRow>
        <SummaryLabel variant="body2">Kind</SummaryLabel>
        <SummaryValue variant="body2">{label}</SummaryValue>
      </SummaryRow>
      {isKubernetes ? (
        <>
          <SummaryRow>
            <SummaryLabel variant="body2">File</SummaryLabel>
            <SummaryValue variant="body2">
              {kubeconfigFile?.name || 'No kubeconfig selected'}
            </SummaryValue>
          </SummaryRow>
          <SummaryRow>
            <SummaryLabel variant="body2">Behavior</SummaryLabel>
            <SummaryValue variant="body2">
              Every reachable context is discovered and registered as a Kubernetes connection.
            </SummaryValue>
          </SummaryRow>
          <SummaryRow>
            <SummaryLabel variant="body2">Credential</SummaryLabel>
            <SummaryValue variant="body2">
              The kubeconfig content is associated with each imported connection.
            </SummaryValue>
          </SummaryRow>
        </>
      ) : (
        <>
          <SummaryRow>
            <SummaryLabel variant="body2">Name</SummaryLabel>
            <SummaryValue variant="body2">{connectionName}</SummaryValue>
          </SummaryRow>
          <SummaryRow>
            <SummaryLabel variant="body2">Credential</SummaryLabel>
            <SummaryValue variant="body2">
              {credentialName || 'No credential selected'}
            </SummaryValue>
          </SummaryRow>
          <SummaryRow>
            <SummaryLabel variant="body2">Verification</SummaryLabel>
            <SummaryValue variant="body2">
              {skipCredentialVerification ? 'Bypassed' : 'Enabled'}
            </SummaryValue>
          </SummaryRow>
        </>
      )}
    </SummaryCard>
  </StepLayout>
);

// Re-exported so step modules can compose the shared header for steps that
// don't use one of the components above.
export { StepHeader };
