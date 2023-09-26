import { Card, CardContent, makeStyles } from '@material-ui/core';
import { useEffect, useState } from 'react';
import { UnControlled as CodeMirror } from 'react-codemirror2';

const useStyles = makeStyles((theme) => ({
  cardRoot: {
    position: 'sticky',
  },
  codeMirror: {
    '& .CodeMirror': {
      minHeight: '300px',
      height: ({ scrollPos }) => getDynamicVh(scrollPos),
    },
  },
  wrapper: {
    [theme.breakpoints.up('md')]: {
      top: ({ scrollPos }) =>
        scrollPos >= 106 ? 106 : window.scrollY > 0 ? 208 - scrollPos : 'auto',
      position: ({ fullWidth }) => (fullWidth ? 'inherit' : 'fixed'),
      minWidth: ({ fullWidth }) => (fullWidth ? undefined : 'calc(50% - 175px)'),
      maxWidth: ({ fullWidth }) => (fullWidth ? undefined : 'calc(50% - 175px)'),
    },
  },
  icon: {
    position: 'absolute',
    right: '24px',
    bottom: '30px',
    color: '#fff',
    zIndex: 11,
    backgroundColor: 'rgb(255,255,255,0.05)',
  },
}));

export default function CodeEditor({ yaml, saveCodeEditorChanges, fullWidth, onChange }) {
  const [style, setStyle] = useState(67);
  const classes = useStyles({ scrollPos: style, fullWidth });

  useEffect(() => {
    window.addEventListener('scroll', handleScroll);

    return () => {
      window.removeEventListener('scroll', () => {});
    };
  }, []);

  const handleScroll = () => {
    setStyle(window.scrollY);
  };

  return (
    <div className={classes.wrapper} style={{}}>
      <Card
        elevation={0}
        // @ts-ignore
        className={classes.cardRoot}
      >
        <CardContent>
          <CodeMirror
            value={yaml}
            className={classes.codeMirror}
            options={{
              theme: 'material',
              lineNumbers: true,
              lineWrapping: true,
              gutters: ['CodeMirror-lint-markers'],
              mode: 'text/x-yaml',
            }}
            onChange={(a, b, c) => {
              console.log('onChange', a, b, c);
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
      </Card>
    </div>
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
