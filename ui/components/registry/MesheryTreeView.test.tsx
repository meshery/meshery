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
  Switch: ({ checked, onClick, disabled, ...rest }: any) => (
    <input
      type="checkbox"
      data-testid="switch"
      checked={!!checked}
      onClick={onClick}
      disabled={disabled}
      onChange={() => {}}
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
});
