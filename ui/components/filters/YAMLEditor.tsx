import React, { useState } from 'react';
import { Close as CloseIcon, Delete as DeleteIcon, Save as SaveIcon } from '@/assets/icons';
import {
  CustomTooltip,
  Dialog,
  DialogActions,
  DialogTitle,
  Divider,
  IconButton,
  styled,
  FullScreenIcon,
  FullScreenExitIcon,
} from '@sistent/sistent';
import { UnControlled as CodeMirror } from '../CodeMirror';
import { FILE_OPS } from '../../utils/Enum';
import { iconMedium } from '../../css/icons.styles';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import TooltipIcon from './TooltipIcon';
import { YmlDialogTitle, YmlDialogTitleText } from './Filters.styled';
import type { YAMLEditorProps } from './Filters.types';

function YAMLEditor({ filter, onClose, onSubmit }: YAMLEditorProps) {
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  const FullScreenCodeMirrorWrapper = styled('div')(() => ({
    height: '100%',
    '& .cm-editor': {
      minHeight: '300px',
      height: fullScreen ? '80vh' : '100%',
    },
  }));

  let resourceData;
  try {
    resourceData = JSON.parse(filter.filter_resource);
  } catch (error) {
    // Handling the error or provide a default value
    console.error('Error parsing JSON:', error);
    resourceData = {}; // Setting a default value if parsing fails
  }

  const config = resourceData?.settings?.config || '';
  const [yaml, setYaml] = useState(config);

  return (
    <Dialog
      onClose={onClose}
      aria-labelledby="filter-dialog-title"
      open
      maxWidth="md"
      fullScreen={fullScreen}
      fullWidth={!fullScreen}
    >
      <YmlDialogTitle>
        <DialogTitle
          disableTypography
          id="filter-dialog-title"
          style={{ width: '100%', display: 'flex' }}
        >
          <YmlDialogTitleText variant="h6">{filter.name}</YmlDialogTitleText>
          <TooltipIcon
            title={fullScreen ? 'Exit Fullscreen' : 'Enter Fullscreen'}
            onClick={toggleFullScreen}
          >
            {fullScreen ? (
              <FullScreenExitIcon style={iconMedium} />
            ) : (
              <FullScreenIcon style={iconMedium} />
            )}
          </TooltipIcon>
          <TooltipIcon title="Exit" onClick={onClose}>
            <CloseIcon style={iconMedium} />
          </TooltipIcon>
        </DialogTitle>
      </YmlDialogTitle>
      <Divider variant="fullWidth" light />
      <FullScreenCodeMirrorWrapper>
        <CodeMirror
          value={config}
          options={{
            theme: 'material',
            lineNumbers: true,
            lineWrapping: true,
            gutters: ['CodeMirror-lint-markers'],
            lint: true,
            mode: 'text/x-yaml',
          }}
          onChange={(_, data, val) => setYaml(val)}
        />
      </FullScreenCodeMirrorWrapper>
      <Divider variant="fullWidth" light />
      <DialogActions>
        <CustomTooltip title="Update Filter">
          <IconButton
            aria-label="Update"
            disabled={
              !CAN(
                Keys.CatalogManagementEditWasmFilter.id,
                Keys.CatalogManagementEditWasmFilter.function,
              )
            }
            onClick={() =>
              onSubmit({
                data: yaml,
                id: filter.id,
                name: filter.name,
                type: FILE_OPS.UPDATE,
                catalogData: filter.catalogData,
              })
            }
          >
            <SaveIcon style={iconMedium} />
          </IconButton>
        </CustomTooltip>
        <CustomTooltip title="Delete Filter">
          <IconButton
            aria-label="Delete"
            disabled={
              !CAN(
                Keys.CatalogManagementDeleteWasmFilter.id,
                Keys.CatalogManagementDeleteWasmFilter.function,
              )
            }
            onClick={() =>
              onSubmit({
                data: yaml,
                id: filter.id,
                name: filter.name,
                type: FILE_OPS.DELETE,
                catalogData: filter.catalogData,
              })
            }
          >
            <DeleteIcon style={iconMedium} />
          </IconButton>
        </CustomTooltip>
      </DialogActions>
    </Dialog>
  );
}

export default YAMLEditor;
