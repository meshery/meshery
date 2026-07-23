import { useEffect } from 'react';
import Link from 'next/link';
import { v4 as uuidv4 } from 'uuid';
import {
  Alert,
  Box,
  CheckCircleIcon,
  Typography,
  SettingsIcon,
  LockIcon,
  AssignmentTurnedInIcon,
} from '@sistent/sistent';
import { alpha, styled } from '@/theme';
import { EVENT_TYPES } from 'lib/event-types';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import {
  buildCredentialSecret,
  CONNECTIONS_PATH,
  connectionCreatedNotify,
  DEFAULT_CONNECTION_DOCS_URL,
  filterCredentialsForKind,
  resolveConnectionName,
  type ConnectionWizardKindConfig,
} from '../ConnectionWizard.helpers';
import { formatWizardError } from './errors';
import {
  ConnectionKindSelectionStep,
  ConnectionReviewStep,
  CredentialAssociationStep,
  GenericConnectionDetailsStep,
} from '../ConnectionWizardStepContent';
import type { WizardContext, WizardStep } from './types';

export const kindPermission = (config?: ConnectionWizardKindConfig | null) => {
  if (!config) {
    return false;
  }

  return config.flow === 'kubernetes'
    ? CAN(Keys.LifecycleManagementAddCluster.id, Keys.LifecycleManagementAddCluster.function)
    : CAN(Keys.MesherySystemConnectMetrics.id, Keys.MesherySystemConnectMetrics.function);
};

const existingCredentialsFor = (ctx: WizardContext) =>
  filterCredentialsForKind(ctx.services.credentials, ctx.data.kindConfig?.kind);

// ---------------------------------------------------------------------------
// 1. Choose connection
// ---------------------------------------------------------------------------

const SelectStepBody = ({ ctx }: { ctx: WizardContext }) => (
  <ConnectionKindSelectionStep
    kinds={ctx.data.availableKinds}
    isLoading={ctx.data.isLoadingKinds}
    selectedKind={ctx.data.kindConfig?.kind ?? null}
    connectionIconMap={ctx.data.connectionIconMap}
    onSelectKind={(kind) => {
      if (kind === ctx.data.kindConfig?.kind) {
        return;
      }
      const kindConfig = ctx.data.availableKinds.find((config) => config.kind === kind) || null;
      // Reset every downstream field so switching kinds never leaks state.
      ctx.patch({
        kindConfig,
        connectionFormData: {},
        credentialFormData: {},
        selectedCredentialId: '',
        credentialMode: 'existing',
        skipCredentialVerification: false,
        kubeconfigFile: null,
        registrationId: null,
        registrationResult: null,
        registrationError: null,
        postConfig: {},
      });
      ctx.advance();
    }}
    canUseKind={kindPermission}
  />
);

export const selectStep: WizardStep = {
  id: 'select',
  label: 'Choose Connection',
  Component: SelectStepBody,
  canProceed: (ctx) => Boolean(ctx.data.kindConfig && kindPermission(ctx.data.kindConfig)),
  helpText: `Choose the type of connection you want to create. Each connection type lets Meshery manage a different kind of infrastructure and interface with another managed service. [Learn more about connections](${DEFAULT_CONNECTION_DOCS_URL}).`,
};

// ---------------------------------------------------------------------------
// 2. Configure connection (connectionSchema)
// ---------------------------------------------------------------------------

const DetailsStepBody = ({ ctx }: { ctx: WizardContext }) => (
  <GenericConnectionDetailsStep
    label={ctx.data.kindConfig?.label}
    isInitializing={false}
    schema={ctx.data.kindConfig?.connectionSchema ?? null}
    formData={ctx.data.connectionFormData}
    formRef={ctx.formRefs.connection}
    onChange={(connectionFormData) => ctx.patch({ connectionFormData })}
  />
);

export const genericDetailsStep: WizardStep = {
  id: 'details',
  label: 'Configure Connection',
  icon: SettingsIcon,
  Component: DetailsStepBody,
  canProceed: (ctx) => Boolean(ctx.data.kindConfig?.connectionSchema),
  onNext: (ctx) => Boolean(ctx.formRefs.connection.current?.validateForm()),
  helpText: (ctx) => {
    const label = ctx.data.kindConfig?.label || 'connection';
    const docsUrl = ctx.data.kindConfig?.docsUrl || DEFAULT_CONNECTION_DOCS_URL;
    return `Fill in the details for your ${label} connection. These fields are specific to this connection type. [Learn more](${docsUrl}).`;
  },
};

// ---------------------------------------------------------------------------
// 3. Associate credential (credentialSchema)
// ---------------------------------------------------------------------------

