import React, { useState } from 'react';

// Create a context to hold the dynamic component
const DynamicComponentContext = React.createContext();

// Custom hook to use the dynamic component
export const useDynamicComponent = () => {
  const context = React.useContext(DynamicComponentContext);
  if (!context) {
    throw new Error('useDynamicComponent must be used within a DynamicComponentProvider');
  }
  return context;
};

// Provider component to set the dynamic component
export const DynamicComponentProvider = ({ children }) => {
  const [DynamicComponent, setComponent] = useState(null);
  return (
    <DynamicComponentContext.Provider value={{ setComponent, DynamicComponent }}>
      {children}
    </DynamicComponentContext.Provider>
  );
};
