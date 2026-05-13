import React, { useContext } from 'react';
import { reactJsonTheme } from '../registry/helper';
import dynamic from 'next/dynamic';
import { SectionBody, ArrayFormatter, SectionHeading } from '../data-formatter';
import _ from 'lodash';
import { Grid2, useTheme, Typography, styled, Box } from '@sistent/sistent';

type FormatterMap = Record<string, (value: any) => React.ReactNode>;
type ResourceFormatterProps = {
  data: any;
};

const splitCamelCaseString = (value: string) => value.replace(/([a-z0-9])([A-Z])/g, '$1 $2');

const ReactJson = dynamic(() => import('@microlink/react-json-view'), { ssr: false });
const FormatterContext = React.createContext<{ propertyFormatters: FormatterMap }>({
  propertyFormatters: {},
});
const LevelContext = React.createContext(0);

export const ColourContainer = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.background.card,
  padding: '1rem',
  [theme.breakpoints.down(599)]: {
    width: '60vw',
  },
  width: '100%',
  overflow: 'scroll',
}));

export const JSONViewFormatter = ({ data }: ResourceFormatterProps) => {
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

type KeyValueInRowProps = {
  Key: string;
  Value: React.ReactNode;
  showFold?: boolean;
};

type TableDataFormatterProps = {
  data?: unknown;
  mainTableData?: unknown;
  mainTableCols?: unknown;
};

const formatDisplayValue = (value: unknown) =>
  typeof value === 'string' ? value : JSON.stringify(value, null, 2);

const KeyValueInRow = ({ Key, Value }: KeyValueInRowProps) => (
  <Box sx={{ display: 'flex', flexDirection: 'column', gap: '0.35rem' }}>
    <Typography variant="body2" sx={{ fontWeight: 600 }}>
      {Key}
    </Typography>
    {Value}
  </Box>
);

const StatusFormatter = ({ status }: { status: unknown }) => (
  <Typography variant="body2">{formatDisplayValue(status)}</Typography>
);

const MemoryUsage = ({ allocatable, capacity }: { allocatable: unknown; capacity: unknown }) => (
  <JSONViewFormatter data={{ allocatable, capacity }} />
);

const TextWithLinkFormatter = ({
  title,
  value,
  onClick,
}: {
  title: string;
  value: React.ReactNode;
  onClick: () => void;
}) => (
  <Box sx={{ display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
    <Typography variant="body2" sx={{ fontWeight: 600 }}>
      {title}
    </Typography>
    <Typography
      variant="body2"
      sx={{ textDecoration: 'underline', cursor: 'pointer' }}
      onClick={onClick}
    >
      {value}
    </Typography>
  </Box>
);

const CodeFormatter = ({ data }: { data: unknown }) => (
  <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{formatDisplayValue(data)}</pre>
);

const LabelFormatter = ({ data }: { data: unknown }) => <JSONViewFormatter data={data} />;

const TableDataFormatter = ({ data, mainTableData, mainTableCols }: TableDataFormatterProps) => (
  <JSONViewFormatter
    data={mainTableData ? { rows: mainTableData, columns: mainTableCols } : data}
  />
);

const ListFormatter = ({ data }: { data: unknown }) =>
  Array.isArray(data) ? (
    <ArrayFormatter items={data} style={{}} />
  ) : (
    <JSONViewFormatter data={data} />
  );

const NumberStateFormatter = ({ data }: { data: unknown }) => <JSONViewFormatter data={data} />;

const CollapsibleSectionFormatter = ({
  title,
  children,
}: {
  title: React.ReactNode;
  children: React.ReactNode;
}) => (
  <details>
    <summary>{title}</summary>
    <Box sx={{ mt: 1 }}>{children}</Box>
  </details>
);

const ContainerFormatter = ({
  containerSpec,
  containerStatus,
}: {
  containerSpec: unknown;
  containerStatus: unknown;
}) => <JSONViewFormatter data={{ containerSpec, containerStatus }} />;

const OperatorDynamicFormatter = ({ data }: { data: unknown }) => <JSONViewFormatter data={data} />;

const SecretFormatter = ({ data }: { data: unknown }) => <JSONViewFormatter data={data} />;

const extractPodVolumnTables = (value: Record<string, unknown> | null | undefined) => {
  if (!value || typeof value !== 'object') {
    return [];
  }

  return Object.entries(value).map(([key, rows]) => ({
    key,
    rows,
    columns:
      Array.isArray(rows) && rows.length > 0 && typeof rows[0] === 'object' && rows[0] !== null
        ? Object.keys(rows[0] as Record<string, unknown>)
        : [],
  }));
};

const Level = ({ children }: { children: React.ReactNode }) => {
  const level = useContext(LevelContext);
  return <LevelContext.Provider value={level + 1}> {children} </LevelContext.Provider>;
};

const ResourceDynamicFormatter = ({ data }: ResourceFormatterProps) => {
  const { propertyFormatters } = useContext(FormatterContext);
  const level = useContext(LevelContext);

  if (_.isString(data)) {
    return <SectionBody body={data}></SectionBody>;
  }

  if (_.isArray(data)) {
    return <ArrayFormatter items={data} style={{}} />;
  }

  if (_.isObject(data)) {
    return Object.entries(data).map(([title, data]) => {
      if (!title.trim() || !data || _.isEmpty(data)) {
        return null;
      }

      if (propertyFormatters?.[title]) {
        return (
          <Grid2 key={title} size={{ sm: 12 }}>
            {propertyFormatters[title](data)}
          </Grid2>
        );
      }

      if (typeof data == 'string') {
        return <KeyValueInRow key={title} Key={title} Value={data} />;
      }

      return (
        <Grid2
          key={title}
          size={{
            sm: 12,
          }}
          style={{
            marginBlock: '0.25rem',
          }}
        >
          <SectionHeading level={level} isLevel>
            {title}
          </SectionHeading>
          <Level>
            <ResourceDynamicFormatter data={data} />
          </Level>
        </Grid2>
      );
    });
  }

  return null;
};

const propertyFormatter: FormatterMap = {
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
        {value.links.map((linkObj: any) => {
          const { label, nodeName, namespace, serviceAccount, resourceCategory } = linkObj;
          const name = nodeName || namespace || serviceAccount;
          if (!name) return null;
          return (
            <TextWithLinkFormatter
              key={label}
              title={label}
              value={name}
              onClick={() => {
                return value.router.push(
                  {
                    pathname: value.router.pathname,
                    query: {
                      resourceCategory: resourceCategory || label,
                      resourceName: name,
                    },
                  },
                  undefined,
                  { shallow: true },
                );
              }}
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
      Value={<LabelFormatter data={value?.data} />}
      showFold={value?.data?.length > 7}
    />
  ),
  annotations: (value: string | any[]) => (
    <KeyValueInRow
      Key={'Annotations'}
      Value={<StatusFormatter status={value} />}
      showFold={value?.length > 7}
    />
  ),
  totalCapacity: (value) => {
    // const readableData = Object.fromEntries(
    // Object.entries(value).map(([key, val]) => [key, convertToReadableUnit(parseInt(val))]),
    // );
    const readableData = value;
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
    // const readableData = Object.fromEntries(
    //   Object.entries(value).map(([key, val]) => [key, convertToReadableUnit(parseInt(val))]),
    // );

    const readableData = value;
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
            {spec.map((containerSpec: any, index: number) => (
              <Box
                key={index}
                sx={{ px: 1, display: 'flex', flexDirection: 'column', gap: '0.3rem' }}
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
            {spec.map((containerSpec: any, index: number) => (
              <Box
                key={index}
                sx={{ px: 1, display: 'flex', flexDirection: 'column', gap: '0.3rem' }}
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
            {tables.map((table: any, index: number) => {
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
const ResourceDetailFormatData = ({ data }: ResourceFormatterProps) => {
  const formatterContextValue = React.useMemo(
    () => ({
      propertyFormatters: propertyFormatter,
    }),
    [],
  );

  return (
    <FormatterContext.Provider value={formatterContextValue}>
      <Grid2
        container
        style={{
          overflowWrap: 'break-word',
          gap: '0.3rem 1rem',
          flexDirection: 'column',
        }}
      >
        <ResourceDynamicFormatter data={data} />
      </Grid2>
    </FormatterContext.Provider>
  );
};

export default ResourceDetailFormatData;
