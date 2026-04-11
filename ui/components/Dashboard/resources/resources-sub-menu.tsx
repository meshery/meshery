// @ts-nocheck

import React from 'react';
import ResourcesTable from './resources-table';
import { TabPanel } from '../tabpanel';

import { SecondaryTab, SecondaryTabs, WrapperPaper } from '../style';
import GetKubernetesNodeIcon from '../utils';
import { iconMedium } from 'css/icons.styles';
import { styled } from '@sistent/sistent';
import { tabsClasses } from '@mui/material';

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

  // Call tableConfig() unconditionally to keep hook count stable across renders.
  // Config functions (e.g. WorkloadTableConfig) contain React hooks internally,
  // so they must always be called regardless of whether the result is used.
  const tableConfigResult = resource.tableConfig();
  const TABS = isCRDS ? CRDsKind : Object.keys(tableConfigResult);

  if (!selectedResource && TABS.length > 0) {
    handleChangeSelectedResource(TABS[0]);
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
      <WrapperPaper>
        <SecondaryTabs
          sx={{
            [`& .${tabsClasses.scrollButtons}`]: {
              '&.Mui-disabled': { display: 'none' },
            },
            '& .MuiTabs-scroller': {
              flexGrow: '0',
            },
            justifyContent: 'center',
          }}
          value={getResourceCategoryIndex(selectedResource)}
          onChange={(_e, v) => handleChangeSelectedResource(getResourceCategory(v))}
          variant={'scrollable'}
          allowScrollButtonsMobile
          scrollButtons="auto"
          centered
        >
          {TABS.map((key, index) => {
            const title = isCRDS ? key : tableConfigResult[key].name;
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
            k8sConfig={k8sConfig}
            resourceConfig={resource.tableConfig}
            submenu={resource.submenu}
            selectedK8sContexts={selectedK8sContexts}
          />
        </TabPanel>
      ))}
    </>
  );
};

export default ResourcesSubMenu;
