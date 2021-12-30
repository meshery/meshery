import React, { useEffect, useState } from 'react'
import { SchemaContext } from "../utils/context/schemaSet"
import { groupWorkloadByType } from '../utils/workloadFilter';
import { createWorkloadTraitSets } from './MesheryMeshInterface/helpers';

function ConfigurationWrapper({ children }) {

  const [workloadTraitSet, setWorkloadTraitsSet] = useState([])
  const [meshWorkloads, setMeshWorkloads] = useState([])

  async function initialiser() {
    const wtsets = await createWorkloadTraitSets("");
    const mwloads = await groupWorkloadByType(wtsets);
    setWorkloadTraitsSet(wtsets);
    setMeshWorkloads(mwloads);
  }

  useEffect(() => {
    initialiser();
  }, []);

  return (
    <SchemaContext.Provider value={{ workloadTraitSet, meshWorkloads }}>
      {children}
    </SchemaContext.Provider>
  )
}

export default ConfigurationWrapper;