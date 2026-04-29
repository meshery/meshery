import React from 'react';
import { ChipWrapper } from '../../connections/styles';
import { KeyValue } from '../../DataFormatter';

/**
 * Helper function to convert camelCase to Title Case with spaces
 * Example: "connectionID" -> "Connection ID", "k8sContextName" -> "K8s Context Name"
 */
export const humanizeFieldName = (fieldName) => {
  return fieldName
    .replace(/([A-Z])/g, ' $1') // Add space before capital letters
    .replace(/^./, (str) => str.toUpperCase()) // Capitalize first letter
    .trim();
};

/**
 * Formatter for connection-related fields (connectionID, k8sContextID, k8sContextName)
 * Displays as a clickable chip that links to the connections page with search filter
 *
 * @param {Object} props - Component props
 * @param {string} props.value - The value to display
 * @param {string} props.fieldName - The field name to humanize
 * @returns {React.ReactElement|null} Formatted connection field or null if no value
 */
export const ConnectionFieldFormatter = ({ value, fieldName }) => {
  if (!value) {
    return null;
  }

  const humanizedName = humanizeFieldName(fieldName);

  return (
    <KeyValue
      Key={humanizedName}
      Value={
        <ChipWrapper
          label={value}
          clickable
          component="a"
          href={`/management/connections?tab=connections&searchText=${encodeURIComponent(value)}`}
          target="_self"
          style={{
            marginBlock: '0.25rem',
          }}
        />
      }
    />
  );
};

/**
 * Property formatters for MeshSync-related metadata fields
 * These formatters handle connection, context, and deployment information
 */
export const MeshSyncPropertyFormatters = {
  connectionID: (value) => <ConnectionFieldFormatter value={value} fieldName="connectionID" />,
  k8sContextID: (value) => <ConnectionFieldFormatter value={value} fieldName="k8sContextID" />,
  k8sContextName: (value) => <ConnectionFieldFormatter value={value} fieldName="k8sContextName" />,
  meshsyncDeploymentMode: (value) => (
    <KeyValue Key={humanizeFieldName('meshsyncDeploymentMode')} Value={value} />
  ),
  operatorStatus: (value) => <KeyValue Key={humanizeFieldName('operatorStatus')} Value={value} />,
  brokerEndpoint: (value) => (
    <KeyValue
      Key={humanizeFieldName('brokerEndpoint')}
      Value={value}
      style={{
        fontFamily: 'monospace',
        fontSize: '0.85rem',
      }}
    />
  ),
};
