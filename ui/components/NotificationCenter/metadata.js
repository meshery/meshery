import React, { useEffect, useState } from 'react';
import { Typography, Grid, Box, styled, List, ListItem } from '@layer5/sistent';
import { Launch as LaunchIcon } from '@material-ui/icons';
import { alpha } from '@mui/material';
import { SEVERITY_STYLE, SEVERITY } from '../NotificationCenter/constants';
import {
  FormatStructuredData,
  SectionBody,
  reorderObjectProperties,
  TextWithLinks,
} from '../DataFormatter';
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

const UnsuccessfulEntityWithError = ({ modelName, error }) => {
  const entityTypesAndQuantities = {};
  error.name.forEach((name, idx) => {
    if (name === modelName) {
      entityTypesAndQuantities[error.entityType[idx]] =
        entityTypesAndQuantities[error.entityType[idx]] + 1 || 1;
    }
  });

  const isEntityFile =
    modelName.includes('.yaml') || modelName.includes('.yml') || modelName.includes('.json');
  const message = `Import did not occur for ${Object.entries(entityTypesAndQuantities)
    .map(([key, value]) => `${value} ${value > 1 ? 'entities' : 'entity'} of type ${key}`)
    .join(', ')}.`;
  if (isEntityFile) {
    return (
      <>
        <Typography variant="body1">
          Import process for file {modelName} encountered error
        </Typography>
        <ErrorMetadataFormatter metadata={error.error} event={{}} />
      </>
    );
  }

  return (
    <>
      <Typography variant="body1">{message}</Typography>
      <ErrorMetadataFormatter metadata={error.error} event={{}} />
    </>
  );
};

const imageCache = {};

const checkImageExists = async (url) => {
  try {
    const response = await fetch(url);
    return response.ok;
  } catch (error) {
    return false;
  }
};

