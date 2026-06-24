import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: any) => <div>{children}</div>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  Paper: ({ children }: any) => <div>{children}</div>,
  Button: ({ children, onClick, disabled, className }: any) => (
    <button data-testid={className || 'button'} onClick={onClick} disabled={!!disabled}>
      {children}
    </button>
  ),
  IconButton: ({ children, onClick, disabled, className }: any) => (
    <button data-testid={className || 'icon-button'} onClick={onClick} disabled={!!disabled}>
      {children}
    </button>
  ),
  Typography: ({ children }: any) => <span>{children}</span>,
  useTheme: () => ({ palette: { mode: 'light', error: { main: '#f00' } } }),
}));

vi.mock('../../../../assets/icons/AddIcon', () => ({
  default: () => <svg data-testid="add-icon" />,
}));

vi.mock('./Accordion', () => ({
  default: ({ children, heading, childProps }: any) => (
    <div
      data-testid="accordion"
      data-heading={heading}
      data-has-remove={String(!!childProps?.hasRemove)}
    >
      {children}
    </div>
  ),
}));

vi.mock('../CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
}));

vi.mock('../../../../assets/icons/HelpOutlineIcon', () => ({
  default: () => <svg data-testid="help-icon" />,
}));

vi.mock('../../../../assets/icons/ErrorOutlineIcon', () => ({
  default: () => <svg data-testid="error-icon" />,
}));

vi.mock('@rjsf/utils', () => ({
  isMultiSelect: (schema: any) => !!schema?.uniqueItems && Array.isArray(schema.items?.enum),
  getDefaultFormState: () => ({ rootSchema: {} }),
}));

vi.mock('../../../../css/icons.styles', () => ({ iconSmall: {} }));

vi.mock('pluralize', () => ({
  default: { singular: (s: string) => (s.endsWith('s') ? s.slice(0, -1) : s) },
}));

vi.mock('../helper', () => ({
  safeDisplayValue: (v: any) => (v == null ? '' : String(v)),
  safeStringTitle: (v: any) => (v == null ? '' : String(v)),
}));

import ArrayFieldTemplate from './ArrayFieldTemlate';

describe('ArrayFieldTemplate', () => {
  const baseProps = {
    schema: { items: { type: 'string' } },
    uiSchema: {},
    idSchema: { $id: 'array1' },
    items: [
      {
        key: 'a',
        children: <div>item-a</div>,
        index: 0,
        hasMoveUp: false,
        hasMoveDown: true,
        hasToolbar: false,
        onReorderClick: () => vi.fn(),
      },
      {
        key: 'b',
        children: <div>item-b</div>,
        index: 1,
        hasMoveUp: true,
        hasMoveDown: false,
        hasToolbar: false,
        onReorderClick: () => vi.fn(),
      },
    ],
    canAdd: true,
    onAddClick: vi.fn(),
    title: 'Items',
    rawErrors: [],
  };

  it('renders a normal array template with the title and item list', () => {
    render(<ArrayFieldTemplate {...(baseProps as any)} />);
    expect(screen.getByText('Items')).toBeInTheDocument();
    const accordions = screen.getAllByTestId('accordion');
    expect(accordions).toHaveLength(2);
  });

  it('renders an Add icon button when canAdd is true', () => {
    render(<ArrayFieldTemplate {...(baseProps as any)} />);
    expect(screen.getByTestId('array-item-add')).toBeInTheDocument();
  });

  it('invokes onAddClick when the add button is clicked', () => {
    const onAddClick = vi.fn();
    render(<ArrayFieldTemplate {...(baseProps as any)} onAddClick={onAddClick} />);
    fireEvent.click(screen.getByTestId('array-item-add'));
    expect(onAddClick).toHaveBeenCalled();
  });

  it('shows an error tooltip when rawErrors are present', () => {
    render(<ArrayFieldTemplate {...(baseProps as any)} rawErrors={['oops']} />);
    expect(screen.getByTestId('error-icon')).toBeInTheDocument();
  });

  it('renders the fixed template variant when isMultiSelect returns true', () => {
    render(
      <ArrayFieldTemplate
        {...(baseProps as any)}
        schema={{ uniqueItems: true, items: { enum: ['a', 'b'] } }}
      />,
    );
    expect(screen.getByTestId('array-item-add')).toBeInTheDocument();
  });

  it('shows the help tooltip when description is provided via uiSchema', () => {
    render(
      <ArrayFieldTemplate {...(baseProps as any)} uiSchema={{ 'ui:description': 'array desc' }} />,
    );
    expect(screen.getByTestId('help-icon')).toBeInTheDocument();
  });
});
