import React from "react";
import { connect } from "react-redux";
import { createUseRemoteComponent, getDependencies, createRequires } from "@paciolan/remote-component";
import { bindActionCreators } from "redux";
import { updateLoadTestData, setK8sContexts } from "../lib/store";
import GrafanaCustomCharts from "./GrafanaCustomCharts";
import MesheryPerformanceComponent from "./MesheryPerformance";
import dataFetch from "../lib/data-fetch";
import PatternServiceForm from "./MesheryMeshInterface/PatternServiceForm";
import PatternServiceFormCore from "./MesheryMeshInterface/PatternServiceFormCore";
import environment, { subscriptionClient } from "../lib/relayEnvironment";
import subscribeMeshSyncStatusEvents from "../components/graphql/subscriptions/MeshSyncStatusSubscription"
import LoadingScreen from "./LoadingComponents/LoadingComponent";
import usePreventUserFromLeavingPage from "../utils/hooks/usePreventUserFromLeavingPage";
import { getK8sClusterIdsFromCtxId } from "../utils/multi-ctx";
import ConfirmationModal from "./ConfirmationModal"
import { getComponentsinFile } from "../utils/utils";
import UploadImport from "./UploadImport";
import ConfigurationSubscription from "../components/graphql/subscriptions/ConfigurationSubscription";

const requires = createRequires(getDependencies);
const useRemoteComponent = createUseRemoteComponent({ requires });

function Extension({ grafana, prometheus, updateLoadTestData, url, isDrawerCollapsed, selectedK8sContexts, k8sconfig }) {
  const [loading, err, RemoteComponent] = useRemoteComponent(url);

  if (loading) {
    return <LoadingScreen message="Loading Meshery Extension" />;
  }

  if (err != null) {
    return <div>Unknown Error: {err.toString()}</div>;
  }

  const getSelectedK8sClusters = () => {
    return getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sconfig);
  }

  return (
    <RemoteComponent
      injectProps={{
        GrafanaCustomCharts,
        updateLoadTestData,
        PatternServiceForm,
        PatternServiceFormCore,
        grafana,
        prometheus,
        MesheryPerformanceComponent,
        dataFetch,
        environment,
        subscriptionClient,
        isDrawerCollapsed,
        LoadingScreen,
        preventLeavingHook : usePreventUserFromLeavingPage,
        getSelectedK8sClusters,
        selectedK8sContexts,
        setK8sContexts,
        k8sconfig,
        resolver : {
          query : {},
          mutation : {},
          subscription : {
            subscribeMeshSyncStatusEvents,
            ConfigurationSubscription
          },
        },
        ConfirmationModal,
        getComponentsinFile,
        UploadImport
      }}
    />
  );
}

const mapStateToProps = (st) => {
  const grafana = st.get("grafana").toJS();
  const prometheus = st.get("prometheus").toJS();
  const isDrawerCollapsed = st.get("isDrawerCollapsed");
  const selectedK8sContexts = st.get('selectedK8sContexts');
  const k8sconfig = st.get("k8sConfig");

  return { grafana, prometheus, isDrawerCollapsed, selectedK8sContexts, k8sconfig };
};

const mapDispatchToProps = (dispatch) => ({ updateLoadTestData : bindActionCreators(updateLoadTestData, dispatch),
  setK8sContexts : bindActionCreators(setK8sContexts, dispatch) }
);

export default connect(mapStateToProps, mapDispatchToProps)(Extension);
