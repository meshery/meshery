import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Hoisted query data we can swap between tests.
const queryState: { data: any } = { data: undefined };

vi.mock('../../../rtk-query/notificationCenter', () => ({
  useGetEventFiltersQuery: () => queryState,
}));

vi.mock('@/components/shared/FormFields/typing-filter', () => ({
  default: ({ filterSchema, defaultFilters, placeholder, handleFilter }: any) => (
    <div data-testid="typing-filter" data-placeholder={placeholder}>
      <button type="button" onClick={() => handleFilter([{ type: 'STATUS', value: 'read' }])}>
        invoke-filter
      </button>
      <span data-testid="default-filters">{JSON.stringify(defaultFilters)}</span>
      <span data-testid="schema-json">{JSON.stringify(filterSchema)}</span>
    </div>
  ),
}));

import Filter from './filter';

describe('NotificationCenter Filter', () => {
  beforeEach(() => {
    queryState.data = undefined;
  });

  it('renders a TypingFilter with the right placeholder and default filters', () => {
    const handleFilter = vi.fn();
    render(<Filter handleFilter={handleFilter} />);

    expect(screen.getByTestId('typing-filter')).toHaveAttribute(
      'data-placeholder',
      'Filter Notifications',
    );
    expect(screen.getByTestId('default-filters')).toHaveTextContent('"status: unread"');
  });

  it('forwards handleFilter to TypingFilter', () => {
    const handleFilter = vi.fn();
    render(<Filter handleFilter={handleFilter} />);

    screen.getByText('invoke-filter').click();
    expect(handleFilter).toHaveBeenCalledTimes(1);
    expect(handleFilter).toHaveBeenCalledWith([{ type: 'STATUS', value: 'read' }]);
  });

  it('includes severity, status, action, author, and category filter definitions', () => {
    render(<Filter handleFilter={vi.fn()} />);
    const schema = JSON.parse(screen.getByTestId('schema-json').textContent || '{}');

    expect(schema.SEVERITY.value).toBe('severity');
    expect(Array.isArray(schema.SEVERITY.values)).toBe(true);
    expect(schema.SEVERITY.values).toEqual(
      expect.arrayContaining(['informational', 'error', 'warning', 'success']),
    );
    expect(schema.STATUS.value).toBe('status');
    expect(schema.STATUS.multiple).toBe(false);
    expect(schema.STATUS.values).toEqual(expect.arrayContaining(['read', 'unread']));
    expect(schema.ACTION.value).toBe('action');
    expect(schema.ACTION.values).toEqual([]);
    expect(schema.AUTHOR.value).toBe('author');
    expect(schema.CATEGORY.value).toBe('category');
    expect(schema.CATEGORY.values).toEqual([]);
  });

  it('uses action/category values from the RTK query response when available', () => {
    queryState.data = { action: ['deploy', 'undeploy'], category: ['pattern'] };
    render(<Filter handleFilter={vi.fn()} />);
    const schema = JSON.parse(screen.getByTestId('schema-json').textContent || '{}');

    expect(schema.ACTION.values).toEqual(['deploy', 'undeploy']);
    expect(schema.CATEGORY.values).toEqual(['pattern']);
  });
});
