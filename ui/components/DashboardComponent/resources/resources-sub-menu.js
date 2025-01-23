import React from 'react';
import { withStyles } from '@material-ui/core';
import { withRouter } from 'next/router';
import { withNotify } from '../../../utils/hooks/useNotification';
import ResourcesTable from './resources-table';
import { TabPanel } from '../tabpanel';
import { UsesSistent } from '@/components/SistentWrapper';
import { SecondaryTab, SecondaryTabs, WrapperContainer, WrapperPaper } from '../style';
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
    isCRDS,
  } = props;
  const CRDsModelName = isCRDS && CRDsKeys.map((key) => key.model);
  const CRDsKind = isCRDS && CRDsKeys.map((key) => key.name);
  if (!selectedResource) {
    let resourceNames;
    if (isCRDS) {
      resourceNames = CRDsKind;
    } else {
      resourceNames = Object.keys(resource.tableConfig());
    }
    handleChangeSelectedResource(resourceNames[0]);
  }

  let TABS;
  if (isCRDS) {
    TABS = CRDsKind;
  } else {
    TABS = Object.keys(resource.tableConfig());
  }

  if (TABS.length > 0 && selectedResource && !TABS.includes(selectedResource)) {
    handleChangeSelectedResource(TABS[0]);
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
              <SecondaryTabs
                value={getResourceCategoryIndex(selectedResource)}
                onChange={(_e, v) => handleChangeSelectedResource(getResourceCategory(v))}
                variant="scrollable"
                scrollButtons="on"
                indicatorColor="primary"
                textColor="primary"
              >
                {TABS.map((key, index) => {
                  const title = isCRDS ? key : resource.tableConfig()[key].name;
                  return (
                    <SecondaryTab
                      key={index}
                      value={index}
                      label={
                        <div className={classes.iconText}>
                          <GetKubernetesNodeIcon
                            kind={key}
                            model={CRDsModelName[index]}
                            size={iconMedium}
                          />
                          {title}
                        </div>
                      }
                    />
                  );
                })}
              </SecondaryTabs>
            </div>
          </WrapperPaper>
          {TABS.map((key, index) => (
            <TabPanel value={selectedResource} index={key} key={`${key}-${index}`}>
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
