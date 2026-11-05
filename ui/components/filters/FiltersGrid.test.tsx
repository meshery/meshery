import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const can = vi.fn(() => true);

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => can(...args),
}));

vi.mock('@sistent/sistent', () => ({
  Grid2: ({ children }: any) => <div data-testid="grid">{children}</div>,
  Button: ({ children, onClick, disabled }: any) => (
    <button onClick={onClick} disabled={disabled}>
      {children}
    </button>
  ),
  Pagination: ({ count, page, onChange }: any) => (
    <div data-testid="pagination">
      <button data-testid="pagination-next" onClick={(e) => onChange(e, page + 1)}>
        {`${page}/${count}`}
      </button>
    </div>
  ),
  Modal: ({ children, open, title, closeModal }: any) =>
    open ? (
      <div data-testid="modal" data-title={title}>
        {children}
        <button data-testid="close" onClick={closeModal} />
      </div>
    ) : null,
}));

vi.mock('./FiltersCard', () => ({
  default: ({
    name,
    handleClone,
    handleDownload,
    deleteHandler,
    updateHandler,
    setSelectedFilters,
    handlePublishModal,
    handleUnpublishModal,
    handleInfoModal,
    setYaml,
  }: any) => (
    <div data-testid={`filter-card-${name}`}>
      <span>{name}</span>
      <button data-testid={`clone-${name}`} onClick={handleClone}>
        clone
      </button>
      <button data-testid={`download-${name}`} onClick={(e: any) => handleDownload(e)}>
        download
      </button>
      <button data-testid={`delete-${name}`} onClick={deleteHandler}>
        delete
      </button>
      <button data-testid={`update-${name}`} onClick={updateHandler}>
        update
      </button>
      <button data-testid={`info-${name}`} onClick={handleInfoModal}>
        info
      </button>
      <button data-testid={`publish-${name}`} onClick={handlePublishModal}>
        publish
      </button>
      <button data-testid={`unpublish-${name}`} onClick={handleUnpublishModal}>
        unpublish
      </button>
      <button data-testid={`select-${name}`} onClick={setSelectedFilters}>
        select
      </button>
      <button data-testid={`yaml-${name}`} onClick={() => setYaml('foo: 1')}>
        yaml
      </button>
    </div>
  ),
}));

vi.mock('../../utils/Enum', () => ({
  FILE_OPS: { DELETE: 'DELETE', UPDATE: 'UPDATE' },
}));

