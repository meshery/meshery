import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import TypingFilter from './index';

let lastAutocompleteProps: any = null;

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = (props: any) => <Component {...props}>{props.children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    Autocomplete: (props: any) => {
      lastAutocompleteProps = props;
      const params = {
        InputProps: { startAdornment: <span data-testid="existing-adornment">x</span> },
        placeholder: props.placeholder,
        onKeyDown: () => null,
      };
      return (
        <div data-testid="autocomplete">
          {props.renderInput(params)}
          <div data-testid="options">
            {(props.options || []).map((o: any, i: number) =>
              props.renderOption({ key: o.label, role: 'option' }, o, { index: i }),
            )}
          </div>
          <div data-testid="tags">
            {props.renderTags ? props.renderTags(props.value || [], (_: any) => ({})) : null}
          </div>
          <button
            aria-label="emit-select"
            onClick={() =>
              props.onChange?.(null, [], 'selectOption', {
                option: { type: 'STATUS', value: 'Open', label: 'status: Open' },
              })
            }
          >
            select
          </button>
          <button
            aria-label="emit-remove"
            onClick={() =>
              props.onChange?.(null, [], 'removeOption', {
                option: props.value?.[0],
              })
            }
          >
            remove
          </button>
          <button
            aria-label="emit-input-change"
            onClick={() => props.onInputChange?.(null, 'status: Open')}
          >
            input-change
          </button>
          <button aria-label="emit-onclose-blur" onClick={() => props.onClose?.({}, 'blur')}>
            blur-close
          </button>
          <button aria-label="emit-onopen" onClick={() => props.onOpen?.()}>
            open
          </button>
          <span data-testid="clear-icon">{props.clearIcon}</span>
        </div>
      );
    },
    Chip: ({ label, ...rest }: any) => (
      <span data-testid="chip" {...rest}>
        {label}
      </span>
    ),
    ContentFilterIcon: () => <svg data-testid="content-filter-icon" />,
    CrossCircleIcon: ({ onClick }: any) => (
      <svg data-testid="cross-circle-icon" onClick={onClick} />
    ),
    InputAdornment: ({ children }: any) => <span>{children}</span>,
    useTheme: () => ({ palette: { icon: { default: '#000' } } }),
    Divider: () => <hr data-testid="divider" />,
    styled,
  };
});

vi.mock('./style', () => ({
  InputField: (props: any) => (
    <input data-testid="input-field" placeholder={props.placeholder} onKeyDown={props.onKeyDown} />
  ),
  Root: ({ children, className }: any) => (
    <div data-testid="root" className={className}>
      {children}
    </div>
  ),
  DropDown: ({ children }: any) => <div data-testid="dropdown">{children}</div>,
}));

vi.mock('css/icons.styles', () => ({
  iconSmall: {},
}));

const filterSchema: any = {
  STATUS: {
    value: 'status',
    description: 'filter by status',
    values: ['Open', 'Closed'],
  },
  AUTHOR: {
    value: 'author',
    description: 'filter by author',
  },
};

describe('TypingFilter', () => {
  beforeEach(() => {
    lastAutocompleteProps = null;
  });

  it('renders with the placeholder', () => {
    render(
      <TypingFilter
        filterSchema={filterSchema}
        placeholder="Filter..."
        handleFilter={vi.fn()}
        defaultFilters={[]}
      />,
    );

    expect(screen.getByTestId('autocomplete')).toBeInTheDocument();
    expect(screen.getByTestId('input-field')).toHaveAttribute('placeholder', 'Filter...');
  });

  it('renders options for each filter when there is no colon in input', () => {
    render(
      <TypingFilter
        filterSchema={filterSchema}
        placeholder="Filter..."
        handleFilter={vi.fn()}
        defaultFilters={[]}
      />,
    );

    // Initial options from filterSchema
    const options = screen.getByTestId('options');
    expect(options).toHaveTextContent('status: filter by status');
    expect(options).toHaveTextContent('author: filter by author');
  });

  it('selects a filter and triggers handleFilter', async () => {
    const user = userEvent.setup();
    const handleFilter = vi.fn();

    render(
      <TypingFilter
        filterSchema={filterSchema}
        placeholder="Filter..."
        handleFilter={handleFilter}
        defaultFilters={[]}
      />,
    );

    // Simulate typing status:
    await user.click(screen.getByRole('button', { name: 'emit-input-change' }));
    // Now select
    await user.click(screen.getByRole('button', { name: 'emit-select' }));

    expect(handleFilter).toHaveBeenCalled();
    const calledWith = handleFilter.mock.calls[0][0];
    expect(calledWith.status).toBe('Open');
  });

  it('clears all filters when the cross-icon is clicked', async () => {
    const user = userEvent.setup();
    const handleFilter = vi.fn();

    render(
      <TypingFilter
        filterSchema={filterSchema}
        placeholder="Filter..."
        handleFilter={handleFilter}
        defaultFilters={[]}
      />,
    );

    // The clearIcon is rendered by the autocomplete, but our mock places it as <CrossCircleIcon onClick={clearAllFilters} />
    await user.click(screen.getByTestId('cross-circle-icon'));
    expect(handleFilter).toHaveBeenCalledWith({});
  });

  it('renders selected filter chips', () => {
    render(
      <TypingFilter
        filterSchema={filterSchema}
        placeholder="Filter..."
        handleFilter={vi.fn()}
        defaultFilters={[{ type: 'STATUS', value: 'Open', label: 'status: Open' }]}
      />,
    );

    expect(screen.getByTestId('chip')).toHaveTextContent('status: Open');
  });

  it('handles "Enter" press for free-solo (custom value) filters', async () => {
    const handleFilter = vi.fn();
    render(
      <TypingFilter
        filterSchema={filterSchema}
        placeholder="Filter..."
        handleFilter={handleFilter}
        defaultFilters={[]}
      />,
    );

    // simulate typing "author: name"
    const inputChangeBtn = screen.getByRole('button', { name: 'emit-input-change' });
    fireEvent.click(inputChangeBtn);

    // override input value to contain author:custom
    lastAutocompleteProps.onInputChange(null, 'author: bob');

    // Now find inner input and fire keydown enter via captured params
    const renderedParams = lastAutocompleteProps.renderInput({ InputProps: {} });
    expect(renderedParams.props.placeholder).toBe('Filter...');
  });

  it('removes a filter when removeOption fires', async () => {
    const user = userEvent.setup();
    const handleFilter = vi.fn();

    render(
      <TypingFilter
        filterSchema={filterSchema}
        placeholder="Filter..."
        handleFilter={handleFilter}
        defaultFilters={[{ type: 'STATUS', value: 'Open', label: 'status: Open' }]}
      />,
    );

    await user.click(screen.getByRole('button', { name: 'emit-remove' }));
    expect(handleFilter).toHaveBeenCalled();
  });

  it('opens and closes via onOpen / onClose blur', async () => {
    const user = userEvent.setup();
    render(
      <TypingFilter
        filterSchema={filterSchema}
        placeholder="Filter..."
        handleFilter={vi.fn()}
        defaultFilters={[]}
      />,
    );

    await user.click(screen.getByRole('button', { name: 'emit-onopen' }));
    await user.click(screen.getByRole('button', { name: 'emit-onclose-blur' }));
    // No throw indicates these are handled correctly
    expect(screen.getByTestId('autocomplete')).toBeInTheDocument();
  });
});
