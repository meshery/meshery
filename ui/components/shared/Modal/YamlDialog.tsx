import React from 'react';
import { Divider, IconButton, Tooltip, Box } from '@sistent/sistent';
import Modal from './Modal';
import { UnControlled as CodeMirror } from '@/components/CodeMirror';
import { FullscreenExit, Delete, Fullscreen, Save } from '@/components/icons';
import { StyledCodeMirrorWrapper } from '@/components/MesheryPatterns/Cards.styles';

const YAMLDialog = ({
  fullScreen,
  name,
  toggleFullScreen,
  config_file,
  setYaml,
  deleteHandler,
  updateHandler,
  isReadOnly = false,
}) => {
  const headerRight = (
    <Tooltip title="Exit Fullscreen" arrow placement="bottom">
      <IconButton onClick={toggleFullScreen} size="large">
        {fullScreen ? <FullscreenExit /> : <Fullscreen />}
      </IconButton>
    </Tooltip>
  );

  const footer = !isReadOnly && (
    <Box display="flex" gap={1}>
      <Tooltip title="Update Pattern">
        <IconButton aria-label="Update" color="primary" onClick={updateHandler} size="large">
          <Save />
        </IconButton>
      </Tooltip>
      <Tooltip title="Delete Filter">
        <IconButton aria-label="Delete" color="primary" onClick={deleteHandler} size="large">
          <Delete />
        </IconButton>
      </Tooltip>
    </Box>
  );

  return (
    <Modal
      open={true}
      title={name}
      headerRight={headerRight}
      closeModal={toggleFullScreen}
      maxWidth="md"
      fullScreen={fullScreen}
      footer={footer}
    >
      <Divider />
      <StyledCodeMirrorWrapper fullScreen={fullScreen}>
        <CodeMirror
          value={config_file}
          options={{
            theme: 'material',
            lineNumbers: true,
            lineWrapping: true,
            gutters: ['CodeMirror-lint-markers'],
            lint: true,
            mode: 'text/x-yaml',
            readOnly: isReadOnly,
          }}
          onChange={(_, data, val) => setYaml(val)}
        />
      </StyledCodeMirrorWrapper>
    </Modal>
  );
};

export default YAMLDialog;
