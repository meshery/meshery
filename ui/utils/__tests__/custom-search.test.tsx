import React from 'react';
import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { SistentThemeProvider } from '@sistent/sistent';
import SearchBar from '../custom-search';

const renderWithTheme = (ui: React.ReactElement) => {
  return render(<SistentThemeProvider initialTheme="dark">{ui}</SistentThemeProvider>);
};

describe('SearchBar component', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.runOnlyPendingTimers();
    vi.useRealTimers();
    vi.restoreAllMocks();
  });

  it('renders search icon when collapsed', () => {
    renderWithTheme(
      <SearchBar
        expanded={false}
        setExpanded={vi.fn()}
        onSearch={vi.fn()}
        placeholder="Search here..."
      />,
    );

    const searchIcon = screen.getByTestId('search-icon');
    expect(searchIcon).toBeDefined();
    expect(screen.queryByTestId('close-icon')).toBeNull();
  });

  it('calls setExpanded(true) when search icon is clicked in collapsed state', () => {
    const setExpanded = vi.fn();
    renderWithTheme(
      <SearchBar
        expanded={false}
        setExpanded={setExpanded}
        onSearch={vi.fn()}
        placeholder="Search here..."
      />,
    );

    const searchIcon = screen.getByTestId('search-icon');
    fireEvent.click(searchIcon);

    expect(setExpanded).toHaveBeenCalledWith(true);
    act(() => {
      vi.advanceTimersByTime(350);
    });
  });

  it('renders close icon and input field when expanded', () => {
    renderWithTheme(
      <SearchBar
        expanded={true}
        setExpanded={vi.fn()}
        onSearch={vi.fn()}
        placeholder="Search here..."
        value="initial text"
      />,
    );

    const input = screen.getByPlaceholderText('Search here...') as HTMLInputElement;
    expect(input).toBeDefined();
    expect(input.value).toBe('initial text');
    expect(screen.getByTestId('close-icon')).toBeDefined();
  });

  it('triggers debounced onSearch after typing in input', () => {
    const onSearch = vi.fn();
    renderWithTheme(
      <SearchBar
        expanded={true}
        setExpanded={vi.fn()}
        onSearch={onSearch}
        placeholder="Search here..."
      />,
    );

    const input = screen.getByPlaceholderText('Search here...') as HTMLInputElement;
    fireEvent.change(input, { target: { value: 'kubernetes' } });

    expect(input.value).toBe('kubernetes');
    expect(onSearch).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(onSearch).toHaveBeenCalledWith('kubernetes');
  });

  it('handles clear icon click safely when setModelsFilters is provided', () => {
    const setExpanded = vi.fn();
    const setModelsFilters = vi.fn();
    const onSearch = vi.fn();

    renderWithTheme(
      <SearchBar
        expanded={true}
        setExpanded={setExpanded}
        onSearch={onSearch}
        setModelsFilters={setModelsFilters}
        placeholder="Search here..."
        value="test"
      />,
    );

    const closeButton = screen.getByTestId('close-icon');
    fireEvent.click(closeButton);

    expect(setModelsFilters).toHaveBeenCalledWith({ page: 0 });
    expect(setExpanded).toHaveBeenCalledWith(false);

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(onSearch).toHaveBeenCalledWith('');
  });

  it('handles clear icon click safely when setModelsFilters is NOT provided (no crash)', () => {
    const setExpanded = vi.fn();
    const onSearch = vi.fn();

    renderWithTheme(
      <SearchBar
        expanded={true}
        setExpanded={setExpanded}
        onSearch={onSearch}
        placeholder="Search here..."
        value="test"
      />,
    );

    const closeButton = screen.getByTestId('close-icon');
    expect(() => {
      fireEvent.click(closeButton);
    }).not.toThrow();

    expect(setExpanded).toHaveBeenCalledWith(false);

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(onSearch).toHaveBeenCalledWith('');
  });

  it('handles clear icon click safely when setExpanded is NOT provided', () => {
    const onSearch = vi.fn();

    renderWithTheme(
      <SearchBar expanded={true} onSearch={onSearch} placeholder="Search here..." value="test" />,
    );

    const closeButton = screen.getByTestId('close-icon');
    expect(() => {
      fireEvent.click(closeButton);
    }).not.toThrow();
  });

  it('collapses and clears text when clear icon is clicked while expanded', () => {
    const setExpanded = vi.fn();
    renderWithTheme(
      <SearchBar
        expanded={true}
        setExpanded={setExpanded}
        onSearch={vi.fn()}
        placeholder="Search here..."
        value="active search"
      />,
    );

    const closeButton = screen.getByTestId('close-icon');
    fireEvent.click(closeButton);

    expect(setExpanded).toHaveBeenCalledWith(false);
  });

  it('does not throw when onSearch is not provided', () => {
    expect(() => {
      renderWithTheme(<SearchBar expanded={true} placeholder="Search here..." />);
      const input = screen.getByPlaceholderText('Search here...');
      fireEvent.change(input, { target: { value: 'some query' } });
      act(() => {
        vi.advanceTimersByTime(500);
      });
    }).not.toThrow();
  });

  it('handles mobile screen width (<= 750px)', () => {
    const originalWidth = window.innerWidth;
    Object.defineProperty(window, 'innerWidth', {
      configurable: true,
      writable: true,
      value: 600,
    });

    try {
      renderWithTheme(
        <SearchBar
          expanded={true}
          setExpanded={vi.fn()}
          onSearch={vi.fn()}
          placeholder="Mobile search"
        />,
      );

      const input = screen.getByPlaceholderText('Mobile search');
      expect(input).toBeDefined();
    } finally {
      Object.defineProperty(window, 'innerWidth', {
        configurable: true,
        writable: true,
        value: originalWidth,
      });
    }
  });

  it('cancels pending debounced onSearch on unmount', () => {
    const onSearch = vi.fn();
    const { unmount } = renderWithTheme(
      <SearchBar expanded={true} setExpanded={vi.fn()} onSearch={onSearch} placeholder="Search" />,
    );

    const input = screen.getByPlaceholderText('Search') as HTMLInputElement;
    fireEvent.change(input, { target: { value: 'test' } });

    unmount();

    act(() => {
      vi.advanceTimersByTime(600);
    });

    expect(onSearch).not.toHaveBeenCalled();
  });

  it('calls latest onSearch without duplicate invocations when onSearch prop changes', () => {
    const onSearch1 = vi.fn();
    const onSearch2 = vi.fn();

    const { rerender } = renderWithTheme(
      <SearchBar expanded={true} setExpanded={vi.fn()} onSearch={onSearch1} placeholder="Search" />,
    );

    const input = screen.getByPlaceholderText('Search') as HTMLInputElement;
    fireEvent.change(input, { target: { value: 'query' } });

    // Parent re-renders with a new onSearch function reference
    rerender(
      <SistentThemeProvider initialTheme="dark">
        <SearchBar
          expanded={true}
          setExpanded={vi.fn()}
          onSearch={onSearch2}
          placeholder="Search"
        />
      </SistentThemeProvider>,
    );

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(onSearch1).not.toHaveBeenCalled();
    expect(onSearch2).toHaveBeenCalledWith('query');
    expect(onSearch2).toHaveBeenCalledTimes(1);
  });
});
