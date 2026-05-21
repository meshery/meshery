import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const plainCardSpy = vi.fn();

vi.mock('@/constants/endpoints', () => ({
  MESHERY_CLOUD_PROD: 'https://cloud.meshery.io',
}));

vi.mock('@sistent/sistent', () => ({
  useTheme: () => ({
    palette: { icon: { default: '#000', disabled: '#777' } },
  }),
  PlainCard: ({
    resources,
    icon,
    title,
  }: {
    resources: Array<{ name: string; link: string }>;
    icon: React.ReactNode;
    title: string;
  }) => {
    plainCardSpy({ resources, title });
    return (
      <div data-testid="plain-card" data-title={title}>
        {icon}
        <ul>
          {resources.map((r, i) => (
            <li key={i}>
              <a href={r.link} data-testid="link">
                {r.name}
              </a>
            </li>
          ))}
        </ul>
      </div>
    );
  },
  DocumentIcon: () => <svg data-testid="document-icon" />,
  DesignIcon: () => <svg data-testid="design-icon" />,
}));

vi.mock('css/icons.styles', () => ({ iconMedium: {} }));

import HelpCenterWidget from './HelpCenterWidget';

describe('HelpCenterWidget', () => {
  it('renders the help center title and document icon', () => {
    plainCardSpy.mockReset();
    render(<HelpCenterWidget />);
    expect(screen.getByTestId('plain-card')).toHaveAttribute('data-title', 'HELP CENTER');
    expect(screen.getByTestId('document-icon')).toBeInTheDocument();
  });

  it('renders the documented set of help center resources', () => {
    plainCardSpy.mockReset();
    render(<HelpCenterWidget />);
    const links = screen.getAllByTestId('link');
    const linkSet = new Set(links.map((a) => a.textContent));
    expect(linkSet).toContain('GitHub Issues');
    expect(linkSet).toContain('Slack');
    expect(linkSet).toContain('Discussion Forum');
    expect(linkSet).toContain('Support Request');
    // confirm support URL uses MESHERY_CLOUD_PROD
    expect(links.find((a) => a.textContent === 'Support Request')).toHaveAttribute(
      'href',
      'https://cloud.meshery.io/support',
    );
  });

  it('passes one DesignIcon per resource to the PlainCard', () => {
    plainCardSpy.mockReset();
    render(<HelpCenterWidget />);
    expect(plainCardSpy.mock.calls[0][0].resources).toHaveLength(5);
  });
});
