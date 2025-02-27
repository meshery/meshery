import { Card, CardContent, styled } from '@layer5/sistent';
import { useEffect, useState } from 'react';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import { UsesSistent } from '../SistentWrapper';

export const CardRoot = styled(Card)({
  position: 'sticky',
});

export const CodeMirrorWrapper = styled(CodeMirror)(({ scrollPos }) => ({
  '& .CodeMirror': {
    minHeight: '300px',
    height: `${getDynamicVh(scrollPos)}`,
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
    <UsesSistent>
      <Wrapper scrollPos={scrollPos} fullWidth={fullWidth}>
        <CardRoot elevation={0}>
          <CardContent>
            <CodeMirror
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
              style={{
                '& .CodeMirror': {
                  minHeight: '300px',
                  height: '54vh',
                },
              }}
            />
          </CardContent>
        </CardRoot>
      </Wrapper>
    </UsesSistent>
  );
}

/**
 * Provides dynamic height according to scroll calculations
 *
 * @param {DoubleRange} scrollPos
 * @returns dynamically calcultaed height in vh
 */
function getDynamicVh(scrollPos) {
  if (window.scrollY === 0) {
    return '67vh';
  }
  const per = getScrollPercentage();
  const threshold = 0.06;
  const vh = 67 + 15 * (per / threshold); // calc(67vh)
  if (per < threshold) {
    return scrollPos > 106 ? `${vh}vh` : '67vh';
  } else if (per > 0.95) {
    return 'calc(100vh - 232px)';
  } else {
    return '82vh';
  }
}

function getScrollPercentage() {
  return window.scrollY / (document.body.scrollHeight - window.innerHeight);
}
