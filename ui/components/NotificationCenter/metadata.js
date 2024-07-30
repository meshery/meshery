import * as React from 'react';
import { Typography, Grid } from '@material-ui/core';
import { Launch as LaunchIcon } from '@material-ui/icons';
import { FormatStructuredData, SectionBody, reorderObjectProperties } from '../DataFormatter';
import { isEmptyAtAllDepths } from '../../utils/objects';
import { canTruncateDescription } from './notification';
import { TextWithLinks } from '../DataFormatter';
import { FormatDryRunResponse } from '../DesignLifeCycle/DryRun';
import { formatDryRunResponse } from 'machines/validator/designValidator';
import { DeploymentSummaryFormatter } from '../DesignLifeCycle/DeploymentSummary';

const DryRunResponse = ({ response }) => {
  return <FormatDryRunResponse dryRunErrors={formatDryRunResponse(response)} />;
};

const TitleLink = ({ href, children, ...props }) => {
  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      style={{ color: 'inherit' }}
      {...props}
    >
      <Typography
        variant="h5"
        style={{
          textDecorationLine: 'underline',
          cursor: 'pointer',
          marginBottom: '0.5rem',
          fontWeight: 'bolder !important',
          textTransform: 'uppercase',
          fontSize: '0.9rem',
        }}
      >
        {children}
        <sup>
          <LaunchIcon style={{ width: '1rem', height: '1rem' }} />
        </sup>
      </Typography>
    </a>
  );
};

export const ErrorMetadataFormatter = ({ metadata, event }) => {
  const longDescription = metadata?.LongDescription || [];
  const probableCause = metadata?.ProbableCause || [];
  const suggestedRemediation = metadata?.SuggestedRemediation || [];
  const errorCode = metadata?.error_code || '';
  const code = metadata?.Code || '';
  const formattedErrorCode = errorCode ? `${errorCode}-${code}` : code;
  const errorLink = `https://docs.meshery.io/reference/error-codes#${formattedErrorCode}`;
  return (
    <Grid container>
      <div>
        <TitleLink href={errorLink}> {formattedErrorCode} </TitleLink>
        {event?.description && <FormatStructuredData data={event.description} />}
        <div style={{ marginTop: '1rem' }}>
          <FormatStructuredData
            data={{
              Details: longDescription,
            }}
          />
        </div>
      </div>
      <Grid container spacing={1} style={{ marginTop: '0.5rem' }}>
        <Grid item sm={probableCause?.length > 0 ? 6 : 12}>
          <FormatStructuredData
            data={{
              'Probable Cause': probableCause,
            }}
          />
        </Grid>
        <Grid item sm={suggestedRemediation?.length > 0 ? 6 : 12}>
          <FormatStructuredData
            data={{
              'Suggested Remediation': suggestedRemediation,
            }}
          />
        </Grid>
      </Grid>
    </Grid>
  );
};

const EmptyState = ({ event }) => {
  return (
    <Typography
      variant="body1"
      style={{
        marginBlock: '0.5rem',
      }}
    >
      {' '}
      {
        <TextWithLinks
          variant="body1"
          style={{
            wordWrap: 'break-word',
          }}
          text={event.description || ''}
        ></TextWithLinks>
      }{' '}
    </Typography>
  );
};

export const FormattedMetadata = ({ event }) => {
  const PropertyFormatters = {
    doc: (value) => <TitleLink href={value}>Doc</TitleLink>,
    ShortDescription: (value) => <SectionBody body={value} style={{ marginBlock: '0.5rem' }} />,
    error: (value) => <ErrorMetadataFormatter metadata={value} event={event} />,
    dryRunResponse: (value) => <DryRunResponse response={value} />,
    DownloadLink: (value) => (
      <TitleLink href={'/api/system/fileDownload?file=' + encodeURIComponent(value)}>
        Download
      </TitleLink>
    ),
    ViewLink: (value) => (
      <TitleLink href={'/api/system/fileView?file=' + encodeURIComponent(value)}>View</TitleLink>
    ),
  };

  const EventTypeFormatters = {
    deploy: DeploymentSummaryFormatter,
    undeploy: DeploymentSummaryFormatter,
  };

  if (EventTypeFormatters[event.action]) {
    const Formatter = EventTypeFormatters[event.action];
    return <Formatter event={event} />;
  }

  if (!event || !event.metadata || isEmptyAtAllDepths(event.metadata)) {
    return <EmptyState event={event} />;
  }

  const metadata = {
    ...event.metadata,
    ShortDescription:
      event.metadata.error || !canTruncateDescription(event.description || '')
        ? null
        : event.description,
  };

  const order = [
    'doc',
    'ShortDescription',
    'LongDescription',
    'Summary',
    'SuggestedRemediation',
    'DownloadLink',
    'ViewLink',
  ];
  const orderdMetadata = reorderObjectProperties(metadata, order);
  return (
    <FormatStructuredData
      propertyFormatters={PropertyFormatters}
      data={orderdMetadata}
      order={order}
    />
  );
};
