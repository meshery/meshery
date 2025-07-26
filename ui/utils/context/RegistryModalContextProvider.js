import React, { useState } from 'react';

export const RegistryModalContext = React.createContext({
  open: false,
  openModal: () => {},
  closeModal: () => {},
  selectedView: '',
  setSelectedView: () => {},
});

const RegistryModalContextProvider = ({ children }) => {
  const [registryModal, setRegistryModal] = useState(false);
  const [selectedView, setSelectedView] = useState('Models');

  return (
    <RegistryModalContext.Provider
      value={{
        open: registryModal,
        openModal: () => setRegistryModal(true),
        closeModal: () => {
          setRegistryModal(false);
        },
        selectedView,
        setSelectedView: (view) => {
          setSelectedView(view);
        },
      }}
    >
      {children}
    </RegistryModalContext.Provider>
  );
};

export default RegistryModalContextProvider;
