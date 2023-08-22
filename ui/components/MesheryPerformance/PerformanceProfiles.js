// @ts-nocheck
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
import AddIcon from "@material-ui/icons/AddCircleOutline";
import { bindActionCreators } from "redux";
import { connect } from "react-redux";
import { updateProgress } from "../../lib/store";
import { withSnackbar } from "notistack";
import GenericModal from "../GenericModal";
import MesheryPerformanceComponent from "./index";
import { Paper, Typography, Button,DialogTitle } from "@material-ui/core";
import fetchPerformanceProfiles from "../graphql/queries/PerformanceProfilesQuery";
import { withStyles } from "@material-ui/core/styles";
import { iconMedium, iconXLarge } from "../../css/icons.styles";
import subscribePerformanceProfiles from "../graphql/subscriptions/PerformanceProfilesSubscription";

const MESHERY_PERFORMANCE_URL = "/api/user/performance/profiles";
const styles = (theme) => ({
  title : {
    textAlign : 'center',
    minWidth : 400,
    padding : '10px',
    color : '#fff',
    flexGrow : 1
  },
  dialogHeader : {
    backgroundColor : theme.palette.type === 'dark' ? theme.palette.secondary.headerColor : theme.palette.secondary.mainBackground,
    display : "flex",
    alignItems : "center"
  },
  topToolbar : {
    margin : "2rem auto",
    display : "flex",
    justifyContent : "space-between",
    paddingLeft : "1rem"
  },
  addButton : {
    width : "fit-content",
    alignSelf : "flex-start"
  },
  viewSwitchButton : {
    justifySelf : "flex-end",
    marginLeft : "auto",
    paddingLeft : "1rem"
  },
  pageContainer : {
    padding : "0.5rem"
  },
  noProfileContainer : {
    padding : "2rem",
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    flexDirection : "column",
  },
  noProfilePaper : {
    padding : "0.5rem"
  },
  noProfileText : {
    fontSize : "1.5rem",
    marginBottom : "2rem",
  },
  addProfileModal : {
    margin : "auto",
    maxWidth : "90%",
    outline : "none"
  },
  addIcon : {
    paddingRight : "0.5"
  }
});
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
      <ToggleButton style={iconXLarge} value="grid">
        <GridOnIcon style={iconMedium}/>
      </ToggleButton>
      <ToggleButton style={iconXLarge} value="table">
        <TableChartIcon style={iconMedium}/>
      </ToggleButton>
    </ToggleButtonGroup>
  );
}

function PerformanceProfile({ updateProgress, enqueueSnackbar, closeSnackbar, classes }) {
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
    const subscription = subscribePerformanceProfiles((res) => {
      let result = res?.subscribePerfProfiles;
      if (typeof result !== "undefined") {
        if (result) {
          setCount(result.total_count || 0);
          setPageSize(result.page_size || 0);
          setTestProfiles(result.profiles || []);
          setPage(result.page || 0);
        }
      }
    }, {
      selector : {
        pageSize : `${pageSize}`,
        page : `${page}`,
        search : `${encodeURIComponent(search)}`,
        order : `${encodeURIComponent(sortOrder)}`,
      }
    })
    return () => {
      subscription.dispose();
    };
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
        let result = res?.getPerformanceProfiles;
        updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          if (result) {
            setCount(result.total_count || 0);
            setPageSize(result.page_size || 0);
            setTestProfiles(result.profiles || []);
            // setPage(result.page || 0);
          }
        }
      },
      error : handleError("Failed to Fetch Profiles"),
    });
  }

  async function showModal(count) {
    let response = await modalRef.current.show({
      title : `Delete ${count ? count : ""} Performance Profile${count > 1 ? "s" : ''}?`,
      subtitle : `Are you sure you want to delete ${count > 1 ? "these" : 'this'} ${count ? count : ""} performance profile${count > 1 ? "s" : ''}?`,

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

        enqueueSnackbar("Performance Profile Deleted!", {
          variant : "success",
          autoHideDuration : 2000,
          action : function Action(key) {
            return (
              <IconButton style={iconMedium} key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon style={iconMedium}/>
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
            <IconButton style={iconMedium} key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon style={iconMedium} />
            </IconButton>
          );
        },
        autoHideDuration : 8000,
      });
    };
  }

  return (
    <>
      <div className={classes.pageContainer}>
        <div className={classes.topToolbar}>
          {(testProfiles.length > 0 || viewType == "table") && (
            <div className={classes.addButton}>
              <Button
                aria-label="Add Performance Profile"
                variant="contained"
                color="primary"
                size="large"
                onClick={() => setProfileForModal({})}
              >
                <AddIcon style={iconMedium} className={classes.addIcon} />
                Add Performance Profile
              </Button>
            </div>
          )}
          <div className={classes.viewSwitchButton}>
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
        {testProfiles.length === 0 && viewType === "grid" && (
          <Paper className={classes.noProfilePaper} >
            <div className={classes.noProfileContainer}>
              <Typography className={classes.noProfileText} align="center" color="textSecondary">
                No Performance Profiles Found
              </Typography>
              <Button
                aria-label="Add Performance Profile"
                variant="contained"
                color="primary"
                size="large"
                onClick={() => setProfileForModal({})}
              >
                <Typography className="addIcon">Add Performance Profile</Typography>
              </Button>
            </div>
          </Paper>
        )}
        <GenericModal
          open={!!profileForModal}
          Content={
            <Paper className={classes.addProfileModal} >
              <div className={classes.dialogHeader}>
                <DialogTitle className={classes.title} >Performance Profile Wizard</DialogTitle>
                <IconButton aria-label="close" style={{ color : "white" }} onClick={() => setProfileForModal(undefined) }>
                  <CloseIcon />
                </IconButton>
              </div>

              <MesheryPerformanceComponent
                loadAsPerformanceProfile
                performanceProfileID={profileForModal?.id}
                profileName={profileForModal?.name}
                meshName={profileForModal?.service_mesh}
                url={profileForModal?.endpoints?.[0]}
                qps={profileForModal?.qps}
                loadGenerator={profileForModal?.load_generators?.[0]}
                t={profileForModal?.duration}
                c={profileForModal?.concurrent_request}
                reqBody={profileForModal?.request_body}
                headers={profileForModal?.request_headers}
                cookies={profileForModal?.request_cookies}
                contentType={profileForModal?.content_type}
                runTestOnMount={!!profileForModal?.runTest}
                metadata={profileForModal?.metadata}
              />
            </Paper>
          }
          handleClose={() => {
            setProfileForModal(undefined);
          }}
        />
      </div>

      <PromptComponent ref={modalRef} />
    </>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

export default withStyles(styles)(connect(null, mapDispatchToProps)(withSnackbar(PerformanceProfile)));
