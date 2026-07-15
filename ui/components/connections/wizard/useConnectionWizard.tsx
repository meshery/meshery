import { createRef, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useNotification } from '@/utils/hooks/useNotification';
import {
  useAddKubernetesConfigMutation,
  useConnectToConnectionMutation,
  useDiscoverKubernetesContextsMutation,
  useGetCredentialsQuery,
  usePerformConnectionActionMutation,
  useUpdateConnectionByIdMutation,
  useVerifyAndRegisterConnectionMutation,
} from '@/rtk-query/connection';
import type { ConnectionWizardKindConfig } from '../ConnectionWizard.helpers';
import { buildSteps } from './registry';
import type {
  GenericRecord,
  WizardContext,
  WizardData,
  WizardFormRefs,
  WizardMode,
  WizardServices,
} from './types';

export type UseConnectionWizardParams = {
  mode: WizardMode;
  isOpen: boolean;
  availableKinds?: ConnectionWizardKindConfig[];
  isLoadingKinds?: boolean;
  connectionIconMap?: Record<string, { icon?: string }>;
  /** configure mode: the connection being (re)configured. */
  initialKindConfig?: ConnectionWizardKindConfig | null;
  initialRegistrationResult?: GenericRecord | null;
  onComplete?: () => void;
};

const makeInitialData = (params: UseConnectionWizardParams): WizardData => ({
  availableKinds: params.availableKinds ?? [],
  isLoadingKinds: params.isLoadingKinds ?? false,
  connectionIconMap: params.connectionIconMap ?? {},
  kindConfig: params.initialKindConfig ?? null,
  connectionFormData: {},
  credentialMode: 'existing',
  selectedCredentialId: '',
  credentialName: '',
  credentialFormData: {},
  skipCredentialVerification: false,
  kubeconfigFile: null,
  registrationId: null,
  registrationResult: params.initialRegistrationResult ?? null,
  registrationError: null,
  postConfig: {},
});

