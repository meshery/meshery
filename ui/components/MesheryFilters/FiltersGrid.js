//@ts-check
import { Grid, Paper, Typography } from "@material-ui/core";
import { Pagination } from "@material-ui/lab";
import React, { useState } from "react";
import FiltersCard from "./FiltersCard";
import { makeStyles } from "@material-ui/core/styles";
import FILE_OPS from "../../utils/configurationFileHandlersEnum";
import ConfirmationMsg from "../ConfirmationModal";
import { getComponentsinFile } from "../../utils/utils";
import UploadImport from "../UploadImport";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function FilterCardGridItem({ filter, handleDeploy, handleUndeploy, handleSubmit, setSelectedFilters }) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [yaml, setYaml] = useState(filter.filter_file);

  return (
    <Grid item {...gridProps}>
      <FiltersCard
        name={filter.name}
        updated_at={filter.updated_at}
        created_at={filter.created_at}
        filter_file={filter.filter_file}
        requestFullSize={() => setGridProps({ xl : 12, md : 12, xs : 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
        handleDeploy={handleDeploy}
        handleUndeploy={handleUndeploy}
        deleteHandler={() => handleSubmit({ data : yaml, id : filter.id, type : FILE_OPS.DELETE ,name : filter.name })}
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
  },
  // text : {
  //   padding : "5px"
  // }
  noFilterPaper : {
    padding : "0.5rem",
    fontSize : "3rem"
  },
  noFilterContainer : {
    padding : "2rem",
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    flexDirection : "column",
  },
  noFilterText : {
    fontSize : "2rem",
    marginBottom : "2rem",
  },
}))

function FiltersGrid({ filters=[],handleDeploy, handleUndeploy, handleSubmit,urlUploadHandler,uploadHandler, setSelectedFilter, selectedFilter, pages = 1,setPage, selectedPage }) {

  const classes = useStyles()

  const [modalOpen, setModalOpen] = useState({
    open : false,
    deploy : false,
    filter_file : null,
    name : "",
    count : 0
  });

  const handleModalClose = () => {
    setModalOpen({
      open : false,
      filter_file : null,
      name : "",
      count : 0
    });
  }

  const handleModalOpen = (filter, isDeploy) => {
    setModalOpen({
      open : true,
      deploy : isDeploy,
      filter_file : filter.filter_file,
      name : filter.name,
      count : getComponentsinFile(filter.filter_file)
    });
  }

  return (
    <div>
      {!selectedFilter.show &&
      <Grid container spacing={3} style={{ padding : "1rem" }}>
        {filters.map((filter) => (
          <FilterCardGridItem
            key={filter.id}
            filter={filter}
            handleDeploy={() => handleModalOpen(filter, true)}
            handleUndeploy={() => handleModalOpen(filter, false)}
            handleSubmit={handleSubmit}
            setSelectedFilters={setSelectedFilter}
          />
        ))}

      </Grid>
      }
      {!selectedFilter.show && filters.length === 0 &&
        <Paper className={classes.noFilterPaper}>
          <div className={classes.noFilterContainer}>
            <Typography align="center" color="textSecondary" className={classes.noFilterText}>
            No Filters Found
            </Typography>
            <div>
              <UploadImport aria-label="URL upload button" handleUpload={urlUploadHandler} handleImport={uploadHandler} configuration="Filter"  />

            </div>
          </div>
        </Paper>
      }
      {filters.length
        ? (
          <div className={classes.pagination} >
            <Pagination count={pages} page={selectedPage+1} onChange={(_, page) => setPage(page - 1)} />
          </div>
        )
        : null}
      <ConfirmationMsg
        open={modalOpen.open}
        handleClose={handleModalClose}
        submit={
          { deploy : () => handleDeploy(modalOpen.filter_file), unDeploy : () => handleUndeploy(modalOpen.filter_file) }
        }
        isDelete={!modalOpen.deploy}
        title={ modalOpen.name }
        componentCount = {modalOpen.count}
        tab={modalOpen.deploy ? 0 : 1}
      />
    </div>
  );
}

export default FiltersGrid;