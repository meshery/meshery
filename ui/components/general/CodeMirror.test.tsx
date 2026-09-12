import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

// Capture props passed to the inner ReactCodeMirror so we can poke its callbacks.
let lastReactCMProps: any = null;

vi.mock('next/dynamic', () => ({
  // next/dynamic accepts a loader fn and returns a component; we just return our stub.
  default: () => (props: any) => {
    lastReactCMProps = props;
    return (
      <div data-testid="codemirror">
        <button
          data-testid="trigger-change"
          onClick={() => props.onChange('next-value', { view: 'view-stub' })}
        >
          change
        </button>
        <button data-testid="trigger-blur" onClick={() => props.onBlur?.('blur-event')}>
          blur
        </button>
        <button data-testid="trigger-create" onClick={() => props.onCreateEditor?.('editor-view')}>
          create
        </button>
      </div>
    );
  },
}));

vi.mock('@uiw/codemirror-theme-material', () => ({
  material: {},
}));

vi.mock('@codemirror/lang-json', () => ({ json: () => 'json-ext' }));
vi.mock('@codemirror/lang-yaml', () => ({ yaml: () => 'yaml-ext' }));

vi.mock('@codemirror/view', () => ({
  EditorView: { lineWrapping: 'wrap-ext' },
}));

vi.mock('@codemirror/lint', () => ({
  lintGutter: () => 'lint-gutter',
  // linter receives a fn — just return a marker string
  linter: (fn: any) => ({ fn, marker: 'linter' }),
}));

vi.mock('js-yaml', () => ({
  default: {
    load: (doc: string) => {
      if (doc.includes('BAD')) throw new Error('bad yaml');
      return null;
    },
  },
}));

import CodeMirror, { Controlled, UnControlled } from './CodeMirror';

describe('CodeMirror', () => {
  beforeEach(() => {
    lastReactCMProps = null;
  });

  it('exports Controlled and UnControlled aliases that point at the same component', () => {
    expect(Controlled).toBe(CodeMirror);
    expect(UnControlled).toBe(CodeMirror);
  });

  it('coerces non-string values to strings', () => {
    render(<CodeMirror value={42 as any} />);
    expect(lastReactCMProps.value).toBe('42');
  });

  it('passes empty string for null/undefined values', () => {
    render(<CodeMirror value={null as any} />);
    expect(lastReactCMProps.value).toBe('');
  });

  it('configures readOnly via options', () => {
    render(<CodeMirror value="x" options={{ readOnly: true }} />);
    expect(lastReactCMProps.readOnly).toBe(true);
  });

  it('honors lineNumbers in options', () => {
    render(<CodeMirror value="x" options={{ lineNumbers: false }} />);
    expect(lastReactCMProps.basicSetup.lineNumbers).toBe(false);
  });

  it('adds language extension for YAML mode', () => {
    render(<CodeMirror value="x" options={{ mode: 'text/x-yaml' }} />);
    expect(lastReactCMProps.extensions).toContain('yaml-ext');
  });

  it('adds language extension for JSON mode', () => {
    render(<CodeMirror value="x" options={{ mode: 'application/json' }} />);
    expect(lastReactCMProps.extensions).toContain('json-ext');
  });

  it('enables lineWrapping when option is set', () => {
    render(<CodeMirror value="x" options={{ lineWrapping: true }} />);
    expect(lastReactCMProps.extensions).toContain('wrap-ext');
  });

  it('adds lint gutter and yaml linter when lint mode is yaml', () => {
    render(<CodeMirror value="x" options={{ lint: true, mode: 'text/x-yaml' }} />);
    const exts = lastReactCMProps.extensions;
    expect(exts).toContain('lint-gutter');
    expect(exts.some((e: any) => e?.marker === 'linter')).toBe(true);
  });

  it('forwards onChange / onBeforeChange callbacks with the new value', () => {
    const onChange = vi.fn();
    const onBefore = vi.fn();
    const { getByTestId } = render(
      <CodeMirror value="x" onChange={onChange} onBeforeChange={onBefore} />,
    );
    fireEvent.click(getByTestId('trigger-change'));
    expect(onChange).toHaveBeenCalledWith('view-stub', null, 'next-value');
    expect(onBefore).toHaveBeenCalledWith('view-stub', null, 'next-value');
  });

  it('forwards onBlur and stashes the editor view via onCreateEditor', () => {
    const onBlur = vi.fn();
    const editorDidMount = vi.fn();
    const { getByTestId } = render(
      <CodeMirror value="x" onBlur={onBlur} editorDidMount={editorDidMount} />,
    );
    fireEvent.click(getByTestId('trigger-create'));
    fireEvent.click(getByTestId('trigger-blur'));
    expect(editorDidMount).toHaveBeenCalledWith('editor-view');
    expect(onBlur).toHaveBeenCalledWith('editor-view', 'blur-event');
  });

  it('calls editorWillUnmount during teardown', () => {
    const editorWillUnmount = vi.fn();
    const { getByTestId, unmount } = render(
      <CodeMirror value="x" editorWillUnmount={editorWillUnmount} />,
    );
    fireEvent.click(getByTestId('trigger-create'));
    unmount();
    expect(editorWillUnmount).toHaveBeenCalled();
  });
});
