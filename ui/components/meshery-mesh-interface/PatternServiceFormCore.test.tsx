import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const useStateCBMock = vi.fn();

vi.mock('../../utils/hooks/useStateCB', () => ({
  default: (...args: any[]) => useStateCBMock(...args),
}));

vi.mock('./PatternService', () => ({
  default: (props: any) => (
    <div data-testid="pattern-service" data-type={props.type}>
      pattern-service
    </div>
  ),
}));

vi.mock('./helpers', () => ({
  getPatternAttributeName: () => 'attr-name',
  createPatternFromConfig: (config: any, namespace: string) => ({ config, namespace }),
}));

vi.mock('../../utils/utils', () => ({
  scrollToTop: vi.fn(),
}));

vi.mock('../performance/helper', () => ({
  generateUUID: () => 'fixed-uuid',
}));

import PatternServiceFormCore from './PatternServiceFormCore';

describe('PatternServiceFormCore', () => {
  beforeEach(() => {
    useStateCBMock.mockReset();
    // useStateCB returns [state, setState, getValueRef]
    useStateCBMock.mockReturnValue([{}, vi.fn(), vi.fn(() => ({}))]);
  });

  it('caches the rendered child for workload types', () => {
    const children = vi.fn().mockImplementation((SettingsForm) => (
      <div data-testid="child">
        <SettingsForm />
      </div>
    ));

    render(
      <PatternServiceFormCore
        schemaSet={{ workload: { properties: {} }, type: 'workload' }}
        onSubmit={vi.fn()}
        onDelete={vi.fn()}
        namespace="default"
      >
        {children}
      </PatternServiceFormCore>,
    );

    expect(screen.getByTestId('child')).toBeInTheDocument();
    expect(children).toHaveBeenCalled();
  });

  it('renders an addon variant using the addon branch', () => {
    const children = vi.fn((SettingsForm) => (
      <div data-testid="addon-child">
        <SettingsForm />
      </div>
    ));

    render(
      <PatternServiceFormCore
        schemaSet={{ workload: { properties: {} }, type: 'addon' }}
        onSubmit={vi.fn()}
        onDelete={vi.fn()}
        namespace="ns"
      >
        {children}
      </PatternServiceFormCore>,
    );

    expect(screen.getByTestId('addon-child')).toBeInTheDocument();
  });

  it('attaches imperative handlers when a reference is provided', () => {
    const reference: any = { current: null };
    const onSubmit = vi.fn();
    const onDelete = vi.fn();
    useStateCBMock.mockReturnValue([{ foo: 'bar' }, vi.fn(), vi.fn(() => ({ foo: 'bar' }))]);

    render(
      <PatternServiceFormCore
        schemaSet={{ workload: { properties: {} }, type: 'workload' }}
        onSubmit={onSubmit}
        onDelete={onDelete}
        namespace="ns"
        reference={reference}
      >
        {(SettingsForm: any) => (
          <div data-testid="child">
            <SettingsForm />
          </div>
        )}
      </PatternServiceFormCore>,
    );

    expect(typeof reference.current.submit).toBe('function');
    expect(typeof reference.current.delete).toBe('function');
    expect(typeof reference.current.getSettings).toBe('function');
    expect(reference.current.referKey).toBe('fixed-uuid');

    reference.current.submit((settings: any) => ({ settings }));
    expect(onSubmit).toHaveBeenCalled();

    reference.current.delete((settings: any) => ({ settings }));
    expect(onDelete).toHaveBeenCalled();

    expect(reference.current.getSettings()).toEqual({ foo: 'bar' });
  });
});
