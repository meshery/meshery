import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import CreateDesignBtn from './CreateDesignBtn';

vi.mock('@sistent/sistent', () => ({
  AddCircleIcon: (props: any) => <svg data-testid="add-icon" {...props} />,
  Button: ({ children, ...props }: any) => (
    <button type="button" {...props}>
      {children}
    </button>
  ),
}));

vi.mock('next/link', () => ({
  default: ({ children, href }: any) => <a href={href}>{children}</a>,
}));

vi.mock('css/icons.styles', () => ({
  iconMedium: {},
}));

describe('CreateDesignBtn', () => {
  it('renders a link to the design configuration page', () => {
    render(<CreateDesignBtn />);
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/configuration/design');
  });

  it('renders the Create Design label and an add icon', () => {
    render(<CreateDesignBtn />);
    expect(screen.getByText(/Create Design/i)).toBeInTheDocument();
    expect(screen.getByTestId('add-icon')).toBeInTheDocument();
  });
});
