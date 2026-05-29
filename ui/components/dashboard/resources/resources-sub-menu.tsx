import React, { useEffect, useMemo } from 'react';
import ResourcesTable from './resources-table';
import { TabPanel } from '../tabpanel';
import { TABS_SCROLL_BUTTONS_CLASS } from '../constants';

import { SecondaryTab, SecondaryTabs, WrapperPaper } from '../style';
import GetKubernetesNodeIcon from '../utils';
import { iconMedium } from 'css/icons.styles';
import { styled } from '@sistent/sistent';

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

type ResourcesSubMenuProps = {
  k8sConfig: unknown;
  resource: {
    tableConfig: (...args: any[]) => Record<string, { name: string }>;
    submenu: boolean;
  };
  selectedK8sContexts: unknown;
  selectedResource?: string;
  handleChangeSelectedResource: (resource: string) => void;
  CRDsKeys?: Array<{ name: string; model?: string }>;
  isCRDS?: boolean;
};

const ResourcesSubMenu = ({
  k8sConfig,
  resource,
  selectedK8sContexts,
  selectedResource,
  handleChangeSelectedResource,
  CRDsKeys = [],
  isCRDS = false,
}: ResourcesSubMenuProps) => {
  const crdsModelName = useMemo(
    () => (isCRDS ? CRDsKeys.map((key) => key.model) : []),
    [CRDsKeys, isCRDS],
  );
  const crdsKind = useMemo(
    () => (isCRDS ? CRDsKeys.map((key) => key.name) : []),
    [CRDsKeys, isCRDS],
  );

  // `tableConfig` is a custom hook (it calls useKubernetesHook internally), so it
  // must be invoked unconditionally on every render to keep hook order stable.
  // `resource` is a required prop, so no defensive guard is needed - and a guard
  // here would reintroduce a conditional hook call (the original crash).
  const tableConfigResult = resource.tableConfig();
  // Derive a stable string key so the `tabs` memo below isn't invalidated by the
  // fresh object identity returned on every render.
  const configKeysStr = Object.keys(tableConfigResult).join(',');
  const tabs = useMemo(
    () => (isCRDS ? crdsKind : configKeysStr ? configKeysStr.split(',') : []),
    [crdsKind, isCRDS, configKeysStr],
  );

  useEffect(() => {
    if (!tabs.length) {
      return;
    }

    if (!selectedResource || !tabs.includes(selectedResource)) {
      handleChangeSelectedResource(tabs[0]);
    }
  }, [handleChangeSelectedResource, selectedResource, tabs]);

  const selectedTabIndex = useMemo(
    () => tabs.findIndex((resourceName) => resourceName === selectedResource),
    [selectedResource, tabs],
  );

  return (
    <>
      <WrapperPaper>
        <SecondaryTabs
          sx={{
            [`& .${TABS_SCROLL_BUTTONS_CLASS}`]: {
              '&.Mui-disabled': { display: 'none' },
            },
            '& .MuiTabs-scroller': {
              flexGrow: '0',
            },
            justifyContent: 'center',
          }}
          value={selectedTabIndex}
          onChange={(_e, value) => handleChangeSelectedResource(tabs[value])}
          variant={'scrollable'}
          allowScrollButtonsMobile
          scrollButtons="auto"
          centered
        >
          {tabs.map((key, index) => {
            const title = isCRDS ? key : tableConfigResult[key].name;
            return (
              <SecondaryTab
                key={index}
                value={index}
                label={
                  <DashboardIconText>
                    <GetKubernetesNodeIcon
                      kind={key}
                      model={crdsModelName[index]}
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
      {tabs.map((key, index) => (
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
