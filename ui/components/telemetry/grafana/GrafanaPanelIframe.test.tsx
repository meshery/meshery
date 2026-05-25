import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import GrafanaPanelIframe from './GrafanaPanelIframe';

describe('GrafanaPanelIframe', () => {
  it('renders an iframe with the supplied src and title', () => {
    const { container } = render(
      <GrafanaPanelIframe src="https://example.com/panel" title="Demo Panel" />,
    );
    const iframe = container.querySelector('iframe');
    expect(iframe).not.toBeNull();
    expect(iframe).toHaveAttribute('src', 'https://example.com/panel');
    expect(iframe).toHaveAttribute('title', 'Demo Panel');
    expect(iframe).toHaveAttribute('loading', 'lazy');
  });
});
