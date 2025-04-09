import React from 'react';
import { FormatStructuredData, SectionBody, reorderObjectProperties } from '../DataFormatter';
import { isEmptyAtAllDepths } from '../../utils/objects';
import { canTruncateDescription } from './notification';
import { DeploymentSummaryFormatter } from '../DesignLifeCycle/DeploymentSummary';
import { EVENT_TYPE, eventDetailFormatterKey } from './constants';
import { TitleLink, DataToFileLink, EmptyState } from './formatters/common';
import { ErrorMetadataFormatter } from './formatters/error';
import { DryRunResponse } from './formatters/pattern_dryrun';
import { ModelImportMessages, ModelImportedSection } from './formatters/model_registration';
import { RelationshipEvaluationEventFormatter } from './formatters/relationship_evaluation';
import { useTheme } from '@layer5/sistent';
import _ from 'lodash';
import Chip from '@mui/material/Chip';
import { RegistrantSummaryFormatter } from './formatters/RegistrantSummaryFormatter';

export const PropertyFormatters = {
  //trace can be very large, so we need to convert it to a file
  trace: (value) => <DataToFileLink data={value} />,
  ShortDescription: (value) => {
    const theme = useTheme();
    return (
      <SectionBody
        body={value}
        style={{ marginBlock: '0.5rem', color: theme.palette.text.default, fontWeight: 'normal' }}
      />
    );
  },
  designId: (value) => {
    const theme = useTheme();
    // console.log('Here is the value', value);
    const parts = value.split('+');
    const designName = parts[0];
    const designId = parts[1];

    return (
      <TitleLink
        href={'/extension/meshmap?mode=design&design=' + encodeURIComponent(designId)}
        style={{
          color: theme.palette.text.default,
          fontWeight: 'normal',
          textDecoration: 'none',
        }}
        target="_self"
      >
        Saved design {designName}
      </TitleLink>
    );
  },
  connectionName: (value) => {
    const theme = useTheme();

    return (
      <Chip
        label={value}
        clickable
        component="a"
        href={`/management/connections?tab=connections&searchText=${value}`}
        style={{
          backgroundColor: theme.palette.primary.light,
          color: theme.palette.primary.contrastText,
          textDecoration: 'none',
          fontWeight: 'bold',
        }}
        target="_self"
      />
    );
  },
  error: (value) => <ErrorMetadataFormatter metadata={value} event={event} />,
  dryRunResponse: (value) => <DryRunResponse response={value} />,
  ModelImportMessage: (value) => value && <ModelImportMessages message={value} />,
  ModelDetails: (value) => value && <ModelImportedSection modelDetails={value} />,
};

export const PropertyLinkFormatters = {
  doc: (value) => ({
    label: 'Doc',
    href: value,
  }),
  DownloadLink: (value) => ({
    label: 'Download File',
    href: '/api/system/fileDownload?file=' + encodeURIComponent(value),
  }),
  ViewLink: (value) => ({
    label: 'View File',
    href: '/api/system/fileView?file=' + encodeURIComponent(value),
  }),
};

const linkOrder = [];

const EventTypeFormatters = {
  [eventDetailFormatterKey(EVENT_TYPE.DEPLOY_DESIGN)]: DeploymentSummaryFormatter,
  [eventDetailFormatterKey(EVENT_TYPE.UNDEPLOY_DESIGN)]: DeploymentSummaryFormatter,
  [eventDetailFormatterKey(EVENT_TYPE.EVALUATE_DESIGN)]: RelationshipEvaluationEventFormatter,
  [eventDetailFormatterKey(EVENT_TYPE.REGISTRANT_SUMMARY)]: RegistrantSummaryFormatter,
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
    ..._.omit(event.metadata, [...linkOrder, 'kind', 'ViewLink']),
    ShortDescription:
      event.metadata.error || !canTruncateDescription(event.description || '')
        ? null
        : event.description,
  };
  console.log(metadata);
  const order = [
    'ShortDescription',
    'connectionName',
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
  console.log(orderedMetadata);
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
  return (
    <FormatStructuredData propertyFormatters={PropertyLinkFormatters} data={filteredMetadata} />
  );
};
