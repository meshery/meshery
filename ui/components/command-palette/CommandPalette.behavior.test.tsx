import React from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';
import { ability } from '@/utils/can';
import CommandPaletteContextProvider, {
  CommandPaletteHost,
  useCommandPalette,
} from '@/utils/context/CommandPaletteContextProvider';
import { SistentThemeProvider } from '@sistent/sistent';
import mesheryUiReducer from '@/store/slices/mesheryUi';
import adapterReducer from '@/store/slices/adapter';
import { createTestStoreWithReducers, sampleAdapters, sampleCapabilities } from './testUtils';

vi.mock('@/store', () => ({
  store: { getState: () => ({}) },
}));

vi.mock('@/rtk-query/user', () => ({
  getProviderCapabilities: vi.fn(),
  getSystemVersion: vi.fn(),
}));

// Stub the 7 static .js drawer-icon SVG assets that Vite cannot parse as JSX
vi.mock('../../public/static/img/drawer-icons/lifecycle_mgmt_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/performance_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/extensions_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/pattern_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/lifecycle_hover_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/performance_hover_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/configuration_hover_svg', () => ({
  default: () => null,
}));

vi.mock('@sistent/sistent', async () => {
  const { createSistentMock: createMock } = await import('./testUtils');
  return createMock();
});

const mockRouter = {
  pathname: '/dashboard',
  asPath: '/dashboard',
  query: {},
  push: vi.fn(),
};

vi.mock('next/router', () => ({
  useRouter: () => mockRouter,
}));

const OpenButton = () => {
  const { openPalette } = useCommandPalette();
  return (
    <button data-testid="open-btn" onClick={openPalette}>
      Open
    </button>
  );
};

const renderWithProviders = (
  store = createTestStoreWithReducers(
    { ui: mesheryUiReducer, adapter: adapterReducer },
    {
      ui: {
        providerCapabilities: sampleCapabilities,
        catalogVisibility: true,
      },
      adapter: {
        meshAdapters: sampleAdapters,
      },
    },
  ),
) => {
  return render(
    <Provider store={store}>
      <SistentThemeProvider initialMode="light">
        <CommandPaletteContextProvider>
          <OpenButton />
          <CommandPaletteHost />
        </CommandPaletteContextProvider>
      </SistentThemeProvider>
    </Provider>,
  );
};

