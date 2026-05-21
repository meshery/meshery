import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const getTemplateMock = vi.fn();

vi.mock('@rjsf/utils', () => ({
  getTemplate: (...args: any[]) => getTemplateMock(...args),
  getUiOptions: (uiSchema: any) => uiSchema || {},
}));

import { CustomFieldTemplate } from './FieldTemplate';

describe('CustomFieldTemplate', () => {
  beforeEach(() => {
    getTemplateMock.mockReset();
  });

  it('returns the children alone when no WrapIfAdditionalTemplate is registered', () => {
    getTemplateMock.mockReturnValue(undefined);
    render(
      <CustomFieldTemplate
        children={<div data-testid="child">hi</div>}
        registry={{} as any}
        uiSchema={{}}
        schema={{}}
      />,
    );
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('delegates to the registered template when available', () => {
    const FakeWrap = ({ children }: any) => <div data-testid="wrap">{children}</div>;
    getTemplateMock.mockReturnValue(FakeWrap);
    render(
      <CustomFieldTemplate
        children={<div data-testid="child">hi</div>}
        registry={{} as any}
        uiSchema={{}}
        schema={{}}
      />,
    );
    expect(screen.getByTestId('wrap')).toBeInTheDocument();
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });
});
