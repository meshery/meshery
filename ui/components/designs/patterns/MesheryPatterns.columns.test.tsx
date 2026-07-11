import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const canMock = vi.fn(() => true);

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: any) => <div>{children}</div>,
  crimson: { 40: '#F91313' },
  InfoOutlinedIcon: () => <svg data-testid="info-outlined" />,
  AccountTreeIcon: () => <svg data-testid="account-tree" />,
  EditIcon: () => <svg data-testid="edit" />,
}));

vi.mock('react-moment', () => ({
  default: ({ children }: any) => <span data-testid="moment">{children}</span>,
}));

vi.mock('@/assets/icons', () => ({
  GetApp: () => <svg data-testid="get-app" />,
  DoneAll: () => <svg data-testid="done-all" />,
  Public: () => <svg data-testid="public" />,
}));

vi.mock('../../../public/static/img/UndeployIcon', () => ({
  default: () => <svg data-testid="undeploy" />,
}));

vi.mock('../../../public/static/img/CloneIcon', () => ({
  default: () => <svg data-testid="clone" />,
}));

vi.mock('../../connections/common', () => ({
  DefaultTableCell: ({ columnData }: any) => (
    <div data-testid="default-cell">{columnData?.label}</div>
  ),
  SortableTableCell: ({ columnData, onSort }: any) => (
    <div data-testid="sortable-cell" onClick={onSort}>
      {columnData?.label}
    </div>
  ),
}));

vi.mock('@/utils/can', () => ({
  default: (...args: any[]) => canMock(...args),
}));

vi.mock('@/assets/icons/CheckIcon', () => ({
  default: () => <svg data-testid="check-icon" />,
}));

vi.mock('@/assets/icons/DryRunIcon', () => ({
  default: () => <svg data-testid="dry-run-icon" />,
}));

vi.mock('@/assets/icons/PatternConfigure', () => ({
  default: () => <svg data-testid="pattern-configure" />,
}));

vi.mock('./ActionPopover', () => ({
  default: ({ actions }: any) => (
    <div
      data-testid="action-popover"
      data-actions={JSON.stringify(actions.map((a: any) => a.label))}
    />
  ),
}));

vi.mock('./CustomToolbarSelect', () => ({
  default: () => <div data-testid="toolbar-select" />,
}));

vi.mock('../../../utils/Enum', () => ({
  VISIBILITY: { PUBLISHED: 'published', PRIVATE: 'private', PUBLIC: 'public' },
}));

vi.mock('./MesheryPatterns.constants', () => ({
  genericClickHandler: (ev: any, fn: any) => {
    ev.stopPropagation?.();
    fn(ev);
  },
}));

import {
  buildPatternActions,
  buildPatternColumns,
  buildPatternsTableOptions,
  PATTERN_COL_VIEWS,
} from './MesheryPatterns.columns';

const makeHandlers = () => ({
  handleOpenInConfigurator: vi.fn(),
  handleClone: vi.fn(),
  openValidateModal: vi.fn(),
  openDryRunModal: vi.fn(),
  openUndeployModal: vi.fn(),
  openDeployModal: vi.fn(),
  handleDesignDownloadModal: vi.fn(),
  handleInfoModal: vi.fn(),
  handleUnpublishModal: vi.fn(),
  handleEvaluateRelationship: vi.fn(),
  userCanEdit: (_: any) => true,
});

describe('PATTERN_COL_VIEWS', () => {
  it('declares the responsive column view config', () => {
    expect(PATTERN_COL_VIEWS).toEqual([
      ['name', 'xs'],
      ['created_at', 'm'],
      ['updated_at', 'm'],
      ['visibility', 's'],
      ['Actions', 'xs'],
    ]);
  });
});

