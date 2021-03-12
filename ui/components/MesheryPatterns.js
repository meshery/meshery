import React, { useState, useEffect, useRef } from "react";
import { withStyles } from "@material-ui/core/styles";
import {
  NoSsr,
  TableCell,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Divider,
  Tooltip
} from "@material-ui/core";
import { UnControlled as CodeMirror } from "react-codemirror2";
import DeleteIcon from "@material-ui/icons/Delete";
import UpdateIcon from "@material-ui/icons/Update";
import UploadIcon from "@material-ui/icons/Publish";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import MoreHorizIcon from "@material-ui/icons/MoreHoriz";
import { updateProgress } from "../lib/store";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import dataFetch from "../lib/data-fetch";

const styles = (theme) => ({
  grid: {
    padding: theme.spacing(2),
  },
  tableHeader: {
    fontWeight: "bolder",
    fontSize: 18,
  },
});

function CustomToolbar(onClick) {
  return function Toolbar() {
    return (
      <>
        <label htmlFor="upload-button">
          <input type="file" accept=".yaml, .yml" hidden onChange={onClick} id="upload-button" name="upload-button" />
          <Tooltip title="Upload Pattern">
            <IconButton aria-label="Upload" component="span">
              <UploadIcon />
            </IconButton>
          </Tooltip>
        </label>
      </>
    );
  };
}

function YAMLEditor({ pattern, onClose, onSubmit }) {
  const [yaml, setYaml] = useState("");

  return (
    <Dialog onClose={onClose} aria-labelledby="pattern-dialog-title" open fullWidth maxWidth="md">
      <DialogTitle id="pattern-dialog-title">{pattern.name}</DialogTitle>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <CodeMirror
          value={pattern.pattern_file}
          options={{
            theme: "material",
            lineNumbers: true,
            lineWrapping: true,
            gutters: ["CodeMirror-lint-markers"],
            lint: true,
            mode: "text/x-yaml",
          }}
          onChange={(_, data, val) => setYaml(val)}
        />
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        <Tooltip title="Update Pattern">
          <IconButton
            aria-label="Update"
            color="primary"
            onClick={() => onSubmit(yaml, pattern.id, pattern.name, "update")}
          >
            <UpdateIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete Pattern">
          <IconButton
            aria-label="Delete"
            color="primary"
            onClick={() => onSubmit(yaml, pattern.id, pattern.name, "delete")}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </DialogActions>
    </Dialog>
  );
}

