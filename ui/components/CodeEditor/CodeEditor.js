import React from "react";
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-yaml";
import "ace-builds/src-noconflict/theme-twilight";
import yaml from "js-yaml";

// eslint-disable-next-line react/prop-types
const CodeEditor = ({ defaultValue = "", name = "YAML Editor" }) => {
  const [value, setValue] = React.useState(defaultValue);
  const [annotations, setAnnotations] = React.useState([]);
  const checkerFunction = (value) => {
    try {
      console.log(yaml.load(value));
    } catch (error) {
      console.log(error.mark.line, error.mark.column, error.message, error.name);
      setAnnotations((annotations) => [
        {
          row: error.mark.line,
          column: error.mark.column,
          text: error.message,
          type: error.name,
        },
        ...annotations,
      ]);
    }
  };
  return (
    <div>
      <AceEditor
        mode="yaml"
        theme="twilight"
        name={name}
        onChange={(value) => {
          setValue(value);
          checkerFunction(value);
        }}
        editorProps={{
          $blockScrolling: true,
        }}
        fontSize={12}
        value={value}
        setOptions={{
          useWorker: false,
        }}
        enableBasicAutocompletion={false}
        annotations={annotations}
        // markers={annotations ? [{ startRow: annotations?.row, startCol: annotations?.column, endRow: annotations?.row, endCol: annotations?.column, className: 'error-marker', type: 'background' }] : []}
        tabSize={8}
      />
    </div>
  );
};

export default CodeEditor;
