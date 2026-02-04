import React, { useState, type ReactNode } from 'react';

type DynamicComponentType = React.ComponentType<any> | null;

type DynamicComponentContextValue = {
  DynamicComponent: DynamicComponentType;
  setComponent: (component: DynamicComponentType) => void;
};

// Create a context to hold the dynamic component
const DynamicComponentContext = React.createContext<DynamicComponentContextValue | undefined>(
  undefined,
);

// Custom hook to use the dynamic component
export const useDynamicComponent = (): DynamicComponentContextValue => {
  const context = React.useContext(DynamicComponentContext);
  if (!context) {
    throw new Error('useDynamicComponent must be used within a DynamicComponentProvider');
  }
  return context;
};

type DynamicComponentProviderProps = {
  children: ReactNode;
};

// Provider component to set the dynamic component
export const DynamicComponentProvider = ({ children }: DynamicComponentProviderProps) => {
  const [DynamicComponent, setComponent] = useState<DynamicComponentType>(null);
  return (
    <DynamicComponentContext.Provider value={{ setComponent, DynamicComponent }}>
      {children}
    </DynamicComponentContext.Provider>
  );
};
