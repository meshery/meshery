import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const useMeshModelComponentsMock = vi.fn();

vi.mock('@sistent/sistent', () => ({
  Accordion: ({ children, expanded, onChange }: any) => (
    <div data-testid="accordion" data-expanded={String(!!expanded)}>
      <button data-testid="accordion-toggle" onClick={(e) => onChange?.(e, !expanded)}>
        toggle
      </button>
      {children}
    </div>
  ),
  AccordionDetails: ({ children }: any) => <div data-testid="accordion-details">{children}</div>,
  AccordionSummary: ({ children }: any) => <div data-testid="accordion-summary">{children}</div>,
  Box: ({ children }: any) => <div>{children}</div>,
  ExpandMoreIcon: () => <svg />,
  FormControl: ({ children }: any) => <div>{children}</div>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  MenuItem: ({ children, value, disabled }: any) => (
    <option value={value} disabled={!!disabled}>
      {children}
    </option>
  ),
  ModalButtonPrimary: ({ children, onClick }: any) => (
    <button data-testid="primary-btn" onClick={onClick}>
      {children}
    </button>
  ),
  ModalButtonSecondary: ({ children, onClick, color }: any) => (
    <button data-testid={`secondary-btn-${color || 'default'}`} onClick={onClick}>
      {children}
    </button>
  ),
  Tab: ({ label }: any) => <span>{label}</span>,
  Tabs: ({ children, value, onChange }: any) => (
    <div data-testid="tabs" data-value={value}>
      <button onClick={(e) => onChange?.(e, 0)} data-testid="tab-0">
        allow
      </button>
      <button onClick={(e) => onChange?.(e, 1)} data-testid="tab-1">
        deny
      </button>
      {children}
    </div>
  ),
  TextField: ({ value, onChange, label, children, select }: any) => (
    <div>
      <label>{label}</label>
      {select ? (
        <select data-testid={`select-${label}`} value={value || ''} onChange={onChange}>
          {children}
        </select>
      ) : (
        <input data-testid={`input-${label}`} value={value || ''} onChange={onChange} />
      )}
    </div>
  ),
  Typography: ({ children }: any) => <span>{children}</span>,
}));

vi.mock('../meshery-mesh-interface/PatternService/RJSF_wrapper', () => ({
  default: ({ formData, onChange }: any) => (
    <div data-testid="rjsf-wrapper">
      <button
        data-testid="trigger-match-change"
        onClick={() => onChange({ ...formData, match: { property: 'value' } })}
      >
        change-match
      </button>
    </div>
  ),
}));

vi.mock('@/utils/hooks/useMeshModelComponents', () => ({
  useMeshModelComponents: () => useMeshModelComponentsMock(),
}));

import SelectorsForm from './SelectorsForm';

describe('SelectorsForm', () => {
  const selectorsSchema = {
    items: {
      properties: {
        allow: {
          properties: {
            from: {
              items: {
                properties: {
                  kind: {},
                  model: {},
                  matchLabels: { type: 'string' },
                },
              },
            },
          },
        },
      },
    },
  };

  beforeEach(() => {
    useMeshModelComponentsMock.mockReturnValue({
      models: {
        Kubernetes: [{ name: 'pod', id: '1', displayName: 'Pod', registrant: { kind: 'native' } }],
      },
      categories: [{ name: 'Kubernetes' }, { name: 'Networking' }],
      getModelFromCategory: vi.fn(),
    });
  });

  it('renders the empty state message when no selectors exist', () => {
    render(<SelectorsForm selectorsSchema={selectorsSchema} formData={{}} onChange={vi.fn()} />);
    expect(screen.getAllByText(/No from selectors added yet/i).length).toBeGreaterThan(0);
  });

  it('adds a selector when the Add button is clicked', () => {
    const onChange = vi.fn();
    render(<SelectorsForm selectorsSchema={selectorsSchema} formData={{}} onChange={onChange} />);

    // The first "Add from" button is the from-direction add button
    const addButtons = screen.getAllByTestId('primary-btn');
    fireEvent.click(addButtons[0]);

    expect(onChange).toHaveBeenCalled();
    const lastCall = onChange.mock.calls.at(-1)![0];
    expect(lastCall.selectors.allow.from).toHaveLength(1);
  });

  it('switches to the Deny tab', () => {
    render(<SelectorsForm selectorsSchema={selectorsSchema} formData={{}} onChange={vi.fn()} />);

    fireEvent.click(screen.getByTestId('tab-1'));

    // We should now see "No from selectors added yet" under Deny tab
    expect(screen.getAllByText(/No from selectors added yet/i).length).toBeGreaterThan(0);
  });

  it('removes a selector when the Remove button is clicked', () => {
    const onChange = vi.fn();
    render(
      <SelectorsForm
        selectorsSchema={selectorsSchema}
        formData={{
          selectors: {
            allow: {
              from: [{ kind: 'ClusterRole', model: { name: '', registrant: { kind: '' } } }],
              to: [],
            },
            deny: { from: [], to: [] },
          },
        }}
        onChange={onChange}
      />,
    );

    const removeButtons = screen.getAllByTestId('secondary-btn-error');
    fireEvent.click(removeButtons[0]);

    expect(onChange).toHaveBeenCalled();
    const lastCall = onChange.mock.calls.at(-1)![0];
    expect(lastCall.selectors.allow.from).toHaveLength(0);
  });

  it('updates a selector kind via TextField change', () => {
    const onChange = vi.fn();
    render(
      <SelectorsForm
        selectorsSchema={selectorsSchema}
        formData={{
          selectors: {
            allow: { from: [{ kind: '', model: { name: '', registrant: { kind: '' } } }], to: [] },
            deny: { from: [], to: [] },
          },
        }}
        onChange={onChange}
      />,
    );

    const kindInputs = screen.getAllByTestId('input-Kind');
    fireEvent.change(kindInputs[0], { target: { value: 'Role' } });

    expect(onChange).toHaveBeenCalled();
    const lastCall = onChange.mock.calls.at(-1)![0];
    expect(lastCall.selectors.allow.from[0].kind).toBe('Role');
  });

  it('updates match settings via the embedded RJSF wrapper', () => {
    const onChange = vi.fn();
    render(
      <SelectorsForm
        selectorsSchema={selectorsSchema}
        formData={{
          selectors: {
            allow: { from: [{ kind: '', model: { name: '', registrant: { kind: '' } } }], to: [] },
            deny: { from: [], to: [] },
          },
        }}
        onChange={onChange}
      />,
    );

    const triggers = screen.getAllByTestId('trigger-match-change');
    fireEvent.click(triggers[0]);

    expect(onChange).toHaveBeenCalled();
  });
});
