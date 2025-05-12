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
  useLegacySelector,
  LegacyStoreContext,
  actionTypes,
} from '../lib/store';
import GrafanaCustomCharts from './telemetry/grafana/GrafanaCustomCharts';
import MesheryPerformanceComponent from './MesheryPerformance';
import dataFetch from '../lib/data-fetch';
import PatternServiceForm from './MesheryMeshInterface/PatternServiceForm';
import PatternServiceFormCore from './MesheryMeshInterface/PatternServiceFormCore';
import RJSFWrapper from './MesheryMeshInterface/PatternService/RJSF_wrapper';
import { createRelayEnvironment, subscriptionClient } from '../lib/relayEnvironment';
import LoadingScreen from './LoadingComponents/LoadingComponent';
import usePreventUserFromLeavingPage from '../utils/hooks/usePreventUserFromLeavingPage';
import { getK8sClusterIdsFromCtxId } from '../utils/multi-ctx';
import ConfirmationModal, { SelectDeploymentTarget } from './ConfirmationModal';
import { getComponentsinFile, generateValidatePayload } from '../utils/utils';
import InfoModal from '../components/Modals/Information/InfoModal';
import ConfigurationSubscription from '../components/graphql/subscriptions/ConfigurationSubscription';
import _PromptComponent from './PromptComponent';
import { CapabilitiesRegistry } from '../utils/disabledComponents';
import { useNotification } from '../utils/hooks/useNotification';
import Modal, { RJSFModalWrapper } from './Modal';
import ExportModal from './ExportModal';
import { MDEditor } from './Markdown';
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
import { RTKContext, useSelectorRtk } from '@/store/hooks';
import {
  selectK8sConfig,
  selectSelectedK8sClusters,
  setK8sContexts,
} from '@/store/slices/mesheryUi';

const requires = createRequires(getDependencies);
const useRemoteComponent = createUseRemoteComponent({ requires });
function NavigatorExtension({ grafana, prometheus, updateLoadTestData, url, isDrawerCollapsed }) {
  const { k8sConfig } = useSelectorRtk((state) => state.ui);
  const { capabilitiesRegistry } = useSelectorRtk((state) => state.ui);
  const { selectedK8sContexts } = useSelectorRtk((state) => state.ui);
  const [loading, err, RemoteComponent] = useRemoteComponent(url);
  const currentOrganization = useGetCurrentOrganization();
  const { store: legacyStore } = useContext(LegacyStoreContext);
  const { store: rtkStore } = useContext(RTKContext);

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

  const getSelectedK8sClusters = () => {
    return getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sConfig);
  };

  const extensionExposedMesheryStore = {
    currentOrganization: {
      set: (organization) =>
        legacyStore.dispatch({ type: actionTypes.SET_ORGANIZATION, organization }),
      get: () => legacyStore.getState().organization,
      useCurrentOrg: () => useLegacySelector((state) => state.organization),
    },
    selectedK8sClusters: {
      get: () => selectSelectedK8sClusters(rtkStore.getState()),
      useSelectedK8sClusters: () => useSelectorRtk(selectSelectedK8sClusters),
    },
    k8sConfig: {
      get: () => selectK8sConfig(rtkStore.getState()),
      useK8sConfig: () => useSelectorRtk(selectK8sConfig),
    },
  };

  const PerformanceTestComponent = (props) => <MesheryPerformanceComponent {...props} />;

  return (
    <DynamicFullScrrenLoader isLoading={loading}>
      <RemoteComponent
        injectProps={{
          GrafanaCustomCharts,
          updateLoadTestData,
          PatternServiceForm,
          RJSFWrapper,
          PatternServiceFormCore,
          grafana,
          prometheus,
          MesheryPerformanceComponent: PerformanceTestComponent,
          dataFetch,
          createRelayEnvironment,
          subscriptionClient,
          isDrawerCollapsed,
          LoadingScreen,
          preventLeavingHook: usePreventUserFromLeavingPage,
          getSelectedK8sClusters,
          selectedK8sContexts,
          setK8sContexts,
          k8sconfig: k8sConfig,
          resolver: {
            query: {},
            mutation: {},
            subscription: {
              ConfigurationSubscription,
            },
          },
          ConfirmationModal,
          SelectDeploymentTarget: SelectDeploymentTarget,
          getComponentsinFile,
          InfoModal,
          ViewInfoModal,
          ExportModal,
          GenericRJSFModal: Modal,
          RJSFModalWrapper: RJSFModalWrapper,
          _PromptComponent,
          generateValidatePayload,
          capabilitiesRegistry,
          CapabilitiesRegistryClass: CapabilitiesRegistry,
          TypingFilter: TypingFilter,
          useNotificationHook: useNotification,
          MDEditor: MDEditor,
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
  const grafana = st.get('grafana').toJS();
  const prometheus = st.get('prometheus').toJS();
  const isDrawerCollapsed = st.get('isDrawerCollapsed');

  return {
    grafana,
    prometheus,
    isDrawerCollapsed,
  };
};

const mapDispatchToProps = (dispatch) => ({
  updateLoadTestData: bindActionCreators(updateLoadTestData, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(NavigatorExtension);
