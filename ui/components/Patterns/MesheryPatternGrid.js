import React, {useState} from 'react'
import {Grid, Pagination} from "@mui/material"
import MesheryDesignCard from "./MesheryDesignCard";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function MesheryPatternsGridItem({ patterns, handleDeploy, handleUnDeploy, handleSubmit, setSelectedPatterns }) {

    const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
    // const [yaml, setYaml] = useState(application.application_file);

  return (
    <Grid item {...gridProps}>
        <MesheryDesignCard 
         name={patterns.name}
         updated_at={patterns.updated_at}
         created_at={patterns.created_at}
         patterns_file={patterns.patterns_file}
         requestFullSize={() => setGridProps({ xl : 12, md : 12, xs : 12 })}
         requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
         handleDeploy={handleDeploy}
         handleUnDeploy={handleUnDeploy}
        //  deleteHandler={() => handleSubmit({ data : yaml, id : patterns.id, type : FILE_OPS.DELETE ,name : patterns.name })}
        //  updateHandler={() => handleSubmit({ data : yaml, id : patterns.id, type : FILE_OPS.UPDATE ,name : patterns.name })}
         setSelectedPatterns={() => setSelectedPatterns({ patterns : patterns, show : true })}
        //  setYaml={setYaml}
        />
    </Grid>
  )
}

/**
 * MesheryPatternsGrid is the react component for rendering grid
 * @param {{
 *  patterns:Array<{
 *  id:string,
 *  created_at: string,
 *  updated_at: string,
 *  patterns_file: string,
 * }>,
 *  handleDeploy: (patterns_file: any) => void,
 *  handleUnDeploy: (patterns_file: any) => void,
 *  handleSubmit: (data: any, id: string, name: string, type: string) => void,
 *  setSelectedPatterns : ({show: boolean, patterns:any}) => void,
 *  selectedPatterns: {show : boolean, patterns : any},
 *  pages?: number,
 *  selectedPage?: number,
 *  setPage: (page: number) => void
 * }} props props
 */

function MesheryPatternGrid({ patterns=[],handleDeploy, handleUnDeploy, handleSubmit, setSelectedPatterns, selectedPattern, pages = 1,setPage, selectedPage }) {
  
    return (
        <div>
        {!selectedPattern.show &&
        <Grid container spacing={3} sx={{ padding : "1rem" }}>
          {patterns.map((patterns) => (
            <MesheryPatternsGridItem
              key={patterns.id}
              patterns={patterns}
            //   handleDeploy={() => handleModalOpen(pattern, true)}
            // handleUnDeploy={() => handleModalOpen(pattern, false)}
              setSelectedPatterns={setSelectedPatterns}
            />
          ))}
  
        </Grid>
        }
        {patterns.length
        ? (
          <div  >
            <Pagination count={pages} page={selectedPage+1} onChange={(_, page) => setPage(page - 1)} />  
          </div>
        )
        : null}
      </div>
        );
}

export default MesheryPatternGrid