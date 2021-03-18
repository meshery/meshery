import React from "react";
import { connect } from "react-redux";
import { createUseRemoteComponent, getDependencies, createRequires } from "@paciolan/remote-component";
import CircularProgress from "@material-ui/core/CircularProgress";
import { bindActionCreators } from "redux";
import { updateLoadTestData } from "../lib/store";
import GrafanaCustomCharts from "./GrafanaCustomCharts";
import MesheryPerformanceComponent from "./MesheryPerformance";
import dataFetch from "../lib/data-fetch"

const requires = createRequires(getDependencies);
const useRemoteComponent = createUseRemoteComponent({ requires });

function Extension({ grafana, updateLoadTestData, url }) {
  const [loading, err, RemoteComponent] = useRemoteComponent(url);

  if (loading) {
    return <CircularProgress />;
  }

  if (err != null) {
    return <div>Unknown Error: {err.toString()}</div>;
  }

  return (
    <RemoteComponent
      injectProps={{
        GrafanaCustomCharts,
        updateLoadTestData,
        grafana,
        MesheryPerformanceComponent,
        dataFetch
      }}
    />
  );
}

const mapStateToProps = (st) => {
  const grafana = st.get("grafana").toJS();
  return { grafana };
};

const mapDispatchToProps = (dispatch) => ({
  updateLoadTestData: bindActionCreators(updateLoadTestData, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(Extension);