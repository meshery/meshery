import React from 'react';
import { act, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const queryState: { data: any } = { data: undefined };
const sliceState = {
  // Match events slice: current_view.filters starts as { initial: true } until loadEvents resolves.
  filters: { initial: true } as Record<string, unknown>,
};
let typingFilterMountCount = 0;

vi.mock('react-redux', () => ({
  useSelector: (sel: any) =>
    sel({
      events: {
        current_view: { filters: sliceState.filters },
      },
    }),
}));

vi.mock('../../../rtk-query/notificationCenter', () => ({
  useGetEventFiltersQuery: () => queryState,
}));

vi.mock('@/components/shared/FormFields/typing-filter', () => ({
  default: ({ filterSchema, defaultFilters, placeholder, handleFilter }: any) => {
    React.useEffect(() => {
      typingFilterMountCount += 1;
    }, []);
    return (
      <div data-testid="typing-filter" data-placeholder={placeholder}>
        <button type="button" onClick={() => handleFilter({})}>
          clear-filters
        </button>
        <button type="button" onClick={() => handleFilter({ status: 'read', severity: ['error'] })}>
          invoke-filter
        </button>
        <button type="button" onClick={() => handleFilter({ status: ['unread', 'read'] })}>
          replace-status
        </button>
        <button type="button" onClick={() => handleFilter({ status: [] })}>
          empty-status
        </button>
        <span data-testid="default-filters">{JSON.stringify(defaultFilters)}</span>
        <span data-testid="schema-json">{JSON.stringify(filterSchema)}</span>
      </div>
    );
  },
}));

import Filter, { filtersToChips, normalizeFilterPayload } from './filter';

describe('NotificationCenter Filter', () => {
  beforeEach(() => {
    queryState.data = undefined;
    sliceState.filters = { initial: true };
    typingFilterMountCount = 0;
  });

  it('shows unread chips on startup when current_view.filters is still { initial: true }', () => {
    render(<Filter handleFilter={vi.fn()} />);

    expect(screen.getByTestId('typing-filter')).toHaveAttribute(
      'data-placeholder',
      'Filter Notifications',
    );
    const chips = JSON.parse(screen.getByTestId('default-filters').textContent || '[]');
    expect(chips).toEqual([{ type: 'STATUS', value: 'unread', label: 'status: unread' }]);
  });

  it('maps post-fetch Redux filters into TypingFilter chips', () => {
    sliceState.filters = {
      status: 'read',
      severity: ['error', 'warning'],
    };
    render(<Filter handleFilter={vi.fn()} />);

    const chips = JSON.parse(screen.getByTestId('default-filters').textContent || '[]');
    expect(chips).toEqual(
      expect.arrayContaining([
        { type: 'STATUS', value: 'read', label: 'status: read' },
        { type: 'SEVERITY', value: 'error', label: 'severity: error' },
        { type: 'SEVERITY', value: 'warning', label: 'severity: warning' },
      ]),
    );
    expect(chips).toHaveLength(3);
  });

  it('allows clearing all filters including unread', () => {
    sliceState.filters = { status: 'unread' };
    const handleFilter = vi.fn();
    render(<Filter handleFilter={handleFilter} />);

    act(() => {
      screen.getByText('clear-filters').click();
    });
    expect(handleFilter).toHaveBeenCalledWith({});
  });

  it('shows no chips after user clears filters in Redux', () => {
    sliceState.filters = {};
    render(<Filter handleFilter={vi.fn()} />);

    const chips = JSON.parse(screen.getByTestId('default-filters').textContent || '[]');
    expect(chips).toEqual([]);
  });

  it('remounts TypingFilter when clearing so local chip state resyncs', () => {
    sliceState.filters = { status: 'unread' };
    render(<Filter handleFilter={vi.fn()} />);
    expect(typingFilterMountCount).toBe(1);

    act(() => {
      screen.getByText('clear-filters').click();
    });
    expect(typingFilterMountCount).toBe(2);
  });

  it('forwards non-empty filter changes', () => {
    const handleFilter = vi.fn();
    render(<Filter handleFilter={handleFilter} />);

    screen.getByText('invoke-filter').click();
    expect(handleFilter).toHaveBeenCalledWith({ status: 'read', severity: ['error'] });
  });

  it('normalizes array status from TypingFilter to scalar read when replacing unread', () => {
    const handleFilter = vi.fn();
    render(<Filter handleFilter={handleFilter} />);

    screen.getByText('replace-status').click();
    expect(handleFilter).toHaveBeenCalledWith({ status: 'read' });
  });

  it('forwards empty payload when normalized status array is empty', () => {
    const handleFilter = vi.fn();
    render(<Filter handleFilter={handleFilter} />);

    act(() => {
      screen.getByText('empty-status').click();
    });
    expect(handleFilter).toHaveBeenCalledWith({});
  });

  it('includes severity, status, action, author, and category filter definitions', () => {
    render(<Filter handleFilter={vi.fn()} />);
    const schema = JSON.parse(screen.getByTestId('schema-json').textContent || '{}');

    expect(schema.SEVERITY.values).toEqual(
      expect.arrayContaining(['informational', 'error', 'warning', 'success']),
    );
    expect(schema.STATUS.multiple).toBe(false);
    expect(schema.STATUS.values).toEqual(expect.arrayContaining(['read', 'unread']));
    expect(schema.ACTION.value).toBe('action');
    expect(schema.AUTHOR.value).toBe('author');
    expect(schema.CATEGORY.value).toBe('category');
  });

  it('uses action/category values from the RTK query response when available', () => {
    queryState.data = { action: ['deploy', 'undeploy'], category: ['pattern'] };
    render(<Filter handleFilter={vi.fn()} />);
    const schema = JSON.parse(screen.getByTestId('schema-json').textContent || '{}');

    expect(schema.ACTION.values).toEqual(['deploy', 'undeploy']);
    expect(schema.CATEGORY.values).toEqual(['pattern']);
  });
});

describe('filtersToChips', () => {
  const schema = {
    STATUS: { value: 'status', multiple: false },
    SEVERITY: { value: 'severity' },
  };

  it('shows unread chip only for the initial pre-fetch sentinel', () => {
    expect(filtersToChips({ initial: true }, schema)).toEqual([
      expect.objectContaining({ type: 'STATUS', value: 'unread' }),
    ]);
  });

  it('returns no chips for an empty cleared filter object', () => {
    expect(filtersToChips({}, schema)).toEqual([]);
  });

  it('maps status and severity filters to chips', () => {
    expect(filtersToChips({ status: 'read', severity: ['warning'] }, schema)).toEqual([
      { type: 'STATUS', value: 'read', label: 'status: read' },
      { type: 'SEVERITY', value: 'warning', label: 'severity: warning' },
    ]);
  });
});

describe('normalizeFilterPayload', () => {
  const schema = {
    STATUS: { value: 'status', multiple: false },
    SEVERITY: { value: 'severity' },
  };

  it('coerces array status to the last selected value', () => {
    expect(normalizeFilterPayload({ status: ['unread', 'read'] }, schema)).toEqual({
      status: 'read',
    });
  });

  it('leaves scalar status and array severity unchanged', () => {
    expect(
      normalizeFilterPayload({ status: 'unread', severity: ['error', 'warning'] }, schema),
    ).toEqual({ status: 'unread', severity: ['error', 'warning'] });
  });

  it('removes empty array status so clear yields an empty payload', () => {
    expect(normalizeFilterPayload({ status: [] }, schema)).toEqual({});
  });
});
