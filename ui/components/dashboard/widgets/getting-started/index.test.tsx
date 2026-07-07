import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi, beforeEach } from 'vitest';

let loggedIn: { data?: { id?: string } } = { data: { id: 'user-1' } };
let profile: {
  data?: { preferences?: { remoteProviderPreferences?: { getstarted?: unknown[] } } };
} = {
  data: { preferences: { remoteProviderPreferences: { getstarted: ['step-1'] } } },
};
let currentOrg: { id?: string } | null = { id: 'org-1' };

const actionCardSpy = vi.fn();
const modalSpy = vi.fn();

vi.mock('@/rtk-query/user', () => ({
  useGetLoggedInUserQuery: () => loggedIn,
  useGetUserByIdQuery: () => profile,
  useHandleUserInviteMutation: () => [vi.fn()],
  useLazyGetTeamsQuery: () => [vi.fn()],
  useUpdateUserPrefMutation: () => [vi.fn()],
}));

vi.mock('@/rtk-query/orgRoles', () => ({
  useGetUserOrgRolesQuery: () => ({ data: [] }),
}));

vi.mock('@/rtk-query/organization', () => ({
  useGetOrgsQuery: () => ({ data: [] }),
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotificationHandlers: () => ({}),
}));

vi.mock('@/utils/can', () => ({ default: () => true }));

vi.mock('@/utils/permission_constants', () => ({
  keys: { ASSIGN_USER_ROLES: { action: 'assign', subject: 'roles' } },
}));

vi.mock('react-redux', () => ({
  useSelector: (selector: (state: unknown) => unknown) =>
    selector({ ui: { organization: currentOrg } }),
}));

vi.mock('css/icons.styles', () => ({ iconMedium: {} }));

vi.mock('@sistent/sistent', () => ({
  Link: ({ children, href }: { children?: React.ReactNode; href?: string }) => (
    <a href={href}>{children}</a>
  ),
  useTheme: () => ({ palette: { icon: { default: '#000' } } }),
  GetStartedIcon: () => <svg data-testid="get-started-icon" />,
  ActionButtonCard: (props: {
    title: string;
    description: string;
    btnTitle: string;
    onClick?: () => void;
    showProgress?: boolean;
    completedSteps?: unknown[];
    totalSteps?: number;
  }) => {
    actionCardSpy(props);
    return (
      <div data-testid="action-card" data-title={props.title}>
        <button onClick={props.onClick} type="button">
          {props.btnTitle}
        </button>
      </div>
    );
  },
  GetStartedModal: (props: {
    open: boolean;
    handleClose?: () => void;
    isFromMeshery?: boolean;
  }) => {
    modalSpy(props);
    return (
      <div
        data-testid="modal"
        data-open={String(props.open)}
        data-meshery={String(props.isFromMeshery)}
      >
        <button type="button" onClick={props.handleClose}>
          close
        </button>
      </div>
    );
  },
}));

import GetStarted from './index';

describe('GetStarted', () => {
  beforeEach(() => {
    actionCardSpy.mockReset();
    modalSpy.mockReset();
    loggedIn = { data: { id: 'user-1' } };
    profile = {
      data: { preferences: { remoteProviderPreferences: { getstarted: ['step-1'] } } },
    };
    currentOrg = { id: 'org-1' };
  });

  it('renders the action card with the GETTING STARTED title and step progress', () => {
    render(<GetStarted />);
    const card = screen.getByTestId('action-card');
    expect(card).toHaveAttribute('data-title', 'GETTING STARTED');
    const props = actionCardSpy.mock.calls[0][0];
    expect(props.showProgress).toBe(true);
    expect(props.completedSteps).toEqual(['step-1']);
    expect(props.totalSteps).toBeGreaterThanOrEqual(1);
  });

  it('opens the modal when the Start button is clicked', async () => {
    render(<GetStarted />);
    expect(screen.getByTestId('modal')).toHaveAttribute('data-open', 'false');
    await userEvent.click(screen.getByRole('button', { name: 'Start' }));
    expect(screen.getByTestId('modal')).toHaveAttribute('data-open', 'true');
  });

  it('closes the modal when the close button is clicked', async () => {
    render(<GetStarted />);
    await userEvent.click(screen.getByRole('button', { name: 'Start' }));
    expect(screen.getByTestId('modal')).toHaveAttribute('data-open', 'true');
    await userEvent.click(screen.getByRole('button', { name: 'close' }));
    expect(screen.getByTestId('modal')).toHaveAttribute('data-open', 'false');
  });

  it('forwards isFromMeshery=true to the GetStartedModal', () => {
    render(<GetStarted />);
    expect(screen.getByTestId('modal')).toHaveAttribute('data-meshery', 'true');
  });

  it('defaults completedSteps to [] when remoteProviderPreferences.getstarted is missing', () => {
    profile = { data: { preferences: {} } };
    render(<GetStarted />);
    const props = actionCardSpy.mock.calls[0][0];
    expect(props.completedSteps).toEqual([]);
  });

  it('does not crash and passes an undefined currentOrgId when the organization is null', () => {
    currentOrg = null;
    render(<GetStarted />);
    expect(screen.getByTestId('action-card')).toBeInTheDocument();
    const modalProps = modalSpy.mock.calls[0][0];
    expect(modalProps.currentOrgId).toBeUndefined();
  });
});
