import React from 'react';
import { act, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (_tag: any) => () => {
    const Styled = ({ children, onClick, value, onChange, autoFocus }: any) => (
      <div onClick={onClick}>
        {typeof value === 'string' ? (
          <input
            value={value}
            onChange={onChange}
            autoFocus={autoFocus}
            data-testid="breadcrumb-input"
          />
        ) : (
          children
        )}
      </div>
    );
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    ClickAwayListener: ({ children, onClickAway }: any) => (
      <div>
        {children}
        <button type="button" onClick={() => onClickAway && onClickAway({})}>
          click-away
        </button>
      </div>
    ),
  };
});

import CustomBreadCrumb from './CustomBreadCrumb';

describe('CustomBreadCrumb', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('renders the title and "Designs" link by default', () => {
    render(<CustomBreadCrumb title="My Design" onBack={vi.fn()} titleChangeHandler={vi.fn()} />);
    expect(screen.getByText('Designs')).toBeInTheDocument();
    expect(screen.getByText('My Design')).toBeInTheDocument();
  });

  it('invokes onBack when "Designs" is clicked', () => {
    const onBack = vi.fn();
    render(<CustomBreadCrumb title="My Design" onBack={onBack} titleChangeHandler={vi.fn()} />);
    fireEvent.click(screen.getByText('Designs'));
    expect(onBack).toHaveBeenCalled();
  });

  it('shows the input editor when the title is clicked, and propagates changes after debounce', () => {
    const titleChangeHandler = vi.fn();
    render(
      <CustomBreadCrumb title="Old" onBack={vi.fn()} titleChangeHandler={titleChangeHandler} />,
    );

    // Click title to enter edit mode.
    fireEvent.click(screen.getByText('Old'));
    const input = screen.getByTestId('breadcrumb-input') as HTMLInputElement;
    expect(input.value).toBe('Old');

    fireEvent.change(input, { target: { value: 'New' } });
    expect(input.value).toBe('New');

    // titleChangeHandler should be debounced (400ms).
    act(() => {
      vi.advanceTimersByTime(400);
    });
    expect(titleChangeHandler).toHaveBeenLastCalledWith('New');
  });

  it('exits edit mode when ClickAwayListener fires', () => {
    render(<CustomBreadCrumb title="Hello" onBack={vi.fn()} titleChangeHandler={vi.fn()} />);

    fireEvent.click(screen.getByText('Hello'));
    expect(screen.getByTestId('breadcrumb-input')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'click-away' }));
    expect(screen.queryByTestId('breadcrumb-input')).not.toBeInTheDocument();
    expect(screen.getByText('Hello')).toBeInTheDocument();
  });
});
