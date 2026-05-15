import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ImportModal from './ImportModal';

let lastModalProps: any = null;

vi.mock('./Modal', () => ({
  default: (props: any) => {
    lastModalProps = props;
    return (
      <div data-testid="rjsf-modal" data-title={props.title} data-submit-btn={props.submitBtnText}>
        {props.children}
      </div>
    );
  },
}));

vi.mock('@sistent/sistent', () => ({
  Button: ({ children, onClick, fullWidth, variant, color, ...rest }: any) => (
    <button
      onClick={onClick}
      data-full-width={String(fullWidth)}
      data-variant={variant}
      data-color={color}
      {...rest}
    >
      {children}
    </button>
  ),
}));

describe('ImportModal', () => {
  beforeEach(() => {
    lastModalProps = null;
  });

  it('renders an RJSFModal with the Import title for the given type', () => {
    render(
      <ImportModal
        importType="design"
        rjsfSchema={{ title: 'schema' }}
        uiSchema={{}}
        handleSubmit={vi.fn()}
        handleClose={vi.fn()}
      />,
    );

    expect(screen.getByTestId('rjsf-modal')).toHaveAttribute('data-title', 'Import Design');
    expect(screen.getByTestId('rjsf-modal')).toHaveAttribute('data-submit-btn', 'Import Design');
  });

  it('renders the inner button labelled with the capitalized type', () => {
    render(
      <ImportModal
        importType="application"
        rjsfSchema={{}}
        uiSchema={{}}
        handleSubmit={vi.fn()}
        handleClose={vi.fn()}
      />,
    );

    expect(screen.getByRole('button', { name: /Import Application/ })).toBeInTheDocument();
  });

  it('calls handleClose and handleSubmit when the inner button is clicked', async () => {
    const user = userEvent.setup();
    const handleSubmit = vi.fn();
    const handleClose = vi.fn();

    render(
      <ImportModal
        importType="filter"
        rjsfSchema={{}}
        uiSchema={{}}
        handleSubmit={handleSubmit}
        handleClose={handleClose}
      />,
    );

    await user.click(screen.getByRole('button', { name: /Import Filter/ }));

    expect(handleClose).toHaveBeenCalledTimes(1);
    expect(handleSubmit).toHaveBeenCalledTimes(1);
  });

  it('defaults uiSchema to an empty object when not provided', () => {
    render(
      <ImportModal
        importType="design"
        rjsfSchema={{}}
        uiSchema={undefined}
        handleSubmit={vi.fn()}
        handleClose={vi.fn()}
      />,
    );

    expect(lastModalProps.uiSchema).toEqual({});
  });

  it('passes the rjsf schema and import type through', () => {
    const schema = { type: 'object', properties: { name: { type: 'string' } } };
    render(
      <ImportModal
        importType="design"
        rjsfSchema={schema}
        uiSchema={{ x: 1 }}
        handleSubmit={vi.fn()}
        handleClose={vi.fn()}
      />,
    );

    expect(lastModalProps.schema).toBe(schema);
    expect(lastModalProps.uiSchema).toEqual({ x: 1 });
    expect(lastModalProps.type).toBe('design');
    expect(lastModalProps.formData).toEqual({});
  });
});
