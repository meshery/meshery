import dataFetch from "@/lib/data-fetch";

export const fetchKuberernetesClusters = () =>
  new Promise((res) => {
    res("Dummy");
  });

export const submitKubernetesClusterConfig = (formData) => new Promise((res, rej) => 
dataFetch(
      "/api/system/kubernetes",
      {
        credentials : "same-origin",
        method : "POST",
        body : formData,
      },
      (result) => res(result),
      (err) => rej(err)
    )
);
