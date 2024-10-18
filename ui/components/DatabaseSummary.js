import React, { useState } from 'react';
import { updateProgress } from '../lib/store';
import { Button, Typography, useTheme } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';
import { Provider, connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import resetDatabase from './graphql/queries/ResetDatabaseQuery';
import debounce from '../utils/debounce';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import ResponsiveDataTable from '../utils/data-table';
import SearchBar from '../utils/custom-search';
import { PROMPT_VARIANTS } from './PromptComponent';
import { store } from '../store';
import { useGetDatabaseSummaryQuery } from '@/rtk-query/system';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

const DatabaseSummary = (props) => {
  const theme = useTheme();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchText, setSearchText] = useState('');
  const { notify } = useNotification();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [tableCols, updateCols] = useState(columns);
  const [columnVisibility] = useState(() => {
    const initialVisibility = {};
    columns.forEach((col) => {
      initialVisibility[col.name] = col.options?.display !== false;
    });
    return initialVisibility;
  });

  const handleError = (msg) => (error) => {
    props.updateProgress({ showProgress: false });
    notify({
      message: `${msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const { data: databaseSummary, refetch } = useGetDatabaseSummaryQuery({
    page: page + 1,
    pagesize: rowsPerPage,
    search: searchText,
  });

  const handleResetDatabase = () => {
    return async () => {
      let responseOfResetDatabase = await props.promptRef.current.show({
        title: 'Reset Meshery Database?',
        subtitle: 'Are you sure that you want to purge all data?',
        options: ['RESET', 'CANCEL'],
        variant: PROMPT_VARIANTS.DANGER,
      });
      if (responseOfResetDatabase === 'RESET') {
        props.updateProgress({ showProgress: true });
        resetDatabase({
          selector: {
            clearDB: 'true',
            ReSync: 'true',
            hardReset: 'true',
          },
          k8scontextID: '',
        }).subscribe({
          next: (res) => {
            props.updateProgress({ showProgress: false });
            if (res.resetStatus === 'PROCESSING') {
              notify({ message: 'Database reset successful.', event_type: EVENT_TYPES.SUCCESS });
              refetch();
            }
          },
          error: handleError('Database is not reachable, try restarting server.'),
        });
      }
    };
  };

  const table_options = {
    filter: false,
    sort: false,
    selectableRows: 'none',
    responsive: 'scrollMaxHeight',
    print: false,
    download: false,
    viewColumns: false,
    search: false,
    fixedHeader: true,
    serverSide: true,
    rowsPerPage: rowsPerPage,
    count: databaseSummary?.total_tables,
    page: page,
    onChangePage: debounce((p) => setPage(p), 200),
    onChangeRowsPerPage: debounce((p) => setRowsPerPage(p), 200),
    onSearchChange: debounce((searchText) => {
      if (searchText) setPage(0);
      setSearchText(searchText != null ? searchText : '');
    }),
  };

  const columns = [
    { name: 'name', label: 'Name' },
    { name: 'count', label: 'Count' },
  ];

  return (
    <UsesSistent>
      <div style={{ marginBottom: theme.spacing(1), marginTop: theme.spacing(2) }}>
        <Button
          variant="contained"
          color="error"
          disabled={!CAN(keys.RESET_DATABASE.action, keys.RESET_DATABASE.subject)}
          onClick={handleResetDatabase()}
          data-cy="btnResetDatabase"
        >
          <Typography variant="subtitle2">RESET DATABASE</Typography>
        </Button>
        <SearchBar
          onSearch={(value) => {
            setSearchText(value);
          }}
          expanded={isSearchExpanded}
          setExpanded={setIsSearchExpanded}
          placeholder="Search"
        />
      </div>
      <ResponsiveDataTable
        data={databaseSummary?.tables}
        options={table_options}
        columns={columns}
        tableCols={tableCols}
        updateCols={updateCols}
        columnVisibility={columnVisibility}
      />
    </UsesSistent>
  );
};

DatabaseSummary.propTypes = {
  promptRef: PropTypes.object.isRequired,
  updateProgress: PropTypes.func.isRequired,
};

const DatabaseSummaryTable = (props) => (
  <Provider store={store}>
    <DatabaseSummary {...props} />
  </Provider>
);

export default connect(null, (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
}))(DatabaseSummaryTable);
