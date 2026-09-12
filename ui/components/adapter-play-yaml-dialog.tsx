import React from 'react';
import {
  DeleteIcon,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  Grid,
  IconButton,
  PlayArrowIcon as PlayIcon,
} from '@sistent/sistent';
import { Controlled as CodeMirror } from './general/CodeMirror';
import ReactSelectWrapper from './general/ReactSelectWrapper';
import { iconMedium } from '../css/icons.styles';

interface AdapterYamlDialogProps {
  open: boolean;
  isDelete: boolean;
  adapterName: string;
  namespace: { value: string; label: string };
  namespaceError: boolean;
  namespaceList: Array<{ value: string; label: string }>;
  onNamespaceChange: (newValue: { value: string; label: string }) => void;
  version: { value: string; label: string };
  versionError: boolean;
  versionList: Array<{ value: string; label: string }>;
  onVersionChange: (newValue: { value: string; label: string }) => void;
  value: string;
  onBeforeChange: (editor: unknown, data: unknown, value: string) => void;
  onClose: () => void;
  onApply: () => void;
}

/**
 * Renders the per-category YAML editor dialog for the "custom"
 * adapter operation. Extracted from MesheryAdapterPlayComponent's
 * generateYAMLEditor() with identical markup and identical handler
 * call signatures.
 */
const AdapterYamlDialog: React.FC<AdapterYamlDialogProps> = ({
  open,
  isDelete,
  adapterName,
  namespace,
  namespaceError,
  namespaceList,
  onNamespaceChange,
  version,
  versionError,
  versionList,
  onVersionChange,
  value,
  onBeforeChange,
  onClose,
  onApply,
}) => {
  return (
    <Dialog
      onClose={onClose}
      aria-labelledby="adapter-dialog-title"
      open={open}
      fullWidth
      maxWidth="md"
    >
      <DialogTitle id="adapter-dialog-title" onClose={onClose}>
        {adapterName} Adapter - Custom YAML
        {isDelete ? '(delete)' : ''}
      </DialogTitle>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <Grid container spacing={5}>
          <Grid item xs={6}>
            <ReactSelectWrapper
              label="Namespace"
              value={namespace}
              error={namespaceError}
              options={namespaceList}
              onChange={onNamespaceChange}
            />
          </Grid>
          <Grid item xs={6}>
            <ReactSelectWrapper
              label="Version"
              value={version}
              error={versionError}
              options={versionList}
              onChange={onVersionChange}
            />
          </Grid>
          <Grid item xs={12}>
            <CodeMirror
              value={value}
              options={{
                theme: 'material',
                lineNumbers: true,
                lineWrapping: true,
                gutters: ['CodeMirror-lint-markers'],
                lint: true,
                mode: 'text/x-yaml',
              }}
              onBeforeChange={onBeforeChange}
            />
          </Grid>
        </Grid>
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        <IconButton aria-label="Apply" onClick={onApply}>
          {!isDelete && <PlayIcon style={iconMedium} />}
          {isDelete && <DeleteIcon style={iconMedium} />}
        </IconButton>
      </DialogActions>
    </Dialog>
  );
};

export default AdapterYamlDialog;
