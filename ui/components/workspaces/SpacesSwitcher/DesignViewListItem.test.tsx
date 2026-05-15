import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const getUserProfileSummaryByIdQuery = vi.fn();
const can = vi.fn(() => true);

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => can(...args),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetUserProfileSummaryByIdQuery: (...args: unknown[]) =>
    getUserProfileSummaryByIdQuery(...args),
}));

vi.mock('@/utils/Enum', () => ({
  RESOURCE_TYPE: { DESIGN: 'design', VIEW: 'view' },
}));

vi.mock('@/utils/context/WorkspaceModalContextProvider', () => ({
  WorkspaceModalContext: React.createContext({
    multiSelectedContent: [],
    setMultiSelectedContent: vi.fn(),
  }),
}));

vi.mock('../../shared/Modal/Information/InfoModal', () => ({
  VIEW_VISIBILITY: { PUBLIC: 'public', PRIVATE: 'private' },
}));

vi.mock('@sistent/sistent', () => {
  return {
    AvatarGroup: ({ children }: any) => <div data-testid="avatar-group">{children}</div>,
    Checkbox: ({ checked, onChange }: any) => (
      <input type="checkbox" checked={!!checked} onChange={onChange} />
    ),
    CustomTooltip: ({ children, title }: any) => (
      <div data-testid={`tooltip-${title}`}>{children}</div>
    ),
    Divider: () => <hr />,
    FormControlLabel: ({ control }: any) => control,
    FormGroup: ({ children }: any) => <div>{children}</div>,
    getFullFormattedTime: (v: any) => `full(${v})`,
    getRelativeTime: (v: any) => `rel(${v})`,
    LockIcon: () => <svg data-testid="lock" />,
    PublicIcon: () => <svg data-testid="public" />,
    Skeleton: ({ animation }: any) => <div data-testid={`skel-${animation || 'pulse'}`} />,
    Typography: ({ children }: any) => <span>{children}</span>,
    VisibilityChipMenu: ({ value, enabled, onChange }: any) => (
      <button
        data-testid="visibility-chip"
        data-enabled={String(enabled)}
        onClick={() => onChange('public')}
      >
        {value}
      </button>
    ),
    Grid2: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  };
});

