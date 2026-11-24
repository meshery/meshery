import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';
import { NavItem, WorkspacesSection } from './WorkspaceFormModalSections';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';

let mockIsTruncated = false;

vi.mock('@/utils/hooks', () => ({
  useIsTextTruncated: () => mockIsTruncated,
}));

vi.mock('@/utils/context/WorkspaceModalContextProvider', () => ({
  WorkspaceModalContext: React.createContext({
    setMultiSelectedContent: vi.fn(),
  }),
}));

vi.mock('css/icons.styles', () => ({
  iconSmall: { height: 20, width: 20 },
}));

vi.mock('@sistent/sistent', () => ({
  CustomTooltip: ({ title, children }: any) => (
    <div data-testid="custom-tooltip" data-title={title}>
      {children}
    </div>
  ),
  ListItem: ({ children, sx }: any) => <li style={sx}>{children}</li>,
  ListItemButton: ({ children, onClick, selected, sx }: any) => (
    <button type="button" onClick={onClick} data-selected={selected} style={sx}>
      {children}
    </button>
  ),
  ListItemIcon: ({ children }: any) => <span>{children}</span>,
  ListItemText: ({ primary }: any) => <span>{primary}</span>,
  WorkspaceIcon: () => <svg data-testid="workspace-icon" />,
}));

vi.mock('@/theme', () => ({
  useTheme: () => ({
    palette: {
      icon: { default: '#000' },
      background: { secondary: '#f5f5f5' },
    },
  }),
}));

describe('WorkspaceFormModalSections - CustomTooltip Conditional Rendering', () => {
  const renderWithContext = (ui: React.ReactElement) =>
    render(
      <WorkspaceModalContext.Provider value={{ setMultiSelectedContent: vi.fn() } as any}>
        {ui}
      </WorkspaceModalContext.Provider>,
    );

  beforeEach(() => {
    vi.clearAllMocks();
    mockIsTruncated = false;
  });

  describe('NavItem', () => {
    const item = {
      id: 'recents',
      label: 'Recents',
      icon: <span data-testid="recents-icon" />,
      content: <div>Recents Content</div>,
    };

    it('wraps with CustomTooltip when drawer is collapsed (open === false)', () => {
      renderWithContext(<NavItem item={item} open={false} selectedId="other" onSelect={vi.fn()} />);

      const tooltip = screen.getByTestId('custom-tooltip');
      expect(tooltip).toBeInTheDocument();
      expect(tooltip).toHaveAttribute('data-title', 'Recents');
    });

    it('renders without CustomTooltip when drawer is expanded (open === true)', () => {
      renderWithContext(<NavItem item={item} open={true} selectedId="other" onSelect={vi.fn()} />);

      expect(screen.queryByTestId('custom-tooltip')).not.toBeInTheDocument();
      expect(screen.getByText('Recents')).toBeInTheDocument();
    });
  });

  describe('WorkspacesSection', () => {
    const workspacesData = {
      workspaces: [
        { id: 'ws-1', name: 'Alpha Workspace' },
        { id: 'ws-2', name: 'Beta Workspace' },
      ],
    };

    it('renders All Workspaces with tooltip when drawer is collapsed (open === false)', () => {
      renderWithContext(
        <WorkspacesSection
          open={false}
          selectedId=""
          onSelect={vi.fn()}
          workspacesData={workspacesData}
          isLoading={false}
        />,
      );

      const tooltips = screen.getAllByTestId('custom-tooltip');
      expect(tooltips[0]).toHaveAttribute('data-title', 'All Workspaces');
    });

    it('renders All Workspaces without tooltip when drawer is expanded (open === true)', () => {
      renderWithContext(
        <WorkspacesSection
          open={true}
          selectedId=""
          onSelect={vi.fn()}
          workspacesData={{ workspaces: [] }}
          isLoading={false}
        />,
      );

      expect(screen.queryByTestId('custom-tooltip')).not.toBeInTheDocument();
      expect(screen.getByText('All Workspaces')).toBeInTheDocument();
    });

    it('wraps workspace items with CustomTooltip when drawer is collapsed (open === false)', () => {
      renderWithContext(
        <WorkspacesSection
          open={false}
          selectedId=""
          onSelect={vi.fn()}
          workspacesData={workspacesData}
          isLoading={false}
        />,
      );

      const tooltips = screen.getAllByTestId('custom-tooltip');
      expect(tooltips).toHaveLength(3); // All Workspaces + 2 workspace items
      expect(tooltips[1]).toHaveAttribute('data-title', 'Alpha Workspace');
      expect(tooltips[2]).toHaveAttribute('data-title', 'Beta Workspace');
    });

    it('renders workspace items without CustomTooltip when drawer is open and text is not truncated', () => {
      mockIsTruncated = false;
      renderWithContext(
        <WorkspacesSection
          open={true}
          selectedId=""
          onSelect={vi.fn()}
          workspacesData={workspacesData}
          isLoading={false}
        />,
      );

      expect(screen.queryByTestId('custom-tooltip')).not.toBeInTheDocument();
      expect(screen.getByText('Alpha Workspace')).toBeInTheDocument();
      expect(screen.getByText('Beta Workspace')).toBeInTheDocument();
    });

    it('wraps workspace items with CustomTooltip when drawer is open and text is truncated', () => {
      mockIsTruncated = true;
      renderWithContext(
        <WorkspacesSection
          open={true}
          selectedId=""
          onSelect={vi.fn()}
          workspacesData={workspacesData}
          isLoading={false}
        />,
      );

      const tooltips = screen.getAllByTestId('custom-tooltip');
      expect(tooltips).toHaveLength(2); // Only the 2 workspace items
      expect(tooltips[0]).toHaveAttribute('data-title', 'Alpha Workspace');
      expect(tooltips[1]).toHaveAttribute('data-title', 'Beta Workspace');
    });
  });
});
