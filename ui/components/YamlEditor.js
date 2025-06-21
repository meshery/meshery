import React, { useState } from 'react';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import {
  CustomTooltip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  IconButton,
  Typography,
  styled,
} from '@layer5/sistent';
import CloseIcon from '@mui/icons-material/Close';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import SaveIcon from '@mui/icons-material/Save';
import DeleteIcon from '@mui/icons-material/Delete';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { FILE_OPS } from '../utils/Enum';

const YamlDialogTitle = styled(DialogTitle)(() => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'end',
}));

const YamlDialogTitleText = styled(Typography)(() => ({
  flexGrow: 1,
}));

/**
 * Reusable YAML editor component for editing YAML content in a dialog
 * @param {Object} props - Component props
 * @param {string} props.title - Title of the dialog
 * @param {string} props.content - YAML content to edit
 * @param {boolean} [props.isReadOnly=false] - Whether the editor is in read-only mode
 * @param {Function} props.onClose - Function to call when the dialog is closed
 * @param {Function} props.onSubmit - Function to call when the content is submitted
 * @param {Object} [props.additionalButtons] - Additional buttons configuration
 * @param {string} [props.id] - ID of the resource being edited
 * @param {string} [props.resourceType='pattern'] - Type of resource (pattern or filter)
 * @param {Object} [props.metadata] - Additional metadata to pass to onSubmit
 * @returns {JSX.Element} YamlEditor component
 */
function YamlEditor({
  title,
  content,
  isReadOnly = false,
  onClose,
  onSubmit,
  additionalButtons,
  id,
  resourceType = 'pattern',
  metadata = {},
}) {
  const [yamlContent, setYamlContent] = useState(content);
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  const handleUpdate = () => {
    onSubmit({
      data: yamlContent,
      id: id,
      name: title,
      type: FILE_OPS.UPDATE,
      catalog_data: metadata.catalog_data,
      ...metadata,
    });
  };

  const handleDelete = () => {
    onSubmit({
      data: yamlContent,
      id: id,
      name: title,
      type: FILE_OPS.DELETE,
      catalog_data: metadata.catalog_data,
      ...metadata,
    });
  };

  const FullScreenCodeMirrorWrapper = styled('div')(() => ({
    height: '100%',
    '& .CodeMirror': {
      minHeight: '300px',
      height: fullScreen ? '80vh' : '30vh',
    },
  }));

  // Determine permissions keys based on resource type
  const editPermissionKey = resourceType === 'pattern' ? keys.EDIT_DESIGN : keys.EDIT_WASM_FILTER;
  const deletePermissionKey =
    resourceType === 'pattern' ? keys.DELETE_A_DESIGN : keys.DELETE_WASM_FILTER;

  return (
    <Dialog
      onClose={onClose}
      aria-labelledby={`${resourceType}-dialog-title`}
      open
      maxWidth="md"
      fullScreen={fullScreen}
      fullWidth={!fullScreen}
    >
      <YamlDialogTitle disableTypography id={`${resourceType}-dialog-title`}>
        <YamlDialogTitleText variant="h6">{title}</YamlDialogTitleText>
        <div>
          <CustomTooltip
            placement="top"
            title={fullScreen ? 'Exit Fullscreen' : 'Enter Fullscreen'}
            interactive
          >
            <IconButton onClick={toggleFullScreen}>
              {fullScreen ? <FullscreenExitIcon /> : <FullscreenIcon />}
            </IconButton>
          </CustomTooltip>
          <CustomTooltip placement="top" title="Exit" interactive>
            <IconButton onClick={onClose}>
              <CloseIcon />
            </IconButton>
          </CustomTooltip>
        </div>
      </YamlDialogTitle>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <FullScreenCodeMirrorWrapper>
          <CodeMirror
            value={content || ''}
            options={{
              theme: 'material',
              lineNumbers: true,
              lineWrapping: true,
              gutters: ['CodeMirror-lint-markers'],
              lint: true,
              mode: 'text/x-yaml',
              readOnly: isReadOnly,
            }}
            onChange={(_, data, val) => setYamlContent(val)}
          />
        </FullScreenCodeMirrorWrapper>
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        {!isReadOnly && (
          <>
            <CustomTooltip title={`Update ${resourceType === 'pattern' ? 'Design' : 'Filter'}`}>
              <div>
                <IconButton
                  aria-label="Update"
                  disabled={!CAN(editPermissionKey.action, editPermissionKey.subject)}
                  onClick={handleUpdate}
                >
                  <SaveIcon />
                </IconButton>
              </div>
            </CustomTooltip>
            <CustomTooltip title={`Delete ${resourceType === 'pattern' ? 'Design' : 'Filter'}`}>
              <div>
                <IconButton
                  aria-label="Delete"
                  disabled={!CAN(deletePermissionKey.action, deletePermissionKey.subject)}
                  onClick={handleDelete}
                >
                  <DeleteIcon />
                </IconButton>
              </div>
            </CustomTooltip>
            {additionalButtons}
          </>
        )}
      </DialogActions>
    </Dialog>
  );
}

export default YamlEditor;
