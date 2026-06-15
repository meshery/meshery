/**
 * Tests for the dynamic connection type registration flow in StepperContent.
 *
 * Covers:
 * - deriveConnectionTypes: pure unit tests (no rendering needed)
 * - SelectConnection component: dynamic rendering, loading/error/empty states,
 *   stable kind passed to registration, no hardcoded list
 */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi, beforeEach } from 'vitest';
import { deriveConnectionTypes } from './StepperContent';

// ---------------------------------------------------------------------------
// Unit tests for deriveConnectionTypes (pure function, no mocks needed)
// ---------------------------------------------------------------------------

describe('deriveConnectionTypes', () => {
  it('returns an empty array for empty input', () => {
    expect(deriveConnectionTypes([])).toEqual([]);
  });

  it('returns an empty array for null/undefined input', () => {
    expect(deriveConnectionTypes(null as any)).toEqual([]);
    expect(deriveConnectionTypes(undefined as any)).toEqual([]);
  });

  it('filters out components whose kind does not end with "Connection"', () => {
    const components = [
      { component: { kind: 'PrometheusCredential' } },
      { component: { kind: 'SomeRandomComponent' } },
      { component: { kind: 'GrafanaCredential' } },
    ];
    expect(deriveConnectionTypes(components)).toEqual([]);
  });

  it('extracts kind by removing the "Connection" suffix — not by slicing a label', () => {
    const components = [{ component: { kind: 'PrometheusConnection' } }];
    const result = deriveConnectionTypes(components);
    expect(result).toHaveLength(1);
    expect(result[0].kind).toBe('Prometheus');
    expect(result[0].label).toBe('Prometheus Connection');
  });

  it('handles multiple connection types correctly', () => {
    const components = [
      { component: { kind: 'PrometheusConnection' } },
      { component: { kind: 'GrafanaConnection' } },
      { component: { kind: 'PrometheusCredential' } }, // should be filtered out
    ];
    const result = deriveConnectionTypes(components);
    expect(result).toHaveLength(2);
    expect(result.map((r) => r.kind)).toEqual(['Prometheus', 'Grafana']);
    expect(result.map((r) => r.label)).toEqual(['Prometheus Connection', 'Grafana Connection']);
  });

  it('deduplicates components with the same kind', () => {
    const components = [
      { component: { kind: 'PrometheusConnection' } },
      { component: { kind: 'PrometheusConnection' } }, // duplicate
    ];
    const result = deriveConnectionTypes(components);
    expect(result).toHaveLength(1);
    expect(result[0].kind).toBe('Prometheus');
  });

  it('handles components with missing or null kind gracefully', () => {
    const components = [
      { component: { kind: null as any } },
      { component: {} },
      {},
      { component: { kind: 'GrafanaConnection' } },
    ];
    const result = deriveConnectionTypes(components);
    expect(result).toHaveLength(1);
    expect(result[0].kind).toBe('Grafana');
  });

  it('automatically includes newly registered {Kind}Connection schemas', () => {
    // Simulates a new "Jaeger" connection type being registered in the backend
    // without any frontend code change.
    const components = [
      { component: { kind: 'PrometheusConnection' } },
      { component: { kind: 'GrafanaConnection' } },
      { component: { kind: 'JaegerConnection' } }, // newly registered
    ];
    const result = deriveConnectionTypes(components);
    expect(result.map((r) => r.kind)).toContain('Jaeger');
    expect(result.map((r) => r.label)).toContain('Jaeger Connection');
  });

  it('does NOT derive kind by splitting on spaces in the label', () => {
    // Ensures the implementation uses suffix removal, not indexOf(' ')
    const components = [{ component: { kind: 'PrometheusConnection' } }];
    const result = deriveConnectionTypes(components);
    // If kind were derived by label.slice(0, label.indexOf(' ')), it would be
    // empty string because "PrometheusConnection" has no space. The correct
    // implementation returns "Prometheus".
    expect(result[0].kind).toBe('Prometheus');
    expect(result[0].kind).not.toBe('');
  });
});

