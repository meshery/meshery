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
