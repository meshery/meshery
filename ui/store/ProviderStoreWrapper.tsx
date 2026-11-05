import React from 'react';
import { Provider } from 'react-redux';
import { store } from '.';

const ProviderStoreWrapper = ({ children }) => {
  return <Provider store={store}>{children}</Provider>;
};

export default ProviderStoreWrapper;
