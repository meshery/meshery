import React, { useContext } from 'react';
import { reactJsonTheme } from '../MeshModelRegistry/helper';
import dynamic from 'next/dynamic';
import { SectionBody, ArrayFormatter } from '../DataFormatter';
import _ from 'lodash';
import {
  Grid,
  useTheme,
  Typography,
  StatusFormatter,
  MemoryUsage,
  TextWithLinkFormatter,
  CodeFormatter,
  LabelFormatter,
  TableDataFormatter,
  styled,
  ListFormatter,
  NumberStateFormatter,
  CollapsibleSectionFormatter,
  ContainerFormatter,
  Box,
  OperatorDynamicFormatter,
  SecretFormatter,
  extractPodVolumnTables,
  splitCamelCaseString,
  KeyValueInRow,
  convertToReadableUnit,
} from '@layer5/sistent';
import { SectionHeading } from '../DataFormatter';
import { UsesSistent } from '../SistentWrapper';

const ReactJson = dynamic(() => import('react-json-view'), { ssr: false });
const FormatterContext = React.createContext({
  propertyFormatters: {},
});
const LevelContext = React.createContext(0);

export const ColourContainer = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#212121' : '#e9eff1',
  padding: '1rem',
}));

export const JSONViewFormatter = ({ data }) => {
  const theme = useTheme();
  const rjvTheme = reactJsonTheme(theme.palette.mode);

  return (
    <ReactJson
      theme={rjvTheme}
      name={false}
      iconStyle="circle"
      src={data}
      style={{
        fontSize: 'inherit',
        minHeight: 'inherit',
        padding: '1.1rem',
      }}
      displayObjectSize={false}
      collapsed={3}
      displayDataTypes={false}
    />
  );
};

const Level = ({ children }) => {
  const level = useContext(LevelContext);
  return <LevelContext.Provider value={level + 1}> {children} </LevelContext.Provider>;
};

const ResourceDynamicFormatter = ({ data }) => {
  const { propertyFormatters } = useContext(FormatterContext);
  const level = useContext(LevelContext);

  if (_.isString(data)) {
    return <SectionBody body={data}></SectionBody>;
  }

  if (_.isArray(data)) {
    return <ArrayFormatter items={data} />;
  }

  if (_.isObject(data)) {
    return Object.entries(data).map(([title, data]) => {
      if (!title.trim() || !data || _.isEmpty(data)) {
        return null;
      }

      if (propertyFormatters?.[title]) {
        return (
          <Grid key={title} sm={12}>
            {propertyFormatters[title](data, data)}
          </Grid>
        );
      }

      if (typeof data == 'string') {
        return <KeyValueInRow key={title} Key={title} Value={data} />;
      }

      return (
        <Grid
          item
          key={title}
          sm={12}
          style={{
            marginBlock: '0.25rem',
          }}
        >
          <SectionHeading level={level}>{title}</SectionHeading>
          <Level>
            <ResourceDynamicFormatter data={data} />
          </Level>
        </Grid>
      );
    });
  }

  return null;
};

