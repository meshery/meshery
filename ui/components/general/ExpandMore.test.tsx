import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

// style.tsx pulls a number of primitives from the design system at module load. Stub them
// (as the other component tests do) so we can exercise the real ExpandMore ARIA logic
// without a theme provider. Only IconButton and CaretDownIcon are actually rendered here;
// everything reached through styled() is collapsed to its children.
vi.mock('@sistent/sistent', () => {
  const styled =
    () =>
    () =>
    ({ children }: any) => <>{children}</>;
  return {
    styled,
    IconButton: ({ children, ...props }: any) => <button {...props}>{children}</button>,
    CaretDownIcon: (props: any) => <svg data-testid="caret" {...props} />,
    alpha: () => 'rgba(0,0,0,0.1)',
    charcoal: {},
    Box: 'div',
    Button: 'button',
    ButtonGroup: 'div',
    Divider: 'hr',
    Drawer: 'div',
    List: 'ul',
    ListItem: 'li',
    ListItemButton: 'button',
    ListItemIcon: 'div',
    ListItemText: 'span',
  };
});

vi.mock('../../css/disableComponent.styles', () => ({
  disabledStyleWithOutOpacity: {},
}));

import { ExpandMore } from './style';

const theme = { palette: { icon: { brand: '#000' } } } as any;

describe('ExpandMore', () => {
  it('announces a collapsed caret as expandable', () => {
    render(<ExpandMore isExpanded={false} hasChildren theme={theme} />);
    const button = screen.getByRole('button');
    expect(button).toHaveAttribute('aria-expanded', 'false');
    expect(button).toHaveAttribute('aria-label', 'Expand');
  });

  it('announces an expanded caret as collapsible', () => {
    render(<ExpandMore isExpanded hasChildren theme={theme} />);
    const button = screen.getByRole('button');
    expect(button).toHaveAttribute('aria-expanded', 'true');
    expect(button).toHaveAttribute('aria-label', 'Collapse');
  });
});
