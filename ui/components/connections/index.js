import React, { useRef, useState } from 'react';
import { NoSsr, AppBar, Tabs, Tab } from '@material-ui/core';
import { ErrorBoundary } from '@layer5/sistent';
import { withStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateProgress } from '../../lib/store';
import Modal from '../Modal';
import styles from './styles';
import MeshSyncTable from './meshSync';
import ConnectionIcon from '../../assets/icons/Connection';
import MeshsyncIcon from '../../assets/icons/Meshsync';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import DefaultError from '../General/error-404/index';
import { useGetSchemaQuery } from '@/rtk-query/schema';
import { withRouter } from 'next/router';
import CustomErrorFallback from '../General/ErrorBoundary';
import ConnectionTable from './ConnectionTable';

/**
 * Parent Component for Connection Component
 *
 * @important
 * - Keep the component's responsibilities focused on connection management. Avoid adding unrelated functionality and state.
 */

function ConnectionManagementPage(props) {
  const [createConnectionModal, setCreateConnectionModal] = useState({
    open: false,
  });

  const { data: schemaResponse } = useGetSchemaQuery({
    schemaName: 'helmRepo',
  });

  const createConnection = schemaResponse ?? {};

  const handleCreateConnectionModalOpen = () => {
    setCreateConnectionModal({ open: true });
  };

  const handleCreateConnectionModalClose = () => {
    setCreateConnectionModal({ open: false });
  };

  const handleCreateConnectionSubmit = () => {};

  return (
    <>
      <Connections
        createConnectionModal={createConnectionModal}
        onOpenCreateConnectionModal={handleCreateConnectionModalOpen}
        onCloseCreateConnectionModal={handleCreateConnectionModalClose}
        {...props}
      />
      {createConnectionModal.open && (
        <Modal
          open={true}
          schema={createConnection.rjsfSchema}
          uiSchema={createConnection.uiSchema}
          handleClose={handleCreateConnectionModalClose}
          handleSubmit={handleCreateConnectionSubmit}
          title="Connect Helm Repository"
          submitBtnText="Connect"
        />
      )}
    </>
  );
}
function Connections(props) {
  const {
    classes,
    updateProgress,
    /*onOpenCreateConnectionModal,*/ operatorState,
    selectedK8sContexts,
    k8sconfig,
    connectionMetadataState,
    meshsyncControllerState,
  } = props;
  const [_operatorState] = useState(operatorState || []);
  const [tab, setTab] = useState(0);
  const _operatorStateRef = useRef(_operatorState);
  _operatorStateRef.current = _operatorState;

  return (
    <NoSsr>
      {CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject) ? (
        <>
          <AppBar position="static" color="default" className={classes.appBar}>
            <Tabs
              value={tab}
              className={classes.tabs}
              onChange={(e, newTab) => {
                e.stopPropagation();
                setTab(newTab);
              }}
              indicatorColor="primary"
              textColor="primary"
              variant="fullWidth"
              sx={{
                height: '10%',
              }}
            >
              <Tab
                className={classes.tab}
                label={
                  <div className={classes.iconText}>
                    <span style={{ marginRight: '0.3rem' }}>Connections</span>
                    <ConnectionIcon width="20" height="20" />
                  </div>
                }
              />
              <Tab
                className={classes.tab}
                label={
                  <div className={classes.iconText}>
                    <span style={{ marginRight: '0.3rem' }}>MeshSync</span>
                    <MeshsyncIcon width="20" height="20" />
                  </div>
                }
              />
            </Tabs>
          </AppBar>

          {tab === 0 && CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject) && (
            <ConnectionTable
              meshsyncControllerState={meshsyncControllerState}
              connectionMetadataState={connectionMetadataState}
            />
          )}
          {tab === 1 && (
            <MeshSyncTable
              classes={classes}
              updateProgress={updateProgress}
              selectedK8sContexts={selectedK8sContexts}
              k8sconfig={k8sconfig}
            />
          )}
        </>
      ) : (
        <DefaultError />
      )}
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const selectedK8sContexts = state.get('selectedK8sContexts');
  const operatorState = state.get('operatorState');
  const connectionMetadataState = state.get('connectionMetadataState');
  const meshsyncControllerState = state.get('controllerState');

  return {
    k8sconfig,
    selectedK8sContexts,
    operatorState,
    connectionMetadataState,
    meshsyncControllerState,
  };
};

const ConnectionManagementPageWithErrorBoundary = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary customFallback={CustomErrorFallback}>
        <ConnectionManagementPage {...props} />
      </ErrorBoundary>
    </NoSsr>
  );
};

// @ts-ignore
export default withStyles(styles)(
  connect(
    mapStateToProps,
    mapDispatchToProps,
  )(withRouter(ConnectionManagementPageWithErrorBoundary)),
);
