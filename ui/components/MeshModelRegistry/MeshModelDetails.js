import React from 'react';
import useStyles from '../../assets/styles/general/tool.styles';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import { FormatStructuredData, reorderObjectProperties } from '../DataFormatter';
import { groupPossibleRelationsByKind, RelationInitiationType } from './helper';
import { Box, Typography, IconButton, Tooltip } from '@material-ui/core';
import InfoIcon from '@/assets/icons/InfoIcon';
import ConnectionIcon from '@/assets/icons/ConnectionIcon';
import { iconMedium, iconSmall } from '../../css/icons.styles';

const KeyValue = ({ property, value }) => (
  <>
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        margin: '0.3rem 0',
      }}
    >
      <p
        style={{
          padding: '0',
          margin: '0 0.5rem 0 0',
          fontSize: '16px',
          fontWeight: '600',
        }}
      >
        {property}
      </p>
      <p style={{ padding: '0', margin: '0', fontSize: '16px' }}>{value}</p>
    </div>
  </>
);

const Title = ({ title }) => (
  <>
    <p
      style={{
        marginTop: '0px',
        fontSize: '19px',
        fontWeight: 'bold',
      }}
    >
      {title}
    </p>
  </>
);

const ModelContents = ({ model }) => {
  const StyleClass = useStyles();
  const PropertyFormattersLeft = {
    version: (value) => <KeyValue property="API Version" value={value} />,
    hostname: (value) => <KeyValue property="Registrant" value={value} />,
    components: (value) => <KeyValue property="Components" value={value} />,
  };

  const metaDataLeft = {
    version: model.version,
    hostname: model.hostname,
    components: model.components === null ? '0' : model.components?.length.toString(),
  };

  const orderLeft = ['version', 'hostname', 'components'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    category: (value) => <KeyValue property="Category" value={value} />,
    duplicates: (value) => <KeyValue property="Duplicates" value={value} />,
    relationships: (value) => <KeyValue property="Relationships" value={value} />,
  };

  const metaDataRight = {
    category: model.category?.name,
    duplicates: model.duplicates.toString(),
    relationships: model.relationships === null ? '0' : model.relationships?.length.toString(),
  };

  const orderRight = ['category', 'duplicates', 'relationships'];
  const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);

  return (
    <>
      <div className={StyleClass.segment}>
        <div
          className={StyleClass.fullWidth}
          style={{
            display: 'flex',
            flexDirection: 'column',
            paddingRight: '1rem',
          }}
        >
          <FormatStructuredData
            data={orderdMetadataLeft}
            propertyFormatters={PropertyFormattersLeft}
            order={orderLeft}
          />
        </div>

        <div className={StyleClass.fullWidth} style={{ display: 'flex', flexDirection: 'column' }}>
          <FormatStructuredData
            data={orderdMetadataRight}
            propertyFormatters={PropertyFormattersRight}
            order={orderRight}
          />
        </div>
      </div>
    </>
  );
};

const ChildTypeTitle = ({ title }) => (
  <p
    style={{
      fontSize: '18px',
      margin: '0.7rem 0',
    }}
  >
    {title}
  </p>
);

const ChildTitle = ({ title }) => (
  <>
    <p
      style={{
        fontSize: '18px',
        fontWeight: 'bold',
        margin: '0.4rem 0',
      }}
    >
      {title}
    </p>
  </>
);

