import {
  DeleteIcon,
  Dialog,
  DialogActions,
  DialogContent,
  Divider,
  FullScreenIcon,
  FullScreenExitIcon,
  IconButton,
  SaveIcon,
  Tooltip,
} from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import { UnControlled as CodeMirror } from './CodeMirror';
import { YamlDialogTitleText, StyledDialog } from './YamlDialog.styles';
import { StyledCodeMirrorWrapper } from '../designs/patterns/Cards.styles';

const YAMLDialog = ({
  fullScreen,
  name,
  toggleFullScreen,
  config_file,
  setYaml,
  deleteHandler,
  updateHandler,
  isReadOnly = false,
  type,
  updatePermissionKey,
  deletePermissionKey,
}) => {
  const defaultUpdateKey =
    type === 'pattern' ? Keys.CatalogManagementEditDesign : Keys.CatalogManagementEditWasmFilter;
  const defaultDeleteKey =
    type === 'pattern' ? Keys.CatalogManagementDeleteADesign : Keys.CatalogManagementDeleteWasmFilter;

  const resolvedUpdateKey = updatePermissionKey || defaultUpdateKey;
  const resolvedDeleteKey = deletePermissionKey || defaultDeleteKey;
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
            {fullScreen ? (
              <FullScreenExitIcon fill="currentColor" />
            ) : (
              <FullScreenIcon fill="currentColor" />
            )}
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
            <IconButton
              aria-label="Update"
              color="primary"
              onClick={updateHandler}
              size="large"
              permissionKey={resolvedUpdateKey}
            >
              <SaveIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete Filter">
            <IconButton
              aria-label="Delete"
              color="primary"
              onClick={deleteHandler}
              size="large"
              permissionKey={resolvedDeleteKey}
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
