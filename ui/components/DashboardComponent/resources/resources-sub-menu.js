import React from 'react';
import { withRouter } from 'next/router';
import { withNotify } from '../../../utils/hooks/useNotification';
import ResourcesTable from './resources-table';
import { TabPanel } from '../tabpanel';

import { SecondaryTab, SecondaryTabs, WrapperContainer, WrapperPaper } from '../style';
import GetKubernetesNodeIcon from '../utils';
import { iconMedium } from 'css/icons.styles';
import { styled } from '@layer5/sistent';

const DashboardIconText = styled('div')({
  display: 'flex',
  flexWrap: 'no-wrap',
  justifyContent: 'center',
  gap: '1rem',
  alignItems: 'center',
  '& svg': {
    verticalAlign: 'middle',
    marginRight: '.5rem',
  },
});

const ResourcesSubMenu = (props) => {
  const {
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
      <WrapperContainer>
        <WrapperPaper>
          <SecondaryTabs
            value={getResourceCategoryIndex(selectedResource)}
            onChange={(_e, v) => handleChangeSelectedResource(getResourceCategory(v))}
            variant="scrollable"
            scrollButtons="on"
            indicatorColor="primary"
            textColor="primary"
            centered={true}
          >
            {TABS.map((key, index) => {
              const title = isCRDS ? key : resource.tableConfig()[key].name;
              return (
                <SecondaryTab
                  key={index}
                  value={index}
                  label={
                    <DashboardIconText>
                      <GetKubernetesNodeIcon
                        kind={key}
                        model={CRDsModelName[index]}
                        size={iconMedium}
                      />
                      {title}
                    </DashboardIconText>
                  }
                />
              );
            })}
          </SecondaryTabs>
        </WrapperPaper>
        {TABS.map((key, index) => (
          <TabPanel value={selectedResource} index={key} key={`${key}-${index}`}>
            <ResourcesTable
              key={index}
              workloadType={key}
              updateProgress={updateProgress}
              k8sConfig={k8sConfig}
              resourceConfig={resource.tableConfig}
              submenu={resource.submenu}
              selectedK8sContexts={selectedK8sContexts}
            />
          </TabPanel>
        ))}
      </WrapperContainer>
    </>
  );
};

export default withRouter(withNotify(ResourcesSubMenu));
