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
import { useTheme, DownloadIcon, InfoIcon, IconButton, Box, Typography } from '@layer5/sistent';
import _ from 'lodash';
import { ChipWrapper } from '../connections/styles';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import { CustomTooltip } from '@layer5/sistent';

const DesignFormatter = ({ value }) => {
  const theme = useTheme();
  const { name, id } = value;
  const [copySuccess, setCopySuccess] = React.useState('');

  const copyToClipboard = (e) => {
    e.stopPropagation();
    navigator.clipboard.writeText(id).then(
      () => {
        setCopySuccess('Design ID copied!');
        setTimeout(() => setCopySuccess(''), 2000);
      },
      (err) => {
        console.error('Could not copy text: ', err);
      },
    );
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
      <TitleLink
        href={'/extension/meshmap?mode=design&design=' + encodeURIComponent(id)}
        style={{
          color: theme.palette.primary.main,
          fontWeight: 'normal',
          textDecoration: 'underline',
        }}
        target="_self"
      >
        {name}
      </TitleLink>
      <CustomTooltip title={copySuccess || 'Copy Design ID'}>
        <IconButton
          size="small"
          onClick={copyToClipboard}
          style={{ alignSelf: 'flex-start', padding: '4px' }}
        >
          <ContentCopyIcon fontSize="small" />
          <Typography variant="caption" style={{ marginLeft: '4px' }}>
            Copy ID
          </Typography>
        </IconButton>
      </CustomTooltip>
    </Box>
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
        href={`/management/connections?tab=connections&searchText=${value}`}
        target="_self"
      />
    );
  },
  error: (value) => <ErrorMetadataFormatter metadata={value} event={event} />,
  dryRunResponse: (value) => <DryRunResponse response={value} />,
  ModelImportMessage: (value) => value && <ModelImportMessages message={value} />,
  ModelDetails: (value) => value && <ModelImportedSection modelDetails={value} />,
};

export const LinkFormatters = {
  doclink: (value) => {
    return (
      <TitleLink href={value} style={{ textAlign: 'end', color: 'inherit' }}>
        Documentation
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
    ..._.omit(event.metadata, [...linkOrder, 'kind', 'ViewLink', 'DownloadLink']),
    ShortDescription:
      event.metadata.error || !canTruncateDescription(event.description || '')
        ? null
        : event.description,
  };
  const order = [
    'doclink',
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
