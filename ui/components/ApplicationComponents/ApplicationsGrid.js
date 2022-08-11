import React, {useState} from 'react'
import {Grid, Pagination} from "@mui/material"
import MesheryApplicationCard from "./ApplicationCard";
const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };


function ApplicationsGridItem({ application,  handleDeploy, handleUnDeploy, handleSubmit, setSelectedApplications }) {
    const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
    // const [yaml, setYaml] = useState(application.application_file);

  return (
    <Grid item {...gridProps}>
   <MesheryApplicationCard
   name={application.name}
   updated_at={application.updated_at}
   created_at={application.created_at}
   application_file={application.application_file}
   requestFullSize={() => setGridProps({ xl : 12, md : 12, xs : 12 })}
   requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
   handleDeploy={handleDeploy}
   handleUnDeploy={handleUnDeploy}
  //  deleteHandler={() => handleSubmit({ data : yaml, id : application.id, type : FILE_OPS.DELETE ,name : application.name })}
  //  updateHandler={() => handleSubmit({ data : yaml, id : application.id, type : FILE_OPS.UPDATE ,name : application.name })}
   setSelectedApplications={() => setSelectedApplications({ application : application, show : true })}
  //  setYaml={setYaml}
    />
   </Grid>
  )
}

/**
 * MesheryApplicationGrid is the react component for rendering grid
 * @param {{
 *  applications:Array<{
 *  id:string,
 *  created_at: string,
 *  updated_at: string,
 *  application_file: string,
 * }>,
 *  handleDeploy: (application_file: any) => void,
 *  handleUnDeploy: (application_file: any) => void,
 *  handleSubmit: (data: any, id: string, name: string, type: string) => void,
 *  setSelectedApplication : ({show: boolean, application:any}) => void,
 *  selectedApplication: {show : boolean, application : any},
 *  pages?: number,
 *  selectedPage?: number,
 *  setPage: (page: number) => void
 * }} props props
 */

 function MesheryApplicationGrid({ applications=[],handleDeploy, handleUnDeploy, handleSubmit, setSelectedApplication, selectedApplication, pages = 1,setPage, selectedPage }) {

  return (
    <div>
      {!selectedApplication.show &&
      <Grid container spacing={3} sx={{ padding : "1rem" }}>
        {applications.map((application) => (
          <ApplicationsGridItem
            key={application.id}
            application={application}
            handleDeploy={handleDeploy}
            // handleUnDeploy={handleUnDeploy}
            // handleSubmit={handleSubmit}
            setSelectedApplications={setSelectedApplication}
          />
        ))}

      </Grid>
      }
        {applications.length
        ? (
          <div >
            <Pagination count={pages} page={selectedPage+1} onChange={(_, page) => setPage(page - 1)} />  
          </div>
        )
        : null}
    </div>
  );
}

export default MesheryApplicationGrid