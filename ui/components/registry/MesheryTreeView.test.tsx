import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const useRegistryRouterMock = vi.fn();
const handleUpdateSelectedRoute = vi.fn();
let windowWidth = 1920;

vi.mock('@sistent/sistent', () => ({
  IconButton: ({ children, onClick, disabled, ...rest }: any) => (
    <button onClick={onClick} disabled={disabled} {...rest}>
      {children}
    </button>
  ),
  FormControlLabel: ({ control, label }: any) => (
    <label>
      {control}
      <span>{label}</span>
    </label>
  ),
  Switch: ({ checked, onClick, onChange, disabled, ...rest }: any) => (
    <input
      type="checkbox"
      data-testid="switch"
      checked={!!checked}
      onClick={onClick}
      disabled={disabled}
      onChange={onChange}
      {...rest}
    />
  ),
  CircularProgress: () => <div data-testid="loading-spinner" />,
  Typography: ({ children }: any) => <span>{children}</span>,
  InfoOutlinedIcon: () => <svg data-testid="info-icon" />,
}));

vi.mock('../../constants/navigator', () => ({
  MODELS: 'Models',
  COMPONENTS: 'Components',
  RELATIONSHIPS: 'Relationships',
  REGISTRANTS: 'Registrants',
  CONNECTIONS: 'Connections',
}));

vi.mock('@/utils/custom-search', () => ({
  default: ({ onSearch, expanded, setExpanded, value }: any) => (
    <div data-testid="search-bar">
      <input
        data-testid="search-input"
        value={value || ''}
        onChange={(e) => onSearch(e.target.value)}
      />
      <button data-testid="toggle-search" onClick={() => setExpanded(!expanded)} type="button">
        toggle:{String(expanded)}
      </button>
    </div>
  ),
}));

vi.mock('@/utils/debounce', () => ({
  default: (fn: any) => fn,
}));

vi.mock('@/utils/dimension', () => ({
  useWindowDimensions: () => ({ width: windowWidth }),
}));

vi.mock('./helper', () => ({
  getFilteredDataForDetailsComponent: () => ({ type: '', data: {} }),
}));

vi.mock('@/components/meshery-mesh-interface/PatternService/CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children }: any) => <div data-testid="ctt">{children}</div>,
}));

vi.mock('@/assets/icons/CollapseAll', () => ({
  default: () => <svg data-testid="collapse-all" />,
}));

vi.mock('@/assets/icons/ExpandAll', () => ({
  default: () => <svg data-testid="expand-all" />,
}));

vi.mock('@/theme', () => ({
  useTheme: () => ({
    palette: { primary: { main: 'primary' } },
  }),
}));

vi.mock('./MeshModel.style', () => ({
  JustifyAndAlignCenter: ({ children }: any) => <div data-testid="center-wrap">{children}</div>,
  MesheryTreeViewWrapper: ({ children }: any) => (
    <div data-testid="tree-view-wrapper">{children}</div>
  ),
}));

vi.mock('./hooks', () => ({
  useRegistryRouter: () => useRegistryRouterMock(),
}));

vi.mock('./MesheryTreeViewModel', () => ({
  default: (props: any) => <div data-testid="models-tree" data-len={props.data.length} />,
}));

vi.mock('./MesheryTreeViewRegistrants', () => ({
  default: (props: any) => <div data-testid="registrants-tree" data-len={props.data.length} />,
}));

vi.mock('./ComponentTree', () => ({
  default: (props: any) => <div data-testid="components-tree" data-len={props.data.length} />,
}));

vi.mock('./RelationshipTree', () => ({
  default: (props: any) => <div data-testid="relationships-tree" data-len={props.data.length} />,
}));

vi.mock('./ConnectionDefinitionTree', () => ({
  default: (props: any) => <div data-testid="connections-tree" data-len={props.data.length} />,
}));

import MesheryTreeView from './MesheryTreeView';

