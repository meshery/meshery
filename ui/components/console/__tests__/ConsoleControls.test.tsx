import React from 'react';
import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import ConsoleControls from '../ConsoleControls';
import { ConsoleHeaderSlotContext } from '../header-slot';

const props = {
  containers: ['app', 'sidecar'],
  container: 'app',
  onContainerChange: vi.fn(),
  query: '',
  onQueryChange: vi.fn(),
  onSearch: vi.fn(),
};

/** Renders the controls with a header slot standing in for the panel's header bar. */
const renderWithSlot = (focused: boolean) => {
  const slot = document.createElement('div');
  slot.setAttribute('data-testid', 'header-slot');
  document.body.appendChild(slot);

  render(
    <ConsoleHeaderSlotContext.Provider value={slot}>
      <div data-testid="console-toolbar">
        <ConsoleControls {...props} focused={focused} />
      </div>
    </ConsoleHeaderSlotContext.Provider>,
  );

  return slot;
};

describe('ConsoleControls', () => {
  it('projects the focused console’s controls into the panel header', () => {
    const slot = renderWithSlot(true);

    const search = screen.getByPlaceholderText('Search');
    expect(slot).toContainElement(search);
    // Rendered from inside the console, but not *in* the console's own toolbar.
    expect(screen.getByTestId('console-toolbar')).not.toContainElement(search);
  });

  it('renders nothing for a console that is open but not on show', () => {
    // The panel has one header; three backgrounded consoles must not stack three
    // container selects into it.
    renderWithSlot(false);

    expect(screen.queryByPlaceholderText('Search')).not.toBeInTheDocument();
  });

  it('falls back to rendering inline when there is no panel header', () => {
    // An extension embedding a console on its own still needs to switch container.
    render(
      <div data-testid="console-toolbar">
        <ConsoleControls {...props} focused />
      </div>,
    );

    expect(screen.getByTestId('console-toolbar')).toContainElement(
      screen.getByPlaceholderText('Search'),
    );
  });
});
