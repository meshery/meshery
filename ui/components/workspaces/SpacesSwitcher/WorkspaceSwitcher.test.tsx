import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { EVENT_TYPES } from 'lib/event-types';

const mockNotify = vi.fn();
const mockUpdateSelectedWorkspace = vi.fn();
const mockSetSelectedWorkspace = vi.fn();
const mockOpenModal = vi.fn();
const mockOnMobileSelect = vi.fn();

let mockSelectedOrganization: { id: string; name: string } | undefined = {
  id: 'org-1',
  name: 'Default Org',
};

const mockWorkspaces = [
  { id: 'ws-1', name: 'Workspace One' },
  { id: 'ws-2', name: 'Workspace Two' },
];

vi.mock('@/utils/context/WorkspaceModalContextProvider', async () => {
  const ReactActual = await vi.importActual<typeof import('react')>('react');
  return {
    WorkspaceModalContext: ReactActual.createContext({
      setSelectedWorkspace: vi.fn(),
      openModal: vi.fn(),
      setCreateNewWorkspaceModalOpen: vi.fn(),
      currentLoadedResource: null,
    }),
  };
});

import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: mockNotify }),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetSelectedOrganization: () => ({ selectedOrganization: mockSelectedOrganization }),
  useGetSelectedWorkspace: () => ({
    selectedWorkspace: mockWorkspaces[0],
    allWorkspaces: mockWorkspaces,
    error: null,
    isLoading: false,
  }),
  useUpdateSelectedWorkspaceMutation: () => [mockUpdateSelectedWorkspace, { isLoading: false }],
}));

vi.mock('css/icons.styles', () => ({
  iconMedium: {},
}));

vi.mock('@sistent/sistent', () => ({
  FormControl: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  FormControlLabel: ({ control }: { control: React.ReactNode }) => <div>{control}</div>,
  FormGroup: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  Grid2: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  styled: (Comp: React.ComponentType<any>) => () => Comp,
  MenuItem: ({ children, onClick }: { children: React.ReactNode; onClick?: () => void }) => (
    <div onClick={onClick}>{children}</div>
  ),
  CircularProgress: () => <div data-testid="loading" />,
  WorkspaceIcon: () => <span />,
  Divider: () => <hr />,
  Button: ({ children, onClick }: { children: React.ReactNode; onClick?: (e: any) => void }) => (
    <button type="button" onClick={onClick}>
      {children}
    </button>
  ),
  Box: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  NoSsr: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  useMediaQuery: () => false,
  useTheme: () => ({
    palette: { icon: { default: '#000', secondary: '#555' }, action: { selected: '#eee' } },
    zIndex: { modal: 1300 },
  }),
}));

vi.mock('./SpaceSwitcher', () => ({
  StyledSelect: ({
    value,
    onChange,
    children,
  }: {
    value: string;
    onChange: (e: { target: { value: string } }) => void;
    children: React.ReactNode;
  }) => (
    <div>
      <select
        data-testid="workspace-select"
        value={value}
        onChange={(e) => onChange({ target: { value: e.target.value } })}
      >
        {mockWorkspaces.map((w) => (
          <option key={w.id} value={w.id}>
            {w.name}
          </option>
        ))}
      </select>
      {children}
    </div>
  ),
}));

import WorkspaceSwitcher from './WorkspaceSwitcher';

const renderSwitcher = () =>
  render(
    <WorkspaceModalContext.Provider
      value={
        {
          setSelectedWorkspace: mockSetSelectedWorkspace,
          openModal: mockOpenModal,
          setCreateNewWorkspaceModalOpen: vi.fn(),
          currentLoadedResource: null,
        } as any
      }
    >
      <WorkspaceSwitcher
        open={true}
        fromMobileView={false}
        expanded={false}
        onExpandedChange={vi.fn()}
        onMobileSelect={mockOnMobileSelect}
      />
    </WorkspaceModalContext.Provider>,
  );

describe('WorkspaceSwitcher', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockSelectedOrganization = { id: 'org-1', name: 'Default Org' };
  });

  it('does not throw or call mutation when selectedOrganization is undefined', async () => {
    mockSelectedOrganization = undefined;
    renderSwitcher();

    const select = screen.getByTestId('workspace-select');
    expect(() => {
      fireEvent.change(select, { target: { value: 'ws-2' } });
    }).not.toThrow();

    await waitFor(() => {
      expect(mockUpdateSelectedWorkspace).not.toHaveBeenCalled();
      expect(mockSetSelectedWorkspace).not.toHaveBeenCalled();
      expect(mockOpenModal).not.toHaveBeenCalled();
    });
  });

  it('reverts workspace selection, skips opening modal, and notifies error when updateSelectedWorkspace rejects', async () => {
    mockUpdateSelectedWorkspace.mockRejectedValueOnce({
      data: { message: 'Forbidden' },
    });
    renderSwitcher();

    const select = screen.getByTestId('workspace-select');
    fireEvent.change(select, { target: { value: 'ws-2' } });

    await waitFor(() => {
      expect(mockUpdateSelectedWorkspace).toHaveBeenCalledWith('org-1', 'ws-2');
      expect(mockSetSelectedWorkspace).toHaveBeenNthCalledWith(1, mockWorkspaces[1]);
      expect(mockSetSelectedWorkspace).toHaveBeenNthCalledWith(2, mockWorkspaces[0]);
      expect(mockOpenModal).not.toHaveBeenCalled();
      expect(mockOnMobileSelect).not.toHaveBeenCalled();
      expect(mockNotify).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'Failed to switch workspace',
          event_type: EVENT_TYPES.ERROR,
          details: 'Forbidden',
        }),
      );
    });
  });

  it('persists workspace selection and opens modal when updateSelectedWorkspace succeeds', async () => {
    mockUpdateSelectedWorkspace.mockResolvedValueOnce({ ok: true });
    renderSwitcher();

    const select = screen.getByTestId('workspace-select');
    fireEvent.change(select, { target: { value: 'ws-2' } });

    await waitFor(() => {
      expect(mockUpdateSelectedWorkspace).toHaveBeenCalledWith('org-1', 'ws-2');
      expect(mockSetSelectedWorkspace).toHaveBeenCalledTimes(1);
      expect(mockSetSelectedWorkspace).toHaveBeenCalledWith(mockWorkspaces[1]);
      expect(mockOpenModal).toHaveBeenCalledWith(true);
      expect(mockNotify).not.toHaveBeenCalled();
    });
  });
});
