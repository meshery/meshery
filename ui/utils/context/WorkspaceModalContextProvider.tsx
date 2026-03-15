import React, { useState } from 'react';
import { useGetSelectedOrganization } from '@/rtk-query/user';
import { useLazyGetWorkspacesQuery } from '@/rtk-query/workspace';

type WorkspaceRef = { id: string; name: string };
type OrgRef = { id: string; name: string };
type CurrentLoadedResource = { id: string; org: OrgRef; workspace: WorkspaceRef };

type WorkspaceModalContextType = {
  open: boolean;
  openModal: () => void;
  closeModal: () => void;
  selectedWorkspace: WorkspaceRef;
  setSelectedWorkspace: (_value: WorkspaceRef) => void;
  openModalWithDefault: () => void;
  multiSelectedContent: any[];
  setMultiSelectedContent: React.Dispatch<React.SetStateAction<any[]>>;
  createNewWorkspaceModalOpen: boolean;
  setCreateNewWorkspaceModalOpen: React.Dispatch<React.SetStateAction<boolean>>;
  currentLoadedResource: CurrentLoadedResource;
  onLoadResource: (_args: { id: string; workspaceId?: string; orgId?: string }) => Promise<void>;
};

export const WorkspaceModalContext = React.createContext<WorkspaceModalContextType>({
  open: false,
  openModal: () => undefined,
  closeModal: () => undefined,
  selectedWorkspace: { id: '', name: '' },
  setSelectedWorkspace: () => undefined,
  openModalWithDefault: () => undefined,
  multiSelectedContent: [],
  setMultiSelectedContent: (() => undefined) as unknown as React.Dispatch<
    React.SetStateAction<any[]>
  >,
  createNewWorkspaceModalOpen: false,
  setCreateNewWorkspaceModalOpen: (() => undefined) as unknown as React.Dispatch<
    React.SetStateAction<boolean>
  >,
  currentLoadedResource: {
    id: '',
    org: { id: '', name: '' },
    workspace: { id: '', name: '' },
  },
  onLoadResource: async () => undefined,
});

const WorkspaceModalContextProvider = ({ children }: React.PropsWithChildren) => {
  const { allOrganizations } = useGetSelectedOrganization();
  const [getWorkspaces] = useLazyGetWorkspacesQuery();
  const [workspaceModal, setWorkspaceModal] = useState(false);
  const [selectedWorkspace, setSelectedWorkspace] = useState({ id: '', name: '' });
  const [multiSelectedContent, setMultiSelectedContent] = useState<any[]>([]);
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
      setCurrentLoadedResource({
        id,
        workspace: { id: '', name: '' },
        org: { id: '', name: '' },
      });
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
        const workspaces = (await getWorkspaces({
          page: 0,
          pagesize: 'all',
          orgID: orgId,
        }).unwrap()) as { workspaces?: Array<{ id: string; name: string }> };
        const workspace = (workspaces?.workspaces ?? []).find((w) => w.id == workspaceId);
        if (workspace) {
          resource.workspace = workspace;
        }
      }
      console.log('onloadResource', workspaceId, orgId, resource);
      setCurrentLoadedResource(resource);
    } catch (e) {
      console.log('[onLoadResource] failed set orgWorkspace context', e);
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
