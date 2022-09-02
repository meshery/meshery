//@ts-check
import { Grid, Paper, Typography, Button } from "@material-ui/core";
import { Pagination } from "@material-ui/lab";
import React, { useState } from "react";
import FiltersCard from "./FiltersCard";
import FILE_OPS from "../../utils/configurationFileHandlersEnum";
import ConfirmationMsg from "../ConfirmationModal";
import { getComponentsinFile } from "../../utils/utils";
import PublishIcon from "@material-ui/icons/Publish";
import useStyles from "../MesheryPatterns/Grid.styles";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function FilterCardGridItem({ filter, handleDeploy, handleUndeploy, handleSubmit, setSelectedFilters, handleClone }) {
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
        handleClone={handleClone}
        deleteHandler={() => handleSubmit({ data : yaml, id : filter.id, type : FILE_OPS.DELETE ,name : filter.name })}
        setSelectedFilters={() => setSelectedFilters({ filter : filter, show : true })}
        setYaml={setYaml}
        description={filter.desciption}
        visibility={filter.visibility}
      />
    </Grid>
  );
}

function FiltersGrid({ filters=[],handleDeploy, handleUndeploy, handleClone, handleSubmit,urlUploadHandler,uploadHandler, setSelectedFilter, selectedFilter, pages = 1,setPage, selectedPage, UploadImport }) {

  const classes = useStyles()

  const [importModal, setImportModal] = useState({
    open : false
  });

  const handleUploadImport = () => {
    setImportModal({
      open : true
    });
  }

  const handleUploadImportClose = () => {
    setImportModal({
      open : false
    });
  }

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
            handleClone={() => handleClone(filter.id)}
            handleDeploy={() => handleModalOpen(filter, true)}
            handleUndeploy={() => handleModalOpen(filter, false)}
            handleSubmit={handleSubmit}
            setSelectedFilters={setSelectedFilter}
          />
        ))}

      </Grid>
      }
      {!selectedFilter.show && filters.length === 0 &&
        <Paper className={classes.noPaper}>
          <div className={classes.noContainer}>
            <Typography align="center" color="textSecondary" className={classes.noText}>
            No Filters Found
            </Typography>
            <div>
              <Button
                aria-label="Add Application"
                variant="contained"
                color="primary"
                size="large"
                // @ts-ignore
                onClick={handleUploadImport}
                style={{ marginRight : "2rem" }}
              >
                <PublishIcon className={classes.addIcon} />
              Import Filter
              </Button>
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
      <UploadImport open={importModal.open} handleClose={handleUploadImportClose} aria-label="URL upload button" handleUrlUpload={urlUploadHandler} handleUpload={uploadHandler} configuration="Filter"  />
    </div>
  );
}

export default FiltersGrid;