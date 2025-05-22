import React, { useState } from 'react';

export const WorkspaceModalContext = React.createContext({
  open: false,
  openModal: () => {},
  closeModal: () => {},
  selectedWorkspace: { id: '', name: '' },
  setSelectedWorkspace: () => {},
  openModalWithDefault: () => {},
  multiSelectedContent: [],
  setMultiSelectedContent: () => {},
});

const WorkspaceModalContextProvider = ({ children }) => {
  const [workspaceModal, setWorkspaceModal] = useState(false);
  const [selectedWorkspace, setSelectedWorkspace] = useState({ id: '', name: '' });
  const [multiSelectedContent, setMultiSelectedContent] = useState([]);
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
        multiSelectedContent,
        setMultiSelectedContent,
      }}
    >
      {children}
    </WorkspaceModalContext.Provider>
  );
};

export default WorkspaceModalContextProvider;