// ---------------------------------------------------------------------------
// Component tests for SelectConnection
// ---------------------------------------------------------------------------

// Mock RTK query hooks
const mockUseGetComponentsQuery = vi.fn();
const mockUseVerifyAndRegisterConnectionMutation = vi.fn();

vi.mock('@/rtk-query/meshModel', () => ({
  useGetComponentsQuery: (...args) => mockUseGetComponentsQuery(...args),
}));

vi.mock('@/rtk-query/connection', () => ({
  useVerifyAndRegisterConnectionMutation: () => mockUseVerifyAndRegisterConnectionMutation(),
  useConnectToConnectionMutation: () => [vi.fn(), {}],
}));

vi.mock('@/rtk-query/credentials', () => ({
  useGetCredentialsQuery: () => ({ data: { credentials: [] } }),
}));

vi.mock('@/components/shared/FormFields/rjsf-utils/common', () => ({
  selectCompSchema: (enums: string[]) => ({
    properties: {
      selectedConnectionType: {
        enum: enums,
        type: 'string',
        title: 'test',
        uniqueItems: true,
        'x-rjsf-grid-area': 12,
        description: 'test',
      },
    },
    required: ['selectedConnectionType'],
    type: 'object',
  }),
}));

vi.mock('../../../meshery-mesh-interface/PatternService/RJSF_wrapper', () => ({
  default: ({ jsonSchema, onChange }) => (
    <div data-testid="rjsf-wrapper">
      {jsonSchema?.properties?.selectedConnectionType?.enum?.map((label: string) => (
        <button
          key={label}
          data-testid={`option-${label}`}
          onClick={() => onChange({ selectedConnectionType: label })}
        >
          {label}
        </button>
      ))}
    </div>
  ),
}));

vi.mock('./StepperContentWrapper', () => ({
  default: ({ children }) => <div data-testid="stepper-content">{children}</div>,
}));

vi.mock('../../../../utils/utils', () => ({
  JsonParse: (v) => v,
  randomPatternNameGenerator: () => 'test-name',
}));

vi.mock('@sistent/sistent', () => ({
  Typography: ({ children, color }) => (
    <span data-testid="typography" data-color={color}>
      {children}
    </span>
  ),
  Checkbox: () => <input type="checkbox" />,
  MenuItem: ({ children }) => <div>{children}</div>,
  ListItemText: ({ primary }) => <span>{primary}</span>,
  Select: ({ children }) => <div>{children}</div>,
  FormControl: ({ children }) => <div>{children}</div>,
  InputLabel: ({ children }) => <label>{children}</label>,
  OutlinedInput: () => <input />,
  Box: ({ children }) => <div>{children}</div>,
}));

// Import after mocks are set up
const { SelectConnection } = await import('./StepperContent');

