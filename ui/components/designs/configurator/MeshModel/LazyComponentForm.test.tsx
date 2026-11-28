import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const getMeshModelComponentMock = vi.fn();
const formMounts = vi.fn();
const keyCounter = vi.hoisted(() => ({ value: 0 }));

vi.mock('@sistent/sistent', () => {
  const styled = (_tag: any) => () => {
    const Styled = ({ children, ...rest }: any) => <div {...rest}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    Typography: ({ children }: any) => <span>{children}</span>,
    CircularProgress: () => <div data-testid="loading" />,
    ExpandMoreIcon: () => null,
    Accordion: ({ children, expanded, onChange }: any) => (
      <div data-testid="accordion" data-expanded={String(expanded)}>
        <button data-testid="accordion-toggle" onClick={onChange} />
        {children}
      </div>
    ),
    AccordionSummary: ({ children }: any) => <div>{children}</div>,
    AccordionDetails: ({ children }: any) => <div>{children}</div>,
  };
});

vi.mock('../../../../api/meshmodel', () => ({
  getMeshModelComponent: (...args: any[]) => getMeshModelComponentMock(...args),
}));

vi.mock('../../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('../../../../css/icons.styles', () => ({ iconMedium: {} }));

vi.mock('../../../../lib/event-types', () => ({ EVENT_TYPES: { ERROR: 'error' } }));

// Mirrors PatternServiceFormCore: a fresh reference key is generated on every mount.
vi.mock('../../../meshery-mesh-interface/PatternServiceForm', () => ({
  default: function PatternServiceFormMock({ reference, onSettingsChange }: any) {
    const [referKey] = React.useState(() => `key-${++keyCounter.value}`);
    React.useEffect(() => {
      formMounts(referKey);
    }, []);
    if (!reference.current) reference.current = {};
    reference.current.referKey = referKey;
    return (
      <button data-testid="form" onClick={() => onSettingsChange?.({ name: 'edited' })}>
        {referKey}
      </button>
    );
  },
}));

vi.mock('@/rtk-query/design', () => ({
  useUpdatePatternFileMutation: () => [vi.fn()],
  useDeletePatternFileMutation: () => [vi.fn()],
  useLazyGetDesignQuery: () => [vi.fn()],
}));

import LazyComponentForm from './LazyComponentForm';
import useDesignLifecycle from './hooks/useDesignLifecycle';

const makeComponent = (kind: string) => ({
  displayName: kind,
  component: { kind, version: 'apps/v1' },
  model: { name: 'kubernetes', model: { version: 'v1.0.0' } },
});

describe('LazyComponentForm', () => {
  beforeEach(() => {
    formMounts.mockReset();
    keyCounter.value = 0;
    getMeshModelComponentMock.mockReset();
    getMeshModelComponentMock.mockResolvedValue({
      components: [{ component: { schema: '{"properties":{}}' } }],
    });
  });

  it('requests the schema for the definition api version recorded in the design', async () => {
    // A differing top-level apiVersion must not change the requested schema version.
    const deployment = {
      ...makeComponent('Deployment'),
      component: { kind: 'Deployment', version: 'apps/v1beta1' },
      apiVersion: 'apps/v1',
    };
    render(<LazyComponentForm component={deployment} />);

    fireEvent.click(screen.getByTestId('accordion-toggle'));

    await waitFor(() =>
      expect(getMeshModelComponentMock).toHaveBeenCalledWith(
        'kubernetes',
        'Deployment',
        'v1.0.0',
        'apps/v1beta1',
      ),
    );
  });

  it('keeps the same form (and component id) when collapsed and re-expanded', async () => {
    const handler = vi.fn();
    const createSettingsChangeHandler = vi.fn(
      (reference) => (formData: any) => handler(reference.current.referKey, formData),
    );
    render(
      <LazyComponentForm
        component={makeComponent('Namespace')}
        createSettingsChangeHandler={createSettingsChangeHandler}
      />,
    );

    const toggle = screen.getByTestId('accordion-toggle');
    fireEvent.click(toggle);
    await screen.findByTestId('form');

    fireEvent.click(toggle); // collapse
    fireEvent.click(toggle); // expand again
    fireEvent.click(screen.getByTestId('form'));

    expect(formMounts).toHaveBeenCalledTimes(1);
    expect(handler).toHaveBeenCalledWith('key-1', { name: 'edited' });
    // The settings handler is created once, not on every render.
    expect(createSettingsChangeHandler).toHaveBeenCalledTimes(1);
  });

  it('gives every component form its own reference', async () => {
    const references: any[] = [];
    const createSettingsChangeHandler = (reference: any) => {
      if (!references.includes(reference)) references.push(reference);
      return vi.fn();
    };
    render(
      <>
        <LazyComponentForm
          component={makeComponent('Namespace')}
          createSettingsChangeHandler={createSettingsChangeHandler}
        />
        <LazyComponentForm
          component={makeComponent('Service')}
          createSettingsChangeHandler={createSettingsChangeHandler}
        />
      </>,
    );

    const [first, second] = screen.getAllByTestId('accordion-toggle');
    fireEvent.click(first);
    fireEvent.click(second);
    const forms = await screen.findAllByTestId('form');

    expect(forms).toHaveLength(2);
    expect(references).toHaveLength(2);
    expect(references[0].current.referKey).not.toBe(references[1].current.referKey);
  });

  it('routes each mounted form edit to its own component in the design', async () => {
    let design: any;
    function Harness() {
      const { onSettingsChange, designJson } = useDesignLifecycle();
      design = designJson;
      return (
        <>
          {['Namespace', 'Service'].map((kind) => {
            const definition = makeComponent(kind);
            return (
              <LazyComponentForm
                key={kind}
                component={definition}
                createSettingsChangeHandler={(reference) => onSettingsChange(definition, reference)}
              />
            );
          })}
        </>
      );
    }
    render(<Harness />);

    const [first, second] = screen.getAllByTestId('accordion-toggle');
    fireEvent.click(first);
    fireEvent.click(second);
    const [namespaceForm, serviceForm] = await screen.findAllByTestId('form');

    fireEvent.click(namespaceForm);
    fireEvent.click(serviceForm);
    fireEvent.click(namespaceForm);

    await waitFor(() => expect(design.components).toHaveLength(2));
    const byKind = Object.fromEntries(design.components.map((c: any) => [c.component.kind, c]));
    expect(byKind.Namespace.id).toBe(namespaceForm.textContent);
    expect(byKind.Service.id).toBe(serviceForm.textContent);
    expect(byKind.Namespace.id).not.toBe(byKind.Service.id);
  });
});
