import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';

const plainCardSpy = vi.fn();

vi.mock('@sistent/sistent', () => ({
  useTheme: () => ({
    palette: {
      icon: { default: '#000' },
      background: { brand: { default: '#abc' } },
    },
  }),
  PlainCard: ({
    resources,
    icon,
    title,
  }: {
    resources: Array<Record<string, unknown>>;
    icon: React.ReactNode;
    title: string;
  }) => {
    plainCardSpy({ resources, title });
    return (
      <div data-testid="plain-card" data-title={title}>
        {icon}
        <ul>
          {resources.map((r, i) => (
            <li key={i} data-testid="resource">
              {String((r as { name: string }).name)}
            </li>
          ))}
        </ul>
      </div>
    );
  },
  BellIcon: () => <svg data-testid="bell-icon" />,
  DesignIcon: () => <svg data-testid="design-icon" />,
}));

const widgetErrorFallbackSpy = vi.fn();

vi.mock('./WidgetErrorFallback', () => ({
  default: (props: { widgetTitle: string; message?: string }) => {
    widgetErrorFallbackSpy(props);
    return (
      <div data-testid="widget-error-fallback" data-title={props.widgetTitle}>
        {props.message}
      </div>
    );
  },
}));

import LatestBlogs from './LatestBlogs';

describe('LatestBlogs', () => {
  beforeEach(() => {
    plainCardSpy.mockReset();
    widgetErrorFallbackSpy.mockReset();
    vi.restoreAllMocks();
  });

  it('renders the Loading placeholder initially', () => {
    // fetch never resolves -> stays in loading state
    global.fetch = vi.fn().mockImplementation(() => new Promise(() => {})) as never;

    render(<LatestBlogs />);
    expect(screen.getByText('Loading...')).toBeInTheDocument();
    expect(screen.getByTestId('plain-card')).toHaveAttribute('data-title', 'LATEST BLOGS');
  });

  it('renders blog titles when fetch resolves with feed XML', async () => {
    const feedXml = `<?xml version="1.0"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <entry>
    <title>Post One</title>
    <content xml:base="https://meshery.io/posts/one"/>
  </entry>
  <entry>
    <title>Post Two</title>
    <content xml:base="https://meshery.io/posts/two"/>
  </entry>
</feed>`;

    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      text: () => Promise.resolve(feedXml),
    }) as never;

    render(<LatestBlogs />);
    await waitFor(() => {
      expect(screen.getByText('Post One')).toBeInTheDocument();
      expect(screen.getByText('Post Two')).toBeInTheDocument();
    });

    const finalCall = plainCardSpy.mock.calls.at(-1)?.[0] as {
      resources: Array<{ name: string }>;
    };
    expect(finalCall.resources.map((r) => r.name)).toEqual(['Post One', 'Post Two']);
  });

  it('falls back to "No Title" and # link when entry is malformed', async () => {
    const feedXml = `<?xml version="1.0"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <entry></entry>
</feed>`;

    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      text: () => Promise.resolve(feedXml),
    }) as never;

    render(<LatestBlogs />);
    await waitFor(() => {
      expect(screen.getByText('No Title')).toBeInTheDocument();
    });
  });

  it('caps blog list at the first five entries', async () => {
    const entries = Array.from(
      { length: 8 },
      (_, i) =>
        `<entry><title>Post ${i}</title><content xml:base="https://meshery.io/posts/${i}"/></entry>`,
    ).join('');
    const feedXml = `<?xml version="1.0"?><feed xmlns="http://www.w3.org/2005/Atom">${entries}</feed>`;

    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      text: () => Promise.resolve(feedXml),
    }) as never;

    render(<LatestBlogs />);
    await waitFor(() => {
      const finalCall = plainCardSpy.mock.calls.at(-1)?.[0] as {
        resources: Array<{ name: string }>;
      };
      expect(finalCall.resources).toHaveLength(5);
    });
  });

  it('logs the error and shows the error fallback on fetch failure', async () => {
    const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    global.fetch = vi.fn().mockRejectedValue(new Error('network down')) as never;

    render(<LatestBlogs />);
    await waitFor(() => {
      expect(errorSpy).toHaveBeenCalled();
      expect(screen.getByTestId('widget-error-fallback')).toHaveAttribute(
        'data-title',
        'Latest Blogs',
      );
    });
    expect(screen.queryByTestId('plain-card')).not.toBeInTheDocument();
  });

  it('shows the error fallback when the feed responds with a non-ok status', async () => {
    const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 500,
      text: () => Promise.resolve(''),
    }) as never;

    render(<LatestBlogs />);
    await waitFor(() => {
      expect(errorSpy).toHaveBeenCalled();
      expect(screen.getByTestId('widget-error-fallback')).toBeInTheDocument();
    });
  });

  it('shows an empty-state message when the feed returns zero entries', async () => {
    const feedXml = `<?xml version="1.0"?>
<feed xmlns="http://www.w3.org/2005/Atom"></feed>`;

    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      text: () => Promise.resolve(feedXml),
    }) as never;

    render(<LatestBlogs />);
    await waitFor(() => {
      expect(screen.getByText('No blog posts found.')).toBeInTheDocument();
    });
    expect(screen.queryByTestId('widget-error-fallback')).not.toBeInTheDocument();
  });
});
