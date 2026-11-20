import React from 'react';
import { render, screen, cleanup } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  LaunchIcon: (props: any) => <svg data-testid="launch-icon" {...props} />,
  Typography: ({ children, variant, ...props }: any) => (
    <span data-testid="typography" data-variant={variant} {...props}>
      {children}
    </span>
  ),
}));

vi.mock('../../../data-formatter', () => ({
  TextWithLinks: ({ text }: { text: string }) => <span data-testid="text-with-links">{text}</span>,
}));

import { TitleLink, EmptyState, DataToFileLink } from './common';

describe('TitleLink', () => {
  it('renders an anchor with href and target=_blank by default', () => {
    render(<TitleLink href="https://example.com">Click me</TitleLink>);
    const link = screen.getByRole('link', { name: /click me/i });
    expect(link).toHaveAttribute('href', 'https://example.com');
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveAttribute('rel', 'noopener noreferrer');
    expect(screen.getByTestId('launch-icon')).toBeInTheDocument();
  });

  it('allows overriding target via props', () => {
    render(
      <TitleLink href="/internal" target="_self">
        Internal
      </TitleLink>,
    );
    const link = screen.getByRole('link', { name: /internal/i });
    expect(link).toHaveAttribute('target', '_self');
  });

  it('renders the provided children inside the typography', () => {
    render(<TitleLink href="/x">Hello world</TitleLink>);
    expect(screen.getByTestId('typography')).toHaveTextContent('Hello world');
  });
});

describe('EmptyState', () => {
  it('renders the event description via TextWithLinks', () => {
    render(<EmptyState event={{ description: 'No description' }} />);
    expect(screen.getByTestId('text-with-links')).toHaveTextContent('No description');
  });

  it('falls back to empty text when description is missing', () => {
    render(<EmptyState event={{}} />);
    expect(screen.getByTestId('text-with-links')).toBeInTheDocument();
    expect(screen.getByTestId('text-with-links')).toHaveTextContent('');
  });
});

describe('DataToFileLink', () => {
  // jsdom does not implement URL.createObjectURL or URL.revokeObjectURL by default; stub them for these tests.
  const originalCreateObjectURL = (global.URL as any).createObjectURL;
  const originalRevokeObjectURL = (global.URL as any).revokeObjectURL;

  beforeEach(() => {
    (global.URL as any).createObjectURL = vi.fn((file: File) => `blob:fake-url-${file.name}`);
    (global.URL as any).revokeObjectURL = vi.fn();
  });

  afterEach(() => {
    cleanup();
    (global.URL as any).createObjectURL = originalCreateObjectURL;
    (global.URL as any).revokeObjectURL = originalRevokeObjectURL;
  });

  it('creates a download link pointing to a generated blob URL', () => {
    render(<DataToFileLink data="hello world" />);
    const link = screen.getByRole('link', { name: /download trace/i });
    expect(link).toHaveAttribute('href', 'blob:fake-url-trace.txt');
    expect(link).toHaveAttribute('download', 'trace.txt');
  });

  it('stringifies object data when rendering a download link', () => {
    render(<DataToFileLink data={{ foo: 'bar' }} />);
    const createSpy = (global.URL as any).createObjectURL as ReturnType<typeof vi.fn>;
    expect(createSpy).toHaveBeenCalled();
    const blobArg = createSpy.mock.calls[0][0];
    expect(blobArg).toBeInstanceOf(File);
    expect(blobArg.name).toBe('trace.txt');
  });

  it('revokes the blob URL on component unmount to prevent memory leaks', () => {
    const { unmount } = render(<DataToFileLink data="trace info" />);
    const revokeSpy = (global.URL as any).revokeObjectURL as ReturnType<typeof vi.fn>;

    expect(revokeSpy).not.toHaveBeenCalled();

    unmount();

    expect(revokeSpy).toHaveBeenCalledTimes(1);
    expect(revokeSpy).toHaveBeenCalledWith('blob:fake-url-trace.txt');
  });

  it('revokes previous blob URL and creates a new one when data prop updates', () => {
    const { rerender } = render(<DataToFileLink data="first-trace" />);
    const createSpy = (global.URL as any).createObjectURL as ReturnType<typeof vi.fn>;
    const revokeSpy = (global.URL as any).revokeObjectURL as ReturnType<typeof vi.fn>;

    expect(createSpy).toHaveBeenCalledTimes(1);
    expect(revokeSpy).not.toHaveBeenCalled();

    rerender(<DataToFileLink data="second-trace" />);

    expect(revokeSpy).toHaveBeenCalledTimes(1);
    expect(revokeSpy).toHaveBeenCalledWith('blob:fake-url-trace.txt');
    expect(createSpy).toHaveBeenCalledTimes(2);
  });
});
