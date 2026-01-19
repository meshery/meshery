import { Card, CardContent, styled } from '@sistent/sistent';
import CodeMirror from '@uiw/react-codemirror';
import { codeMirrorTheme, yamlExtensions } from '@/utils/codemirror';

export const CodeMirrorWrapper = styled(CodeMirror)(() => ({
  '& .cm-editor': {
    minHeight: '300px',
    height: '54vh',
  },
  '& .cm-scroller': {
    minHeight: '300px',
  },
}));

export default function CodeEditor({ yaml, saveCodeEditorChanges, onChange }) {
  return (
    <Card elevation={0}>
      <CardContent>
        <CodeMirrorWrapper
          value={yaml}
          theme={codeMirrorTheme}
          basicSetup={{ lineNumbers: true, highlightActiveLine: false }}
          extensions={yamlExtensions}
          onChange={(value, viewUpdate) => {
            onChange?.(value, viewUpdate);
          }}
          onBlur={() => saveCodeEditorChanges?.(yaml)}
        />
      </CardContent>
    </Card>
  );
}
