import dataFetch from "../../../lib/data-fetch";

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

