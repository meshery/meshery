import React from 'react';

/** Design or view row in the workspace explorer multi-select toolbar. */
export type WorkspaceMultiSelectItem = {
  id: string;
};

/** Default value; keep in sync with WorkspaceModalContextProvider state shape. */
export const workspaceModalContextDefault = {
  open: false,
  openModal: () => {},
  closeModal: () => {},
  selectedWorkspace: { id: '', name: '' },
  setSelectedWorkspace: () => {},
  openModalWithDefault: () => {},
  multiSelectedContent: [] as WorkspaceMultiSelectItem[],
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
};

export type WorkspaceModalContextValue = typeof workspaceModalContextDefault;

/**
 * App-level workspace explorer modal. Defined in its own module so every consumer
 * and WorkspaceModalContextProvider share one React context instance (avoids duplicate
 * createContext copies when the provider is bundled separately on kanvas builds).
 */
export const WorkspaceModalContext = React.createContext<WorkspaceModalContextValue>(
  workspaceModalContextDefault,
);
