import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/components/general/ConnectClustersBtn', () => ({
  default: () => <button type="button">ConnectClustersBtn</button>,
}));

vi.mock('../style', () => ({
  ConnectClusterWrapper: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="connect-cluster-wrapper">{children}</div>
  ),
  ConnectClusterText: ({
    children,
    variant,
    align,
  }: {
    children?: React.ReactNode;
    variant?: string;
    align?: string;
  }) => (
    <p data-testid="connect-cluster-text" data-variant={variant} data-align={align}>
      {children}
    </p>
  ),
}));

import ConnectCluster from './ConnectCluster';

describe('ConnectCluster', () => {
  it('renders the message inside the wrapper and the connect button', () => {
    render(<ConnectCluster message="No connections found" />);
    expect(screen.getByTestId('connect-cluster-wrapper')).toBeInTheDocument();
    const text = screen.getByTestId('connect-cluster-text');
    expect(text).toHaveTextContent('No connections found');
    expect(text).toHaveAttribute('data-variant', 'h5');
    expect(text).toHaveAttribute('data-align', 'center');
    expect(screen.getByRole('button', { name: 'ConnectClustersBtn' })).toBeInTheDocument();
  });

  it('renders ReactNode messages', () => {
    render(<ConnectCluster message={<span data-testid="rich">Try again</span>} />);
    expect(screen.getByTestId('rich')).toHaveTextContent('Try again');
  });
});
