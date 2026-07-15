import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import PerformanceCalendar from './PerformanceCalendar';
import MesheryPerformanceComponent from './index';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import DefaultError from '@/components/general/error-404/index';
import { Modal, Button, Grid2, Paper, Typography, useTheme, styled } from '@sistent/sistent';
import { updateProgressAction } from '@/store/slices/mesheryUi';
import { useDispatch } from 'react-redux';
import { useGetPerformanceProfilesQuery } from '@/rtk-query/performance-profile';
import { useGetPerformanceResultsQuery } from '@meshery/schemas/mesheryApi';

const StyledPaper = styled(Paper)({
  padding: '1rem',
});

const StyledButton = styled(Button)(() => ({ padding: '0.5rem' }));

const ResultContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  justifyContent: 'space-between',
  [theme.breakpoints.down('md')]: {
    flexDirection: 'column',
  },
}));

const Separator = styled('div')(({ theme, vertical }) => ({
  ...(vertical
    ? {
        height: '10.4rem',
        width: '1px',
        background: 'black',
        marginTop: '1.1rem',
        backgroundColor: theme.palette.border?.normal,
        opacity: '0.7',
        [theme.breakpoints.down('md')]: {
          display: 'none',
        },
      }
    : {
        display: 'none',
        [theme.breakpoints.down('md')]: {
          display: 'block',
          width: '100%',
          height: '1px',
          background: 'black',
          marginTop: '1.1rem',
          backgroundColor: theme.palette.border?.normal,
          opacity: '0.7',
        },
      }),
}));

function Dashboard() {
  const [profiles, setProfiles] = useState({ count: 0, profiles: [] });
  const [tests, setTests] = useState({ count: 0, tests: [] });
  const [runTest, setRunTest] = useState(false);
  const { notify } = useNotification();
  const router = useRouter();
  const dispatch = useDispatch();
  const theme = useTheme();
  const {
    data: performanceProfilesData,
    isFetching: isFetchingProfiles,
    isError: isProfilesError,
    error: profilesError,
  } = useGetPerformanceProfilesQuery({
    page: 0,
    pagesize: 10,
    search: '',
    order: '',
  });
  const {
    data: performanceResultsData,
    isFetching: isFetchingResults,
    isError: isResultsError,
    error: resultsError,
  } = useGetPerformanceResultsQuery({
    page: '0',
    pagesize: '10',
    search: '',
    order: '',
    from: '',
    to: '',
  });

  useEffect(() => {
    dispatch(updateProgressAction({ showProgress: isFetchingProfiles || isFetchingResults }));
  }, [dispatch, isFetchingProfiles, isFetchingResults]);

  useEffect(() => {
    if (!performanceProfilesData) return;

    setProfiles({
      count: performanceProfilesData.totalCount || 0,
      profiles: performanceProfilesData.profiles || [],
    });
  }, [performanceProfilesData]);

  useEffect(() => {
    if (isProfilesError) handleError('Failed to Fetch Profiles')(profilesError);
  }, [isProfilesError, profilesError]);

  useEffect(() => {
    if (!performanceResultsData) return;

    setTests({
      count: performanceResultsData.totalCount || 0,
      tests: performanceResultsData.results || [],
    });
  }, [performanceResultsData]);

  useEffect(() => {
    if (isResultsError) handleError('Failed to Fetch Results')(resultsError);
  }, [isResultsError, resultsError]);

  function handleError(msg) {
    return function (error) {
      dispatch(updateProgressAction({ showProgress: false }));
      notify({
        message: `${msg}: ${error}`,
        event_type: EVENT_TYPES.ERROR,
        details: error.toString(),
      });
    };
  }

  return (
    <>
      {CAN(
        Keys.PerformanceManagementViewPerformanceProfiles.id,
        Keys.PerformanceManagementViewPerformanceProfiles.function,
      ) ? (
        <>
          <Grid2
            container
            spacing={2}
            style={{ padding: '0.5rem' }}
            size="grow"
            sx={{ alignContent: 'space-around' }}
            data-testid="performance-dashboard"
          >
            <Grid2 spacing={1} size={{ xs: 12, lg: 6 }} sx={{ flexDirection: 'column' }}>
              <Grid2>
                <StyledPaper>
                  <ResultContainer>
                    <StyledPaper sx={{ boxShadow: 'none' }}>
                      <div
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          height: '6.8rem',
                          color: theme.palette.text.default,
                        }}
                      >
                        <Typography
                          variant="h2"
                          component="div"
                          style={{ marginRight: '0.75rem', color: theme.palette.icon?.default }}
                        >
                          {tests.count.toLocaleString('en')}
                        </Typography>
                        <Typography
                          variant="body1"
                          component="div"
                          style={{ color: theme.palette.text.default }}
                        >
                          Total Tests Run
                        </Typography>
                      </div>
                      <div style={{ margin: '2rem 0 0 auto', width: 'fit-content' }}>
                        <StyledButton
                          onClick={() => setRunTest(true)}
                          disabled={
                            !CAN(
                              Keys.PerformanceManagementRunTest.id,
                              Keys.PerformanceManagementRunTest.function,
                            )
                          }
                          variant="contained"
                        >
                          Run Test
                        </StyledButton>
                      </div>
                    </StyledPaper>
                    <Separator vertical />
                    <Separator />
                    <StyledPaper sx={{ boxShadow: 'none' }}>
                      <div style={{ display: 'flex', alignItems: 'center', height: '6.8rem' }}>
                        <Typography
                          variant="h2"
                          component="div"
                          style={{ marginRight: '0.75rem', color: theme.palette.icon?.default }}
                        >
                          {profiles.count}
                        </Typography>
                        <Typography
                          variant="body1"
                          component="div"
                          style={{ color: theme.palette.text.default }}
                        >
                          Profiles
                        </Typography>
                      </div>
                      <div style={{ margin: '2rem 0 0 auto', width: 'fit-content' }}>
                        <StyledButton
                          variant="contained"
                          onClick={() => router.push('/performance/profiles')}
                        >
                          Manage Profiles
                        </StyledButton>
                      </div>
                    </StyledPaper>
                  </ResultContainer>
                </StyledPaper>
              </Grid2>
              <Grid2>
                <StyledPaper>
                  <PerformanceCalendar style={{ height: '40rem', margin: '2rem 0 0' }} />
                </StyledPaper>
              </Grid2>
            </Grid2>
          </Grid2>

          <Modal
            open={!!runTest}
            closeModal={() => setRunTest(false)}
            maxWidth="md"
            title="Performance Profile Wizard"
          >
            <MesheryPerformanceComponent />
          </Modal>
        </>
      ) : (
        <DefaultError />
      )}
    </>
  );
}

export default Dashboard;
