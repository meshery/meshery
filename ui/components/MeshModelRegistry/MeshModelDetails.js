import React from 'react';
import useStyles from '../../assets/styles/general/tool.styles';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import { FormatStructuredData, reorderObjectProperties } from '../DataFormatter';
import { FormControl, Select, MenuItem, Chip, CircularProgress, useTheme } from '@material-ui/core';
import { withStyles } from '@material-ui/core/styles';
import styles from '../connections/styles';
import { CONNECTION_STATES, CONNECTION_STATE_TO_TRANSITION_MAP } from '../../utils/Enum';
import classNames from 'classnames';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import {
  useUpdateEntityStatusMutation,
  useGetModelByNameQuery,
  useGetComponentByNameQuery,
} from '@/rtk-query/meshModel';
import _ from 'lodash';
import { JustifyAndAlignCenter } from './MeshModel.style';
import { withSuppressedErrorBoundary } from '../General/ErrorBoundary';

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

const RenderContents = ({
  metaDataLeft,
  metaDataRight,
  PropertyFormattersLeft,
  PropertyFormattersRight,
  orderLeft,
  orderRight,
}) => {
  const StyleClass = useStyles();

  return (
    <div className={StyleClass.segment}>
      <div
        className={StyleClass.fullWidth}
        style={{ display: 'flex', flexDirection: 'column', paddingRight: '1rem' }}
      >
        <FormatStructuredData
          data={reorderObjectProperties(metaDataLeft, orderLeft)}
          propertyFormatters={PropertyFormattersLeft}
          order={orderLeft}
        />
      </div>

      <div className={StyleClass.fullWidth} style={{ display: 'flex', flexDirection: 'column' }}>
        <FormatStructuredData
          data={reorderObjectProperties(metaDataRight, orderRight)}
          propertyFormatters={PropertyFormattersRight}
          order={orderRight}
        />
      </div>
    </div>
  );
};

const ModelContents = withSuppressedErrorBoundary(({ model }) => {
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
  const isShowStatusSelector = !Array.isArray(model.version);

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between' }}>
        <TitleWithImg displayName={model.displayName} iconSrc={model?.metadata?.svgColor} />
        {isShowStatusSelector && <StatusChip entityData={model} entityType="models" />}
      </div>
      <RenderContents
        metaDataLeft={metaDataLeft}
        metaDataRight={metaDataRight}
        PropertyFormattersLeft={PropertyFormattersLeft}
        PropertyFormattersRight={PropertyFormattersRight}
        orderLeft={orderdMetadataLeft}
        orderRight={orderdMetadataRight}
      />
    </div>
  );
});

const ComponentContents = withSuppressedErrorBoundary(({ component }) => {
  const { data, isSuccess } = useGetComponentByNameQuery({
    name: component.kind,
    params: {
      apiVersion: component.apiVersion,
      trim: false,
    },
  });
  const componentData = data?.components?.find((comp) => comp.id === component.id);
  const PropertyFormattersLeft = {
    version: (value) => <KeyValue property="API Version" value={value} />,
    modelName: (value) => <KeyValue property="Model Name" value={value} />,
    kind: (value) => <KeyValue property="Kind" value={value} />,
    subCategory: (value) => <KeyValue property="Sub Category" value={value} />,
  };

  const metaDataLeft = {
    version: componentData?.apiVersion,
    modelName: componentData?.metadata?.modelDisplayName,
    kind: componentData?.kind,
    subCategory: componentData?.metadata?.subCategory,
  };

  const orderLeft = ['version', 'modelName', 'kind'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    registrant: (value) => <KeyValue property="Registrant" value={value} />,
    duplicates: (value) => <KeyValue property="Duplicates" value={value} />,
    category: (value) => <KeyValue property="Category" value={value} />,
  };

  const metaDataRight = {
    registrant: component?.hostname,
    duplicates: componentData?.duplicates?.toString(),
    category: componentData?.model?.category?.name,
  };

  const orderRight = ['registrant', 'duplicates'];
  const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);

  return (
    <>
      {isSuccess && componentData ? (
        <div>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <TitleWithImg
              displayName={componentData?.displayName}
              iconSrc={componentData?.metadata?.svgColor}
            />
          </div>
          <Description description={JSON.parse(componentData?.schema)?.description} />
          <RenderContents
            metaDataLeft={metaDataLeft}
            metaDataRight={metaDataRight}
            PropertyFormattersLeft={PropertyFormattersLeft}
            PropertyFormattersRight={PropertyFormattersRight}
            orderLeft={orderdMetadataLeft}
            orderRight={orderdMetadataRight}
          />
        </div>
      ) : (
        <JustifyAndAlignCenter>
          <CircularProgress size={24} />
        </JustifyAndAlignCenter>
      )}
    </>
  );
});

