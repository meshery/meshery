import * as React from 'react';
import { Typography, Grid } from '@material-ui/core';
import { Launch as LaunchIcon } from '@material-ui/icons';
import { FormatStructuredData, LinkFormatters } from '../DataFormatter';
import { isEmptyAtAllDepths } from '../../utils/objects';

const DryRunResponse = ({ response }) => {
  const cleanedResponse = {};
  Object.entries(response).forEach(([componentKind, components]) => {
    Object.entries(components).forEach(([, data]) => {
      cleanedResponse[`Component: ${componentKind}`] = {
        ...data.error,
      };
    });
  });

  if (isEmptyAtAllDepths(cleanedResponse)) {
    return (
      <Typography variant="h6" style={{ textAlign: 'center', marginBlock: '1rem' }}>
        No Errors Found
      </Typography>
    );
  }
  return <FormatStructuredData data={cleanedResponse} />;
};

export const ErrorMetadataFormatter = ({ metadata, event, classes }) => {
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
        <a href={errorLink} target="_blank" rel="noopener noreferrer" style={{ color: 'inherit' }}>
          <Typography
            variant="h5"
            className={classes.descriptionHeading}
            style={{ textDecorationLine: 'underline', cursor: 'pointer', marginBottom: '0.5rem' }}
          >
            {formattedErrorCode}
            <sup>
              <LaunchIcon style={{ width: '1rem', height: '1rem' }} />
            </sup>
          </Typography>
        </a>
        <FormatStructuredData data={event.description} />
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
  return <Typography variant="body1"> {event.description} </Typography>;
};

export const FormattedMetadata = ({ event, classes }) => {
  const PropertyFormatters = {
    Doc: (value) => LinkFormatters.DOC.formatter(value),
    error: (value) => <ErrorMetadataFormatter metadata={value} event={event} classes={classes} />,
    dryRunResponse: (value) => <DryRunResponse response={value} />,
  };
  if (!event || !event.metadata || isEmptyAtAllDepths(event.metadata)) {
    return <EmptyState event={event} />;
  }

  return <FormatStructuredData propertyFormatters={PropertyFormatters} data={event.metadata} />;
};
