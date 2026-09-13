import React from 'react';
import { render, screen, act } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const rjsfFormMock = vi.fn();
const notifyMock = vi.fn();

vi.mock('../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: notifyMock }),
}));

vi.mock('../helpers', () => ({
  buildUiSchema: () => ({ 'ui:order': [] }),
}));

vi.mock('./helper', () => ({
  getRefinedJsonSchema: (schema: any) => ({ ...schema, _refined: true }),
}));

vi.mock('./RJSF', () => ({
  default: (props: any) => {
    rjsfFormMock(props);
    return (
      <div
        data-testid="rjsf-form"
        data-is-loading={String(props.isLoading)}
        data-live-validate={String(props.liveValidate)}
      />
    );
  },
}));

import RJSFWrapper from './RJSF_wrapper';

describe('RJSFWrapper', () => {
  beforeEach(() => {
    rjsfFormMock.mockClear();
    notifyMock.mockClear();
  });

  it('renders the RJSF form initially as loading', () => {
    render(
      <RJSFWrapper
        jsonSchema={{ type: 'object', properties: {} }}
        formData={{}}
        onChange={vi.fn()}
      />,
    );
    expect(screen.getByTestId('rjsf-form')).toHaveAttribute('data-is-loading', 'true');
  });

  it('forwards live-validate from props', () => {
    render(
      <RJSFWrapper
        jsonSchema={{ type: 'object', properties: {} }}
        formData={{}}
        liveValidate={false}
      />,
    );
    expect(screen.getByTestId('rjsf-form')).toHaveAttribute('data-live-validate', 'false');
  });

  it('wraps the form with RJSFWrapperComponent when provided', () => {
    const Wrapper = ({ children }: any) => <div data-testid="custom-wrapper">{children}</div>;
    render(
      <RJSFWrapper
        jsonSchema={{ type: 'object', properties: {} }}
        formData={{}}
        onChange={vi.fn()}
        RJSFWrapperComponent={Wrapper as any}
      />,
    );
    expect(screen.getByTestId('custom-wrapper')).toBeInTheDocument();
  });

  it('calls onChange after the debounce timer fires', async () => {
    vi.useFakeTimers();
    const onChange = vi.fn();
    render(
      <RJSFWrapper
        jsonSchema={{ type: 'object', properties: {} }}
        formData={{ name: 'a' }}
        onChange={onChange}
      />,
    );
    await act(async () => {
      vi.advanceTimersByTime(500);
    });
    expect(onChange).toHaveBeenCalled();
    vi.useRealTimers();
  });
});