vi.mock('../designs/patterns/Grid.styles', () => ({
  GridAddIconStyles: () => <svg data-testid="grid-add-icon" />,
  GridNoContainerStyles: ({ children }: any) => <div>{children}</div>,
  GridNoPapperStyles: ({ children }: any) => <div data-testid="empty-paper">{children}</div>,
  GridNoTextStyles: ({ children }: any) => <div data-testid="no-text">{children}</div>,
  GridPaginationStyles: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('../shared/Modal/Modal', () => ({
  RJSFModalWrapper: ({ submitBtnText, handleSubmit, handleClose }: any) => (
    <div data-testid="rjsf-wrapper">
      <button data-testid="submit-rjsf" onClick={() => handleSubmit({ foo: 'bar' })}>
        {submitBtnText}
      </button>
      <button data-testid="cancel-rjsf" onClick={handleClose} />
    </div>
  ),
}));

vi.mock('../../public/static/img/drawer-icons/filter_svg', () => ({
  default: () => <svg data-testid="filter-icon" />,
}));

import FiltersGrid from './FiltersGrid';

const filtersData = [
  {
    id: '1',
    name: 'one',
    filter_resource: JSON.stringify({ settings: { config: 'a: 1' } }),
  },
  {
    id: '2',
    name: 'two',
    filter_resource: '',
  },
];

describe('FiltersGrid', () => {
  let setSelectedFilter: ReturnType<typeof vi.fn>;
  let setPublishModal: ReturnType<typeof vi.fn>;
  let setPage: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    can.mockReset();
    can.mockReturnValue(true);
    setSelectedFilter = vi.fn();
    setPublishModal = vi.fn();
    setPage = vi.fn();
  });

  const renderComponent = (overrides: any = {}) =>
    render(
      <FiltersGrid
        filters={filtersData}
        handleClone={vi.fn()}
        handleDownload={vi.fn()}
        handleSubmit={vi.fn()}
        setSelectedFilter={setSelectedFilter}
        selectedFilter={{ show: false, filter: null }}
        pages={3}
        setPage={setPage}
        selectedPage={1}
        canPublishFilter={true}
        handleUploadImport={vi.fn()}
        handlePublish={vi.fn()}
        handleUnpublishModal={() => () => Promise.resolve()}
        publishModal={{ open: false, filter: {}, name: '' }}
        setPublishModal={setPublishModal}
        publishSchema={{ rjsfSchema: {}, uiSchema: {} }}
        handleInfoModal={vi.fn()}
        {...overrides}
      />,
    );

  it('renders one card per filter in the array', () => {
    renderComponent();
    expect(screen.getByTestId('filter-card-one')).toBeInTheDocument();
    expect(screen.getByTestId('filter-card-two')).toBeInTheDocument();
  });

  it('hides the filter cards when selectedFilter.show is true', () => {
    renderComponent({ selectedFilter: { show: true, filter: { id: 1 } } });
    expect(screen.queryByTestId('filter-card-one')).not.toBeInTheDocument();
  });

  it('shows the empty state when filters is empty and no selected filter', () => {
    renderComponent({ filters: [] });
    expect(screen.getByTestId('empty-paper')).toBeInTheDocument();
    expect(screen.getByTestId('no-text')).toHaveTextContent('No Filters Found');
  });

  it('renders an "Import Filter" button in the empty state and calls handleUploadImport', async () => {
    const user = userEvent.setup();
    const handleUploadImport = vi.fn();
    renderComponent({ filters: [], handleUploadImport });
    await user.click(screen.getByRole('button', { name: /import filter/i }));
    expect(handleUploadImport).toHaveBeenCalledTimes(1);
  });

  it('renders pagination when there are filters', () => {
    renderComponent();
    expect(screen.getByTestId('pagination')).toBeInTheDocument();
  });

  it('calls setPage when pagination changes', async () => {
    const user = userEvent.setup();
    renderComponent();
    await user.click(screen.getByTestId('pagination-next'));
    // page-1 should be set (1-indexed becomes 0-indexed)
    expect(setPage).toHaveBeenCalled();
  });

  it('opens publish modal when handlePublishModal is invoked and canPublishFilter is true', async () => {
    const user = userEvent.setup();
    renderComponent();
    await user.click(screen.getByTestId('publish-one'));
    expect(setPublishModal).toHaveBeenCalledWith({
      open: true,
      filter: filtersData[0],
      name: '',
    });
  });

  it('does NOT open publish modal when canPublishFilter is false', async () => {
    const user = userEvent.setup();
    renderComponent({ canPublishFilter: false });
    await user.click(screen.getByTestId('publish-one'));
    expect(setPublishModal).not.toHaveBeenCalled();
  });

  it('renders the publish RJSF modal when publishModal.open is true', () => {
    renderComponent({
      publishModal: { open: true, filter: { name: 'one' }, name: '' },
    });
    expect(screen.getByTestId('modal')).toBeInTheDocument();
  });

  it('closes the publish modal when the close button is clicked', async () => {
    const user = userEvent.setup();
    renderComponent({
      publishModal: { open: true, filter: { name: 'one' }, name: '' },
    });
    await user.click(screen.getByTestId('cancel-rjsf'));
    expect(setPublishModal).toHaveBeenCalledWith({
      open: false,
      filter: {},
      name: '',
    });
  });
});
