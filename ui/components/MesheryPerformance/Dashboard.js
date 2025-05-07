import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { updateProgress } from '../../lib/store';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import MesheryMetrics from '../MesheryMetrics';
import PerformanceCalendar from './PerformanceCalendar';
import MesheryPerformanceComponent from './index';
import fetchPerformanceProfiles from '../graphql/queries/PerformanceProfilesQuery';
import fetchAllResults from '../graphql/queries/FetchAllResultsQuery';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import DefaultError from '@/components/General/error-404/index';
import {
  Modal,
  Button,
  Grid,
  Paper,
  Typography,
  useTheme,
  styled,
  useMediaQuery,
} from '@layer5/sistent';

// const MESHERY_PERFORMANCE_URL = "/api/user/performance/profiles";
// const MESHERY_PERFORMANCE_TEST_URL = "/api/user/performance/profiles/results";

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

function Dashboard({ updateProgress, grafana, router }) {
  const [profiles, setProfiles] = useState({ count: 0, profiles: [] });
  const [tests, setTests] = useState({ count: 0, tests: [] });
  const [runTest, setRunTest] = useState(false);
  const { notify } = useNotification();

  const theme = useTheme();
  const matches = useMediaQuery(theme.breakpoints.up('xs'));
  if (matches) {
    console.log('HIT');
  }

  /**
   * fetch performance profiles when the page loads
   */
  useEffect(() => {
    fetchTestProfiles();
    fetchTests();
  }, []);

  function fetchTestProfiles() {
    updateProgress({ showProgress: true });

    fetchPerformanceProfiles({
      selector: {
        // default
        pageSize: `10`,
        page: `0`,
        search: ``,
        order: ``,
      },
    }).subscribe({
      next: (res) => {
        // @ts-ignore
        let result = res?.getPerformanceProfiles;
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          if (result) {
            setProfiles({
              count: result.total_count || 0,
              profiles: result.profiles || [],
            });
          }
        }
      },
      error: handleError('Failed to Fetch Profiles'),
    });
  }

  function fetchTests() {
    updateProgress({ showProgress: true });

    fetchAllResults({
      selector: {
        // default
        pageSize: `10`,
        page: `0`,
        search: ``,
        order: ``,
        from: ``,
        to: ``,
      },
    }).subscribe({
      next: (res) => {
        // @ts-ignore
        let result = res?.fetchAllResults;
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          if (result) {
            setTests({
              count: result.total_count || 0,
              tests: result.results || [],
            });
          }
        }
      },
      error: handleError('Failed to Fetch Results'),
    });
  }

  function handleError(msg) {
    return function (error) {
      updateProgress({ showProgress: false });
      notify({
        message: `${msg}: ${error}`,
        event_type: EVENT_TYPES.ERROR,
        details: error.toString(),
      });
    };
  }

  return (
    <>
      {CAN(keys.VIEW_PERFORMANCE_PROFILES.action, keys.VIEW_PERFORMANCE_PROFILES.subject) ? (
        <>
          <Grid container spacing={2} style={{ padding: '0.5rem' }} alignContent="space-around">
            <Grid container item spacing={1} direction="column" lg xs={12}>
              <Grid item>
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
                          disabled={!CAN(keys.RUN_TEST.action, keys.RUN_TEST.subject)}
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
              </Grid>
              <Grid item>
                <StyledPaper>
                  <PerformanceCalendar style={{ height: '40rem', margin: '2rem 0 0' }} />
                </StyledPaper>
              </Grid>
            </Grid>
            <Grid item lg xs={12}>
              <StyledPaper style={{ height: '100%' }}>
                <MesheryMetrics
                  boardConfigs={grafana.selectedBoardsConfigs}
                  grafanaURL={grafana.grafanaURL}
                  grafanaAPIKey={grafana.grafanaAPIKey}
                  handleGrafanaChartAddition={() => router.push('/settings/#metrics')}
                />
              </StyledPaper>
            </Grid>
          </Grid>

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

const mapStateToProps = (st) => {
  const grafana = st.get('grafana').toJS();
  return { grafana: { ...grafana, ts: new Date(grafana.ts) } };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(Dashboard));