const CredentialStepBody = ({ ctx }: { ctx: WizardContext }) => {
  const existingCredentials = existingCredentialsFor(ctx);
  const connectionName = resolveConnectionName(
    ctx.data.kindConfig?.kind || 'connection',
    ctx.data.connectionFormData,
  );

  // Default the credential name to the connection name. Seeded once when the
  // step mounts and only while the field is still untouched, so the default
  // tracks a freshly configured connection name without clobbering a name the
  // user has typed.
  useEffect(() => {
    if (!ctx.data.credentialName) {
      ctx.patch({ credentialName: connectionName });
    }
  }, []);

  return (
    <CredentialAssociationStep
      label={ctx.data.kindConfig?.label}
      existingCredentials={existingCredentials}
      credentialMode={ctx.data.credentialMode}
      selectedCredentialId={ctx.data.selectedCredentialId}
      credentialName={ctx.data.credentialName}
      onCredentialNameChange={(credentialName) => ctx.patch({ credentialName })}
      onCredentialModeChange={(credentialMode) => ctx.patch({ credentialMode })}
      onSelectedCredentialChange={(selectedCredentialId) => ctx.patch({ selectedCredentialId })}
      credentialSchema={ctx.data.kindConfig?.credentialSchema ?? null}
      credentialFormData={ctx.data.credentialFormData}
      formRef={ctx.formRefs.credential}
      onCredentialFormChange={(credentialFormData) => ctx.patch({ credentialFormData })}
      skipCredentialVerification={ctx.data.skipCredentialVerification}
      onSkipCredentialVerificationChange={(skipCredentialVerification) =>
        ctx.patch({ skipCredentialVerification })
      }
    />
  );
};

export const genericCredentialStep: WizardStep = {
  id: 'credential',
  label: 'Associate Credential',
  icon: LockIcon,
  Component: CredentialStepBody,
  hidden: (ctx) => !ctx.data.kindConfig?.credentialSchema,
  canProceed: (ctx) => {
    const hasExisting = existingCredentialsFor(ctx).length > 0;
    if (ctx.data.credentialMode === 'existing' && hasExisting) {
      return Boolean(ctx.data.selectedCredentialId);
    }
    return true;
  },
  onNext: (ctx) => {
    const hasExisting = existingCredentialsFor(ctx).length > 0;
    if (ctx.data.credentialMode === 'existing' && hasExisting) {
      return Boolean(ctx.data.selectedCredentialId);
    }
    return Boolean(ctx.formRefs.credential.current?.validateForm());
  },
  helpText: `Credentials store secrets that authenticate this connection. Reuse an existing credential or create a new one. [Learn more about connections](${DEFAULT_CONNECTION_DOCS_URL}).`,
};

// ---------------------------------------------------------------------------
// 4. Review & register (generic: register + connect)
// ---------------------------------------------------------------------------

const RegisterStepBody = ({ ctx }: { ctx: WizardContext }) => {
  const { kindConfig, registrationError } = ctx.data;
  const connectionName = resolveConnectionName(
    kindConfig?.kind || 'connection',
    ctx.data.connectionFormData,
  );
  const selectedCredential = existingCredentialsFor(ctx).find(
    (credential) => credential.id === ctx.data.selectedCredentialId,
  );
  const credentialName =
    selectedCredential?.name || ctx.data.credentialName.trim() || connectionName;

  return (
    <Box sx={{ display: 'grid', gap: 2 }}>
      <ConnectionReviewStep
        isKubernetes={false}
        label={kindConfig?.label}
        kubeconfigFile={null}
        connectionName={connectionName}
        credentialName={credentialName}
        skipCredentialVerification={ctx.data.skipCredentialVerification}
      />
      {Boolean(registrationError) && (
        <Alert severity="error" variant="filled">
          {formatWizardError(registrationError)}
        </Alert>
      )}
    </Box>
  );
};

const buildGenericPayload = (ctx: WizardContext) => {
  const { kindConfig } = ctx.data;
  // A connection carries its own kind/type/subType; use them directly rather
  // than deriving type/subType from a model's category/subCategory.
  const connectionType = kindConfig?.type;
  const connectionSubType = kindConfig?.subType;
  const selectedCredential = existingCredentialsFor(ctx).find(
    (credential) => credential.id === ctx.data.selectedCredentialId,
  );
  const connectionName = resolveConnectionName(
    kindConfig?.kind || 'connection',
    ctx.data.connectionFormData,
  );
  const credentialSecret = buildCredentialSecret(selectedCredential, ctx.data.credentialFormData);

  // For a newly created credential, stamp the user-supplied name (falling back
  // to the connection name) so it is persisted as a first-class, identifiable
  // credential rather than an unnamed secret.
  if (!selectedCredential) {
    credentialSecret.name = ctx.data.credentialName.trim() || connectionName;
  }

  return {
    id: ctx.data.registrationId,
    kind: kindConfig?.kind,
    name: connectionName,
    type: typeof connectionType === 'string' ? connectionType.toLowerCase() : undefined,
    subType: typeof connectionSubType === 'string' ? connectionSubType.toLowerCase() : undefined,
    metadata: ctx.data.connectionFormData,
    credentialSecret,
    skipCredentialVerification: ctx.data.skipCredentialVerification,
  };
};

