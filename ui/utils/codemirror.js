import { EditorView } from '@codemirror/view';
import { yaml } from '@codemirror/lang-yaml';
import { json } from '@codemirror/lang-json';
import { material } from '@uiw/codemirror-theme-material';
import jsyaml from 'js-yaml';

export const codeMirrorTheme = material;
export const yamlExtensions = [yaml(), EditorView.lineWrapping];
export const jsonExtensions = [json(), EditorView.lineWrapping];

export const isValidYaml = (value) => {
  if (!value) return false;
  try {
    jsyaml.load(value);
    return true;
  } catch (e) {
    return false;
  }
};

export const isValidJson = (value) => {
  if (!value) return false;
  try {
    JSON.parse(value);
    return true;
  } catch (e) {
    return false;
  }
};