export const useConnectionWizard = (params: UseConnectionWizardParams) => {
  const { mode, isOpen, onComplete } = params;
  const { notify } = useNotification();
  const [verifyAndRegisterConnection] = useVerifyAndRegisterConnectionMutation();
  const [connectToConnection] = useConnectToConnectionMutation();
  const [addKubernetesConfig] = useAddKubernetesConfigMutation();
  const [discoverKubernetesContexts] = useDiscoverKubernetesContextsMutation();
  const [updateConnectionById] = useUpdateConnectionByIdMutation();
  const [performConnectionAction] = usePerformConnectionActionMutation();
  const { data: credentialsResponse } = useGetCredentialsQuery(undefined, { skip: !isOpen });

  const formRefs = useRef<WizardFormRefs>({
    connection: createRef(),
    credential: createRef(),
  }).current;

  const [data, setData] = useState<WizardData>(() => makeInitialData(params));
  const [activeIndex, setActiveIndex] = useState(0);
  const [isBusy, setIsBusy] = useState(false);

  const patch = useCallback((partial: Partial<WizardData>) => {
    setData((current) => ({ ...current, ...partial }));
  }, []);
  const patchPostConfig = useCallback((partial: GenericRecord) => {
    setData((current) => ({ ...current, postConfig: { ...current.postConfig, ...partial } }));
  }, []);

  // Keep the externally-owned inputs (kind list, icons) in sync with the store.
  const { availableKinds, isLoadingKinds, connectionIconMap } = params;
  useEffect(() => {
    setData((current) => ({
      ...current,
      availableKinds: availableKinds ?? [],
      isLoadingKinds: isLoadingKinds ?? false,
      connectionIconMap: connectionIconMap ?? {},
    }));
  }, [availableKinds, isLoadingKinds, connectionIconMap]);

  // In configure mode the connection (kind + details) is resolved
  // asynchronously by the caller; seed it once it arrives.
  const { initialKindConfig, initialRegistrationResult } = params;
  useEffect(() => {
    if (mode !== 'configure') {
      return;
    }
    setData((current) => ({
      ...current,
      kindConfig: initialKindConfig ?? null,
      registrationResult: initialRegistrationResult ?? current.registrationResult,
    }));
  }, [mode, initialKindConfig, initialRegistrationResult]);

  // Keep `reset` stable while still resetting from the latest params: read them
  // from a ref updated in the render body (not an effect, so children never see
  // stale values during commit).
  const paramsRef = useRef(params);
  paramsRef.current = params;
  const reset = useCallback(() => {
    setData(makeInitialData(paramsRef.current));
    setActiveIndex(0);
    setIsBusy(false);
  }, []);

  const services = useMemo<WizardServices>(
    () => ({
      notify,
      registerConnection: (body) => verifyAndRegisterConnection({ body }).unwrap(),
      connectConnection: (body) => connectToConnection({ body }).unwrap(),
      discoverKubeContexts: async (file) => {
        const formData = new FormData();
        formData.append('k8sfile', file);
        const result = await discoverKubernetesContexts({ body: formData }).unwrap();
        return Array.isArray(result) ? result : [];
      },
      uploadKubeconfig: (file, options) => {
        const formData = new FormData();
        formData.append('k8sfile', file);
        if (options?.selectedContextIds) {
          formData.append('selectedContexts', JSON.stringify(options.selectedContextIds));
        }
        // The backend keys per-context options (name + meshsyncDeploymentMode)
        // under `contexts` by context id, persisting the mode to each
        // connection's metadata before it connects.
        if (options?.contexts && Object.keys(options.contexts).length > 0) {
          formData.append('contexts', JSON.stringify(options.contexts));
        }
        return addKubernetesConfig({ body: formData }).unwrap();
      },
      updateConnectionById: (connectionId, body) =>
        updateConnectionById({ connectionId, body }).unwrap(),
      setMeshsyncMode: (connectionId, mode) =>
        performConnectionAction({
          connectionId,
          body: { action: 'setMeshsyncMode', mode },
        }).unwrap(),
      flushMeshsync: (connectionId) =>
        performConnectionAction({
          connectionId,
          // `flushMeshsync` is a valid server action but is not yet part of the
          // published schemas action enum (a closed union of `setMeshsyncMode`),
          // so cast until the schemas enum change ships. FOLLOW-UP: drop the cast
          // once @meshery/schemas exposes the `flushMeshsync` action.
          body: { action: 'flushMeshsync' } as unknown as { action: 'setMeshsyncMode' },
        }).unwrap(),
      credentials: credentialsResponse?.credentials || [],
    }),
    [
      notify,
      verifyAndRegisterConnection,
      connectToConnection,
      addKubernetesConfig,
      discoverKubernetesContexts,
      updateConnectionById,
      performConnectionAction,
      credentialsResponse,
    ],
  );

  const ctx: WizardContext = { mode, data, patch, patchPostConfig, services, formRefs };

  const steps = buildSteps(data.kindConfig, mode).filter((step) => !step.hidden?.(ctx));
  const safeIndex = steps.length === 0 ? 0 : Math.min(activeIndex, steps.length - 1);
  const activeStep = steps[safeIndex] ?? null;
  const isLast = safeIndex >= steps.length - 1;

  const canProceed = activeStep?.canProceed ? activeStep.canProceed(ctx) : true;
  const canGoBack = safeIndex > 0;
  const nextLabel = activeStep?.nextLabel ? activeStep.nextLabel(ctx) : isLast ? 'Finish' : 'Next';

  const next = async () => {
    if (!activeStep || isBusy) {
      return;
    }
    if (activeStep.onNext) {
      setIsBusy(true);
      let ok = false;
      try {
        ok = await activeStep.onNext(ctx);
      } catch {
        ok = false;
      }
      setIsBusy(false);
      if (!ok) {
        return;
      }
    }
    if (isLast) {
      onComplete?.();
      return;
    }
    setActiveIndex(safeIndex + 1);
  };

  const back = () => setActiveIndex((index) => Math.max(0, index - 1));

  return {
    ctx,
    steps,
    stepLabels: steps.map((step) => step.label),
    activeIndex: safeIndex,
    activeStep,
    isBusy,
    isLast,
    canProceed,
    canGoBack,
    nextLabel,
    next,
    back,
    reset,
  };
};