const ComponentsContents = ({ components }) => {
  const StyleClass = useStyles();
  return (
    <div>
      {components.map((component, index) => {
        const PropertyFormattersLeft = {
          version: (value) => <KeyValue property="API Version" value={value} />,
        };

        const metaDataLeft = {
          version: component.apiVersion,
        };

        const orderLeft = ['version'];
        const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

        const PropertyFormattersRight = {
          subCategory: (value) => <KeyValue property="Sub Category" value={value} />,
        };

        const metaDataRight = {
          subCategory: component.kind,
        };

        const orderRight = ['subCategory'];
        const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);

        return (
          <div
            key={index}
            style={{
              margin: '0.9rem 0',
            }}
          >
            <ChildTitle title={component.displayName} />
            <div className={StyleClass.segment}>
              <div
                style={{
                  display: 'flex',
                  width: '50%',
                  paddingRight: '1rem',
                }}
              >
                <FormatStructuredData
                  data={orderdMetadataLeft}
                  propertyFormatters={PropertyFormattersLeft}
                  order={orderLeft}
                />
              </div>
              <div style={{ display: 'flex', width: '50%' }}>
                <FormatStructuredData
                  data={orderdMetadataRight}
                  propertyFormatters={PropertyFormattersRight}
                  order={orderRight}
                />
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

const formatRelationTitle = (relationInitiationType, relatesTo) => {
  relationInitiationType = relationInitiationType === RelationInitiationType.FROM ? 'From' : 'To';
  if (relatesTo.kind === '*') {
    return `${relationInitiationType} to all components in model ${relatesTo?.rel?.model} (${relatesTo?.rel?.subType})`;
  }

  if (relatesTo.model === '*') {
    return `${relationInitiationType} to all components  (${relatesTo?.rel?.subType})`;
  }

  return `${relationInitiationType} ${relatesTo.kind} (${relatesTo?.rel?.subType})`;
};

const RelationListItem = ({ relatesTo, relationInitiationType }) => {
  const title = formatRelationTitle(relationInitiationType, relatesTo);
  const description = relatesTo.rel?.metadata?.description;

  const infoIconColor = '#FFF';

  return (
    <Box
      sx={{
        marginTop: '0.5rem',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
      }}
    >
      <Box
        sx={{
          display: 'flex',
          gap: '0.5rem',
          alignItems: 'center',
        }}
      >
        <ConnectionIcon {...iconSmall} fill="#fff" />

        <Typography variant="body2">{title}</Typography>
      </Box>
      <Box
        sx={{
          display: 'flex',
          gap: '0.5rem',
          alignItems: 'center',
        }}
      >
        <Tooltip title={description || ''} placement="top" arrow>
          <IconButton size="small">
            <InfoIcon {...iconSmall} fill={infoIconColor} />
          </IconButton>
        </Tooltip>
      </Box>
    </Box>
  );
};

const ComponentContents = ({ component, relationships }) => {
  const StyleClass = useStyles();
  const PropertyFormattersLeft = {
    version: (value) => <KeyValue property="API Version" value={value} />,
    modelName: (value) => <KeyValue property="Model Name" value={value} />,
    kind: (value) => <KeyValue property="Kind" value={value} />,
  };

  const relationsGroupedByKind = groupPossibleRelationsByKind(component, relationships);
  const availableRelationKinds = Object.keys(relationsGroupedByKind).filter(
    (kind) =>
      relationsGroupedByKind[kind]?.[RelationInitiationType.FROM]?.length > 0 ||
      relationsGroupedByKind[kind]?.[RelationInitiationType.TO]?.length > 0,
  );

  console.log('relationships', relationships);
  console.log('relationshipByKind', relationsGroupedByKind);

  const metaDataLeft = {
    version: component.apiVersion,
    modelName: component.model?.displayName,
    kind: component.kind,
  };

  const orderLeft = ['version', 'modelName', 'kind'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    registrant: (value) => <KeyValue property="Registrant" value={value} />,
    duplicates: (value) => <KeyValue property="Duplicates" value={value} />,
  };

  const metaDataRight = {
    registrant: component.displayhostname,
    duplicates: component.duplicates.toString(),
  };

  const componentDescription = JSON.parse(component.schema).description;

  const orderRight = ['registrant', 'duplicates'];
  const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);

  return (
    <>
      <p style={{ fontWeight: '600', margin: '0', fontSize: '16px' }}>Description</p>
      <p style={{ margin: '0', fontSize: '16px' }}>{componentDescription}</p>
      <div className={StyleClass.segment}>
        <div
          className={StyleClass.fullWidth}
          style={{
            display: 'flex',
            flexDirection: 'column',
            paddingRight: '1rem',
          }}
        >
          <FormatStructuredData
            data={orderdMetadataLeft}
            propertyFormatters={PropertyFormattersLeft}
            order={orderLeft}
          />
        </div>

        <div className={StyleClass.fullWidth} style={{ display: 'flex', flexDirection: 'column' }}>
          <FormatStructuredData
            data={orderdMetadataRight}
            propertyFormatters={PropertyFormattersRight}
            order={orderRight}
          />
        </div>
      </div>
      <ChildTitle title="Available Relationships" />
      {availableRelationKinds.map((kind) => (
        <Box
          key={kind}
          sx={{
            display: 'flex',
            flexDirection: 'column',
            borderRadius: '0.25rem',
            marginTop: '0.5rem',
            // padding: '0.5rem',
            // backgroundColor: uiTheme.theme.palette.secondary.mainBackground2
          }}
        >
          <p style={{ fontWeight: '600', margin: '0', fontSize: '14px' }}>{kind}</p>

          {Object.values(RelationInitiationType).map((relationInitiationType) =>
            relationsGroupedByKind[kind][relationInitiationType].map((from) => (
              <RelationListItem
                key={from.kind}
                relatesTo={from}
                backgroundColor="#fff"
                relationInitiationType={relationInitiationType}
              />
            )),
          )}
        </Box>
      ))}
    </>
  );
};

const RelationshipContents = ({ relationship }) => {
  const PropertyFormattersLeft = {
    version: (value) => <KeyValue property="API Version" value={value} />,
    modelName: (value) => <KeyValue property="Model Name" value={value} />,
    kind: (value) => <KeyValue property="Kind" value={value} />,
  };

  const metaDataLeft = {
    version: relationship.apiVersion,
    modelName: relationship.model?.displayName,
  };

  const orderLeft = ['version', 'modelName'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    registrant: (value) => <KeyValue property="Registrant" value={value} />,
    subType: (value) => <KeyValue property="Sub Type" value={value} />,
  };

  const metaDataRight = {
    registrant: relationship.displayhostname,
    subType: relationship.subType,
  };

  const orderRight = ['subType', 'registrant'];
  const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);

  return (
    <>
      <div
        style={{
          display: 'flex',
          marginTop: '12px',
          flexDirection: 'row',
        }}
      >
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            width: '50%',
            paddingRight: '1rem',
          }}
        >
          <FormatStructuredData
            data={orderdMetadataLeft}
            propertyFormatters={PropertyFormattersLeft}
            order={orderLeft}
          />
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
          <FormatStructuredData
            data={orderdMetadataRight}
            propertyFormatters={PropertyFormattersRight}
            order={orderRight}
          />
        </div>
      </div>
    </>
  );
};

const RelationshipsContents = ({ relationships }) => {
  const StyleClass = useStyles();
  return (
    <div>
      {relationships.map((rela, index) => {
        const PropertyFormattersLeft = {
          version: (value) => <KeyValue property="API Version" value={value} />,
        };

        const metaDataLeft = {
          version: rela.apiVersion,
        };

        const orderLeft = ['version'];
        const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

        const PropertyFormattersRight = {
          subType: (value) => <KeyValue property="Sub Category" value={value} />,
        };

        const metaDataRight = {
          subType: rela.subType,
        };

        const orderRight = ['subType'];
        const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);
        return (
          <div
            key={index}
            style={{
              margin: '0.9rem 0',
            }}
          >
            <ChildTitle title={rela.kind} />
            <div className={StyleClass.segment}>
              <div
                style={{
                  display: 'flex',
                  width: '50%',
                  paddingRight: '1rem',
                }}
              >
                <FormatStructuredData
                  data={orderdMetadataLeft}
                  propertyFormatters={PropertyFormattersLeft}
                  order={orderLeft}
                />
              </div>
              <div style={{ display: 'flex', width: '50%' }}>
                <FormatStructuredData
                  data={orderdMetadataRight}
                  propertyFormatters={PropertyFormattersRight}
                  order={orderRight}
                />
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

const RegistrantContent = ({ registrant }) => {
  const PropertyFormattersLeft = {
    models: (value) => <KeyValue property="Models" value={value} />,
    components: (value) => <KeyValue property="Components" value={value} />,
  };

  const metaDataLeft = {
    models: registrant.summary?.models.toString(),
    components: registrant.summary?.components.toString(),
  };

  const orderLeft = ['models', 'components'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    relationships: (value) => <KeyValue property="Relationships" value={value} />,
    policies: (value) => <KeyValue property="Policies" value={value} />,
  };

  const metaDataRight = {
    relationships: registrant.summary?.relationships.toString(),
    policies: registrant.summary?.policies.toString(),
  };

  const orderRight = ['relationships', 'policies'];
  const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);
  return (
    <>
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          width: '50%',
          paddingRight: '1rem',
        }}
      >
        <FormatStructuredData
          data={orderdMetadataLeft}
          propertyFormatters={PropertyFormattersLeft}
          order={orderLeft}
        />
      </div>
      <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
        <FormatStructuredData
          data={orderdMetadataRight}
          propertyFormatters={PropertyFormattersRight}
          order={orderRight}
        />
      </div>
    </>
  );
};

const MeshModelDetails = ({ view, show, rela, regi, comp, relationships }) => {
  const StyleClass = useStyles();

  return (
    <div
      className={
        (view === MODELS && !show.model.displayName) ||
        (view === COMPONENTS && !comp.displayName) ||
        (view === RELATIONSHIPS && !rela.kind) ||
        (view === REGISTRANTS && !regi.hostname)
          ? StyleClass.emptyDetailsContainer
          : StyleClass.detailsContainer
      }
    >
      {((view === MODELS && !show.model.displayName) ||
        (view === COMPONENTS && !comp.displayName) ||
        (view === RELATIONSHIPS && !rela.kind)) && (
        <p style={{ color: '#969696', margin: 'auto' }}>No {view} selected</p>
      )}
      {view === REGISTRANTS && regi.hostname && (
        <div>
          <Title title={regi.hostname} />

          <RegistrantContent registrant={regi} />

          {show.model.displayName && <hr style={{ margin: '1rem 0' }} />}
        </div>
      )}
      {(view === MODELS || view === REGISTRANTS) && (
        <>
          {show.model.displayName && (
            <div>
              <Title title={show.model.displayName} />
              <ModelContents model={show.model} />
            </div>
          )}
          {show.components.length !== 0 && (
            <div>
              <hr style={{ margin: '1rem 0' }} />
              <ChildTypeTitle title="Components" />
              <ComponentsContents components={show.components} />
            </div>
          )}
          {show.relationships.length !== 0 && (
            <div>
              <hr style={{ marginTop: '1rem 0' }} />
              <ChildTypeTitle title="Relationships" />
              <RelationshipsContents relationships={show.relationships} />
            </div>
          )}
        </>
      )}
      {view === COMPONENTS && comp.displayName && (
        <>
          <div>
            <img src={comp.metadata.svgColor} height="80px" width="80px"/>
            <ChildTitle title={comp.displayName} />
            <ComponentContents component={comp} relationships={relationships} />
          </div>
          {/* <div>
            <hr style={{ marginTop: '1rem 0' }} />
            <ChildTypeTitle title="Relationships" />
            <RelationshipsContents relationships={relationships} component={comp}/>
          </div> */}
        </>
      )}
      {view === RELATIONSHIPS && rela.kind && (
        <div>
          <ChildTitle title={rela.kind} />
          <p style={{ fontWeight: '600', margin: '0', fontSize: '16px' }}>Description</p>
          <p style={{ margin: '0', fontSize: '16px' }}>{rela.metadata?.description}</p>
          <RelationshipContents relationship={rela} />
        </div>
      )}
    </div>
  );
};

export default MeshModelDetails;
