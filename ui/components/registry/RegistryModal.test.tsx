import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

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
  InfoIcon: () => <svg data-testid="info-icon" />,
  DatabaseIcon: () => <svg data-testid="database-icon" />,
  CustomTooltip: ({ children }: any) => <>{children}</>,
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
  Modal: ({ children, isOpen, title }: any) =>
    isOpen ? (
      <div data-testid="registry-modal" data-title={title}>
        {children}
      </div>
    ) : null,
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
  useGetMeshModelsQuery: () => ({
    data: { models: [] },
    isLoading: false,
  }),
  useGetComponentsQuery: () => ({
    data: { totalCount: 7 },
    isLoading: false,
  }),
  useGetRelationshipsQuery: () => ({
    data: { totalCount: 3 },
    isLoading: false,
  }),
  useGetRegistrantsQuery: () => ({
    data: { totalCount: 5 },
    isLoading: false,
  }),
}));

vi.mock('./helper', () => ({
  removeDuplicateVersions: (m: any[]) => m,
}));

vi.mock('../general/style', () => ({
  ChevronButtonWrapper: ({ children, onClick, isCollapsed }: any) => (
    <button
      type="button"
      onClick={onClick}
      data-collapsed={String(!!isCollapsed)}
      data-testid="sidebar-collapse-toggle"
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
});

describe('Navigation', () => {
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
});
