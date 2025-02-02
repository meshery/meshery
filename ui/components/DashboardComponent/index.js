import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { useRouter, withRouter } from 'next/router';
import { withNotify } from '../../utils/hooks/useNotification';
import { updateProgress } from '../../lib/store';
import { ResourcesConfig } from './resources/config';
import ResourcesTable from './resources/resources-table';
import ResourcesSubMenu from './resources/resources-sub-menu';
import Overview from './overview';
import KubernetesIcon from '../../assets/icons/technology/kubernetes';
import MesheryIcon from './images/meshery-icon.js';
import { TabPanel } from './tabpanel';
import { iconLarge } from '../../css/icons.styles';
import { useWindowDimensions } from '@/utils/dimension';
import { Tab, Tabs, CustomTooltip } from '@layer5/sistent';
import { UsesSistent } from '../SistentWrapper';
import { WrapperContainer, WrapperPaper } from './style';
import _ from 'lodash';

const useDashboardRouter = () => {
  const router = useRouter();
  const { query, push: pushRoute, route } = router;

  const resourceCategory = query.resourceCategory || 'Overview';
  const selectedResource = query.resource;

  const changeResourceTab = (resourceCategory) => {
    if (query.resourceCategory === resourceCategory) {
      return;
    }
    pushRoute(
      `${route}?resourceCategory=${resourceCategory || query.resourceCategory}`,
      undefined,
      { shallow: true },
    );
  };

  const handleChangeSelectedResource = (resource) => {
    if (query.resource === resource) {
      return;
    }
    pushRoute(`${route}?resourceCategory=${resourceCategory}&resource=${resource}`, undefined, {
      shallow: true,
    });
  };

  return { resourceCategory, changeResourceTab, selectedResource, handleChangeSelectedResource };
};

const ResourceCategoryTabs = ['Overview', ...Object.keys(ResourcesConfig)];

const DashboardComponent = ({ k8sconfig, selectedK8sContexts, updateProgress }) => {
  const { resourceCategory, changeResourceTab, selectedResource, handleChangeSelectedResource } =
    useDashboardRouter();

  const getResourceCategoryIndex = (resourceCategory) => {
    return ResourceCategoryTabs.findIndex((resource) => resource === resourceCategory);
  };

  const getResourceCategory = (index) => {
    return ResourceCategoryTabs[index];
  };

  const { width } = useWindowDimensions();

  if (!ResourceCategoryTabs.includes(resourceCategory)) {
    changeResourceTab('Overview');
  }
  return (
    <>
      <UsesSistent>
        <WrapperContainer>
          <WrapperPaper square>
            <Tabs
              value={getResourceCategoryIndex(resourceCategory)}
              indicatorColor="primary"
              onChange={(_e, val) => {
                changeResourceTab(getResourceCategory(val));
              }}
              variant={width < 1280 ? 'scrollable' : 'fullWidth'}
              scrollButtons="on"
              textColor="primary"
            >
              {ResourceCategoryTabs.map((resource, idx) => {
                return (
                  <CustomTooltip
                    key={`${resource}-${idx}`}
                    title={`View ${resource}`}
                    placement="top"
                  >
                    <Tab
                      key={resource}
                      icon={
                        resource === 'Overview' ? (
                          <MesheryIcon style={iconLarge} />
                        ) : (
                          <KubernetesIcon style={iconLarge} />
                        )
                      }
                      label={resource}
                    />
                  </CustomTooltip>
                );
              })}
            </Tabs>
          </WrapperPaper>

          <TabPanel value={resourceCategory} index={'Overview'}>
            <Overview />
          </TabPanel>

          {Object.keys(ResourcesConfig).map((resource, idx) => {
            let CRDsKeys = [];
            const isCRDS = resource === 'CRDS';
            if (isCRDS) {
              const TableValue = Object.values(
                ResourcesConfig[resource].tableConfig(
                  null,
                  null,
                  k8sconfig,
                  null,
                  resource,
                  selectedK8sContexts,
                ),
              );
              CRDsKeys = TableValue.map((item) => _.pick(item, ['name', 'model']));
            }

            return (
              <TabPanel value={resourceCategory} index={resource} key={`${resource}-${idx}`}>
                {ResourcesConfig[resource].submenu ? (
                  <ResourcesSubMenu
                    key={idx}
                    resource={ResourcesConfig[resource]}
                    selectedResource={selectedResource}
                    handleChangeSelectedResource={handleChangeSelectedResource}
                    updateProgress={updateProgress}
                    k8sConfig={k8sconfig}
                    selectedK8sContexts={selectedK8sContexts}
                    CRDsKeys={CRDsKeys}
                    isCRDS={isCRDS}
                  />
                ) : (
                  <ResourcesTable
                    key={idx}
                    workloadType={resource}
                    k8sConfig={k8sconfig}
                    selectedK8sContexts={selectedK8sContexts}
                    resourceConfig={ResourcesConfig[resource].tableConfig}
                    menu={ResourcesConfig[resource].submenu}
                    updateProgress={updateProgress}
                  />
                )}
              </TabPanel>
            );
          })}
        </WrapperContainer>
      </UsesSistent>
    </>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    k8sconfig,
    selectedK8sContexts,
  };
};

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withRouter(withNotify(DashboardComponent)));
