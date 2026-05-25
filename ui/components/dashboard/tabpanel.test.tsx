import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { TabPanel } from './tabpanel';

vi.mock('@sistent/sistent', () => ({
  Typography: ({ children, ...props }: any) => (
    <div data-testid="typography" {...props}>
      {children}
    </div>
  ),
}));

describe('TabPanel', () => {
  it('renders children when value matches index', () => {
    render(
      <TabPanel value={1} index={1}>
        <span>active content</span>
      </TabPanel>,
    );
    expect(screen.getByText('active content')).toBeInTheDocument();
    expect(screen.getByRole('tabpanel')).not.toHaveAttribute('hidden');
  });

  it('hides and omits children when value does not match index', () => {
    render(
      <TabPanel value={1} index={2}>
        <span>inactive content</span>
      </TabPanel>,
    );
    expect(screen.queryByText('inactive content')).not.toBeInTheDocument();
    const panel = screen.getByRole('tabpanel', { hidden: true });
    expect(panel).toHaveAttribute('hidden');
  });

  it('sets accessible ids and aria-labelledby from index', () => {
    render(
      <TabPanel value="Overview" index="Overview">
        <span>panel</span>
      </TabPanel>,
    );
    const panel = screen.getByRole('tabpanel');
    expect(panel).toHaveAttribute('id', 'simple-tabpanel-Overview');
    expect(panel).toHaveAttribute('aria-labelledby', 'simple-tab-Overview');
  });

  it('forwards additional props to the root tabpanel div', () => {
    render(
      <TabPanel value={0} index={0} data-extra="present">
        <span>x</span>
      </TabPanel>,
    );
    expect(screen.getByRole('tabpanel')).toHaveAttribute('data-extra', 'present');
  });
});
