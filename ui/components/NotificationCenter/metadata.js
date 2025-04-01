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
  error: (value) => <ErrorMetadataFormatter metadata={value} event={event} />,
  dryRunResponse: (value) => <DryRunResponse response={value} />,
  ModelImportMessage: (value) => value && <ModelImportMessages message={value} />,
  ModelDetails: (value) => value && <ModelImportedSection modelDetails={value} />,
};

const PropertyLinkFormatters = {
  doc: (value) => (
    <TitleLink href={value} style={{ textAlign: 'end', color: 'inherit' }}>
      Doc
    </TitleLink>
  ),
  DownloadLink: (value) => (
    <TitleLink
      href={'/api/system/fileDownload?file=' + encodeURIComponent(value)}
      style={{ textAlign: 'end', color: 'inherit' }}
    >
      Download File
    </TitleLink>
  ),
  ViewLink: (value) => (
    <TitleLink
      href={'/api/system/fileView?file=' + encodeURIComponent(value)}
      style={{ textAlign: 'end', color: 'inherit' }}
    >
      View File
    </TitleLink>
  ),
};

const linkOrder = ['doc', 'DownloadLink', 'ViewLink'];

const EventTypeFormatters = {
  [eventDetailFormatterKey(EVENT_TYPE.DEPLOY_DESIGN)]: DeploymentSummaryFormatter,
  [eventDetailFormatterKey(EVENT_TYPE.UNDEPLOY_DESIGN)]: DeploymentSummaryFormatter,
  [eventDetailFormatterKey(EVENT_TYPE.EVALUATE_DESIGN)]: RelationshipEvaluationEventFormatter,
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
    ..._.omit(event.metadata, linkOrder),
    ShortDescription:
      event.metadata.error || !canTruncateDescription(event.description || '')
        ? null
        : event.description,
  };

  const order = [
    'ShortDescription',
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
    />
  );
};

export const FormattedLinkMetadata = ({ event }) => {
  const filteredMetadata = _.pick(event.metadata, linkOrder);
  return (
    <FormatStructuredData propertyFormatters={PropertyLinkFormatters} data={filteredMetadata} />
  );
};
