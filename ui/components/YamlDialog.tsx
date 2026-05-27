import {
  DeleteIcon,
  Dialog,
  DialogActions,
  DialogContent,
  Divider,
  FullScreenIconExit as FullscreenExitIcon,
  FullScreenIcon as FullscreenIcon,
  IconButton,
  SaveIcon,
  Tooltip,
} from '@sistent/sistent';
import { UnControlled as CodeMirror } from './CodeMirror';
import { YamlDialogTitleText, StyledDialog } from './designs/patterns/style';
import { StyledCodeMirrorWrapper } from './designs/patterns/Cards.styles';

const YAMLDialog = ({
  fullScreen,
  name,
  toggleFullScreen,
  config_file,
  setYaml,
  deleteHandler,
  updateHandler,
  type = 'filter',
  isReadOnly = false,
}) => {
  const resourceLabel = type === 'pattern' ? 'pattern' : 'filter';
  const resourceTitle = type === 'pattern' ? 'Pattern' : 'Filter';

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
          <IconButton
            aria-label={fullScreen ? 'Exit fullscreen YAML editor' : 'Enter fullscreen YAML editor'}
            onClick={toggleFullScreen}
            size="large"
          >
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
          <Tooltip title={`Update ${resourceTitle}`}>
            <IconButton
              aria-label={`Update ${resourceLabel} YAML`}
              color="primary"
              onClick={updateHandler}
              size="large"
            >
              <SaveIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title={`Delete ${resourceTitle}`}>
            <IconButton
              aria-label={`Delete ${resourceLabel} YAML`}
              color="primary"
              onClick={deleteHandler}
              size="large"
            >
              <DeleteIcon />
            </IconButton>
          </Tooltip>
        </DialogActions>
      )}
    </Dialog>
  );
};

export default YAMLDialog;
