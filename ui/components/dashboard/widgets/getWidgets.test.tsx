import React from 'react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('./Thumbnails/WorkspaceActivity.svg', () => ({
  default: { src: 'workspace-activity.svg' },
}));
vi.mock('./Thumbnails/GettingStartedProgress.svg', () => ({
  default: { src: 'getting-started.svg' },
}));
vi.mock('./Thumbnails/ClusterOverview.png', () => ({
  default: { src: 'cluster-overview.png' },
}));
vi.mock('./Thumbnails/HelpCenter.svg', () => ({
  default: { src: 'help-center.svg' },
}));
vi.mock('./Thumbnails/MyDesigns.png', () => ({
  default: { src: 'my-designs.png' },
}));
vi.mock('./Thumbnails/ClusterStatus.png', () => ({
  default: { src: 'cluster-status.png' },
}));
vi.mock('./Thumbnails/LatestBlogs.png', () => ({
  default: { src: 'latest-blogs.png' },
}));

vi.mock('../overview', () => ({
  default: ({ isEditMode }: { isEditMode?: boolean }) => (
    <div data-testid="overview" data-edit={String(isEditMode)} />
  ),
}));

vi.mock('./getting-started', () => ({
  default: () => <div data-testid="getstarted" />,
}));

vi.mock('./HelpCenterWidget', () => ({
  default: () => <div data-testid="help-center" />,
}));

vi.mock('./LatestBlogs', () => ({
  default: () => <div data-testid="latest-blogs" />,
}));

vi.mock('./MyDesignsWidget', () => ({
  default: () => <div data-testid="my-designs" />,
}));

vi.mock('./WorkspaceActivityWidget', () => ({
  default: () => <div data-testid="workspace-activity" />,
}));

vi.mock('../charts/KubernetesConnectionChart', () => ({
  default: () => <div data-testid="k8s-conn" />,
}));

import getWidgets from './getWidgets';

describe('getWidgets', () => {
  it('returns the expected widget keys with thumbnails and defaultSizing', () => {
    const widgets = getWidgets({ iconsProps: {}, isEditMode: false });

    expect(Object.keys(widgets).sort()).toEqual(
      [
        'CONNECTIONS_STATUS_CHART',
        'GETTING_STARTED',
        'HELP_CENTER',
        'LATEST_BLOGS',
        'MY_DESIGNS',
        'OVERVIEW',
        'WORKSPACE_ACTIVITY',
      ].sort(),
    );

    expect(widgets.OVERVIEW.thumbnail).toBe('cluster-overview.png');
    expect(widgets.GETTING_STARTED.thumbnail).toBe('getting-started.svg');
    expect(widgets.HELP_CENTER.thumbnail).toBe('help-center.svg');
    expect(widgets.MY_DESIGNS.thumbnail).toBe('my-designs.png');
    expect(widgets.WORKSPACE_ACTIVITY.thumbnail).toBe('workspace-activity.svg');
    expect(widgets.CONNECTIONS_STATUS_CHART.thumbnail).toBe('cluster-status.png');
    expect(widgets.LATEST_BLOGS.thumbnail).toBe('latest-blogs.png');
  });

  it('marks every widget as enabled via the alwaysShown predicate', () => {
    const widgets = getWidgets({ iconsProps: {}, isEditMode: false });
    Object.values(widgets).forEach((w) => {
      expect(w.isEnabled()).toBe(true);
    });
  });

  it('exposes default sizing for each widget', () => {
    const widgets = getWidgets({ iconsProps: {}, isEditMode: true });
    expect(widgets.OVERVIEW.defaultSizing).toEqual({ w: 12, h: 2.565 });
    expect(widgets.GETTING_STARTED.defaultSizing).toEqual({ w: 3, h: 2 });
    expect(widgets.HELP_CENTER.defaultSizing).toEqual({ w: 3, h: 2 });
    expect(widgets.MY_DESIGNS.defaultSizing).toEqual({ w: 3, h: 2 });
    expect(widgets.WORKSPACE_ACTIVITY.defaultSizing).toEqual({ w: 6, h: 2 });
    expect(widgets.CONNECTIONS_STATUS_CHART.defaultSizing).toEqual({ w: 6, h: 2 });
    expect(widgets.LATEST_BLOGS.defaultSizing).toEqual({ w: 3, h: 2 });
  });

  it('forwards isEditMode into the Overview widget element', () => {
    const widgets = getWidgets({ iconsProps: {}, isEditMode: true });
    // The component itself is a JSX element; assert its props.
    expect((widgets.OVERVIEW.component as React.ReactElement).props.isEditMode).toBe(true);
  });

  it('the Help Center widget element receives "Get help" description and iconsProps', () => {
    const iconsProps = { width: 24 };
    const widgets = getWidgets({ iconsProps, isEditMode: false });
    const helpElement = widgets.HELP_CENTER.component as React.ReactElement<{
      description?: string;
      iconsProps?: { width?: number };
    }>;
    expect(helpElement.props.description).toBe('Get help');
    expect(helpElement.props.iconsProps).toBe(iconsProps);
  });
});
