import * as React from 'react';
import { Typography, Grid, Box, Stack, styled, useTheme } from '@layer5/sistent';
import { Launch as LaunchIcon } from '@material-ui/icons';
import {
  FormatStructuredData,
  SectionBody,
  reorderObjectProperties,
  TextWithLinks,
} from '../DataFormatter';
import { alpha } from '@mui/material';
import { SEVERITY_STYLE, SEVERITY } from '../NotificationCenter/constants';

import { isEmptyAtAllDepths } from '../../utils/objects';
import { canTruncateDescription } from './notification';
import { FormatDryRunResponse } from '../DesignLifeCycle/DryRun';
import { formatDryRunResponse } from 'machines/validator/designValidator';
import { DeploymentSummaryFormatter } from '../DesignLifeCycle/DeploymentSummary';

const StyledDetailBox = styled(Box)(({ theme, severityColor, bgOpacity }) => ({
  padding: theme.spacing(2),
  backgroundColor: alpha(severityColor, bgOpacity),
  border: `1px solid ${theme.palette.divider}`,
  display: 'flex',
  justifyContent: 'space-between',
}));

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

const ComponentWithIcon = ({ component }) => {
  const theme = useTheme();
  const { DisplayName, Metadata, Model } = component;
  const svgColor = theme.palette.mode === 'dark' ? Metadata.svgWhite : Metadata.svgColor;

  return (
    <>
      <Typography variant="body1">{`Model ${Model.displayName} imported Component`}</Typography>
      <Grid container alignItems="center" spacing={1}>
        <Grid item>
          <Typography variant="body1">{DisplayName}</Typography>
        </Grid>
        <Grid item>
          <Typography variant="body2">
            {Model.model ? Model.model.version : 'Version not available'}
          </Typography>
        </Grid>
        <Grid item>
          <div
            style={{
              maxWidth: '30px',
              maxHeight: '30px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
            dangerouslySetInnerHTML={{ __html: svgColor }}
          />
        </Grid>
      </Grid>
    </>
  );
};

const RelationshipDetail = ({ relationship }) => {
  const { Model, Kind, Subtype } = relationship;
  return (
    <Typography variant="body1">
      {Kind} {Subtype} {Model.displayName}
    </Typography>
  );
};

const ModelDetail = ({ model }) => {
  const { Subtype } = model;
  if (Subtype !== undefined) {
    console.log(Subtype);
  }
  return <></>;
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

const ModelImportedSection = ({ modelName }) => (
  <Typography variant="h6">{`Imported Model(s) ${modelName}`}</Typography>
);

const ComponentsSection = ({ components }) => (
  <>
    {components.map((component, index) => (
      <ComponentWithIcon key={index} component={component} />
    ))}
  </>
);

const RelationshipsSection = ({ relationships }) => (
  <>
    {relationships.map((relationship, index) => (
      <RelationshipDetail key={index} relationship={relationship} />
    ))}
  </>
);

const ImportedModel = ({ models }) => (
  <>
    {models.map((model, index) => (
      <ModelDetail key={index} model={model} />
    ))}
  </>
);

const UnsuccessfulEntityWithError = ({ errors }) => {
  const parsedErrors = Array.isArray(errors) ? errors : [];

  return (
    <Stack spacing={2}>
      {parsedErrors.map((entry, index) => (
        <Box key={index}>
          <Typography variant="h6">
            {`Entity with filename `}
            <span style={{ textDecoration: 'underline' }}>{entry.name}</span>
            {entry.entityType ? ` of type ${entry.entityType}` : ''}
            {` failed to import`}
          </Typography>
          <ErrorMetadataFormatter metadata={entry.error} event={{}} />
        </Box>
      ))}
    </Stack>
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
    ImportedModelName: (value) =>
      value && (
        <StyledDetailBox
          severityColor={SEVERITY_STYLE[SEVERITY.INFO].color}
          bgOpacity={0.1}
          display="flex"
          gap={2}
          flexDirection="column"
        >
          <ModelImportedSection modelName={value} />
        </StyledDetailBox>
      ),
    ImportedComponent: (value) =>
      value && (
        <StyledDetailBox
          severityColor={SEVERITY_STYLE[SEVERITY.INFO].color}
          bgOpacity={0.1}
          display="flex"
          gap={2}
          flexDirection="column"
        >
          <ComponentsSection components={value} />
        </StyledDetailBox>
      ),
    ImportedRelationship: (value) => value && <RelationshipsSection relationships={value} />,
    ImportedModel: (value) =>
      value && (
        // <StyledDetailBox
        //   severityColor={SEVERITY_STYLE[SEVERITY.INFO].color}
        //   bgOpacity={0.1}
        //   display="flex"
        //   gap={2}
        //   flexDirection="column"
        // >
        <ImportedModel models={value} />
        // </StyledDetailBox>
      ),
    UnsuccessfulEntityNameWithError: (value) =>
      value && (
        <StyledDetailBox
          severityColor={SEVERITY_STYLE[SEVERITY.ERROR].color}
          bgOpacity={0.1}
          display="flex"
          gap={2}
          flexDirection="column"
        >
          <UnsuccessfulEntityWithError errors={value} />
        </StyledDetailBox>
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
    'ImportedModelName',
    'ImportedComponent',
    'ImportedRelationship',
    'ImportedModel',
    'UnsuccessfulEntityNameWithError',
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
    />
  );
};
