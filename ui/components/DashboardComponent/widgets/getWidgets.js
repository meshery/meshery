import WorkspaceActivityThumbnail from './Thumbnails/WorkspaceActivity.svg';
import GettingStartedThumbnail from './Thumbnails/GettingStartedProgress.svg';
import ClusterOverviewThumbnail from './Thumbnails/ClusterOverview.png';
import HelpCenterThumbnail from './Thumbnails/HelpCenter.svg';
import MyDesignthumbnail from './Thumbnails/MyDesigns.png';
import ClusterStatusThumbnail from './Thumbnails/ClusterStatus.png';
import LatestBlogsThumbnail from './Thumbnails/LatestBlogs.png';
import Overview from '../overview';
import GetStarted from './getting-started';
import HelpCenterWidget from './HelpCenterWidget';
import LatestBlogs from './LatestBlogs';
import MyDesignsWidget from './MyDesignsWidget';
import WorkspaceActivityWidget from './WorkspaceActivityWidget';
import KubernetesConnectionStatsChart from '../charts/KubernetesConnectionChart';
import React from 'react';

const getWidgets = ({ iconsProps, isEditMode }) => {
  const alwaysShown = () => true;

  const WIDGETS = {
    OVERVIEW: {
      title: 'Overview',
      isEnabled: alwaysShown,
      component: <Overview isEditMode={isEditMode} />,
      thumbnail: ClusterOverviewThumbnail?.src,
      defaultSizing: {
        w: 12,
        h: 2.565,
      },
    },
    GETTING_STARTED: {
      title: 'Getting Started',
      component: <GetStarted iconsProps={iconsProps} />,
      thumbnail: GettingStartedThumbnail?.src,
      isEnabled: alwaysShown,
      defaultSizing: {
        w: 3,
        h: 2,
      },
    },
    HELP_CENTER: {
      title: 'Help Center',
      isEnabled: alwaysShown,
      component: (
        <HelpCenterWidget href="" title="" description="Get help" iconsProps={iconsProps} />
      ),
      defaultSizing: {
        w: 3,
        h: 2,
      },
      thumbnail: HelpCenterThumbnail?.src,
    },
    MY_DESIGNS: {
      isEnabled: alwaysShown,
      title: 'My Designs',
      component: <MyDesignsWidget iconsProps={iconsProps} />,
      thumbnail: MyDesignthumbnail?.src,
      defaultSizing: {
        w: 3,
        h: 2,
      },
    },
    WORKSPACE_ACTIVITY: {
      title: 'Workspace Activity',
      isEnabled: alwaysShown,
      component: <WorkspaceActivityWidget />,
      thumbnail: WorkspaceActivityThumbnail?.src,
      defaultSizing: {
        w: 6,
        h: 2,
      },
    },
    CONNECTIONS_STATUS_CHART: {
      title: 'Connections Status Chart',
      isEnabled: alwaysShown,
      component: <KubernetesConnectionStatsChart />,
      thumbnail: ClusterStatusThumbnail?.src,
      defaultSizing: {
        w: 6,
        h: 2,
      },
    },
  };

  return WIDGETS;
};

export default getWidgets;
