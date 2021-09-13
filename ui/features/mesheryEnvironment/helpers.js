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

/**
 *
 * @returns {Promise<{grafana: [any], prometheus: [any]}>} MeshScanData - Returns an array of Mesh Scan Data for Prometheus and Grafana
 */

export const promGrafMeshScan = () =>
  new Promise((res, rej) => {
    dataFetch(
      "/api/system/meshsync/grafana",
      {
        credentials: "same-origin",
        method: "GET",
      },
      (result) => {
        // use error code utility for generating errors
        if (!result) rej("Result is undefined");
        res(result);
      },
      (err) => rej(err)
    );
  });

/**
 * extractURLFromScanData scans the ingress urls from the
 * mesh scan data and returns an array of the response
 * @param {object[]} scannedData
 * @returns {string[]}
 */
export const extractURLFromScanData = (scannedData) => {
  const result = [];
  scannedData.forEach((data) => {
    // Add loadbalancer based url
    if (Array.isArray(data.status?.loadBalancer?.ingress)) {
      data.status.loadBalancer.ingress.forEach((lbdata) => {
        let protocol = "http";

        // Iterate over ports exposed by the service
        if (Array.isArray(data.spec.ports)) {
          data.spec.ports.forEach(({ port }) => {
            if (port === 443) protocol = "https";

            // From kubernetes v1.19 docs
            // Hostname is set for load-balancer ingress points that are DNS based (typically AWS load-balancers)
            // IP is set for load-balancer ingress points that are IP based (typically GCE or OpenStack load-balancers)
            let address = lbdata.ip || lbdata.hostname;
            if (address) result.push(`${protocol}://${address}:${port}`);
          });
        }
      });
    }

    // Add clusterip based url
    // As per kubernetes v1.19 api, "None", "" as well as a valid ip is a valid clusterIP
    // Looking for valid ipv4 address
    if (data.spec.clusterIP?.match(/^((25[0-5]|(2[0-4]|1[0-9]|[1-9]|)[0-9])(\.(?!$)|$)){4}$/g)?.[0]) {
      let protocol = "http";
      if (Array.isArray(data.spec.ports)) {
        data.spec.ports.forEach(({ port }) => {
          if (port === 443) protocol = "https";
          result.push(`${protocol}://${data.spec.clusterIP}:${port}`);
        });
      }
    }
  });

  return result;
};
