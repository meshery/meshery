import React, { useContext } from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { SchemaContext } from '../schemaSet';

describe('SchemaContext', () => {
  it('provides the documented default value', () => {
    let captured: unknown;
    const Probe = () => {
      captured = useContext(SchemaContext);
      return null;
    };
    render(<Probe />);
    expect(captured).toEqual({ workloadTraitSet: null, meshWorkloads: null });
  });

  it('allows consumers to override the value via a provider', () => {
    const value = { workloadTraitSet: { id: 'a' }, meshWorkloads: ['x'] };

    const Probe = () => {
      const ctx = useContext(SchemaContext);
      return <span data-testid="probe">{JSON.stringify(ctx)}</span>;
    };

    render(
      <SchemaContext.Provider value={value as never}>
        <Probe />
      </SchemaContext.Provider>,
    );

    expect(JSON.parse(screen.getByTestId('probe').textContent || 'null')).toEqual(value);
  });
});
