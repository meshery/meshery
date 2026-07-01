import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi, beforeEach } from 'vitest';

let loggedInReturn: { data?: { id?: string } } = { data: { id: 'user-1' } };
let patternsReturn: {
  data?: { patterns?: Array<{ name: string; id: string; updatedAt?: string }> };
  isFetching?: boolean;
} = { data: { patterns: [] }, isFetching: false };

const designCardSpy = vi.fn();
const useGetUserDesignsQuerySpy = vi.fn();

vi.mock('@/rtk-query/user', () => ({
  useGetLoggedInUserQuery: () => loggedInReturn,
}));

vi.mock('@/rtk-query/design', () => ({
  useGetUserDesignsQuery: (...args: unknown[]) => {
    useGetUserDesignsQuerySpy(...args);
    return patternsReturn;
  },
}));

vi.mock('@/constants/endpoints', () => ({
  MESHERY_CLOUD_PROD: 'https://cloud.meshery.io',
}));

vi.mock('@sistent/sistent', () => ({
  Box: ({
    children,
    sx: _sx,
    ...rest
  }: { children: React.ReactNode; sx?: unknown } & React.HTMLAttributes<HTMLDivElement>) => (
    <div {...rest}>{children}</div>
  ),
  CatalogIcon: () => <svg data-testid="catalog-icon" />,
  DesignIcon: () => <svg data-testid="design-icon" />,
  EditIcon: () => <svg data-testid="edit-icon" />,
  GithubIcon: () => <svg data-testid="github-icon" />,
  useTheme: () => ({
    palette: { icon: { default: '#000' } },
  }),
  DesignCard: (props: {
    resources: Array<{ name: string; link: string }>;
    isPatternsFetching: boolean;
    title: string;
    btnTitle: string;
    href: string;
    cardData: Array<{ title: string }>;
    sortOrder: string;
    setSortOrder: (order: string) => void;
    actionButton: boolean;
  }) => {
    designCardSpy(props);
    return (
      <div
        data-testid="design-card"
        data-title={props.title}
        data-btn={props.btnTitle}
        data-href={props.href}
        data-loading={String(props.isPatternsFetching)}
        data-sort={props.sortOrder}
      >
        <button type="button" onClick={() => props.setSortOrder('created_at desc')}>
          changeOrder
        </button>
        <ul>
          {props.resources.map((r, i) => (
            <li key={i} data-testid="resource-link">
              <a href={r.link}>{r.name}</a>
            </li>
          ))}
        </ul>
        <ul>
          {props.cardData.map((c) => (
            <li key={c.title} data-testid="card-title">
              {c.title}
            </li>
          ))}
        </ul>
      </div>
    );
  },
}));

import MyDesignsWidget from './MyDesignsWidget';

describe('MyDesignsWidget', () => {
  beforeEach(() => {
    designCardSpy.mockReset();
    useGetUserDesignsQuerySpy.mockReset();
    loggedInReturn = { data: { id: 'user-1' } };
    patternsReturn = { data: { patterns: [] }, isFetching: false };
  });

  it('renders the DesignCard with default props', () => {
    render(<MyDesignsWidget />);
    const card = screen.getByTestId('design-card');
    expect(card).toHaveAttribute('data-title', 'MY RECENT DESIGNS');
    expect(card).toHaveAttribute('data-btn', 'See All Designs');
    expect(card).toHaveAttribute(
      'data-href',
      'https://cloud.meshery.io/catalog/content/my-designs',
    );
  });

  it('renders three pre-defined card data items (scratch / template / github)', () => {
    render(<MyDesignsWidget />);
    const titles = screen.getAllByTestId('card-title').map((el) => el.textContent);
    expect(titles).toContain('Create a design from scratch');
    expect(titles).toContain('Choose a template to start with');
    expect(titles).toContain('Import design from github');
  });

  it('maps patterns into resource entries with the expected link format', () => {
    patternsReturn = {
      data: {
        patterns: [
          { id: 'd1', name: 'Design One', updatedAt: '2024-01-01' },
          { id: 'd2', name: 'Design Two', updatedAt: '2024-01-02' },
        ],
      },
      isFetching: false,
    };
    render(<MyDesignsWidget />);
    const resources = screen.getAllByTestId('resource-link');
    expect(resources).toHaveLength(2);
    expect(resources[0]).toHaveTextContent('Design One');
    expect(resources[0].querySelector('a')).toHaveAttribute(
      'href',
      '/extension/meshmap?mode=design&design=d1',
    );
  });

  it('reflects the loading state via DesignCard prop', () => {
    patternsReturn = { isFetching: true };
    render(<MyDesignsWidget />);
    expect(screen.getByTestId('design-card')).toHaveAttribute('data-loading', 'true');
  });

  it('lets users change the sort order via the setSortOrder callback', async () => {
    render(<MyDesignsWidget />);
    expect(screen.getByTestId('design-card')).toHaveAttribute('data-sort', 'updated_at desc');
    await userEvent.click(screen.getByRole('button', { name: 'changeOrder' }));
    expect(screen.getByTestId('design-card')).toHaveAttribute('data-sort', 'created_at desc');
  });

  it('skips fetching designs until the logged-in user id is available', () => {
    loggedInReturn = {};
    render(<MyDesignsWidget />);
    const [, options] = useGetUserDesignsQuerySpy.mock.calls[0];
    expect(options).toEqual({ skip: true });
  });

  it('does not skip fetching designs once the logged-in user id is available', () => {
    render(<MyDesignsWidget />);
    const [, options] = useGetUserDesignsQuerySpy.mock.calls[0];
    expect(options).toEqual({ skip: false });
  });
});