export const genericRegisterStep: WizardStep = {
  id: 'register',
  label: 'Review & Create',
  icon: AssignmentTurnedInIcon,
  Component: RegisterStepBody,
  nextLabel: () => 'Create Connection',
  helpText: (ctx) => {
    const label = ctx.data.kindConfig?.label || 'connection';
    const docsUrl = ctx.data.kindConfig?.docsUrl || DEFAULT_CONNECTION_DOCS_URL;
    return `Review your ${label} details. When you click "Create Connection", Meshery will register the new connection, then verify reachability of the system/service, verify credentials, and run discovery (if applicable). [Learn more](${docsUrl}).`;
  },
  onNext: async (ctx) => {
    const { kindConfig } = ctx.data;
    if (!kindConfig) {
      return false;
    }
    ctx.patch({ registrationError: null });

    try {
      // Connection definitions (not registry components) are the source of
      // truth for the form schema, so the legacy `initialize` round-trip — which
      // looks up a `{Kind}Connection` component that no longer exists — is
      // skipped. The registration id is just a process tracker, so mint it
      // client-side; register/connect initialize the state machine themselves.
      const registrationId = ctx.data.registrationId || uuidv4();
      if (registrationId !== ctx.data.registrationId) {
        ctx.patch({ registrationId });
      }

      const basePayload = { ...buildGenericPayload(ctx), id: registrationId };
      await ctx.services.registerConnection({ ...basePayload, status: 'register' });
      const result = await ctx.services.connectConnection({ ...basePayload, status: 'connect' });

      ctx.patch({ registrationResult: result ?? basePayload });
      ctx.services.notify(connectionCreatedNotify(kindConfig.label));
      return true;
    } catch (error) {
      ctx.patch({ registrationError: error });
      ctx.services.notify({
        message: `Failed to create ${kindConfig.label} connection: ${formatWizardError(error)}`,
        event_type: EVENT_TYPES.ERROR,
      });
      return false;
    }
  },
};

// ---------------------------------------------------------------------------
// 5. Receipt
// ---------------------------------------------------------------------------

const SuccessBadge = styled(Box)(({ theme }) => ({
  width: 88,
  height: 88,
  borderRadius: '50%',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: alpha(theme.palette.background.brand?.default, 0.12),
}));

const SuccessIcon = styled(CheckCircleIcon)(({ theme }) => ({
  width: 52,
  height: 52,
  fill: theme.palette.background.brand?.default,
}));

const ConnectionsLink = styled(Link)(({ theme }) => ({
  color: theme.palette.background.brand?.default || theme.palette.primary.main,
  fontWeight: 600,
  textDecoration: 'underline',
  textUnderlineOffset: 2,
}));

const ReceiptStepBody = ({ ctx }: { ctx: WizardContext }) => {
  const { kindConfig } = ctx.data;
  const connectionName = resolveConnectionName(
    kindConfig?.kind || 'connection',
    ctx.data.connectionFormData,
  );

  return (
    <Box
      sx={{ display: 'grid', gap: 2, justifyItems: 'center', textAlign: 'center', py: 5, px: 2 }}
    >
      <SuccessBadge>
        <SuccessIcon />
      </SuccessBadge>
      <Typography variant="h6" sx={{ fontWeight: 600 }}>
        {kindConfig?.label} connection created
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ maxWidth: 420 }}>
        <strong>{connectionName}</strong> is now a registered for use. Manage your new connection
        anytime on <ConnectionsLink href={CONNECTIONS_PATH}>Connections</ConnectionsLink>.
      </Typography>
    </Box>
  );
};

export const genericReceiptStep: WizardStep = {
  id: 'receipt',
  label: 'Done',
  Component: ReceiptStepBody,
  nextLabel: () => 'Finish',
  helpText: `Your connection is ready. Manage it anytime from [connections](${DEFAULT_CONNECTION_DOCS_URL})—change its status, assign it to environments, or remove it. [Learn more](${DEFAULT_CONNECTION_DOCS_URL}).`,
};
