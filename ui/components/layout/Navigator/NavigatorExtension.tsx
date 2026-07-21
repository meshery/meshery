import React, { useContext, useMemo } from 'react';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import MesheryPerformanceComponent from '../../performance';
import PatternServiceFormCore from '../../meshery-mesh-interface/PatternServiceFormCore';
import InfoModal from '../../shared/Modal/Information/InfoModal';
import _PromptComponent from '../../general/PromptComponent';
import { ProviderUiAccessControl } from '../../../utils/disabledComponents';
import { useNotification } from '../../../utils/hooks/useNotification';
import Modal from '../../shared/Modal/Modal';
import ExportDesignModal from '../../designs/export/ExportDesignModal';
import { FormatStructuredData } from '../../data-formatter';
import { useFilterK8sContexts } from '@/utils/hooks/useKubernetesHook';
import { useDynamicComponent } from '@/utils/context/dynamicContext';
import { ValidateDesign } from '../../designs/lifecycle/ValidateDesign';
import { DryRunDesign } from '../../designs/lifecycle/DryRun';
import { DeployStepper, UnDeployStepper } from '../../designs/lifecycle/DeployStepper';
import { designValidationMachine } from 'machines/validator/designValidator';
import CAN from '@/utils/can';
import { mesheryEventBus } from '@/utils/eventBus';
import { ThemeTogglerCore } from '@/theme/hooks';
import RJSFForm from '../../meshery-mesh-interface/PatternService/RJSF';
import { DynamicFullScreenLoader } from '../../shared/LoadingState/DynamicFullscreenLoader';
import Troubleshoot from '../../TroubleshootingComponent';
import TypingFilter from '@/components/shared/FormFields/typing-filter';
import CreateModelModal from '../../registry/CreateModelModal';
import ImportModelModal from '../../registry/ImportModelModal';
import { ViewInfoModal } from '../../workspaces/ViewInfoModal';
import { selectK8sConfig, selectSelectedK8sClusters } from '@/store/slices/mesheryUi';
import { useSelector } from 'react-redux';
import { store } from '../../../store';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import { RelationshipEvaluationTraceFormatter } from '../NotificationCenter/formatters/relationship_evaluation';
import { useRegistryModal } from '@/utils/hooks/useRegistryModal';

const requires = createRequires(getDependencies);
const useRemoteComponent = createUseRemoteComponent({ requires });

type NavigatorExtensionProps = {
  url: string;
};

const extensionExposedMesheryStore = {
  selectedK8sClusters: {
    get: () => selectSelectedK8sClusters(store.getState()),
  },
  k8sConfig: {
    get: () => selectK8sConfig(store.getState()),
  },
};

/**
 * Inert stand-in for the removed `subscribeConfiguration` GraphQL subscription.
 * Matches the old call shape — `(onNext, variables) => ({ dispose })` — so an
 * extension bundle published against the previous host contract still mounts and
 * unmounts cleanly. It never emits: `onNext` is intentionally never invoked.
 */
const noopSubscription = () => ({ dispose: () => {} });

function PerformanceTestComponent(props: React.ComponentProps<typeof MesheryPerformanceComponent>) {
  return (
    <ProviderStoreWrapper>
      <MesheryPerformanceComponent {...props} />
    </ProviderStoreWrapper>
  );
}

function NavigatorExtensionError({ error }: { error: unknown }) {
  return (
    <div role="alert">
      <h2>Uh-oh!😔 Please pardon our mesh.</h2>
      <div
        style={{
          backgroundColor: '#1E2117',
          color: '#FFFFFF',
          padding: '.85rem',
          borderRadius: '.2rem',
        }}
      >
        <code>{String(error)}</code>
      </div>
      <div style={{ marginTop: '1rem' }}>
        <Troubleshoot showDesignerButton={false} />
      </div>
    </div>
  );
}

function NavigatorExtension({ url }: NavigatorExtensionProps) {
  const {
    providerCapabilities,
    selectedK8sContexts,
    organization: currentOrganization,
  } = useSelector((state) => state.ui);
  // `useRemoteComponent` fetches the bundle, evaluates it as a CommonJS module,
  // and returns `module.exports.default` as the component. The remote bundle MUST
  // therefore export the component as a CommonJS default
  // (`module.exports = { default: Component, __esModule: true }`). NOTE: the hook
  // does NOT throw when `.default` is missing — it returns `RemoteComponent ===
  // undefined` with `err === undefined`, so a mis-built bundle fails silently
  // (no loader error; React then throws "Element type is invalid … got: undefined").
  // If an extension renders as undefined with no error, inspect the bundle's export
  // shape, not this code. See meshery-extensions/docs/troubleshoot.md.
  const [loading, err, RemoteComponent] = useRemoteComponent(url);
  const { openModalWithDefault, onLoadResource } = useContext(WorkspaceModalContext);
  const registryModal = useRegistryModal();

  const injectProps = useMemo(
    () => ({
      PatternServiceFormCore,
      RelationshipEvaluationResponseFormatter: RelationshipEvaluationTraceFormatter,
      MesheryPerformanceComponent: PerformanceTestComponent,
      selectedK8sContexts,
      // Meshery Server no longer exposes any GraphQL subscription. `resolver` is
      // kept only for backward compatibility with already-published extension
      // bundles: they call `resolver.subscription.ConfigurationSubscription(cb, vars)`
      // and later `.dispose()` on the result. Handing them `undefined` would throw
      // at mount, so we hand them an inert subscription that never emits. New
      // extensions should read designs/filters from the REST API instead.
      resolver: {
        query: {},
        mutation: {},
        subscription: {
          ConfigurationSubscription: noopSubscription,
        },
      },
      InfoModal,
      ViewInfoModal,
      ExportModal: ExportDesignModal,
      GenericRJSFModal: Modal,
      _PromptComponent,
      providerCapabilities,
      ProviderUiAccessControlClass: ProviderUiAccessControl,
      // Backward-compatible alias for already-published remote extensions.
      CapabilitiesRegistryClass: ProviderUiAccessControl,
      TypingFilter,
      useNotificationHook: useNotification,
      StructuredDataFormatter: FormatStructuredData,
      CreateModelModal,
      ImportModelModal,
      ValidateDesign,
      DryRunDesign,
      DeployStepper,
      UnDeployStepper,
      designValidationMachine,
      mesheryEventBus,
      ThemeTogglerCore,
      RJSForm: RJSFForm,
      hooks: {
        CAN,
        useFilterK8sContexts,
        useDynamicComponent,
      },
      mesheryStore: extensionExposedMesheryStore,
      currentOrganization,
      openWorkspaceModal: openModalWithDefault,
      openRegistryModal: registryModal,
      SetCurrentLoadedResourceInOrgWorkspaceSession: onLoadResource,
    }),
    [
      providerCapabilities,
      currentOrganization,
      onLoadResource,
      openModalWithDefault,
      registryModal,
      selectedK8sContexts,
    ],
  );

  if (err != null) {
    return <NavigatorExtensionError error={err} />;
  }

  return (
    <DynamicFullScreenLoader isLoading={loading}>
      <RemoteComponent injectProps={injectProps} />
    </DynamicFullScreenLoader>
  );
}

export default NavigatorExtension;
