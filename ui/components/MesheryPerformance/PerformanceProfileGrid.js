//@ts-check
import { Grid } from "@material-ui/core";
import { Pagination } from "@material-ui/lab";
import React, { useState } from "react";
import PerformanceCard from "./PerformanceCard";
import { makeStyles } from "@material-ui/core/styles";


const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function PerformanceCardGridItem({ profile, deleteHandler, setProfileForModal }) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);

  return (
    <Grid item {...gridProps}>
      <PerformanceCard
        id={profile.id}
        name={profile.name}
        endpoints={profile.endpoints}
        loadGenerators={profile.load_generators}
        reqHeaders={profile.request_headers}
        results={profile.total_results || 0}
        testRunDuration={profile.duration}
        concurrentRequest={profile.concurrent_request}
        qps={profile.qps}
        serviceMesh={profile.service_mesh}
        contentType={profile.content_type}
        requestBody={profile.request_body}
        requestCookies={profile.request_cookies}
        requestHeaders={profile.request_headers}
        lastRun={profile.last_run}
        handleEdit={() => setProfileForModal(profile)}
        handleDelete={() => deleteHandler(profile.id)}
        handleRunTest={() => setProfileForModal({ ...profile, runTest : true })}
        requestFullSize={() => setGridProps({ xl : 12, md : 12, xs : 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
      />
    </Grid>
  );
}
const useStyles = makeStyles(() => ({
  pagination : {
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    marginTop : "2rem"
  }
}))
/**
 * PerformanceProfileGrid is the react component for rendering grid
 * @param {{
 *  profiles: Array<{
 *    id: string,
 *    created_at: string,
 *    updated_at: string,
 *    endpoints: Array<string>,
 *    load_generators: Array<string>,
 *    name: string,
 *    user_id: string,
 *    duration: string,
 *  }>,
 *  deleteHandler: (id: string) => void,
 *  setProfileForModal: (profile: any) => void,
 *  pages?: number,
 *  setPage: (page: number) => void
 * }} props props
 */


function PerformanceProfileGrid({
  profiles = [], deleteHandler, setProfileForModal, pages = 1, setPage
}) {
  const classes = useStyles()
  return (
    <div>
      <Grid container spacing={2} style={{ padding : "1rem" }}>
        {profiles.map((profile) => (
          <PerformanceCardGridItem
            key={profile.id}
            profile={profile}
            deleteHandler={deleteHandler}
            setProfileForModal={setProfileForModal}
          />
        ))}
      </Grid>
      {profiles.length
        ? (
          <div className={classes.pagination} >
            <Pagination count={pages} onChange={(_, page) => setPage(page - 1)} />
          </div>
        )
        : null}
    </div>
  );
}

export default PerformanceProfileGrid;