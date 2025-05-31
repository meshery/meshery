import React, { useContext } from 'react';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import MesheryPerformanceComponent from './Performance';
import PatternServiceFormCore from './MesheryMeshInterface/PatternServiceFormCore';
import InfoModal from '../components/General/Modals/Information/InfoModal';
import ConfigurationSubscription from '../components/graphql/subscriptions/ConfigurationSubscription';
import _PromptComponent from './PromptComponent';
import { CapabilitiesRegistry } from '../utils/disabledComponents';
import { useNotification } from '../utils/hooks/useNotification';
import Modal from './General/Modals/Modal';
import ExportModal from './ExportModal';
import { FormatStructuredData } from './DataFormatter';
import { useFilterK8sContexts } from './hooks/useKubernetesHook';
import { useDynamicComponent } from '@/utils/context/dynamicContext';
import { ValidateDesign } from './DesignLifeCycle/ValidateDesign';
import { DryRunDesign } from './DesignLifeCycle/DryRun';
import { DeployStepper, UnDeployStepper } from './DesignLifeCycle/DeployStepper';
import { designValidationMachine } from 'machines/validator/designValidator';
import CAN from '@/utils/can';
import { mesheryEventBus } from '@/utils/eventBus';
import { ThemeTogglerCore } from '@/themes/hooks';
import RJSFForm from './MesheryMeshInterface/PatternService/RJSF';
import { DynamicFullScrrenLoader } from './LoadingComponents/DynamicFullscreenLoader';
import Troubleshoot from './TroubleshootingComponent';
import TypingFilter from './TypingFilter';
import CreateModelModal from './Settings/Registry/CreateModelModal';
import ImportModelModal from './Settings/Registry/ImportModelModal';
import { ViewInfoModal } from './ViewInfoModal';
import { selectK8sConfig, selectSelectedK8sClusters } from '@/store/slices/mesheryUi';
import { useSelector } from 'react-redux';
import { store } from '../store';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';

const requires = createRequires(getDependencies);
const useRemoteComponent = createUseRemoteComponent({ requires });

function NavigatorExtension({ url }) {
  const { capabilitiesRegistry } = useSelector((state) => state.ui);
  const { selectedK8sContexts } = useSelector((state) => state.ui);
  const [loading, err, RemoteComponent] = useRemoteComponent(url);
  const { organization: currentOrganization } = useSelector((state) => state.ui);
  const { openModalWithDefault } = useContext(WorkspaceModalContext);

  if (err != null) {
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
          <code>{err.toString()}</code>
        </div>
        <div style={{ marginTop: '1rem' }}>
          <Troubleshoot showDesignerButton={false} />
        </div>
      </div>
    );
  }
  const extensionExposedMesheryStore = {
    selectedK8sClusters: {
      get: () => selectSelectedK8sClusters(store.getState()),
    },
    k8sConfig: {
      get: () => selectK8sConfig(store.getState()),
    },
  };

  const PerformanceTestComponent = (props) => (
    <ProviderStoreWrapper>
      <MesheryPerformanceComponent {...props} />
    </ProviderStoreWrapper>
  );

  return (
    <DynamicFullScrrenLoader isLoading={loading}>
      <RemoteComponent
        injectProps={{
          PatternServiceFormCore,
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
          TypingFilter: TypingFilter,
          useNotificationHook: useNotification,
          StructuredDataFormatter: FormatStructuredData,
          CreateModelModal: CreateModelModal,
          ImportModelModal: ImportModelModal,
          ValidateDesign,
          DryRunDesign,
          DeployStepper,
          UnDeployStepper,
          designValidationMachine,
          mesheryEventBus: mesheryEventBus,
          ThemeTogglerCore,
          RJSForm: RJSFForm,
          hooks: {
            CAN: CAN,
            useFilterK8sContexts,
            useDynamicComponent,
          },
          mesheryStore: extensionExposedMesheryStore,
          currentOrganization,
          openWorkspaceModal: openModalWithDefault,
        }}
      />
    </DynamicFullScrrenLoader>
  );
}

export default NavigatorExtension;
