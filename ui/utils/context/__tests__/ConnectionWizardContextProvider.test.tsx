import React, { useContext } from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/components/connections/ConnectionWizardModal', () => ({
  default: () => null,
}));

import ConnectionWizardContextProvider, {
  ConnectionWizardContext,
} from '../ConnectionWizardContextProvider';

const Probe = () => {
  const ctx = useContext(ConnectionWizardContext);
  return (
    <div>
      <span data-testid="open">{String(ctx.open)}</span>
      <span data-testid="kind">{ctx.presetKind ?? 'none'}</span>
      <span data-testid="skip">{String(ctx.skipKindSelection)}</span>
      <button
        type="button"
        onClick={() => ctx.openCreateConnection({ kind: 'kubernetes', skipKindSelection: true })}
      >
        openK8s
      </button>
      <button type="button" onClick={() => ctx.openCreateConnection()}>
        openPlain
      </button>
      <button type="button" onClick={() => ctx.closeCreateConnection()}>
        close
      </button>
    </div>
  );
};

describe('ConnectionWizardContextProvider', () => {
  it('opens with kind preset and clears on close', async () => {
    const user = userEvent.setup();
    render(
      <ConnectionWizardContextProvider>
        <Probe />
      </ConnectionWizardContextProvider>,
    );

    expect(screen.getByTestId('open')).toHaveTextContent('false');

    await user.click(screen.getByRole('button', { name: 'openK8s' }));
    expect(screen.getByTestId('open')).toHaveTextContent('true');
    expect(screen.getByTestId('kind')).toHaveTextContent('kubernetes');
    expect(screen.getByTestId('skip')).toHaveTextContent('true');

    await user.click(screen.getByRole('button', { name: 'close' }));
    expect(screen.getByTestId('open')).toHaveTextContent('false');
    expect(screen.getByTestId('kind')).toHaveTextContent('none');
    expect(screen.getByTestId('skip')).toHaveTextContent('false');
  });

  it('opens without preset for the toolbar create path', async () => {
    const user = userEvent.setup();
    render(
      <ConnectionWizardContextProvider>
        <Probe />
      </ConnectionWizardContextProvider>,
    );

    await user.click(screen.getByRole('button', { name: 'openPlain' }));
    expect(screen.getByTestId('open')).toHaveTextContent('true');
    expect(screen.getByTestId('kind')).toHaveTextContent('none');
    expect(screen.getByTestId('skip')).toHaveTextContent('false');
  });
});