const propertyFormatter = {
  status: (status) => <KeyValueInRow Key={'Status'} Value={<StatusFormatter status={status} />} />,
  usage: (value) => {
    if (value.allocatable && value.capacity) {
      return (
        <KeyValueInRow
          Key={'Resource Usage'}
          Value={<MemoryUsage allocatable={value.allocatable} capacity={value.capacity} />}
        />
      );
    }
    return null;
  },
  deeplinks: (value) => {
    return (
      <>
        {value.links.map((linkObj) => {
          const { label, nodeName, namespace, serviceAccount } = linkObj;
          const name = nodeName || namespace || serviceAccount;
          if (!name) return null;
          return (
            <TextWithLinkFormatter
              key={label}
              title={label}
              value={name}
              // onClick={() => {
              //   return value.router.push(
              //     {
              //       pathname: value.router.pathname,
              //       query: {
              //         resourceCategory: resourceCategory || label,
              //         resourceName: name,
              //       },
              //     },
              //     undefined,
              //     { shallow: true },
              //   );
              // }}
            />
          );
        })}
      </>
    );
  },
  conditions: (value) => (
    <KeyValueInRow Key={'Conditions'} Value={<StatusFormatter status={value} />} />
  ),
  configData: (value) => (
    <KeyValueInRow Key={'Config Data'} Value={<CodeFormatter data={value} />} />
  ),
  labels: (value) => (
    <KeyValueInRow
      Key={'Labels'}
      Value={<LabelFormatter data={value?.data} selectedLabels={[]} />}
    />
  ),
  annotations: (value) => (
    <KeyValueInRow Key={'Annotations'} Value={<StatusFormatter status={value} />} />
  ),
  totalCapacity: (value) => {
    const readableData = Object.fromEntries(
      Object.entries(value).map(([key, val]) => [key, convertToReadableUnit(parseInt(val))]),
    );
    return (
      <KeyValueInRow
        Key={'Capacity'}
        Value={
          <ColourContainer>
            <TableDataFormatter data={readableData} />
          </ColourContainer>
        }
      />
    );
  },
  totalAllocatable: (value) => {
    const readableData = Object.fromEntries(
      Object.entries(value).map(([key, val]) => [key, convertToReadableUnit(parseInt(val))]),
    );
    return (
      <KeyValueInRow
        Key={'Allocatable'}
        Value={
          <ColourContainer>
            <TableDataFormatter data={readableData} />
          </ColourContainer>
        }
      />
    );
  },
  tolerations: (value) => (
    <KeyValueInRow
      Key={'Tolerations'}
      Value={
        <ColourContainer>
          <TableDataFormatter data={value} />
        </ColourContainer>
      }
    />
  ),
  nodeSelector: (value) => (
    <KeyValueInRow Key={'Node Selector'} Value={<ListFormatter data={value} />} />
  ),
  selector: (value) => <KeyValueInRow Key={'Selector'} Value={<ListFormatter data={value} />} />,
  images: (value) => <KeyValueInRow Key={'Images'} Value={<ListFormatter data={value} />} />,
  finalizers: (value) => (
    <KeyValueInRow Key={'Finalizers'} Value={<ListFormatter data={value} />} />
  ),
  accessModes: (value) => (
    <KeyValueInRow Key={'Access Modes'} Value={<ListFormatter data={value} />} />
  ),
  loadBalancer: (value) => (
    <KeyValueInRow Key={'Load Balancer'} Value={<ListFormatter data={value} />} />
  ),
  rules: (value) => <KeyValueInRow Key={'Rules'} Value={<ListFormatter data={value} />} />,
  numberStates: (value) => {
    return <KeyValueInRow Key={'Number States'} Value={<NumberStateFormatter data={value} />} />;
  },
  initContainers: (value) => {
    const spec = value.spec;
    const status = value.status;
    return (
      <KeyValueInRow
        Key={'Init Containers'}
        Value={
          <>
            {spec.map((containerSpec, index) => (
              <Box
                paddingInline={1}
                key={index}
                display={'flex'}
                flexDirection={'column'}
                gap={'0.3rem'}
              >
                <CollapsibleSectionFormatter title={containerSpec.name}>
                  <ContainerFormatter containerSpec={spec[index]} containerStatus={status[index]} />
                </CollapsibleSectionFormatter>
              </Box>
            ))}
          </>
        }
      />
    );
  },
  containers: (value) => {
    const spec = value.spec;
    const status = value.status;
    return (
      <KeyValueInRow
        Key={'Containers'}
        Value={
          <>
            {spec.map((containerSpec, index) => (
              <Box
                paddingInline={1}
                key={index}
                display={'flex'}
                flexDirection={'column'}
                gap={'0.3rem'}
              >
                <CollapsibleSectionFormatter title={containerSpec.name}>
                  <ContainerFormatter containerSpec={spec[index]} containerStatus={status[index]} />
                </CollapsibleSectionFormatter>
              </Box>
            ))}
          </>
        }
      />
    );
  },
  podVolumes: (value) => {
    const tables = extractPodVolumnTables(value);
    return (
      <KeyValueInRow
        Key={'Pod Volumes'}
        Value={
          <>
            {tables.map((table, index) => {
              return (
                <ColourContainer style={{ paddingInline: '1rem' }} key={index}>
                  <Typography variant="body1">{splitCamelCaseString(table.key)}</Typography>
                  <TableDataFormatter mainTableData={table.rows} mainTableCols={table.columns} />
                </ColourContainer>
              );
            })}
          </>
        }
      />
    );
  },
  connections: (value) => {
    return <KeyValueInRow Key={'Connections'} Value={<JSONViewFormatter data={value} />} />;
  },
  commands: (value) => {
    return <OperatorDynamicFormatter data={value} />;
  },
  ingressRules: (value) => {
    return <KeyValueInRow Key={'Ingress Rules'} Value={<JSONViewFormatter data={value} />} />;
  },
  secret: (value) => {
    return <KeyValueInRow Key={'Secret'} Value={<SecretFormatter data={value} />} />;
  },
};
const ResourceDetailFormatData = ({ data }) => {
  return (
    <UsesSistent>
      <FormatterContext.Provider
        value={{
          propertyFormatters: propertyFormatter,
        }}
      >
        <Grid
          container
          style={{
            wordBreak: 'break-word',
            overflowWrap: 'break-word',
            gap: '0.3rem 1rem',
          }}
        >
          <ResourceDynamicFormatter data={data} />
        </Grid>
      </FormatterContext.Provider>
    </UsesSistent>
  );
};

export default ResourceDetailFormatData;
