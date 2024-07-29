import * as React from 'react';
import { Typography, Grid, Box, styled, useTheme } from '@layer5/sistent';
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
  const svgColor =
    theme.palette.mode === 'dark' && Metadata.svgWhite && Metadata.svgColor
      ? Metadata.svgWhite
      : Metadata.svgColor || Metadata.svgWhite;

  return (
    <Grid container alignItems="center" spacing={1}>
      <Grid item>
        <Typography variant="body1">{`${DisplayName}/${Model.displayName}`}</Typography>
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
  );
};

const RelationshipDetail = ({ relationship }) => {
  const { Model, Kind, Subtype, Selectors } = relationship;

  const renderSelectors = (selectors) => {
    return selectors.map((selector, index) => {
      const from = selector.allow.from.map((f) => `${f.kind}/${f.model}`).join(', ');
      const to = selector.allow.to.map((t) => `${t.kind}/${t.model}`).join(', ');
      return (
        <Typography key={index} variant="body1">
          <strong>Model Name:</strong> {Model.displayName} <strong>Kind:</strong> {Kind}{' '}
          <strong>SubType:</strong> {Subtype}
          <br />
          <strong>From</strong> {from}
          <strong> To</strong> {to}
        </Typography>
      );
    });
  };

  return <div>{Selectors && Selectors.length > 0 && <div>{renderSelectors(Selectors)}</div>}</div>;
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
  <Typography variant="body1">{`${modelName}`}</Typography>
);
const ModelImportMessage = ({ message }) => <Typography variant="body1">{message}</Typography>;
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

const UnsuccessfulEntityWithError = ({ errors }) => {
  const parsedErrors = Array.isArray(errors) ? errors : [];

  return (
    <>
      {parsedErrors.map((entry, index) => (
        <StyledDetailBox
          key={index}
          severityColor={SEVERITY_STYLE[SEVERITY.ERROR].color}
          bgOpacity={0.1}
          display="flex"
          gap={2}
          flexDirection="column"
        >
          <Box key={index}>
            <Typography variant="h6">Import encountered an error:</Typography>
            <Typography variant="body1">
              <ul>
                {entry.name.map((name, idx) => (
                  <li key={idx}>
                    <strong>Entity Type:</strong>{' '}
                    {entry.entityType && entry.entityType[idx]
                      ? entry.entityType[idx].charAt(0).toUpperCase() +
                        entry.entityType[idx].slice(1)
                      : 'N/A'}
                    ,<strong> Model Name:</strong> {name}
                  </li>
                ))}
              </ul>
            </Typography>
            <ErrorMetadataFormatter metadata={entry.error} event={{}} />
          </Box>
        </StyledDetailBox>
      ))}
    </>
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
    ModelImportMessage: (value) =>
      value && (
        <StyledDetailBox
          severityColor={SEVERITY_STYLE[SEVERITY.INFO].color}
          bgOpacity={0.1}
          display="flex"
          gap={2}
          flexDirection="column"
        >
          <ModelImportMessage message={value} />
        </StyledDetailBox>
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
          <Typography variant="h6">Imported Model(s)</Typography>
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
          <Typography variant="h6">Imported Component(s)</Typography>
          <ComponentsSection components={value} />
        </StyledDetailBox>
      ),
    ImportedRelationship: (value) =>
      value && (
        <StyledDetailBox
          severityColor={SEVERITY_STYLE[SEVERITY.INFO].color}
          bgOpacity={0.1}
          display="flex"
          gap={2}
          flexDirection="column"
        >
          <Typography variant="h6">Imported Relationship(s)</Typography>
          <RelationshipsSection relationships={value} />
        </StyledDetailBox>
      ),
    UnsuccessfulEntityNameWithError: (value) =>
      value && <UnsuccessfulEntityWithError errors={value} />,
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
    'ModelImportMessage',
    'ImportedModelName',
    'ImportedComponent',
    'ImportedRelationship',
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
