import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { afterAll, afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

// Pins the request the component details panel ACTUALLY sends, through the
// real registry hooks and a real store. The generated getRegistryComponents
// operation has no `id` param, so an `id` passed to it is silently dropped and
// the panel downloads every untrimmed component of that apiVersion.

// The schemas client reads its base URL at module load.
const { previousEndpointPrefix } = vi.hoisted(() => {
  const previous = process.env.RTK_MESHERY_ENDPOINT_PREFIX;
  process.env.RTK_MESHERY_ENDPOINT_PREFIX = 'http://localhost';
  return { previousEndpointPrefix: previous };
});

afterAll(() => {
  if (previousEndpointPrefix === undefined) {
    delete process.env.RTK_MESHERY_ENDPOINT_PREFIX;
  } else {
    process.env.RTK_MESHERY_ENDPOINT_PREFIX = previousEndpointPrefix;
  }
});

// rtk-query/utils imports the app store at module scope.
vi.mock('../../store', () => ({ store: { dispatch: vi.fn() } }));

vi.mock('@/theme', () => ({
  useTheme: () => ({
    palette: { text: { default: 'black', secondary: 'gray' } },
    typography: { fontFamily: 'sans-serif', caption: { fontSize: '0.75rem' } },
  }),
}));

vi.mock('@sistent/sistent', () => ({
  FormControl: ({ children }: any) => <div>{children}</div>,
  Select: ({ children }: any) => <select>{children}</select>,
  MenuItem: ({ children }: any) => <option>{children}</option>,
  CircularProgress: () => <div data-testid="loading" />,
  Button: ({ children }: any) => <button>{children}</button>,
  DownloadIcon: () => <svg />,
  ExpandMoreIcon: () => <svg />,
  Accordion: ({ children }: any) => <div>{children}</div>,
  AccordionSummary: ({ children }: any) => <div>{children}</div>,
  AccordionDetails: ({ children }: any) => <div>{children}</div>,
  styled: (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) =>
      React.createElement(Component, props, children);
    return StyledComponent;
  },
}));

vi.mock('next/dynamic', () => ({
  default: () => () => <div data-testid="react-json" />,
}));

vi.mock('@/assets/styles/general/tool.styles', () => ({
  DetailsContainer: ({ children }: any) => <div>{children}</div>,
  Segment: ({ children }: any) => <div>{children}</div>,
  FullWidth: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('@/components/data-formatter', () => ({
  FormatStructuredData: () => <div />,
  reorderObjectProperties: (o: any) => o,
}));

vi.mock('./MeshModel.style', () => ({
  JustifyAndAlignCenter: ({ children }: any) => <div>{children}</div>,
  StyledKeyValueFormattedValue: ({ children }: any) => <div>{children}</div>,
  StyledKeyValuePropertyDiv: ({ children }: any) => <div>{children}</div>,
  StyledKeyValueProperty: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('./helper', () => ({ reactJsonTheme: () => ({}) }));
vi.mock('css/icons.styles', () => ({ iconSmall: {} }));

import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import { api } from '@/rtk-query/index';
import MeshModelDetails from './MeshModelDetails';

const okJsonResponse = (body: unknown) => ({
  ok: true,
  status: 200,
  redirected: false,
  headers: new Headers({ 'content-type': 'application/json' }),
  url: '',
  text: () => Promise.resolve(JSON.stringify(body)),
  json: () => Promise.resolve(body),
  clone() {
    return this;
  },
});

const requestUrl = (input: RequestInfo | URL) =>
  typeof input === 'string' ? input : input instanceof URL ? input.href : input.url;

const componentDef = {
  id: '6f1d1a8e-6b0e-4c1e-9b7a-2f4a0c9d1e01',
  displayName: 'Deployment',
  component: { kind: 'Deployment', version: 'apps/v1' },
  model: { name: 'kubernetes', displayName: 'Kubernetes' },
};

describe('MeshModelDetails component request', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn(async () =>
      okJsonResponse({
        page: 0,
        pageSize: 1,
        totalCount: 1,
        components: [
          {
            ...componentDef,
            component: {
              ...componentDef.component,
              schema: JSON.stringify({ description: 'Declarative updates for Pods' }),
            },
            registrant: { hostname: 'artifacthub' },
            duplicates: 0,
          },
        ],
      }),
    );
    vi.spyOn(globalThis, 'fetch').mockImplementation(fetchMock as unknown as typeof fetch);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('fetches only the selected component, scoped by model, kind and apiVersion', async () => {
    const store = configureStore({
      reducer: { [api.reducerPath]: api.reducer },
      middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware),
    });

    render(
      <Provider store={store}>
        <MeshModelDetails
          view="Components"
          showDetailsData={{ type: 'Components', data: componentDef }}
        />
      </Provider>,
    );

    await waitFor(() => expect(fetchMock).toHaveBeenCalled());
    const urls = fetchMock.mock.calls.map(([input]) => new URL(requestUrl(input)));

    // Never the registry-wide listing, which cannot be narrowed to one component.
    expect(urls.map((url) => url.pathname)).not.toContain('/api/registry/components');
    expect(urls).toHaveLength(1);
    expect(urls[0].pathname).toBe('/api/registry/models/kubernetes/components/Deployment');
    expect(urls[0].searchParams.get('apiVersion')).toBe('apps/v1');

    expect(await screen.findByText('Declarative updates for Pods')).toBeInTheDocument();
  });
});
