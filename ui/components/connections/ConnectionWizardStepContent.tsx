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
import { ReactNode, useMemo } from 'react';
import { alpha, styled, useTheme } from '@/theme';
import RJSFWrapper from '@/components/meshery-mesh-interface/PatternService/RJSF_wrapper';
import ConnectionIcon from '@/assets/icons/Connection';
import { getFallbackImageBasedOnKind, normalizeStaticImagePath } from '@/utils/fallback';
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
  // A native <button> defaults to the UA text color; set it explicitly so the
  // Typography children (which inherit) follow the theme in dark mode.
  color: theme.palette.text.primary,
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
  // A light, neutral tile so colored integration logos read clearly on the
  // dark card surface.
  background: theme.palette.common.white,
  border: `1px solid ${theme.palette.divider}`,
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
  color: theme.palette.text.primary,
  cursor: 'pointer',
  textAlign: 'center',
  transition: 'border-color 0.15s ease, background 0.15s ease',
  '&:hover': {
    borderColor: theme.palette.background.brand.default,
    background: alpha(theme.palette.background.brand.default, 0.04),
  },
}));

const UploadIconCircle = styled(Box)(({ theme }) => ({
  width: 56,
  height: 56,
  borderRadius: '50%',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: alpha(theme.palette.background.brand.default, 0.12),
}));

const UploadIcon = styled(CloudUploadIcon)(({ theme }) => ({
  width: 28,
  height: 28,
  fill: theme.palette.background.brand.default,
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
}: ConnectionKindSelectionStepProps) => {
  const theme = useTheme();
  return (
    <StepLayout>
      <StepHeader
        title="Choose a connection type"
        subtitle="Select the type of managed or unmanaged connection that you would like to register."
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
            // Prefer the SVG carried on the connection definition, then the redux
            // connection metadata icon, then the per-kind static asset.
            // normalizeStaticImagePath turns inline SVG markup into a data URI and
            // normalizes repo-relative paths.
            //
            // The icon tile (KindIconWrap) is always a white surface, so always
            // use the COLOR variant — the white variant (e.g. Kubernetes' all-white
            // logo) would render white-on-white and disappear.
            const definitionSvg = config.svgColor || config.svgWhite;
            const iconSrc =
              normalizeStaticImagePath(definitionSvg) ||
              normalizeStaticImagePath(connectionIconMap?.[config.kind]?.icon) ||
              normalizeStaticImagePath(getFallbackImageBasedOnKind(config.kind));
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
                    <ConnectionIcon
                      height={28}
                      width={28}
                      fill={theme.palette.background.brand.default}
                    />
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
};

// Derive an RJSF uiSchema from a connection/credential JSON schema so each
// field gets a placeholder hint. Meshery's RJSF widgets already surface a
// field's `description` as a hover info-button; we mirror that text into a
// `ui:placeholder` so the guidance is also visible inline in the empty input.
const buildPlaceholderUiSchema = (
  schema: Record<string, unknown> | null,
): Record<string, { 'ui:placeholder': string }> => {
  const properties = (schema?.properties ?? {}) as Record<
    string,
    { description?: string; title?: string }
  >;
  return Object.entries(properties).reduce<Record<string, { 'ui:placeholder': string }>>(
    (acc, [key, prop]) => {
      const placeholder = prop?.description || prop?.title;
      if (placeholder) {
        acc[key] = { 'ui:placeholder': placeholder };
      }
      return acc;
    },
    {},
  );
};

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
}: GenericConnectionDetailsStepProps) => {
  const uiSchema = useMemo(() => buildPlaceholderUiSchema(schema), [schema]);
  return (
    <StepLayout>
      <StepHeader
        title={`Configure ${label ?? 'connection'}`}
        subtitle="Fill in the fields below to set up this connection."
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
            uiSchema={uiSchema}
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
};

type CredentialAssociationStepProps = {
  label?: string;
  existingCredentials: CredentialOption[];
  credentialMode: 'existing' | 'new';
  selectedCredentialId: string;
  credentialName: string;
  onCredentialNameChange: (name: string) => void;
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
  credentialName,
  onCredentialNameChange,
  onCredentialModeChange,
  onSelectedCredentialChange,
  credentialSchema,
  credentialFormData,
  formRef,
  onCredentialFormChange,
  skipCredentialVerification,
  onSkipCredentialVerificationChange,
}: CredentialAssociationStepProps) => {
  const credentialUiSchema = useMemo(
    () => buildPlaceholderUiSchema(credentialSchema),
    [credentialSchema],
  );
  return (
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
        <>
          <TextField
            fullWidth
            label="Credential name"
            value={credentialName}
            onChange={(event) => onCredentialNameChange(event.target.value)}
            helperText="A name to identify this credential. Defaults to the connection name."
          />
          <FormContainer>
            <RJSFWrapper
              key={`${label}-credential-form`}
              jsonSchema={credentialSchema}
              uiSchema={credentialUiSchema}
              formData={credentialFormData}
              formRef={formRef}
              liveValidate={false}
              hideTitle
              onChange={onCredentialFormChange}
            />
          </FormContainer>
        </>
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
            Skip the connectivity check and register this connection without verifying the
            credential.
          </Typography>
        </Box>
      </InlineNotice>
    </StepLayout>
  );
};

type KubernetesImportStepProps = {
  kubeconfigFile: File | null;
  onPickFile: (file: File | null) => void;
};

export const KubernetesImportStep = ({ kubeconfigFile, onPickFile }: KubernetesImportStepProps) => (
  <StepLayout>
    <StepHeader
      title="Upload a kubeconfig"
      subtitle="Upload a kubeconfig file. Meshery will read the Kubernetes contexts inside and let you choose which ones to import."
    />
    <UploadDropzone
      type="button"
      onClick={() => document.getElementById('connection-wizard-kubeconfig-input')?.click()}
    >
      <UploadIconCircle>
        <UploadIcon />
      </UploadIconCircle>
      {kubeconfigFile ? (
        <Typography
          variant="body1"
          title={kubeconfigFile.name}
          sx={{
            fontWeight: 600,
            maxWidth: '100%',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          }}
        >
          {kubeconfigFile.name}
        </Typography>
      ) : (
        <Typography variant="body1" sx={{ fontWeight: 600 }}>
          Click to choose a kubeconfig file
        </Typography>
      )}
      <Typography variant="caption" color="text.secondary">
        {kubeconfigFile
          ? 'Click to replace the selected file'
          : 'Accepts kubeconfigs with embedded certificates'}
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
