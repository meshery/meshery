import { Card, CardContent, styled } from '@sistent/sistent';
import { useEffect, useState } from 'react';
import { UnControlled as CodeMirror } from 'react-codemirror2';

export const CardRoot = styled(Card)({
  position: 'sticky',
});

export const CodeMirrorWrapper = styled(CodeMirror)(() => ({
  '& .CodeMirror': {
    height: '54vh',
  },
}));

export const Wrapper = styled('div')(({ theme, scrollPos, fullWidth }) => ({
  [theme.breakpoints.up('md')]: {
    top: scrollPos >= 106 ? 106 : window.scrollY > 0 ? 208 - scrollPos : 'auto',
    position: fullWidth ? 'inherit' : 'fixed',
    minWidth: fullWidth ? undefined : 'calc(50% - 175px)',
    maxWidth: fullWidth ? undefined : 'calc(50% - 175px)',
  },
}));

export const Icon = styled('div')({
  position: 'absolute',
  right: '24px',
  bottom: '30px',
  color: '#fff',
  zIndex: 11,
  backgroundColor: 'rgba(255, 255, 255, 0.05)',
});

export default function CodeEditor({ yaml, saveCodeEditorChanges, fullWidth, onChange }) {
  const [scrollPos, setScrollPos] = useState(67);

  useEffect(() => {
    const handleScroll = () => setScrollPos(window.scrollY);

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <Wrapper scrollPos={scrollPos} fullWidth={fullWidth}>
      <CardRoot elevation={0}>
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
      </CardRoot>
    </Wrapper>
  );
}
