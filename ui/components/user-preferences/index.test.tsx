import React from 'react';
import { act, render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Keys } from '@meshery/schemas/permissions';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const notify = vi.fn();
const dispatch = vi.fn();
const mockUseHasPermission = vi.fn(() => true);
const mockUseGetProviderCapabilitiesQuery = vi.fn(() => capabilitiesReturn);
const mockUseGetUserPrefQuery = vi.fn(() => userDataReturn);
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

let catalogMutationReturn: any = {};

vi.mock('@sistent/sistent', () => ({
  // Default to granted so existing toggle specs aren't masked; deny specs override mockUseHasPermission.
  useHasPermission: (permissionKey: any) => mockUseHasPermission(permissionKey),
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
      disabled={rest.disabled}
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
  default: ({ permissionKey }: any) => (
    <div data-testid="default-error" data-permission-key={JSON.stringify(permissionKey)} />
  ),
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
  useGetProviderCapabilitiesQuery: (...args: any[]) => mockUseGetProviderCapabilitiesQuery(...args),
  useGetUserPrefQuery: (...args: any[]) => mockUseGetUserPrefQuery(...args),
  useUpdateUserPrefMutation: () => [updateUserPref, catalogMutationReturn],
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
    mockUseHasPermission.mockReset();
    mockUseHasPermission.mockReturnValue(true);
    mockUseGetUserPrefQuery.mockClear();
    mockUseGetProviderCapabilitiesQuery.mockClear();
    updateUserPref.mockClear();
    updateUserPrefWithContext.mockClear();
    catalogMutationReturn = {};

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

  it.each([true, false])(
    'restores the saved Catalog Content value %s on failure',
    async (saved) => {
      userDataReturn.data.usersExtensionPreferences.catalogContent = saved;
      updateUserPref.mockImplementationOnce(() => ({
        unwrap: () => Promise.reject(new Error('network down')),
      }));
      render(<UserPreference />);
      fireEvent.click(screen.getByTestId('CatalogContentPreference'));

      await waitFor(() =>
        expect(notify).toHaveBeenCalledWith({
          message: 'There was an error sending your preference',
          event_type: 'error',
        }),
      );
      expect(screen.getByTestId('CatalogContentPreference')).toHaveProperty('checked', saved);
      expect(dispatch).toHaveBeenLastCalledWith({ catalogVisibility: saved });
    },
  );

  it('rolls back to a value successfully saved earlier in this session', async () => {
    render(<UserPreference />);
    const toggle = screen.getByTestId('CatalogContentPreference');
    fireEvent.click(toggle);
    await waitFor(() =>
      expect(notify).toHaveBeenCalledWith({
        message: 'Catalog Content was disabled',
        event_type: 'success',
      }),
    );

    updateUserPref.mockImplementationOnce(() => ({
      unwrap: () => Promise.reject(new Error('network down')),
    }));
    fireEvent.click(toggle);
    await waitFor(() =>
      expect(notify).toHaveBeenCalledWith({
        message: 'There was an error sending your preference',
        event_type: 'error',
      }),
    );
    expect(toggle).not.toBeChecked();
    expect(dispatch).toHaveBeenLastCalledWith({ catalogVisibility: false });
  });

  it('blocks repeated toggles until the active save finishes', async () => {
    let resolveSave!: () => void;
    const pendingSave = new Promise((resolve) => {
      resolveSave = () => resolve({});
    });
    updateUserPref.mockImplementationOnce(() => ({ unwrap: () => pendingSave }));
    const { rerender } = render(<UserPreference />);
    const toggle = screen.getByTestId('CatalogContentPreference');
    fireEvent.click(toggle);
    fireEvent.click(toggle); // Exercise the ref guard before mutation state rerenders.
    catalogMutationReturn = { isLoading: true, startedTimeStamp: 2 };
    userDataReturn = {
      ...userDataReturn,
      data: { usersExtensionPreferences: { catalogContent: true, otherPreference: 'changed' } },
    };
    rerender(<UserPreference />);

    try {
      expect(toggle).toBeDisabled();
      fireEvent.click(toggle);
      expect(updateUserPref).toHaveBeenCalledTimes(1);
      expect(toggle).not.toBeChecked();
    } finally {
      await act(async () => resolveSave());
    }

    catalogMutationReturn = { isLoading: false, fulfilledTimeStamp: 3 };
    userDataReturn = {
      ...userDataReturn,
      data: { usersExtensionPreferences: { catalogContent: false } },
      startedTimeStamp: 4,
    };
    rerender(<UserPreference />);
    expect(toggle).toBeEnabled();
  });

  it('blocks a new toggle while saved preferences are being refetched', () => {
    userDataReturn.isFetching = true;
    render(<UserPreference />);
    const toggle = screen.getByTestId('CatalogContentPreference');
    expect(toggle).toBeDisabled();
    fireEvent.click(toggle);
    expect(updateUserPref).not.toHaveBeenCalled();
  });

  it('skips preference queries, ignores cached preferences, and renders DefaultError when permission is denied', () => {
    mockUseHasPermission.mockReturnValue(false);
    render(<UserPreference />);

    expect(mockUseHasPermission).toHaveBeenCalledWith(Keys.ExtensibilityViewMesheryUserPreferences);
    expect(mockUseGetUserPrefQuery).toHaveBeenCalledWith(undefined, { skip: true });
    expect(mockUseGetProviderCapabilitiesQuery).toHaveBeenCalledWith(undefined, { skip: true });
    expect(screen.getByTestId('default-error')).toHaveAttribute(
      'data-permission-key',
      JSON.stringify(Keys.ExtensibilityViewMesheryUserPreferences),
    );
    expect(screen.queryByTestId('CatalogContentPreference')).not.toBeInTheDocument();
  });
});
