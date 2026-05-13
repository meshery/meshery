import React, { useContext, useMemo } from 'react';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import MesheryPerformanceComponent from '../../performance';
import PatternServiceFormCore from '../../meshery-mesh-interface/PatternServiceFormCore';
import InfoModal from '../../shared/Modal/Information/InfoModal';
import ConfigurationSubscription from '@/graphql/subscriptions/ConfigurationSubscription';
import _PromptComponent from '../../PromptComponent';
import { CapabilitiesRegistry } from '../../../utils/disabledComponents';
import { useNotification } from '../../../utils/hooks/useNotification';
import Modal from '../../shared/Modal/Modal';
import ExportModal from '../../shared/Modal/ExportModal';
import { FormatStructuredData } from '../../data-formatter';
import { useFilterK8sContexts } from '@/utils/hooks/useKubernetesHook';
import { useDynamicComponent } from '@/utils/context/dynamicContext';
import { ValidateDesign } from '../../designs/lifecycle/ValidateDesign';
import { DryRunDesign } from '../../designs/lifecycle/DryRun';
import { DeployStepper, UnDeployStepper } from '../../designs/lifecycle/DeployStepper';
import { designValidationMachine } from 'machines/validator/designValidator';
import CAN from '@/utils/can';
import { mesheryEventBus } from '@/utils/eventBus';
import { ThemeTogglerCore } from '@/themes/hooks';
import RJSFForm from '../../meshery-mesh-interface/PatternService/RJSF';
import { DynamicFullScreenLoader } from '../../shared/LoadingState/DynamicFullscreenLoader';
import Troubleshoot from '../../TroubleshootingComponent';
import TypingFilter from '../../typing-filter';
import CreateModelModal from '../../registry/CreateModelModal';
import ImportModelModal from '../../registry/ImportModelModal';
import { ViewInfoModal } from '../../shared/Modal/ViewInfoModal';
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
    capabilitiesRegistry,
    selectedK8sContexts,
    organization: currentOrganization,
  } = useSelector((state) => state.ui);
  const [loading, err, RemoteComponent] = useRemoteComponent(url);
  const { openModalWithDefault, onLoadResource } = useContext(WorkspaceModalContext);
  const registryModal = useRegistryModal();

  const injectProps = useMemo(
    () => ({
      PatternServiceFormCore,
      RelationshipEvaluationResponseFormatter: RelationshipEvaluationTraceFormatter,
      MesheryPerformanceComponent: PerformanceTestComponent,
      selectedK8sContexts,
      resolver: {
        query: {},
        mutation: {},
        subscription: {
          ConfigurationSubscription,
        },
      },
      InfoModal,
      ViewInfoModal,
      ExportModal,
      GenericRJSFModal: Modal,
      _PromptComponent,
      capabilitiesRegistry,
      CapabilitiesRegistryClass: CapabilitiesRegistry,
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
      capabilitiesRegistry,
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
