import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

const resourcesTableSpy = vi.fn();

vi.mock('./resources-table', () => ({
  default: (props: { workloadType: string }) => {
    resourcesTableSpy(props);
    return <div data-testid="resources-table" data-kind={props.workloadType} />;
  },
}));

vi.mock('../tabpanel', () => ({
  TabPanel: ({
    children,
    value,
    index,
  }: {
    children?: React.ReactNode;
    value?: unknown;
    index?: unknown;
  }) => (value === index ? <div data-testid={`panel-${String(index)}`}>{children}</div> : null),
}));

vi.mock('../constants', () => ({
  TABS_SCROLL_BUTTONS_CLASS: 'scroll-buttons',
}));

vi.mock('../style', () => ({
  SecondaryTab: ({
    label,
    onClick,
    value,
  }: {
    label?: React.ReactNode;
    onClick?: React.MouseEventHandler;
    value?: number;
  }) => (
    <button type="button" data-testid={`tab-${value}`} onClick={onClick}>
      {label}
    </button>
  ),
  SecondaryTabs: ({
    children,
    onChange,
  }: {
    children?: React.ReactNode;
    onChange?: (event: unknown, value: number) => void;
  }) => (
    <div data-testid="tabs">
      <div data-onchange={!!onChange}>{children}</div>
      <button type="button" onClick={(e) => onChange?.(e, 1)}>
        change-to-1
      </button>
    </div>
  ),
  WrapperPaper: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="wrapper">{children}</div>
  ),
}));

vi.mock('../utils', () => ({
  default: ({ kind }: { kind: string }) => <svg data-testid={`icon-${kind}`} />,
}));

vi.mock('css/icons.styles', () => ({ iconMedium: {} }));

vi.mock('@sistent/sistent', () => ({
  styled: (tag: string | React.ComponentType<unknown>) => () => {
    const Wrapped = ({
      children,
      ...rest
    }: React.HTMLAttributes<HTMLElement> & { children?: React.ReactNode }) =>
      typeof tag === 'string'
        ? React.createElement(tag, rest as never, children)
        : React.createElement(tag as React.ComponentType<unknown>, rest as never, children);
    return Wrapped;
  },
}));

import ResourcesSubMenu from './resources-sub-menu';

describe('ResourcesSubMenu', () => {
  const baseResource = {
    submenu: true,
    tableConfig: () => ({
      Pod: { name: 'Pod' },
      Deployment: { name: 'Deployment' },
    }),
  };

  it('renders a tab per table-config key (non-CRDs path)', () => {
    const handleChange = vi.fn();
    render(
      <ResourcesSubMenu
        k8sConfig={null}
        resource={baseResource}
        selectedK8sContexts={null}
        selectedResource="Pod"
        handleChangeSelectedResource={handleChange}
      />,
    );
    expect(screen.getByTestId('tab-0')).toHaveTextContent('Pod');
    expect(screen.getByTestId('tab-1')).toHaveTextContent('Deployment');
    expect(screen.getByTestId('panel-Pod')).toBeInTheDocument();
  });

  it('invokes handleChangeSelectedResource when no selected resource is set (effect)', () => {
    const handleChange = vi.fn();
    render(
      <ResourcesSubMenu
        k8sConfig={null}
        resource={baseResource}
        selectedK8sContexts={null}
        handleChangeSelectedResource={handleChange}
      />,
    );
    expect(handleChange).toHaveBeenCalledWith('Pod');
  });

  it('invokes handleChangeSelectedResource with the next tab when the user changes tabs', async () => {
    const handleChange = vi.fn();
    render(
      <ResourcesSubMenu
        k8sConfig={null}
        resource={baseResource}
        selectedK8sContexts={null}
        selectedResource="Pod"
        handleChangeSelectedResource={handleChange}
      />,
    );
    await userEvent.click(screen.getByRole('button', { name: 'change-to-1' }));
    expect(handleChange).toHaveBeenCalledWith('Deployment');
  });

  it('renders the CRDs path when isCRDS is true', () => {
    const handleChange = vi.fn();
    render(
      <ResourcesSubMenu
        k8sConfig={null}
        resource={baseResource}
        selectedK8sContexts={null}
        selectedResource="Foo"
        handleChangeSelectedResource={handleChange}
        isCRDS={true}
        CRDsKeys={[
          { name: 'Foo', model: 'foo-model' },
          { name: 'Bar', model: 'bar-model' },
        ]}
      />,
    );
    expect(screen.getByTestId('tab-0')).toHaveTextContent('Foo');
    expect(screen.getByTestId('tab-1')).toHaveTextContent('Bar');
    // Panel for selected resource only
    expect(screen.getByTestId('panel-Foo')).toBeInTheDocument();
  });
});
