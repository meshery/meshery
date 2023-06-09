//@ts-check
import React, { useEffect, useState } from "react";
import {
  Avatar,
  Button, Grid, IconButton, Link, Table, TableCell, TableRow, Typography
} from "@material-ui/core";
import DeleteIcon from "@material-ui/icons/Delete";
import EditIcon from "@material-ui/icons/Edit";
import Moment from "react-moment";
import PerformanceResults from "./PerformanceResults";
import FlipCard from "../FlipCard";
import { makeStyles } from "@material-ui/core/styles";
import { iconMedium } from "../../css/icons.styles";
import { useTheme } from '@material-ui/core/styles';
import moment from "moment";
import dataFetch from "../../lib/data-fetch";

const useStyles = makeStyles((theme) => ({
  cardButtons : {
    display : "flex",
    justifyContent : "flex-end",
    alignItems : "center"
  },
  testsButton : {
    marginRight : "0.5rem"
  },
  perfResultsContainer : {
    marginTop : "0.5rem"
  },
  backGrid : {
    marginBottom : "0.25rem",
    minHeight : "6rem"
  },
  deleteEditButton : {
    width : "fit-content",
    margin : "0 0 0 auto"
  },
  noOfResultsContainer : {
    margin : "0 0 1rem",
    '& div' : {
      display : "flex",
      alignItems : "center"
    },
  },
  bottomPart : {
    display : "flex",
    justifyContent : "space-between",
    alignItems : "center",
  },
  lastRunText : {
    marginRight : "0.5rem",
    marginLeft : "0.5rem",
  },
  resultText : {
    color : theme.palette.secondary.lightText,
  },
}))

const avatarHandler = (function AvatarHandler() {
const idToAvatarMap = {}

function fetchUserAvatarLink(userId, setterCallbackFunction) {
  if (!userId) return null
  dataFetch(`/api/user/profile/${userId}`, {
    credentials : "include"
  },
  function assignAvatarLinkToId(result) {
    if (result.avatar_url) {
      idToAvatarMap[userId] = result.avatar_url
      setterCallbackFunction(result.avatar_url)
    }
  },
  function handleProfileFetchError(error) {
    console.error("failed to fetch user profile with ID", userId, error)
  }
  )
}

return {
  getAvatar : async function _getAvatar(userId, setterFn) {
    if (idToAvatarMap[userId]) {
      return setterFn(idToAvatarMap[userId])
    }

    fetchUserAvatarLink(userId, setterFn)
  }
}
})()

