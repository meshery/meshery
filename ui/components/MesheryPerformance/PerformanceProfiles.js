//@ts-check
import { ToggleButton, ToggleButtonGroup } from "@material-ui/lab";
import React, { useEffect, useState, useRef } from "react";
import PromptComponent from "../PromptComponent";
import GridOnIcon from "@material-ui/icons/GridOn";
import CloseIcon from "@material-ui/icons/Close";
import TableChartIcon from "@material-ui/icons/TableChart";
import PerformanceProfileTable from "./PerformanceProfileTable";
import PerformanceProfileGrid from "./PerformanceProfileGrid";
import dataFetch from "../../lib/data-fetch";
import IconButton from "@material-ui/core/IconButton";
import AddIcon from "@material-ui/icons/Add";
import { bindActionCreators } from "redux";
import { connect } from "react-redux";
import { updateProgress } from "../../lib/store";
import { withSnackbar } from "notistack";
import GenericModal from "../GenericModal";
import MesheryPerformanceComponent from "./index";
import { Paper, Typography, Button } from "@material-ui/core";
import fetchPerformanceProfiles from "../graphql/queries/PerformanceProfilesQuery";

const MESHERY_PERFORMANCE_URL = "/api/user/performance/profiles";

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
      onChange={(_, newView) => changeView(newView)}
      aria-label="Switch View"
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
  const modalRef = useRef(null)

  const [page, setPage] = useState(0);
  const [search, setSearch] = useState("");
  const [sortOrder, setSortOrder] = useState("");
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [testProfiles, setTestProfiles] = useState([]);
  const [profileForModal, setProfileForModal] = useState();
  // const [loading, setLoading] = useState(false);

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
    // const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
    //   sortOrder
    // )}`;

    updateProgress({ showProgress : true });
    fetchPerformanceProfiles({
      selector : {
        pageSize : `${pageSize}`,
        page : `${page}`,
        search : `${encodeURIComponent(search)}`,
        order : `${encodeURIComponent(sortOrder)}`,
      },
    }).subscribe({
      next : (res) => {
        // @ts-ignore
        let result = res?.getPerformanceProfiles;
        updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          if (result) {
            setCount(result.total_count || 0);
            setPageSize(result.page_size || 0);
            setTestProfiles(result.profiles || []);
            setPage(result.page || 0);
          }
        }
      },
      error : handleError("Failed to Fetch Profiles"),
    });
  }

  async function showModal() {
    let response = await modalRef.current.show({
      title : "Delete Performance Profile?",

      subtitle : "Are you sure you want to delete this performance profile?",

      options : ["Yes", "No"],
    })
    return response;
  }

  function deleteProfile(id) {
    dataFetch(
      `${MESHERY_PERFORMANCE_URL}/${id}`,
      {
        method : "DELETE",
        credentials : "include",
      },
      () => {
        updateProgress({ showProgress : false });

        enqueueSnackbar("Performance Profile Successfully Deleted!", {
          variant : "success",
          autoHideDuration : 2000,
          action : function Action(key) {
            return (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            );
          },
        });

        fetchTestProfiles(page, pageSize, search, sortOrder);
      },
      handleError("Failed To Delete Profile")
    );
  }


  function handleError(msg) {
    return function (error) {
      updateProgress({ showProgress : false });

      enqueueSnackbar(`${msg} : ${error}`, {
        variant : "error",
        action : function Action(key) {
          return (
            <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon />
            </IconButton>
          );
        },
        autoHideDuration : 8000,
      });
    };
  }

  return (
    <>
      <div style={{ padding : "0.5rem" }}>
        <div
          style={{
            margin : "2rem auto",
            display : "flex",
            justifyContent : "space-between",
            paddingLeft : "1rem"
          }}
        >
          {testProfiles.length > 0 && (
            <div style={{ width : "fit-content", alignSelf : "flex-start" }}>
              <Button
                aria-label="Add Performance Profile"
                variant="contained"
                color="primary"
                size="large"
                // @ts-ignore
                onClick={() => setProfileForModal({})}
              >
                <AddIcon />
                Add Performance Profile
              </Button>
            </div>
          )}
          <div style={{ justifySelf : "flex-end", marginLeft : "auto", paddingLeft : "1rem" }}>
            <ViewSwitch view={viewType} changeView={setViewType} />
          </div>
        </div>
        {viewType === "grid"
          ? (
            <PerformanceProfileGrid
              profiles={testProfiles}
              deleteHandler={deleteProfile}
              setProfileForModal={setProfileForModal}
              pages={Math.ceil(count / pageSize)}
              setPage={setPage}
            />
          )
          : (
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
              setProfileForModal={setProfileForModal}
              handleDelete={deleteProfile}
              showModal={showModal}
              fetchTestProfiles={fetchTestProfiles}
            />
          )}
        {testProfiles.length == 0 && (
          <Paper style={{ padding : "0.5rem" }}>
            <div
              style={{
                padding : "2rem",
                display : "flex",
                justifyContent : "center",
                alignItems : "center",
                flexDirection : "column",
              }}
            >
              <Typography style={{ fontSize : "1.5rem", marginBottom : "2rem" }} align="center" color="textSecondary">
                No Performance Profiles Found
              </Typography>
              <Button
                aria-label="Add Performance Profile"
                variant="contained"
                color="primary"
                size="large"
                // @ts-ignore
                onClick={() => setProfileForModal({})}
              >
                <AddIcon />
                Add Performance Profile
              </Button>
            </div>
          </Paper>
        )}
        <GenericModal
          open={!!profileForModal}
          Content={
            <Paper style={{ margin : "auto", maxWidth : "90%", outline : "none" }}>
              <MesheryPerformanceComponent
                // @ts-ignore
                loadAsPerformanceProfile
                // @ts-ignore
                performanceProfileID={profileForModal?.id}
                // @ts-ignore
                profileName={profileForModal?.name}
                // @ts-ignore
                meshName={profileForModal?.service_mesh}
                // @ts-ignore
                url={profileForModal?.endpoints?.[0]}
                // @ts-ignore
                qps={profileForModal?.qps}
                // @ts-ignore
                loadGenerator={profileForModal?.load_generators?.[0]}
                // @ts-ignore
                t={profileForModal?.duration}
                // @ts-ignore
                c={profileForModal?.concurrent_request}
                // @ts-ignore
                reqBody={profileForModal?.request_body}
                // @ts-ignore
                headers={profileForModal?.request_headers}
                // @ts-ignore
                cookies={profileForModal?.request_cookies}
                // @ts-ignore
                contentType={profileForModal?.content_type}
                // @ts-ignore
                runTestOnMount={!!profileForModal?.runTest}
              />
            </Paper>
          }
          handleClose={() => {
            fetchTestProfiles(page, pageSize, search, sortOrder);
            setProfileForModal(undefined);
          }}
        />
      </div>

      <PromptComponent ref={modalRef} />
    </>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

export default connect(null, mapDispatchToProps)(withSnackbar(PerformanceProfile));
