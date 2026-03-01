import { Card, CardContent, styled } from '@sistent/sistent';
import { UnControlled as CodeMirror } from '../CodeMirror';

export const CodeMirrorWrapper = styled('div')(() => ({
  '& .cm-editor': {
    minHeight: '300px',
    height: '54vh',
  },
}));

export default function CodeEditor({ yaml, saveCodeEditorChanges, onChange }) {
  return (
    <Card elevation={0}>
      <CardContent>
        <CodeMirrorWrapper>
          <CodeMirror
            value={yaml}
            options={{
              theme: 'material',
              lineNumbers: true,
              lineWrapping: true,
              gutters: ['CodeMirror-lint-markers'],
              mode: 'text/x-yaml',
            }}
            onChange={(value, viewUpdate) => {
              onChange?.(value, viewUpdate);
            }}
            onBlur={(value) => saveCodeEditorChanges?.(value)}
          />
        </CodeMirrorWrapper>
      </CardContent>
    </Card>
  );
}
