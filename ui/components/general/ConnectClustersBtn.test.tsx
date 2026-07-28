import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';

const h = vi.hoisted(() => ({
  openCreateConnectionMock: vi.fn(),
}));

vi.mock('@/utils/context/ConnectionWizardContextProvider', () => ({
  useConnectionWizardModal: () => ({
    openCreateConnection: h.openCreateConnectionMock,
    closeCreateConnection: vi.fn(),
    open: false,
    presetKind: null,
    skipKindSelection: false,
  }),
}));

vi.mock('@sistent/sistent', () => ({
  AddCircleIcon: (props: any) => <svg data-testid="add-icon" {...props} />,
  Button: ({ children, onClick, ...props }: any) => (
    <button type="button" onClick={onClick} {...props}>
      {children}
    </button>
  ),
  useTheme: () => ({ spacing: (n: number) => `${n * 8}px` }),
}));

vi.mock('../../css/icons.styles', () => ({
  iconMedium: {},
}));

import ConnectClustersBtn from './ConnectClustersBtn';

describe('ConnectClustersBtn', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('opens the Create Connection wizard in place with Kubernetes pre-selected', () => {
    render(<ConnectClustersBtn />);

    fireEvent.click(screen.getByRole('button', { name: /Connect Clusters/i }));

    expect(h.openCreateConnectionMock).toHaveBeenCalledWith({
      kind: 'kubernetes',
      skipKindSelection: true,
    });
  });

  it('renders the Connect Clusters label and add icon', () => {
    render(<ConnectClustersBtn />);
    expect(screen.getByText(/Connect Clusters/i)).toBeInTheDocument();
    expect(screen.getByTestId('add-icon')).toBeInTheDocument();
  });
});
