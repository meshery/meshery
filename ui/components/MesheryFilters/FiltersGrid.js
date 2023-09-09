//@ts-check
import { Grid, Paper, Typography, Button } from "@material-ui/core";
import { Pagination } from "@material-ui/lab";
import React, { useState } from "react";
import FiltersCard from "./FiltersCard";
import { FILE_OPS } from "../../utils/Enum";
import ConfirmationMsg from "../ConfirmationModal";
import { getComponentsinFile } from "../../utils/utils";
import PublishIcon from "@material-ui/icons/Publish";
import useStyles from "../MesheryPatterns/Grid.styles";
import Modal from "../Modal";
import Filter from "../../public/static/img/drawer-icons/filter_svg.js";
import PublicIcon from '@material-ui/icons/Public';
import safeJsonParse from "../ConnectionWizard/helpers/safeJsonParse";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function FilterCardGridItem({
  filter,
  yamlConfig,
  handleDeploy,
  handleDownload,
  handleUndeploy,
  handleSubmit,
  setSelectedFilters,
  handleClone,
  handlePublishModal,
  handleUnpublishModal,
  canPublishFilter,
}) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [yaml, setYaml] = useState(yamlConfig);

  return (
    <Grid item {...gridProps}>
      <FiltersCard
        name={filter.name}
        updated_at={filter.updated_at}
        created_at={filter.created_at}
        filter_resource={yaml}
        canPublishFilter={canPublishFilter}
        handlePublishModal={handlePublishModal}
        handleUnpublishModal={handleUnpublishModal}
        requestFullSize={() => setGridProps({ xl : 12, md : 12, xs : 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
        handleDeploy={handleDeploy}
        handleUndeploy={handleUndeploy}
        handleClone={handleClone}
        handleDownload={(ev) => handleDownload(ev, filter.id, filter.name)}
        deleteHandler={() => handleSubmit({ data : yaml, id : filter.id, type : FILE_OPS.DELETE, name : filter.name })}
        updateHandler={() => handleSubmit({ data : yaml, id : filter.id, type : FILE_OPS.UPDATE, name : filter.name })}
        setSelectedFilters={() => setSelectedFilters({ filter : filter, show : true })}
        setYaml={setYaml}
        description={filter.desciption}
        visibility={filter.visibility}
      />
    </Grid>
  );
}

/**
 * FilterGrid is the react component for rendering grid
 * @param {{
 *  filters:Array<{
 *  id:string,
 *  created_at: string,
 *  updated_at: string,
 *  filter_file: string,
 * }>,
 *  handleVerify: (e: Event, filter_file: any, filter_id: string) => void,
 *  handlePublish: (catalog_data : any) => void,
 *  handleUnpublishModal: (ev: Event, filter: any) => (() => Promise<void>),
 *  handleDeploy: (filter_file: any) => void,
 *  handleUnDeploy: (filter_file: any) => void,
 *  handleDownload: (e : Event, id : string , name : string ) => void,
 *  handleSubmit: (data: any, id: string, name: string, type: string) => void,
 *  setSelectedFilter : ({show: boolean, filter:any}) => void,
 *  selectedFilter: {show : boolean, filter : any},
 *  pages?: number,
 *  selectedPage?: number,
 *  setPage: (page: number) => void,
 *  filterErrors: Map,
 *  canPublishFilter: boolean,
 *  publishModal: {
 *   open: boolean,
 *   filter: any,
 *   name: string,
 *  },
 *  setPublishModal: (publishModal: { open: boolean, filter: any, name: string }) => void,
 *  publishSchema: object,
 * }} props props
 */

function FiltersGrid({
  filters = [],
  handleDeploy,
  handleUndeploy,
  handleClone,
  handleDownload,
  handleSubmit,
  setSelectedFilter,
  selectedFilter,
  pages = 1,
  setPage,
  selectedPage,
  importSchema,
  canPublishFilter,
  handlePublish,
  handleUnpublishModal,
  handleImportFilter,
  publishModal,
  setPublishModal,
  publishSchema
}) {
  const classes = useStyles();

  const [importModal, setImportModal] = useState({
    open : false,
  });

  const handlePublishModal = (filter) => {
    if (canPublishFilter) {
      setPublishModal({
        open : true,
        filter : filter,
        name : "",
      });
    }
  };

  const handlePublishModalClose = () => {
    setPublishModal({
      open : false,
      filter : {},
      name : "",
    });
  };

  const handleUploadImport = () => {
    setImportModal({
      open : true,
    });
  };

  const handleUploadImportClose = () => {
    setImportModal({
      open : false,
    });
  };

  const [modalOpen, setModalOpen] = useState({
    open : false,
    deploy : false,
    filter_file : null,
    name : "",
    count : 0,
  });

  const handleModalClose = () => {
    setModalOpen({
      open : false,
      filter_file : null,
      name : "",
      count : 0,
    });
  };

  const handleModalOpen = (filter, isDeploy) => {
    setModalOpen({
      open : true,
      deploy : isDeploy,
      filter_file : filter.filter_file,
      name : filter.name,
      count : getComponentsinFile(filter.filter_file),
    });
  };

  const getYamlConfig = (filter_resource) => {
    if (filter_resource) {
      return safeJsonParse(filter_resource).settings.config;
    }

    return "";
  };

  return (
    <div>
      {!selectedFilter.show && (
        <Grid container spacing={3} style={{ padding : "1rem" }}>
          {filters.map((filter) => (
            <FilterCardGridItem
              key={filter.id}
              filter={filter}
              yamlConfig={getYamlConfig(filter.filter_resource)}
              handleClone={() => handleClone(filter.id, filter.name)}
              handleDownload={handleDownload}
              handleDeploy={() => handleModalOpen(filter, true)}
              handleUndeploy={() => handleModalOpen(filter, false)}
              handleSubmit={handleSubmit}
              setSelectedFilters={setSelectedFilter}
              canPublishFilter={canPublishFilter}
              handlePublishModal={() => handlePublishModal(filter)}
              handleUnpublishModal={(e) => handleUnpublishModal(e, filter)()}
            />
          ))}
        </Grid>
      )}
      {!selectedFilter.show && filters.length === 0 && (
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
      )}
      {filters.length ? (
        <div className={classes.pagination}>
          <Pagination count={pages} page={selectedPage + 1} onChange={(_, page) => setPage(page - 1)} />
        </div>
      ) : null}
      <ConfirmationMsg
        open={modalOpen.open}
        handleClose={handleModalClose}
        submit={{
          deploy : () => handleDeploy(modalOpen.filter_file),
          unDeploy : () => handleUndeploy(modalOpen.filter_file),
        }}
        isDelete={!modalOpen.deploy}
        title={modalOpen.name}
        componentCount={modalOpen.count}
        tab={modalOpen.deploy ? 2 : 1}
      />
      {canPublishFilter && (
        <Modal
          open={publishModal.open}
          schema={publishSchema.rjsfSchema}
          uiSchema={publishSchema.uiSchema}
          title={publishModal.filter?.name}
          handleClose={handlePublishModalClose}
          handleSubmit={handlePublish}
          showInfoIcon={{ text : "Upon submitting your catalog item, an approval flow will be initiated.", link : "https://docs.meshery.io/concepts/catalog" }}
          submitBtnText="Submit for Approval"
          submitBtnIcon={<PublicIcon/>}
        />
      )}
      <Modal
        open={importModal.open}
        schema={importSchema.rjsfSchema}
        uiSchema={importSchema.uiSchema}
        handleClose={handleUploadImportClose}
        handleSubmit={handleImportFilter}
        title="Import Filter"
        submitBtnText="Import"
        leftHeaderIcon={<Filter fill="#fff" style={{ height : "24px", width : "24px", fonSize : "1.45rem" }} />}
        submitBtnIcon={<PublishIcon/>}
      />
    </div>
  );
}

export default FiltersGrid;