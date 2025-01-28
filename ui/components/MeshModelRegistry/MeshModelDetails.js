import React from 'react';
import { DetailsContainer, Segment, FullWidth } from '@/assets/styles/general/tool.styles';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import { FormatStructuredData, reorderObjectProperties } from '../DataFormatter';
import {
  FormControl,
  Select,
  MenuItem,
  Chip,
  CircularProgress,
  useTheme,
  Button,
} from '@material-ui/core';
import DownloadIcon from '@mui/icons-material/Download';
import { withStyles } from '@material-ui/core/styles';
import styles from '../connections/styles';
import { REGISTRY_ITEM_STATES, REGISTRY_ITEM_STATES_TO_TRANSITION_MAP } from '../../utils/Enum';
import classNames from 'classnames';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import {
  useUpdateEntityStatusMutation,
  useGetComponentsQuery,
  useGetMeshModelsQuery,
} from '@/rtk-query/meshModel';
import _ from 'lodash';
import { JustifyAndAlignCenter } from './MeshModel.style';
import { reactJsonTheme } from './helper';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { Accordion, AccordionDetails, AccordionSummary, styled } from '@layer5/sistent';
import dynamic from 'next/dynamic';
import { UsesSistent } from '../SistentWrapper';
import {
  StyledKeyValueFormattedValue,
  StyledKeyValuePropertyDiv,
  StyledKeyValueProperty,
} from './MeshModel.style';
const ReactJson = dynamic(() => import('react-json-view'), { ssr: false });

const ExportAvailable = true;
const KeyValue = ({ property, value }) => {
  let formattedValue = value;

  if (Array.isArray(value)) {
    formattedValue = value.join(', ');
  }

  return (
    <StyledKeyValuePropertyDiv>
      <StyledKeyValueProperty>{property}</StyledKeyValueProperty>
      <StyledKeyValueFormattedValue>{formattedValue}</StyledKeyValueFormattedValue>
    </StyledKeyValuePropertyDiv>
  );
};

const StyledTitle = styled('div')(({ theme }) => ({
  fontSize: '1.25rem',
  fontFamily: theme.typography.fontFamily,
  textAlign: 'left',
  lineHeight: '1.3rem',
}));

const RenderContents = ({
  metaDataLeft,
  metaDataRight,
  PropertyFormattersLeft,
  PropertyFormattersRight,
  orderLeft,
  orderRight,
  jsonData,
}) => {
  const theme = useTheme();
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
      <Segment>
        <FullWidth style={{ display: 'flex', flexDirection: 'column', paddingRight: '1rem' }}>
          <FormatStructuredData
            data={reorderObjectProperties(metaDataLeft, orderLeft)}
            propertyFormatters={PropertyFormattersLeft}
            order={orderLeft}
          />
        </FullWidth>
        <FullWidth style={{ display: 'flex', flexDirection: 'column' }}>
          <FormatStructuredData
            data={reorderObjectProperties(metaDataRight, orderRight)}
            propertyFormatters={PropertyFormattersRight}
            order={orderRight}
          />
        </FullWidth>
      </Segment>
      {jsonData && (
        <Accordion
          style={{
            borderRadius: '6px',
            backgroundColor: theme.palette.secondary.toolbarBg2,
            color: theme.palette.secondary.text,
            margin: '0 -1rem',
            padding: '0',
          }}
        >
          <AccordionSummary
            expandIcon={<ExpandMoreIcon style={{ fill: theme.palette.secondary.text }} />}
          >
            Advanced Details
          </AccordionSummary>
          <AccordionDetails
            style={{
              padding: '0',
              fontSize: '0.85rem',
            }}
          >
            <ReactJson
              theme={reactJsonTheme(theme.palette.type)}
              name={false}
              displayDataTypes={false}
              iconStyle="circle"
              src={jsonData}
              style={{
                fontSize: '.85rem',
                minHeight: 'inherit',
                padding: '1.1rem',
                margin: '0rem',
              }}
              collapsed={1} // expanded upto 1 level
            />
          </AccordionDetails>
        </Accordion>
      )}
    </div>
  );
};

