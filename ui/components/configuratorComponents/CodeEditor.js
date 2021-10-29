import { Card, makeStyles, CardContent } from "@material-ui/core";
import { UnControlled as CodeMirror } from "react-codemirror2";

const useStyles = makeStyles(() => ({
  root : {
    position : "sticky"
  },
  codeMirror : {
    '& .CodeMirror' : {
      minHeight : "300px",
      height : '60vh',
    }
  },
}))

export default function CodeEditor({ yaml, saveCodeEditorChanges }) {
  const classes = useStyles();

  return (
    <div>
      <Card
        elevation={0}
        // @ts-ignore
        className={classes.root}>
        <CardContent >
          <CodeMirror
            value={yaml}
            className={classes.codeMirror}
            options={{
              theme : "material",
              lineNumbers : true,
              lineWrapping : true,
              gutters : ["CodeMirror-lint-markers"],
              mode : "text/x-yaml",
            }}
            onBlur={(a) => saveCodeEditorChanges(a)}
          />
        </CardContent>
      </Card>
    </div>
  )
}