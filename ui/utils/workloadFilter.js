/**
 * Filter array of workloads by mesh-type
 *
 * @param {Array.<Object>} workloads array of raw workloads
 *
 * @returns {Array.Array.<Object>} filteredWorkloads: workloads filtered by Mesh Type
 */
export function groupWorkloadByType(workloads) {
  let filteredWorkloads = {}
  workloads.map((workload) => {
    if (workload.metadata != null) {
      const adapterName = workload.metadata["adapter.meshery.io/name"];
      let wl = filteredWorkloads[adapterName] || []
      wl.push(workload)
      filteredWorkloads[adapterName] = wl;
    }
  })

  return filteredWorkloads;
}

/**
 * Filter-workloads by versions
 *
 * @param {Object} meshfilteredWorkloads array of mesh-filtered-workloads
 *
 * @returns {Array.Array.<Object>} versioned and typed filtered workloads
 */
export function groupWorkloadByVersion(meshfilteredWorkloads) {
  let versionedFilteredWorkloads = {}
  meshfilteredWorkloads.map(workload => {
    const version = workload?.oam_definition?.spec?.metadata?.meshVersion;
    if (version) {
      let versionedFilteredMesh = versionedFilteredWorkloads[version] || []
      versionedFilteredMesh.push(workload)
      versionedFilteredWorkloads[version] = versionedFilteredMesh
    }
  })

  return versionedFilteredWorkloads
}
