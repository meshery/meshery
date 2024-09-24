import React from 'react';
import { connect } from 'react-redux';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import { bindActionCreators } from 'redux';
import { updateLoadTestData, setK8sContexts } from '../lib/store';
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
import UploadImport from './UploadImport';
import InfoModal from '../components/Modals/Information/InfoModal';
import ConfigurationSubscription from '../components/graphql/subscriptions/ConfigurationSubscription';
import PromptComponent from './PromptComponent';
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
import Troubleshoot from './TroubleshootingComponent';
import CAN from '@/utils/can';
import { mesheryEventBus } from '@/utils/can';
import { ThemeTogglerCore } from '@/themes/hooks';
import RJSFForm from './MesheryMeshInterface/PatternService/RJSF';

const requires = createRequires(getDependencies);
const useRemoteComponent = createUseRemoteComponent({ requires });
function NavigatorExtension({
  grafana,
  prometheus,
  updateLoadTestData,
  url,
  isDrawerCollapsed,
  selectedK8sContexts,
  k8sconfig,
  capabilitiesRegistry,
}) {
  const [loading, err, RemoteComponent] = useRemoteComponent(url);
  console.log(err);
  if (loading) {
    return <LoadingScreen animatedIcon="AnimatedMeshery" message="Loading Meshery Extension" />;
  }

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
    return getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sconfig);
  };

  return (
    <RemoteComponent
      injectProps={{
        GrafanaCustomCharts,
        updateLoadTestData,
        PatternServiceForm,
        RJSFWrapper,
        PatternServiceFormCore,
        grafana,
        prometheus,
        MesheryPerformanceComponent,
        dataFetch,
        createRelayEnvironment,
        subscriptionClient,
        isDrawerCollapsed,
        LoadingScreen,
        preventLeavingHook: usePreventUserFromLeavingPage,
        getSelectedK8sClusters,
        selectedK8sContexts,
        setK8sContexts,
        k8sconfig,
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
        UploadImport,
        InfoModal,
        ExportModal,
        GenericRJSFModal: Modal,
        RJSFModalWrapper: RJSFModalWrapper,
        PromptComponent,
        generateValidatePayload,
        capabilitiesRegistry,
        CapabilitiesRegistryClass: CapabilitiesRegistry,
        useNotificationHook: useNotification,
        MDEditor: MDEditor,
        StructuredDataFormatter: FormatStructuredData,
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
      }}
    />
  );
}

const mapStateToProps = (st) => {
  const grafana = st.get('grafana').toJS();
  const prometheus = st.get('prometheus').toJS();
  const isDrawerCollapsed = st.get('isDrawerCollapsed');
  const selectedK8sContexts = st.get('selectedK8sContexts');
  const k8sconfig = st.get('k8sConfig');
  const capabilitiesRegistry = st.get('capabilitiesRegistry');

  return {
    grafana,
    prometheus,
    isDrawerCollapsed,
    selectedK8sContexts,
    k8sconfig,
    capabilitiesRegistry,
  };
};

const mapDispatchToProps = (dispatch) => ({
  updateLoadTestData: bindActionCreators(updateLoadTestData, dispatch),
  setK8sContexts: bindActionCreators(setK8sContexts, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(NavigatorExtension);
