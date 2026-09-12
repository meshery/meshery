import React, { isValidElement } from 'react';
import { describe, it, expect } from 'vitest';
import { buildCustomInputProps } from '../buildCustomInputProps';

const leading = <span data-testid="leading">filter</span>;

describe('buildCustomInputProps', () => {
  // Regression: NotificationCenter threw "undefined is not an object
  // (evaluating 'e.InputProps.startAdornment')" because params.InputProps
  // could be undefined when MUI's Autocomplete renderInput callback fired
  // before InputProps was wired through.
  it('does not throw when params.InputProps is undefined', () => {
    expect(() => buildCustomInputProps(undefined, leading)).not.toThrow();
    const result = buildCustomInputProps(undefined, leading);
    expect(result).toBeDefined();
    expect(isValidElement(result.startAdornment)).toBe(true);
  });

  it('preserves existing startAdornment after the leading adornment', () => {
    const existing = <span data-testid="existing">existing</span>;
    const result = buildCustomInputProps({ startAdornment: existing }, leading);
    const fragment = result.startAdornment as React.ReactElement;
    expect(isValidElement(fragment)).toBe(true);
    const children = (fragment.props as { children: React.ReactNode[] }).children;
    expect(children[0]).toBe(leading);
    expect(children[1]).toBe(existing);
  });

  it('forwards arbitrary InputProps fields untouched', () => {
    const ref = React.createRef<HTMLDivElement>();
    const result = buildCustomInputProps(
      { ref, className: 'foo', startAdornment: undefined },
      leading,
    );
    expect(result.ref).toBe(ref);
    expect(result.className).toBe('foo');
  });

  it('handles an empty InputProps object', () => {
    expect(() => buildCustomInputProps({}, leading)).not.toThrow();
    const result = buildCustomInputProps({}, leading);
    expect(isValidElement(result.startAdornment)).toBe(true);
  });
});
