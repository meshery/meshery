import React, { useState } from 'react';
import { useGetSelectedOrganization } from '@/rtk-query/user';
import { useLazyGetWorkspacesQuery } from '@/rtk-query/workspace';

export const WorkspaceModalContext = React.createContext({
  open: false,
  openModal: () => {},
  closeModal: () => {},
  selectedWorkspace: { id: '', name: '' },
  setSelectedWorkspace: () => {},
  openModalWithDefault: () => {},
  multiSelectedContent: [],
  setMultiSelectedContent: () => {},
  createNewWorkspaceModalOpen: false,
  setCreateNewWorkspaceModalOpen: () => {},
  currentLoadedResource: {
    id: '',
    org: {
      id: '',
      name: '',
    },
    workspace: {
      id: '',
      name: '',
    },
  },
  onLoadResource: () => {},
});

const WorkspaceModalContextProvider = ({ children }) => {
  const { allOrganizations } = useGetSelectedOrganization();
  const [getWorkspaces] = useLazyGetWorkspacesQuery();
  const [workspaceModal, setWorkspaceModal] = useState(false);
  const [selectedWorkspace, setSelectedWorkspace] = useState({ id: '', name: '' });
  const [multiSelectedContent, setMultiSelectedContent] = useState([]);
  const [createNewWorkspaceModalOpen, setCreateNewWorkspaceModalOpen] = useState(false);
  // stores the context for currently loaded resource . sometime a user might be viewing a resource
  // from differnt org/workpace than the currently selected one
  const [currentLoadedResource, setCurrentLoadedResource] = useState({
    id: '',
    org: {
      id: '',
      name: '',
    },
    workspace: {
      id: '',
      name: '',
    },
  });

  const onLoadResource = async ({ id, workspaceId, orgId }) => {
    if (!workspaceId && !orgId) {
      setCurrentLoadedResource({ id, workspace: { id: '' }, org: { id: '' } });
      return;
    }

    try {
      const resource = {
        id,
        workspace: { id: workspaceId, name: 'Private Workspace' },
        org: { id: orgId, name: 'Private Org' },
      };
      const org = allOrganizations.find((o) => o.id == orgId);
      console.log('onload resource invoked', workspaceId, orgId, org, allOrganizations);
      if (org) {
        resource.org = org;
        const workspaces = await getWorkspaces({
          page: 0,
          pagesize: 'all',
          orgID: orgId,
        }).unwrap();
        const workspace = (workspaces?.workspaces ?? []).find((w) => w.id == workspaceId);
        if (workspace) {
          resource.workspace = workspace;
        }
      }
      console.log('onloadResource', workspaceId, orgId, resource);
      setCurrentLoadedResource(resource);
    } catch (e) {
      console.error('[onLoadResource] failed set orgWorkspace context', e);
    }
  };

  return (
    <WorkspaceModalContext.Provider
      value={{
        open: workspaceModal,
        openModal: () => setWorkspaceModal(true),
        closeModal: () => {
          setWorkspaceModal(false);
          setSelectedWorkspace({ id: '', name: '' });
          setMultiSelectedContent([]);
        },
        openModalWithDefault: () => {
          setWorkspaceModal(true);
          setSelectedWorkspace({ id: '', name: '' });
        },
        selectedWorkspace,
        setSelectedWorkspace: ({ id, name }) => {
          setMultiSelectedContent([]);
          setSelectedWorkspace({ id, name });
        },
        multiSelectedContent,
        setMultiSelectedContent,
        currentLoadedResource,
        onLoadResource,
        createNewWorkspaceModalOpen,
        setCreateNewWorkspaceModalOpen,
      }}
    >
      {children}
    </WorkspaceModalContext.Provider>
  );
};

export default WorkspaceModalContextProvider;
