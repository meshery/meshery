import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { ResponsiveImage } from './ResponsiveImage';

describe('ResponsiveImage', () => {
  it('renders an img with the provided src and forwards the test id', () => {
    render(
      <ResponsiveImage src="/extension-icon.png" alt="Extension icon" testId="extension-icon" />,
    );

    const img = screen.getByTestId('extension-icon');
    expect(img.tagName).toBe('IMG');
    expect(img).toHaveAttribute('src', '/extension-icon.png');
    expect(img).toHaveAttribute('alt', 'Extension icon');
  });

  it('defaults alt to an empty string when none is provided', () => {
    render(<ResponsiveImage src="/extension-icon.png" testId="extension-icon" />);

    expect(screen.getByTestId('extension-icon')).toHaveAttribute('alt', '');
  });
});
