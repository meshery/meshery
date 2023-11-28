import {
  Dialog,
  DialogActions,
  DialogContent,
  Divider,
  IconButton,
  Tooltip,
} from '@layer5/sistent-components';
import useStyles from './MesheryPatterns/Cards.styles';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import FullscreenExit from '@material-ui/icons/FullscreenExit';
import DeleteIcon from '@material-ui/icons/Delete';
import Fullscreen from '@material-ui/icons/Fullscreen';
import Save from '@material-ui/icons/Save';
import { StyledDialog, YamlDialogTitleText } from './MesheryPatterns/style';

const YAMLDialog = ({
  fullScreen,
  name,
  toggleFullScreen,
  config_file,
  setYaml,
  deleteHandler,
  updateHandler,
}) => {
  const classes = useStyles();
  return (
    <Dialog
      aria-labelledby="filter-dialog-title"
      open
      maxWidth="md"
      fullScreen={fullScreen}
      fullWidth={!fullScreen}
    >
      <StyledDialog disableTypography id="filter-dialog-title">
        <YamlDialogTitleText variant="h6" className={classes.yamlDialogTitleText}>
          {name}
        </YamlDialogTitleText>
        <Tooltip title="Exit Fullscreen" arrow interactive placement="bottom">
          <IconButton onClick={toggleFullScreen}>
            {fullScreen ? <FullscreenExit /> : <Fullscreen />}
          </IconButton>
        </Tooltip>
      </StyledDialog>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <CodeMirror
          value={config_file}
          className={fullScreen ? classes.fullScreenCodeMirror : ''}
          options={{
            theme: 'material',
            lineNumbers: true,
            lineWrapping: true,
            gutters: ['CodeMirror-lint-markers'],
            // @ts-ignore
            lint: true,
            mode: 'text/x-yaml',
          }}
          onChange={(_, data, val) => setYaml(val)}
        />
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        <Tooltip title="Update Pattern">
          <IconButton aria-label="Update" color="primary" onClick={updateHandler}>
            <Save />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete Filter">
          <IconButton aria-label="Delete" color="primary" onClick={deleteHandler}>
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </DialogActions>
    </Dialog>
  );
};

export default YAMLDialog;