describe('SelectConnection component', () => {
  const mockHandleNext = vi.fn();
  const mockSetSharedData = vi.fn();
  const mockRegisterConnection = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    mockUseVerifyAndRegisterConnectionMutation.mockReturnValue([mockRegisterConnection, {}]);
  });

  it('shows a loading state while fetching connection types', () => {
    mockUseGetComponentsQuery.mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
    });

    render(<SelectConnection setSharedData={mockSetSharedData} handleNext={mockHandleNext} />);

    expect(screen.getByText(/loading available connection types/i)).toBeInTheDocument();
  });

  it('shows an error state when the fetch fails', () => {
    mockUseGetComponentsQuery.mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: true,
    });

    render(<SelectConnection setSharedData={mockSetSharedData} handleNext={mockHandleNext} />);

    expect(screen.getByText(/failed to load connection types/i)).toBeInTheDocument();
  });

  it('renders cached connection types even if a refetch fails (isError is true)', () => {
    mockUseGetComponentsQuery.mockReturnValue({
      data: {
        components: [{ component: { kind: 'PrometheusConnection' } }],
      },
      isLoading: false,
      isError: true,
    });

    render(<SelectConnection setSharedData={mockSetSharedData} handleNext={mockHandleNext} />);

    expect(screen.queryByText(/failed to load connection types/i)).not.toBeInTheDocument();
    expect(screen.getByTestId('option-Prometheus Connection')).toBeInTheDocument();
  });

  it('shows an empty state when no connection types are available', () => {
    mockUseGetComponentsQuery.mockReturnValue({
      data: { components: [] },
      isLoading: false,
      isError: false,
    });

    render(<SelectConnection setSharedData={mockSetSharedData} handleNext={mockHandleNext} />);

    expect(screen.getByText(/no connection types are currently available/i)).toBeInTheDocument();
  });

  it('renders connection types dynamically from registry data', () => {
    mockUseGetComponentsQuery.mockReturnValue({
      data: {
        components: [
          { component: { kind: 'PrometheusConnection' } },
          { component: { kind: 'GrafanaConnection' } },
        ],
      },
      isLoading: false,
      isError: false,
    });

    render(<SelectConnection setSharedData={mockSetSharedData} handleNext={mockHandleNext} />);

    expect(screen.getByTestId('option-Prometheus Connection')).toBeInTheDocument();
    expect(screen.getByTestId('option-Grafana Connection')).toBeInTheDocument();
  });

  it('passes stable kind (not a sliced label) to the registration payload', async () => {
    const user = userEvent.setup();

    // RTK Query mutations return an object with an `.unwrap()` method.
    mockRegisterConnection.mockReturnValue({
      unwrap: () =>
        Promise.resolve({
          connection: { schema: '{}', id: 'conn-1' },
          credential: { schema: '{}' },
        }),
    });

    mockUseGetComponentsQuery.mockReturnValue({
      data: {
        components: [{ component: { kind: 'PrometheusConnection' } }],
      },
      isLoading: false,
      isError: false,
    });

    render(<SelectConnection setSharedData={mockSetSharedData} handleNext={mockHandleNext} />);

    await user.click(screen.getByTestId('option-Prometheus Connection'));

    await waitFor(() => {
      expect(mockRegisterConnection).toHaveBeenCalledWith({
        body: {
          kind: 'Prometheus', // stable kind, NOT "Prometheus Connection" or a sliced label
          status: 'initialize',
        },
      });
    });
  });

  it('automatically renders a newly registered connection type without frontend changes', () => {
    // Simulates Jaeger being registered in the backend registry
    mockUseGetComponentsQuery.mockReturnValue({
      data: {
        components: [
          { component: { kind: 'PrometheusConnection' } },
          { component: { kind: 'GrafanaConnection' } },
          { component: { kind: 'JaegerConnection' } }, // new — no frontend change needed
        ],
      },
      isLoading: false,
      isError: false,
    });

    render(<SelectConnection setSharedData={mockSetSharedData} handleNext={mockHandleNext} />);

    expect(screen.getByTestId('option-Jaeger Connection')).toBeInTheDocument();
  });

  it('does not render non-connection components (e.g. credentials)', () => {
    mockUseGetComponentsQuery.mockReturnValue({
      data: {
        components: [
          { component: { kind: 'PrometheusConnection' } },
          { component: { kind: 'PrometheusCredential' } }, // should NOT appear
        ],
      },
      isLoading: false,
      isError: false,
    });

    render(<SelectConnection setSharedData={mockSetSharedData} handleNext={mockHandleNext} />);

    expect(screen.getByTestId('option-Prometheus Connection')).toBeInTheDocument();
    expect(screen.queryByTestId('option-PrometheusCredential')).not.toBeInTheDocument();
  });

  it('queries the components API with the correct search param', () => {
    mockUseGetComponentsQuery.mockReturnValue({
      data: { components: [] },
      isLoading: false,
      isError: false,
    });

    render(<SelectConnection setSharedData={mockSetSharedData} handleNext={mockHandleNext} />);

    expect(mockUseGetComponentsQuery).toHaveBeenCalledWith({
      params: { pagesize: 'all', search: 'Connection' },
    });
  });
});
