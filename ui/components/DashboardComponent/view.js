import React from 'react';
import { ArrowBack } from '@material-ui/icons';
import TooltipButton from '../../utils/TooltipButton';
import { Paper, IconButton } from '@material-ui/core';
import { useRouter } from 'next/router';
// import { FormatStructuredData } from '../DataFormatter';
import NameValueTable from '../DataFormatter/NameValueTable';
import ResponsiveDataTable from '../../utils/data-table';

const View = (props) => {
  const { type, setView, resource, classes } = props;
  console.log(type, setView, resource);
  console.log('inside view component');

  // TODO: handle the condition for images
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
    };

    return (
      <>
        {key}

        <ResponsiveDataTable
          classes={classes.muiRow}
          data={value}
          columns={columns}
          options={options}
          tableCols={columns}
          updateCols={() => {}}
          columnVisibility={{}}
        />
      </>
    );
  }

  const RenderObject = (obj) => {
    function processObjForKeyValTable(obj, parentKey = '') {
      return Object.entries(obj).flatMap(([key, value]) => {
        const currentKey = parentKey ? `${parentKey}.${key}` : key;
        if (Array.isArray(value)) {
          return []; // Skip processing arrays
        }
        if (typeof value === 'object' && value !== null) {
          return processObjForKeyValTable(value, currentKey);
        } else {
          if (key === 'attribute') {
            return processObjForKeyValTable(JSON.parse(value), currentKey);
          }
          if (key === 'id') {
            return { name: currentKey, value: value, hide: true };
          }
          return { name: currentKey, value };
        }
      });
    }

    const processObjForKeyDataTable = (obj, parentKey = '') => {
      let results = [];
      for (const [key, value] of Object.entries(obj)) {
        console.log('key', key);
        console.log('value', value);
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
          results.push(processObjForKeyDataTable(value, currentKey));
        } else {
          if (key === 'attribute') {
            results.push(processObjForKeyDataTable(JSON.parse(value), currentKey));
          }
        }
      }
      return results;
    };

    return (
      <>
        {processObjForKeyDataTable(obj)}
        <NameValueTable rows={processObjForKeyValTable(obj)} />
      </>
    );
  };

  const router = useRouter();

  const HeaderComponent = () => {
    return (
      <>
        <TooltipButton title="Back" placement="left">
          <IconButton onClick={() => router.back()}>
            <ArrowBack />
          </IconButton>
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
          <div style={{ margin: '1rem 7rem' }}>
            <ResourceMetrics />
            <RenderObject obj={resource} />
          </div>
        </Paper>
      </div>
    </>
  );
};

export default View;
