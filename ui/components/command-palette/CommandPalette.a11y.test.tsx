import React from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
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

const renderWithTheme = (
  themeMode: 'light' | 'dark' = 'light',
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
      <SistentThemeProvider initialMode={themeMode}>
        <CommandPaletteContextProvider>
          <OpenButton />
          <CommandPaletteHost />
        </CommandPaletteContextProvider>
      </SistentThemeProvider>
    </Provider>,
  );
};

describe('CommandPalette - Accessibility & States', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockRouter.pathname = '/dashboard';
    mockRouter.asPath = '/dashboard';
    mockRouter.query = {};
    ability.update([{ action: 'manage', subject: 'all' }]);
  });

  it('has accessible name "Command palette" on role="dialog"', async () => {
    renderWithTheme('light');
    fireEvent.click(screen.getByTestId('open-btn'));

    const dialog = await screen.findByRole('dialog', { name: 'Command palette' });
    expect(dialog).toBeInTheDocument();
    expect(dialog).toHaveAttribute('data-testid', 'command-palette');
  });

  it('contains input with role="combobox" and correct aria attributes', async () => {
    renderWithTheme('light');
    fireEvent.click(screen.getByTestId('open-btn'));

    const combobox = await screen.findByRole('combobox', { name: 'Search commands' });
    expect(combobox).toBeInTheDocument();
    expect(screen.getByRole('combobox', { name: 'Search commands' })).toBe(combobox);
    expect(combobox.tagName).toBe('INPUT');
    expect(combobox).toHaveAttribute('data-testid', 'command-palette-input');
    expect(combobox).toHaveAttribute('aria-autocomplete', 'list');
    expect(combobox).toHaveAttribute('aria-expanded', 'true');
    expect(combobox).toHaveAttribute('aria-label', 'Search commands');

    const searchIcon = screen.getByTestId('search-icon');
    expect(searchIcon).toBeInTheDocument();
    expect(combobox.parentElement).toContainElement(searchIcon);
  });

  it('renders commands as role="option" elements', async () => {
    renderWithTheme('light');
    fireEvent.click(screen.getByTestId('open-btn'));

    await waitFor(() => {
      const options = screen.getAllByRole('option');
      expect(options.length).toBeGreaterThan(0);
    });
  });

  it('marks disabled items with aria-disabled="true" and displays visible "Requires: ..."', async () => {
    ability.update([]); // Disallow all permissions
    renderWithTheme('light');
    fireEvent.click(screen.getByTestId('open-btn'));

    const input = await screen.findByTestId('command-palette-input');
    fireEvent.change(input, { target: { value: 'designs' } });

    const item = await screen.findByTestId('command-palette-item-nav:/configuration/designs');
    expect(item).toBeInTheDocument();
    expect(item).toHaveAttribute('aria-disabled', 'true');
    expect(item).toHaveTextContent(/Requires:/i);
  });

  it('announces filtered results through visually hidden role="status" announcer', async () => {
    renderWithTheme('light');
    fireEvent.click(screen.getByTestId('open-btn'));

    const statusAnnouncer = await screen.findByRole('status');
    expect(statusAnnouncer).toBeInTheDocument();
    expect(statusAnnouncer.textContent).toMatch(/\d+ results?/i);

    const input = screen.getByTestId('command-palette-input');
    fireEvent.change(input, { target: { value: 'catalog' } });

    await waitFor(() => {
      expect(statusAnnouncer.textContent).toMatch(/1 result/i);
    });
  });

  it('shows "No commands available" when providerCapabilities is null', async () => {
    const storeWithoutCapabilities = createTestStoreWithReducers(
      { ui: mesheryUiReducer, adapter: adapterReducer },
      {
        ui: {
          providerCapabilities: null,
          catalogVisibility: true,
        },
        adapter: {
          meshAdapters: sampleAdapters,
        },
      },
    );

    renderWithTheme('light', storeWithoutCapabilities);
    fireEvent.click(screen.getByTestId('open-btn'));

    await waitFor(() => {
      expect(screen.getAllByText('No commands available').length).toBeGreaterThanOrEqual(1);
    });
  });

  it('shows "No commands found" when query matches nothing', async () => {
    renderWithTheme('light');
    fireEvent.click(screen.getByTestId('open-btn'));

    const input = await screen.findByTestId('command-palette-input');
    fireEvent.change(input, { target: { value: 'xyz-nonexistent-command' } });

    await waitFor(() => {
      expect(screen.getAllByText('No commands found').length).toBeGreaterThanOrEqual(1);
      const statusAnnouncer = screen.getByRole('status');
      expect(statusAnnouncer.textContent).toBe('No commands found');
    });
  });

  it('renders section group headers with role="presentation" wrapper and role="group" labelled ul', async () => {
    renderWithTheme('light');
    fireEvent.click(screen.getByTestId('open-btn'));

    await waitFor(() => {
      const groups = screen.getAllByRole('group');
      expect(groups.length).toBeGreaterThan(0);
      groups.forEach((groupUl) => {
        const labelledBy = groupUl.getAttribute('aria-labelledby');
        expect(labelledBy).toBeTruthy();
        const heading = document.getElementById(labelledBy!);
        expect(heading).toBeInTheDocument();
      });
    });
  });

  it('renders correctly in light theme', async () => {
    renderWithTheme('light');
    fireEvent.click(screen.getByTestId('open-btn'));

    const dialog = await screen.findByRole('dialog', { name: 'Command palette' });
    expect(dialog).toBeInTheDocument();
  });

  it('renders correctly in dark theme', async () => {
    renderWithTheme('dark');
    fireEvent.click(screen.getByTestId('open-btn'));

    const dialog = await screen.findByRole('dialog', { name: 'Command palette' });
    expect(dialog).toBeInTheDocument();
  });
});
