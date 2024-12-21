import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { useRouter, withRouter } from 'next/router';
import { withNotify } from '../../utils/hooks/useNotification';
import { Tabs, Tab, Paper, styled, useTheme } from '@layer5/sistent';
import { updateProgress } from '../../lib/store';
import { ResourcesConfig } from './resources/config';
import ResourcesTable from './resources/resources-table';
import ResourcesSubMenu from './resources/resources-sub-menu';
import Overview from './overview';
import KubernetesIcon from '../../assets/icons/technology/kubernetes';
import MesheryIcon from './images/meshery-icon.js';
import { TabPanel } from './tabpanel';
import { CustomTextTooltip } from '../MesheryMeshInterface/PatternService/CustomTextTooltip';
import { iconLarge } from '../../css/icons.styles';
import { useWindowDimensions } from '@/utils/dimension';
import { UsesSistent } from '@/components/SistentWrapper';

const StyledWrapper = styled('div')(() => ({
  flexGrow: 1,
  maxWidth: '100vw',
  height: 'auto',
}));

const StyledTabs = styled(Tabs)(() => {
  const theme = useTheme();
  return {
    width: '100%',
    '& .MuiTabs-indicator': {
      backgroundColor:
        theme.palette.mode === 'dark'
          ? theme.palette.background.neutral.default
          : theme.palette.primary.main,
    },
  };
});

const StyledTab = styled(Tab)(({ theme }) => ({
  width: 'max(6rem, 20%)',
  margin: 0,
  minWidth: 40,
  paddingLeft: 0,
  paddingRight: 0,
  '&.Mui-selected': {
    color:
      theme.palette.mode === 'dark'
        ? theme.palette.background.neutral.default
        : theme.palette.primary.main,
  },
}));

const StyledPaper = styled(Paper)(() => ({
  maxWidth: '90%',
  margin: 'auto',
  overflow: 'hidden',
}));

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
const DashboardComponent = ({ classes, k8sconfig, selectedK8sContexts, updateProgress }) => {
  const { resourceCategory, changeResourceTab, selectedResource, handleChangeSelectedResource } =
    useDashboardRouter();

  const getResourceCategoryIndex = (resourceCategory) => {
    return ResourceCategoryTabs.findIndex((resource) => resource === resourceCategory);
  };

  const getResourceCategory = (index) => {
    return ResourceCategoryTabs[index];
  };
  const { width } = useWindowDimensions();

  return (
    <UsesSistent>
      <StyledWrapper>
        <StyledPaper square>
          <StyledTabs
            value={getResourceCategoryIndex(resourceCategory)}
            indicatorColor="primary"
            onChange={(_e, val) => {
              changeResourceTab(getResourceCategory(val));
            }}
            variant={width < 1280 ? 'scrollable' : 'fullWidth'}
            scrollButtons="on"
            textColor="primary"
            centered
          >
            {ResourceCategoryTabs.map((resource, idx) => {
              return (
                <CustomTextTooltip key={idx} title={`View ${resource}`} placement="top">
                  <StyledTab
                    value={idx}
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
                </CustomTextTooltip>
              );
            })}
          </StyledTabs>
        </StyledPaper>

        <TabPanel value={resourceCategory} index={'Overview'}>
          <Overview />
        </TabPanel>
        {Object.keys(ResourcesConfig).map((resource, idx) => (
          <TabPanel value={resourceCategory} index={resource} key={resource}>
            {ResourcesConfig[resource].submenu ? (
              <ResourcesSubMenu
                key={idx}
                resource={ResourcesConfig[resource]}
                selectedResource={selectedResource}
                handleChangeSelectedResource={handleChangeSelectedResource}
                updateProgress={updateProgress}
                classes={classes}
                k8sConfig={k8sconfig}
                selectedK8sContexts={selectedK8sContexts}
              />
            ) : (
              <ResourcesTable
                key={idx}
                workloadType={resource}
                classes={classes}
                k8sConfig={k8sconfig}
                selectedK8sContexts={selectedK8sContexts}
                resourceConfig={ResourcesConfig[resource].tableConfig}
                menu={ResourcesConfig[resource].submenu}
                updateProgress={updateProgress}
              />
            )}
          </TabPanel>
        ))}
      </StyledWrapper>
    </UsesSistent>
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
