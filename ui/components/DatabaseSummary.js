import React, { useState } from 'react';
import { updateProgress } from '../lib/store';
import { Button, Typography, withStyles } from '@material-ui/core';
import { Provider, connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import resetDatabase from './graphql/queries/ResetDatabaseQuery';
import debounce from '../utils/debounce';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import ResponsiveDataTable from '../utils/data-table';
import SearchBar from '../utils/custom-search';
import useStyles from '../assets/styles/general/tool.styles';
import { PROMPT_VARIANTS } from './PromptComponent';
import { store } from '../store';
import { useGetDatabaseSummaryQuery } from '@/rtk-query/system';

const styles = (theme) => ({
  textCenter: {
    textAlign: 'center',
  },
  textEnd: {
    textAlign: 'end',
  },
  gapBottom: {
    paddingBottom: '0.5rem',
  },
  DBBtn: {
    margin: theme.spacing(0.5),
    padding: theme.spacing(1),
    borderRadius: 5,
    backgroundColor: '#8F1F00',
    '&:hover': {
      backgroundColor: '#B32700',
    },
  },
  container: {
    display: 'flex',
    justifyContent: 'center',
    marginTop: theme.spacing(2),
  },
});

const DatabaseSummary = (props) => {
  const { classes } = props;
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchText, setSearchText] = useState('');
  const { notify } = useNotification();
  const StyleClass = useStyles();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);

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
    {
      name: 'name',
      label: 'Name',
    },
    {
      name: 'count',
      label: 'Count',
    },
  ];

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility] = useState(() => {
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    columns.forEach((col) => {
      initialVisibility[col.name] = col.options?.display !== false;
    });
    return initialVisibility;
  });

  const customInlineStyle = {
    marginBottom: '0.5rem',
    marginTop: '1rem',
  };

  return (
    <>
      <div className={StyleClass.toolWrapper} style={customInlineStyle}>
        <div>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            size="medium"
            onClick={handleResetDatabase()}
            className={classes.DBBtn}
            data-cy="btnResetDatabase"
          >
            <Typography align="center" variant="subtitle2">
              {' '}
              RESET DATABASE{' '}
            </Typography>
          </Button>
        </div>
        <div>
          <SearchBar
            onSearch={(value) => {
              setSearchText(value);
            }}
            expanded={isSearchExpanded}
            setExpanded={setIsSearchExpanded}
            placeholder="Search"
          />
        </div>
      </div>
      <ResponsiveDataTable
        data={databaseSummary?.tables}
        options={table_options}
        columns={columns}
        tableCols={tableCols}
        updateCols={updateCols}
        columnVisibility={columnVisibility}
      />
    </>
  );
};

const mapStateToProps = () => ({});

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

DatabaseSummary.propTypes = {
  promptRef: PropTypes.object.isRequired,
};

const DatabaseSummaryTable = (props) => {
  return (
    <Provider store={store}>
      <DatabaseSummary {...props} />
    </Provider>
  );
};

export default withStyles(styles, { withTheme: true })(
  connect(mapStateToProps, mapDispatchToProps)(DatabaseSummaryTable),
);
