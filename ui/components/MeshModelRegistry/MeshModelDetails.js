import React from 'react';
import useStyles from '../../assets/styles/general/tool.styles';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import { FormatStructuredData, reorderObjectProperties } from '../DataFormatter';

const KeyValue = ({ property, value }) => {
  let formattedValue = value;

  if (Array.isArray(value)) {
    formattedValue = value.join(', ');
  }

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        margin: '0.6rem 0',
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
      <p style={{ padding: '0', margin: '0', fontSize: '16px' }}>{formattedValue}</p>
    </div>
  );
};

const Title = ({ title }) => (
  <p
    style={{
      fontSize: '19px',
      fontWeight: 'bold',
    }}
  >
    {title}
  </p>
);

const ModelContents = ({ model }) => {
  const StyleClass = useStyles();
  const PropertyFormattersLeft = {
    version: (value) => <KeyValue property="API Version" value={value} />,
    hostname: (value) => <KeyValue property="Registrant" value={value} />,
    components: (value) => <KeyValue property="Components" value={value} />,
    subCategory: (value) => <KeyValue property="Sub Category" value={value} />,
  };

  const getCompRelValue = () => {
    let components = 0;
    let relationships = 0;

    if (model.versionBasedData) {
      model?.versionBasedData.forEach((model) => {
        components = components + (model?.components === null ? 0 : model.components.length);
        relationships =
          relationships + (model?.relationships === null ? 0 : model.relationships.length);
      });
    } else {
      components = model?.components === null ? 0 : model?.components?.length;
      relationships = model?.relationships === null ? 0 : model?.relationships?.length;
    }
    return {
      components,
      relationships,
    };
  };

  const metaDataLeft = {
    version: model?.version,
    hostname: model?.hostname,
    components: getCompRelValue()?.components?.toString(),
    subCategory: model?.subCategory,
  };

  const orderLeft = ['version', 'hostname', 'components', 'subCategory'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    category: (value) => <KeyValue property="Category" value={value} />,
    duplicates: (value) => <KeyValue property="Duplicates" value={value} />,
    relationships: (value) => <KeyValue property="Relationships" value={value} />,
  };

  const metaDataRight = {
    category: model?.category?.name,
    duplicates: model?.duplicates?.toString(),
    relationships: getCompRelValue().relationships.toString(),
  };

  const orderRight = ['category', 'duplicates', 'relationships'];
  const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);

  return (
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
  );
};

const ComponentContents = ({ component }) => {
  const StyleClass = useStyles();
  const PropertyFormattersLeft = {
    version: (value) => <KeyValue property="API Version" value={value} />,
    modelName: (value) => <KeyValue property="Model Name" value={value} />,
    kind: (value) => <KeyValue property="Kind" value={value} />,
    subCategory: (value) => <KeyValue property="Sub Category" value={value} />,
  };

  const metaDataLeft = {
    version: component?.apiVersion,
    modelName: component?.metadata?.modelDisplayName,
    kind: component?.kind,
    subCategory: component?.metadata?.subCategory,
  };

  const orderLeft = ['version', 'modelName', 'kind'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    registrant: (value) => <KeyValue property="Registrant" value={value} />,
    duplicates: (value) => <KeyValue property="Duplicates" value={value} />,
    category: (value) => <KeyValue property="Category" value={value} />,
  };

  const metaDataRight = {
    registrant: component?.displayhostname,
    duplicates: component?.duplicates?.toString(),
    category: component?.model?.category?.name,
  };

  const orderRight = ['registrant', 'duplicates'];
  const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);

  return (
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

const RegistrantContent = ({ registrant }) => {
  const PropertyFormattersLeft = {
    models: (value) => <KeyValue property="Models" value={value} />,
    components: (value) => <KeyValue property="Components" value={value} />,
  };

  const metaDataLeft = {
    models: registrant.summary?.models?.toString(),
    components: registrant.summary?.components?.toString(),
  };

  const orderLeft = ['models', 'components'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    relationships: (value) => <KeyValue property="Relationships" value={value} />,
    policies: (value) => <KeyValue property="Policies" value={value} />,
  };

  const metaDataRight = {
    relationships: registrant.summary?.relationships?.toString(),
    policies: registrant.summary?.policies?.toString(),
  };

  const orderRight = ['relationships', 'policies'];
  const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);
  return (
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
  );
};

const Description = ({ description }) => (
  <div style={{ margin: '0.6rem 0' }}>
    <p style={{ fontWeight: '600', margin: '0', fontSize: '16px' }}>Description</p>
    <p style={{ margin: '0', fontSize: '16px' }}>{description}</p>
  </div>
);

const TitleWithImg = ({ displayName, iconSrc }) => (
  <div style={{ display: 'flex', alignItems: 'center' }}>
    <img src={iconSrc} height="55px" width="55px" style={{ marginRight: '0.6rem' }} />
    <Title title={displayName} />
  </div>
);

const MeshModelDetails = ({ view, showDetailsData }) => {
  const StyleClass = useStyles();
  const isEmptyDetails =
    Object.keys(showDetailsData.data).length === 0 || showDetailsData.type === 'none';

  return (
    <div
      className={isEmptyDetails ? StyleClass.emptyDetailsContainer : StyleClass.detailsContainer}
    >
      {isEmptyDetails && <p style={{ color: '#969696', margin: 'auto' }}>No {view} selected</p>}
      {showDetailsData.type === MODELS && (
        <div>
          <TitleWithImg
            displayName={showDetailsData.data.displayName}
            iconSrc={showDetailsData.data.metadata?.svgColor}
          />
          <ModelContents model={showDetailsData.data} />
        </div>
      )}
      {showDetailsData.type === RELATIONSHIPS && (
        <div>
          <Title title={showDetailsData.data.kind} />
          <Description description={showDetailsData.data.metadata?.description} />
          <RelationshipContents relationship={showDetailsData.data} />
        </div>
      )}
      {showDetailsData.type === COMPONENTS && (
        <div>
          <TitleWithImg
            displayName={showDetailsData.data.displayName}
            iconSrc={showDetailsData.data.metadata?.svgColor}
          />
          {showDetailsData.data.schema && (
            <Description description={JSON.parse(showDetailsData.data.schema)?.description} />
          )}
          <ComponentContents component={showDetailsData.data} />
        </div>
      )}
      {showDetailsData.type === REGISTRANTS && (
        <div>
          <Title title={showDetailsData.data.hostname} />
          <RegistrantContent registrant={showDetailsData.data} />
        </div>
      )}
    </div>
  );
};

export default MeshModelDetails;
