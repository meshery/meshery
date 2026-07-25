import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';
import { K8sEmptyState } from './K8sContextEmptyState';

const h = vi.hoisted(() => ({
  openCreateConnectionMock: vi.fn(),
}));

let mockMode: 'light' | 'dark' = 'light';

vi.mock('@/utils/context/ConnectionWizardContextProvider', () => ({
  useConnectionWizardModal: () => ({
    openCreateConnection: h.openCreateConnectionMock,
    closeCreateConnection: vi.fn(),
    open: false,
    presetKind: null,
    skipKindSelection: false,
  }),
}));

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = (props: any) => <Component {...props}>{props.children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };

  return {
    AddIcon: (props: any) => <svg data-testid="add-icon" {...props} />,
    Button: ({ children, type, variant, color, sx, onClick, ...rest }: any) => (
      <button
        type={type}
        data-variant={variant}
        data-color={color}
        data-sx={JSON.stringify(sx || {})}
        onClick={onClick}
        {...rest}
      >
        {children}
      </button>
    ),
    Typography: ({ children, variant }: any) => (
      <p data-variant={variant} data-testid="typography">
        {children}
      </p>
    ),
    styled,
    useTheme: () => ({ palette: { mode: mockMode } }),
  };
});

vi.mock('../../../assets/img/OperatorLight', () => ({
  default: () => <div data-testid="operator-light" />,
}));

vi.mock('../../../assets/img/Operator', () => ({
  default: () => <div data-testid="operator" />,
}));

describe('K8sEmptyState', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders the default empty state in light mode', () => {
    mockMode = 'light';

    render(<K8sEmptyState message={undefined} />);

    expect(screen.getByTestId('operator')).toBeInTheDocument();
    expect(screen.queryByTestId('operator-light')).not.toBeInTheDocument();
    expect(screen.getByTestId('typography')).toHaveTextContent('No cluster connected yet');
    expect(screen.getByRole('button')).toHaveTextContent('Connect Clusters');
  });

  it('opens the Create Connection wizard in place with Kubernetes pre-selected', () => {
    mockMode = 'light';
    render(<K8sEmptyState message={undefined} />);

    fireEvent.click(screen.getByRole('button', { name: /Connect Clusters/i }));

    expect(h.openCreateConnectionMock).toHaveBeenCalledWith({
      kind: 'kubernetes',
      skipKindSelection: true,
    });
  });

  it('renders the Operator dark variant in dark mode', () => {
    mockMode = 'dark';

    render(<K8sEmptyState message="No active cluster found" />);

    expect(screen.getByTestId('operator-light')).toBeInTheDocument();
    expect(screen.queryByTestId('operator')).not.toBeInTheDocument();
  });

  it('uses the provided message when present', () => {
    mockMode = 'light';
    render(<K8sEmptyState message="No active cluster found" />);

    expect(screen.getByTestId('typography')).toHaveTextContent('No active cluster found');
  });
});
