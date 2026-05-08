import React, { useCallback, useMemo, useState } from 'react';
import { NoSsr } from '@sistent/sistent';
import { ErrorBoundary, AppBar } from '@sistent/sistent';
import Modal from '../General/Modals/Modal';
import { ConnectionIconText, ConnectionTab, ConnectionTabs } from './styles';
import MeshSyncTable from './meshSync';
import ConnectionIcon from '../../assets/icons/Connection';
import MeshsyncIcon from '../../assets/icons/Meshsync';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import DefaultError from '../General/error-404/index';
import { useGetSchemaQuery } from '@/rtk-query/schema';
import CustomErrorFallback from '../General/ErrorBoundary';
import ConnectionTable from './ConnectionTable';
import { useRouter } from 'next/router';

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
function Connections() {
  const router = useRouter();
  const { query, pathname, push, isReady } = router;
  const tabParam = typeof query.tab === 'string' ? query.tab.toLowerCase() : undefined;
  const connectionId = typeof query.connectionId === 'string' ? query.connectionId : undefined;

  const tab = useMemo(() => (tabParam === 'meshsync' ? 1 : 0), [tabParam]);

  const updateUrlParams = useCallback(
    (params) => {
      const newQuery = { ...query, ...params };

      Object.keys(newQuery).forEach((key) => {
        if (newQuery[key] === undefined || newQuery[key] === '') {
          delete newQuery[key];
        }
      });

      push({ pathname, query: newQuery }, undefined, { shallow: true });
    },
    [pathname, push, query],
  );

  // Handle tab change and update URL
  const handleTabChange = useCallback(
    (event, newTab) => {
      event.stopPropagation();

      if (newTab !== tab) {
        updateUrlParams({
          tab: newTab === 0 ? 'connections' : 'meshsync',
          connectionId: undefined,
        });
      }
    },
    [tab, updateUrlParams],
  );
  // Update URL with connection ID
  const updateUrlWithConnectionId = useCallback(
    (id) => {
      if (id && id === connectionId) {
        return;
      }

      updateUrlParams({ connectionId: id || undefined });
    },
    [connectionId, updateUrlParams],
  );

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
              selectedConnectionId={connectionId}
              updateUrlWithConnectionId={updateUrlWithConnectionId}
            />
          )}
          {tab === 1 && (
            <MeshSyncTable
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

const ConnectionManagementPageWithErrorBoundary = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary customFallback={CustomErrorFallback}>
        <ConnectionManagementPage {...props} />
      </ErrorBoundary>
    </NoSsr>
  );
};

export default ConnectionManagementPageWithErrorBoundary;
