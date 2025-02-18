import {
  Dialog,
  DialogActions,
  DialogContent,
  Divider,
  IconButton,
  Tooltip,
  styled,
} from '@layer5/sistent';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import DeleteIcon from '@mui/icons-material/Delete';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import SaveIcon from '@mui/icons-material/Save';
import { YamlDialogTitleText, StyledDialog } from './MesheryPatterns/style';

const StyledCodeMirrorWrapper = styled('div')(({ fullScreen }) => ({
  height: fullScreen ? '100%' : 'auto',
  width: '100%',
  '& .CodeMirror': {
    minHeight: '300px',
    height: fullScreen ? '100%' : 'auto',
    width: '100%',
  },
}));

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
  return (
    <Dialog
      aria-labelledby="filter-dialog-title"
      open
      maxWidth="md"
      fullScreen={fullScreen}
      fullWidth={!fullScreen}
    >
      <StyledDialog disableTypography id="filter-dialog-title">
        <YamlDialogTitleText variant="h6">{name}</YamlDialogTitleText>
        <Tooltip title="Exit Fullscreen" arrow placement="bottom">
          <IconButton onClick={toggleFullScreen} size="large">
            {fullScreen ? <FullscreenExitIcon /> : <FullscreenIcon />}
          </IconButton>
        </Tooltip>
      </StyledDialog>
      <Divider />
      <DialogContent>
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
      </DialogContent>
      <Divider />
      {!isReadOnly && (
        <DialogActions>
          <Tooltip title="Update Pattern">
            <IconButton aria-label="Update" color="primary" onClick={updateHandler} size="large">
              <SaveIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete Filter">
            <IconButton aria-label="Delete" color="primary" onClick={deleteHandler} size="large">
              <DeleteIcon />
            </IconButton>
          </Tooltip>
        </DialogActions>
      )}
    </Dialog>
  );
};

export default YAMLDialog;
