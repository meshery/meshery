// @ts-nocheck
import React, { useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/router';
import { useNotificationHandlers } from '../../utils/hooks/useNotification';
import { ResourcesConfig } from './resources/config';
import ResourcesTable from './resources/resources-table';
import ResourcesSubMenu from './resources/resources-sub-menu';
import KubernetesIcon from '../../assets/icons/technology/kubernetes';
import MesheryIcon from './images/meshery-icon';
import { TabPanel } from './tabpanel';
import { iconLarge } from '../../css/icons.styles';
import { useWindowDimensions } from '@/utils/dimension';
import {
  Tab,
  Tabs,
  CustomTooltip,
  Box,
  Stack,
  EditIcon,
  CloseIcon,
  SaveAsIcon,
  OutlinedValidateIcon,
  OutlinedResetIcon,
  useTheme,
  ErrorBoundary,
} from '@sistent/sistent';
import { WrapperPaper } from './style';
import _ from 'lodash';
import { AddWidgetsToLayoutPanel, LayoutActionButton, LayoutWidget } from './components';
import { Responsive, WidthProvider } from 'react-grid-layout/legacy';
import { DEFAULT_LAYOUT, LOCAL_PROVIDER_LAYOUT, OVERVIEW_LAYOUT } from './defaultLayout';
import { applyMinSizeConstraints } from './layoutConstraints';
import Popup from '../General/Popup';
import { useGetUserPrefQuery, useUpdateUserPrefMutation } from '@/rtk-query/user';
import getWidgets from './widgets/getWidgets';
import { tabsClasses } from '@mui/material';
import { useSelector } from 'react-redux';
import useUnsavedChanges from './useUnsavedChanges';
import UnsavedChangesModal from './UnsavedChangesModal';

const ResponsiveReactGridLayout = WidthProvider(Responsive);

const cols = { lg: 12, md: 10, sm: 6, xs: 4, xxs: 2 };

const useDashboardRouter = () => {
  const router = useRouter();
  const { query, push: pushRoute, route, isReady } = router;

  // Use 'Overview' until router is hydrated with query params
  const resourceCategory = isReady ? query.resourceCategory || 'Overview' : 'Overview';
  const selectedResource = isReady ? query.resource : undefined;

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
const GRID_COLS = { lg: 12, md: 10, sm: 6, xs: 4, xxs: 2 };
const BREAKPOINTS = Object.keys(GRID_COLS);

const overlaps = (a, b) => a.x < b.x + b.w && b.x < a.x + a.w && a.y < b.y + b.h && b.y < a.y + a.h;

const findCollision = (candidate, items) => items.find((item) => overlaps(candidate, item)) || null;

const compactLayout = (items) => {
  const sorted = [...items].sort((a, b) => (a.y !== b.y ? a.y - b.y : a.x - b.x));

  const placeAtFirstFreeRow = (item, placed, y = 0) => {
    const collision = findCollision({ ...item, y }, placed);
    if (!collision) {
      return { ...item, y };
    }
    return placeAtFirstFreeRow(item, placed, Math.max(y + 1, collision.y + collision.h));
  };

  return sorted.reduce((placed, item) => {
    return [...placed, placeAtFirstFreeRow(item, placed)];
  }, []);
};

const isLayoutValid = (items, breakpoint) => {
  if (!Array.isArray(items)) return false;
  const maxCols = GRID_COLS[breakpoint];

  const coordinatesAreValid = items.every(({ x, y, w, h }) => {
    return (
      Number.isFinite(x + y + w + h) && x >= 0 && y >= 0 && w >= 1 && h >= 1 && x + w <= maxCols
    );
  });

  if (!coordinatesAreValid) {
    return false;
  }

  const hasCollisions = items.some((item, index) => {
    return items.slice(index + 1).some((peer) => overlaps(item, peer));
  });

  return !hasCollisions;
};

const Dashboard = () => {
  const { data: userData, isLoading } = useGetUserPrefQuery();
  const [updateUserPref] = useUpdateUserPrefMutation();
  const defaultLayout = isLoading
    ? OVERVIEW_LAYOUT
    : userData?.remoteProviderPreferences
      ? DEFAULT_LAYOUT
      : LOCAL_PROVIDER_LAYOUT; //TODO: Use capability to determine default layout
  const { resourceCategory, changeResourceTab, selectedResource, handleChangeSelectedResource } =
    useDashboardRouter();

  const getResourceCategoryIndex = (resourceCategory) => {
    return ResourceCategoryTabs.findIndex((resource) => resource === resourceCategory);
  };

  const getResourceCategory = (index) => {
    return ResourceCategoryTabs[index];
  };

  const { width } = useWindowDimensions();
  const theme = useTheme();

  if (!ResourceCategoryTabs.includes(resourceCategory)) {
    changeResourceTab('Overview');
  }
  const [currentBreakPoint, setCurrentBreakpoint] = useState('lg');
  const { selectedK8sContexts } = useSelector((state) => state.ui);
  const { k8sConfig } = useSelector((state) => state.ui);
  const [isEditMode, setIsEditMode] = useState(false);

  const iconsProps = useMemo(
    () => ({
      fill: theme.palette.icon.default,
      primaryFill: theme.palette.icon.default,
      secondaryFill: theme.palette.icon.secondary,
      width: '40',
    }),
    [theme.palette.icon.default, theme.palette.icon.secondary],
  );

  const WIDGETS = useMemo(() => getWidgets({ iconsProps, isEditMode }), [iconsProps, isEditMode]);
  const widgetSizing = useMemo(
    () =>
      Object.fromEntries(
        Object.entries(WIDGETS).map(([key, widget]) => [key, widget.defaultSizing]),
      ),
    [WIDGETS],
  );
  const availableHandles = ['s', 'w', 'e', 'n', 'sw', 'nw', 'se', 'ne'];
  const cols = GRID_COLS;

  const clampLayoutItem = (item, breakpoint) => {
    const maxCols = GRID_COLS[breakpoint];
    const sizing = WIDGETS[item.i]?.defaultSizing;

    if (sizing) {
      const w = Math.min(maxCols, Math.max(1, sizing.w || 1));
      const h = Math.max(1, sizing.h || 1);
      const x = Math.min(Math.max(item.x || 0, 0), maxCols - w);
      return { ...item, w, h, x, minW: w, maxW: w, minH: h, maxH: h };
    }

    const w = Math.min(maxCols, Math.max(1, item.w || 1));
    const x = Math.min(Math.max(item.x || 0, 0), maxCols - w);
    const h = Math.max(1, item.h || 1);
    return { ...item, w, h, x, minW: 1, maxW: maxCols, minH: 1, maxH: Infinity };
  };

  const normalizeLayouts = (layouts) =>
    Object.fromEntries(
      BREAKPOINTS.map((bp) => [bp, (layouts[bp] || []).map((item) => clampLayoutItem(item, bp))]),
    );

  const sanitizeLayouts = (layouts) =>
    Object.fromEntries(
      BREAKPOINTS.map((bp) => {
        const items = layouts[bp] || [];
        const normalized = items.map((item) => clampLayoutItem(item, bp));
        const compacted = compactLayout(normalized);
        const compactedById = Object.fromEntries(compacted.map((item) => [item.i, item]));
        const hasGaps = normalized.some(
          (item) => item.y - (compactedById[item.i]?.y ?? item.y) >= 2,
        );
        return [bp, !isLayoutValid(items, bp) || hasGaps ? compacted : normalized];
      }),
    );
  const isWidgetAlreadyAdded = (key, layout, breakpoint) => {
    return (layout[breakpoint] || []).some((item) => item.i === key);
  };

  const getWidgetsAvailableToBeAdded = (layout, breakpoint) => {
    return Object.entries(WIDGETS)
      .map(([key, value]) => ({ key, ...value }))
      .filter(
        (widget) => widget?.isEnabled?.() && !isWidgetAlreadyAdded(widget.key, layout, breakpoint),
      );
  };
  const orgDashboardLayout = useMemo(
    () => sanitizeLayouts(userData?.dashboardPreferences || defaultLayout),
    [userData?.dashboardPreferences, defaultLayout],
  );
  const [dashboardLayout, setDashboardLayout] = useState(orgDashboardLayout);

  const {
    showModal: showUnsavedModal,
    confirmNavigation,
    cancelNavigation,
  } = useUnsavedChanges({
    isEditMode,
    dashboardLayout,
    savedLayout: orgDashboardLayout,
  });

  const handleDiscardAndNavigate = () => {
    cancelEditing();
    confirmNavigation();
  };

  const handleSaveAndNavigate = async () => {
    const isSaved = await saveLayout();
    if (!isSaved) {
      return;
    }
    setIsEditMode(false);
    confirmNavigation();
  };

  const widgetsToAdd = getWidgetsAvailableToBeAdded(dashboardLayout, currentBreakPoint);

  const editModeStyles = {
    padding: '0.87rem',
    paddingBottom: '4rem',
    backgroundColor: theme.palette?.background?.tabs,
  };

  const onAddWidget = (widget, key) => {
    const updated = Object.fromEntries(
      BREAKPOINTS.map((bp) => {
        const newItem = clampLayoutItem(
          { i: key, x: 0, y: Infinity, static: false, moved: false, ...widget.defaultSizing },
          bp,
        );
        return [bp, [...(dashboardLayout[bp] || []), newItem]];
      }),
    );
    setDashboardLayout(sanitizeLayouts(updated));
  };
  const { handleError, handleSuccess } = useNotificationHandlers();

  const updateLayout = async (dashboardLayout) => {
    const constrainedLayoutToSave = applyMinSizeConstraints(
      dashboardLayout,
      defaultLayout,
      cols,
      widgetSizing,
    );
    const res = await updateUserPref({ dashboardPreferences: constrainedLayoutToSave });
    if (res.error) {
      handleError('failed to save layout');
      return false;
    }
    handleSuccess('Layout saved');
    return true;
  };

  const toggleEditMode = () => {
    setIsEditMode(!isEditMode);
  };
  const cancelEditing = () => {
    setIsEditMode(false);
    setDashboardLayout(orgDashboardLayout);
  };

  const saveLayout = () => {
    return updateLayout(sanitizeLayouts(dashboardLayout));
  };

  const resetLayout = () => {
    const normalizedDefaultLayout = sanitizeLayouts(defaultLayout);
    updateLayout(normalizedDefaultLayout);
  };
  const LayoutActions = {
    START_EDIT: {
      label: 'Edit',
      Icon: EditIcon,
      action: toggleEditMode,
      description: 'Configure dashboard layout for the current organization',
      isShown: !isEditMode,
    },
    CANCEL_EDIT: {
      label: 'Cancel',
      Icon: CloseIcon,
      action: cancelEditing,
      description: 'Stop editing and discard any unsaved changes',
      isShown: isEditMode,
    },

    SAVE_LAYOUT: {
      label: 'Save',
      Icon: SaveAsIcon,
      description: 'Save the current layout',
      action: saveLayout,
      isShown: isEditMode,
    },
    SAVE_AND_CLOSE: {
      label: 'Save and Close',
      Icon: OutlinedValidateIcon,
      action: async () => {
        const isSaved = await saveLayout();
        if (!isSaved) {
          return;
        }
        toggleEditMode();
      },

      isShown: isEditMode,
      description: 'Save the current layout and close edit mode',
    },

    RESET_LAYOUT: {
      label: 'Reset',
      description: 'Delete custom configuration and reset to default layout',
      Icon: OutlinedResetIcon,
      action: resetLayout,
      isShown: isEditMode,
    },
  };

  const topBarActions = Object.entries(_.omit(LayoutActions, 'START_EDIT'))
    .filter(([, action]) => action.isShown)
    .map(([key, layoutAction]) => ({ key, ...layoutAction }));

  const onBreakpointChange = (breakpoint) => {
    setCurrentBreakpoint(breakpoint);
  };
  useEffect(() => {
    setDashboardLayout(orgDashboardLayout);
  }, [orgDashboardLayout]);

  const onLayoutChange = (layout, layouts) => {
    setDashboardLayout(normalizeLayouts(layouts));
  };

  const widgetsToRenderForLayout = (layout, breakpoint) => {
    return (layout[breakpoint] || [])
      .filter((layoutItem) => WIDGETS[layoutItem.i]?.isEnabled?.())
      .map((layoutItem) => ({
        key: layoutItem.i,
        ...(WIDGETS[layoutItem.i] || {}),
      }));
  };

  const removeWidget = (key) => {
    setDashboardLayout((currentLayouts) => {
      const updatedLayout = (currentLayouts[currentBreakPoint] || []).filter((w) => w.i !== key);
      const updatedLayouts = _.set(currentLayouts, currentBreakPoint, updatedLayout);
      return { ...updatedLayouts };
    });
  };

  const constrainedLayouts = useMemo(
    () => applyMinSizeConstraints(dashboardLayout, defaultLayout, cols, widgetSizing),
    [dashboardLayout, defaultLayout, widgetSizing],
  );

  return (
    <>
      <>
        <WrapperPaper>
          <Tabs
            sx={{
              [`& .${tabsClasses.scrollButtons}`]: {
                '&.Mui-disabled': { display: 'none' },
              },
            }}
            value={getResourceCategoryIndex(resourceCategory)}
            indicatorColor="primary"
            onChange={(_e, val) => {
              changeResourceTab(getResourceCategory(val));
            }}
            variant={width < 1080 ? 'scrollable' : 'fullWidth'}
            allowScrollButtonsMobile
            scrollButtons
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
                    style={{
                      gap: '0.4rem',
                    }}
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
          <Box display="flex" flexDirection={'column'} gap="1rem">
            <Box padding={0} width={'100%'}>
              <Stack
                direction="row"
                useFlexGap
                gap="0rem 2rem"
                justifyContent="end"
                flexWrap={'wrap-reverse'}
              >
                {topBarActions.map(({ key, ...layoutAction }) => (
                  <LayoutActionButton {...layoutAction} key={key} />
                ))}
              </Stack>

              <ResponsiveReactGridLayout
                layouts={constrainedLayouts}
                resizeHandles={availableHandles}
                isResizable={isEditMode}
                isDraggable={isEditMode}
                cols={cols}
                draggableHandle=".react-grid-dragHandleExample"
                breakpoints={{ lg: 1200, md: 996, sm: 768, xs: 480, xxs: 0 }}
                onBreakpointChange={onBreakpointChange}
                onLayoutChange={onLayoutChange}
                compactType="vertical"
                measureBeforeMount={false}
                style={{
                  backgroundColor: 'transparent',
                }}
                containerPadding={[0, 8]}
              >
                {widgetsToRenderForLayout(dashboardLayout, currentBreakPoint).map((widget) => {
                  return (
                    <div key={widget.key} style={isEditMode ? editModeStyles : {}}>
                      <ErrorBoundary>
                        <LayoutWidget
                          isEditMode={isEditMode}
                          key={widget.key}
                          widget={widget}
                          removeWidget={removeWidget}
                        />
                      </ErrorBoundary>
                    </div>
                  );
                })}
              </ResponsiveReactGridLayout>
              <LayoutActionButton {...LayoutActions.START_EDIT} />
            </Box>
            <AddWidgetsToLayoutPanel
              editMode={isEditMode}
              widgetsToAdd={widgetsToAdd}
              onAddWidget={onAddWidget}
            />
          </Box>
        </TabPanel>

        {Object.keys(ResourcesConfig).map((resource, idx) => {
          let CRDsKeys = [];
          const isCRDS = resource === 'CRDS';
          if (isCRDS) {
            const TableValue = Object.values(
              ResourcesConfig[resource].tableConfig(
                null,
                null,
                k8sConfig,
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
                  k8sConfig={k8sConfig}
                  selectedK8sContexts={selectedK8sContexts}
                  CRDsKeys={CRDsKeys}
                  isCRDS={isCRDS}
                />
              ) : (
                <ResourcesTable
                  key={idx}
                  workloadType={resource}
                  k8sConfig={k8sConfig}
                  selectedK8sContexts={selectedK8sContexts}
                  resourceConfig={ResourcesConfig[resource].tableConfig}
                  menu={ResourcesConfig[resource].submenu}
                />
              )}
            </TabPanel>
          );
        })}
      </>
      <Popup />
      <UnsavedChangesModal
        open={showUnsavedModal}
        onClose={cancelNavigation}
        onDiscard={handleDiscardAndNavigate}
        onSave={handleSaveAndNavigate}
      />
    </>
  );
};

export default Dashboard;
