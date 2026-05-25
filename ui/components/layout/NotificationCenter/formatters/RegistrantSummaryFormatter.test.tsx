import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
  LaunchIcon: () => <svg data-testid="launch-icon" />,
}));

vi.mock('./common', () => ({
  TitleLink: ({ href, children, target, rel }: any) => (
    <a data-testid="title-link" href={href} target={target} rel={rel}>
      {children}
    </a>
  ),
}));

import { RegistrantSummaryFormatter } from './RegistrantSummaryFormatter';

describe('RegistrantSummaryFormatter', () => {
  it('renders the event description and a documentation link', () => {
    render(<RegistrantSummaryFormatter event={{ description: 'New registrant processed' }} />);

    expect(screen.getByText('New registrant processed')).toBeInTheDocument();
    const link = screen.getByTestId('title-link');
    expect(link).toHaveAttribute(
      'href',
      'https://docs.meshery.io/concepts/logical#logical-concepts',
    );
    expect(link).toHaveTextContent(/understanding models/i);
  });

  it('renders the link even when description is missing', () => {
    render(<RegistrantSummaryFormatter event={{}} />);
    expect(screen.getByTestId('title-link')).toBeInTheDocument();
  });
});
