import React, { useRef, useState } from 'react';
import { NoSsr } from '@layer5/sistent';
import { ErrorBoundary, AppBar } from '@layer5/sistent';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateProgress } from '../../lib/store';
import Modal from '../Modal';
import { ConnectionIconText, ConnectionTab, ConnectionTabs } from './styles';
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
    updateProgress,
    operatorState,
    connectionMetadataState,
    meshsyncControllerState,
    router,
  } = props;
  const [_operatorState] = useState(operatorState || []);
  const _operatorStateRef = useRef(_operatorState);
  _operatorStateRef.current = _operatorState;

  const { query, pathname, push, isReady } = router;
  const tabParam = query.tab?.toLowerCase();
  const connectionId = query.connectionId;

  const tab = tabParam === 'meshsync' ? 1 : 0;

  const updateUrlParams = (params) => {
    const newQuery = { ...query, ...params };

    Object.keys(newQuery).forEach((key) => {
      if (newQuery[key] === undefined || newQuery[key] === '') {
        delete newQuery[key];
      }
    });

    push({ pathname, query: newQuery }, undefined, { shallow: true });
  };

  // Handle tab change and update URL
  const handleTabChange = (e, newTab) => {
    e.stopPropagation();

    if (newTab !== tab) {
      updateUrlParams({
        tab: newTab === 0 ? 'connections' : 'meshsync',
        connectionId: undefined, // Clear the connection ID when switching tabs
      });
    }
  };
  // Update URL with connection ID
  const updateUrlWithConnectionId = (id) => {
    if (id && id === connectionId) return;

    updateUrlParams({ connectionId: id || undefined });
  };

  if (!isReady) return null;
  return (
    <NoSsr>
      {CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject) ? (
        <>
          <AppBar position="static" color="default" style={{ marginBottom: '3rem' }}>
            <ConnectionTabs
              value={tab}
              onChange={handleTabChange}
              indicatorColor="primary"
              textColor="primary"
              variant="fullWidth"
              sx={{
                height: '10%',
              }}
            >
              <ConnectionTab
                label={
                  <ConnectionIconText>
                    <span style={{ marginRight: '0.3rem' }}>Connections</span>
                    <ConnectionIcon width="20" height="20" />
                  </ConnectionIconText>
                }
              />
              <ConnectionTab
                label={
                  <ConnectionIconText>
                    <span style={{ marginRight: '0.3rem' }}>MeshSync</span>
                    <MeshsyncIcon width="20" height="20" />
                  </ConnectionIconText>
                }
              />
            </ConnectionTabs>
          </AppBar>

          {tab === 0 && CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject) && (
            <ConnectionTable
              meshsyncControllerState={meshsyncControllerState}
              connectionMetadataState={connectionMetadataState}
              selectedConnectionId={connectionId}
              updateUrlWithConnectionId={updateUrlWithConnectionId}
            />
          )}
          {tab === 1 && (
            <MeshSyncTable
              updateProgress={updateProgress}
              selectedResourceId={connectionId}
              updateUrlWithResourceId={updateUrlWithConnectionId}
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
  const operatorState = state.get('operatorState');
  const connectionMetadataState = state.get('connectionMetadataState');
  const meshsyncControllerState = state.get('controllerState');

  return {
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
export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withRouter(ConnectionManagementPageWithErrorBoundary));
