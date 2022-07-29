import React, {useState} from 'react'
import { Button, Grid, IconButton, Typography } from "@mui/material";
import PerformanceCard from "./PerformanceCard";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function PerformanceCardGridItem({ profile }) {
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
        />
      </Grid>
    );
  }

  function PerformanceProfileGrid({
    profiles = [], deleteHandler, setProfileForModal,
  }) {

    return (
      <div>
        <Grid container spacing={2} style={{ padding : "1rem" }}>
          {profiles.map((profile) => (
            <PerformanceCardGridItem
              key={profile.id}
              profile={profile}
            />
          ))}
        </Grid>
      </div>
    );
  }

export default PerformanceProfileGrid