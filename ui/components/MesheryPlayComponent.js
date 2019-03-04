import {connect} from "react-redux";
import NoSsr from '@material-ui/core/NoSsr';
import dataFetch from '../lib/data-fetch';
import {Controlled as CodeMirror} from 'react-codemirror2'
import { withStyles, Grid, FormControlLabel, Switch, FormGroup, TextField, Button, Snackbar, RadioGroup, Radio, FormLabel, FormControl } from '@material-ui/core';
import { blue } from '@material-ui/core/colors';
import MesherySnackbarWrapper from './MesherySnackbarWrapper';
import PropTypes from 'prop-types';
import { withRouter } from 'next/router';
import { updateK8SConfig } from "../lib/store";
import { bindActionCreators } from "redux";

const styles = theme => ({
  root: {
    padding: theme.spacing(10),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  margin: {
    margin: theme.spacing(1),
  },
  alreadyConfigured: {
    textAlign: 'center',
    padding: theme.spacing(20),
  },
  colorSwitchBase: {
    color: blue[300],
    '&$colorChecked': {
      color: blue[500],
      '& + $colorBar': {
        backgroundColor: blue[500],
      },
    },
  },
  colorBar: {},
  colorChecked: {},
  uploadButton: {
    margin: theme.spacing(1),
    marginTop: theme.spacing(3),
  },
  rightIcon: {
    // marginLeft: theme.spacing(1),
  },
  fileLabel: {
    width: '100%',
  },
  fileLabelText: {
    // width: '79%',
  },
  editorContainer: {
    width: '100%',
  },
  deleteLabel: {
    paddingRight: theme.spacing(2),
  },
  alignRight: {
    textAlign: 'right',
  }
});

class MesheryPlayComponent extends React.Component {
  constructor(props){
    super(props);
    
    this.cmEditor = null;

    const {Name, Ops} = props;
    if (Name === '' || Object.keys(Ops).length === 0) {
      props.router.push('/configure'); // TODO: need to figure out a better way to do this
    }

    this.state = {
      showSnackbar: false,
      snackbarVariant: '',
      snackbarMessage: '',

      selectedOp: '',
      cmEditorVal: '',
      cmEditorValError: false,
      selectionError: false,

      namespace: 'default',
      namespaceError: false,
      deleteOp: false,
    }
  }

  handleSnackbarClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }

    this.setState({ showSnackbar: false });
  };

  handleChange = name => event => {
    if (name === 'deleteOp'){
      this.setState({ [name]: event.target.checked });
      return;
    }
    if (name === 'namespace' && event.target.value !== '') {
      this.setState({ namespaceError: false });  
    }
    
    if (name === 'selectedOp' && event.target.value !== '') {
      if (event.target.value === 'custom'){ 
          if(this.state.cmEditorVal !== '' && this.cmEditor.state.lint.marked.length === 0) {
            this.setState({ selectionError: false, cmEditorValError: false });  
          }
        } else {
          this.setState({ selectionError: false });  
        }
    } 

    this.setState({ [name]: event.target.value });
  };

  handleSubmit = () => {
    const { selectedOp, selectionError, namespace, namespaceError, cmEditorVal, cmEditorValError } = this.state;
    const {Ops} = this.props;
    if (selectedOp === '' || typeof Ops[selectedOp] === 'undefined') {
        this.setState({selectionError: true});
        return;
    }
    if (selectedOp === 'custom' && (cmEditorVal === '' || this.cmEditor.state.lint.marked.length > 0)) {
      this.setState({cmEditorValError: true, selectionError: true});
      return
    }
    if (namespace === '') {
      this.setState({namespaceError: true});
      return
    }
    this.submitOp()
  }

  submitOp = () => {
    const { selectedOp, cmEditorVal, deleteOp } = this.state;
    // const fileInput = document.querySelector('#k8sfile') ;

    const data = {
      'query': selectedOp,
      'customBody': cmEditorVal,
      'deleteOp': deleteOp? 'on':'',
    }
    
    const params = Object.keys(data).map((key) => {
      return encodeURIComponent(key) + '=' + encodeURIComponent(data[key]);
    }).join('&');
    console.log(`data to be submitted for load test: ${params}`);

    // console.log(`data to be submitted for load test: ${params}`);
    let self = this;
    dataFetch('/api/mesh', { 
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
      },
      body: params
    }, result => {
      if (typeof result !== 'undefined'){
        this.setState({showSnackbar: true, snackbarVariant: 'success', snackbarMessage: 'Operation success!'});
        // this.props.updateK8SConfig({k8sConfig: {inClusterConfig, k8sfile, meshLocationURL, contextName, reconfigureCluster: false}});
        // this.props.updateMeshInfo({mesh: result});
      }
    }, self.handleError);
  }

  handleError = error => {
    this.setState({showSnackbar: true, snackbarVariant: 'error', snackbarMessage: `Operation failed: ${error}`});
  }

  handleReconfigure = () => {
    const { k8sConfig } = this.props;
    k8sConfig.reconfigureCluster = true;
    this.props.updateK8SConfig({k8sConfig});
    this.props.router.push('/configure');
  }

  render() {
    const {classes, color, iconButtonClassName, avatarClassName, ...other} = this.props;
    const {
      showSnackbar, 
      snackbarVariant, 
      snackbarMessage, 

      selectedOp,
      cmEditorVal,
      deleteOp,
      namespace,

      selectionError,
      namespaceError,
      cmEditorValError,
     } = this.state;
     const {Name, Ops} = this.props;
    var self = this;
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.root}>
          <Grid container spacing={5}>
          <Grid item xs={12} className={classes.alignRight}>
            <Button variant="contained" color="secondary" size="large" onClick={this.handleReconfigure}>
            Reconfigure
            </Button>
          </Grid>
          <Grid item xs={12}>
            <TextField
              required
              id="namespace"
              name="namespace"
              label="Namespace"
              fullWidth
              value={namespace}
              error={namespaceError}
              margin="normal"
              variant="outlined"
              onChange={this.handleChange('namespace')}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl required error={selectionError} component="fieldset">
            <FormLabel component="legend">{`Play with ${Name}`}</FormLabel>
            <RadioGroup
            aria-label={`Play with ${Name}`}
            name="query"
            className={classes.group}
            value={selectedOp}
            onChange={this.handleChange('selectedOp')}
            >
            {Object.keys(Ops).map(key => (
              
              <FormControlLabel key={key} value={key} control={<Radio />} label={Ops[key]} />
              
            ))}
           </RadioGroup>
           </FormControl>
           </Grid>
           <Grid item xs={12} sm={6} className={classes.alignRight}>
            <FormControlLabel
                  key="addUpdate"
                  control={
                    <FormControlLabel
                    key="delete"
                    control={
                        <Switch
                            checked={deleteOp}
                            onChange={this.handleChange('deleteOp')}
                        />
                    }
                    label="Add/Update"
                    labelPlacement="start"
                    className={classes.deleteLabel}
                />
                  }
                  label="Delete"
                  labelPlacement="end"
              />
            </Grid>
           
           
            <Grid item xs={12} hidden={selectedOp != 'custom'}>
            <FormControl required error={cmEditorValError} component="fieldset" className={classes.editorContainer}>
            <FormLabel component="legend">Custom yaml</FormLabel>
            <CodeMirror
                editorDidMount={editor => { this.cmEditor = editor }}
                value={cmEditorVal}
                options={{
                  mode: 'yaml',
                  theme: 'material',
                  lineNumbers: true,
                  lineWrapping: true,
                  gutters: ["CodeMirror-lint-markers"],
                  // lint: {
                  //   "getAnnotations": self.customBodyLinter,
                  //   "async": true 
                  // },
                  lint: true,
                  mode: "text/x-yaml"
                }}
                onBeforeChange={(editor, data, value) => {
                  this.setState({cmEditorVal: value});
                  if(value !== '' && this.cmEditor.state.lint.marked.length === 0) {
                    this.setState({ selectionError: false, cmEditorValError: false });  
                  }
                }}
                onChange={(editor, data, value) => {
                }}
              />
              </FormControl>
            </Grid>
          </Grid>
          <React.Fragment>
            <div className={classes.buttons}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                size="large"
                onClick={this.handleSubmit}
                className={classes.button}
              >
              Submit
              </Button>
            </div>
          </React.Fragment>
          </div>
        </React.Fragment>

          <Snackbar
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'left',
            }}
            open={showSnackbar}
            autoHideDuration={6000}
            onClose={this.handleSnackbarClose}
          >
          <MesherySnackbarWrapper 
            variant={snackbarVariant}
            message={snackbarMessage}
            onClose={this.handleSnackbarClose}
            />
        </Snackbar>
      </NoSsr>
    )
  }
}

MesheryPlayComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
      updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
  }
}

const mapStateToProps = state => {
    // console.log("header - mapping state to props. . . new title: "+ state.get("page").get("title"));
    // console.log("state: " + JSON.stringify(state));
    const mesh = state.get("mesh");
    // const k8sConfig = state.get("k8sConfig").toObject();
    let newprops = {};
    if (typeof mesh !== 'undefined'){
      newprops = { 
        Name: mesh.get("Name"),
        Ops: mesh.get("Ops").toObject(),
        k8sConfig: state.get("k8sConfig").toObject(),
      }
    }
    return newprops;
}

export default withStyles(styles)(connect(
    mapStateToProps,
    mapDispatchToProps
  )(withRouter(MesheryPlayComponent)));
