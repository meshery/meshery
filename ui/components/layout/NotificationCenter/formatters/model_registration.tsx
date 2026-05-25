import React, { useEffect, useState } from 'react';
import { Typography, Grid2, Box } from '@sistent/sistent';
import { ErrorMetadataFormatter } from './error';
import { TitleLink } from './common';
import { FALLBACK_MESHERY_IMAGE_PATH } from '@/constants/common';
import { normalizeStaticImagePath } from '@/utils/fallback';
import { iconMedium } from 'css/icons.styles';

type ImportError = {
  name: string[];
  entityType: string[];
  error: Record<string, unknown>;
};

const UnsuccessfulEntityWithError = ({
  modelName,
  error,
}: {
  modelName: string;
  error: ImportError;
}) => {
  const entityTypesAndQuantities: Record<string, number> = {};
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

const checkImageExists = async (url: string): Promise<boolean> => {
  try {
    const response = await fetch(url);
    return response.ok;
  } catch {
    return false;
  }
};

type ModelComponent = {
  DisplayName: string;
  Metadata: string;
  Model: string;
  Version?: string;
};

const ComponentWithIcon = ({ component }: { component: ModelComponent }) => {
  const { DisplayName, Metadata, Model, Version } = component;
  const modelname = Model;
  const kind = Metadata.toLowerCase();

  const paths = [
    `/ui/public/static/img/meshmodels/${modelname}/color/${kind}-color.svg`,
    `/ui/public/static/img/meshmodels/${modelname}/white/${kind}-white.svg`,
    `/ui/public/static/img/meshmodels/${modelname}/color/${modelname}-color.svg`,
    `/ui/public/static/img/meshmodels/${modelname}/white/${modelname}-white.svg`,
  ].map((path) => normalizeStaticImagePath(path));

  const defaultPath = normalizeStaticImagePath(FALLBACK_MESHERY_IMAGE_PATH);

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
    <Grid2
      container
      alignItems="center"
      spacing={1}
      style={{
        marginBottom: '8px',
        marginLeft: '1rem',
        marginTop: '8px',
      }}
      size="grow"
    >
      <Grid2>
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
      </Grid2>
      <Grid2>
        <Typography variant="body1" component="span">
          {DisplayName}
        </Typography>
        <Typography variant="body1" component="span" style={{ marginLeft: '1rem' }}>
          {version}
        </Typography>
      </Grid2>
    </Grid2>
  );
};

type Selector = {
  allow: { from: { kind: string }[]; to: { kind: string }[] };
};

type Relationship = {
  Kind: string;
  Subtype: string;
  RelationshipType: string;
  Selectors?: Selector[];
};

const RelationshipDetail = ({ relationship }: { relationship: Relationship }) => {
  const { Kind, Subtype, Selectors, RelationshipType } = relationship;

  const renderSelectors = (selectors: Selector[]) => {
    return selectors.map((selector, index) => {
      return (
        <div key={index} style={{ marginLeft: '2rem', marginTop: '0.5rem' }}>
          <React.Fragment key={index}>
            <Typography variant="body1">
              Kind of {Kind}, sub type {Subtype} and type {RelationshipType}
            </Typography>
            <Grid2 container size="grow">
              <Grid2 size={{ xs: 6 }}>
                <Typography variant="body2">
                  <strong>FROM</strong>
                </Typography>
                {selector.allow.from.map((f, idx) => (
                  <Typography key={idx} variant="body2">
                    {`${f.kind}`}
                  </Typography>
                ))}
              </Grid2>
              <Grid2 size={{ xs: 6 }}>
                <Typography variant="body2">
                  <strong>TO</strong>
                </Typography>
                {selector.allow.to.map((t, idx) => (
                  <Typography key={idx} variant="body2">
                    {`${t.kind}`}
                  </Typography>
                ))}
              </Grid2>
            </Grid2>
          </React.Fragment>
        </div>
      );
    });
  };

  return <Box>{Selectors && Selectors.length > 0 && renderSelectors(Selectors)}</Box>;
};

type ModelDetail = {
  Components?: ModelComponent[];
  Relationships?: Relationship[];
  Errors?: ImportError[];
};

export const ModelImportedSection = ({
  modelDetails,
}: {
  modelDetails: Record<string, ModelDetail> | unknown;
}) => {
  if (typeof modelDetails !== 'object' || Array.isArray(modelDetails)) {
    return null;
  }

  const details = modelDetails as Record<string, ModelDetail>;

  return (
    <>
      {Object.keys(details).map((modelName, index) => {
        const detail = details[modelName];
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

export const ModelImportMessages = ({ message }: { message: React.ReactNode }) => (
  <Typography data-testid="ModelImportMessages-Wrapper">
    <span style={{ fontWeight: 'bold', fontSize: '17px' }}>{`SUMMARY: `}</span>
    <span style={{ fontSize: '17px' }}>{message}</span>
  </Typography>
);