const ComponentWithIcon = ({ component }) => {
  const { DisplayName, Metadata, Model, Version } = component;
  const modelname = Model;
  const kind = Metadata.toLowerCase();

  const paths = [
    `ui/public/static/img/meshmodels/${modelname}/color/${modelname}-color.svg`,
    `ui/public/static/img/meshmodels/${modelname}/white/${modelname}-white.svg`,
    `ui/public/static/img/meshmodels/${modelname}/color/${kind}-color.svg`,
    `ui/public/static/img/meshmodels/${modelname}/white/${kind}-white.svg`,
  ];

  const [finalPath, setFinalPath] = useState(
    'ui/public/static/img/meshmodels/meshery-core/color/meshery-core-color.svg',
  );

  useEffect(() => {
    const loadImages = async () => {
      if (imageCache[modelname]) {
        setFinalPath(imageCache[modelname]);
        return;
      }

      for (const path of paths) {
        const exists = await checkImageExists(path);
        if (exists) {
          imageCache[modelname] = path;
          setFinalPath(path);
          return;
        }
      }
    };

    loadImages();
  }, [modelname, paths]);

  const version = Version
    ? Version.startsWith('v')
      ? Version
      : `v${Version}`
    : 'Version not available';

  return (
    <Grid
      container
      alignItems="center"
      spacing={1}
      style={{ marginBottom: '8px', marginLeft: '1rem', marginTop: '8px' }}
    >
      <Grid item>
        <div
          style={{
            maxWidth: '30px',
            maxHeight: '30px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <img src={finalPath} style={{ width: '30px', height: '30px' }} alt={DisplayName} />
        </div>
      </Grid>
      <Grid item>
        <Typography variant="body1" component="span">
          {DisplayName}
        </Typography>
        <Typography variant="body1" component="span" style={{ marginLeft: '1rem' }}>
          {version}
        </Typography>
      </Grid>
    </Grid>
  );
};

const RelationshipDetail = ({ relationship }) => {
  const { Kind, Subtype, Selectors, RelationshipType } = relationship;

  const renderSelectors = (selectors) => {
    return selectors.map((selector, index) => {
      return (
        <div key={index} style={{ marginLeft: '2rem', marginTop: '0.5rem' }}>
          <React.Fragment key={index}>
            <Typography variant="body1">
              Kind of {Kind}, sub type {Subtype} and type {RelationshipType}
            </Typography>
            <Grid container>
              <Grid item xs={6}>
                <Typography variant="body2">
                  <strong>FROM</strong>
                </Typography>
                {selector.allow.from.map((f, idx) => (
                  <Typography key={idx} variant="body2">
                    {`${f.kind}`}
                  </Typography>
                ))}
              </Grid>
              <Grid item xs={6}>
                <Typography variant="body2">
                  <strong>TO</strong>
                </Typography>
                {selector.allow.to.map((t, idx) => (
                  <Typography key={idx} variant="body2">
                    {`${t.kind}`}
                  </Typography>
                ))}
              </Grid>
            </Grid>
          </React.Fragment>
        </div>
      );
    });
  };

  return <Box>{Selectors && Selectors.length > 0 && renderSelectors(Selectors)}</Box>;
};

export const ErrorMetadataFormatter = ({ metadata, event }) => {
  const longDescription = metadata?.LongDescription || [];
  const probableCause = metadata?.ProbableCause || [];
  const suggestedRemediation = metadata?.SuggestedRemediation || [];
  const errorCode = metadata?.error_code || '';
  const code = metadata?.Code || '';
  const formattedErrorCode = errorCode ? `${errorCode}-${code}` : code;
  const errorLink = `https://docs.meshery.io/reference/error-codes#${formattedErrorCode}`;
  const ErrorDetailsObjectFormatter = ({ heading, value }) => {
    return (
      <Box>
        <Typography variant="body1">
          <strong>{heading}</strong>
        </Typography>
        <List sx={{}}>
          {value.map((error, idx) => (
            <ListItem key={idx}>
              <Box>
                <Typography variant="body1">{error}</Typography>
              </Box>
            </ListItem>
          ))}
        </List>
      </Box>
    );
  };
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
            propertyFormatters={{
              Details: (value) => <ErrorDetailsObjectFormatter heading="Details" value={value} />,
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
            propertyFormatters={{
              'Probable Cause': (value) => (
                <ErrorDetailsObjectFormatter heading="Probable Cause" value={value} />
              ),
            }}
          />
        </Grid>
        <Grid item sm={suggestedRemediation?.length > 0 ? 6 : 12}>
          <FormatStructuredData
            data={{
              'Suggested Remediation': suggestedRemediation,
            }}
            propertyFormatters={{
              'Suggested Remediation': (value) => (
                <ErrorDetailsObjectFormatter heading="Suggested Remediation" value={value} />
              ),
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

const ModelImportedSection = ({ modelDetails }) => {
  if (typeof modelDetails !== 'object' || Array.isArray(modelDetails)) {
    return null;
  }

  return (
    <>
      {Object.keys(modelDetails).map((modelName, index) => {
        const detail = modelDetails[modelName];
        const isEntityFile =
          modelName.includes('.yaml') || modelName.includes('.yml') || modelName.includes('.json');

        const isMultipleComponents =
          Array.isArray(detail.Components) && detail.Components.length > 1;
        const isMultipleRelationships =
          Array.isArray(detail.Relationships) && detail.Relationships.length > 1;
        const hasComponents = Array.isArray(detail.Components) && detail.Components.length > 0;
        const hasRelationships =
          Array.isArray(detail.Relationships) && detail.Relationships.length > 0;
        const hasErrors = Array.isArray(detail.Errors) && detail.Errors.length > 0;

        return (
          <Box key={index} mb={2}>
            <Typography gutterBottom>
              <span style={{ fontWeight: 'bold', fontSize: '17px' }}>
                {isEntityFile ? 'FILE NAME:' : 'MODEL:'}{' '}
              </span>
              <span style={{ fontSize: '18px' }}>{modelName}</span>
            </Typography>
            {hasComponents && (
              <>
                <StyledDetailBox
                  severityColor={SEVERITY_STYLE[SEVERITY.SUCCESS].color}
                  bgOpacity={0.2}
                  display="flex"
                  gap={2}
                  flexDirection="column"
                >
                  <Typography variant="body1">
                    <span style={{ fontWeight: 'bold', fontSize: '14px', marginLeft: '1rem' }}>
                      {isMultipleComponents ? 'COMPONENTS:' : 'COMPONENT:'}
                    </span>
                  </Typography>
                  {detail.Components.map((component, idx) => (
                    <Box key={idx} ml={2}>
                      <ComponentWithIcon component={component} />
                    </Box>
                  ))}
                </StyledDetailBox>
              </>
            )}
            {hasRelationships && (
              <>
                <StyledDetailBox
                  severityColor={SEVERITY_STYLE[SEVERITY.SUCCESS].color}
                  bgOpacity={0.2}
                  display="flex"
                  gap={2}
                  flexDirection="column"
                >
                  <Typography variant="body1">
                    <span style={{ fontWeight: 'bold', fontSize: '14px', marginLeft: '1rem' }}>
                      {isMultipleRelationships ? 'RELATIONSHIPS:' : 'RELATIONSHIP:'}
                    </span>
                  </Typography>
                  {detail.Relationships.map((relationship, idx) => (
                    <Box key={idx} ml={2}>
                      <RelationshipDetail relationship={relationship} />
                    </Box>
                  ))}
                </StyledDetailBox>
              </>
            )}
            {hasErrors &&
              detail.Errors.map((error, idx) => (
                <Box key={idx}>
                  <StyledDetailBox
                    severityColor={SEVERITY_STYLE[SEVERITY.ERROR].color}
                    bgOpacity={0.2}
                    display="flex"
                    gap={2}
                    flexDirection="column"
                  >
                    <UnsuccessfulEntityWithError modelName={modelName} error={error} />
                  </StyledDetailBox>
                </Box>
              ))}
          </Box>
        );
      })}
    </>
  );
};

const ModelImportMessages = ({ message }) => (
  <Typography>
    <span style={{ fontWeight: 'bold', fontSize: '17px' }}>{`SUMMARY: `}</span>
    <span style={{ fontSize: '17px' }}>{message}</span>
  </Typography>
);

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
          <ModelImportMessages message={value} />
        </StyledDetailBox>
      ),

    ModelDetails: (value) =>
      value && (
        <StyledDetailBox
          severityColor={SEVERITY_STYLE[SEVERITY.INFO].color}
          bgOpacity={0.1}
          display="flex"
          gap={2}
          flexDirection="column"
        >
          <ModelImportedSection modelDetails={value} />
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
    />
  );
};
