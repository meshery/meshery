import React from 'react';
import { withStyles } from '@material-ui/core';
import KubernetesIcon from '../../../assets/icons/technology/kubernetes';
import { withRouter } from 'next/router';
import { withNotify } from '../../../utils/hooks/useNotification';
import ResourcesTable from './resources-table';
import { TabPanel } from '../tabpanel';
import { Box, CustomTooltip, Tab, Tabs } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';
import { WrapperContainer, WrapperPaper } from '../style';

const styles = (theme) => ({
  icon: {
    display: 'inline',
    verticalAlign: 'text-top',
    width: theme.spacing(1.75),
    marginLeft: theme.spacing(0.5),
  },

  iconText: {
    display: 'flex',
    flexWrap: 'no-wrap',
    justifyContent: 'center',
    alignItems: 'center',
    '& svg': {
      verticalAlign: 'middle',
      marginRight: '.5rem',
    },
  },
  backToPlay: { margin: theme.spacing(2) },
  link: { cursor: 'pointer' },
  container: {
    display: 'flex',
    justifyContent: 'center',
    marginTop: theme.spacing(2),
  },
  topToolbar: {
    marginBottom: '2rem',
    display: 'flex',
    justifyContent: 'space-between',
    paddingLeft: '1rem',
    maxWidth: '90%',
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
  const {
    classes,
    updateProgress,
    k8sConfig,
    resource,
    selectedK8sContexts,
    selectedResource,
    handleChangeSelectedResource,
  } = props;

  if (!selectedResource) {
    handleChangeSelectedResource(Object.keys(resource.tableConfig())[0]);
  }

  const TABS = Object.keys(resource.tableConfig());

  const getResourceCategoryIndex = (resourceCategory) => {
    return TABS.findIndex((resource) => resource === resourceCategory);
  };

  const getResourceCategory = (index) => {
    return TABS[index];
  };

  return (
    <>
      <UsesSistent>
        <WrapperContainer>
          <WrapperPaper>
            <div>
              <Box
                sx={{ margin: '0 auto', width: '100%', maxWidth: { xs: 800, sm: 880, md: 1200 } }}
              >
                <Tabs
                  value={getResourceCategoryIndex(selectedResource)}
                  onChange={(_e, v) => handleChangeSelectedResource(getResourceCategory(v))}
                  variant="scrollable"
                  scrollButtons="on"
                  indicatorColor="primary"
                  textColor="primary"
                  // centered
                >
                  {TABS.map((key, index) => (
                    <CustomTooltip
                      key={index}
                      title={`${resource.tableConfig()[key].name}`}
                      placement="top"
                    >
                      <Tab
                        key={index}
                        value={index}
                        label={
                          <div className={classes.iconText}>
                            <KubernetesIcon
                              className={classes.iconText}
                              width="22px"
                              height="22px"
                            />
                            {resource.tableConfig()[key].name}
                          </div>
                        }
                      />
                    </CustomTooltip>
                  ))}
                </Tabs>
              </Box>
            </div>
          </WrapperPaper>
          {TABS.map((key, index) => (
            <TabPanel value={selectedResource} index={key} key={index}>
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
            </TabPanel>
          ))}
        </WrapperContainer>
      </UsesSistent>
    </>
  );
};

export default withStyles(styles, { withTheme: true })(withRouter(withNotify(ResourcesSubMenu)));