const ModelContents = ({ modelDef }) => {
  const PropertyFormattersLeft = {
    version: (value) => <KeyValue property="API Version" value={value} />,
    hostname: (value) => <KeyValue property="Registrant" value={value} />,
    components: (value) => <KeyValue property="Components" value={value} />,
    subCategory: (value) => <KeyValue property="Sub-Category" value={value} />,
    modelVersion: (value) => <KeyValue property="Model Version" value={value} />,
    registrant: (value) => <KeyValue property="Registrant" value={value} />,
  };

  const getCompRelValue = () => {
    let components = 0;
    let relationships = 0;
    if (modelDef?.versionBasedData) {
      modelDef?.versionBasedData.forEach((modelDefVersion) => {
        components = components + modelDefVersion?.components_count;
        relationships = relationships + modelDefVersion?.relationships_count;
      });
    } else {
      components = modelDef?.components_count;
      relationships = modelDef?.relationships_count;
    }
    return {
      components,
      relationships,
    };
  };

  const metaDataLeft = {
    version: modelDef?.model?.version,
    modelVersion: modelDef?.model?.modelVersion,
    hostname: modelDef?.registrant?.hostname,
    components: getCompRelValue()?.components?.toString(),
    subCategory: modelDef?.model?.subCategory,
    registrant: modelDef?.registrant?.name,
  };

  const orderLeft = ['version', 'hostname', 'components', 'subCategory'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    category: (value) => <KeyValue property="Category" value={value} />,
    duplicates: (value) => <KeyValue property="Duplicates" value={value} />,
    relationships: (value) => <KeyValue property="Relationships" value={value} />,
  };

  const metaDataRight = {
    category: modelDef?.category?.name,
    duplicates: modelDef?.duplicates?.toString(),
    relationships: getCompRelValue().relationships.toString(),
  };
  const handleExport = () => {
    const a = document.createElement('a');
    a.href = '/api/meshmodels/export?id=' + modelDef.id;
    document.body.appendChild(a);
    a.click();
    a.remove();
  };
  const orderRight = ['category', 'duplicates', 'relationships'];
  const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);
  const isShowStatusSelector = !Array.isArray(modelDef?.model.version);

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <TitleWithImg displayName={modelDef.displayName} iconSrc={modelDef?.metadata?.svgColor} />
        <div style={{ display: 'block' }}>
          {ExportAvailable ? (
            <Button
              aria-label="Export Model"
              variant="contained"
              color="primary"
              size="medium"
              alt="Export Model to OCI Image"
              onClick={handleExport}
              style={{ display: 'flex', width: '100%', marginBottom: '.25rem' }}
            >
              <DownloadIcon style={{ fontSize: '1.2rem' }} />
              Export
            </Button>
          ) : null}
          {isShowStatusSelector && <StatusChip entityData={modelDef} entityType="models" />}
        </div>
      </div>
      <RenderContents
        metaDataLeft={metaDataLeft}
        metaDataRight={metaDataRight}
        PropertyFormattersLeft={PropertyFormattersLeft}
        PropertyFormattersRight={PropertyFormattersRight}
        orderLeft={orderdMetadataLeft}
        orderRight={orderdMetadataRight}
        jsonData={modelDef}
      />
    </div>
  );
};

const ComponentContents = ({ componentDef }) => {
  const { data, isSuccess } = useGetComponentsQuery({
    params: {
      id: componentDef.id,
      apiVersion: componentDef.component.version,
      trim: false,
    },
  });
  const componentData = data?.components?.find((comp) => comp.id === componentDef.id);
  const PropertyFormattersLeft = {
    version: (value) => <KeyValue property="API Version" value={value} />,
    modelName: (value) => <KeyValue property="Model Name" value={value} />,
    kind: (value) => <KeyValue property="Kind" value={value} />,
    subCategory: (value) => <KeyValue property="Sub Category" value={value} />,
  };

  const metaDataLeft = {
    version: componentData?.component?.version,
    modelName: componentData?.model?.displayName,
    kind: componentData?.component.kind,
    subCategory: componentData?.model?.subCategory,
  };

  const orderLeft = ['version', 'modelName', 'kind'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    registrant: (value) => <KeyValue property="Registrant" value={value} />,
    duplicates: (value) => <KeyValue property="Duplicates" value={value} />,
    category: (value) => <KeyValue property="Category" value={value} />,
  };

  const metaDataRight = {
    registrant: componentData?.registrant?.hostname,
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
              iconSrc={componentData?.styles?.svgColor}
            />
          </div>
          <Description description={JSON.parse(componentData?.component?.schema)?.description} />
          <RenderContents
            metaDataLeft={metaDataLeft}
            metaDataRight={metaDataRight}
            PropertyFormattersLeft={PropertyFormattersLeft}
            PropertyFormattersRight={PropertyFormattersRight}
            orderLeft={orderdMetadataLeft}
            orderRight={orderdMetadataRight}
            jsonData={componentData}
          />
        </div>
      ) : (
        <JustifyAndAlignCenter>
          <CircularProgress size={24} />
        </JustifyAndAlignCenter>
      )}
    </>
  );
};

