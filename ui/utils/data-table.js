// Data table for parent tables with custom column visibility control
import MUIDataTable from 'mui-datatables';
import React, { useEffect, useCallback } from 'react';
import { useWindowDimensions } from './dimension';
import { styled } from '@mui/system';

const StyledMUIDataTable = styled(MUIDataTable)(() => ({
  padding: '5px 18px 0 18px',
}));

const ResponsiveDataTable = ({
  data,
  columns,
  options,
  tableCols,
  updateCols,
  columnVisibility,
  ...props
}) => {
  const { width } = useWindowDimensions();

  const formatDate = useCallback(
    (date, width) => {
      const dateOptions = {
        day: 'numeric',
        weekday: 'short',
        month: 'long',
        year: 'numeric',
      };

      if (width < 1240 && width >= 915) {
        dateOptions.month = 'short';
        dateOptions.day = 'numeric';
        dateOptions.year = 'numeric';
        dateOptions.weekday = 'short';
      } else if (width < 915) {
        dateOptions.month = 'short';
        dateOptions.day = '2-digit';
        dateOptions.year = 'numeric';
        dateOptions.weekday = undefined;
      }

      return date.toLocaleDateString('en-US', dateOptions);
    },
    [width],
  );

  const updatedOptions = {
    ...options,
    // viewColumns: false,
    onViewColumnsChange: (column, action) => {
      let colToChange;
      switch (action) {
        case 'add':
          colToChange = columns.find((obj) => obj.name === column);
          colToChange.options.display = true;
          updateCols([...columns]);
          break;
        case 'remove':
          colToChange = columns.find((obj) => obj.name === column);
          colToChange.options.display = false;
          updateCols([...columns]);
          break;
      }
    },
  };

  useEffect(() => {
    columns?.forEach((col) => {
      if (!col.options) {
        col.options = {};
      }
      // Set the display option based on columnVisibility state
      col.options.display = columnVisibility[col.name];

      if (
        ['updated_at', 'created_at', 'deleted_at', 'last_login_time', 'joined_at'].includes(
          col.name,
        )
      ) {
        col.options.customBodyRender = (value) => {
          if (value === 'NA') {
            return value;
          } else if (value?.Valid === true) {
            const date = new Date(value.Time);
            return formatDate(date, width);
          } else if (value?.Valid === false) {
            return 'NA';
          } else {
            const date = new Date(value);
            return formatDate(date, width);
          }
        };
      }
    });
    updateCols([...columns]);
  }, [width, columnVisibility]);

  const components = {
    ExpandButton: function () {
      return '';
    },
  };

  return (
    <StyledMUIDataTable
      components={components}
      columns={tableCols}
      data={data}
      options={updatedOptions}
      {...props}
    />
  );
};

export default ResponsiveDataTable;
