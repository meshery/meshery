import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { updateGrafanaConfig, updatePrometheusConfig, updateTelemetryUrls } from '../../lib/store';
import { withStyles } from '@material-ui/core/styles';
import { withNotify, useNotification } from '../../utils/hooks/useNotification';
import { Tooltip, Tabs, Tab, Paper, Typography } from '@material-ui/core';
import { updateProgress } from '../../lib/store';
import { ResourcesConfig } from './resources/config';
import ResourcesTable from './resources/resources-table';
import ResourcesSubMenu from './resources/resources-sub-menu';
import Overview from './overview';
import KubernetesIcon from '../../assets/icons/technology/kubernetes';
import MesheryIcon from './images/meshery-icon.js';
import { EVENT_TYPES } from '../../lib/event-types';

const styles = (theme) => ({
  wrapperClss: {
    flexGrow: 1,
    maxWidth: '100%',
    height: 'auto',
  },
  tab: {
    minWidth: 40,
    paddingLeft: 0,
    paddingRight: 0,
    '&.Mui-selected': {
      color: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
  subMenuTab: {
    backgroundColor: theme.palette.type === 'dark' ? '#212121' : '#f5f5f5',
  },
  tabs: {
    '& .MuiTabs-indicator': {
      backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
  icon: {
    display: 'inline',
    verticalAlign: 'text-top',
    width: theme.spacing(1.75),
    marginLeft: theme.spacing(0.5),
  },
  iconText: {
    display: 'inline',
    verticalAlign: 'middle',
  },
  backToPlay: { margin: theme.spacing(2) },
  link: { cursor: 'pointer' },
  container: {
    display: 'flex',
    justifyContent: 'center',
    marginTop: theme.spacing(2),
  },
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  },
  topToolbar: {
    marginBottom: '2rem',
    display: 'flex',
    justifyContent: 'space-between',
    paddingLeft: '1rem',
    maxWidth: '90%',
  },
  dashboardSection: {
    backgroundColor: theme.palette.secondary.elevatedComponents,
    padding: theme.spacing(2),
    borderRadius: 4,
    height: '100%',
  },
  cardHeader: { fontSize: theme.spacing(2) },
  card: {
    height: '100%',
    marginTop: theme.spacing(2),
  },
  cardContent: { height: '100%' },
  boxWrapper: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'end',
    flexDirection: 'row',
    flexWrap: 'wrap',
    height: '60vh',
    borderRadius: 0,
    color: 'white',
    ['@media (max-width: 455px)']: {
      width: '100%',
    },
    zIndex: 5,
  },
  box: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    flexDirection: 'column',
    width: 300,
    height: 300,
    backgroundColor: theme.palette.secondary.dark,
    border: '0px solid #000',
    boxShadow: theme.shadows[5],
    margin: theme.spacing(2),
    cursor: 'pointer',
  },
});

const DashboardComponent = (props) => {
  const { classes, k8sconfig, selectedK8sContexts, updateProgress } = props;
  const { notify } = useNotification();

  const [tabVal, setTabVal] = React.useState(0);

  const handleChange = () => (event, newValue) => {
    setTabVal(newValue);
  };
  const fetchData = async () => {
    try {
      const response = await fetch('api/meshmodels/registrants');
      const data = await response.json();

      if (data.total_count > 0) {
        for (const registrant of data.registrants) {
          const { hostname, summary } = registrant;

          let successMessage = `For registrant ${hostname} successfully imported `;
          if (summary.models > 0) successMessage += ` ${summary.models} models`;
          if (summary.components > 0) successMessage += ` ${summary.components} components`;
          if (summary.relationships > 0)
            successMessage += ` ${summary.relationships} relationships`;
          if (summary.policies > 0) successMessage += ` ${summary.policies} policies`;

          notify({
            message: successMessage,
            event_type: EVENT_TYPES.SUCCESS,
          });

          const nonRegisteredResponse = await fetch('api/meshmodels/nonRegisterEntity', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ Hostname: hostname }),
          });

          const nonRegisteredData = await nonRegisteredResponse.json();
          console.log(nonRegisteredData);

          let errorMessage = `For registrant ${hostname} failed to import`;
          if (nonRegisteredData.summary.models > 0)
            errorMessage += ` ${nonRegisteredData.summary.models} models`;
          if (nonRegisteredData.summary.components > 0)
            errorMessage += ` ${nonRegisteredData.summary.components} components`;
          if (nonRegisteredData.summary.relationships > 0)
            errorMessage += ` ${nonRegisteredData.summary.relationships} relationships`;
          if (nonRegisteredData.summary.policies > 0)
            errorMessage += ` ${nonRegisteredData.summary.policies} policies`;

          notify({
            message: errorMessage,
            event_type: EVENT_TYPES.ERROR,
          });
        }
      }
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  function TabContainer(props) {
    return (
      <Typography component="div" style={{ paddingTop: 2 }}>
        {props.children}
      </Typography>
    );
  }
  return (
    <>
      <div className={classes.wrapperClss}>
        <Paper square className={classes.wrapperClss}>
          <Tabs
            value={tabVal}
            className={classes.tabs}
            onChange={handleChange()}
            variant="fullWidth"
            indicatorColor="primary"
            textColor="primary"
            allowScrollButtonsMobile
            scrollButtons="auto"
          >
            <Tooltip title={`View Overview`} placement="top">
              <Tab
                className={classes.tab}
                scrollButtons
                icon={<MesheryIcon style={{ width: '28px', height: '28px' }} />}
                label={'Overview'}
              />
            </Tooltip>

            {Object.keys(ResourcesConfig).map((resource, idx) => {
              return (
                <Tooltip title={`View ${resource}`} placement="top">
                  <Tab
                    key={idx}
                    className={classes.tab}
                    icon={<KubernetesIcon style={{ width: '28px', height: '28px' }} />}
                    label={resource}
                  />
                </Tooltip>
              );
            })}
          </Tabs>
        </Paper>
        {tabVal === 0 && (
          <TabContainer>
            <Overview />
          </TabContainer>
        )}
        {Object.keys(ResourcesConfig).map((resource, idx) => {
          return (
            tabVal === idx + 1 &&
            (ResourcesConfig[resource].submenu ? (
              <TabContainer>
                <ResourcesSubMenu
                  key={idx}
                  resource={ResourcesConfig[resource]}
                  updateProgress={updateProgress}
                  classes={classes}
                  k8sConfig={k8sconfig}
                  selectedK8sContexts={selectedK8sContexts}
                />
              </TabContainer>
            ) : (
              <TabContainer>
                <ResourcesTable
                  key={idx}
                  workloadType={resource}
                  classes={classes}
                  k8sConfig={k8sconfig}
                  selectedK8sContexts={selectedK8sContexts}
                  resourceConfig={ResourcesConfig[resource].tableConfig}
                  menu={ResourcesConfig[resource].submenu}
                  updateProgress={updateProgress}
                />
              </TabContainer>
            ))
          );
        })}
      </div>
    </>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  updateGrafanaConfig: bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig: bindActionCreators(updatePrometheusConfig, dispatch),
  updateTelemetryUrls: bindActionCreators(updateTelemetryUrls, dispatch),
});

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const meshAdapters = state.get('meshAdapters');
  const meshAdaptersts = state.get('meshAdaptersts');
  const grafana = state.get('grafana').toJS();
  const prometheus = state.get('prometheus').toJS();
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    meshAdapters,
    meshAdaptersts,
    k8sconfig,
    grafana,
    prometheus,
    selectedK8sContexts,
  };
};

export default withStyles(styles, { withTheme: true })(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withNotify(DashboardComponent))),
);
