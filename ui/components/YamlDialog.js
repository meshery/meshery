import {
  Dialog,
  DialogActions,
  DialogContent,
  Divider,
  IconButton,
  Tooltip,
} from '@sistent/sistent';
import CodeMirror from '@uiw/react-codemirror';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import DeleteIcon from '@mui/icons-material/Delete';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import SaveIcon from '@mui/icons-material/Save';
import { YamlDialogTitleText, StyledDialog } from './MesheryPatterns/style';
import { StyledCodeMirrorWrapper } from './MesheryPatterns/Cards.styles';
import { codeMirrorTheme, yamlExtensions } from '@/utils/codemirror';

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
            theme={codeMirrorTheme}
            basicSetup={{ lineNumbers: true, highlightActiveLine: false }}
            extensions={yamlExtensions}
            editable={!isReadOnly}
            onChange={(value) => setYaml(value)}
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
