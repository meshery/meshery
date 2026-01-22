import React, { useMemo, useRef, useEffect } from 'react';
import dynamic from 'next/dynamic';
import { material } from '@uiw/codemirror-theme-material';
import { json } from '@codemirror/lang-json';
import { yaml } from '@codemirror/lang-yaml';
import { EditorView } from '@codemirror/view';
import { lintGutter, linter } from '@codemirror/lint';
import jsyaml from 'js-yaml';

const ReactCodeMirror = dynamic(() => import('@uiw/react-codemirror').then((mod) => mod.default), {
  ssr: false,
});

const getSafeValue = (value) => {
  if (typeof value === 'string') {
    return value;
  }
  if (value) {
    return String(value);
  }
  return '';
};

const getDiagnosticRange = (pos, docLength) => {
  const safePos = Math.max(0, Math.min(pos, docLength));
  const end = Math.min(safePos + 1, docLength);
  return { from: safePos, to: end };
};

const yamlLinter = linter((view) => {
  const doc = view.state.doc.toString();
  try {
    jsyaml.load(doc);
    return [];
  } catch (error) {
    const pos = error?.mark?.position ?? 0;
    const { from, to } = getDiagnosticRange(pos, doc.length);
    return [
      {
        from,
        to,
        severity: 'error',
        message: error?.reason || error?.message || 'Invalid YAML',
      },
    ];
  }
});

const jsonLinter = linter((view) => {
  const doc = view.state.doc.toString();
  try {
    JSON.parse(doc);
    return [];
  } catch (error) {
    const match = /position (\d+)/i.exec(error?.message || '');
    const pos = match ? Number(match[1]) : 0;
    const { from, to } = getDiagnosticRange(pos, doc.length);
    return [
      {
        from,
        to,
        severity: 'error',
        message: error?.message || 'Invalid JSON',
      },
    ];
  }
});

const getLanguageExtension = (mode) => {
  switch (mode) {
    case 'text/x-yaml':
      return yaml();
    case 'application/json':
      return json();
    default:
      return null;
  }
};

type CodeMirrorOptions = {
  mode?: string;
  lineWrapping?: boolean;
  [key: string]: any;
};

type CodeMirrorProps = {
  value?: string;
  options?: CodeMirrorOptions;
  onBeforeChange?: (_editor: any, _data: any, _value: string) => void;
  onChange?: (_value: string, _viewUpdate: any) => void;
  onBlur?: (_editor: any, _event: any) => void;
  editorDidMount?: (_editor: any, _value?: string) => void;
  editorWillUnmount?: (_editor: any) => void;
  [key: string]: any;
};

const CodeMirror = ({
  value,
  options = {},
  onBeforeChange,
  onChange,
  onBlur,
  editorDidMount,
  editorWillUnmount,
  ...rest
}: CodeMirrorProps) => {
  const editorViewRef = useRef<any>(null);

  useEffect(() => {
    return () => {
      if (editorWillUnmount) {
        editorWillUnmount(editorViewRef.current);
      }
    };
  }, [editorWillUnmount]);

  const extensions = useMemo(() => {
    const editorExtensions: any[] = [];
    const languageExtension = options.mode ? getLanguageExtension(options.mode) : null;

    if (languageExtension) {
      editorExtensions.push(languageExtension);
    }

    if (options.lineWrapping) {
      editorExtensions.push(EditorView.lineWrapping);
    }

    if (options.lint) {
      editorExtensions.push(lintGutter());
      if (options.mode === 'text/x-yaml') {
        editorExtensions.push(yamlLinter);
      }
      if (options.mode === 'application/json') {
        editorExtensions.push(jsonLinter);
      }
    }

    return editorExtensions;
  }, [options.lineWrapping, options.lint, options.mode, options.readOnly]);

  return (
    <ReactCodeMirror
      value={getSafeValue(value)}
      theme={material}
      basicSetup={{ lineNumbers: options.lineNumbers ?? true }}
      readOnly={options.readOnly}
      extensions={extensions}
      onChange={(nextValue, viewUpdate) => {
        if (onBeforeChange) {
          onBeforeChange(viewUpdate?.view, null, nextValue);
        }
        if (onChange) {
          onChange(nextValue, viewUpdate);
        }
      }}
      onBlur={(event) => {
        if (onBlur) {
          onBlur(editorViewRef.current, event);
        }
      }}
      onCreateEditor={(view) => {
        editorViewRef.current = view;
        if (editorDidMount) {
          editorDidMount(view, getSafeValue(value));
        }
      }}
      {...rest}
    />
  );
};

export { CodeMirror as Controlled, CodeMirror as UnControlled };
export default CodeMirror;
