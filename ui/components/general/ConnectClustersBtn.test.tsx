import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import ConnectClustersBtn from './ConnectClustersBtn';

vi.mock('@sistent/sistent', () => ({
  AddCircleIcon: (props: any) => <svg data-testid="add-icon" {...props} />,
  Button: ({ children, ...props }: any) => (
    <button type="button" {...props}>
      {children}
    </button>
  ),
  useTheme: () => ({ spacing: (n: number) => `${n * 8}px` }),
}));

vi.mock('next/link', () => ({
  default: ({ children, href }: any) => <a href={href}>{children}</a>,
}));

vi.mock('../../css/icons.styles', () => ({
  iconMedium: {},
}));

describe('ConnectClustersBtn', () => {
  it('renders a link to the connections page', () => {
    render(<ConnectClustersBtn />);
    expect(screen.getByRole('link')).toHaveAttribute('href', '/management/connections');
  });

  it('renders the Connect Clusters label and add icon', () => {
    render(<ConnectClustersBtn />);
    expect(screen.getByText(/Connect Clusters/i)).toBeInTheDocument();
    expect(screen.getByTestId('add-icon')).toBeInTheDocument();
  });
});
