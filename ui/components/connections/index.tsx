import React, { useCallback, useMemo, useRef, useState } from 'react';
import { NoSsr } from '@sistent/sistent';
import { ErrorBoundary, AppBar } from '@sistent/sistent';
import Modal from '../shared/Modal/Modal';
import { ConnectionIconText, ConnectionTab, ConnectionTabs } from './styles';
import MeshSyncTable from './meshSync';
import ConnectionIcon from '../../assets/icons/Connection';
import MeshsyncIcon from '../../assets/icons/Meshsync';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import DefaultError from '../general/error-404/index';
import { useGetSchemaQuery } from '@/rtk-query/schema';
import CustomErrorFallback from '../shared/ErrorBoundary/ErrorBoundary';
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

  // Next.js's pages-router `router.query` and `router.push` get fresh
  // references on each render, which previously cascaded into a new
  // `updateUrlWithConnectionId` every commit. That prop is a dep of
  // ConnectionTable's `options` memo and of an in-table useEffect, so the
  // unstable reference forced both to invalidate every render, contributing
  // to the connections-page update-depth loop (React error #185). Mirror the
  // router state into refs so the callbacks below stay referentially stable.
  //
  // Assigning ref.current during render (rather than in a useEffect) is the
  // documented "latest value" pattern. Effects run child-first in the commit
  // phase, so deferring the sync to a parent useEffect would leave child
  // effects reading a stale `query`/`push` on the same commit they fire — and
  // ConnectionTable's expansion-sync effect does call `updateUrlParams`
  // through this ref. Writing in render keeps the ref in lockstep with the
  // values React just rendered with, before any child effect can read it.
  const routerStateRef = useRef({ query, pathname, push });
  routerStateRef.current = { query, pathname, push };

  const updateUrlParams = useCallback((params) => {
    const {
      query: currentQuery,
      pathname: currentPathname,
      push: currentPush,
    } = routerStateRef.current;
    const newQuery = { ...currentQuery, ...params };

    Object.keys(newQuery).forEach((key) => {
      if (newQuery[key] === undefined || newQuery[key] === '') {
        delete newQuery[key];
      }
    });

    currentPush({ pathname: currentPathname, query: newQuery }, undefined, { shallow: true });
  }, []);

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

  // Read latest selected connection id without re-creating the callback when
  // the URL changes — the dedupe guard would otherwise destabilize the prop.
  // Synced in render for the same reason as `routerStateRef`: child effects
  // run before parent effects in the commit phase.
  const connectionIdRef = useRef(connectionId);
  connectionIdRef.current = connectionId;

  // Update URL with connection ID
  const updateUrlWithConnectionId = useCallback(
    (id) => {
      if (id && id === connectionIdRef.current) {
        return;
      }

      updateUrlParams({ connectionId: id || undefined });
    },
    [updateUrlParams],
  );

  if (!isReady) return null;
  return (
    <NoSsr>
      {CAN(
        Keys.WorkspaceManagementViewConnections.id,
        Keys.WorkspaceManagementViewConnections.function,
      ) ? (
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

          {tab === 0 &&
            CAN(
              Keys.WorkspaceManagementViewConnections.id,
              Keys.WorkspaceManagementViewConnections.function,
            ) && (
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
