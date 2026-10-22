import React from 'react';
import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { render } from '@testing-library/react';
import { hasClass, AddClassRecursively } from '../Elements';

describe('hasClass', () => {
  it('returns true when the element itself has the class', () => {
    const el = { className: 'foo bar baz', parentElement: null };
    expect(hasClass(el, 'bar')).toBe(true);
  });

  it('returns true when any ancestor has the class', () => {
    const grandparent = { className: 'has-target', parentElement: null };
    const parent = { className: 'middle', parentElement: grandparent };
    const child = { className: 'leaf', parentElement: parent };
    expect(hasClass(child, 'has-target')).toBe(true);
  });

  it('returns false when neither the element nor its ancestors have the class', () => {
    const parent = { className: 'foo', parentElement: null };
    const child = { className: 'bar', parentElement: parent };
    expect(hasClass(child, 'baz')).toBe(false);
  });

  it('returns false for null / undefined element', () => {
    expect(hasClass(null, 'whatever')).toBe(false);
    expect(hasClass(undefined, 'whatever')).toBe(false);
  });

  it('does not throw if className is not a string (SVG-like DOM nodes)', () => {
    const svgLike = { className: { baseVal: 'foo' }, parentElement: null };
    // The function only matches when className is a string, so this should return false
    // but should not throw.
    expect(() => hasClass(svgLike, 'foo')).not.toThrow();
    expect(hasClass(svgLike, 'foo')).toBe(false);
  });

  it('logs and swallows errors when traversal raises', () => {
    const spy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const trap = {
      get className() {
        throw new Error('boom');
      },
      parentElement: null,
    };
    expect(hasClass(trap, 'foo')).toBe(false);
    expect(spy).toHaveBeenCalled();
    spy.mockRestore();
  });
});

describe('AddClassRecursively', () => {
  beforeEach(() => {
    vi.spyOn(console, 'error').mockImplementation(() => {});
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('appends the className to every nested valid React element', () => {
    const { container } = render(
      <AddClassRecursively className="injected">
        <div className="outer">
          <span className="inner">hello</span>
        </div>
      </AddClassRecursively>,
    );

    const outer = container.querySelector('div')!;
    const inner = container.querySelector('span')!;

    expect(outer.className).toContain('outer');
    expect(outer.className).toContain('injected');
    expect(inner.className).toContain('inner');
    expect(inner.className).toContain('injected');
  });

  it('preserves non-element children (strings, numbers) untouched', () => {
    const { container } = render(
      <AddClassRecursively className="x">
        <div>
          some text
          <span>child</span>
        </div>
      </AddClassRecursively>,
    );

    expect(container.textContent).toContain('some text');
    expect(container.textContent).toContain('child');
  });

  it('handles elements with no existing className', () => {
    const { container } = render(
      <AddClassRecursively className="extra">
        <div>
          <span>leaf</span>
        </div>
      </AddClassRecursively>,
    );

    const div = container.querySelector('div')!;
    expect(div.className).toContain('extra');
  });
});
