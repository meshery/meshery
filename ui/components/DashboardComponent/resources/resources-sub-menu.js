import React from 'react';
import { withStyles } from '@material-ui/core';
import { withRouter } from 'next/router';
import { withNotify } from '../../../utils/hooks/useNotification';
import ResourcesTable from './resources-table';
import { TabPanel } from '../tabpanel';
import { Box, CustomTooltip, Tab, Tabs } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';
import { WrapperContainer, WrapperPaper } from '../style';
import GetKubernetesNodeIcon from '../utils';
import { iconMedium } from 'css/icons.styles';

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
    gap: '1rem',
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
    CRDsKeys,
  } = props;
  const isCRD = CRDsKeys.length > 0;
  if (!selectedResource) {
    let resourceNames = Object.keys(resource.tableConfig());
    if (isCRD) {
      resourceNames = CRDsKeys;
    }
    handleChangeSelectedResource(resourceNames[0]);
  }

  let TABS;
  if (isCRD) {
    TABS = CRDsKeys;
  } else {
    TABS = Object.keys(resource.tableConfig());
  }

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
                  {TABS.map((key, index) => {
                    const title = isCRD ? key : resource.tableConfig()[key].name;
                    return (
                      <CustomTooltip key={index} title={title} placement="top">
                        <Tab
                          key={index}
                          value={index}
                          label={
                            <div className={classes.iconText}>
                              <GetKubernetesNodeIcon kind={key} isCRD={isCRD} size={iconMedium} />
                              {title}
                            </div>
                          }
                        />
                      </CustomTooltip>
                    );
                  })}
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
