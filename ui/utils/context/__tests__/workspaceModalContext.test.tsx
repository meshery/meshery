import React, { useContext } from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { WorkspaceModalContext, workspaceModalContextDefault } from '../workspaceModalContext';
import WorkspaceModalContextProvider from '../WorkspaceModalContextProvider';

vi.mock('@/rtk-query/user', () => ({
  useGetSelectedOrganization: () => ({ allOrganizations: [] }),
}));

vi.mock('@/rtk-query/workspace', () => ({
  useLazyGetWorkspacesQuery: () => [vi.fn()],
}));

describe('workspaceModalContext', () => {
  it('resolves to one context instance across import paths (prevents duplicate createContext)', async () => {
    const fromRelative = WorkspaceModalContext;
    const { WorkspaceModalContext: fromAlias } =
      await import('@/utils/context/workspaceModalContext');
    expect(fromAlias).toBe(fromRelative);
  });

  it('provider wires openModal through the shared context object (not the default no-op)', () => {
    let openModal: (() => void) | undefined;

    const Probe = () => {
      openModal = useContext(WorkspaceModalContext).openModal;
      return null;
    };

    render(
      <WorkspaceModalContextProvider>
        <Probe />
      </WorkspaceModalContextProvider>,
    );

    expect(openModal).toBeDefined();
    expect(openModal).not.toBe(workspaceModalContextDefault.openModal);
  });
});
