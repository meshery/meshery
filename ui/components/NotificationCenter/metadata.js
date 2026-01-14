import React from 'react';
import { FormatStructuredData, SectionBody, reorderObjectProperties, KeyValue } from '../DataFormatter';
import { isEmptyAtAllDepths } from '../../utils/objects';
import { canTruncateDescription } from './notification';
import { DeploymentSummaryFormatter } from '../DesignLifeCycle/DeploymentSummary';
import { EVENT_TYPE, eventDetailFormatterKey } from './constants';
import { TitleLink, DataToFileLink, EmptyState } from './formatters/common';
import { ErrorMetadataFormatter } from './formatters/error';
import { DryRunResponse, SchemaValidationFormatter } from './formatters/pattern_dryrun';
import { ModelImportMessages, ModelImportedSection } from './formatters/model_registration';
import { RelationshipEvaluationEventFormatter } from './formatters/relationship_evaluation';
import { useTheme, DownloadIcon, InfoIcon } from '@sistent/sistent';
import _ from 'lodash';
import { ChipWrapper } from '../connections/styles';
import { AcademyEventsFormatter } from './formatters/academy_events';

const DesignFormatter = ({ value }) => {
  const theme = useTheme();
  const { name, id } = value;

  return (
    <TitleLink
      href={'/extension/meshmap?mode=design&design=' + encodeURIComponent(id)}
      style={{
        color: theme.palette.text.default,
        fontWeight: 'normal',
        textDecoration: 'none',
      }}
      target="_self"
    >
      Saved design {name}
    </TitleLink>
  );
};

const ShortDescriptionFormatter = ({ value }) => {
  const theme = useTheme();
  return (
    <SectionBody
      body={value}
      style={{ marginBlock: '0.5rem', color: theme.palette.text.default, fontWeight: 'normal' }}
    />
  );
};

/**
 * Helper function to convert camelCase to Title Case with spaces
 * Example: "connectionID" -> "Connection ID", "k8sContextName" -> "K8s Context Name"
 */
const humanizeFieldName = (fieldName) => {
  return fieldName
    .replace(/([A-Z])/g, ' $1') // Add space before capital letters
    .replace(/^./, (str) => str.toUpperCase()) // Capitalize first letter
    .trim();
};

/**
 * Formatter for connection-related fields (connectionID, k8sContextID, k8sContextName)
 * Displays as a clickable chip that links to the connections page with search filter
 */
const ConnectionFieldFormatter = ({ value, fieldName }) => {
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

export const PropertyFormatters = {
  trace: (value) => <DataToFileLink data={value} />,
  ShortDescription: (value) => <ShortDescriptionFormatter value={value} />,
  design: (value) => <DesignFormatter value={value} />,
  connectionName: (value) => {
    return (
      <ChipWrapper
        label={value}
        clickable
        component="a"
        href={`/management/connections?tab=connections&searchText=${encodeURIComponent(value)}`}
        target="_self"
      />
    );
  },
  // MeshSync event field formatters with humanized labels
  connectionID: (value) => <ConnectionFieldFormatter value={value} fieldName="connectionID" />,
  k8sContextID: (value) => <ConnectionFieldFormatter value={value} fieldName="k8sContextID" />,
  k8sContextName: (value) => <ConnectionFieldFormatter value={value} fieldName="k8sContextName" />,
  meshsyncDeploymentMode: (value) => (
    <KeyValue
      Key={humanizeFieldName('meshsyncDeploymentMode')}
      Value={value}
    />
  ),
  operatorStatus: (value) => (
    <KeyValue
      Key={humanizeFieldName('operatorStatus')}
      Value={value}
    />
  ),
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
  error: (value) => <ErrorMetadataFormatter metadata={value} event={event} />,
  dryRunResponse: (value) => <DryRunResponse response={value} />,
  ModelImportMessage: (value) => value && <ModelImportMessages message={value} />,
  ModelDetails: (value) => value && <ModelImportedSection modelDetails={value} />,
  history_title: () => null,
};

export const LinkFormatters = {
  doclink: (value) => {
    return (
      <TitleLink href={value} style={{ textAlign: 'end', color: 'inherit' }}>
        Doc
      </TitleLink>
    );
  },
};
export const PropertyLinkFormatters = {
  doc: (value) => ({
    label: 'Doc',
    href: value,
  }),
  DownloadLink: (value) => ({
    label: 'Download File',
    href: '/api/system/fileDownload?file=' + encodeURIComponent(value),
    icon: DownloadIcon,
  }),
  ViewLink: (value) => ({
    label: 'Get Logs',
    href: '/api/system/fileView?file=' + encodeURIComponent(value),
    icon: InfoIcon,
  }),
};

const linkOrder = ['doclink'];

const EventTypeFormatters = {
  [eventDetailFormatterKey(EVENT_TYPE.DEPLOY_DESIGN)]: DeploymentSummaryFormatter,
  [eventDetailFormatterKey(EVENT_TYPE.UNDEPLOY_DESIGN)]: DeploymentSummaryFormatter,
  [eventDetailFormatterKey(EVENT_TYPE.EVALUATE_DESIGN)]: RelationshipEvaluationEventFormatter,
  [eventDetailFormatterKey(EVENT_TYPE.ACADEMY_QUIZ_EVALUATION)]: AcademyEventsFormatter,
  [eventDetailFormatterKey(EVENT_TYPE.VALIDATE_DESIGN)]: SchemaValidationFormatter,
  // [eventDetailFormatterKey(EVENT_TYPE.REGISTRANT_SUMMARY)]: RegistrantSummaryFormatter,
};

export const FormattedMetadata = ({ event }) => {
  const theme = useTheme();
  if (EventTypeFormatters[eventDetailFormatterKey(event)]) {
    const Formatter = EventTypeFormatters[eventDetailFormatterKey(event)];
    return <Formatter event={event} />;
  }

  if (!event || !event.metadata || isEmptyAtAllDepths(event.metadata)) {
    return <EmptyState event={event} />;
  }

  const metadata = {
    ..._.omit(event.metadata, [...linkOrder, 'id', 'kind', 'ViewLink', 'DownloadLink']),
    ShortDescription:
      event.metadata.error || !canTruncateDescription(event.description || '')
        ? null
        : event.description,
  };
  const order = [
    'doclink',
    'ShortDescription',
    'connectionName',
    'connectionID',
    'k8sContextName',
    'k8sContextID',
    'meshsyncDeploymentMode',
    'operatorStatus',
    'brokerEndpoint',
    'LongDescription',
    'Summary',
    'SuggestedRemediation',
    'ModelImportMessage',
    'ModelDetails',
  ];
  const hasImportedModelName = !!metadata.ImportedModelName;

  const orderedMetadata = hasImportedModelName
    ? reorderObjectProperties({ ...metadata, ShortDescription: null }, order) // Exclude ShortDescription
    : reorderObjectProperties(metadata, order);
  return (
    <FormatStructuredData
      propertyFormatters={PropertyFormatters}
      data={orderedMetadata}
      order={order}
      style={{
        fontWeight: 'normal',
        color: theme.palette.text.default,
      }}
      event={event}
    />
  );
};

export const FormattedLinkMetadata = ({ event }) => {
  const filteredMetadata = _.pick(event.metadata, linkOrder);
  return <FormatStructuredData propertyFormatters={LinkFormatters} data={filteredMetadata} />;
};
