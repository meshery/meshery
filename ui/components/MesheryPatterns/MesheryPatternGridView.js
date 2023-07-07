//@ts-check
import { Grid, Paper, Typography, Button } from "@material-ui/core";
import { Pagination } from "@material-ui/lab";
import React, { useState } from "react";
import MesheryPatternCard from "./MesheryPatternCard";
import DesignConfigurator from "../configuratorComponents/MeshModel";
import { FILE_OPS, ACTIONS } from "../../utils/Enum";
import ConfirmationMsg from "../ConfirmationModal";
import { getComponentsinFile } from "../../utils/utils";
import PublicIcon from '@material-ui/icons/Public';
import PublishIcon from "@material-ui/icons/Publish";
import useStyles from "./Grid.styles";
import Validation from "../Validation";
import { publish_schema } from "../schemas/publish_schema";
import Modal from "../Modal";
import _ from "lodash";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function PatternCardGridItem({ pattern, handleDeploy, handleVerify, handlePublishModal, handleUnpublishModal, handleUnDeploy, handleClone, handleSubmit, setSelectedPatterns, canPublishPattern = false }) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [yaml, setYaml] = useState(pattern.pattern_file);

  return (
    <Grid item {...gridProps}>
      <MesheryPatternCard
        // id={pattern.id}
        canPublishPattern={canPublishPattern}
        name={pattern.name}
        updated_at={pattern.updated_at}
        created_at={pattern.created_at}
        pattern_file={pattern.pattern_file}
        requestFullSize={() => setGridProps({ xl : 12, md : 12, xs : 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
        handleDeploy={handleDeploy}
        handleVerify={handleVerify}
        handlePublishModal={handlePublishModal}
        handleUnDeploy={handleUnDeploy}
        handleUnpublishModal={handleUnpublishModal}
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
 *  handleVerify: (e: Event, pattern_file: any, pattern_id: string) => void,
 *  handlePublish: (catalog_data : any) => void,
 *  handleUnpublishModal: (ev: Event, pattern: any) => (() => Promise<void>),
 *  handleDeploy: (pattern_file: any) => void,
 *  handleUnDeploy: (pattern_file: any) => void,
 *  handleSubmit: (data: any, id: string, name: string, type: string) => void,
 *  setSelectedPattern : ({show: boolean, pattern:any}) => void,
 *  selectedPattern: {show : boolean, pattern : any},
 *  pages?: number,
 *  selectedPage?: number,
 *  setPage: (page: number) => void
 *  patternErrors: Map
 *  canPublishPattern: boolean
 * }} props props
 */

function MesheryPatternGrid({ patterns=[], handleVerify, handlePublish, handleUnpublishModal, handleDeploy, handleUnDeploy, urlUploadHandler, handleClone, uploadHandler, handleSubmit, setSelectedPattern, selectedPattern, pages = 1,setPage, selectedPage, UploadImport, fetch, patternErrors, canPublishPattern = false }) {

  const classes = useStyles()

  const [importModal, setImportModal] = useState({
    open : false
  });
  const [publishModal, setPublishModal] = useState({
    open : false,
    pattern : {}
  });
  const handlePublishModal = (pattern) => {
    if (canPublishPattern) {
      setPublishModal({
        open : true,
        pattern : pattern
      });
    }
  };
  const handlePublishModalClose = () => {
    setPublishModal({
      open : false,
      pattern : {}
    });
  };

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
  const [payload, setPayload] = useState({
    id : "",
    catalog_data : {}
  });

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

  const onChange = (e) => {
    setPayload({
      id : publishModal.pattern?.id,
      catalog_data : e
    })
  }

  const handleModalOpen = (pattern, action) => {
    const compCount = getComponentsinFile(pattern.pattern_file);
    const validationBody = (
      <Validation
        errors={patternErrors.get(pattern.id)}
        compCount={compCount}
        handleClose={() => setModalOpen({ ...modalOpen, open : false })}
      />
    )
    setModalOpen({
      open : true,
      action : action,
      pattern_file : pattern.pattern_file,
      name : pattern.name,
      count : compCount,
      validationBody : validationBody
    });
  }

  return (
    <div>
      {selectedPattern.show &&
      <DesignConfigurator pattern={selectedPattern.pattern} show={setSelectedPattern}  onSubmit={handleSubmit} />
      }
      {!selectedPattern.show &&
      <Grid container spacing={3} style={{ padding : "1rem" }}>
        {patterns.map((pattern) => (
          <PatternCardGridItem
            key={pattern.id}
            pattern={pattern}
            canPublishPattern={canPublishPattern}
            handleClone={() => handleClone(pattern.id, pattern.name)}
            handleDeploy={() => handleModalOpen(pattern, ACTIONS.DEPLOY)}
            handleUnDeploy={() => handleModalOpen(pattern, ACTIONS.UNDEPLOY)}
            handleVerify={(e) => handleVerify(e, pattern.pattern_file, pattern.id)}
            handlePublishModal={() => handlePublishModal(pattern)}
            handleUnpublishModal={(e) => handleUnpublishModal(e, pattern)()}
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
          { deploy : () => handleDeploy(modalOpen.pattern_file, modalOpen.name), unDeploy : () => handleUnDeploy(modalOpen.pattern_file, modalOpen.name) }
        }
        title={ modalOpen.name }
        componentCount={modalOpen.count}
        tab={modalOpen.action}
        validationBody={modalOpen.validationBody}
      />
      {canPublishPattern &&
      <Modal open={publishModal.open} schema={publish_schema} onChange={onChange} handleClose={handlePublishModalClose} formData={_.isEmpty(payload.catalog_data) ?publishModal?.pattern?.catalog_data : payload.catalog_data} aria-label="catalog publish" title={publishModal.pattern?.name}>
        <Button
          title="Publish"
          variant="contained"
          color="primary"
          className={classes.testsButton}
          onClick={() => {
            handlePublishModalClose();
            handlePublish(payload)
          }}
        >
          <PublicIcon className={classes.iconPatt} />
          <span className={classes.btnText}> Publish </span>
        </Button>
      </Modal>
      }
      <UploadImport open={importModal.open} handleClose={handleUploadImportClose} aria-label="URL upload button" handleUrlUpload={urlUploadHandler} handleUpload={uploadHandler} fetch={async() => await fetch()} configuration="Designs"  />
    </div>
  );
}

export default MesheryPatternGrid;