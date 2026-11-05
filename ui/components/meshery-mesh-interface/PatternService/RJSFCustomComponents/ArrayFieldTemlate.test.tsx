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

vi.mock('../../../../css/icons.styles', () => ({ iconSmall: {} }));

vi.mock('pluralize', () => ({
  default: { singular: (s: string) => (s.endsWith('s') ? s.slice(0, -1) : s) },
}));

vi.mock('../helper', () => ({
  safeDisplayValue: (v: any) => (v == null ? '' : String(v)),
  safeStringTitle: (v: any) => (v == null ? '' : String(v)),
}));

import ArrayFieldTemplate from './ArrayFieldTemlate';

// A stand-in for the element RJSF builds per array item: a keyed element carrying the item's
// props and rendering its own content, exactly like `ArrayFieldItem` does.
const ArrayItem = (props: any) => <div data-testid={props.itemKey}>{props.children}</div>;

const makeArrayItem = (key: string, index: number) => (
  <ArrayItem key={key} itemKey={key} index={index} disabled={false} readonly={false}>
    <span>{`item ${key}`}</span>
  </ArrayItem>
);

describe('ArrayFieldTemplate', () => {
  const baseProps = {
    schema: { items: { type: 'string' } },
    uiSchema: {},
    fieldPathId: { $id: 'array1', path: ['array1'] },
    items: [makeArrayItem('a', 0), makeArrayItem('b', 1)],
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

  it('renders the content of each item element', () => {
    render(<ArrayFieldTemplate {...(baseProps as any)} />);

    // Each entry of `items` is rendered as an element, so its own content shows up.
    expect(screen.getByText('item a')).toBeInTheDocument();
    expect(screen.getByText('item b')).toBeInTheDocument();
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

  it('renders the fixed template variant for a tuple-items (fixed) array', () => {
    // RJSF routes a multi-select array to its own widget before the template
    // runs, so only a tuple `items` reaches the template as a fixed array.
    const { container } = render(
      <ArrayFieldTemplate
        {...(baseProps as any)}
        schema={{
          title: 'Items',
          type: 'array',
          description: 'pick some',
          items: [{ type: 'string' }, { type: 'string' }],
        }}
      />,
    );
    // The fixed variant inlines the description and wraps items in a fieldset;
    // the normal variant shows it behind a help tooltip instead.
    expect(container.querySelector('fieldset')).toBeInTheDocument();
    expect(container.querySelector('.field-description')).toHaveTextContent('pick some');
    expect(screen.queryByTestId('help-icon')).not.toBeInTheDocument();
  });

  it('renders the normal template variant for a single-item schema', () => {
    const { container } = render(
      <ArrayFieldTemplate
        {...(baseProps as any)}
        schema={{
          title: 'Items',
          type: 'array',
          description: 'pick some',
          items: { type: 'string' },
        }}
      />,
    );
    expect(screen.getByTestId('help-icon')).toBeInTheDocument();
    expect(container.querySelector('.field-description')).not.toBeInTheDocument();
    expect(container.querySelector('fieldset')).not.toBeInTheDocument();
  });

  it('shows the help tooltip when description is provided via uiSchema', () => {
    render(
      <ArrayFieldTemplate {...(baseProps as any)} uiSchema={{ 'ui:description': 'array desc' }} />,
    );
    expect(screen.getByTestId('help-icon')).toBeInTheDocument();
  });
});
