import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const can = vi.fn(() => true);

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => can(...args),
}));

vi.mock('@/assets/icons', () => ({
  GetApp: ({ ...props }: any) => <svg data-testid="download-icon" {...props} />,
  Public: ({ ...props }: any) => <svg data-testid="public-icon" {...props} />,
}));

vi.mock('@sistent/sistent', () => ({
  InfoOutlinedIcon: ({ ...props }: any) => <svg data-testid="info-icon" {...props} />,
  EditIcon: ({ ...props }: any) => <svg data-testid="edit-icon" {...props} />,
}));

vi.mock('react-moment', () => ({
  default: ({ children }: any) => <span data-testid="moment">{String(children)}</span>,
}));

vi.mock('../../public/static/img/CloneIcon', () => ({
  default: ({ ...props }: any) => <svg data-testid="clone-icon" {...props} />,
}));

vi.mock('../../css/icons.styles', () => ({
  iconMedium: {},
}));

vi.mock('../../utils/Enum', () => ({
  VISIBILITY: { PUBLISHED: 'published', PUBLIC: 'public', PRIVATE: 'private' },
}));

vi.mock('../connections/common/index', () => ({
  DefaultTableCell: ({ columnData }: any) => <th>{columnData?.label}</th>,
  SortableTableCell: ({ columnData, onSort }: any) => (
    <th>
      <button data-testid={`sort-${columnData?.label}`} onClick={onSort}>
        {columnData?.label}
      </button>
    </th>
  ),
}));

vi.mock('./TooltipIcon', () => ({
  default: ({ children, onClick, title, disabled }: any) => (
    <button data-testid={`action-${title}`} onClick={onClick} disabled={disabled}>
      {children}
    </button>
  ),
}));

vi.mock('./Filters.styled', () => ({
  ActionsBox: ({ children }: any) => <div data-testid="actions-box">{children}</div>,
}));

import { buildFiltersColumns } from './Filters.columns';

const handleClone = vi.fn();
const handleDownload = vi.fn();
const handleInfoModal = vi.fn();
const handlePublishModal = vi.fn();
const handleUnpublishModal = vi.fn();
const setSelectedRowData = vi.fn();

const filters = [
  { id: 'a', name: 'first', visibility: 'public' },
  { id: 'b', name: 'second', visibility: 'published' },
];

describe('buildFiltersColumns', () => {
  beforeEach(() => {
    handleClone.mockReset();
    handleDownload.mockReset();
    handleInfoModal.mockReset();
    handlePublishModal.mockReset();
    handleUnpublishModal.mockReset();
    setSelectedRowData.mockReset();
    can.mockReset();
    can.mockReturnValue(true);
  });

  const build = (overrides: any = {}) =>
    buildFiltersColumns({
      filters,
      canPublishFilter: true,
      handleClone,
      handleDownload,
      handleInfoModal,
      handlePublishModal,
      handleUnpublishModal: () => () => Promise.resolve(),
      setSelectedRowData,
      sortOrder: 'updated_at desc',
      ...overrides,
    });

  it('returns the canonical 5-column schema', () => {
    const columns = build();
    const names = columns.map((c) => c.name);
    expect(names).toEqual(['name', 'created_at', 'updated_at', 'visibility', 'Actions']);
  });

  it('applies sortDirection to the column matching sortOrder', () => {
    const columns = build({ sortOrder: 'name asc' });
    const nameCol = columns.find((c) => c.name === 'name')!;
    expect(nameCol.options.sortDirection).toBe('asc');
  });

  it('renders moment for created_at and updated_at body cells', () => {
    const columns = build();
    const createdCol = columns.find((c) => c.name === 'created_at')!;
    const { getByTestId, rerender } = render(createdCol.options.customBodyRender('2026-01-01'));
    expect(getByTestId('moment')).toHaveTextContent('2026-01-01');

    const updatedCol = columns.find((c) => c.name === 'updated_at')!;
    rerender(updatedCol.options.customBodyRender('2026-01-02'));
    expect(getByTestId('moment')).toHaveTextContent('2026-01-02');
  });

  it('renders an Edit button for non-published filters and triggers setSelectedRowData', async () => {
    const user = userEvent.setup();
    const columns = build();
    const actionsCol = columns.find((c) => c.name === 'Actions')!;
    render(actionsCol.options.customBodyRender(null, { rowIndex: 0 }));

    // visibility !== published -> Config button (Edit)
    const editBtn = screen.getByTestId('action-Config');
    await user.click(editBtn);
    expect(setSelectedRowData).toHaveBeenCalledWith(filters[0]);
  });

  it('renders a Clone button for published filters and triggers handleClone', async () => {
    const user = userEvent.setup();
    const columns = build();
    const actionsCol = columns.find((c) => c.name === 'Actions')!;
    render(actionsCol.options.customBodyRender(null, { rowIndex: 1 }));

    await user.click(screen.getByTestId('action-Clone'));
    expect(handleClone).toHaveBeenCalledWith('b', 'second');
  });

  it('invokes handleDownload with the row id and name', async () => {
    const user = userEvent.setup();
    const columns = build();
    const actionsCol = columns.find((c) => c.name === 'Actions')!;
    render(actionsCol.options.customBodyRender(null, { rowIndex: 0 }));

    await user.click(screen.getByTestId('action-Download'));
    expect(handleDownload).toHaveBeenCalledWith(expect.anything(), 'a', 'first');
  });

  it('invokes handleInfoModal with the row data', async () => {
    const user = userEvent.setup();
    const columns = build();
    const actionsCol = columns.find((c) => c.name === 'Actions')!;
    render(actionsCol.options.customBodyRender(null, { rowIndex: 0 }));

    await user.click(screen.getByTestId('action-Filter Information'));
    expect(handleInfoModal).toHaveBeenCalledWith(filters[0]);
  });

  it('shows Publish action for non-published filters when canPublishFilter is true', () => {
    const columns = build();
    const actionsCol = columns.find((c) => c.name === 'Actions')!;
    render(actionsCol.options.customBodyRender(null, { rowIndex: 0 }));
    expect(screen.getByTestId('action-Publish')).toBeInTheDocument();
    expect(screen.queryByTestId('action-Unpublish')).not.toBeInTheDocument();
  });

  it('shows Unpublish action for published filters', () => {
    const columns = build();
    const actionsCol = columns.find((c) => c.name === 'Actions')!;
    render(actionsCol.options.customBodyRender(null, { rowIndex: 1 }));
    expect(screen.getByTestId('action-Unpublish')).toBeInTheDocument();
  });

  it('renders Sortable headers for sortable columns', () => {
    const columns = build();
    const nameCol = columns.find((c) => c.name === 'name')!;
    render(nameCol.options.customHeadRender({ index: 0, label: 'Filter Name' }, () => {}, {}));
    expect(screen.getByTestId('sort-Filter Name')).toBeInTheDocument();
  });
});
