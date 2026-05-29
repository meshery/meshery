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

// The *TableConfig functions are custom hooks (they call useKubernetesHook
// internally), so the field is `use`-prefixed to make that explicit and to let
// eslint-plugin-react-hooks enforce correct usage at every call site.
type ResourceTableConfigHook = (...args: any[]) => Record<string, { name: string; model?: string }>;

type ResourcesSubMenuProps = {
  k8sConfig: unknown;
  resource: {
    useTableConfig: ResourceTableConfigHook;
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

  // useTableConfig is a custom hook, so it is called unconditionally at the top
  // level to keep hook order stable across renders.
  const tableConfigResult = resource.useTableConfig();
  const tabs = useMemo(
    () => (isCRDS ? crdsKind : Object.keys(tableConfigResult)),
    [crdsKind, isCRDS, tableConfigResult],
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
            useResourceConfig={resource.useTableConfig}
            submenu={resource.submenu}
            selectedK8sContexts={selectedK8sContexts}
          />
        </TabPanel>
      ))}
    </>
  );
};

type CRDsResourcesSubMenuProps = Omit<ResourcesSubMenuProps, 'CRDsKeys' | 'isCRDS'>;

// Resolves the available CRD kinds by invoking the CRDs table-config hook at the
// top level of its own component. The parent renders this for the CRDs tab so
// the hook is never called from inside the parent's render loop/condition (which
// would violate the Rules of Hooks).
export const CRDsResourcesSubMenu = (props: CRDsResourcesSubMenuProps) => {
  const { resource, k8sConfig, selectedK8sContexts } = props;
  const crdsConfig = resource.useTableConfig(
    null,
    null,
    k8sConfig,
    null,
    'CRDS',
    selectedK8sContexts,
  );
  const CRDsKeys = Object.values(crdsConfig).map((item) => ({
    name: item.name,
    model: item.model,
  }));
  return <ResourcesSubMenu {...props} CRDsKeys={CRDsKeys} isCRDS />;
};

export default ResourcesSubMenu;
