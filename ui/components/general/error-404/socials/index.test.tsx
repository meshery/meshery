import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import Socials from './index';

vi.mock('@sistent/sistent', () => ({
  Tooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
  Fade: ({ children }: any) => <>{children}</>,
}));

vi.mock('./styles', () => ({
  SocialMain: ({ children }: any) => <div data-testid="social-main">{children}</div>,
  SocialContainer: ({ children }: any) => <div data-testid="social-container">{children}</div>,
  MessageIcon: (props: any) => <svg data-testid="message-icon" {...props} />,
  SlackIcon: (props: any) => <svg data-testid="slack-icon" {...props} />,
  TwitterIcon: (props: any) => <svg data-testid="twitter-icon" {...props} />,
  GithubIcon: (props: any) => <svg data-testid="github-icon" {...props} />,
  YoutubeIcon: (props: any) => <svg data-testid="youtube-icon" {...props} />,
  DockerIcon: (props: any) => <svg data-testid="docker-icon" {...props} />,
}));

describe('Socials', () => {
  it('renders the social main container and inner container', () => {
    render(<Socials />);
    expect(screen.getByTestId('social-main')).toBeInTheDocument();
    expect(screen.getByTestId('social-container')).toBeInTheDocument();
  });

  it('renders all six social icons with tooltip titles', () => {
    render(<Socials />);
    expect(screen.getByTestId('message-icon')).toBeInTheDocument();
    expect(screen.getByTestId('slack-icon')).toBeInTheDocument();
    expect(screen.getByTestId('twitter-icon')).toBeInTheDocument();
    expect(screen.getByTestId('github-icon')).toBeInTheDocument();
    expect(screen.getByTestId('youtube-icon')).toBeInTheDocument();
    expect(screen.getByTestId('docker-icon')).toBeInTheDocument();

    const tooltips = screen.getAllByTestId('tooltip');
    expect(tooltips).toHaveLength(6);

    const titles = tooltips.map((node) => node.getAttribute('data-title'));
    expect(titles).toEqual(
      expect.arrayContaining([
        expect.stringMatching(/Meshery community/i),
        expect.stringMatching(/Slack/i),
        expect.stringMatching(/Meshery on X/i),
        expect.stringMatching(/Contribute to Meshery/i),
        expect.stringMatching(/community meeting/i),
        expect.stringMatching(/Docker images/i),
      ]),
    );
  });

  it('renders external links to each social platform', () => {
    render(<Socials />);
    const links = screen.getAllByRole('link');
    const hrefs = links.map((link) => link.getAttribute('href'));
    expect(hrefs).toEqual(
      expect.arrayContaining([
        'mailto:maintainers@meshery.io',
        'https://slack.meshery.io',
        'https://x.com/mesheryio',
        'https://github.com/meshery',
        expect.stringContaining('youtube.com'),
        'https://hub.docker.com/u/meshery/',
      ]),
    );
  });
});
