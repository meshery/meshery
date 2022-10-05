import React, { useState } from "react";
import { Grid, Box, Pagination,Typography } from "@mui/material";
import MesheryCard from "../MesheryCard";

const INITIAL_GRID_SIZE = { xl: 4, md: 6, xs: 12 };

function MesheryFilterGridItem({ filter, handleDeploy, handleUnDeploy, handleSubmit, setSelectedFilter }) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  // const [yaml, setYaml] = useState(filetrs.filter_file);

  return (
    <Grid item {...gridProps}>
      <MesheryCard
        name={filter.name}
        updated_at={filter.updated_at}
        created_at={filter.created_at}
         file={filter.filetr_file}
        type="filter"
        requestFullSize={() => setGridProps({ xl: 12, md: 12, xs: 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
         handleDeploy={() => handleDeploy(filters.filterss_file)}
         handleUnDeploy={() => handleUnDeploy(filters.filters_file)}
        //  deleteHandler={() => handleSubmit({ data : yaml, id : filters.id, type : FILE_OPS.DELETE ,name : filters.name })}
        //  updateHandler={() => handleSubmit({ data : yaml, id : filters.id, type : FILE_OPS.UPDATE ,name : filters.name })}
        // setSelectedFilter={() => setSelectedFilter({ fileter : filter, show : true })}
        //  setYaml={setYaml}
      />
    </Grid>
  );
}

/**
 * MesheryPatternsGrid is the react component for rendering grid
 * @param {{
 *  filters:Array<{
 *  id:string,
 *  created_at: string,
 *  updated_at: string,
 *  patterns_file: string,
 * }>,
 *  handleDeploy: (fileter_file: any) => void,
 *  handleUnDeploy: (filter_file: any) => void,
 *  handleSubmit: (data: any, id: string, name: string, type: string) => void,
 *  setSelectedFilter : ({show: boolean, filter:any}) => void,
 *  selectedFilter: {show : boolean, filter : any},
 *  pages?: number,
 *  selectedPage?: number,
 *  setPage: (page: number) => void
 * }} props props
 */

function MesheryFilterGrid({
  filters = [],
  handleDeploy,
  handleUnDeploy,
  handleSubmit,
  setSelectedFilter,
  selectedFilter,
  pages = 1,
  setPage,
  selectedPage,
}) {
  return (
    <div>
      {!selectedFilter.show && (
        <Grid container spacing={3} sx={{ padding: "1rem" }}>
          {filters.map((filter) => (
            <MesheryFilterGridItem
              key={filter.id}
              filters={filter}
              handleDeploy={handleDeploy}
              // handleUnDeploy={handleUnDeploy}
              // handleSubmit={handleSubmit}
              setSelectedFilter={setSelectedFilter}
            />
          ))}
          {filters.length ? (
            <Box display="flex" justifyContent="center" marginTop={2} alignItems="center">
              <Pagination count={pages} page={selectedPage + 1} onChange={(_, page) => setPage(page - 1)} />
            </Box>
          ) : null}
        </Grid>
      )}
    </div>
  );
}

export default MesheryFilterGrid;

