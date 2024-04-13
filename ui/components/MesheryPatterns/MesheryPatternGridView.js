import { Grid, Paper, Typography } from '@material-ui/core';
import { Pagination } from '@material-ui/lab';
import React, { useState } from 'react';
import MesheryPatternCard from './MesheryPatternCard';
import DesignConfigurator from '../configuratorComponents/MeshModel';
import { FILE_OPS, ACTIONS } from '../../utils/Enum';
import { EVENT_TYPES } from '../../lib/event-types';
import ConfirmationMsg from '../ConfirmationModal';
import { getComponentsinFile } from '../../utils/utils';
import useStyles from './Grid.styles';
import Validation from '../Validation';
import Modal from '../Modal';
import PublicIcon from '@material-ui/icons/Public';
import DryRunComponent from '../DryRun/DryRunComponent';
import { withSnackbar } from 'notistack';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateProgress } from '../../lib/store';
import ExportModal from '../ExportModal';
import downloadContent from '@/utils/fileDownloader';
import { useNotification } from '@/utils/hooks/useNotification';
const INITIAL_GRID_SIZE = { xl: 4, md: 6, xs: 12 };

function PatternCardGridItem({
  pattern,
  handleDeploy,
  handleVerify,
  handlePublishModal,
  handleUnpublishModal,
  handleUnDeploy,
  handleClone,
  handleSubmit,
  handleDownload,
  setSelectedPatterns,
  canPublishPattern = false,
  user,
  handleInfoModal,
}) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [yaml, setYaml] = useState(pattern.pattern_file);

  return (
    <Grid item {...gridProps}>
      <MesheryPatternCard
        id={pattern.id}
        user={user}
        canPublishPattern={canPublishPattern}
        name={pattern.name}
        updated_at={pattern.updated_at}
        created_at={pattern.created_at}
        pattern_file={pattern.pattern_file}
        requestFullSize={() => setGridProps({ xl: 12, md: 12, xs: 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
        handleDeploy={handleDeploy}
        handleVerify={handleVerify}
        handlePublishModal={handlePublishModal}
        handleUnDeploy={handleUnDeploy}
        handleUnpublishModal={handleUnpublishModal}
        handleClone={handleClone}
        handleInfoModal={handleInfoModal}
        handleDownload={handleDownload}
        deleteHandler={() =>
          handleSubmit({
            data: yaml,
            id: pattern.id,
            type: FILE_OPS.DELETE,
            name: pattern.name,
            catalog_data: pattern.catalog_data,
          })
        }
        updateHandler={() =>
          handleSubmit({
            data: yaml,
            id: pattern.id,
            type: FILE_OPS.UPDATE,
            name: pattern.name,
            catalog_data: pattern.catalog_data,
          })
        }
        setSelectedPatterns={() => setSelectedPatterns({ pattern: pattern, show: true })}
        setYaml={setYaml}
        description={pattern.description}
        visibility={pattern.visibility}
        pattern={pattern}
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
 *  canPublishPattern: boolean,
 *  publishModal: {
 *   open: boolean,
 *   filter: any,
 *   name: string
 *  },
 *  setPublishModal: (publishModal: { open: boolean, filter: any, name: string }) => void
 * }} props props
 */

function MesheryPatternGrid({
  patterns = [],
  handleVerify,
  handlePublish,
  handleUnpublishModal,
  handleDeploy,
  handleUnDeploy,
  handleClone,
  handleSubmit,
  setSelectedPattern,
  selectedPattern,
  pages = 1,
  setPage,
  selectedPage,
  patternErrors,
  canPublishPattern = false,
  publishModal,
  setPublishModal,
  selectedK8sContexts,
  publishSchema,
  user,
  updateProgress,
  handleInfoModal,
}) {
  const classes = useStyles();
  const { notify } = useNotification();
  const handlePublishModal = (pattern) => {
    if (canPublishPattern) {
      setPublishModal({
        open: true,
        pattern: pattern,
        name: '',
      });
    }
  };
  const handlePublishModalClose = () => {
    setPublishModal({
      open: false,
      pattern: {},
      name: '',
    });
  };

  const [modalOpen, setModalOpen] = useState({
    open: false,
    deploy: false,
    pattern_file: null,
    pattern_id: '',
    name: '',
    count: 0,
    dryRunComponent: null,
  });

  const [downloadModal, setDownloadModal] = useState({
    open: false,
    content: null,
  });
  const handleDownload = (e, design, source_type, params) => {
    e.stopPropagation();
    updateProgress({ showProgress: true });
    try {
      let id = design.id;
      let name = design.name;
      downloadContent({ id, name, type: 'pattern', source_type, params });
      updateProgress({ showProgress: false });
      notify({ message: `"${name}" design downloaded`, event_type: EVENT_TYPES.INFO });
    } catch (e) {
      console.error(e);
    }
  };
  const handleDownloadDialogClose = () => {
    setDownloadModal((prevState) => ({
      ...prevState,
      open: false,
      content: null,
    }));
  };

  const handleDesignDownloadModal = (e, pattern) => {
    e.stopPropagation();
    setDownloadModal((prevState) => ({
      ...prevState,
      open: true,
      content: pattern,
    }));
  };

  const handleModalClose = () => {
    setModalOpen({
      open: false,
      pattern_file: null,
      name: '',
      pattern_id: '',
      count: 0,
    });
  };

  const handleModalOpen = (pattern, action) => {
    const compCount = getComponentsinFile(pattern.pattern_file);
    const validationBody = (
      <Validation
        errors={patternErrors.get(pattern.id)}
        compCount={compCount}
        handleClose={() => setModalOpen({ ...modalOpen, open: false })}
      />
    );

    const dryRunComponent = (
      <DryRunComponent
        design={JSON.stringify({
          pattern_file: pattern.pattern_file,
          pattern_id: pattern.id,
        })}
        noOfElements={compCount}
        selectedContexts={selectedK8sContexts}
      />
    );
    setModalOpen({
      open: true,
      action: action,
      pattern_file: pattern.pattern_file,
      name: pattern.name,
      pattern_id: pattern.id,
      count: compCount,
      validationBody: validationBody,
      dryRunComponent: dryRunComponent,
    });
  };

  return (
    <div>
      {selectedPattern.show && (
        <DesignConfigurator
          pattern={selectedPattern.pattern}
          show={setSelectedPattern}
          onSubmit={handleSubmit}
        />
      )}
      {!selectedPattern.show && (
        <Grid container spacing={3}>
          {patterns.map((pattern) => (
            <PatternCardGridItem
              key={pattern.id}
              user={user}
              pattern={pattern}
              canPublishPattern={canPublishPattern}
              handleClone={() => handleClone(pattern.id, pattern.name)}
              handleDeploy={() => handleModalOpen(pattern, ACTIONS.DEPLOY)}
              handleUnDeploy={() => handleModalOpen(pattern, ACTIONS.UNDEPLOY)}
              handleVerify={(e) => handleVerify(e, pattern.pattern_file, pattern.id)}
              handlePublishModal={() => handlePublishModal(pattern)}
              handleUnpublishModal={(e) => handleUnpublishModal(e, pattern)()}
              handleInfoModal={() => handleInfoModal(pattern)}
              handleSubmit={handleSubmit}
              handleDownload={(e) => handleDesignDownloadModal(e, pattern)}
              setSelectedPatterns={setSelectedPattern}
            />
          ))}
        </Grid>
      )}

      {!selectedPattern.show && patterns.length === 0 && (
        <Paper className={classes.noPaper}>
          <div className={classes.noContainer}>
            <Typography align="center" color="textSecondary" className={classes.noText}>
              No Designs Found
            </Typography>
            <div>
              {/* <Button
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
                </Button> */}
            </div>
          </div>
        </Paper>
      )}
      {patterns.length ? (
        <div className={classes.pagination}>
          <Pagination
            count={pages}
            page={selectedPage + 1}
            onChange={(_, page) => setPage(page - 1)}
          />
        </div>
      ) : null}
      <ConfirmationMsg
        open={modalOpen.open}
        handleClose={handleModalClose}
        submit={{
          deploy: () => handleDeploy(modalOpen.pattern_file, modalOpen.pattern_id, modalOpen.name),
          unDeploy: () =>
            handleUnDeploy(modalOpen.pattern_file, modalOpen.pattern_id, modalOpen.name),
        }}
        title={modalOpen.name}
        componentCount={modalOpen.count}
        tab={modalOpen.action}
        dryRunComponent={modalOpen.dryRunComponent}
        validationBody={modalOpen.validationBody}
      />
      {canPublishPattern && publishModal.open && (
        <Modal
          open={true}
          schema={publishSchema.rjsfSchema}
          uiSchema={publishSchema.uiSchema}
          handleClose={handlePublishModalClose}
          aria-label="catalog publish"
          title={publishModal.pattern?.name}
          handleSubmit={handlePublish}
          showInfoIcon={{
            text: 'Upon submitting your catalog item, an approval flow will be initiated.',
            link: 'https://docs.meshery.io/concepts/catalog',
          }}
          submitBtnText="Submit for Approval"
          submitBtnIcon={<PublicIcon />}
        />
      )}
      <ExportModal
        downloadModal={downloadModal}
        handleDownloadDialogClose={handleDownloadDialogClose}
        handleDesignDownload={handleDownload}
      />
    </div>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

// @ts-ignore
export default connect(mapDispatchToProps)(withSnackbar(MesheryPatternGrid));