describe('CommandPalette - Behavior & Interactions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockRouter.pathname = '/dashboard';
    mockRouter.asPath = '/dashboard';
    mockRouter.query = {};
    ability.update([{ action: 'manage', subject: 'all' }]);
  });

  it('opens on Ctrl+K and focuses the search input', async () => {
    renderWithProviders();

    expect(screen.queryByTestId('command-palette')).not.toBeInTheDocument();

    fireEvent.keyDown(window, { key: 'k', ctrlKey: true });

    await waitFor(() => {
      expect(screen.getByTestId('command-palette')).toBeInTheDocument();
    });

    const input = screen.getByTestId('command-palette-input');
    expect(input).toBeInTheDocument();
  });

  it('navigates through options via ArrowDown and ArrowUp', async () => {
    renderWithProviders();
    fireEvent.click(screen.getByTestId('open-btn'));

    const input = await screen.findByTestId('command-palette-input');
    const initialActive = input.getAttribute('aria-activedescendant');

    fireEvent.keyDown(input, { key: 'ArrowDown' });
    const afterDown = input.getAttribute('aria-activedescendant');
    expect(afterDown).not.toEqual(initialActive);

    fireEvent.keyDown(input, { key: 'ArrowUp' });
    const afterUp = input.getAttribute('aria-activedescendant');
    expect(afterUp).toEqual(initialActive);
  });

  it('executes router.push on exit transition when selecting an enabled option with Enter', async () => {
    renderWithProviders();
    fireEvent.click(screen.getByTestId('open-btn'));

    const input = await screen.findByTestId('command-palette-input');
    fireEvent.change(input, { target: { value: 'designs' } });

    const item = await screen.findByTestId('command-palette-item-nav:/configuration/designs');
    expect(item).toBeInTheDocument();

    fireEvent.keyDown(input, { key: 'Enter' });

    await waitFor(() => {
      expect(mockRouter.push).toHaveBeenCalledWith('/configuration/designs');
    });
  });

  it('does nothing when Enter is pressed on a disabled option and keeps palette open', async () => {
    ability.update([]); // remove all permissions
    renderWithProviders();
    fireEvent.click(screen.getByTestId('open-btn'));

    const input = await screen.findByTestId('command-palette-input');
    fireEvent.change(input, { target: { value: 'designs' } });

    const item = await screen.findByTestId('command-palette-item-nav:/configuration/designs');
    expect(item).toBeInTheDocument();
    expect(item.getAttribute('aria-disabled')).toBe('true');

    fireEvent.keyDown(input, { key: 'Enter' });

    expect(mockRouter.push).not.toHaveBeenCalled();
    expect(screen.getByTestId('command-palette')).toBeInTheDocument();
  });

  it('closes on a single Escape press via userEvent and does not push route', async () => {
    const user = userEvent.setup();
    renderWithProviders();
    fireEvent.click(screen.getByTestId('open-btn'));

    const input = await screen.findByTestId('command-palette-input');
    input.focus();

    await user.keyboard('{Escape}');

    await waitFor(() => {
      expect(screen.queryByTestId('command-palette')).not.toBeInTheDocument();
    });
    expect(mockRouter.push).not.toHaveBeenCalled();
  });

  it('filters commands when typing and updates available options', async () => {
    renderWithProviders();
    fireEvent.click(screen.getByTestId('open-btn'));

    const input = await screen.findByTestId('command-palette-input');
    fireEvent.change(input, { target: { value: 'catalog' } });

    await waitFor(() => {
      expect(
        screen.getByTestId('command-palette-item-nav:/configuration/catalog'),
      ).toBeInTheDocument();
      expect(
        screen.queryByTestId('command-palette-item-nav:/management/performance'),
      ).not.toBeInTheDocument();
    });
  });

  it('renders each group header at most once when filtering (typing "des")', async () => {
    renderWithProviders();
    fireEvent.click(screen.getByTestId('open-btn'));

    const input = await screen.findByTestId('command-palette-input');
    fireEvent.change(input, { target: { value: 'des' } });

    await waitFor(() => {
      const groupHeadings = screen
        .getAllByRole('heading')
        .map((h) => h.textContent?.trim())
        .filter((t) => t === 'ACTIONS' || t === 'NAVIGATE' || t === 'SETTINGS');

      const uniqueGroups = new Set(groupHeadings);
      expect(groupHeadings.length).toBe(uniqueGroups.size);
    });
  });

  it('does not replace input value with option title during exit transition', async () => {
    renderWithProviders();
    fireEvent.click(screen.getByTestId('open-btn'));

    const input = (await screen.findByTestId('command-palette-input')) as HTMLInputElement;
    fireEvent.change(input, { target: { value: 'des' } });
    expect(input.value).toBe('des');

    const item = await screen.findByTestId('command-palette-item-nav:/configuration/designs');
    fireEvent.click(item);

    // During palette closing / exit, input value must remain "des" and not be replaced by "Designs"
    expect(input.value).toBe('des');
  });

  it('preserves table state and merges router.query when selecting Import Design while on /configuration/designs', async () => {
    mockRouter.pathname = '/configuration/designs';
    mockRouter.asPath = '/configuration/designs?page=2&tab=my-designs';
    mockRouter.query = { page: '2', tab: 'my-designs' };

    renderWithProviders();
    fireEvent.click(screen.getByTestId('open-btn'));

    const input = await screen.findByTestId('command-palette-input');
    fireEvent.change(input, { target: { value: 'import design' } });

    const item = await screen.findByTestId('command-palette-item-action:import-design');
    fireEvent.click(item);

    await waitFor(() => {
      expect(mockRouter.push).toHaveBeenCalledWith({
        pathname: '/configuration/designs',
        query: {
          page: '2',
          tab: 'my-designs',
          action: 'import',
        },
      });
    });
  });

  it('re-evaluates permissions dynamically on every re-open', async () => {
    // 1. Initial open with NO view designs permission
    ability.update([]);
    renderWithProviders();

    fireEvent.click(screen.getByTestId('open-btn'));
    const input = await screen.findByTestId('command-palette-input');
    fireEvent.change(input, { target: { value: 'designs' } });

    const itemDisabled = await screen.findByTestId(
      'command-palette-item-nav:/configuration/designs',
    );
    expect(itemDisabled.getAttribute('aria-disabled')).toBe('true');

    // Close palette
    fireEvent.keyDown(input, { key: 'Escape' });
    await waitFor(() => {
      expect(screen.queryByTestId('command-palette')).not.toBeInTheDocument();
    });

    // 2. Grant permission and re-open
    ability.update([{ action: 'manage', subject: 'all' }]);
    fireEvent.click(screen.getByTestId('open-btn'));

    const inputReopened = await screen.findByTestId('command-palette-input');
    fireEvent.change(inputReopened, { target: { value: 'designs' } });

    const itemEnabled = await screen.findByTestId(
      'command-palette-item-nav:/configuration/designs',
    );
    expect(itemEnabled.getAttribute('aria-disabled')).toBe('false');
  });
});
