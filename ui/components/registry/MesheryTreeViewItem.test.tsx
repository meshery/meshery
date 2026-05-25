import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import MesheryTreeViewItem from './MesheryTreeViewItem';

vi.mock('../../constants/navigator', () => ({
  MODELS: 'Models',
}));

const normalizeStaticImagePath = vi.fn((s?: string) => s || null);
vi.mock('@/utils/fallback', () => ({
  normalizeStaticImagePath: (...args: any[]) => normalizeStaticImagePath(...args),
}));

vi.mock('./StyledTreeItem', () => ({
  default: ({ itemId, labelText, onClick, children }: any) => (
    <div data-testid={`tree-item-${itemId}`}>
      <button onClick={onClick} data-testid={`tree-button-${itemId}`}>
        {labelText}
      </button>
      {children}
    </div>
  ),
}));

vi.mock('./MeshModel.style', () => ({
  StyledTreeItemDiv: ({ children }: any) => <div>{children}</div>,
  StyledTreeItemNameDiv: ({ children }: any) => <span>{children}</span>,
}));

vi.mock('./VersionedModelComponentTree', () => ({
  default: ({ modelDef, versionedModelDef }: any) => (
    <div data-testid={`v-comp-${modelDef.id}-${versionedModelDef.id}`}>component tree</div>
  ),
}));

vi.mock('./VersionedModelRelationshipTree', () => ({
  default: ({ modelDef, versionedModelDef }: any) => (
    <div data-testid={`v-rel-${modelDef.id}-${versionedModelDef.id}`}>relationship tree</div>
  ),
}));

describe('MesheryTreeViewItem', () => {
  const baseProps = {
    setShowDetailsData: vi.fn(),
    showDetailsData: { type: '', data: {} },
    handleToggle: vi.fn(),
    handleSelect: vi.fn(),
    selected: [],
    expanded: [],
  };

  it('renders the model display name when present', () => {
    const modelDef = {
      id: 'm1',
      displayName: 'Kubernetes',
      name: 'kubernetes',
      metadata: { svgColor: '/img.svg' },
    };

    render(<MesheryTreeViewItem {...baseProps} modelDef={modelDef} />);

    expect(screen.getByText('Kubernetes')).toBeInTheDocument();
    expect(normalizeStaticImagePath).toHaveBeenCalledWith('/img.svg');
  });

  it('falls back to name when no displayName is set', () => {
    const modelDef = { id: 'm1', name: 'fallback-name', metadata: {} };

    render(<MesheryTreeViewItem {...baseProps} modelDef={modelDef} />);

    expect(screen.getByText('fallback-name')).toBeInTheDocument();
  });

  it('fires setShowDetailsData with model data when clicked', async () => {
    const setShowDetailsData = vi.fn();
    const user = userEvent.setup();
    const modelDef = { id: 'm1', displayName: 'M', metadata: {} };

    render(
      <MesheryTreeViewItem
        {...baseProps}
        modelDef={modelDef}
        setShowDetailsData={setShowDetailsData}
      />,
    );

    await user.click(screen.getByTestId('tree-button-m1'));

    expect(setShowDetailsData).toHaveBeenCalledWith({
      type: 'Models',
      data: modelDef,
    });
  });

  it('prefixes the registrant ID into the itemId', () => {
    const modelDef = { id: 'm1', displayName: 'M', metadata: {} };

    render(<MesheryTreeViewItem {...baseProps} modelDef={modelDef} registrantID="reg-1" />);

    expect(screen.getByTestId('tree-item-reg-1.1.m1')).toBeInTheDocument();
  });

  it('renders versionBasedData child tree items', () => {
    const modelDef = {
      id: 'm1',
      displayName: 'M',
      metadata: {},
      versionBasedData: [
        { id: 'v1', model: { version: 'v1.0' }, name: 'm' },
        { id: 'v2', model: { version: '2.0' }, name: 'm' },
      ],
    };

    render(<MesheryTreeViewItem {...baseProps} modelDef={modelDef} />);

    expect(screen.getByText('v1.0')).toBeInTheDocument();
    expect(screen.getByText('v2.0')).toBeInTheDocument();
    expect(screen.getByTestId('v-comp-m1-v1')).toBeInTheDocument();
    expect(screen.getByTestId('v-rel-m1-v2')).toBeInTheDocument();
  });

  it('renders versions without crashing when versionBasedData is missing', () => {
    const modelDef = { id: 'm1', displayName: 'M', metadata: {} };
    const { container } = render(<MesheryTreeViewItem {...baseProps} modelDef={modelDef} />);
    expect(container.querySelector('[data-testid="tree-item-m1"]')).toBeInTheDocument();
  });

  it('clicking a version child calls setShowDetailsData with the versioned model', async () => {
    const setShowDetailsData = vi.fn();
    const user = userEvent.setup();
    const modelDef = {
      id: 'm1',
      displayName: 'M',
      metadata: {},
      versionBasedData: [{ id: 'v1', model: { version: 'v1' }, name: 'm' }],
    };

    render(
      <MesheryTreeViewItem
        {...baseProps}
        modelDef={modelDef}
        setShowDetailsData={setShowDetailsData}
      />,
    );

    await user.click(screen.getByTestId('tree-button-m1.v1'));

    expect(setShowDetailsData).toHaveBeenCalledWith({
      type: 'Models',
      data: modelDef.versionBasedData[0],
    });
  });
});
