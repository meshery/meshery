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
      // return Object.entries(obj).flatMap(([key, value]) => {
      //   const currentKey = parentKey ? `${parentKey}.${key}` : key;
      //   if (Array.isArray(value)) {
      //     return []; // Skip processing arrays
      //   }
      //   if (typeof value === 'object' && value !== null) {
      //     return processObjForKeyValTable(value, currentKey);
      //   } else {
      //     if (key === 'attribute') {
      //       return processObjForKeyValTable(JSON.parse(value), currentKey);
      //     }
      //     if (key === 'id') {
      //       return { name: currentKey, value: value, hide: true };
      //     }
      //     return { name: currentKey, value };
      //   }
      // });
      //  let rows = [];

      // for (const [key, value] of Object.entries(obj)) {
      //   const currentKey = parentKey ? `${parentKey}.${key}` : key;

      //   if (Array.isArray(value)) {
      //     // Skip processing arrays
      //     continue;
      //   }

      //   let row = [];

      // if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
      //   row = [{ [currentKey]: processObjForKeyValTable(value, currentKey) }];
      // } else {
      //   if (key === 'attribute') {
      //     row = processObjForKeyValTable(JSON.parse(value), currentKey);
      //   } else if (key === 'id') {
      //     row.push({ name: currentKey, value: value, hide: true });
      //   } else {
      //     row.push({ name: currentKey, value });
      //   }
      // }

      // rows.push(...row); // Push the row into the main rows array
      // }

      let rows = [];
      let currentGroup = [];

      for (const [key, value] of Object.entries(obj)) {
        const currentKey = parentKey ? `${parentKey}.${key}` : key;

        if (Array.isArray(value)) {
          // Skip the key if the value is an array
          continue;
        } else if (typeof value === 'object' && value !== null) {
          // For objects, recursively process and add to the current group
          currentGroup.push(...processObjForKeyValTable(value, currentKey));
        } else {
          // For non-objects, add to the rows directly
          if (key === 'attribute') {
            currentGroup.push(...processObjForKeyValTable(JSON.parse(value), currentKey));
          } else if (key === 'id') {
            currentGroup.push({ name: currentKey, value: value, hide: true });
          } else {
            currentGroup.push({ name: currentKey, value });
          }
        }
      }

      // Group by the parent key
      if (parentKey !== '' && currentGroup.length > 0) {
        rows.push({ [parentKey]: currentGroup });
      } else {
        rows.push(...currentGroup);
      }

      return rows;
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

    let cleanedObj = processObjForKeyValTable(obj)[0]?.obj;

    const processObjectForTable = (obj, parentKey = '') => {
      console.log('inside processObjectForTable', obj, parentKey);
      for (const entry of obj) {
        const currentKey = entry.name;

        if (Array.isArray(entry.value)) {
          console.log('inside array', entry.value);
          // If the value is an array, call the Table component
          // with the parentKey and the array as data
          return <NameValueTable rows={entry.value} />;
        }

        if (typeof entry.value === 'object' && entry.value !== null) {
          // If the value is an object, recursively process it
          return processObjectForTable(entry.value, currentKey);
        }
      }
    };

    return (
      <>
        {processObjForKeyDataTable(obj)}
        {/* <NameValueTable rows={processObjForKeyValTable(obj)[0]?.obj} /> */}
        {cleanedObj.map((entry) => {
          return processObjectForTable([entry]);
        })}
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
