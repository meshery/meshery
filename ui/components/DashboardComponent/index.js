import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { withStyles } from '@material-ui/core/styles';
import { withNotify } from '../../utils/hooks/useNotification';
import { Tooltip, Tabs, Tab, Paper, Typography } from '@material-ui/core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPoll, faDatabase, faFileInvoice } from '@fortawesome/free-solid-svg-icons';
import { faMendeley } from '@fortawesome/free-brands-svg-icons';
import { updateProgress } from '../../lib/store';
import { iconMedium } from '../../css/icons.styles';
import Clusters from './clusters';
import Namespaces from './namespaces';
import Nodes from './nodes';
import Workloads from './workloads';
import Storage from './storage';
import Network from './network';
import Security from './security';
import Configuration from './configuration';

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
  console.log('DashboardComponent props: ', props);
  const { classes, k8sconfig } = props;

  const [tabVal, setTabVal] = React.useState(0);

  const handleChange = () => (event, newValue) => {
    setTabVal(newValue);
  };

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
          >
            <Tooltip title="Connect Meshery Adapters" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faMendeley} style={iconMedium} />}
                label="Clusters"
                data-cy="tabServiceMeshes"
              />
            </Tooltip>
            <Tooltip title="Configure Metrics backends" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faPoll} style={iconMedium} />}
                label="Nodes"
                tab="tabMetrics"
              />
            </Tooltip>
            <Tooltip title="Configure Metrics backends" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faPoll} style={iconMedium} />}
                label="Namespaces"
                tab="tabMetrics"
              />
            </Tooltip>
            <Tooltip title="Reset System" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faDatabase} style={iconMedium} />}
                label="Workloads"
                tab="systemReset"
              />
            </Tooltip>
            <Tooltip title="Registry" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faFileInvoice} style={iconMedium} />}
                label="Storage"
                tab="registry"
              />
            </Tooltip>
            <Tooltip title="Registry" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faFileInvoice} style={iconMedium} />}
                label="Network"
                tab="registry"
              />
            </Tooltip>
            <Tooltip title="Registry" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faFileInvoice} style={iconMedium} />}
                label="Security"
                tab="registry"
              />
            </Tooltip>
            <Tooltip title="Registry" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faFileInvoice} style={iconMedium} />}
                label="Configuration"
                tab="registry"
              />
            </Tooltip>
          </Tabs>
        </Paper>
        {tabVal === 0 && (
          <TabContainer>
            <Clusters />
          </TabContainer>
        )}
        {tabVal === 1 && (
          <TabContainer>
            <Nodes updateProgress={updateProgress} classes={classes} k8sConfig={k8sconfig} />
          </TabContainer>
        )}
        {tabVal === 2 && (
          <TabContainer>
            <Namespaces updateProgress={updateProgress} classes={classes} k8sConfig={k8sconfig} />
          </TabContainer>
        )}
        {tabVal === 3 && (
          <TabContainer>
            <Workloads />
          </TabContainer>
        )}
        {tabVal === 4 && (
          <TabContainer>
            <Storage />
          </TabContainer>
        )}
        {tabVal === 5 && (
          <TabContainer>
            <Network />
          </TabContainer>
        )}
        {tabVal === 6 && (
          <TabContainer>
            <Security />
          </TabContainer>
        )}
        {tabVal === 7 && (
          <TabContainer>
            <Configuration />
          </TabContainer>
        )}
      </div>
    </>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  // updateGrafanaConfig: bindActionCreators(updateGrafanaConfig, dispatch),
  // updatePrometheusConfig: bindActionCreators(updatePrometheusConfig, dispatch),
  // updateTelemetryUrls: bindActionCreators(updateTelemetryUrls, dispatch),
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