describe('MesheryTreeView', () => {
  const makeProps = (overrides: Partial<any> = {}) => ({
    data: [{ id: 'm1' }],
    view: 'Models',
    setSearchText: vi.fn(),
    searchText: null,
    setPage: vi.fn(),
    checked: false,
    setChecked: vi.fn(),
    setShowDetailsData: vi.fn(),
    showDetailsData: { type: '', data: {} },
    setResourcesDetail: vi.fn(),
    setModelsFilters: vi.fn(),
    lastItemRef: {
      Models: { current: null },
      Components: { current: null },
      Relationships: { current: null },
      Registrants: { current: null },
    } as any,
    isFetching: { Models: false, Components: false, Relationships: false, Registrants: false },
    isLoading: { Models: false, Components: false, Relationships: false, Registrants: false },
    ...overrides,
  });

  beforeEach(() => {
    windowWidth = 1920;
    handleUpdateSelectedRoute.mockReset();
    useRegistryRouterMock.mockReturnValue({
      handleUpdateSelectedRoute,
      selectedItemUUID: '',
    });
  });

  it('renders the models tree for view=Models', () => {
    render(<MesheryTreeView {...makeProps({ view: 'Models' })} />);
    expect(screen.getByTestId('models-tree')).toHaveAttribute('data-len', '1');
  });

  it('renders the registrants tree for view=Registrants', () => {
    render(<MesheryTreeView {...makeProps({ view: 'Registrants' })} />);
    expect(screen.getByTestId('registrants-tree')).toBeInTheDocument();
  });

  it('renders the components tree for view=Components', () => {
    render(<MesheryTreeView {...makeProps({ view: 'Components' })} />);
    expect(screen.getByTestId('components-tree')).toBeInTheDocument();
  });

  it('renders the relationships tree for view=Relationships', () => {
    render(<MesheryTreeView {...makeProps({ view: 'Relationships' })} />);
    expect(screen.getByTestId('relationships-tree')).toBeInTheDocument();
  });

  it('shows the loading spinner when data is empty and isLoading is true', () => {
    render(
      <MesheryTreeView
        {...makeProps({
          view: 'Models',
          data: [],
          isLoading: {
            Models: true,
            Components: false,
            Relationships: false,
            Registrants: false,
          },
        })}
      />,
    );
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
  });

  it('shows "No result found" when search returns empty', () => {
    render(
      <MesheryTreeView
        {...makeProps({
          view: 'Models',
          data: [],
          searchText: 'k8s',
        })}
      />,
    );
    expect(screen.getByText('No result found')).toBeInTheDocument();
  });

  it('disables expand/collapse buttons for COMPONENTS view', () => {
    render(<MesheryTreeView {...makeProps({ view: 'Components' })} />);
    const buttons = screen.getAllByRole('button');
    const expand = buttons.find((b) => b.querySelector('[data-testid="expand-all"]'));
    expect(expand).toBeDisabled();
  });

  it('routes via setSearchText when the search input changes', async () => {
    const setSearchText = vi.fn();
    render(<MesheryTreeView {...makeProps({ setSearchText })} />);

    fireEvent.change(screen.getByTestId('search-input'), {
      target: { value: 'kube' },
    });

    expect(setSearchText).toHaveBeenCalledWith('kube');
  });

  it('renders a Duplicates switch only on the MODELS view', () => {
    const { rerender } = render(<MesheryTreeView {...makeProps({ view: 'Models' })} />);
    expect(screen.getByTestId('switch')).toBeInTheDocument();

    rerender(<MesheryTreeView {...makeProps({ view: 'Relationships' })} />);
    expect(screen.queryByTestId('switch')).not.toBeInTheDocument();
  });

  it('disables the Duplicates switch when there are no records', () => {
    render(<MesheryTreeView {...makeProps({ view: 'Models', data: [], searchText: 'x' })} />);
    // Search returns no results — the switch should still render disabled
    // Note: with no records, "No result found" is shown but the header still renders
    expect(screen.getByTestId('switch')).toBeDisabled();
  });

  it('toggles the Duplicates switch via setChecked', async () => {
    const setChecked = vi.fn();
    const user = userEvent.setup();

    render(<MesheryTreeView {...makeProps({ setChecked, view: 'Models' })} />);

    await user.click(screen.getByTestId('switch'));
    expect(setChecked).toHaveBeenCalled();
  });

  describe('scroll position across infinite-scroll page loads', () => {
    it('restores the scroll offset after more records arrive', () => {
      const { container, rerender } = render(
        <MesheryTreeView {...makeProps({ view: 'Models', data: [{ id: 'm1' }] })} />,
      );

      const scroller = container.querySelector('.scrollElement') as HTMLDivElement;
      expect(scroller).not.toBeNull();

      fireEvent.scroll(scroller, { target: { scrollTop: 120 } });

      // Stand in for the browser dropping the offset when the list re-renders
      // with the next page of records.
      scroller.scrollTop = 0;

      rerender(
        <MesheryTreeView {...makeProps({ view: 'Models', data: [{ id: 'm1' }, { id: 'm2' }] })} />,
      );

      expect(scroller.scrollTop).toBe(120);
    });

    it('does not throw when the list empties and the scroll container unmounts', () => {
      const { container, rerender } = render(
        <MesheryTreeView {...makeProps({ view: 'Models', data: [{ id: 'm1' }] })} />,
      );

      const scroller = container.querySelector('.scrollElement') as HTMLDivElement;
      fireEvent.scroll(scroller, { target: { scrollTop: 120 } });

      // A search that matches nothing swaps the scrollable container out for the
      // "No result found" state, so the restore must tolerate its absence.
      expect(() =>
        rerender(
          <MesheryTreeView
            {...makeProps({ view: 'Models', data: [], searchText: 'no-such-model' })}
          />,
        ),
      ).not.toThrow();
      expect(screen.getByText('No result found')).toBeInTheDocument();
    });

    it('restores the offset of the view being returned to, not the last one scrolled', () => {
      // The reproduction from the issue: switching tabs unmounts and remounts
      // the scroll container. Each view keeps its own offset, so coming back to
      // Models must not land on wherever Relationships was left.
      const { container, rerender } = render(
        <MesheryTreeView {...makeProps({ view: 'Models', data: [{ id: 'm1' }] })} />,
      );

      const modelsScroller = container.querySelector('.scrollElement') as HTMLDivElement;
      fireEvent.scroll(modelsScroller, { target: { scrollTop: 120 } });

      // Switch to Relationships and scroll it somewhere else entirely.
      rerender(<MesheryTreeView {...makeProps({ view: 'Relationships', data: [{ id: 'r1' }] })} />);
      const relScroller = container.querySelector('.scrollElement') as HTMLDivElement;
      fireEvent.scroll(relScroller, { target: { scrollTop: 40 } });

      // Back to Models: the remounted container starts at 0 and must be
      // restored to Models' own 120, not Relationships' 40.
      rerender(<MesheryTreeView {...makeProps({ view: 'Models', data: [{ id: 'm1' }] })} />);
      const restored = container.querySelector('.scrollElement') as HTMLDivElement;

      expect(restored.scrollTop).toBe(120);
    });
    // The `!== undefined` comparison exists so that a saved offset of 0 is
    // restored rather than skipped. Every other test here saves a non-zero
    // offset, so a regression to a truthiness check would leave them green.
    it('restores a saved offset of 0', () => {
      const { container, rerender } = render(
        <MesheryTreeView {...makeProps({ view: 'Models', data: [{ id: 'm1' }] })} />,
      );

      const scroller = container.querySelector('.scrollElement') as HTMLDivElement;
      fireEvent.scroll(scroller, { target: { scrollTop: 120 } });
      fireEvent.scroll(scroller, { target: { scrollTop: 0 } });

      // Stand in for the browser leaving the list somewhere else before the
      // re-render; the saved 0 must win.
      scroller.scrollTop = 99;

      rerender(
        <MesheryTreeView {...makeProps({ view: 'Models', data: [{ id: 'm1' }, { id: 'm2' }] })} />,
      );

      expect(scroller.scrollTop).toBe(0);
    });

    // The container is not merely re-rendered here: the empty state replaces it
    // entirely, so the ref is detached and a new node is attached on the way
    // back. Distinct from the tab switch above, which swaps one view's container
    // for another's.
    it('restores the offset after the empty state detaches and reattaches the container', () => {
      const { container, rerender } = render(
        <MesheryTreeView {...makeProps({ view: 'Models', data: [{ id: 'm1' }] })} />,
      );

      const before = container.querySelector('.scrollElement') as HTMLDivElement;
      fireEvent.scroll(before, { target: { scrollTop: 120 } });

      // A search matching nothing swaps the list out for "No result found".
      rerender(
        <MesheryTreeView
          {...makeProps({ view: 'Models', data: [], searchText: 'no-such-model' })}
        />,
      );
      expect(container.querySelector('.scrollElement')).toBeNull();

      // Clearing the search brings the list, and a fresh container, back.
      rerender(<MesheryTreeView {...makeProps({ view: 'Models', data: [{ id: 'm1' }] })} />);
      const after = container.querySelector('.scrollElement') as HTMLDivElement;

      expect(after).not.toBe(before);
      expect(after.scrollTop).toBe(120);
    });
  });
});
