import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Typography,
  Divider,
} from '@sistent/sistent';
import { Close as CloseIcon, Fullscreen as FullscreenIcon, FullscreenExit as FullscreenExitIcon, Save as SaveIcon } from '@/components/icons';
import { UnControlled as CodeMirror } from '../CodeMirror';

function YAMLEditor({ pattern, onClose, onSubmit, isReadOnly = false }) {
  const [fullScreen, setFullScreen] = useState(false);
  const [yamlContent, setYamlContent] = useState(pattern.pattern_file);

  const handleSave = () => {
    onSubmit({
      data: yamlContent,
      name: pattern.name,
      id: pattern.id,
      type: 'update',
      metadata: pattern.metadata,
      catalog_data: pattern.catalog_data,
    });
    onClose();
  };

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  return (
    <Dialog
      open={true}
      onClose={onClose}
      fullScreen={fullScreen}
      maxWidth="md"
      fullWidth={!fullScreen}
      aria-labelledby="yamleditor-dialog-title"
    >
      <DialogTitle id="yamleditor-dialog-title" style={{ padding: '12px 24px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6">YAML Editor - {pattern.name}</Typography>
          <div>
            <IconButton onClick={toggleFullScreen} size="small" style={{ marginRight: '8px' }}>
              {fullScreen ? <FullscreenExitIcon /> : <FullscreenIcon />}
            </IconButton>
            <IconButton onClick={onClose} size="small">
              <CloseIcon />
            </IconButton>
          </div>
        </div>
      </DialogTitle>
      <Divider />
      <DialogContent
        style={{
          padding: 0,
          overflow: 'hidden',
          height: fullScreen ? 'calc(100vh - 120px)' : '500px',
        }}
      >
        <CodeMirror
          value={yamlContent}
          options={{
            lineNumbers: true,
            mode: 'yaml',
            theme: 'material',
            lineWrapping: true,
            readOnly: isReadOnly,
          }}
          onChange={(_, __, value) => setYamlContent(value)}
          style={{ height: '100%' }}
        />
      </DialogContent>
      <Divider />
      {!isReadOnly && (
        <DialogActions style={{ padding: '8px 24px' }}>
          <IconButton
            onClick={handleSave}
            color="primary"
            title="Save changes"
            style={{ borderRadius: '4px', padding: '8px 16px' }}
          >
            <SaveIcon style={{ marginRight: '8px' }} />
            <Typography variant="button">Save</Typography>
          </IconButton>
        </DialogActions>
      )}
    </Dialog>
  );
}

YAMLEditor.displayName = 'YAMLEditor';

export default YAMLEditor;
