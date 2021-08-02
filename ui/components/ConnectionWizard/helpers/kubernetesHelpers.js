import dataFetch from "../../../lib/data-fetch";


/** 
  * Pings kuberenetes server endpoint
  * @param  {(res) => void} successHandler
  * @param  {(err) => void} errorHandler
*/
export const pingKubernetes = (successHandler,errorHandler) => {
  dataFetch(
    "/api/k8sconfig/ping",
    {
      credentials: "same-origin",
      credentials: "include",
    },
    successHandler,
    errorHandler
  );
}

export const pingKubernetesWithNotification = (updateProgress, snackbar, action) => {

  const successHandler = (res) => {
    updateProgress({showProgress: false})  
    if(res !== undefined)
      snackbar(`Kubernetes sucessfully pinged`, {
        variant: "success",
        action,
        autoHideDuration: 7000,
      })
  }

  const errorHandler = (err) => {
    updateProgress({showProgress: false})  
    snackbar(`Failed to ping kubernetes: ${err}`, {
      variant: "error",
      action,
      autoHideDuration: 7000,
    })
  }

  updateProgress({showProgress : true})
  pingKubernetes(successHandler, errorHandler)

}
  
/**
  * Figures out if kubernetes connection is established or not
  *
  * @param {true|false} isClusterConfigured - data received from meshery server
  * as to whether or not the server config is found
  * @param {true|false} kubernetesPingStatus - found after pinging the kubernetes
  * server endpoint
  *
  * @return {true|false}
*/
export const isKubernetesConnected = (isClusterConfigured,kubernetesPingStatus) => {

  if(isClusterConfigured){
    if(kubernetesPingStatus) return true
  }

  return false
}


export const deleteKubernetesConfig = (successCb,errorCb) => 
  dataFetch(
    "/api/k8sconfig",
    {
      credentials: "same-origin",
      method: "DELETE",
      credentials: "include",
    },
    successCb,
    errorCb
  )

export const reconfigureKubernetes = (updateProgress, enqueueSnackbar, action, updateK8SConfig) => {

  const successCb = (result) => {
    updateProgress({ showProgress: false });
    if (typeof result !== "undefined") {
      
      updateK8SConfig({
        k8sConfig: {
          inClusterConfig: false,
          k8sfile: "",
          contextName: "",
          clusterConfigured: false,
        },
      });
      enqueueSnackbar("Kubernetes config was successfully removed!", {
        variant: "success",
        autoHideDuration: 2000,
        action
      });
    }
  }
      
  const errorCb = (err) => {

    updateProgress({ showProgress: false });
    enqueueSnackbar("Not able to remove kubernetes cluster: "+err, {
      variant: "error",
      autoHideDuration: 2000,
      action 
    });
  }

  deleteKubernetesConfig(
    successCb,
    errorCb
  )

}


export const fetchContexts = (updateProgress, k8sfile) => {

  const formData = new FormData();
  // if (inClusterConfigForm) {
  //   return;
  // }
    
  // formData.append('contextName', contextName);
  formData.append("k8sfile", k8sfile);

  updateProgress({ showProgress: true });

  return new Promise((res, rej) => {
    dataFetch(
      "/api/k8sconfig/contexts",
      {
        credentials: "same-origin",
        method: "POST",
        credentials: "include",
        body: formData,
      },
      (result) => {
        updateProgress({ showProgress: false });

        if (typeof result !== "undefined") {
          let ctName = "";
          result.forEach(({ contextName, currentContext }) => {
            if (currentContext) {
              ctName = contextName;
            }
          });

          res({result, currentContextName: ctName})
        }
      },
      (err) => rej(err)
    );
  })
    
};


export const submitConfig = (enqueueSnackbar, updateProgress, updateK8SConfig, action, contextName, k8sfile) => {

  const inClusterConfigForm = false
  const formData = new FormData();
  formData.append("inClusterConfig", inClusterConfigForm ? "on" : ""); // to simulate form behaviour of a checkbox
  if (!inClusterConfigForm) {
    formData.append("contextName", contextName);
    formData.append("k8sfile", k8sfile);
  }
  updateProgress({ showProgress: true });
  dataFetch(
    "/api/k8sconfig",
    {
      credentials: "same-origin",
      method: "POST",
      credentials: "include",
      body: formData,
    },
    (result) => {
      updateProgress({ showProgress: false });
      if (typeof result !== "undefined") {
        enqueueSnackbar("Kubernetes config was successfully validated!", {
          variant: "success",
          autoHideDuration: 2000,
          action
        });
        updateK8SConfig({
          k8sConfig: {
            inClusterConfig: inClusterConfigForm,
            k8sfile,
            contextName: result.contextName,
            clusterConfigured: true,
            configuredServer: result.configuredServer,
          },
        });
      }
    },
    (err) => alert(err)
  );
};
