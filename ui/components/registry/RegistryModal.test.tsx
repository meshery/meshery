import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { Keys } from '@meshery/schemas/permissions';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const {
  mockUseHasPermission,
  mockUseGetMeshModelsQuery,
  mockUseGetComponentsQuery,
  mockUseGetRelationshipsQuery,
  mockUseGetRegistrantsQuery,
} = vi.hoisted(() => ({
  mockUseHasPermission: vi.fn(),
  mockUseGetMeshModelsQuery: vi.fn(),
  mockUseGetComponentsQuery: vi.fn(),
  mockUseGetRelationshipsQuery: vi.fn(),
  mockUseGetRegistrantsQuery: vi.fn(),
}));

let mediaQueryReturn = false;
const RegistryModalContext = React.createContext<any>({});

// Stable theme reference so setHeaderInfo effects don't loop.
function createStableTheme() {
  return {
    palette: {
      icon: { default: 'i', white: 'w' },
      common: { white: 'w' },
      mode: 'light',
      background: { paper: '#fff' },
    },
    transitions: {
      create: () => '',
      easing: { sharp: '' },
      duration: { enteringScreen: 0, leavingScreen: 0 },
    },
    breakpoints: {
      up: () => '@media (min-width:600px)',
      down: () => '@media (max-width:600px)',
    },
    spacing: (n: number) => n,
  };
}

const stableTheme = createStableTheme();

vi.mock('@sistent/sistent', () => ({
  ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
  Box: ({ children, sx }: any) => (
    <div data-testid="box" data-sx={JSON.stringify(sx || {})}>
      {children}
    </div>
  ),
  List: ({ children }: any) => <ul>{children}</ul>,
  ListItem: ({ children }: any) => <li>{children}</li>,
  ListItemButton: ({ children, onClick, selected }: any) => (
    <button onClick={onClick} data-selected={String(!!selected)}>
      {children}
    </button>
  ),
  ListItemIcon: ({ children }: any) => <span>{children}</span>,
  ListItemText: ({ primary }: any) => <span>{primary}</span>,
  IconButton: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
  ChevronLeftIcon: (props) =>
    React.createElement('svg', { 'data-component': 'chevron-left-icon', ...props }),
  ChevronRightIcon: (props) =>
    React.createElement('svg', { 'data-component': 'chevron-right-icon', ...props }),
  LeftArrowIcon: (props) =>
    React.createElement('svg', { 'data-component': 'left-arrow-icon', ...props }),
  FileIcon: () => <svg data-testid="file-icon" />,
  AppRegistrationIcon: () => <svg data-testid="app-registration-icon" />,
  InfoIcon: () => <svg data-testid="info-icon" />,
  DatabaseIcon: () => <svg data-testid="database-icon" />,
  CustomTooltip: ({ children }: any) => <>{children}</>,
  useHasPermission: (permissionKey: any) => mockUseHasPermission(permissionKey),
  useMediaQuery: () => mediaQueryReturn,
  useTheme: () => stableTheme,
  alpha: (value: string) => value,
  lighten: (value: string) => value,
  ThemeProvider: ({ children }: any) => <>{children}</>,
  ErrorBoundary: ({ children }: any) => <>{children}</>,
  CssBaseline: () => null,
  NoSsr: ({ children }: any) => <>{children}</>,
  createTheme: (...overrides: any[]) =>
    overrides.reduce(
      (theme, override) => ({
        ...theme,
        ...override,
        palette: { ...theme.palette, ...(override?.palette || {}) },
        transitions: { ...theme.transitions, ...(override?.transitions || {}) },
        breakpoints: { ...theme.breakpoints, ...(override?.breakpoints || {}) },
        spacing: override?.spacing || theme.spacing,
      }),
      createStableTheme(),
    ),
  Modal: ({ children, open, title }: any) =>
    open ? (
      <div data-testid="registry-modal" data-title={title}>
        {children}
      </div>
    ) : null,
  Drawer: ({ children }: any) => <div>{children}</div>,
  DARK_BLUE_GRAY: '#000',
  styled: (Component: any) => () => {
    const StyledComponent = ({ children, ...rest }: any) =>
      React.createElement(Component, rest, children);
    return StyledComponent;
  },
}));

vi.mock('@/assets/icons/Connection', () => ({ default: () => <svg /> }));
vi.mock('@/assets/icons/Component', () => ({ default: () => <svg /> }));
vi.mock('@/components/shared/Modal', () => ({
  Modal: ({ children, isOpen, onClose, title }: any) =>
    isOpen ? (
      <div data-testid="registry-modal" data-title={title}>
        <button type="button" data-testid="close-registry-modal" onClick={onClose}>
          Close
        </button>
        {children}
      </div>
    ) : null,
}));

vi.mock('../general/error-404/index', () => ({
  default: ({ permissionKey }: any) => (
    <div data-testid="default-error" data-permission-key={JSON.stringify(permissionKey)} />
  ),
}));

