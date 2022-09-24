/**
 * Filter array of workloads by mesh-type
 *
 * @param {Array.<Object>} workloads array of raw workloads
 *
 * @returns {Array.Array.<Object>} filteredWorkloads: workloads filtered by Mesh Type
 */
export function groupWorkloadByType(workloads) {
  let filteredWorkloads = {};
  workloads.map((wtSet) => {
    if (wtSet.workload.metadata != null) {
      const adapterName = wtSet.workload.metadata["adapter.meshery.io/name"];
      let wl = filteredWorkloads[adapterName] || [];
      wl.push(wtSet);
      filteredWorkloads[adapterName] = wl;
    }
  });
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
  let versionedFilteredWorkloads = {};
  meshfilteredWorkloads.map(wtSet => {
    const version = wtSet.workload?.oam_definition?.spec?.metadata?.meshVersion ||  wtSet.workload.oam_definition?.spec?.metadata?.version || "Meshery";
    if (version) {
      let versionedFilteredMesh = versionedFilteredWorkloads[version] || [];
      versionedFilteredMesh.push(wtSet);
      versionedFilteredWorkloads[version] = versionedFilteredMesh;
    }
  });

  return versionedFilteredWorkloads;
}

export function getUnformattedName(oamWorkloadOrTrait) {
  return oamWorkloadOrTrait.workload.oam_definition?.metadata?.name || "Un-Named";
}

export function findWorkloadByName(name, workloads) {
  return workloads?.find(workload => getUnformattedName(workload) == name);
}