const RelationshipContents = ({ relationshipDef }) => {
  const PropertyFormattersLeft = {
    version: (value) => <KeyValue property="API Version" value={value} />,
    registrant: (value) => <KeyValue property="Registrant" value={value} />,
  };

  const metaDataLeft = {
    registrant: relationshipDef.model.registrant.name,
    modelName: relationshipDef.model?.displayName,
    version: relationshipDef.schemaVersion,
  };

  const orderLeft = ['registrant', 'version'];
  const orderdMetadataLeft = reorderObjectProperties(metaDataLeft, orderLeft);

  const PropertyFormattersRight = {
    registrant: (value) => <KeyValue property="Registrant" value={value} />,
    subType: (value) => <KeyValue property="Sub Type" value={value} />,
  };

  const metaDataRight = {
    registrant: relationshipDef.model.registrant.hostname,
    subType: relationshipDef.subType,
  };

  const orderRight = ['subType', 'registrant'];
  const orderdMetadataRight = reorderObjectProperties(metaDataRight, orderRight);

  return (
    <div>
      <div style={{ display: 'flex', flexDirection: 'column' }}>
        <StyledTitle>{`${relationshipDef?.kind} :: ${relationshipDef.type} :: ${relationshipDef.subType}`}</StyledTitle>
        <Description description={relationshipDef?.metadata?.description} />
      </div>
      <RenderContents
        metaDataLeft={metaDataLeft}
        metaDataRight={metaDataRight}
        PropertyFormattersLeft={PropertyFormattersLeft}
        PropertyFormattersRight={PropertyFormattersRight}
        orderLeft={orderdMetadataLeft}
        orderRight={orderdMetadataRight}
        jsonData={relationshipDef}
      />
    </div>
  );
};

const RegistrantContent = ({ registrant }) => {
  const PropertyFormattersLeft = {
    models: (value) => <KeyValue property="Models" value={value} />,
    components: (value) => <KeyValue property="Components" value={value} />,
  };

  const metaDataLeft = {
    models: registrant?.summary?.models?.toString(),
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
        <StyledTitle>{registrant?.hostname}</StyledTitle>
      </div>
      <RenderContents
        metaDataLeft={metaDataLeft}
        metaDataRight={metaDataRight}
        PropertyFormattersLeft={PropertyFormattersLeft}
        PropertyFormattersRight={PropertyFormattersRight}
        orderLeft={orderdMetadataLeft}
        orderRight={orderdMetadataRight}
        jsonData={registrant}
      />
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
  <div style={{ display: 'flex', alignItems: 'center', flexBasis: '60%' }}>
    {iconSrc && <img src={iconSrc} height="32px" width="32px" style={{ marginRight: '0.6rem' }} />}
    <StyledTitle>{displayName}</StyledTitle>
  </div>
);

// TODO: remove with styles and use styled component
const StatusChip = withStyles(styles)(({ classes, entityData, entityType }) => {
  const nextStatus = Object.values(REGISTRY_ITEM_STATES);
  const [updateEntityStatus] = useUpdateEntityStatusMutation();
  const { data: modelData, isSuccess } = useGetMeshModelsQuery({
    params: {
      id: entityData.model.id,
      version: entityData.model.version,
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
    [REGISTRY_ITEM_STATES_TO_TRANSITION_MAP.IGNORED]: () => <RemoveCircleIcon />,
    [REGISTRY_ITEM_STATES_TO_TRANSITION_MAP.ENABLED]: () => <AssignmentTurnedInIcon />,
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
                    : REGISTRY_ITEM_STATES_TO_TRANSITION_MAP?.[status] || status
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
});

const MeshModelDetails = ({ view, showDetailsData }) => {
  const theme = useTheme();
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
        return <ModelContents modelDef={showDetailsData.data} />;
      case RELATIONSHIPS:
        return <RelationshipContents relationshipDef={showDetailsData.data} />;
      case COMPONENTS:
        return <ComponentContents componentDef={showDetailsData.data} />;
      case REGISTRANTS:
        return <RegistrantContent registrant={showDetailsData.data} />;
      default:
        return null;
    }
  };

  return (
    <UsesSistent>
      <DetailsContainer isEmpty={isEmptyDetails}>
        {isEmptyDetails ? renderEmptyDetails() : getContent(showDetailsData.type)}
      </DetailsContainer>
    </UsesSistent>
  );
};

export default MeshModelDetails;
