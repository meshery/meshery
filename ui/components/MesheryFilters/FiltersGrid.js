//@ts-check
import { Grid } from "@material-ui/core";
import { Pagination } from "@material-ui/lab";
import React, { useState } from "react";
import FiltersCard from "./FiltersCard";
import { makeStyles } from "@material-ui/core/styles";
import FILE_OPS from "../../utils/configurationFileHandlersEnum";


const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function FilterCardGridItem({ filter, handleDeploy, handleSubmit, setSelectedFilters }) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [yaml, setYaml] = useState(filter.filter_file);

  return (
    <Grid item {...gridProps}>
      <FiltersCard
        // id={filter.id}
        name={filter.name}
        updated_at={filter.updated_at}
        created_at={filter.created_at}
        filter_file={filter.filter_file}
        requestFullSize={() => setGridProps({ xl : 12, md : 12, xs : 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
        handleDeploy={() => handleDeploy(filter.filter_file)}
        deleteHandler={() => handleSubmit({ data : yaml, id : filter.id, type : FILE_OPS.DELETE ,name : filter.name })}
        updateHandler={() => handleSubmit({ data : yaml, id : filter.id, type : FILE_OPS.UPDATE ,name : filter.name })}
        setSelectedFilters={() => setSelectedFilters({ filter : filter, show : true })}
        setYaml={setYaml}
      />
    </Grid>
  );
}
const useStyles = makeStyles(() => ({
  pagination : {
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    marginTop : "2rem"
  }
}))

/**
 * FiltersGrid is the react component for rendering grid
 * @param {{
 *  filters:Array<{
 *  id:string,
 *  created_at: string,
 *  updated_at: string,
 *  filter_file: string,
 * }>,
 *  handleDeploy: (filter_file: any) => void,
 *  handleUnDeploy: (filter_file: any) => void,
 *  handleSubmit: (data: any, id: string, name: string, type: string) => void,
 *  setSelectedFilter : ({show: boolean, filter:any}) => void,
 *  selectedFilter: {show : boolean, filter : any},
 *  pages?: number,
 *  selectedPage?: number,
 *  setPage: (page: number) => void
 * }} props props
 */

function FiltersGrid({ filters=[],handleDeploy, handleUnDeploy, handleSubmit, setSelectedFilter, selectedFilter, pages = 1,setPage, selectedPage }) {

  const classes = useStyles()
  return (
    <div>
      {!selectedFilter.show &&
      <Grid container spacing={3} style={{ padding : "1rem" }}>
        {filters.map((filter) => (
          <FilterCardGridItem
            key={filter.id}
            filter={filter}
            handleDeploy={handleDeploy}
            handleSubmit={handleSubmit}
            setSelectedFilters={setSelectedFilter}
          />
        ))}

      </Grid>
      }
      {filters.length
        ? (
          <div className={classes.pagination} >
            <Pagination count={pages} page={selectedPage+1} onChange={(_, page) => setPage(page - 1)} />
          </div>
        )
        : null}
    </div>
  );
}

export default FiltersGrid;