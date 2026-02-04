import React, { useState, type ReactNode } from 'react';

type RegistryModalParams = {
  tab?: string;
  searchText?: string;
  selectedItemUUID?: string;
};

type RegistryModalContextValue = {
  open: boolean;
  openModal: () => void;
  closeModal: () => void;
  openModalWithParams: (_params?: RegistryModalParams) => void;
  selectedView: string;
  setSelectedView: (_view: string) => void;
  searchText: string;
  setSearchText: (_text: string) => void;
  selectedItemUUID: string;
  setSelectedItemUUID: (_id: string) => void;
};

export const RegistryModalContext = React.createContext<RegistryModalContextValue>({
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

type RegistryModalContextProviderProps = {
  children: ReactNode;
};

const RegistryModalContextProvider = ({ children }: RegistryModalContextProviderProps) => {
  const [registryModal, setRegistryModal] = useState(false);
  const [selectedView, setSelectedView] = useState('Models');
  const [searchText, setSearchText] = useState('');
  const [selectedItemUUID, setSelectedItemUUID] = useState('');

  const openModalWithParams = (params: RegistryModalParams = {}) => {
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
        setSelectedView: (view: string) => {
          setSelectedView(view);
        },
        searchText,
        setSearchText: (text: string) => setSearchText(text),
        selectedItemUUID,
        setSelectedItemUUID: (id: string) => setSelectedItemUUID(id),
      }}
    >
      {children}
    </RegistryModalContext.Provider>
  );
};

export default RegistryModalContextProvider;
