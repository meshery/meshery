import React from 'react';
import FinishFlagIcon from '../../../assets/icons/FinishFlagIcon';
import AnalyticsIcon from '@mui/icons-material/Analytics';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';

// Import step components
import { ImportStep } from './steps/ImportStep';
import { DiscoveryStep } from './steps/DiscoveryStep';
import { SelectionStep } from './steps/SelectionStep';
import { ResultsStep } from './steps/ResultsStep';

/**
 * Generic wizard steps that work for all connection types
 */
const GENERIC_WIZARD_STEPS = [
  'Import Connection',
  'Discovery & Analysis',
  'Configuration & Selection',
  'Results',
];

/**
 * Generic wizard icons
 */
const GENERIC_WIZARD_ICONS = {
  1: <CloudUploadIcon width={26} height={26} />,
  2: <AnalyticsIcon width={26} height={26} />,
  3: <CheckCircleIcon width={26} height={26} />,
  4: <FinishFlagIcon width="2rem" />,
};

/**
 * Generic wizard content - step components
 */
const GENERIC_WIZARD_CONTENT = [
  <ImportStep key="import" />,
  <DiscoveryStep key="discovery" />,
  <SelectionStep key="selection" />,
  <ResultsStep key="results" />,
];

/**
 * Connection type specific configurations
 */
const CONNECTION_TYPE_CONFIGS = {
  kubernetes: {
    steps: ['Import Kubeconfig', 'Cluster Discovery', 'Cluster Selection', 'Connection Results'],
    icons: {
      1: <CloudUploadIcon width={26} height={26} />,
      2: <AnalyticsIcon width={26} height={26} />,
      3: <CheckCircleIcon width={26} height={26} />,
      4: <FinishFlagIcon width="2rem" />,
    },
    content: [
      <ImportStep key="kubernetes-import" connectionType="kubernetes" />,
      <DiscoveryStep key="kubernetes-discovery" connectionType="kubernetes" />,
      <SelectionStep key="kubernetes-selection" connectionType="kubernetes" />,
      <ResultsStep key="kubernetes-results" connectionType="kubernetes" />,
    ],
  },
};

/**
 * Get wizard steps for a specific connection type
 */
export const getWizardSteps = (connectionType) => {
  const config = CONNECTION_TYPE_CONFIGS[connectionType];
  return config ? config.steps : GENERIC_WIZARD_STEPS;
};

/**
 * Get wizard icons for a specific connection type
 */
export const getWizardIcons = (connectionType) => {
  const config = CONNECTION_TYPE_CONFIGS[connectionType];
  return config ? config.icons : GENERIC_WIZARD_ICONS;
};

/**
 * Get wizard content components for a specific connection type
 */
export const getWizardContent = (connectionType) => {
  const config = CONNECTION_TYPE_CONFIGS[connectionType];
  return config ? config.content : GENERIC_WIZARD_CONTENT;
};


export const isConnectionTypeSupported = (connectionType) => {
  return Object.prototype.hasOwnProperty.call(CONNECTION_TYPE_CONFIGS, connectionType);
};
