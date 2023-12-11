import React, { useState } from 'react';
import { ArrowBack } from '@material-ui/icons';
import TooltipButton from '../../utils/TooltipButton';
import { Paper, Typography } from '@material-ui/core';
import NameValueTable from '../DataFormatter/NameValueTable';
import { ResponsiveDataTable } from '@layer5/sistent-components';
import { ALL_VIEW } from './resources/config';
import GetNodeIcon from '../configuratorComponents/MeshModel/NodeIcon';
import { JsonParse } from '../../utils/utils';

const View = (props) => {
  const {
    // type,
    setView,
    resource,
    classes,
  } = props;

  function RenderDynamicTable(key, value) {
    const allKeys = value.reduce((keys, obj) => {
      Object.keys(obj).forEach((key) => {
        if (!keys.includes(key)) {
          keys.push(key);
        }
      });
      return keys;
    }, []);

    const columns = allKeys.map((key) => ({
      name: key,
      label: key,
      options: {
        filter: false,
        sort: false,
        display: key == 'id' ? false : true,
        customBodyRender: function CustomBody(value) {
          let parentStyle = {
            position: 'absolute',
            top: 0,
            right: 0,
            bottom: 0,
            left: 0,
            boxSizing: 'border-box',
            display: 'block',
            width: '100%',
          };
          let cellStyle = {
            boxSizing: 'border-box',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          };
          return (
            <>
              <div style={{ position: 'relative', height: '20px' }}>
                <div style={parentStyle}>
                  <div style={cellStyle}>
                    {typeof value === 'object' && value !== null ? JSON.stringify(value) : value}
                  </div>
                </div>
              </div>
            </>
          );
        },
      },
    }));

    let options = {
      filter: false,
      download: false,
      print: false,
      search: false,
      viewColumns: false,
      selectableRows: 'none',
      pagination: false,
      responsive: 'standard',
      fixedHeader: true,
      resizableColumns: true,
    };

    return (
      <>
        <div style={{ margin: '2rem 0' }}>
          <Typography style={{ fontSize: '1.2rem', marginBottom: '1rem' }} align="left">
            {key.toUpperCase()}
          </Typography>

          <ResponsiveDataTable
            classes={classes.muiRow}
            data={value}
            columns={columns}
            options={options}
            tableCols={columns}
            updateCols={() => {}}
            columnVisibility={{}}
          />
        </div>
      </>
    );
  }

  const RenderObject = (obj) => {
    function ProcessObjForKeyValTable(obj) {
      const [processedData, setProcessedData] = React.useState([]);

      function processObj(obj, parentKey = '') {
        let rows = [];
        let currentGroup = [];

        for (const [key, value] of Object.entries(obj)) {
          const currentKey = parentKey ? `${parentKey}.${key}` : key;

          if (Array.isArray(value)) {
            // Skip the key if the value is an array
            continue;
          } else if (typeof value === 'object' && value !== null) {
            // For objects, recursively process and add to the current group
            currentGroup.push(...processObj(value, currentKey));
          } else {
            // For non-objects, add to the rows directly
            if (key === 'attribute') {
              currentGroup.push(...processObj(JSON.parse(value), currentKey));
            } else if (key === 'id') {
              currentGroup.push({ name: currentKey, value: value, hide: true });
            } else {
              currentGroup.push({ name: currentKey, value });
            }
          }
        }

        // Group by the parent key
        if (parentKey !== '' && currentGroup.length > 0) {
          if (Array.isArray(currentGroup)) {
            setProcessedData((prev) => [...prev, { [parentKey]: currentGroup }]);
          }
        }

        return rows;
      }

      React.useEffect(() => {
        processObj(obj);
      }, [obj]);

      return (
        <>
          {processedData.map((obj, index) => (
            <div key={index}>
              {Object.entries(obj).map(([key, value], innerIndex) => {
                const parts = key.split('.');
                const lastPart = parts[parts.length - 1];
                const heading = lastPart.replace('_', ' ');
                return value.length == 1 && value[0].hide == true ? null : (
                  <div style={{ margin: '2rem 0' }} key={innerIndex}>
                    <Typography
                      style={{
                        fontSize: '.9rem',
                        marginBottom: '1rem',
                        textTransform: 'uppercase',
                      }}
                      align="left"
                    >
                      {heading}
                    </Typography>
                    <NameValueTable rows={value} />
                  </div>
                );
              })}
            </div>
          ))}
        </>
      );
    }

    const ProcessObjForKeyDataTable = (obj, parentKey = '') => {
      let results = [];
      for (const [key, value] of Object.entries(obj)) {
        const currentKey = parentKey ? `${parentKey}.${key}` : key;
        if (
          Array.isArray(value) &&
          value.length > 0 &&
          typeof value[0] === 'object' &&
          value[0] !== null
        ) {
          results.push(RenderDynamicTable(key, value));
        }
        if (typeof value === 'object' && value !== null) {
          results.push(ProcessObjForKeyDataTable(value, currentKey));
        } else {
          if (key === 'attribute') {
            results.push(ProcessObjForKeyDataTable(JSON.parse(value), currentKey));
          }
        }
      }
      return results;
    };

    return (
      <>
        {ProcessObjForKeyValTable(obj)}
        {ProcessObjForKeyDataTable(obj)}
      </>
    );
  };

  const HeaderComponent = () => {
    return (
      <>
        <TooltipButton title="Back" placement="left" style={{ padding: '10px' }}>
          <ArrowBack onClick={() => setView(ALL_VIEW)} />
        </TooltipButton>
      </>
    );
  };

  const ResourceMetrics = () => {
    return <></>;
  };

  return (
    <>
      <div
        style={{
          margin: '1rem auto',
        }}
      >
        <Paper>
          <HeaderComponent />
          <div style={{ margin: '1rem 7rem', paddingBottom: '1rem' }}>
            <ResourceMetrics />
            <RenderObject obj={resource} />
          </div>
        </Paper>
      </div>
    </>
  );
};

export default View;

export const Title = ({ onClick, data, value }) => {
  const [isHovered, setHovered] = useState(false);
  return (
    <div
      style={{
        color: 'inherit',
        cursor: 'pointer',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        textDecoration: isHovered ? 'underline' : 'none',
      }}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <div onClick={onClick} style={{ display: 'inherit', alignItems: 'center' }}>
        <GetNodeIcon metadata={JsonParse(data)} />
        <Typography style={{ marginLeft: '0.25rem' }} variant="body2">
          {value}
        </Typography>
      </div>
    </div>
  );
};
