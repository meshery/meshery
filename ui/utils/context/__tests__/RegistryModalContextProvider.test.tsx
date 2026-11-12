import React, { useContext } from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';
import RegistryModalContextProvider, {
  RegistryModalContext,
} from '../RegistryModalContextProvider';

type Ctx = React.ContextType<typeof RegistryModalContext>;

const renderWithContext = (renderProp: (ctx: Ctx) => React.ReactNode) => {
  const Consumer = () => {
    const ctx = useContext(RegistryModalContext);
    return <>{renderProp(ctx)}</>;
  };

  return render(
    <RegistryModalContextProvider>
      <Consumer />
    </RegistryModalContextProvider>,
  );
};

describe('RegistryModalContextProvider', () => {
  it('initialises with sensible defaults', () => {
    renderWithContext((ctx) => (
      <div>
        <span data-testid="open">{String(ctx.open)}</span>
        <span data-testid="view">{ctx.selectedView}</span>
        <span data-testid="search">{ctx.searchText}</span>
        <span data-testid="uuid">{ctx.selectedItemUUID}</span>
      </div>
    ));

    expect(screen.getByTestId('open')).toHaveTextContent('false');
    expect(screen.getByTestId('view')).toHaveTextContent('Models');
    expect(screen.getByTestId('search')).toHaveTextContent('');
    expect(screen.getByTestId('uuid')).toHaveTextContent('');
  });

  it('openModal flips open to true and closeModal resets state', async () => {
    renderWithContext((ctx) => (
      <div>
        <span data-testid="open">{String(ctx.open)}</span>
        <span data-testid="search">{ctx.searchText}</span>
        <span data-testid="uuid">{ctx.selectedItemUUID}</span>
        <button type="button" onClick={ctx.openModal}>
          open
        </button>
        <button
          type="button"
          onClick={() => {
            ctx.setSearchText('hello');
            ctx.setSelectedItemUUID('uuid-1');
          }}
        >
          seed
        </button>
        <button type="button" onClick={ctx.closeModal}>
          close
        </button>
      </div>
    ));

    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: 'open' }));
    expect(screen.getByTestId('open')).toHaveTextContent('true');

    await user.click(screen.getByRole('button', { name: 'seed' }));
    expect(screen.getByTestId('search')).toHaveTextContent('hello');
    expect(screen.getByTestId('uuid')).toHaveTextContent('uuid-1');

    await user.click(screen.getByRole('button', { name: 'close' }));
    expect(screen.getByTestId('open')).toHaveTextContent('false');
    expect(screen.getByTestId('search')).toHaveTextContent('');
    expect(screen.getByTestId('uuid')).toHaveTextContent('');
  });

  it('openModalWithParams seeds tab/search/selectedItemUUID when provided', async () => {
    renderWithContext((ctx) => (
      <div>
        <span data-testid="open">{String(ctx.open)}</span>
        <span data-testid="view">{ctx.selectedView}</span>
        <span data-testid="search">{ctx.searchText}</span>
        <span data-testid="uuid">{ctx.selectedItemUUID}</span>
        <button
          type="button"
          onClick={() =>
            ctx.openModalWithParams({
              tab: 'Components',
              searchText: 'istio',
              selectedItemUUID: 'abc',
            })
          }
        >
          openWithParams
        </button>
      </div>
    ));

    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: 'openWithParams' }));
    expect(screen.getByTestId('open')).toHaveTextContent('true');
    expect(screen.getByTestId('view')).toHaveTextContent('Components');
    expect(screen.getByTestId('search')).toHaveTextContent('istio');
    expect(screen.getByTestId('uuid')).toHaveTextContent('abc');
  });

  it('openModalWithParams accepts undefined params (defaults preserved)', async () => {
    renderWithContext((ctx) => (
      <div>
        <span data-testid="open">{String(ctx.open)}</span>
        <span data-testid="view">{ctx.selectedView}</span>
        <button type="button" onClick={() => ctx.openModalWithParams()}>
          openNoArgs
        </button>
      </div>
    ));
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: 'openNoArgs' }));
    expect(screen.getByTestId('open')).toHaveTextContent('true');
    expect(screen.getByTestId('view')).toHaveTextContent('Models');
  });

  it('setSelectedView updates the view', async () => {
    renderWithContext((ctx) => (
      <div>
        <span data-testid="view">{ctx.selectedView}</span>
        <button type="button" onClick={() => ctx.setSelectedView('Registrants')}>
          set
        </button>
      </div>
    ));
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: 'set' }));
    expect(screen.getByTestId('view')).toHaveTextContent('Registrants');
  });
});
