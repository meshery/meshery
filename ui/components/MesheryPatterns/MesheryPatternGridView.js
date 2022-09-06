//@ts-check
import { Grid, Paper, Typography, Button } from "@material-ui/core";
import { Pagination } from "@material-ui/lab";
import React, { useState } from "react";
import MesheryPatternCard from "./MesheryPatternCard";
import PatternConfiguratorComponent from "../configuratorComponents/patternConfigurator"
import FILE_OPS from "../../utils/configurationFileHandlersEnum";
import ConfirmationMsg from "../ConfirmationModal";
import { getComponentsinFile } from "../../utils/utils";
import PublishIcon from "@material-ui/icons/Publish";
import useStyles from "./Grid.styles";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function PatternCardGridItem({ pattern, handleDeploy, handleUnDeploy, handleClone, handleSubmit, setSelectedPatterns }) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [yaml, setYaml] = useState(pattern.pattern_file);

  return (
    <Grid item {...gridProps}>
      <MesheryPatternCard
        // id={pattern.id}
        name={pattern.name}
        updated_at={pattern.updated_at}
        created_at={pattern.created_at}
        pattern_file={pattern.pattern_file}
        requestFullSize={() => setGridProps({ xl : 12, md : 12, xs : 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
        handleDeploy={handleDeploy}
        handleUnDeploy={handleUnDeploy}
        handleClone={handleClone}
        deleteHandler={() => handleSubmit({ data : yaml, id : pattern.id, type : FILE_OPS.DELETE ,name : pattern.name })}
        updateHandler={() => handleSubmit({ data : yaml, id : pattern.id, type : FILE_OPS.UPDATE ,name : pattern.name })}
        setSelectedPatterns={() => setSelectedPatterns({ pattern : pattern, show : true })}
        setYaml={setYaml}
        description={pattern.description}
        visibility={pattern.visibility}
      />
    </Grid>
  );
}

/**
 * MesheryPatternGrid is the react component for rendering grid
 * @param {{
 *  patterns:Array<{
 *  id:string,
 *  created_at: string,
 *  updated_at: string,
 *  pattern_file: string,
 * }>,
 *  handleDeploy: (pattern_file: any) => void,
 *  handleUnDeploy: (pattern_file: any) => void,
 *  handleSubmit: (data: any, id: string, name: string, type: string) => void,
 *  setSelectedPattern : ({show: boolean, pattern:any}) => void,
 *  selectedPattern: {show : boolean, pattern : any},
 *  pages?: number,
 *  selectedPage?: number,
 *  setPage: (page: number) => void
 * }} props props
 */

function MesheryPatternGrid({ patterns=[], handleDeploy, handleUnDeploy, urlUploadHandler, handleClone, uploadHandler, handleSubmit, setSelectedPattern, selectedPattern, pages = 1,setPage, selectedPage, UploadImport  }) {

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
    pattern_file : null,
    name : "",
    count : 0
  });

  const handleModalClose = () => {
    setModalOpen({
      open : false,
      pattern_file : null,
      name : "",
      count : 0
    });
  }

  const handleModalOpen = (pattern, isDeploy) => {
    setModalOpen({
      open : true,
      deploy : isDeploy,
      pattern_file : pattern.pattern_file,
      name : pattern.name,
      count : getComponentsinFile(pattern.pattern_file)
    });
  }

  return (
    <div>
      {selectedPattern.show &&
      <PatternConfiguratorComponent pattern={selectedPattern.pattern} show={setSelectedPattern}  onSubmit={handleSubmit} />
      }
      {!selectedPattern.show &&
      <Grid container spacing={3} style={{ padding : "1rem" }}>
        {patterns.map((pattern) => (
          <PatternCardGridItem
            key={pattern.id}
            pattern={pattern}
            handleClone={() => handleClone(pattern.id)}
            handleDeploy={() => handleModalOpen(pattern, true)}
            handleUnDeploy={() => handleModalOpen(pattern, false)}
            handleSubmit={handleSubmit}
            setSelectedPatterns={setSelectedPattern}
          />
        ))}

      </Grid>
      }
      {!selectedPattern.show && patterns.length === 0 &&
          <Paper className={classes.noPaper}>
            <div className={classes.noContainer}>
              <Typography align="center" color="textSecondary" className={classes.noText}>
                No Designs Found
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
              Import Design
                </Button>
              </div>
            </div>
          </Paper>
      }
      {patterns.length
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
          { deploy : () => handleDeploy(modalOpen.pattern_file), unDeploy : () => handleUnDeploy(modalOpen.pattern_file) }
        }
        isDelete={!modalOpen.deploy}
        title={ modalOpen.name }
        componentCount={modalOpen.count}
        tab={modalOpen.deploy ? 0 : 1}
      />
      <UploadImport open={importModal.open} handleClose={handleUploadImportClose} aria-label="URL upload button" handleUrlUpload={urlUploadHandler} handleUpload={uploadHandler} configuration="Designs"  />
    </div>
  );
}

export default MesheryPatternGrid;