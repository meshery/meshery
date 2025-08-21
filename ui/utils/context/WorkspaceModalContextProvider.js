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
  createNewWorkspaceModalOpen: false,
  setCreateNewWorkspaceModalOpen: () => {},
});

const WorkspaceModalContextProvider = ({ children }) => {
  const [workspaceModal, setWorkspaceModal] = useState(false);
  const [selectedWorkspace, setSelectedWorkspace] = useState({ id: '', name: '' });
  const [multiSelectedContent, setMultiSelectedContent] = useState([]);
  const [createNewWorkspaceModalOpen, setCreateNewWorkspaceModalOpen] = useState(false);

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

        createNewWorkspaceModalOpen,
        setCreateNewWorkspaceModalOpen,
      }}
    >
      {children}
    </WorkspaceModalContext.Provider>
  );
};

export default WorkspaceModalContextProvider;
