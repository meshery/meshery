import React from "react";
import { connect } from "react-redux";
import { createUseRemoteComponent, getDependencies, createRequires } from "@paciolan/remote-component";
import { bindActionCreators } from "redux";
import { updateLoadTestData } from "../lib/store";
import GrafanaCustomCharts from "./GrafanaCustomCharts";
import MesheryPerformanceComponent from "./MesheryPerformance";
import dataFetch from "../lib/data-fetch";
import PatternServiceForm from "./MesheryMeshInterface/PatternServiceForm";
import PatternServiceFormCore from "./MesheryMeshInterface/PatternServiceFormCore";
import environment, { subscriptionClient } from "../lib/relayEnvironment";
import subscribeMeshSyncStatusEvents from "../components/graphql/subscriptions/MeshSyncStatusSubscription"
import LoadingScreen from "./LoadingComponents/LoadingComponent";
import usePreventUserFromLeavingPage from "../utils/hooks/usePreventUserFromLeavingPage";

const requires = createRequires(getDependencies);
const useRemoteComponent = createUseRemoteComponent({ requires });

function Extension({ grafana, prometheus, updateLoadTestData, url, isDrawerCollapsed }) {
  const [loading, err, RemoteComponent] = useRemoteComponent(url);

  if (loading) {
    return <LoadingScreen message="Loading Meshery Extension" />;
  }

  if (err != null) {
    return <div>Unknown Error: {err.toString()}</div>;
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
        resolver : {
          query : {},
          mutation : {},
          subscription : {
            subscribeMeshSyncStatusEvents,
          },
        },
      }}
    />
  );
}

const mapStateToProps = (st) => {
  const grafana = st.get("grafana").toJS();
  const prometheus = st.get("prometheus").toJS();
  const isDrawerCollapsed = st.get("isDrawerCollapsed");
  return { grafana, prometheus, isDrawerCollapsed };
};

const mapDispatchToProps = (dispatch) => ({ updateLoadTestData : bindActionCreators(updateLoadTestData, dispatch) });

export default connect(mapStateToProps, mapDispatchToProps)(Extension);
