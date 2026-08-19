import React from 'react';
import { render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
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

  // Sistent's DataTableToolbar (>= 0.22.3, see layer5io/sistent#1790) now
  // auto-compacts its own trailing slots (filter / column-visibility /
  // view-switch) once the *actual* browser viewport drops below the MUI `sm`
  // breakpoint, independent of the `width` prop this component uses to hide
  // the Create/Import buttons. These two mechanisms are separate and can
  // combine on a real narrow device, so we cover that combination here
  // rather than relying only on the mocked `width` prop above.
  describe('on a real narrow browser viewport (sistent auto-compaction)', () => {
    const originalInnerWidth = window.innerWidth;

    const setWindowWidth = (value: number) => {
      Object.defineProperty(window, 'innerWidth', {
        writable: true,
        configurable: true,
        value,
      });
      window.dispatchEvent(new Event('resize'));
    };

    afterEach(() => {
      setWindowWidth(originalInnerWidth);
    });

    it('auto-hides the universal filter once the real window is narrow, even with the actions row visible', async () => {
      setWindowWidth(375);
      renderToolbar({ width: 375, isSearchExpanded: false });
      // sistent syncs real dimensions in a mount effect; flush it
      await screen.findByTestId('meshery-patterns-create-design-btn');
      expect(screen.queryByTestId('meshery-patterns-universal-filter')).toBeNull();
    });

    it('keeps the universal filter visible on a wide real window', async () => {
      setWindowWidth(1200);
      renderToolbar({ width: 1200, isSearchExpanded: false });
      expect(await screen.findByTestId('meshery-patterns-universal-filter')).toBeInTheDocument();
    });
  });
});
