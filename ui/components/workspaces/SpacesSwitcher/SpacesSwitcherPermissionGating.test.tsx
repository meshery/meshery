import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Keys } from '@meshery/schemas/permissions';

let mockHasPermission = vi.fn();
let mockUseFetchViewsQuery = vi.fn();
let mockUseGetUserDesignsQuery = vi.fn();
let mockUseGetViewsOfWorkspaceQuery = vi.fn();
let mockUseGetDesignsOfWorkspaceQuery = vi.fn();
let mockIsInOperatorMode = vi.fn();

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@sistent/sistent')>();
  return {
    ...actual,
    useHasPermission: (key: any) => mockHasPermission(key),
    useTheme: () => ({
      palette: {
        text: { default: '#000' },
        icon: { default: '#000', brand: '#000' },
        background: { default: '#fff' },
        error: { dark: '#d32f2f' },
        common: { white: '#fff' },
      },
    }),
    Select: ({ children, value, onChange }: any) => (
      <select data-testid="type-select" value={value} onChange={onChange}>
        {children}
      </select>
    ),
    MenuItem: ({ children, value }: any) => (
      <option data-testid={`menu-item-${value}`} value={value}>
        {children}
      </option>
    ),
    StyledSearchBar: () => <input data-testid="search-bar" />,
    PromptComponent: React.forwardRef(() => null),
    WorkspaceContentMoveModal: () => null,
  };
});

vi.mock('@/utils/utils', () => ({
  isInOperatorMode: () => mockIsInOperatorMode(),
  JsonParse: (s: string) => JSON.parse(s),
}));

vi.mock('@/rtk-query/view', () => ({
  useFetchViewsQuery: (...args: any[]) => mockUseFetchViewsQuery(...args),
  useDeleteViewMutation: () => [vi.fn()],
  useUpdateViewVisibilityMutation: () => [vi.fn()],
}));

vi.mock('@/rtk-query/design', () => ({
  useGetUserDesignsQuery: (...args: any[]) => mockUseGetUserDesignsQuery(...args),
  useDeletePatternFileMutation: () => [vi.fn()],
}));

vi.mock('@/rtk-query/workspace', () => ({
  useGetViewsOfWorkspaceQuery: (...args: any[]) => mockUseGetViewsOfWorkspaceQuery(...args),
  useGetDesignsOfWorkspaceQuery: (...args: any[]) => mockUseGetDesignsOfWorkspaceQuery(...args),
  useGetWorkspacesQuery: () => ({ data: { workspaces: [] }, isLoading: false }),
  useAssignDesignToWorkspaceMutation: () => [vi.fn()],
  useAssignViewToWorkspaceMutation: () => [vi.fn()],
}));