vi.mock('./MeshModelComponent', () => ({
  default: (props: any) => (
    <div
      data-testid="mesh-model-component"
      data-view={props.externalView}
      data-search={props.externalSearchText}
      data-uuid={props.externalSelectedItemUUID}
    />
  ),
}));

vi.mock('css/icons.styles', () => ({
  iconMedium: {},
  iconSmall: {},
}));

vi.mock('../../constants/navigator', () => ({
  MODELS: 'Models',
  COMPONENTS: 'Components',
  RELATIONSHIPS: 'Relationships',
  REGISTRANTS: 'Registrants',
}));

vi.mock('@/utils/context/RegistryModalContextProvider', () => ({
  get RegistryModalContext() {
    return RegistryModalContext;
  },
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useGetMeshModelsQuery: (...args: any[]) => mockUseGetMeshModelsQuery(...args),
  useGetComponentsQuery: (...args: any[]) => mockUseGetComponentsQuery(...args),
  useGetRelationshipsQuery: (...args: any[]) => mockUseGetRelationshipsQuery(...args),
  useGetRegistrantsQuery: (...args: any[]) => mockUseGetRegistrantsQuery(...args),
}));

vi.mock('./helper', () => ({
  removeDuplicateVersions: (m: any[]) => m,
}));

vi.mock('../general/style', () => ({
  ChevronButtonWrapper: ({ children, onClick, isCollapsed, sx, ...props }: any) => (
    <button
      type="button"
      onClick={onClick}
      data-collapsed={String(!!isCollapsed)}
      data-testid="sidebar-collapse-toggle"
      {...props}
    >
      {children}
    </button>
  ),
}));

import RegistryModal, { Navigation } from './RegistryModal';

const wrap = (ui: React.ReactElement, ctx: any) => (
  <RegistryModalContext.Provider value={ctx}>{ui}</RegistryModalContext.Provider>
);

describe('RegistryModal', () => {
  const baseCtx = () => ({
    open: false,
    closeModal: vi.fn(),
    selectedView: 'Models',
    searchText: '',
    selectedItemUUID: '',
    setSelectedView: vi.fn(),
  });

  beforeEach(() => {
    mediaQueryReturn = false;
    mockUseHasPermission.mockReset();
    mockUseHasPermission.mockReturnValue(true);
    mockUseGetMeshModelsQuery.mockReset();
    mockUseGetMeshModelsQuery.mockReturnValue({
      data: { models: [] },
      isLoading: false,
    });
    mockUseGetComponentsQuery.mockReset();
    mockUseGetComponentsQuery.mockReturnValue({
      data: { totalCount: 7 },
      isLoading: false,
    });
    mockUseGetRelationshipsQuery.mockReset();
    mockUseGetRelationshipsQuery.mockReturnValue({
      data: { totalCount: 3 },
      isLoading: false,
    });
    mockUseGetRegistrantsQuery.mockReset();
    mockUseGetRegistrantsQuery.mockReturnValue({
      data: { totalCount: 5 },
      isLoading: false,
    });
  });

  it('renders nothing when context.open is false', () => {
    render(wrap(<RegistryModal />, baseCtx()));
    expect(screen.queryByTestId('registry-modal')).not.toBeInTheDocument();
  });

  it('renders the modal when open is true and the inner Navigation drives the title', () => {
    const ctx = { ...baseCtx(), open: true };
    render(wrap(<RegistryModal />, ctx));
    const modal = screen.getByTestId('registry-modal');
    // Navigation calls setHeaderInfo on mount, which updates the modal title
    // to "Registry - <selectedView>". Accept either the default ("Registry")
    // or the updated form.
    expect(modal.getAttribute('data-title')).toMatch(/^Registry/);
  });

  it('renders DefaultError, keeps modal closable, and skips Navigation queries when denied', () => {
    mockUseHasPermission.mockReturnValue(false);
    const ctx = { ...baseCtx(), open: true };

    render(wrap(<RegistryModal />, ctx));

    const modal = screen.getByTestId('registry-modal');
    expect(modal).toHaveAttribute('data-title', 'Registry');
    expect(screen.getByTestId('default-error')).toHaveAttribute(
      'data-permission-key',
      JSON.stringify(Keys.MesherySystemViewRegistry),
    );
    expect(screen.queryByTestId('mesh-model-component')).not.toBeInTheDocument();
    expect(mockUseGetMeshModelsQuery).not.toHaveBeenCalled();
    expect(mockUseGetComponentsQuery).not.toHaveBeenCalled();
    expect(mockUseGetRelationshipsQuery).not.toHaveBeenCalled();
    expect(mockUseGetRegistrantsQuery).not.toHaveBeenCalled();

    fireEvent.click(screen.getByTestId('close-registry-modal'));
    expect(ctx.closeModal).toHaveBeenCalledTimes(1);
  });

  it('resets modal header title and unmounts Navigation when permission is revoked while open', () => {
    const ctx = { ...baseCtx(), open: true };
    const { rerender } = render(wrap(<RegistryModal />, ctx));

    expect(screen.getByTestId('registry-modal')).toHaveAttribute('data-title', 'Registry - Models');
    expect(screen.getByTestId('mesh-model-component')).toBeInTheDocument();

    mockUseHasPermission.mockReturnValue(false);
    rerender(wrap(<RegistryModal />, ctx));

    expect(screen.getByTestId('registry-modal')).toHaveAttribute('data-title', 'Registry');
    expect(screen.getByTestId('default-error')).toBeInTheDocument();
    expect(screen.queryByTestId('mesh-model-component')).not.toBeInTheDocument();
  });
});

