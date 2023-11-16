import React from 'react';
import { ArrowBack } from '@material-ui/icons';
import TooltipButton from '../../utils/TooltipButton';
import { Paper, IconButton } from '@material-ui/core';
import { useRouter } from 'next/router';
// import { FormatStructuredData } from '../DataFormatter';
import NameValueTable from '../DataFormatter/NameValueTable';
import ResponsiveDataTable from '../../utils/data-table';

const View = (props) => {
  const { type, setView, resource } = props;
  console.log(type, setView, resource);
  console.log('inside view component');

  // TODO: handle the condition for images
  function RenderDynamicTable(key, value) {
    console.log('inside render dynamic table');
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

    console.log('columns', columns);
    console.log('value', value);

    return (
      <>
        {key}

        <ResponsiveDataTable
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
    const processObjForKeyValTable = (obj, parentKey = '') => {
      let rows = [];
      for (const [key, value] of Object.entries(obj)) {
        if (Array.isArray(value)) {
          continue;
        }

        const currentKey = parentKey ? `${parentKey}.${key}` : key;
        if (typeof value === 'object' && value !== null) {
          return processObjForKeyValTable(value, currentKey);
        } else {
          if (key === 'attribute') {
            return processObjForKeyValTable(JSON.parse(value), currentKey);
          }
          if (key === 'id') {
            return rows.push({ name: currentKey, value: value, hide: true });
          }
          return rows.push({ name: currentKey, value });
        }
      }

      console.log('rows', rows);

      return (
        <>
          <NameValueTable rows={rows} />
        </>
      );
    };

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
        {processObjForKeyValTable(obj)}
      </>
    );
  };

  //   const RecursiveTable = ({ data }) => {
  //     const [accumulatedData, setAccumulatedData] = useState([]);
  //     let currentArray = [];

  //     const renderTable = () => {
  //       if (accumulatedData.length > 0) {
  //         const result = <NameValueTable rows={accumulatedData} />;
  //         setAccumulatedData([]);
  //         return result;
  //       }
  //       return null;
  //     };

  //     const renderRecursiveTable = (obj, parentKey = '') => {
  //     const rows = [];

  //     for (const [key, value] of Object.entries(obj)) {
  //       const currentKey = parentKey ? `${parentKey}.${key}` : key;

  //       if (Array.isArray(value)) {
  //         // If value is an array, render a subtable
  //         rows.push(<RenderDynamicTable key={currentKey} value={value} />);
  //       } else if (typeof value === 'object' && value !== null) {
  //         // If value is an object, recursively render a subtable
  //         rows.push(renderRecursiveTable(value, currentKey));
  //       } else if (key === 'attribute') {
  //           return rows.push(renderRecursiveTable(JSON.parse(value), currentKey));
  //         }
  //         else if (key === 'id') {
  //           return rows.push({ name: currentKey, value: value, hide: true });
  //         } else {
  //         // For non-object values, accumulate data for rendering a NameValueTable
  //         rows.push({ name: currentKey, value });
  //       }
  //     }

  //     if (rows.length > 0) {
  //       // If there are accumulated rows, render a NameValueTable
  //       return <NameValueTable rows={rows} />;
  //     }

  //     return null;
  //     };

  //      return <>{renderRecursiveTable(data)}</>;
  // };

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

  // const ResourceDetails = (resource) => {
  //   console.log(resource);
  //   return (
  //     <>
  //       <div
  //         style={{
  //           margin: '1rem auto',
  //         }}
  //       >
  //         {/* <FormatStructuredData data={resource} /> */}
  //         <NameValueTable rows={resource.resource} />
  //       </div>
  //     </>
  //   );
  // };

  // let cleanedResourceMetadata = sanitizeResource(resource?.metadata);
  // let cleanedResourceStatus = sanitizeResource(resource?.status);
  // let cleanedResourceSpec = sanitizeResource(resource?.spec);

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
            {/* <ResourceDetails resource={cleanedResourceMetadata} />
            <ResourceDetails resource={cleanedResourceStatus} /> */}
            {/* <SanitizeResourceCopy obj={resource?.status} /> */}
            {/* <RecursiveTable data={resource?.metadata} />
            <RecursiveTable data={resource?.status} /> */}
            <RenderObject obj={resource?.metadata} />
            <RenderObject obj={resource?.status} />
          </div>
        </Paper>
      </div>
    </>
  );
};

export default View;
