import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const notify = vi.fn();
const dispatch = vi.fn();
const updateUserPref = vi.fn().mockImplementation(() => ({
  unwrap: () => Promise.resolve({}),
}));
const updateUserPrefWithContext = vi.fn();

let selectedStore: any = {
  ui: { providerCapabilities: { extensions: {} } },
};

let userDataReturn: any = {
  data: undefined,
  isSuccess: false,
  isError: false,
  error: null,
};

let capabilitiesReturn: any = {
  data: undefined,
  isSuccess: false,
};

vi.mock('@sistent/sistent', () => ({
  // This suite exercises preference toggles, not authorization: grant every
  // capability so the permission gate never masks the behaviour under test.
  useHasPermission: () => true,
  Tab: ({ label, ...rest }: any) => <button {...rest}>{label}</button>,
  Tabs: ({ children, value, onChange }: any) => (
    <div data-testid="tabs" data-value={value} onClick={(e) => onChange?.(e, 0)}>
      {children}
    </div>
  ),
  Typography: ({ children, ...rest }: any) => <span {...rest}>{children}</span>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  FormGroup: ({ children }: any) => <div>{children}</div>,
  FormControlLabel: ({ label, control }: any) => (
    <label>
      {control}
      {label}
    </label>
  ),
  Switch: ({ checked, onChange, ...rest }: any) => (
    <input
      type="checkbox"
      checked={!!checked}
      onChange={onChange}
      data-testid={rest['data-cy'] || 'switch'}
    />
  ),
  IconButton: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
  Select: ({ children, ...rest }: any) => <select {...rest}>{children}</select>,
  Box: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
  TextField: (props: any) => <input {...props} />,
  Chip: ({ label, ...rest }: any) => <span {...rest}>{label}</span>,
  Button: ({ children, onClick, ...rest }: any) => (
    <button onClick={onClick} {...rest}>
      {children}
    </button>
  ),
  Divider: (props: any) => <hr {...props} />,
  LeftArrowIcon: () => <svg />,
  CardContent: ({ children }: any) => <div>{children}</div>,
  CardHeader: ({ children }: any) => <div>{children}</div>,
  CustomTooltip: ({ children }: any) => <>{children}</>,
  NoSsr: ({ children }: any) => <>{children}</>,
  TachometerIcon: () => <svg />,
  SettingsCellIcon: () => <svg />,
  SettingsRemoteIcon: () => <svg />,
  useTheme: () => ({
    palette: { icon: { default: 'icon' }, mode: 'light' },
  }),
  ErrorBoundary: ({ children }: any) => <>{children}</>,
  // `DefaultError` (rendered by the new permission guard) imports `styled`
  // from error-404/styles.tsx. See the same comment in workspaces/index.test.tsx.
  styled: (Component: any) => () => {
    if (typeof Component === 'string') {
      return ({ children, ...rest }: any) => <Component {...rest}>{children}</Component>;
    }
    return Component;
  },
}));

vi.mock('../general/error-404/index', () => ({
  default: ({ permissionKey }: any) => <div data-testid="default-error">{permissionKey}</div>,
}));

vi.mock('../../assets/icons/CopyIcon', () => ({
  default: () => <svg />,
}));

vi.mock('./style', () => {
  const make =
    (tag = 'div') =>
    ({ children, ...rest }: any) =>
      React.createElement(tag, rest, children);

  return {
    StatsWrapper: make('div'),
    ProviderCard: make('div'),
    RootContainer: make('div'),
    BoxWrapper: make('div'),
    Divider: make('hr'),
    GridCapabilityHeader: make('div'),
    GridExtensionHeader: make('div'),
    GridExtensionItem: make('div'),
    StyledPaper: make('div'),
    TabLabel: make('span'),
    HideScrollbar: make('div'),
    IconStyled: make('span'),
    FormLegend: make('legend'),
    FormContainerWrapper: make('div'),
    FormGroupWrapper: make('fieldset'),
  };
});

vi.mock('../ExtensionSandbox', () => ({
  default: () => <div data-testid="extension-sandbox" />,
}));

vi.mock('../general/RemoteComponent', () => ({
  default: () => <div data-testid="remote-component" />,
}));

vi.mock('../../utils/ExtensionPointSchemaValidator', () => ({
  default: () => () => null,
}));

vi.mock('../settings/MesherySettingsPerformanceComponent', () => ({
  default: () => <div data-testid="perf-component" />,
}));

vi.mock('../../css/icons.styles', () => ({
  iconMedium: {},
}));

vi.mock('../../lib/event-types', () => ({
  EVENT_TYPES: { SUCCESS: 'success', ERROR: 'error' },
}));

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('@/utils/dimension', () => ({
  useWindowDimensions: () => ({ width: 1024, height: 768 }),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetProviderCapabilitiesQuery: () => capabilitiesReturn,
  useGetUserPrefQuery: () => userDataReturn,
  useUpdateUserPrefMutation: () => [updateUserPref],
  useUpdateUserPrefWithContextMutation: () => [updateUserPrefWithContext],
}));

vi.mock('@/theme/hooks', () => ({
  ThemeTogglerCore: ({ Component }: any) => (
    <div data-testid="theme-toggler">
      <Component mode="light" toggleTheme={() => {}} />
    </div>
  ),
}));

vi.mock('../dashboard/style', () => ({
  SecondaryTab: ({ label }: any) => <button>{label}</button>,
  SecondaryTabs: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('react-redux', () => ({
  useSelector: (fn: any) => fn(selectedStore),
  useDispatch: () => dispatch,
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  toggleCatalogContent: vi.fn((arg: any) => arg),
  updateProgress: vi.fn(),
}));

// UserPreference component is the default export
import UserPreference from './index';

describe('UserPreference', () => {
  beforeEach(() => {
    notify.mockReset();
    dispatch.mockReset();
    updateUserPref.mockClear();
    updateUserPrefWithContext.mockClear();

    userDataReturn = {
      data: { usersExtensionPreferences: { catalogContent: true } },
      isSuccess: true,
      isError: false,
      error: null,
    };

    capabilitiesReturn = {
      data: {
        providerName: 'Test',
        providerType: 'remote',
        capabilities: [{ feature: 'A', endpoint: '/a' }],
        extensions: {},
      },
      isSuccess: true,
    };

    selectedStore = { ui: { providerCapabilities: { extensions: {} } } };
  });

  it('renders the General and Performance tabs by default', () => {
    render(<UserPreference />);
    expect(screen.getByText('General')).toBeInTheDocument();
    expect(screen.getByText('Performance')).toBeInTheDocument();
  });

  it('renders the Catalog Content toggle on the General tab', () => {
    render(<UserPreference />);
    const toggle = screen.getByTestId('CatalogContentPreference') as HTMLInputElement;
    expect(toggle).toBeInTheDocument();
  });

  it('toggling Catalog Content dispatches and calls updateUserPref', () => {
    render(<UserPreference />);
    const toggle = screen.getByTestId('CatalogContentPreference');
    fireEvent.click(toggle);
    expect(dispatch).toHaveBeenCalled();
    expect(updateUserPref).toHaveBeenCalled();
  });
});