vi.mock('./styles', () => ({
  StyledAvatarContainer: ({ children }: any) => <div>{children}</div>,
  StyledListIcon: ({ children }: any) => <span>{children}</span>,
  StyledListItem: ({ children, onClick, onMouseEnter, onMouseLeave, 'data-testid': tid }: any) => (
    <div
      data-testid={tid}
      onClick={onClick}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      {children}
    </div>
  ),
  StyledListItemText: ({ primary, secondary }: any) => (
    <div>
      <span data-testid="primary">{primary}</span>
      <span data-testid="secondary">{secondary}</span>
    </div>
  ),
  StyledSmallAvatar: ({ alt }: any) => <span data-testid="small-avatar">{alt}</span>,
  StyledSmallAvatarContainer: ({ children, id }: any) => (
    <div id={id} data-testid="small-avatar-container">
      {children}
    </div>
  ),
  StyledTypography: ({ children }: any) => <span>{children}</span>,
  StyledUpdatedText: ({ children }: any) => <span>{children}</span>,
  StyledUserDetailsContainer: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('css/icons.styles', () => ({ iconMedium: {} }));

vi.mock('./UserAvatarComponent', () => ({
  default: ({ userData }: any) => (
    <div data-testid="user-avatar">{userData?.firstName ?? 'unknown'}</div>
  ),
}));

vi.mock('./hooks', () => ({
  useGetIconBasedOnMode: () => <svg data-testid="resource-icon" />,
}));

import DesignViewListItem, { DesignViewListItemSkeleton } from './DesignViewListItem';

describe('DesignViewListItem', () => {
  beforeEach(() => {
    getUserProfileSummaryByIdQuery.mockReset();
    getUserProfileSummaryByIdQuery.mockReturnValue({ data: undefined, isLoading: false });
    can.mockReset();
    can.mockReturnValue(true);
  });

  const sampleItem = {
    id: 'd-1',
    name: 'My Design',
    updatedAt: '2026-01-01',
    visibility: 'public',
    userId: 'u-1',
    organizationName: 'Acme',
    workspaceName: 'ws',
  };

  it('renders the item name, organization, workspace and relative time', () => {
    render(
      <DesignViewListItem
        selectedItem={sampleItem}
        handleItemClick={() => {}}
        MenuComponent={<span>menu</span>}
        onVisibilityChange={() => {}}
        canChangeVisibility
        type="design"
      />,
    );
    expect(screen.getByTestId('primary')).toHaveTextContent('My Design');
    expect(screen.getByTestId('secondary')).toHaveTextContent('rel(2026-01-01)');
    expect(screen.getByText('Acme')).toBeInTheDocument();
    expect(screen.getByText('ws')).toBeInTheDocument();
  });

  it('uses embedded user data when available and skips the lookup query', () => {
    render(
      <DesignViewListItem
        selectedItem={{
          ...sampleItem,
          firstName: 'Alice',
          email: 'a@e.com',
        }}
        handleItemClick={() => {}}
        MenuComponent={<span>menu</span>}
        onVisibilityChange={() => {}}
        canChangeVisibility
        type="design"
      />,
    );
    expect(screen.getByTestId('user-avatar')).toHaveTextContent('Alice');
    expect(getUserProfileSummaryByIdQuery).toHaveBeenCalledWith({ id: 'u-1' }, { skip: true });
  });

  it('uses the fetched user when no embedded user data is provided', () => {
    getUserProfileSummaryByIdQuery.mockReturnValue({
      data: { firstName: 'Bob', id: 'u-1' },
      isLoading: false,
    });

    render(
      <DesignViewListItem
        selectedItem={sampleItem}
        handleItemClick={() => {}}
        MenuComponent={<span>menu</span>}
        onVisibilityChange={() => {}}
        canChangeVisibility
        type="design"
      />,
    );
    expect(screen.getByTestId('user-avatar')).toHaveTextContent('Bob');
    expect(getUserProfileSummaryByIdQuery).toHaveBeenCalledWith({ id: 'u-1' }, { skip: false });
  });

  it('triggers handleItemClick when the row is clicked', async () => {
    const user = userEvent.setup();
    const handleItemClick = vi.fn();

    render(
      <DesignViewListItem
        selectedItem={sampleItem}
        handleItemClick={handleItemClick}
        MenuComponent={<span>menu</span>}
        onVisibilityChange={() => {}}
        canChangeVisibility
        type="design"
      />,
    );

    await user.click(screen.getByTestId('design-tr-d-1'));
    expect(handleItemClick).toHaveBeenCalled();
  });

  it('calls onVisibilityChange with selectedItem when VisibilityChipMenu changes', async () => {
    const user = userEvent.setup();
    const onVisibilityChange = vi.fn();

    render(
      <DesignViewListItem
        selectedItem={sampleItem}
        handleItemClick={() => {}}
        MenuComponent={<span>menu</span>}
        onVisibilityChange={onVisibilityChange}
        canChangeVisibility
        type="design"
      />,
    );

    await user.click(screen.getByTestId('visibility-chip'));
    expect(onVisibilityChange).toHaveBeenCalledWith('public', sampleItem);
  });

  it('renders activeUsers when provided', () => {
    render(
      <DesignViewListItem
        selectedItem={sampleItem}
        handleItemClick={() => {}}
        MenuComponent={<span>menu</span>}
        onVisibilityChange={() => {}}
        canChangeVisibility
        type="design"
        activeUsers={[
          { client_id: 'c1', name: 'one', color: 'red' },
          { client_id: 'c2', name: 'two', color: 'blue' },
        ]}
      />,
    );

    expect(screen.getByTestId('avatar-group')).toBeInTheDocument();
    expect(screen.getAllByTestId('small-avatar')).toHaveLength(2);
  });

  it('renders a checkbox when isMultiSelectMode is true', () => {
    render(
      <DesignViewListItem
        selectedItem={sampleItem}
        handleItemClick={() => {}}
        MenuComponent={<span>menu</span>}
        onVisibilityChange={() => {}}
        canChangeVisibility
        type="design"
        isMultiSelectMode={true}
      />,
    );
    expect(screen.getByRole('checkbox')).toBeInTheDocument();
  });

  it('does not render workspace name when showWorkspaceName=false', () => {
    render(
      <DesignViewListItem
        selectedItem={{ ...sampleItem, workspaceName: 'hidden-ws' }}
        handleItemClick={() => {}}
        MenuComponent={<span>menu</span>}
        onVisibilityChange={() => {}}
        canChangeVisibility
        type="design"
        showWorkspaceName={false}
      />,
    );
    expect(screen.queryByText('hidden-ws')).not.toBeInTheDocument();
  });

  it('does not render organization name when showOrganizationName=false', () => {
    render(
      <DesignViewListItem
        selectedItem={{ ...sampleItem, organizationName: 'hidden-org' }}
        handleItemClick={() => {}}
        MenuComponent={<span>menu</span>}
        onVisibilityChange={() => {}}
        canChangeVisibility
        type="design"
        showOrganizationName={false}
      />,
    );
    expect(screen.queryByText('hidden-org')).not.toBeInTheDocument();
  });
});

describe('DesignViewListItemSkeleton', () => {
  it('renders multiple skeleton blocks', () => {
    const { container } = render(<DesignViewListItemSkeleton />);
    expect(container.querySelectorAll('[data-testid^="skel-"]').length).toBeGreaterThan(0);
  });

  it('shows the checkbox skeleton column when isMultiSelectMode is true', () => {
    const { container } = render(<DesignViewListItemSkeleton isMultiSelectMode={true} />);
    expect(container.querySelectorAll('[data-testid^="skel-"]').length).toBeGreaterThan(0);
  });
});
