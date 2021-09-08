import dataFetch from "@/lib/data-fetch";

/**
 * @typedef {{k8sConfig : import("./mesheryEnvironmentSlice").KubernetesCluster, meshAdapters: import("../mesheryComponents/mesheryComponentsSlice").AdaptersListType, grafana: import("./mesheryEnvironmentSlice").GrafanaType , promethues: import("./mesheryEnvironmentSlice").PrometheusType, loadTestPreferences: [], anonymousUsageStats: boolean, anonymousPerfResults: boolean, updatedAt: string | Date }} MesheryPreferenceType
 */

/**
 *
 * @returns {Promise<MesheryPreferenceType>}
 */
export const mesherySystemSync = () =>
  new Promise((res, rej) => {
    dataFetch(
      "/api/system/sync",
      {
        credentials: "same-origin",
        method: "GET",
      },
      (result) => res(result),
      (err) => rej(err)
    );
  });
