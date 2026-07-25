import React, { useState } from 'react';
import {
  CustomTooltip,
  Dialog,
  DialogActions,
  DialogContent,
  Divider,
  IconButton,
  styled,
  FullScreenIcon,
  FullScreenExitIcon,
} from '@sistent/sistent';
import { Close as CloseIcon, Delete as DeleteIcon, Save as SaveIcon } from '@/assets/icons';
import { UnControlled as CodeMirror } from '../../general/CodeMirror';
import { FILE_OPS } from '../../../utils/Enum';

import { Keys } from '@meshery/schemas/permissions';
import { YamlDialogTitle, YamlDialogTitleText } from './MesheryPatterns.styled';

function YAMLEditor({ pattern, onClose, onSubmit, isReadOnly = false }) {
  const [yaml, setYaml] = useState(pattern.patternFile);
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  const FullScreenCodeMirrorWrapper = styled('div')(() => ({
    height: '100%',
    '& .cm-editor': {
      minHeight: '300px',
      height: fullScreen ? '80vh' : '30vh',
    },
  }));

  return (
    <Dialog
      onClose={onClose}
      aria-labelledby="pattern-dialog-title"
      open
      maxWidth="md"
      fullScreen={fullScreen}
      fullWidth={!fullScreen}
    >
      <YamlDialogTitle
        style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'end' }}
        disableTypography
        id="pattern-dialog-title"
      >
        <div>
          <YamlDialogTitleText variant="h6">{pattern.name}</YamlDialogTitleText>
        </div>
        <div>
          <CustomTooltip
            placement="top"
            title={fullScreen ? 'Exit Fullscreen' : 'Enter Fullscreen'}
            onClick={toggleFullScreen}
          >
            {fullScreen ? <FullScreenExitIcon /> : <FullScreenIcon />}
          </CustomTooltip>
          <CustomTooltip placement="top" title="Exit" onClick={onClose}>
            <CloseIcon />
          </CustomTooltip>
        </div>
      </YamlDialogTitle>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <FullScreenCodeMirrorWrapper>
          <CodeMirror
            value={pattern.patternFile}
            options={{
              theme: 'material',
              lineNumbers: true,
              lineWrapping: true,
              gutters: ['CodeMirror-lint-markers'],
              // @ts-ignore
              lint: true,
              mode: 'text/x-yaml',
              readOnly: isReadOnly,
            }}
            onChange={(_, data, val) => setYaml(val)}
          />
        </FullScreenCodeMirrorWrapper>
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        {isReadOnly ? null : (
          <>
            <CustomTooltip title="Update Design">
              <IconButton
                aria-label="Update"
                permissionKey={Keys.CatalogManagementEditDesign}
                onClick={() =>
                  onSubmit({
                    data: yaml,
                    id: pattern.id,
                    name: pattern.name,
                    type: FILE_OPS.UPDATE,
                    catalogData: pattern.catalogData,
                  })
                }
              >
                <SaveIcon />
              </IconButton>
            </CustomTooltip>
            <CustomTooltip title="Delete Pattern">
              <IconButton
                aria-label="Delete"
                permissionKey={Keys.CatalogManagementDeleteADesign}
                onClick={() =>
                  onSubmit({
                    data: yaml,
                    id: pattern.id,
                    name: pattern.name,
                    type: FILE_OPS.DELETE,
                    catalogData: pattern.catalogData,
                  })
                }
              >
                <DeleteIcon />
              </IconButton>
            </CustomTooltip>
          </>
        )}
      </DialogActions>
    </Dialog>
  );
}

export default YAMLEditor;