const RelationshipContents = withSuppressedErrorBoundary(({ relationship }) => {
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
    <div>
      <div style={{ display: 'flex', flexDirection: 'column' }}>
        <Title title={relationship.subType} />
        <Description description={relationship.metadata.description} />
      </div>
      <RenderContents
        metaDataLeft={metaDataLeft}
        metaDataRight={metaDataRight}
        PropertyFormattersLeft={PropertyFormattersLeft}
        PropertyFormattersRight={PropertyFormattersRight}
        orderLeft={orderdMetadataLeft}
        orderRight={orderdMetadataRight}
      />
    </div>
  );
});

const RegistrantContent = withSuppressedErrorBoundary(({ registrant }) => {
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
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between' }}>
        <Title title={registrant.hostname} />
      </div>
      <RenderContents
        metaDataLeft={metaDataLeft}
        metaDataRight={metaDataRight}
        PropertyFormattersLeft={PropertyFormattersLeft}
        PropertyFormattersRight={PropertyFormattersRight}
        orderLeft={orderdMetadataLeft}
        orderRight={orderdMetadataRight}
      />
    </div>
  );
});

const Description = ({ description }) => (
  <div style={{ margin: '0.6rem 0' }}>
    <p style={{ fontWeight: '600', margin: '0', fontSize: '16px' }}>Description</p>
    <p style={{ margin: '0', fontSize: '16px' }}>{description}</p>
  </div>
);

const TitleWithImg = ({ displayName, iconSrc }) => (
  <div style={{ display: 'flex', alignItems: 'center' }}>
    {iconSrc && <img src={iconSrc} height="55px" width="55px" style={{ marginRight: '0.6rem' }} />}
    <Title title={displayName} />
  </div>
);

// TODO: remove with styles and use either makestyle or styled component
const StatusChip = withSuppressedErrorBoundary(
  withStyles(styles)(({ classes, entityData, entityType }) => {
    const nextStatus = ['registered', 'ignored'];
    const [updateEntityStatus] = useUpdateEntityStatusMutation();
    const { data: modelData, isSuccess } = useGetModelByNameQuery({
      name: entityData.name,
      params: {
        version: entityData.version,
      },
    });

    const data = modelData?.models?.find((model) => model.id === entityData.id);

    const handleStatusChange = (e) => {
      updateEntityStatus({
        entityType: _.toLower(entityType),
        body: {
          id: data.id,
          status: e.target.value,
          displayname: entityData.displayName,
        },
      });
    };

    const icons = {
      [CONNECTION_STATES.IGNORED]: () => <RemoveCircleIcon />,
      [CONNECTION_STATES.REGISTERED]: () => <AssignmentTurnedInIcon />,
    };

    return (
      <FormControl
        className={classes.chipFormControl}
        style={{ minWidth: '0%', flexDirection: 'inherit' }}
      >
        {isSuccess ? (
          <Select
            labelId="entity-status-select-label"
            id={data?.id}
            key={data?.id}
            value={data?.status}
            defaultValue={data?.status}
            onClick={(e) => e.stopPropagation()}
            onChange={(e) => handleStatusChange(e)}
            className={classes.statusSelect}
            disableUnderline
            disabled={!isSuccess} // Disable the select when isSuccess is false
            MenuProps={{
              anchorOrigin: {
                vertical: 'bottom',
                horizontal: 'left',
              },
              transformOrigin: {
                vertical: 'top',
                horizontal: 'left',
              },
              getContentAnchorEl: null,
              MenuListProps: { disablePadding: true },
              PaperProps: { square: true },
            }}
          >
            {nextStatus.map((status) => (
              <MenuItem
                disabled={status === data?.status}
                style={{ padding: '0', display: status === data?.status ? 'none' : 'flex' }}
                value={status}
                key={status}
              >
                <Chip
                  className={classNames(classes.statusChip, classes[status])}
                  avatar={icons[status] ? icons[status]() : ''}
                  label={
                    status === data?.status
                      ? status
                      : CONNECTION_STATE_TO_TRANSITION_MAP?.[status] || status
                  }
                />
              </MenuItem>
            ))}
          </Select>
        ) : (
          <CircularProgress size={24} />
        )}
      </FormControl>
    );
  }),
);

const MeshModelDetails = ({ view, showDetailsData }) => {
  const theme = useTheme();
  const StyleClass = useStyles();
  const isEmptyDetails =
    Object.keys(showDetailsData.data).length === 0 || showDetailsData.type === 'none';

  const renderEmptyDetails = () => (
    <p style={{ color: theme.palette.secondary.menuItemBorder, margin: 'auto' }}>
      No {view} selected
    </p>
  );

  const getContent = (type) => {
    switch (type) {
      case MODELS:
        return <ModelContents model={showDetailsData.data} />;
      case RELATIONSHIPS:
        return <RelationshipContents relationship={showDetailsData.data} />;
      case COMPONENTS:
        return <ComponentContents component={showDetailsData.data} />;
      case REGISTRANTS:
        return <RegistrantContent registrant={showDetailsData.data} />;
      default:
        return null;
    }
  };

  return (
    <div
      className={isEmptyDetails ? StyleClass.emptyDetailsContainer : StyleClass.detailsContainer}
    >
      {isEmptyDetails ? renderEmptyDetails() : getContent(showDetailsData.type)}
    </div>
  );
};

export default MeshModelDetails;