function MesheryPatterns({ updateProgress, enqueueSnackbar, closeSnackbar, user, classes }) {
  const [page, setPage] = useState(0);
  const [search] = useState("");
  const [sortOrder] = useState("");
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [patterns, setPatterns] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);

  const searchTimeout = useRef(null);

  /**
   * fetch patterns when the page loads
   */
  useEffect(() => {
    fetchPatterns(page, pageSize, search, sortOrder);
  }, []);

  /**
   * fetchPatterns constructs the queries based on the parameters given
   * and fetches the patterns
   * @param {number} page current page
   * @param {number} pageSize items per page
   * @param {string} search search string
   * @param {string} sortOrder order of sort
   */
  function fetchPatterns(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
      sortOrder
    )}`;

    updateProgress({ showProgress: true });

    dataFetch(
      `/api/experimental/patternfile${query}`,
      {
        credentials: "include",
      },
      (result) => {
        console.log("PatternFile API", `/api/experimental/patternfile${query}`);
        updateProgress({ showProgress: false });
        if (result) {
          setPatterns(result.patterns || []);
          setPage(result.page || 0);
          setPageSize(result.page_size || 0);
          setCount(result.total_count || 0);
        }
      },
      handleError
    );
  }

  function handleError(error) {
    updateProgress({ showProgress: false });

    enqueueSnackbar(`There was an error fetching results: ${error}`, {
      variant: "error",
      action: function Action(key) {
        return (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        );
      },
      autoHideDuration: 8000,
    });
  }

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  function handleSubmit(data, id, name, type) {
    if (type === "delete") {
      dataFetch(
        `/api/experimental/patternfile/${id}`,
        {
          credentials: "include",
          method: "DELETE",
        },
        () => {
          console.log("PatternFile API", `/api/experimental/patternfile/${id}`);
          updateProgress({ showProgress: false });
          fetchPatterns(page, pageSize, search, sortOrder);
        },
        handleError
      );
    }

    if (type === "update") {
      dataFetch(
        `/api/experimental/patternfile`,
        {
          credentials: "include",
          method: "POST",
          body: JSON.stringify({ id, name, pattern_file: data }),
        },
        () => {
          console.log("PatternFile API", `/api/experimental/patternfile`);
          updateProgress({ showProgress: false });
          fetchPatterns(page, pageSize, search, sortOrder);
        },
        handleError
      );
    }

    if (type === "upload") {
      dataFetch(
        `/api/experimental/patternfile`,
        {
          credentials: "include",
          method: "POST",
          body: JSON.stringify({ name, pattern_file: data }),
        },
        () => {
          console.log("PatternFile API", `/api/experimental/patternfile`);
          updateProgress({ showProgress: false });
          fetchPatterns(page, pageSize, search, sortOrder);
        },
        handleError
      );
    }
  }

  function uploadHandler(ev) {
    if (!ev.target.files?.length) return;

    const file = ev.target.files[0];

    // Create a reader
    const reader = new FileReader();
    reader.addEventListener("load", (event) => {
      handleSubmit(
        event.target.result, 
        "", 
        file?.name || "meshery_" + Math.floor(Math.random() * 100), 
        "upload",
      );
    });
    reader.readAsText(file);
  }

  const columns = [
    {
      name: "name",
      label: "Pattern Name",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: "created_at",
      label: "Upload Timestamp",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: "updated_at",
      label: "Update Timestamp",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: "Details",
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(_, tableMeta) {
          return (
            <IconButton
              aria-label="more"
              color="inherit"
              onClick={() => setSelectedRowData(patterns[tableMeta.rowIndex])}
            >
              <MoreHorizIcon />
            </IconButton>
          );
        },
      },
    },
  ];

  columns.forEach((column, idx) => {
    if (column.name === sortOrder.split(" ")[0]) {
      columns[idx].options.sortDirection = sortOrder.split(" ")[1];
    }
  });

  const options = {
    filter: false,
    sort: !(user && user.user_id === "meshery"),
    search: !(user && user.user_id === "meshery"),
    filterType: "textField",
    responsive: "scrollFullHeight",
    resizableColumns: true,
    serverSide: true,
    selectableRows: "none",
    count,
    rowsPerPage: pageSize,
    rowsPerPageOptions: [10, 20, 25],
    fixedHeader: true,
    page,
    print: false,
    download: false,
    customToolbar: CustomToolbar(uploadHandler),

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(" : ") : [];
      let order = "";
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case "changePage":
          fetchPatterns(tableState.page, pageSize, search, sortOrder);
          break;
        case "changeRowsPerPage":
          fetchPatterns(page, tableState.rowsPerPage, search, sortOrder);
          break;
        case "search":
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              fetchPatterns(page, pageSize, tableState.searchText !== null ? tableState.searchText : "", sortOrder);
            }
          }, 500);
          break;
        case "sort":
          if (sortInfo.length == 2) {
            if (sortInfo[1] === "ascending") {
              order = `${columns[tableState.activeColumn].name} asc`;
            } else {
              order = `${columns[tableState.activeColumn].name} desc`;
            }
          }
          if (order !== sortOrder) {
            fetchPatterns(page, pageSize, search, order);
          }
          break;
      }
    },
  };

  return (
    <NoSsr>
      {selectedRowData && Object.keys(selectedRowData).length > 0 && (
        <YAMLEditor pattern={selectedRowData} onClose={resetSelectedRowData()} onSubmit={handleSubmit} />
      )}
      <MUIDataTable
        title={<div className={classes.tableHeader}>Meshery Patterns</div>}
        data={patterns}
        columns={columns}
        // @ts-ignore
        options={options}
      />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  return {
    user: state.get("user").toObject(),
  };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryPatterns)));
