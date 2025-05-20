import React, { useEffect, useState } from 'react';
import { Typography, Grid, Box } from '@layer5/sistent';
import { ErrorMetadataFormatter } from './error';
import { TitleLink } from './common';
import { FALLBACK_MESHERY_IMAGE_PATH } from '@/constants/common';
import { iconMedium } from 'css/icons.styles';

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

const checkImageExists = async (url) => {
  try {
    const response = await fetch(url);
    return response.ok;
  } catch {
    return false;
  }
};

const ComponentWithIcon = ({ component }) => {
  const { DisplayName, Metadata, Model, Version } = component;
  const modelname = Model;
  const kind = Metadata.toLowerCase();

  const paths = [
    `/static/img/meshmodels/${modelname}/color/${kind}-color.svg`,
    `/static/img/meshmodels/${modelname}/white/${kind}-white.svg`,
    `/static/img/meshmodels/${modelname}/color/${modelname}-color.svg`,
    `/static/img/meshmodels/${modelname}/white/${modelname}-white.svg`,
  ];

  const defaultPath = FALLBACK_MESHERY_IMAGE_PATH;

  const [finalPath, setFinalPath] = useState(defaultPath);

  useEffect(() => {
    const loadImages = async () => {
      for (let i = 0; i < 2; i++) {
        const exists = await checkImageExists(paths[i]);
        if (exists) {
          setFinalPath(paths[i]);
          return;
        }
      }

      for (let i = 2; i < 4; i++) {
        const exists = await checkImageExists(paths[i]);
        if (exists) {
          setFinalPath(paths[i]);
          return;
        }
      }
    };

    loadImages();
  }, [modelname, kind, paths]);

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
      style={{
        marginBottom: '8px',
        marginLeft: '1rem',
        marginTop: '8px',
      }}
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
          <img src={finalPath} {...iconMedium} alt={DisplayName} />
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

export const ModelImportedSection = ({ modelDetails }) => {
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
            <div
              style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}
              data-testid={`ModelImportedSection-ModelHeader-${modelName}`}
            >
              <Typography gutterBottom>
                <span style={{ fontWeight: 'bold', fontSize: '17px' }}>
                  {isEntityFile ? 'FILE NAME:' : 'MODEL:'}{' '}
                </span>
                <span style={{ fontSize: '18px' }}>{modelName}</span>
              </Typography>
              {!isEntityFile && (
                <TitleLink
                  href={`settings?settingsCategory=Registry&tab=Models&searchText=${modelName}`}
                  target="_self"
                >
                  Registry
                </TitleLink>
              )}
            </div>
            {hasComponents && (
              <>
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
              </>
            )}
            {hasRelationships && (
              <>
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
              </>
            )}
            {hasErrors &&
              detail.Errors.map((error, idx) => (
                <Box key={idx}>
                  <UnsuccessfulEntityWithError modelName={modelName} error={error} />
                </Box>
              ))}
          </Box>
        );
      })}
    </>
  );
};

export const ModelImportMessages = ({ message }) => (
  <Typography data-testid="ModelImportMessages-Wrapper">
    <span style={{ fontWeight: 'bold', fontSize: '17px' }}>{`SUMMARY: `}</span>
    <span style={{ fontSize: '17px' }}>{message}</span>
  </Typography>
);
