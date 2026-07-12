import React from 'react';
import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import SessionControls from '../SessionControls';
import { SessionHeaderSlotContext } from '../header-slot';

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
    <SessionHeaderSlotContext.Provider value={slot}>
      <div data-testid="session-toolbar">
        <SessionControls {...props} focused={focused} />
      </div>
    </SessionHeaderSlotContext.Provider>,
  );

  return slot;
};

describe('SessionControls', () => {
  it('projects the focused session’s controls into the panel header', () => {
    const slot = renderWithSlot(true);

    const search = screen.getByPlaceholderText('Search');
    expect(slot).toContainElement(search);
    // Rendered from inside the session, but not *in* the session's own toolbar.
    expect(screen.getByTestId('session-toolbar')).not.toContainElement(search);
  });

  it('renders nothing for a session that is open but not on show', () => {
    // The panel has one header; three backgrounded sessions must not stack three
    // container selects into it.
    renderWithSlot(false);

    expect(screen.queryByPlaceholderText('Search')).not.toBeInTheDocument();
  });

  it('falls back to rendering inline when there is no panel header', () => {
    // An extension embedding a session on its own still needs to switch container.
    render(
      <div data-testid="session-toolbar">
        <SessionControls {...props} focused />
      </div>,
    );

    expect(screen.getByTestId('session-toolbar')).toContainElement(
      screen.getByPlaceholderText('Search'),
    );
  });
});
