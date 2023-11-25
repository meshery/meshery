import React from 'react';
import { withStyles } from '@material-ui/core';
import { Tooltip, Tabs, Tab, Typography } from '@material-ui/core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faMendeley } from '@fortawesome/free-brands-svg-icons';
import { iconSmall } from '../../../css/icons.styles';
import { withRouter } from 'next/router';
import { withNotify } from '../../../utils/hooks/useNotification';
import ResourcesTable from './resources-table';
import { Paper } from '@material-ui/core';
import { Box } from '@material-ui/core';

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

const ResourcesSubMenu = (props) => {
  const { classes, updateProgress, k8sConfig, resource, selectedK8sContexts } = props;

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
        <Paper className={classes.wrapperClss}>
          <div className={classes.subMenuTab}>
            <Box sx={{ margin: '0 auto', width: '100%', maxWidth: { xs: 490, sm: 880 } }}>
              <Tabs
                value={tabVal}
                className={classes.tabs}
                onChange={handleChange()}
                variant="scrollable"
                indicatorColor="primary"
                textColor="primary"
              >
                {Object.keys(resource.tableConfig()).map((key, index) => {
                  return (
                    <Tooltip title={`${resource.tableConfig()[key].name}`} placement="top">
                      <Tab
                        key={index}
                        label={
                          <div className={classes.iconText}>
                            {resource.tableConfig()[key].name}
                            <FontAwesomeIcon
                              icon={faMendeley}
                              style={iconSmall}
                              className={classes.icon}
                            />
                          </div>
                        }
                      />
                    </Tooltip>
                  );
                })}
              </Tabs>
            </Box>
          </div>
        </Paper>
        {Object.keys(resource.tableConfig()).map((key, index) => {
          return (
            tabVal === index && (
              <TabContainer>
                <ResourcesTable
                  key={index}
                  workloadType={key}
                  updateProgress={updateProgress}
                  classes={classes}
                  k8sConfig={k8sConfig}
                  resourceConfig={resource.tableConfig}
                  submenu={resource.submenu}
                  selectedK8sContexts={selectedK8sContexts}
                />
              </TabContainer>
            )
          );
        })}
      </div>
    </>
  );
};

export default withStyles(styles, { withTheme: true })(withRouter(withNotify(ResourcesSubMenu)));
