//@ts-check
import { ToggleButton, ToggleButtonGroup } from "@material-ui/lab";
import React, { useEffect, useState } from "react";
import GridOnIcon from "@material-ui/icons/GridOn";
import CloseIcon from "@material-ui/icons/Close";
import TableChartIcon from "@material-ui/icons/TableChart";
import PerformanceProfileTable from "./PerformanceProfileTable";
import PerformanceProfileGrid from "./PerformanceProfileGrid";
import dataFetch from "../../lib/data-fetch";
import IconButton from "@material-ui/core/IconButton";
import { bindActionCreators } from "redux";
import { connect } from "react-redux";
import { updateProgress } from "../../lib/store";
import { withSnackbar } from "notistack";

/**
 * Type Definition for View Type
 * @typedef {"grid" | "table"} TypeView
 */

/**
 * ViewSwitch component renders a switch for toggling between
 * grid and table views
 * @param {{ view: TypeView, changeView: (view: TypeView) => void }} props
 */
function ViewSwitch({ view, changeView }) {
  return (
    <ToggleButtonGroup
      size="small"
      value={view}
      exclusive
      onChange={(_, newView) => changeView(newView)} aria-label="Switch View"
    >
      <ToggleButton value="grid">
        <GridOnIcon />
      </ToggleButton>
      <ToggleButton value="table">
        <TableChartIcon />
      </ToggleButton>
    </ToggleButtonGroup>
  );
}

function PerformanceProfile({ updateProgress, enqueueSnackbar, closeSnackbar }) {
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ("grid")
  );

  const [page, setPage] = useState(0);
  const [search, setSearch] = useState("");
  const [sortOrder, setSortOrder] = useState("");
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [testProfiles, setTestProfiles] = useState([]);

  /**
   * fetch performance profiles when the page loads
   */
  useEffect(() => {
    fetchTestProfiles(page, pageSize, search, sortOrder);
  }, [page, pageSize, search, sortOrder]);

  /**
   * fetchTestProfiles constructs the queries based on the parameters given
   * and fetches the performance profiles
   * @param {number} page current page
   * @param {number} pageSize items per page
   * @param {string} search search string
   * @param {string} sortOrder order of sort
   */
  function fetchTestProfiles(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
      sortOrder
    )}`;

    const val = "f";
    if (val) {
      updateProgress({ showProgress: true });

      dataFetch(
        `/api/user/performance/profiles${query}`,
        {
          credentials: "include",
        },
        (result) => {
          updateProgress({ showProgress: false });
          if (result) {
            setTestProfiles(result.profiles || []);
            setPage(result.page || 0);
            setPageSize(result.page_size || 0);
            setCount(result.total_count || 0);
          }
        },
        handleError
      );
    }
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

  return (
    <div style={{ padding: "0.5rem" }}>
      <div style={{ margin: "0 0 2rem auto", width: "fit-content" }}>
        <ViewSwitch view={viewType} changeView={setViewType} />
      </div>
      {viewType === "grid" ? (
        <PerformanceProfileGrid profiles={testProfiles} />
      ) : (
        <PerformanceProfileTable
          page={page}
          setPage={setPage}
          search={search}
          setSearch={setSearch}
          sortOrder={sortOrder}
          setSortOrder={setSortOrder}
          count={count}
          pageSize={pageSize}
          setPageSize={setPageSize}
          testProfiles={testProfiles}
        />
      )}
    </div>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default connect(null, mapDispatchToProps)(withSnackbar(PerformanceProfile));
