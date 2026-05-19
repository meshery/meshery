import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const rjsfWrapperMock = vi.fn();

vi.mock('./RJSF_wrapper', () => ({
  default: (props: any) => {
    rjsfWrapperMock(props);
    return (
      <div
        data-testid="rjsf-wrapper"
        data-hide-title={String(props.hideTitle)}
        data-hide-submit={String(props.hideSubmit)}
      />
    );
  },
}));

import PatternService from './index';

describe('PatternService', () => {
  beforeEach(() => {
    rjsfWrapperMock.mockClear();
  });

  it('returns null when the schema has no properties', () => {
    const { container } = render(
      <PatternService
        type="workload"
        jsonSchema={{ properties: {} }}
        onChange={vi.fn()}
        formData={{}}
      />,
    );
    expect(container.firstChild).toBeNull();
  });

  it('renders the RJSFWrapper when the schema has at least one property', () => {
    render(
      <PatternService
        type="workload"
        jsonSchema={{ properties: { name: { type: 'string' } } }}
        onChange={vi.fn()}
        formData={{ name: 'foo' }}
      />,
    );
    expect(screen.getByTestId('rjsf-wrapper')).toBeInTheDocument();
    expect(screen.getByTestId('rjsf-wrapper')).toHaveAttribute('data-hide-title', 'true');
    expect(screen.getByTestId('rjsf-wrapper')).toHaveAttribute('data-hide-submit', 'false');
  });

  it('hides the submit button for trait types', () => {
    render(
      <PatternService
        type="trait"
        jsonSchema={{ properties: { name: { type: 'string' } } }}
        onChange={vi.fn()}
        formData={{}}
      />,
    );
    expect(screen.getByTestId('rjsf-wrapper')).toHaveAttribute('data-hide-submit', 'true');
    expect(screen.getByTestId('rjsf-wrapper')).toHaveAttribute('data-hide-title', 'false');
  });
});
