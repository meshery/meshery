import React from 'react';
import { ArrowBack } from '@material-ui/icons';
import TooltipButton from '../../utils/TooltipButton';
import { Paper, IconButton } from '@material-ui/core';
import { useRouter } from 'next/router';
// import { FormatStructuredData } from '../DataFormatter';
import NameValueTable from '../DataFormatter/NameValueTable';
// import ResponsiveDataTable from '../../utils/data-table';

const View = (props) => {
  const { type, setView, resource } = props;
  console.log(type, setView, resource);
  console.log('inside view component');

  // function renerDynamicTable(key, value) {
  //   console.log('inside render dynamic table');
  //   const allKeys = value.reduce((keys, obj) => {
  //     Object.keys(obj).forEach((key) => {
  //       if (!keys.includes(key)) {
  //         keys.push(key);
  //       }
  //     });
  //     return keys;
  //   }, []);

  //   const columns = allKeys.map((key) => ({
  //     name: key,
  //     label: key,
  //     options: {
  //       filter: false,
  //       sort: false,
  //     },
  //   }));

  //   let options = {
  //     filter: false,
  //     download: false,
  //     print: false,
  //     search: false,
  //     viewColumns: false,
  //     selectableRows: 'none',
  //     pagination: false,
  //     responsive: 'standard',
  //   };

  //   console.log("columns", columns)
  //   console.log("value", value)

  //   return (
  //     <>
  //       {key}

  //       <ResponsiveDataTable
  //         data={value}
  //         columns={columns}
  //         options={options}
  //         tableCols={columns}
  //         updateCols={() => {}}
  //       />
  //     </>
  //   );
  // }

  // function sanitizeResourceCopy(obj, parentKey = '') {
  //   Object.entries(obj).forEach(([key, value]) => {
  //     const currentKey = parentKey ? `${parentKey}.${key}` : key;

  //       //       console.log("val", value)
  //       // console.log("Array.isArray(value)", Array.isArray(value))

  //     if (Array.isArray(value)) {
  //      return renerDynamicTable(key, value);
  //     }
  //       if (typeof value === 'object' && value !== null) {
  //         return sanitizeResource(value, currentKey);
  //       } else {
  //         if (key === 'attribute') {
  //           return sanitizeResource(JSON.parse(value), currentKey);
  //         }
  //         if (key === 'id') {
  //           return { name: currentKey, value: value, hide: true };
  //         }
  //         return { name: currentKey, value };
  //       }

  //   });
  // }

  function sanitizeResource(obj, parentKey = '') {
    return Object.entries(obj).flatMap(([key, value]) => {
      const currentKey = parentKey ? `${parentKey}.${key}` : key;

      //       console.log("val", value)
      // console.log("Array.isArray(value)", Array.isArray(value))

      if (!Array.isArray(value)) {
        if (typeof value === 'object' && value !== null) {
          return sanitizeResource(value, currentKey);
        } else {
          if (key === 'attribute') {
            return sanitizeResource(JSON.parse(value), currentKey);
          }
          if (key === 'id') {
            return { name: currentKey, value: value, hide: true };
          }
          return { name: currentKey, value };
        }
      }
    });
  }

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

  const ResourceDetails = (resource) => {
    console.log(resource);
    return (
      <>
        <div
          style={{
            margin: '1rem auto',
          }}
        >
          {/* <FormatStructuredData data={resource} /> */}
          <NameValueTable rows={resource.resource} />
        </div>
      </>
    );
  };

  let cleanedResourceMetadata = sanitizeResource(resource?.metadata);
  let cleanedResourceStatus = sanitizeResource(resource?.status);
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
            <ResourceDetails resource={cleanedResourceMetadata} />
            <ResourceDetails resource={cleanedResourceStatus} />
          </div>
        </Paper>
      </div>
    </>
  );
};

export default View;
