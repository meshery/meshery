//@ts-check
import { Grid, Typography } from "@material-ui/core";
import { Pagination } from "@material-ui/lab";
import React, { useState } from "react";
import FiltersCard from "./FiltersCard";
import { makeStyles } from "@material-ui/core/styles";
import FILE_OPS from "../../utils/configurationFileHandlersEnum";
import ConfirmationMsg from "../ConfirmationModal";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function FilterCardGridItem({ filter, handleDeploy, handleSubmit, setSelectedFilters }) {
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
        handleDeploy={() => handleDeploy(filter.filter_file)}
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
  text : {
    padding : "5px"
  }
}))

function FiltersGrid({ filters=[],handleDeploy, handleSubmit, setSelectedFilter, selectedFilter, pages = 1,setPage, selectedPage }) {

  const classes = useStyles()

  const [modalOpen, setModalOpen] = useState({
    open : false,
    deploy : false,
    filter_file : null
  });

  const handleModalClose = () => {
    setModalOpen({
      open : false,
      deploy : false,
      filter_file : null
    });
  }

  const handleModalOpen = (filter_file, isDeploy) => {
    setModalOpen({
      open : true,
      deploy : isDeploy,
      filter_file : filter_file
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
            handleDeploy={() => handleModalOpen(filter.filter_file, true)}
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
      <ConfirmationMsg
        open={modalOpen.open}
        handleClose={handleModalClose}
        submit={() => handleDeploy(modalOpen.filter_file)}
        isDelete={!modalOpen.deploy}
        title={<Typography variant="h6" className={classes.text} >The selected operation will be applied to following contexts.</Typography>}
      />
    </div>
  );
}

export default FiltersGrid;