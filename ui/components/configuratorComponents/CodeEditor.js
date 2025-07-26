import { Card, CardContent, styled } from '@sistent/sistent';
import { UnControlled as CodeMirror } from 'react-codemirror2';

export const CodeMirrorWrapper = styled(CodeMirror)(() => ({
  '& .CodeMirror': {
    minHeight: '300px',
    height: '54vh',
  },
}));

export default function CodeEditor({ yaml, saveCodeEditorChanges, onChange }) {
  return (
    <Card elevation={0}>
      <CardContent>
        <CodeMirrorWrapper
          value={yaml}
          options={{
            theme: 'material',
            lineNumbers: true,
            lineWrapping: true,
            gutters: ['CodeMirror-lint-markers'],
            mode: 'text/x-yaml',
          }}
          onChange={(a, b, c) => {
            onChange(a, b, c);
          }}
          onBlur={(a) => saveCodeEditorChanges(a)}
        />
      </CardContent>
    </Card>
  );
}
