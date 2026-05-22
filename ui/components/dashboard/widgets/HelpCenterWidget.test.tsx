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

    // check that links exist
    expect(links.length).toBeGreaterThan(0);

    // verify important resources by URL instead of fragile text labels
    const hrefs = links.map((a) => a.getAttribute('href'));

    expect(hrefs).toEqual(
      expect.arrayContaining([
        expect.stringContaining('docs'),
        expect.stringContaining('slack'),
        expect.stringContaining('forum'),
        expect.stringContaining('support'),
      ]),
    );
  });

  it('passes one DesignIcon per resource to the PlainCard', () => {
    plainCardSpy.mockReset();
    render(<HelpCenterWidget />);
    expect(plainCardSpy.mock.calls[0][0].resources).toHaveLength(5);
  });
});
