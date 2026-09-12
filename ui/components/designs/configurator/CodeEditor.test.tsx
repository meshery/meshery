import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (_tag: any) => () => {
    const Styled = ({ children, ...rest }: any) => <div {...rest}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    Card: ({ children }: any) => <div data-testid="card">{children}</div>,
    CardContent: ({ children }: any) => <div>{children}</div>,
  };
});

vi.mock('../../general/CodeMirror', () => ({
  UnControlled: ({ value, onChange, onBlur, options }: any) => (
    <div>
      <textarea
        data-testid="codemirror"
        defaultValue={value}
        data-mode={options?.mode || ''}
        onChange={(e) => onChange('editor', { data: 1 }, e.target.value)}
        onBlur={() => onBlur('editor')}
      />
    </div>
  ),
}));

import CodeEditor from './CodeEditor';

describe('CodeEditor', () => {
  it('renders the supplied yaml inside the editor', () => {
    render(<CodeEditor yaml="foo: bar" saveCodeEditorChanges={vi.fn()} onChange={vi.fn()} />);
    const editor = screen.getByTestId('codemirror') as HTMLTextAreaElement;
    expect(editor.value).toBe('foo: bar');
    expect(editor).toHaveAttribute('data-mode', 'text/x-yaml');
  });

  it('invokes onChange with the editor change payload', () => {
    const onChange = vi.fn();
    render(<CodeEditor yaml="" saveCodeEditorChanges={vi.fn()} onChange={onChange} />);
    const editor = screen.getByTestId('codemirror') as HTMLTextAreaElement;

    fireEvent.change(editor, { target: { value: 'updated: yaml' } });
    expect(onChange).toHaveBeenCalledWith('editor', { data: 1 }, 'updated: yaml');
  });

  it('invokes saveCodeEditorChanges on blur', () => {
    const saveCodeEditorChanges = vi.fn();
    render(<CodeEditor yaml="" saveCodeEditorChanges={saveCodeEditorChanges} onChange={vi.fn()} />);

    fireEvent.blur(screen.getByTestId('codemirror'));
    expect(saveCodeEditorChanges).toHaveBeenCalledWith('editor');
  });
});
