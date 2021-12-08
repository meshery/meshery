import { Card, makeStyles, CardContent } from '@material-ui/core';
import { useEffect, useState } from 'react';
import { UnControlled as CodeMirror } from 'react-codemirror2';

const useStyles = makeStyles(theme => ({
  root: {
    position: 'sticky'
  },
  codeMirror: {
    '& .CodeMirror': {
      minHeight: '300px',
      height: ({ scrollPos }) => getDynamicVh(scrollPos),
    }
  },
  wrapper: {
    [theme.breakpoints.up('md')]: {
      top: ({ scrollPos }) => scrollPos >= 106 ? 106 : window.scrollY > 0 ? 208 - scrollPos: 'auto',
      position: 'fixed',
      minWidth: 'calc(50% - 175px)',
      maxWidth: 'calc(50% - 175px)',
    },
  }
}));

export default function CodeEditor({ yaml, saveCodeEditorChanges }) {
  const [style, setStyle] = useState(67);
  const classes = useStyles({ scrollPos: style });

  useEffect(() => {
    window.addEventListener('scroll', handleScroll);

    return () => {
      window.removeEventListener('scroll', () => {});
    };
  },[]);

  const handleScroll = () => {
    setStyle(window.scrollY);
  };

  return (
    <div className={classes.wrapper} >
      <Card
        elevation={0}
        // @ts-ignore
        className={classes.root}
      >

        <CardContent >
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
            onBlur={(a) => saveCodeEditorChanges(a)}
            style={{
              '& .CodeMirror': {
                minHeight: '300px',
                height: '67vh',
              }
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
  if (window.scrollY == 0) {
    return '67vh';
  }
  const per = getScrollPercentage();
  const threshold = 0.06;
  const vh = 67 + 15*(per/threshold); // calc(67vh)
  if (per < threshold) {
    return scrollPos > 106 ? `${vh}vh` : '67vh';
  } else if (per > 0.95) {
    return 'calc(100vh - 232px)';
  } else {
    return '82vh';
  }
}

function getScrollPercentage() {
  return window.scrollY/(document.body.scrollHeight - window.innerHeight);
}