describe('Navigation', () => {
  beforeEach(() => {
    mediaQueryReturn = false;
    mockUseHasPermission.mockReset();
    mockUseHasPermission.mockReturnValue(true);
    mockUseGetMeshModelsQuery.mockReset();
    mockUseGetMeshModelsQuery.mockReturnValue({
      data: { models: [] },
      isLoading: false,
    });
    mockUseGetComponentsQuery.mockReset();
    mockUseGetComponentsQuery.mockReturnValue({
      data: { totalCount: 7 },
      isLoading: false,
    });
    mockUseGetRelationshipsQuery.mockReset();
    mockUseGetRelationshipsQuery.mockReturnValue({
      data: { totalCount: 3 },
      isLoading: false,
    });
    mockUseGetRegistrantsQuery.mockReset();
    mockUseGetRegistrantsQuery.mockReturnValue({
      data: { totalCount: 5 },
      isLoading: false,
    });
  });

  const baseCtx = () => ({
    open: true,
    closeModal: vi.fn(),
    selectedView: 'Models',
    searchText: 'hello',
    selectedItemUUID: 'uuid-1',
    setSelectedView: vi.fn(),
  });

  it('renders nav items with counts and the wrapped MeshModelComponent', () => {
    render(wrap(<Navigation setHeaderInfo={vi.fn()} />, baseCtx()));

    expect(screen.getByText('Models (0)')).toBeInTheDocument();
    expect(screen.getByText('Components (7)')).toBeInTheDocument();
    expect(screen.getByText('Relationships (3)')).toBeInTheDocument();
    expect(screen.getByText('Registrants (5)')).toBeInTheDocument();

    const wrapped = screen.getByTestId('mesh-model-component');
    expect(wrapped).toHaveAttribute('data-view', 'Models');
    expect(wrapped).toHaveAttribute('data-search', 'hello');
    expect(wrapped).toHaveAttribute('data-uuid', 'uuid-1');
  });

  it('clicking a nav item invokes setSelectedView and setHeaderInfo', () => {
    const setHeaderInfo = vi.fn();
    const ctx = baseCtx();
    render(wrap(<Navigation setHeaderInfo={setHeaderInfo} />, ctx));

    const componentsButton = screen
      .getByText('Components (7)')
      .closest('button') as HTMLButtonElement;
    fireEvent.click(componentsButton);

    expect(ctx.setSelectedView).toHaveBeenCalledWith('Components');
    expect(setHeaderInfo).toHaveBeenCalled();
  });

  it('skips all registry count queries, ignores cached counts, and renders DefaultError when standalone Navigation is denied', () => {
    mockUseHasPermission.mockReturnValue(false);
    const setHeaderInfo = vi.fn();
    mockUseGetMeshModelsQuery.mockReturnValue({
      data: { models: [{ id: 'm1' }] },
      isLoading: false,
    });

    render(wrap(<Navigation setHeaderInfo={setHeaderInfo} />, baseCtx()));

    expect(mockUseGetMeshModelsQuery).toHaveBeenCalledWith(
      { params: { pagesize: 'all' } },
      { skip: true },
    );
    expect(mockUseGetComponentsQuery).toHaveBeenCalledWith(
      { params: { pagesize: 'all' } },
      { skip: true },
    );
    expect(mockUseGetRelationshipsQuery).toHaveBeenCalledWith(
      { params: { pagesize: 'all' } },
      { skip: true },
    );
    expect(mockUseGetRegistrantsQuery).toHaveBeenCalledWith(
      { params: { pagesize: 'all' } },
      { skip: true },
    );
    expect(setHeaderInfo).not.toHaveBeenCalled();
    expect(screen.getByTestId('default-error')).toHaveAttribute(
      'data-permission-key',
      JSON.stringify(Keys.MesherySystemViewRegistry),
    );
    expect(screen.queryByTestId('mesh-model-component')).not.toBeInTheDocument();
    expect(screen.queryByText(/Components \(7\)/)).not.toBeInTheDocument();
  });
});