function PerformanceCard({
  profile,
  handleDelete,
  handleEdit,
  handleRunTest,
  requestFullSize,
  requestSizeRestore,
}) {
  const classes = useStyles()
  const theme = useTheme()
  const [userAvatar, setUserAvatar] = useState(null)

  useEffect(() => {
    avatarHandler.getAvatar(profile.user_id, setUserAvatar)
  }, [])

  const {
    id,
    name,
    endpoints,
    load_generators : loadGenerators,
    total_results : results,
    duration : testRunDuration,
    concurrent_request : concurrentRequest,
    qps,
    service_mesh : serviceMesh,
    content_type : contentType,
    request_body : requestBody,
    request_cookies : requestCookies,
    request_headers : requestHeaders,
    last_run : lastRun
  } = profile
  const [renderTable, setRenderTable] = useState(false);


  const tableData = [
    {
      name : "Endpoints",
      value : endpoints?.join(", ")
    },
    {
      name : "Load Generators",
      value : loadGenerators?.join(", ")
    },
    {
      name : "Running Duration",
      value : testRunDuration
    },
    {
      name : "Concurrent Requests",
      value : concurrentRequest
    },
    {
      name : "queries/second",
      value : qps
    },
    {
      name : "Service Mesh",
      value : serviceMesh
    },
    {
      name : "Content Type",
      value : contentType,
      omitEmpty : true
    },
    {
      name : "Request Body",
      value : requestBody,
      omitEmpty : true
    },
    {
      name : "Cookies",
      value : requestCookies,
      omitEmpty : true
    },
    {
      name : "Request Headers",
      value : requestHeaders,
      omitEmpty : true
    },
    {
      name : "Created At",
      value : profile.created_at ? moment(profile.created_at).format("LLL"): "unknown"
    },
    {
      name : "Last Updated",
      value : profile.updated_at ? moment(profile.updated_at).format("LLL") : "unknown"
    },
    {
      name : "Last Run",
      value : profile.last_run ? moment(profile.last_run).format("LLL") : "unknown"
    }
  ]

  function genericClickHandler(ev, fn) {
    ev.stopPropagation();
    fn();
  }

  return (
    <FlipCard
      onClick={() => {
        setRenderTable(false);
        requestSizeRestore();
      }}
      duration={600}
    >
      {/* FRONT PART */}
      <>
        <div>
          <div style={{ display : "flex", justifyContent : "space-between" }}>
            <Typography variant="h6" component="div">
              {name}
            </Typography>
            <img src={`/static/img/load-test/${loadGenerators[0]}.svg`} height="24px" />
          </div>
          <div className={classes.noOfResultsContainer} >
            <div >
              <Typography variant="h2" component="div" style={{ marginRight : "0.75rem", color : `${theme.palette.type === "dark" ? "#fff" : "#647881"}` }}>
                {(results || "0").toLocaleString('en')}
              </Typography>
              <Typography variant="body1" className={classes.resultText} component="div">
                Results
              </Typography>
            </div>
          </div>
        </div>
        <div style={{ display : "flex", justifyContent : "space-between" }}>
          <div className={classes.bottomPart} >
            <Link href={`https://meshery.layer5.io/user/${profile.user_id}`} target="_blank">
              <Avatar alt="profile-avatar" src={userAvatar} />
            </Link>
            <div className={classes.lastRunText} >
              <div>
                {lastRun
                  ? (
                    <Typography variant="caption" style={{ fontStyle : "italic", color : `${theme.palette.type === "dark" ? "rgba(255, 255, 255, 0.7)" : "#647881"}` }}>
                      Last Run: <Moment format="LLL">{lastRun}</Moment>
                    </Typography>
                  )
                  : null}
              </div>
            </div>
          </div>
          <div className={classes.cardButtons} >
            <Button
              variant="contained"
              onClick={(ev) =>
                genericClickHandler(ev, () => {
                  setRenderTable((renderTable) => {
                    if (renderTable) {
                      requestSizeRestore();
                      return false;
                    }

                    requestFullSize();
                    return true;
                  });
                })
              }
              className={classes.testsButton}
            >
              {renderTable
                ? "Hide"
                : "View"} Results
            </Button>
            <Button
              color="primary"
              variant="contained"
              onClick={(ev) => genericClickHandler(ev, handleRunTest)}
            >
              Run Test
            </Button>
          </div>
        </div>
        {renderTable ? (
          <div onClick={(ev) => ev.stopPropagation()} className={classes.perfResultsContainer} >
            <PerformanceResults
              // @ts-ignore
              CustomHeader={<Typography variant="h6">Test Results</Typography>}
              // @ts-ignore
              endpoint={`/api/user/performance/profiles/${id}/results`}
              // @ts-ignore
              elevation={0}
            />
          </div>
        ) : null}
      </>

      {/* BACK PART */}
      <>
        <Grid className={classes.backGrid}
          container
          spacing={1}
          alignContent="space-between"
          alignItems="center"
        >
          <Grid item xs={8}>
            <Typography variant="h6" gutterBottom>
              {name} Details
            </Typography>
          </Grid>
          <Grid item xs={4}>
            <div className={classes.deleteEditButton} >
              <IconButton onClick={(ev) => genericClickHandler(ev, handleEdit)}>
                <EditIcon style={iconMedium} />
              </IconButton>
              <IconButton onClick={(ev) => genericClickHandler(ev, handleDelete)}>
                <DeleteIcon style={iconMedium} />
              </IconButton>
            </div>
          </Grid>
        </Grid>
        <Table size="small" dense >
          {
          // tableData.map(function renderDesignTableRow(({name, value, omitEmpty}), idx) {
          //   return <DetailsTable rowKey={} />
          // })
            tableData.map(function renderDesignTableRow(data) {
              const { name, value, omitEmpty } = data;
              return  <DetailsTable key={name} rowKey={name} value={value} omitEmpty={omitEmpty} />
            })
          }
        </Table>

      </>
    </FlipCard >
  );
}

// @ts-ignore
export default PerformanceCard;

function DetailsTable({ rowKey, value, omitEmpty }) {
  if (omitEmpty && (value === undefined || value === null)) {
    return null
  }

  return (
    <TableRow>
      <TableCell><b>{rowKey}</b></TableCell>
      <TableCell>{value || "none"}</TableCell>
    </TableRow>
  )
}
