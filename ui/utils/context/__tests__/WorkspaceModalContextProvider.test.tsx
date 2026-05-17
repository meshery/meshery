import React, { useContext } from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi, beforeEach } from 'vitest';

const getWorkspacesUnwrap = vi.fn();
const getWorkspaces = vi.fn(() => ({ unwrap: getWorkspacesUnwrap }));
const allOrganizations: Array<{ id: string; name: string }> = [];

vi.mock('@/rtk-query/user', () => ({
  useGetSelectedOrganization: () => ({ allOrganizations }),
}));

vi.mock('@/rtk-query/workspace', () => ({
  useLazyGetWorkspacesQuery: () => [getWorkspaces],
}));

import WorkspaceModalContextProvider, {
  WorkspaceModalContext,
} from '../WorkspaceModalContextProvider';

type Ctx = React.ContextType<typeof WorkspaceModalContext>;

const renderWithCtx = (renderProp: (ctx: Ctx) => React.ReactNode) => {
  const Consumer = () => {
    const ctx = useContext(WorkspaceModalContext);
    return <>{renderProp(ctx)}</>;
  };
  return render(
    <WorkspaceModalContextProvider>
      <Consumer />
    </WorkspaceModalContextProvider>,
  );
};

describe('WorkspaceModalContextProvider', () => {
  beforeEach(() => {
    getWorkspacesUnwrap.mockReset();
    getWorkspaces.mockClear();
    allOrganizations.length = 0;
    // Silence the explicit console.log statements in the source under test.
    vi.spyOn(console, 'log').mockImplementation(() => {});
  });

  it('initial state matches defaults', () => {
    renderWithCtx((ctx) => (
      <div>
        <span data-testid="open">{String(ctx.open)}</span>
        <span data-testid="selected">{ctx.selectedWorkspace.id}</span>
        <span data-testid="multi">{ctx.multiSelectedContent.length}</span>
        <span data-testid="createOpen">{String(ctx.createNewWorkspaceModalOpen)}</span>
        <span data-testid="loaded">{ctx.currentLoadedResource.id}</span>
      </div>
    ));

    expect(screen.getByTestId('open')).toHaveTextContent('false');
    expect(screen.getByTestId('selected')).toHaveTextContent('');
    expect(screen.getByTestId('multi')).toHaveTextContent('0');
    expect(screen.getByTestId('createOpen')).toHaveTextContent('false');
    expect(screen.getByTestId('loaded')).toHaveTextContent('');
  });

  it('openModal sets open to true, closeModal resets state', async () => {
    renderWithCtx((ctx) => (
      <div>
        <span data-testid="open">{String(ctx.open)}</span>
        <span data-testid="selected">{ctx.selectedWorkspace.id}</span>
        <span data-testid="multi">{ctx.multiSelectedContent.length}</span>
        <button type="button" onClick={ctx.openModal}>
          open
        </button>
        <button
          type="button"
          onClick={() => {
            ctx.setSelectedWorkspace({ id: 'ws-1', name: 'Workspace 1' });
            ctx.setMultiSelectedContent(['a', 'b', 'c'] as never);
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
    expect(screen.getByTestId('selected')).toHaveTextContent('ws-1');

    await user.click(screen.getByRole('button', { name: 'close' }));
    expect(screen.getByTestId('open')).toHaveTextContent('false');
    expect(screen.getByTestId('selected')).toHaveTextContent('');
    expect(screen.getByTestId('multi')).toHaveTextContent('0');
  });

  it('openModalWithDefault resets selectedWorkspace and opens the modal', async () => {
    renderWithCtx((ctx) => (
      <div>
        <span data-testid="open">{String(ctx.open)}</span>
        <span data-testid="selected">{ctx.selectedWorkspace.id}</span>
        <button type="button" onClick={() => ctx.setSelectedWorkspace({ id: 'ws-2', name: 'WS2' })}>
          seed
        </button>
        <button type="button" onClick={ctx.openModalWithDefault}>
          openDefault
        </button>
      </div>
    ));

    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: 'seed' }));
    expect(screen.getByTestId('selected')).toHaveTextContent('ws-2');

    await user.click(screen.getByRole('button', { name: 'openDefault' }));
    expect(screen.getByTestId('open')).toHaveTextContent('true');
    expect(screen.getByTestId('selected')).toHaveTextContent('');
  });

  it('setSelectedWorkspace also clears multiSelectedContent', async () => {
    renderWithCtx((ctx) => (
      <div>
        <span data-testid="selected">{ctx.selectedWorkspace.id}</span>
        <span data-testid="multi">{ctx.multiSelectedContent.length}</span>
        <button type="button" onClick={() => ctx.setMultiSelectedContent(['x'] as never)}>
          fillMulti
        </button>
        <button type="button" onClick={() => ctx.setSelectedWorkspace({ id: 'ws-3', name: 'WS3' })}>
          setSel
        </button>
      </div>
    ));

    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: 'fillMulti' }));
    expect(screen.getByTestId('multi')).toHaveTextContent('1');

    await user.click(screen.getByRole('button', { name: 'setSel' }));
    expect(screen.getByTestId('selected')).toHaveTextContent('ws-3');
    expect(screen.getByTestId('multi')).toHaveTextContent('0');
  });

  describe('onLoadResource', () => {
    it('clears workspace/org IDs when both ids are absent', async () => {
      renderWithCtx((ctx) => (
        <div>
          <span data-testid="loaded">{ctx.currentLoadedResource.id}</span>
          <button type="button" onClick={() => ctx.onLoadResource({ id: 'res-1' })}>
            load
          </button>
        </div>
      ));

      const user = userEvent.setup();
      await user.click(screen.getByRole('button', { name: 'load' }));
      await waitFor(() => expect(screen.getByTestId('loaded')).toHaveTextContent('res-1'));
      expect(getWorkspaces).not.toHaveBeenCalled();
    });

    it('looks up the matching organization and workspace', async () => {
      allOrganizations.push({ id: 'org-1', name: 'Org One' });
      getWorkspacesUnwrap.mockResolvedValue({
        workspaces: [{ id: 'ws-1', name: 'WS One' }],
      });

      const loaded: Array<{ workspace?: { id?: string; name?: string } }> = [];

      renderWithCtx((ctx) => {
        loaded.push(ctx.currentLoadedResource as never);
        return (
          <div>
            <span data-testid="orgName">{ctx.currentLoadedResource.org.name}</span>
            <span data-testid="wsName">{ctx.currentLoadedResource.workspace.name}</span>
            <button
              type="button"
              onClick={() =>
                ctx.onLoadResource({ id: 'res-1', workspaceId: 'ws-1', orgId: 'org-1' })
              }
            >
              load
            </button>
          </div>
        );
      });

      const user = userEvent.setup();
      await user.click(screen.getByRole('button', { name: 'load' }));

      await waitFor(() => {
        expect(getWorkspaces).toHaveBeenCalledWith({
          page: 0,
          pagesize: 'all',
          orgId: 'org-1',
        });
        expect(screen.getByTestId('orgName')).toHaveTextContent('Org One');
        expect(screen.getByTestId('wsName')).toHaveTextContent('WS One');
      });
    });

    it('falls back to the Private placeholders when the org is not found', async () => {
      // no orgs in allOrganizations
      renderWithCtx((ctx) => (
        <div>
          <span data-testid="orgName">{ctx.currentLoadedResource.org.name}</span>
          <span data-testid="wsName">{ctx.currentLoadedResource.workspace.name}</span>
          <button
            type="button"
            onClick={() => ctx.onLoadResource({ id: 'res-1', workspaceId: 'ws-1', orgId: 'org-x' })}
          >
            load
          </button>
        </div>
      ));

      const user = userEvent.setup();
      await user.click(screen.getByRole('button', { name: 'load' }));

      await waitFor(() => {
        expect(screen.getByTestId('orgName')).toHaveTextContent('Private Org');
        expect(screen.getByTestId('wsName')).toHaveTextContent('Private Workspace');
      });
      expect(getWorkspaces).not.toHaveBeenCalled();
    });

    it('handles unwrap failures gracefully by logging and not throwing', async () => {
      allOrganizations.push({ id: 'org-1', name: 'Org One' });
      getWorkspacesUnwrap.mockRejectedValueOnce(new Error('boom'));

      renderWithCtx((ctx) => (
        <button
          type="button"
          onClick={() => ctx.onLoadResource({ id: 'res-1', workspaceId: 'ws-1', orgId: 'org-1' })}
        >
          load
        </button>
      ));

      const user = userEvent.setup();
      await user.click(screen.getByRole('button', { name: 'load' }));
      // wait for the async error path to settle
      await waitFor(() => expect(getWorkspaces).toHaveBeenCalledTimes(1));
    });
  });

  it('setCreateNewWorkspaceModalOpen toggles the create modal', async () => {
    renderWithCtx((ctx) => (
      <div>
        <span data-testid="createOpen">{String(ctx.createNewWorkspaceModalOpen)}</span>
        <button type="button" onClick={() => ctx.setCreateNewWorkspaceModalOpen(true)}>
          openCreate
        </button>
      </div>
    ));

    await userEvent.click(screen.getByRole('button', { name: 'openCreate' }));
    expect(screen.getByTestId('createOpen')).toHaveTextContent('true');
  });
});
