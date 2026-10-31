import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/components/command-palette/CommandPalette', () => ({
  default: () => null,
}));

import CommandPaletteContextProvider, {
  CommandPaletteHost,
  useCommandPalette,
} from '../CommandPaletteContextProvider';

const Probe = () => {
  const { open, openPalette, closePalette, togglePalette } = useCommandPalette();
  return (
    <div>
      <span data-testid="open">{String(open)}</span>
      <button type="button" onClick={openPalette}>
        open
      </button>
      <button type="button" onClick={closePalette}>
        close
      </button>
      <button type="button" onClick={togglePalette}>
        toggle
      </button>
      <input data-testid="search-input" />
      <textarea data-testid="search-textarea" />
      <div data-testid="editable-element" contentEditable tabIndex={0} />
    </div>
  );
};

describe('CommandPaletteContextProvider and CommandPaletteHost', () => {
  it('opens on Ctrl+K', () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    expect(screen.getByTestId('open')).toHaveTextContent('false');
    fireEvent.keyDown(window, { key: 'k', ctrlKey: true });
    expect(screen.getByTestId('open')).toHaveTextContent('true');
  });

  it('opens on Meta+K', () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    expect(screen.getByTestId('open')).toHaveTextContent('false');
    fireEvent.keyDown(window, { key: 'k', metaKey: true });
    expect(screen.getByTestId('open')).toHaveTextContent('true');
  });

  it('closes on second shortcut press', () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    expect(screen.getByTestId('open')).toHaveTextContent('false');
    fireEvent.keyDown(window, { key: 'k', ctrlKey: true });
    expect(screen.getByTestId('open')).toHaveTextContent('true');

    fireEvent.keyDown(window, { key: 'k', ctrlKey: true });
    expect(screen.getByTestId('open')).toHaveTextContent('false');
  });

  it('ignores Ctrl+Shift+K and Ctrl+Alt+K', () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    expect(screen.getByTestId('open')).toHaveTextContent('false');

    // Ctrl+Shift+K
    fireEvent.keyDown(window, { key: 'k', ctrlKey: true, shiftKey: true });
    expect(screen.getByTestId('open')).toHaveTextContent('false');

    // Ctrl+Alt+K
    fireEvent.keyDown(window, { key: 'k', ctrlKey: true, altKey: true });
    expect(screen.getByTestId('open')).toHaveTextContent('false');
  });

  it('ignores isComposing events (IME composition)', () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    expect(screen.getByTestId('open')).toHaveTextContent('false');
    fireEvent.keyDown(window, { key: 'k', ctrlKey: true, isComposing: true });
    expect(screen.getByTestId('open')).toHaveTextContent('false');
  });

  it('ignores repeated keydown events (event.repeat)', () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    expect(screen.getByTestId('open')).toHaveTextContent('false');
    fireEvent.keyDown(window, { key: 'k', ctrlKey: true, repeat: true });
    expect(screen.getByTestId('open')).toHaveTextContent('false');
  });

  it('ignores defaultPrevented events', () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    expect(screen.getByTestId('open')).toHaveTextContent('false');

    const event = new KeyboardEvent('keydown', {
      key: 'k',
      ctrlKey: true,
      cancelable: true,
      bubbles: true,
    });
    event.preventDefault(); // marks defaultPrevented = true
    fireEvent(window, event);

    expect(screen.getByTestId('open')).toHaveTextContent('false');
  });

  it('removes listener when CommandPaletteHost unmounts', () => {
    const Harness = ({ mountHost }: { mountHost: boolean }) => (
      <CommandPaletteContextProvider>
        {mountHost && <CommandPaletteHost />}
        <Probe />
      </CommandPaletteContextProvider>
    );

    const { rerender } = render(<Harness mountHost={true} />);
    expect(screen.getByTestId('open')).toHaveTextContent('false');

    // Unmount just the host
    rerender(<Harness mountHost={false} />);

    fireEvent.keyDown(window, { key: 'k', ctrlKey: true });
    // Remains false because host listener was removed
    expect(screen.getByTestId('open')).toHaveTextContent('false');
  });

  it('works when an <input> has focus', async () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    const input = screen.getByTestId('search-input');
    input.focus();
    expect(document.activeElement).toBe(input);

    expect(screen.getByTestId('open')).toHaveTextContent('false');
    fireEvent.keyDown(input, { key: 'k', ctrlKey: true });
    expect(screen.getByTestId('open')).toHaveTextContent('true');
  });

  it('works when a <textarea> has focus', async () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    const textarea = screen.getByTestId('search-textarea');
    textarea.focus();
    expect(document.activeElement).toBe(textarea);

    expect(screen.getByTestId('open')).toHaveTextContent('false');
    fireEvent.keyDown(textarea, { key: 'k', ctrlKey: true });
    expect(screen.getByTestId('open')).toHaveTextContent('true');
  });

  it('works when a contentEditable element has focus', async () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    const editable = screen.getByTestId('editable-element');
    editable.focus();
    expect(document.activeElement).toBe(editable);

    expect(screen.getByTestId('open')).toHaveTextContent('false');
    fireEvent.keyDown(editable, { key: 'k', ctrlKey: true });
    expect(screen.getByTestId('open')).toHaveTextContent('true');
  });

  it('calls preventDefault on a handled shortcut', () => {
    render(
      <CommandPaletteContextProvider>
        <CommandPaletteHost />
        <Probe />
      </CommandPaletteContextProvider>,
    );

    const event = new KeyboardEvent('keydown', {
      key: 'k',
      ctrlKey: true,
      cancelable: true,
      bubbles: true,
    });
    const preventDefaultSpy = vi.spyOn(event, 'preventDefault');

    fireEvent(window, event);

    expect(preventDefaultSpy).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('open')).toHaveTextContent('true');
  });

  it('supports programmatic openPalette, closePalette, and togglePalette', async () => {
    const user = userEvent.setup();
    render(
      <CommandPaletteContextProvider>
        <Probe />
      </CommandPaletteContextProvider>,
    );

    expect(screen.getByTestId('open')).toHaveTextContent('false');

    await user.click(screen.getByRole('button', { name: 'open' }));
    expect(screen.getByTestId('open')).toHaveTextContent('true');

    await user.click(screen.getByRole('button', { name: 'close' }));
    expect(screen.getByTestId('open')).toHaveTextContent('false');

    await user.click(screen.getByRole('button', { name: 'toggle' }));
    expect(screen.getByTestId('open')).toHaveTextContent('true');

    await user.click(screen.getByRole('button', { name: 'toggle' }));
    expect(screen.getByTestId('open')).toHaveTextContent('false');
  });
});
