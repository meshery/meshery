/**
 * Filter-workloads by versions
 *
 * @param {Object} meshfilteredWorkloads array of mesh-filtered-workloads
 *
 * @returns {Record<string, Array<Object>>} versioned and typed filtered workloads
 */
export function groupWorkloadByVersion(meshfilteredWorkloads) {
  const versionedFilteredWorkloads = {};
  if (!Array.isArray(meshfilteredWorkloads)) {
    return versionedFilteredWorkloads;
  }

  meshfilteredWorkloads.forEach((wtSet) => {
    if (!wtSet) return;
    const version =
      wtSet.workload?.oam_definition?.spec?.metadata?.meshVersion ||
      wtSet.workload?.oam_definition?.spec?.metadata?.version ||
      'Meshery';
    if (version) {
      const versionedFilteredMesh = versionedFilteredWorkloads[version] || [];
      versionedFilteredMesh.push(wtSet);
      versionedFilteredWorkloads[version] = versionedFilteredMesh;
    }
  });

  return versionedFilteredWorkloads;
}

export function getUnformattedName(oamWorkloadOrTrait) {
  return oamWorkloadOrTrait?.workload?.oam_definition?.metadata?.name || 'Un-Named';
}

export function findWorkloadByName(name, workloads) {
  if (!Array.isArray(workloads)) {
    return undefined;
  }
  return workloads.find((workload) => getUnformattedName(workload) === name);
}
