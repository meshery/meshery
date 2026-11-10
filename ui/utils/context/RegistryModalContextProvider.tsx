import React, { useState } from 'react';

export const RegistryModalContext = React.createContext({
  open: false,
  openModal: () => {},
  closeModal: () => {},
  openModalWithParams: () => {},
  selectedView: '',
  setSelectedView: () => {},
  searchText: '',
  setSearchText: () => {},
  selectedItemUUID: '',
  setSelectedItemUUID: () => {},
});

const RegistryModalContextProvider = ({ children }) => {
  const [registryModal, setRegistryModal] = useState(false);
  const [selectedView, setSelectedView] = useState('Models');
  const [searchText, setSearchText] = useState('');
  const [selectedItemUUID, setSelectedItemUUID] = useState('');

  const openModalWithParams = (params = {}) => {
    if (params.tab) setSelectedView(params.tab);
    if (params.searchText) setSearchText(params.searchText);
    if (params.selectedItemUUID) setSelectedItemUUID(params.selectedItemUUID);
    setRegistryModal(true);
  };

  const closeModal = () => {
    setRegistryModal(false);
    // Reset parameters when closing
    setSearchText('');
    setSelectedItemUUID('');
  };

  return (
    <RegistryModalContext.Provider
      value={{
        open: registryModal,
        openModal: () => setRegistryModal(true),
        closeModal,
        openModalWithParams,
        selectedView,
        setSelectedView: (view) => {
          setSelectedView(view);
        },
        searchText,
        setSearchText,
        selectedItemUUID,
        setSelectedItemUUID,
      }}
    >
      {children}
    </RegistryModalContext.Provider>
  );
};

export default RegistryModalContextProvider;
