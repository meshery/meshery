import React, { useState } from 'react';

export const WorkspaceModalContext = React.createContext({
  open: false,
  openModal: () => {},
  closeModal: () => {},
  selectedWorkspace: { id: '', name: '' },
  setSelectedWorkspace: () => {},
  openModalWithDefault: () => {},
});

const WorkspaceModalContextProvider = ({ children }) => {
  const [workspaceModal, setWorkspaceModal] = useState(false);
  const [selectedWorkspace, setSelectedWorkspace] = useState({ id: '', name: '' });
  return (
    <WorkspaceModalContext.Provider
      value={{
        open: workspaceModal,
        openModal: () => setWorkspaceModal(true),
        closeModal: () => {
          setWorkspaceModal(false);
        },
        openModalWithDefault: () => {
          setWorkspaceModal(true);
          setSelectedWorkspace({ id: '', name: '' });
        },
        selectedWorkspace,
        setSelectedWorkspace,
      }}
    >
      {children}
    </WorkspaceModalContext.Provider>
  );
};

export default WorkspaceModalContextProvider;
