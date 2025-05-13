import React, { useContext } from 'react';
import { connect } from 'react-redux';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import { bindActionCreators } from 'redux';
import {
  updateLoadTestData,
  setK8sContexts,
  useLegacySelector,
  LegacyStoreContext,
  actionTypes,
  selectSelectedK8sClusters,
  selectK8sConfig,
} from '../lib/store';
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
import { useGetCurrentOrganization } from '@/utils/hooks/useStateValue';

const requires = createRequires(getDependencies);
const useRemoteComponent = createUseRemoteComponent({ requires });
function NavigatorExtension({ url, selectedK8sContexts, capabilitiesRegistry }) {
  const [loading, err, RemoteComponent] = useRemoteComponent(url);
  const currentOrganization = useGetCurrentOrganization();
  const { store: legacyStore } = useContext(LegacyStoreContext);
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
      set: (organization) =>
        legacyStore.dispatch({ type: actionTypes.SET_ORGANIZATION, organization }),
      get: () => legacyStore.getState().organization,
      useCurrentOrg: () => useLegacySelector((state) => state.organization),
    },
    selectedK8sClusters: {
      get: () => selectSelectedK8sClusters(legacyStore.getState()),
      useSelectedK8sClusters: () => useLegacySelector(selectSelectedK8sClusters),
    },
    k8sconfig: {
      get: () => selectK8sConfig(legacyStore.getState()),
      useK8sConfig: () => useLegacySelector(selectK8sConfig),
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

const mapStateToProps = (st) => {
  const selectedK8sContexts = st.get('selectedK8sContexts');
  const capabilitiesRegistry = st.get('capabilitiesRegistry');

  return {
    selectedK8sContexts,
    capabilitiesRegistry,
  };
};

const mapDispatchToProps = (dispatch) => ({
  updateLoadTestData: bindActionCreators(updateLoadTestData, dispatch),
  setK8sContexts: bindActionCreators(setK8sContexts, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(NavigatorExtension);
