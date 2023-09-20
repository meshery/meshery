import React from "react";
import { createUseRemoteComponent, getDependencies, createRequires } from "@paciolan/remote-component";
// import CircularProgress from "@material-ui/core/CircularProgress";
import LoadingScreen from "./LoadingComponents/LoadingComponent";

const requires = createRequires(getDependencies);

const useRemoteComponent = createUseRemoteComponent({ requires });

const RemoteComponent = ({ url, loaderType }) => {
  const [loading, err, RemoteComponent] = useRemoteComponent(url.url);
  if (loading) {
    if (loaderType === "circular") {
      return ''
    } else {
      return <LoadingScreen animatedIcon="AnimatedMeshery" message="Establishing Remote Connection" />;
    }
  }
  if (err != null) {
    /* Debugging log */
    console.error(`Unknown Error: ${err.toString()}`)
    return <></>;
  }

  return (
    <div>
      <RemoteComponent />
    </div>
  );
}

export default RemoteComponent;