vi.mock('react-redux', () => ({
  useSelector: () => ({ id: 'org-1' }),
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('./components', () => ({
  ImportButton: () => null,
  SortBySelect: () => null,
  TableListHeader: () => null,
  UserSearchAutoComplete: () => null,
  VisibilitySelect: () => null,
  MultiContentSelectToolbar: () => null,
}));

vi.mock('./MainDesignsContent', () => ({
  default: () => <div data-testid="main-designs-content" />,
}));

vi.mock('./MainViewsContent', () => ({
  default: () => <div data-testid="main-views-content" />,
}));

vi.mock('../../designs/export/ExportDesignModal', () => ({
  default: () => null,
}));

import RecentContent from './RecentContent';
import SharedContent from './SharedContent';
import WorkspaceContent from './WorkspaceContent';

describe('SpacesSwitcher View Permission Gating', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseFetchViewsQuery.mockReturnValue({
      data: { views: [], totalCount: 0 },
      isLoading: false,
      isFetching: false,
      refetch: vi.fn(),
    });
    mockUseGetUserDesignsQuery.mockReturnValue({
      data: { patterns: [], totalCount: 0 },
      isLoading: false,
      isFetching: false,
      refetch: vi.fn(),
    });
    mockUseGetViewsOfWorkspaceQuery.mockReturnValue({
      data: { views: [], totalCount: 0 },
      isLoading: false,
      isFetching: false,
      refetch: vi.fn(),
    });
    mockUseGetDesignsOfWorkspaceQuery.mockReturnValue({
      data: { designs: [], totalCount: 0 },
      isLoading: false,
      isFetching: false,
      refetch: vi.fn(),
    });
  });

  describe('RecentContent', () => {
    it('defaults to DESIGN and skips views query when in operator mode without view permission', () => {
      mockIsInOperatorMode.mockReturnValue(true);
      mockHasPermission.mockImplementation((key) => {
        if (key?.id === Keys.KanvasViewViews.id || key === Keys.KanvasViewViews) {
          return false;
        }
        return true;
      });

      render(<RecentContent />);

      // Select value should be 'design', not 'view'
      const select = screen.getByTestId('type-select') as HTMLSelectElement;
      expect(select.value).toBe('design');
      expect(screen.queryByTestId('menu-item-view')).not.toBeInTheDocument();

      // View query MUST be skipped
      const fetchViewsCall = mockUseFetchViewsQuery.mock.calls.at(-1);
      expect(fetchViewsCall).toBeDefined();
      expect(fetchViewsCall[1]).toEqual(expect.objectContaining({ skip: true }));

      // MainViewsContent MUST NOT be mounted; MainDesignsContent MUST be mounted
      expect(screen.queryByTestId('main-views-content')).not.toBeInTheDocument();
      expect(screen.getByTestId('main-designs-content')).toBeInTheDocument();
    });

    it('defaults to VIEW and runs views query when in operator mode with view permission', () => {
      mockIsInOperatorMode.mockReturnValue(true);
      mockHasPermission.mockReturnValue(true);

      render(<RecentContent />);

      const select = screen.getByTestId('type-select') as HTMLSelectElement;
      expect(select.value).toBe('view');
      expect(screen.getByTestId('menu-item-view')).toBeInTheDocument();

      const fetchViewsCall = mockUseFetchViewsQuery.mock.calls.at(-1);
      expect(fetchViewsCall).toBeDefined();
      expect(fetchViewsCall[1]).toEqual(expect.objectContaining({ skip: false }));

      expect(screen.getByTestId('main-views-content')).toBeInTheDocument();
      expect(screen.queryByTestId('main-designs-content')).not.toBeInTheDocument();
    });
  });

  describe('SharedContent', () => {
    it('defaults to DESIGN and skips views query when in operator mode without view permission', () => {
      mockIsInOperatorMode.mockReturnValue(true);
      mockHasPermission.mockImplementation((key) => {
        if (key?.id === Keys.KanvasViewViews.id || key === Keys.KanvasViewViews) {
          return false;
        }
        return true;
      });

      render(<SharedContent />);

      const select = screen.getByTestId('type-select') as HTMLSelectElement;
      expect(select.value).toBe('design');
      expect(screen.queryByTestId('menu-item-view')).not.toBeInTheDocument();

      const fetchViewsCall = mockUseFetchViewsQuery.mock.calls.at(-1);
      expect(fetchViewsCall).toBeDefined();
      expect(fetchViewsCall[1]).toEqual(expect.objectContaining({ skip: true }));

      expect(screen.queryByTestId('main-views-content')).not.toBeInTheDocument();
      expect(screen.getByTestId('main-designs-content')).toBeInTheDocument();
    });

    it('defaults to VIEW and runs views query when in operator mode with view permission', () => {
      mockIsInOperatorMode.mockReturnValue(true);
      mockHasPermission.mockReturnValue(true);

      render(<SharedContent />);

      const select = screen.getByTestId('type-select') as HTMLSelectElement;
      expect(select.value).toBe('view');
      expect(screen.getByTestId('menu-item-view')).toBeInTheDocument();

      const fetchViewsCall = mockUseFetchViewsQuery.mock.calls.at(-1);
      expect(fetchViewsCall).toBeDefined();
      expect(fetchViewsCall[1]).toEqual(expect.objectContaining({ skip: false }));

      expect(screen.getByTestId('main-views-content')).toBeInTheDocument();
      expect(screen.queryByTestId('main-designs-content')).not.toBeInTheDocument();
    });
  });

  describe('WorkspaceContent', () => {
    const mockWorkspace = { id: 'ws-123', name: 'Test Workspace' };

    it('defaults to DESIGN and skips views query when in operator mode without view permission', () => {
      mockIsInOperatorMode.mockReturnValue(true);
      mockHasPermission.mockImplementation((key) => {
        if (key?.id === Keys.KanvasViewViews.id || key === Keys.KanvasViewViews) {
          return false;
        }
        return true;
      });

      render(<WorkspaceContent workspace={mockWorkspace} />);

      const select = screen.getByTestId('type-select') as HTMLSelectElement;
      expect(select.value).toBe('design');
      expect(screen.queryByTestId('menu-item-view')).not.toBeInTheDocument();

      const getViewsCall = mockUseGetViewsOfWorkspaceQuery.mock.calls.at(-1);
      expect(getViewsCall).toBeDefined();
      expect(getViewsCall[1]).toEqual(expect.objectContaining({ skip: true }));

      expect(screen.queryByTestId('main-views-content')).not.toBeInTheDocument();
      expect(screen.getByTestId('main-designs-content')).toBeInTheDocument();
    });

    it('defaults to VIEW and runs views query when in operator mode with view permission', () => {
      mockIsInOperatorMode.mockReturnValue(true);
      mockHasPermission.mockReturnValue(true);

      render(<WorkspaceContent workspace={mockWorkspace} />);

      const select = screen.getByTestId('type-select') as HTMLSelectElement;
      expect(select.value).toBe('view');
      expect(screen.getByTestId('menu-item-view')).toBeInTheDocument();

      const getViewsCall = mockUseGetViewsOfWorkspaceQuery.mock.calls.at(-1);
      expect(getViewsCall).toBeDefined();
      expect(getViewsCall[1]).toEqual(expect.objectContaining({ skip: false }));

      expect(screen.getByTestId('main-views-content')).toBeInTheDocument();
      expect(screen.queryByTestId('main-designs-content')).not.toBeInTheDocument();
    });

    it('safely resets to DESIGN when permission is dynamically revoked after mount', () => {
      mockIsInOperatorMode.mockReturnValue(true);
      mockHasPermission.mockReturnValue(true);

      const { rerender } = render(<WorkspaceContent workspace={mockWorkspace} />);

      expect((screen.getByTestId('type-select') as HTMLSelectElement).value).toBe('view');
      expect(screen.getByTestId('menu-item-view')).toBeInTheDocument();

      // Revoke permission dynamically
      mockHasPermission.mockReturnValue(false);
      rerender(<WorkspaceContent workspace={mockWorkspace} />);

      const select = screen.getByTestId('type-select') as HTMLSelectElement;
      expect(select.value).toBe('design');
      expect(screen.queryByTestId('menu-item-view')).not.toBeInTheDocument();

      const getViewsCall = mockUseGetViewsOfWorkspaceQuery.mock.calls.at(-1);
      expect(getViewsCall[1]).toEqual(expect.objectContaining({ skip: true }));
      expect(screen.queryByTestId('main-views-content')).not.toBeInTheDocument();
      expect(screen.getByTestId('main-designs-content')).toBeInTheDocument();
    });
  });
});
