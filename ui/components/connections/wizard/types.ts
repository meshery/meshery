import type { ReactNode, RefObject, ElementType } from 'react';
import type { ConnectionWizardKindConfig, CredentialRecord } from '../ConnectionWizard.helpers';

export type GenericRecord = Record<string, unknown>;

/**
 * `create` drives the full wizard (select -> details -> credential -> register
 * -> post-config -> receipt). `configure` mounts only the post-config + receipt
 * steps for an already-registered connection.
 */
export type WizardMode = 'create' | 'configure';

export type RjsfFormHandle = {
  validateForm: () => boolean;
  state: { formData: GenericRecord };
};

export type WizardFormRefs = {
  connection: RefObject<RjsfFormHandle | null>;
  credential: RefObject<RjsfFormHandle | null>;
};

/**
 * The mutable bag of state shared across every step. Steps read from it and
 * mutate it through `patch`; the engine owns the React state.
 */
export type WizardData = {
  // selection (create mode)
  availableKinds: ConnectionWizardKindConfig[];
  isLoadingKinds: boolean;
  connectionIconMap: Record<string, { icon?: string }>;
  kindConfig: ConnectionWizardKindConfig | null;

  // generic details + credential
  connectionFormData: GenericRecord;
  credentialMode: 'existing' | 'new';
  selectedCredentialId: string;
  credentialName: string;
  credentialFormData: GenericRecord;
  skipCredentialVerification: boolean;

  // kubernetes details
  kubeconfigFile: File | null;

  // registration outcome
  registrationId: string | null;
  registrationResult: GenericRecord | null;
  registrationError: unknown;

  // free-form per-extension post-config state
  postConfig: GenericRecord;
};

/** A kubeconfig context as returned by the discovery endpoint (not persisted). */
export type DiscoveredKubeContext = {
  id: string;
  name: string;
  server: string;
  reachable: boolean;
  connectionId?: string;
};

/** Per-context options forwarded to the kubeconfig upload endpoint. */
export type KubeconfigContextOption = {
  /** Name override for the created connection. */
  name?: string;
  /** MeshSync deployment mode for this context's connection. */
  meshsyncDeploymentMode?: string;
};

/** Options forwarded to the kubeconfig upload (create) endpoint. */
export type KubeconfigImportOptions = {
  /** Restrict the import to these discovered context IDs. */
  selectedContextIds?: string[];
  /** Per-context options (name + MeshSync mode), keyed by discovered context ID. */
  contexts?: Record<string, KubeconfigContextOption>;
};

export type WizardServices = {
  notify: (opts: { message: string; event_type: number; details?: string }) => void;
  /** POST /integrations/connections/register with the given body. */
  registerConnection: (body: GenericRecord) => Promise<GenericRecord>;
  /** POST /integrations/connections/register (status: connect). */
  connectConnection: (body: GenericRecord) => Promise<GenericRecord>;
  /** POST /system/kubernetes/contexts — parse contexts without persisting. */
  discoverKubeContexts: (file: File) => Promise<DiscoveredKubeContext[]>;
  /** POST /system/kubernetes — persist the selected contexts as connections. */
  uploadKubeconfig: (file: File, options?: KubeconfigImportOptions) => Promise<GenericRecord>;
  /** PUT /integrations/connections/{id} { status }. */
  updateConnectionById: (connectionId: string, body: GenericRecord) => Promise<GenericRecord>;
  credentials: CredentialRecord[];
};

export type WizardContext = {
  mode: WizardMode;
  data: WizardData;
  patch: (partial: Partial<WizardData>) => void;
  patchPostConfig: (partial: GenericRecord) => void;
  services: WizardServices;
  formRefs: WizardFormRefs;
};

export type WizardStep = {
  id: string;
  label: string;
  icon?: ElementType;
  Component: (props: { ctx: WizardContext }) => ReactNode;
  /** Whether the primary (Next) button is enabled. Defaults to always. */
  canProceed?: (ctx: WizardContext) => boolean;
  /**
   * Runs when leaving the step via the primary button. Return false (or throw)
   * to stay on the step. Async is awaited with a busy state on the button.
   */
  onNext?: (ctx: WizardContext) => Promise<boolean> | boolean;
  /** Override the primary button label for this step. */
  nextLabel?: (ctx: WizardContext) => string;
  /** Skip the step entirely when this returns true. */
  hidden?: (ctx: WizardContext) => boolean;
};

/**
 * A connection-specific contribution to the wizard, resolved by matching the
 * selected connection's kind (optionally narrowed by type/subType). Any step
 * left undefined falls back to the generic default; `credentialStep: null`
 * removes the credential step (e.g. Kubernetes carries its credential inline).
 */
export type ConnectionExtension = {
  match: { kind: string; type?: string; subType?: string };
  detailsStep?: WizardStep;
  credentialStep?: WizardStep | null;
  registerStep?: WizardStep;
  postConfigSteps?: WizardStep[];
  receiptStep?: WizardStep;
};
