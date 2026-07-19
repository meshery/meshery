import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { SistentThemeProviderWithoutBaseLine } from '@sistent/sistent';
import MesheryPatternsToolbar from './MesheryPatternsToolbar';

const baseProps = {
  width: 1200,
  isSearchExpanded: false,
  setIsSearchExpanded: vi.fn(),
  selectedPattern: { show: false },
  viewType: 'table',
  setViewType: vi.fn(),
  disableCreateImportDesignButton: false,
  disableUniversalFilter: false,
  pageTitle: 'Designs',
  router: { push: vi.fn() },
  handleUploadImport: vi.fn(),
  setSearch: vi.fn(),
  filter: {
    visibility: {
      name: 'Visibility',
      options: [
        { label: 'Public', value: 'public' },
        { label: 'Private', value: 'private' },
      ],
    },
  },
  selectedFilters: { visibility: 'All' },
  setSelectedFilters: vi.fn(),
  handleApplyFilter: vi.fn(),
  columns: [],
  columnVisibility: {},
  setColumnVisibility: vi.fn(),
};

const renderToolbar = (props = {}) =>
  render(
    <SistentThemeProviderWithoutBaseLine initialMode="light">
      <MesheryPatternsToolbar {...baseProps} {...props} />
    </SistentThemeProviderWithoutBaseLine>,
  );

describe('MesheryPatternsToolbar', () => {
  it('renders with real (unmocked) Sistent components without throwing', () => {
    renderToolbar();
  });

  it('renders the search bar, create/import buttons, and filter', () => {
    renderToolbar();
    expect(screen.getByTestId('meshery-patterns-search-bar')).toBeInTheDocument();
    expect(screen.getByTestId('meshery-patterns-create-design-btn')).toBeInTheDocument();
    expect(screen.getByTestId('meshery-patterns-import-design-btn')).toBeInTheDocument();
    expect(screen.getByTestId('meshery-patterns-universal-filter')).toBeInTheDocument();
  });

  it('renders the column-visibility control (not reachable by data-testid or title)', () => {
    // NOTE: CustomColumnVisibilityControl (sistent) doesn't declare or forward a
    // data-testid prop, so `data-testid="meshery-patterns-column-visibility-control"`
    // is silently dropped. Its icon button also has no title/aria-label reaching
    // the DOM at rest (MUI's Tooltip only exposes the label on hover/focus), so
    // it isn't reachable by getByTitle/getByRole(name:) either — only by a raw
    // container query against its (now-unique) id. This predates this migration:
    // the original ToolWrapper-based toolbar passed the same data-testid prop to
    // the same component. Flagging as a pre-existing a11y/testability gap in
    // sistent, not a regression.
    const { container } = renderToolbar();
    const columnVisibilityWrapper = container.querySelector('#designs-column-visibility-ref');
    expect(columnVisibilityWrapper).toBeInTheDocument();
    expect(columnVisibilityWrapper?.querySelector('button')).toBeInTheDocument();
    expect(screen.queryByTestId('meshery-patterns-column-visibility-control')).toBeNull();
  });

  it('hides the column-visibility control in grid view', () => {
    // Create/Import buttons are NOT gated by viewType (only by
    // selectedPattern.show / disableCreateImportDesignButton), so this test
    // only covers the column-visibility control, which is table-view-only.
    renderToolbar({ viewType: 'grid' });
    expect(screen.queryByTestId('meshery-patterns-column-visibility-control')).toBeNull();
    expect(document.querySelector('#designs-column-visibility-ref')).toBeNull();
  });

  it('hides the actions row on narrow viewports while search is expanded', () => {
    renderToolbar({ width: 400, isSearchExpanded: true });
    expect(screen.queryByTestId('meshery-patterns-create-design-btn')).toBeNull();
    // search bar itself should still be present so it can take the full width
    expect(screen.getByTestId('meshery-patterns-search-bar')).toBeInTheDocument();
  });
});
