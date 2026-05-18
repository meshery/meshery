import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@rjsf/utils', () => ({
  ADDITIONAL_PROPERTY_FLAG: '__ADDITIONAL_PROPERTY__',
}));

vi.mock('@sistent/sistent', () => ({
  DeleteIcon: () => <svg data-testid="delete-icon" />,
  FormControl: ({ children }: any) => <div data-testid="form-control">{children}</div>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  IconButton: ({ children, onClick, disabled }: any) => (
    <button data-testid="delete-btn" onClick={onClick} disabled={!!disabled}>
      {children}
    </button>
  ),
  Input: ({ defaultValue, onChange, disabled, id }: any) => (
    <input
      data-testid="key-input"
      id={id}
      defaultValue={defaultValue}
      disabled={!!disabled}
      onChange={onChange}
    />
  ),
  InputLabel: ({ children }: any) => <label>{children}</label>,
}));

vi.mock('../../../../css/icons.styles', () => ({ iconMedium: {} }));

import WrapIfAdditionalTemplate from './WrapIfAdditionalTemplate';

describe('WrapIfAdditionalTemplate', () => {
  it('returns children unchanged when the schema is not additional', () => {
    render(
      <WrapIfAdditionalTemplate
        children={<div data-testid="child">child</div>}
        classNames="classy"
        id="id1"
        label="key1"
        onDropPropertyClick={() => vi.fn()}
        onKeyChange={vi.fn()}
        schema={{}}
      />,
    );
    expect(screen.getByTestId('child')).toBeInTheDocument();
    expect(screen.queryByTestId('key-input')).not.toBeInTheDocument();
  });

  it('renders a key editor and delete button when the schema has the additional flag', () => {
    const onDrop = vi.fn();
    render(
      <WrapIfAdditionalTemplate
        children={<div>child</div>}
        classNames="cls"
        id="id1"
        label="key1"
        onDropPropertyClick={(label: string) => () => onDrop(label)}
        onKeyChange={vi.fn()}
        schema={{ __ADDITIONAL_PROPERTY__: true }}
      />,
    );
    expect(screen.getByTestId('key-input')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('delete-btn'));
    expect(onDrop).toHaveBeenCalledWith('key1');
  });

  it('strips the default newKey placeholder', () => {
    render(
      <WrapIfAdditionalTemplate
        children={<div>child</div>}
        classNames="cls"
        id="id1"
        label="newKey-foo"
        onDropPropertyClick={() => vi.fn()}
        onKeyChange={vi.fn()}
        schema={{ __ADDITIONAL_PROPERTY__: true }}
      />,
    );
    expect(screen.getByTestId('key-input')).toHaveValue('');
  });

  it('forwards key edits to onKeyChange', () => {
    const onKeyChange = vi.fn();
    render(
      <WrapIfAdditionalTemplate
        children={<div>child</div>}
        classNames="cls"
        id="id1"
        label="key1"
        onDropPropertyClick={() => vi.fn()}
        onKeyChange={onKeyChange}
        schema={{ __ADDITIONAL_PROPERTY__: true }}
      />,
    );
    fireEvent.change(screen.getByTestId('key-input'), { target: { value: 'newKey' } });
    expect(onKeyChange).toHaveBeenCalledWith('newKey');
  });

  it('disables the input and delete button when readonly is true', () => {
    render(
      <WrapIfAdditionalTemplate
        children={<div>child</div>}
        classNames="cls"
        id="id1"
        label="key1"
        readonly
        onDropPropertyClick={() => vi.fn()}
        onKeyChange={vi.fn()}
        schema={{ __ADDITIONAL_PROPERTY__: true }}
      />,
    );
    expect(screen.getByTestId('key-input')).toBeDisabled();
    expect(screen.getByTestId('delete-btn')).toBeDisabled();
  });
});
