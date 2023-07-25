import React from "react";
import { createUseRemoteComponent, getDependencies, createRequires } from "@paciolan/remote-component";
// import CircularProgress from "@material-ui/core/CircularProgress";
import LoadingScreen from "./LoadingComponents/LoadingComponent";

const requires = createRequires(getDependencies);

const useRemoteComponent = createUseRemoteComponent({ requires });

const RemoteComponent = ({ url }) => {
  const [loading, err, RemoteComponent] = useRemoteComponent(url.url);
  if (loading) {
    return (
      <LoadingScreen animatedIcon="AnimatedMeshery" message="Establishing Remote Connection" />
    );
  }
  if (err != null) {
    return <div>Unknown Error: {err.toString()}</div>;
  }

  return (
    <div>
      <RemoteComponent />
    </div>
  );
}

export default RemoteComponent;
