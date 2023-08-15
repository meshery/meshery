//@ts-check
import { Grid, Paper, Typography, Button } from "@material-ui/core";
import { Pagination } from "@material-ui/lab";
import React, { useState } from "react";
import MesheryApplicationCard from "./ApplicationsCard";
import ConfirmationModal from "../ConfirmationModal";
import { getComponentsinFile } from "../../utils/utils";
import PublishIcon from "@material-ui/icons/Publish";
import useStyles from "../MesheryPatterns/Grid.styles";
import { FILE_OPS } from "../../utils/Enum";
import DryRunComponent from "../DryRun/DryRunComponent";


const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function ApplicationsGridItem({ application,  handleDeploy, handleUnDeploy, handleSubmit, setSelectedApplications, handleAppDownload }) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [yaml, setYaml] = useState(application.application_file);

  return (
    <Grid item {...gridProps}>
      <MesheryApplicationCard
        id={application.id}
        name={application.name}
        updated_at={application.updated_at}
        created_at={application.created_at}
        application_file={application.application_file}
        requestFullSize={() => setGridProps({ xl : 12, md : 12, xs : 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
        handleDeploy={handleDeploy}
        handleUnDeploy={handleUnDeploy}
        deleteHandler={() => handleSubmit({ data : yaml, id : application.id, type : FILE_OPS.DELETE ,name : application.name })}
        updateHandler={() => handleSubmit({ data : yaml, id : application.id, type : FILE_OPS.UPDATE ,name : application.name, source_type : application.type.String })}
        setSelectedApplications={() => setSelectedApplications({ application : application, show : true })}
        setYaml={setYaml}
        source_type={application.type.String}
        handleAppDownload={handleAppDownload}
      />
    </Grid>
  );
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
 *  urlUploadHandler: any
 *  uploadHandler: Function
 *  selectedK8sContexts : Array
 *  UploadImport: any
 * handleAppDownload: Function
 * }} props props
 */

function MesheryApplicationGrid({ applications=[],handleDeploy, handleUnDeploy, handleSubmit,urlUploadHandler,uploadHandler, setSelectedApplication, selectedApplication, pages = 1,setPage, selectedPage, UploadImport, fetch, handleAppDownload ,selectedK8sContexts }) {

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
    application_file : null,
    name : "",
    count : 0,
    dryRunComponent : null
  });

  const handleModalClose = () => {
    setModalOpen((prevState) => ({
      ...prevState,
      open : false,
      application_file : null,
      name : "",
      count : 0
    } ))
  }

  const handleModalOpen = (app, isDeploy) => {

    const dryRunComponent =(
      <DryRunComponent
        design={app.application_file}
        noOfElements={getComponentsinFile(app.application_file)}
        selectedContexts={selectedK8sContexts}
      ></DryRunComponent>)

    setModalOpen({
      open : true,
      deploy : isDeploy,
      application_file : app.application_file,
      name : app.name,
      dryRunComponent,
      count : getComponentsinFile(app.application_file)
    });
  }

  return (
    <div>
      {!selectedApplication.show &&
      <Grid container spacing={3} style={{ padding : "1rem" }}>
        {applications.map((application) => (
          <ApplicationsGridItem
            key={application.id}
            application={application}
            handleDeploy={() => handleModalOpen(application, true)}
            handleUnDeploy={() => handleModalOpen(application, false)}
            handleSubmit={handleSubmit}
            setSelectedApplications={setSelectedApplication}
            handleAppDownload={handleAppDownload}
          />
        ))}

      </Grid>
      }
      {!selectedApplication.show && applications.length === 0 &&
          <Paper className={classes.noPaper}>
            <div className={classes.noContainer}>
              <Typography align="center" color="textSecondary" className={classes.noText}>
                No Applications Found
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
              Import Application
                </Button>
              </div>
            </div>
          </Paper>
      }
      {applications.length
        ? (
          <div className={classes.pagination} >
            <Pagination count={pages} page={selectedPage+1} onChange={(_, page) => setPage(page - 1)} />
          </div>
        )
        : null}
      <ConfirmationModal
        open={modalOpen.open}
        handleClose={handleModalClose}
        submit={
          { deploy : () => handleDeploy(modalOpen.application_file, modalOpen.name), unDeploy : () => handleUnDeploy (modalOpen.application_file, modalOpen.name) }
        }
        isDelete={!modalOpen.deploy}
        title={ modalOpen.name }
        componentCount={ modalOpen.count }
        dryRunComponent={modalOpen.dryRunComponent}
        tab={modalOpen.deploy ? 2 : 1}
      />
      <UploadImport open={importModal.open} handleClose={handleUploadImportClose} isApplication = {true} aria-label="URL upload button" handleUrlUpload={urlUploadHandler} handleUpload={uploadHandler} fetch={() => fetch()} configuration="Application"  />
    </div>
  );
}

export default MesheryApplicationGrid;