describe('buildPatternActions', () => {
  it('returns the canonical action labels for a private design', () => {
    const handlers = makeHandlers();
    const actions = buildPatternActions({
      rowData: { id: 'p1', name: 'p', patternFile: '' },
      visibility: 'private',
      patterns: [{ id: 'p1' }],
      tableMeta: { rowIndex: 0 },
      handlers,
    });
    const labels = actions.map((a: any) => a.label);
    expect(labels).toContain('Edit');
    expect(labels).toContain('Validate Design');
    expect(labels).toContain('Deploy');
    expect(labels).toContain('Download');
    expect(labels).not.toContain('Unpublish');
  });

  it('exposes Unpublish only for published designs', () => {
    const actions = buildPatternActions({
      rowData: { id: 'p1', name: 'p', patternFile: '' },
      visibility: 'published',
      patterns: [{ id: 'p1' }],
      tableMeta: { rowIndex: 0 },
      handlers: makeHandlers(),
    });
    expect(actions.map((a: any) => a.label)).toContain('Unpublish');
  });

  it('routes the Edit action through handleOpenInConfigurator', () => {
    const handlers = makeHandlers();
    const actions = buildPatternActions({
      rowData: { id: 'p1', name: 'p', patternFile: '' },
      visibility: 'private',
      patterns: [{ id: 'p1' }],
      tableMeta: { rowIndex: 0 },
      handlers,
    });
    const edit = actions.find((a: any) => a.label === 'Edit');
    expect(edit).toBeDefined();
    edit!.onClick({ stopPropagation: vi.fn() } as any);
    expect(handlers.handleOpenInConfigurator).toHaveBeenCalled();
  });
});

describe('buildPatternColumns', () => {
  it('returns columns and renders an ActionPopover for the actions column', () => {
    const patterns = [{ id: 'p1', name: 'p', visibility: 'private', patternFile: '' }];
    const columns = buildPatternColumns({ patterns, handlers: makeHandlers() });

    expect(columns.map((c: any) => c.name)).toEqual([
      'name',
      'created_at',
      'updated_at',
      'visibility',
      'Actions',
    ]);

    const actionsColumn = columns.find((c: any) => c.name === 'Actions');
    render(actionsColumn!.options.customBodyRender(null, { rowIndex: 0 }));
    expect(screen.getByTestId('action-popover')).toBeInTheDocument();
  });

  it('renders a Moment-formatted cell for created_at and updated_at', () => {
    const patterns = [{ id: 'p1', name: 'p', visibility: 'private', patternFile: '' }];
    const columns = buildPatternColumns({ patterns, handlers: makeHandlers() });
    const createdAt = columns.find((c: any) => c.name === 'created_at');
    render(createdAt!.options.customBodyRender('2024-01-02'));
    expect(screen.getByTestId('moment')).toHaveTextContent('2024-01-02');
  });
});

describe('buildPatternsTableOptions', () => {
  const build = (overrides: any = {}) =>
    buildPatternsTableOptions({
      patterns: [],
      columns: [],
      count: 42,
      pageSize: 10,
      page: 3,
      search: '',
      sortOrder: 'name asc',
      isLocalProvider: false,
      searchTimeout: { current: null },
      setPage: vi.fn(),
      setPageSize: vi.fn(),
      setSearch: vi.fn(),
      setSortOrder: vi.fn(),
      setSelectedRowData: vi.fn(),
      deletePatterns: vi.fn(),
      showModal: vi.fn(),
      initPatternsSubscription: vi.fn(),
      ...overrides,
    });

  it('returns a config preserving the supplied page, page size and counts', () => {
    const options = build();

    expect(options.count).toBe(42);
    expect(options.rowsPerPage).toBe(10);
    expect(options.page).toBe(3);
    expect(options.sortOrder).toEqual({ name: 'name', direction: 'asc' });
    expect(options.print).toBe(false);
    expect(options.download).toBe(false);
  });

  it('disables sort on the local provider and enables it on a remote provider', () => {
    expect(build({ isLocalProvider: true }).sort).toBe(false);
    expect(build().sort).toBe(true);
  });
});
