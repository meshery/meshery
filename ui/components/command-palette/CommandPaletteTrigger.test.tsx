import React, { createContext, useContext } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import CommandPaletteContextProvider, {
  useCommandPalette,
} from '@/utils/context/CommandPaletteContextProvider';
import { CommandPaletteTrigger } from './CommandPaletteTrigger';

vi.mock('./CommandPalette', () => ({
  CommandPalette: () => null,
}));

const createMockTheme = (mode: 'light' | 'dark' = 'light') => ({
  palette: {
    mode,
    background: {
      card: mode === 'light' ? 'white' : 'black',
    },
    common: {
      white: 'white',
    },
    text: {
      constant: {
        white: 'white',
      },
    },
  },
});

const ThemeCtx = createContext(createMockTheme('light'));

vi.mock('@sistent/sistent', () => {
  const useTheme = () => useContext(ThemeCtx);

  const Box = ({ component = 'div', children, sx: _sx, ...props }: any) =>
    React.createElement(component, props, children);

  const Typography = ({ variant: _variant, children, sx: _sx, ...props }: any) => (
    <span {...props}>{children}</span>
  );

  const SearchIcon = (props: any) => <svg data-testid="search-icon" {...props} />;

  const Button = ({ children, sx: _sx, ...props }: any) => <button {...props}>{children}</button>;

  const CustomTooltip = ({ title, children }: any) => (
    <div data-testid="custom-tooltip" data-tooltip-title={title}>
      {children}
    </div>
  );

  const alpha = (color: string, _opacity: number) => color;

  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) =>
      React.createElement(Component || 'div', props, children);
    return StyledComponent;
  };

  class MockEventBus {
    listeners: Record<string, Function[]> = {};
    on() {}
    off() {}
    emit() {}
  }

  const mockExports = {
    useTheme,
    Box,
    Typography,
    SearchIcon,
    Button,
    CustomTooltip,
    alpha,
    styled,
    EventBus: MockEventBus,
    createCanShow: () => () => true,
    InfoIcon: () => null,
    notificationColors: {
      info: { main: 'teal' },
      error: { main: 'red', dark: 'darkred' },
      warning: { main: 'orange' },
      success: { main: 'green' },
    },
  };

  return new Proxy(mockExports, {
    has(_target, _prop) {
      return true;
    },
    get(target: any, prop: string | symbol) {
      if (
        typeof prop === 'symbol' ||
        prop === 'then' ||
        prop === 'default' ||
        prop === '__esModule'
      ) {
        return target[prop];
      }
      if (prop in target) {
        return target[prop];
      }
      if (typeof prop === 'string') {
        const Fallback = ({ children, ...props }: any) =>
          React.createElement('div', props, children);
        Fallback.displayName = `SistentMock_${prop}`;
        return Fallback;
      }
      return undefined;
    },
  });
});

// Isolate singleton store initialized by can.ts -> store
vi.mock('@/store', () => ({
  store: { getState: () => ({}) },
}));

// Spy on context state
let currentOpenState = false;
const PaletteStateConsumer = () => {
  const { open } = useCommandPalette();
  currentOpenState = open;
  return <div data-testid="palette-open-state">{String(open)}</div>;
};

const renderTrigger = (themeMode: 'light' | 'dark' = 'light') => {
  return render(
    <ThemeCtx.Provider value={createMockTheme(themeMode)}>
      <CommandPaletteContextProvider>
        <PaletteStateConsumer />
        <CommandPaletteTrigger />
      </CommandPaletteContextProvider>
    </ThemeCtx.Provider>,
  );
};

describe('CommandPaletteTrigger', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders trigger with data-testid, aria-label, aria-keyshortcuts, and search icon', () => {
    renderTrigger();
    const trigger = screen.getByTestId('command-palette-trigger');
    expect(trigger).toBeInTheDocument();
    expect(trigger).toHaveAttribute('aria-label', 'Search commands');
    expect(trigger).toHaveAttribute('aria-keyshortcuts', 'Control+K Meta+K');

    const icon = screen.getByTestId('search-icon');
    expect(icon).toBeInTheDocument();
  });

  it('is wrapped in CustomTooltip with default Ctrl K title', () => {
    renderTrigger();
    const tooltip = screen.getByTestId('custom-tooltip');
    expect(tooltip).toHaveAttribute('data-tooltip-title', 'Search commands (Ctrl K)');
  });

  it('clicking trigger opens the palette via context', () => {
    renderTrigger();
    expect(currentOpenState).toBe(false);

    const trigger = screen.getByTestId('command-palette-trigger');
    fireEvent.click(trigger);
    expect(currentOpenState).toBe(true);
  });

  it('renders correctly in both light and dark mode themes', () => {
    const { unmount: unmountLight } = renderTrigger('light');
    expect(screen.getByTestId('command-palette-trigger')).toBeInTheDocument();
    unmountLight();

    const { unmount: unmountDark } = renderTrigger('dark');
    expect(screen.getByTestId('command-palette-trigger')).toBeInTheDocument();
    unmountDark();
  });

  it('detects Mac platform and displays ⌘K shortcut hint', () => {
    const originalNavigator = window.navigator;
    try {
      Object.defineProperty(window, 'navigator', {
        value: {
          platform: 'MacIntel',
          userAgentData: { platform: 'macOS' },
        },
        configurable: true,
      });

      renderTrigger();
      const tooltip = screen.getByTestId('custom-tooltip');
      expect(tooltip).toHaveAttribute('data-tooltip-title', 'Search commands (⌘K)');
      expect(screen.getByText('⌘K')).toBeInTheDocument();
    } finally {
      Object.defineProperty(window, 'navigator', {
        value: originalNavigator,
        configurable: true,
      });
    }
  });

  it('shortcut kbd element has aria-hidden="true"', () => {
    renderTrigger();
    const kbd = screen.getByText('Ctrl K');
    expect(kbd.tagName.toLowerCase()).toBe('kbd');
    expect(kbd).toHaveAttribute('aria-hidden', 'true');
  });
});
