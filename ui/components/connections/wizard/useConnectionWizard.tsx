import { createRef, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useNotification } from '@/utils/hooks/useNotification';
import {
  useAddKubernetesConfigMutation,
  useConnectToConnectionMutation,
  useGetCredentialsQuery,
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
  credentialFormData: {},
  skipCredentialVerification: false,
  kubeconfigFile: null,
  registrationId: null,
  connectionModel: null,
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
  const [updateConnectionById] = useUpdateConnectionByIdMutation();
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

  const reset = useCallback(() => {
    setData(makeInitialData(params));
    setActiveIndex(0);
    setIsBusy(false);
    // params identity is stable enough for a reset trigger; intentionally not a dep.
  }, []);

  const services = useMemo<WizardServices>(
    () => ({
      notify,
      registerConnection: (body) => verifyAndRegisterConnection({ body }).unwrap(),
      connectConnection: (body) => connectToConnection({ body }).unwrap(),
      uploadKubeconfig: (file) => {
        const formData = new FormData();
        formData.append('k8sfile', file);
        return addKubernetesConfig({ body: formData }).unwrap();
      },
      updateConnectionById: (connectionId, body) =>
        updateConnectionById({ connectionId, body }).unwrap(),
      credentials: credentialsResponse?.credentials || [],
    }),
    [
      notify,
      verifyAndRegisterConnection,
      connectToConnection,
      addKubernetesConfig,
      updateConnectionById,
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
