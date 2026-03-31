import React from 'react';
import { createRoot } from 'react-dom/client';
import ExtensionComponent from './components/ExtensionComponent/ExtensionComponent';
import "./styles.css";

const root = createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <ExtensionComponent />
  </React.StrictMode>
);
