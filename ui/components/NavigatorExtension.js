import React from 'react';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import MesheryPerformanceComponent from './MesheryPerformance';
import PatternServiceFormCore from './MesheryMeshInterface/PatternServiceFormCore';
import InfoModal from '../components/Modals/Information/InfoModal';
import ConfigurationSubscription from '../components/graphql/subscriptions/ConfigurationSubscription';
import _PromptComponent from './PromptComponent';
import { CapabilitiesRegistry } from '../utils/disabledComponents';
import { useNotification } from '../utils/hooks/useNotification';
import Modal from './Modal';
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
import CreateModelModal from './Registry/CreateModelModal';
import ImportModelModal from './Registry/ImportModelModal';
import { ViewInfoModal } from './ViewInfoModal';

import {
  selectedOrg,
  selectK8sConfig,
  selectSelectedK8sClusters,
  setK8sContexts,
  setOrganization,
} from '@/store/slices/mesheryUi';
import { useDispatch, useSelector } from 'react-redux';
import { store } from '../store';

const requires = createRequires(getDependencies);
const useRemoteComponent = createUseRemoteComponent({ requires });

function NavigatorExtension({ url }) {
  const { k8sConfig } = useSelector((state) => state.ui);
  const { capabilitiesRegistry } = useSelector((state) => state.ui);
  const { selectedK8sContexts } = useSelector((state) => state.ui);
  const { isDrawerCollapsed } = useSelector((state) => state.ui);
  const { prometheus } = useSelector((state) => state.telemetry);
  const { grafana } = useSelector((state) => state.telemetry);
  
  const [loading, err, RemoteComponent] = useRemoteComponent(url);
  const { organization: currentOrganization } = useSelector((state) => state.ui);
  const dispatch = useDispatch();

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
    // <div>Unknown Error: {err.toString()}</div>;
  }

  const extensionExposedMesheryStore = {
    currentOrganization: {
      set: (organization) => dispatch(setOrganization({ organization })),
      get: () => selectedOrg(store.getState()),
      useCurrentOrg: () => useSelector((state) => state.ui.organization),
    },
    selectedK8sClusters: {
      get: () => selectSelectedK8sClusters(store.getState()),
      useSelectedK8sClusters: () => useSelector(selectSelectedK8sClusters),
    },
    k8sConfig: {
      get: () => selectK8sConfig(store.getState()),
      useK8sConfig: () => useSelector(selectK8sConfig),
    },
  };

  const PerformanceTestComponent = (props) => <MesheryPerformanceComponent {...props} />;

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
        }}
      />
    </DynamicFullScrrenLoader>
  );
}

export default NavigatorExtension;