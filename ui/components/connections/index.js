import React, { useRef, useState, useEffect } from 'react';
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
    selectedK8sContexts,
    k8sconfig,
    connectionMetadataState,
    meshsyncControllerState,
    router,
  } = props;
  const [_operatorState] = useState(operatorState || []);
  const _operatorStateRef = useRef(_operatorState);
  _operatorStateRef.current = _operatorState;

  const isUpdatingUrl = useRef(false);
  const isMounted = useRef(false);
  const isExternalUrlChange = useRef(false);

  const { query } = router;
  const tabParam = query.tab?.toLowerCase();
  const connectionId = query.connectionId;

  const [tab, setTab] = useState(() => {
    if (tabParam === 'meshsync') return 1;
    return 0; // Default to connections tab
  });

  // Mark component as mounted after first render
  useEffect(() => {
    isMounted.current = true;
    return () => {
      isMounted.current = false;
    };
  }, []);

  // Listen to URL changes from outside and prevent re-rendering loops
  useEffect(() => {
    if (isUpdatingUrl.current) return;

    isExternalUrlChange.current = true;

    const timeout = setTimeout(() => {
      isExternalUrlChange.current = false;
    }, 100);

    return () => clearTimeout(timeout);
  }, [router.asPath]);

  // Update URL when tab changes, with safeguard against infinite loops
  useEffect(() => {
    if (!isMounted.current) return;
    if (isExternalUrlChange.current) return;
    // Check for the tab
    const newTab = tab === 0 ? 'connections' : 'meshsync';
    const currentTab = query.tab;

    if (newTab !== currentTab && !isUpdatingUrl.current) {
      isUpdatingUrl.current = true;
      const newQuery = { ...query, tab: newTab };

      // Preserve the connectionId if it exists
      if (!connectionId) {
        delete newQuery.connectionId;
      }

      router
        .push(
          {
            pathname: router.pathname,
            query: newQuery,
          },
          undefined,
          { shallow: true },
        )
        .then(() => {
          setTimeout(() => {
            isUpdatingUrl.current = false;
          }, 100);
        });
    }
  }, [tab]);

  // Listen to URL changes and update the tab state if needed
  useEffect(() => {
    if (!isUpdatingUrl.current) {
      if (tabParam === 'meshsync' && tab !== 1) {
        setTab(1);
      } else if (tabParam === 'connections' && tab !== 0) {
        setTab(0);
      }
    }
  }, [tabParam]);

  const handleTabChange = (e, newTab) => {
    e.stopPropagation();
    setTab(newTab);
  };

  // Function to update URL with connection ID, with safeguards for page load
  const updateUrlWithConnectionId = (id) => {
    if (isUpdatingUrl.current || isExternalUrlChange.current) return;

    // If we're dealing with the initial connection from URL, don't update again
    if (!isMounted.current && id === connectionId) return;

    // Skip update if the ID is already in the URL, but DON'T skip when clearing the ID (id is empty)
    if (id && id === connectionId) return;

    isUpdatingUrl.current = true;
    const newQuery = { ...query };

    if (id) {
      newQuery.connectionId = id;
    } else {
      delete newQuery.connectionId;
    }

    router
      .push(
        {
          pathname: router.pathname,
          query: newQuery,
        },
        undefined,
        { shallow: true },
      )
      .then(() => {
        // Reset the flag after URL update completes
        setTimeout(() => {
          isUpdatingUrl.current = false;
        }, 100);
      });
  };

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
              selectedK8sContexts={selectedK8sContexts}
              k8sconfig={k8sconfig}
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
export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withRouter(ConnectionManagementPageWithErrorBoundary));
