import { Grid } from '@material-ui/core';
import React, { useState } from 'react';
import PerformanceCard from './PerformanceCard';
import { makeStyles } from '@material-ui/core/styles';
import { Pagination } from '@layer5/sistent';
import { UsesSistent } from '../SistentWrapper';

const INITIAL_GRID_SIZE = { xl: 4, md: 6, xs: 12 };

function PerformanceCardGridItem({ profile, deleteHandler, setProfileForModal, testHandler }) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);

  return (
    <Grid item {...gridProps}>
      <PerformanceCard
        profile={profile}
        handleEdit={() => setProfileForModal(profile)}
        handleDelete={() => deleteHandler(profile.id)}
        handleProfile={() => setProfileForModal({ ...profile })}
        handleRunTest={() => testHandler({ ...profile, runTest: true })}
        requestFullSize={() => setGridProps({ xl: 12, md: 12, xs: 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
      />
    </Grid>
  );
}
const useStyles = makeStyles(() => ({
  pagination: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: '2rem',
  },
}));
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
  profiles = [],
  deleteHandler,
  setProfileForModal,
  pages = 1,
  setPage,
  testHandler,
}) {
  const classes = useStyles();
  return (
    <UsesSistent>
      <Grid container spacing={2} style={{ padding: '1rem' }}>
        {profiles.map((profile) => (
          <PerformanceCardGridItem
            key={profile.id}
            profile={profile}
            deleteHandler={deleteHandler}
            testHandler={testHandler}
            setProfileForModal={setProfileForModal}
          />
        ))}
      </Grid>
      {profiles.length ? (
        <div className={classes.pagination}>
          <Pagination count={pages} onChange={(_, page) => setPage(page - 1)} />
        </div>
      ) : null}
    </UsesSistent>
  );
}

export default PerformanceProfileGrid;
