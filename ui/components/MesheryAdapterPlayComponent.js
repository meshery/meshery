import NoSsr from '@material-ui/core/NoSsr';
import dataFetch from '../lib/data-fetch';
import {Controlled as CodeMirror} from 'react-codemirror2'
import { withStyles, Grid, FormControlLabel, TextField, Button, FormLabel, FormControl, IconButton, FormGroup, Dialog, DialogTitle, DialogContent, DialogActions, Divider, Card, CardHeader, Avatar, CardMedia, CardContent, Typography, CardActions, Collapse, Menu, MenuItem, InputAdornment } from '@material-ui/core';
import { blue } from '@material-ui/core/colors';
import PropTypes from 'prop-types';
import { withRouter } from 'next/router';
import { updateProgress } from '../lib/store';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import CloseIcon from '@material-ui/icons/Close';
import { withSnackbar } from 'notistack';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowRight, faTrashAlt, faTerminal, faExternalLinkSquareAlt } from '@fortawesome/free-solid-svg-icons';
import FavoriteIcon from '@material-ui/icons/Favorite';
import ShareIcon from '@material-ui/icons/Share';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
// import MoreVertIcon from '@material-ui/icons/MoreVert';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import PlayIcon from '@material-ui/icons/PlayArrow';

const styles = theme => ({
  root: {
    padding: theme.spacing(10),
    width: '100%',
  },
  buttons: {
    // display: 'flex',
    // justifyContent: 'flex-end',
    width: '100%',
  },
  custom: {
    // float: 'right',
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
  },
  padLeft: {
    paddingLeft: theme.spacing(0.25),
  },
  padRight: {
    paddingRight: theme.spacing(0.25),
  },
  deleteRight: {
    float: 'right',
  },
  expTitleIcon: {
    width: theme.spacing(3),
    display: 'inline',
    verticalAlign: 'middle',
  },
  expIstioTitleIcon: {
    width: theme.spacing(2),
    display: 'inline',
    verticalAlign: 'middle',
    marginLeft: theme.spacing(0.5),
    marginRight: theme.spacing(0.5),
  },
  expTitle: {
    display: 'inline',
    verticalAlign: 'middle',
  }
});


class MesheryAdapterPlayComponent extends React.Component {
  constructor(props){
    super(props);
    
    this.cmEditorAdd = null;
    this.cmEditorDel = null;

    const {adapter} = props;
    
    const menuState = {};

    this.addIconEles = {};
    this.delIconEles = {};
    // initializing menuState;
    if(adapter && adapter.ops){
      // adapter.ops.forEach(({category}) => {
      //   if(typeof category === 'undefined'){
      //     category = 0;
      //   }
      //   // menuState[category] = {
      //   //   add: false,
      //   //   delete: false,
      //   // }
      // })
      // NOTE: this will have to updated to match the categories
      [0,1,2,3,4].forEach(i => {
        menuState[i] = {
          add: false,
          delete: false,
        }
      });
    }

    this.state = {
      selectedOp: '',
      cmEditorValAdd: '',
      cmEditorValAddError: false,

      cmEditorValDel: '',
      cmEditorValDelError: false,

      selectionError: false,

      namespace: 'default',
      namespaceError: false,

      customDialogAdd: false,
      customDialogDel: false,

      menuState, // category: {add: 1, delete: 0}
    }
  }

  handleChange = (name, isDelete=false) => {
    const self = this;
    return event => {
      if (name === 'namespace' && event.target.value !== '') {
        this.setState({ namespaceError: false });  
      }
      
      if (name === 'selectedOp' && event.target.value !== '') {
        if (event.target.value === 'custom'){ 
          if (isDelete){
            if(self.state.cmEditorValDel !== '' && self.cmEditorDel.state.lint.marked.length === 0) {
              self.setState({ selectionError: false, cmEditorValDelError: false });  
            }
          } else {
            if(self.state.cmEditorValAdd !== '' && self.cmEditorAdd.state.lint.marked.length === 0) {
              self.setState({ selectionError: false, cmEditorValAddError: false });  
            }
          }
        } else {
          self.setState({ selectionError: false });  
        }
      } 

      self.setState({ [name]: event.target.value });
    };
  }

  handleModalClose(isDelete){
    const self = this;
    return () => {
      const item = isDelete?'customDialogDel':'customDialogAdd';
      self.setState({[item]: false});
    }
  }

  handleModalOpen(isDelete) {
    const self = this;
    return () => {
      const item = isDelete?'customDialogDel':'customDialogAdd';
      self.setState({[item]: true});
    }
  }

  // handleDelete = (selectedOp) => () => {
  //   this.handleSubmit(selectedOp, true)();
  // }

  handleSubmit = (cat, selectedOp, deleteOp=false) => {
    const self = this;
    return () => {
      const { namespace, namespaceError, cmEditorValAdd, cmEditorValAddError, cmEditorValDel, cmEditorValDelError } = self.state;
      const {adapter} = self.props;
      const filteredOp = adapter.ops.filter(({key}) => key === selectedOp);
      if (selectedOp === '' || typeof filteredOp === 'undefined' || filteredOp.length === 0) {
        self.setState({selectionError: true});
        return;
      }
      if(deleteOp){
        if (selectedOp === 'custom' && (cmEditorValDel === '' || self.cmEditorDel.state.lint.marked.length > 0)) {
          self.setState({cmEditorValDelError: true, selectionError: true});
          return
        }
      } else {
        if (selectedOp === 'custom' && (cmEditorValAdd === '' || self.cmEditorAdd.state.lint.marked.length > 0)) {
          self.setState({cmEditorValAddError: true, selectionError: true});
          return
        }
      }
      if (namespace === '') {
        self.setState({namespaceError: true});
        return
      }
      self.submitOp(cat, selectedOp, deleteOp);
    }
  }

  submitOp = (cat, selectedOp, deleteOp=false) => {
    const { namespace, cmEditorValAdd, cmEditorValDel, menuState} = this.state;
    const { adapter } = this.props;
    // const fileInput = document.querySelector('#k8sfile') ;

    const data = {
      'adapter': adapter.adapter_location,
      'query': selectedOp,
      'namespace': namespace,
      'customBody': deleteOp?cmEditorValDel:cmEditorValAdd,
      'deleteOp': deleteOp? 'on':'',
    }
    
    const params = Object.keys(data).map((key) => {
      return encodeURIComponent(key) + '=' + encodeURIComponent(data[key]);
    }).join('&');
    this.props.updateProgress({showProgress: true});
    let self = this;
    dataFetch('/api/mesh/ops', { 
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
      },
      body: params
    }, result => {
      self.props.updateProgress({showProgress: false});
      menuState[cat][deleteOp?'delete':'add'] = false;
      const dlg = deleteOp?'customDialogDel':'customDialogAdd';
      self.setState({menuState, [dlg]: false});

      if (typeof result !== 'undefined'){
        self.props.enqueueSnackbar('Operation executing...', {
          variant: 'info',
          autoHideDuration: 2000,
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key) }
            >
              <CloseIcon />
            </IconButton>
          ),
        });
      }
    }, self.handleError(cat, deleteOp));
  }

  handleError = (cat, deleteOp) => {
    const self = this;
    return error => {
      const {menuState} = self.state;
      menuState[cat][deleteOp?'delete':'add'] = false;
      const dlg = deleteOp?'customDialogDel':'customDialogAdd';
      self.setState({menuState, [dlg]: false});

      self.props.updateProgress({showProgress: false});
      self.props.enqueueSnackbar(`Operation submission failed: ${error}`, {
        variant: 'error',
        action: (key) => (
          <IconButton
            key="close"
            aria-label="Close"
            color="inherit"
            onClick={() => self.props.closeSnackbar(key) }
          >
            <CloseIcon />
          </IconButton>
        ),
        autoHideDuration: 8000,
      });
    };
  }

  handleExpandClick() {
    // setExpanded(!expanded);
  }

  generateMenu(cat, isDelete, selectedAdapterOps) {
    const {menuState} = this.state;
    const ele = !isDelete?this.addIconEles[cat]:this.delIconEles[cat];
    return (
      <Menu
        id="long-menu"
        anchorEl={ele}
        keepMounted
        open={menuState[cat][isDelete?'delete':'add']}
        onClose={this.addDelHandleClick(cat, isDelete)}
        // PaperProps={{
        //   style: {
        //     maxHeight: ITEM_HEIGHT * 4.5,
        //     width: 200,
        //   },
        // }}
      >
        {selectedAdapterOps.map(({key, value}) => (
          <MenuItem key={`${key}_${new Date().getTime()}`} onClick={this.handleSubmit(cat, key, isDelete)}>
            {value}
          </MenuItem>
        ))}
      </Menu>
    );
  }

  generateYAMLEditor(cat, isDelete) {
    const {adapter} = this.props;
    const {customDialogAdd, customDialogDel, namespace, namespaceError, cmEditorValAdd, cmEditorValDel, 
      cmEditorValAddError, cmEditorValDelError} = this.state;
    const self = this;
    return (
      <Dialog
        onClose={this.handleModalClose(isDelete)}
        aria-labelledby="adapter-dialog-title"
        open={isDelete?customDialogDel:customDialogAdd}
        fullWidth={true}
        maxWidth={'md'}
      >
        <DialogTitle id="adapter-dialog-title" onClose={this.handleModalClose(isDelete)}>
          {adapter.name} Adapter - Custom YAML {isDelete?'(delete)':''}
        </DialogTitle>
        <Divider variant="fullWidth" light />
        <DialogContent>
          <Grid container spacing={5}>
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
            <Grid item xs={12}>
              <CodeMirror
                editorDidMount={editor => { 
                  if(isDelete){
                    self.cmEditorDel = editor;
                  } else {
                    self.cmEditorAdd = editor;
                  } 
}}
                value={isDelete?cmEditorValDel:cmEditorValAdd}
                options={{
                  theme: 'material',
                  lineNumbers: true,
                  lineWrapping: true,
                  gutters: ["CodeMirror-lint-markers"],
                  lint: true,
                  mode: 'text/x-yaml'
                }}
                onBeforeChange={(editor, data, value) => {
                  if(isDelete){
                    self.setState({cmEditorValDel: value});
                  } else {
                    self.setState({cmEditorValAdd: value});
                  }
                  if(isDelete){
                    if(value !== '' && self.cmEditorDel.state.lint.marked.length === 0) {
                      self.setState({ selectionError: false, cmEditorValDelError: false });  
                    } 
                  } else {
                    if(value !== '' && self.cmEditorAdd.state.lint.marked.length === 0) {
                      self.setState({ selectionError: false, cmEditorValAddError: false });  
                    }
                  }
                }
                }
                onChange={(editor, data, value) => {
                }}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <Divider variant="fullWidth" light />
        <DialogActions>
          <IconButton aria-label="Apply" color="primary" onClick={this.handleSubmit(cat, 'custom', isDelete)}>
            {/* <FontAwesomeIcon icon={faArrowRight} transform="shrink-4" fixedWidth /> */}
            {!isDelete && <PlayIcon />}
            {isDelete && <DeleteIcon />}
          </IconButton>
          
          {/* <IconButton aria-label="Delete" color="primary" onClick={this.handleSubmit('custom', isDelete)}>
            <FontAwesomeIcon icon={faTrashAlt} transform="shrink-4" fixedWidth />
          </IconButton> */}

        </DialogActions>
      </Dialog>
    );
  }

  addDelHandleClick = (cat, isDelete) => {
    const self = this;
    return () => {
      let  {menuState, customDialogAdd, customDialogDel} = self.state;
      menuState[cat][isDelete?'delete':'add'] = !menuState[cat][isDelete?'delete':'add'];

      const dlg = isDelete?'customDialogDel':'customDialogAdd';
      let dlgv = isDelete?customDialogDel:customDialogAdd;
      if (cat === 4){
        dlgv = !dlgv;
      }
      self.setState({menuState, [dlg]: dlgv});
    }
  }

  generateCardForCategory(cat) {
    if (typeof cat === 'undefined'){
      cat = 0;
    }
    const {classes, adapter} = this.props;
    const {menuState} = this.state;
    // const expanded = false;

    const selectedAdapterOps = adapter && adapter.ops? adapter.ops.filter(({category}) => typeof category === 'undefined' && cat === 0 || category === cat):[];
    let content, description;
    switch(cat){
    case 0:
      content = 'Manage Service Mesh Lifecycle';
      description = 'Deploy a service mesh or SMI adapter on your cluster.';
      break;

    case 1:
      content = 'Manage Sample Application Lifecycle';
      description = 'Deploy sample applications on/off the service mesh.';
      break;

    case 2:
      content = 'Apply Service Mesh Configuration';
      description = 'Configure your service mesh using some pre-defined options.';
      break;

    case 3:
      content = 'Validate Service Mesh Configuration';
      description = 'Validate your service mesh configuration against best practices.';
      break;

    case 4:
      content = 'Apply Custom Configuration';
      description = 'Customize the configuration of your service mesh.';
      break;
    }
    return (
      <Card className={classes.card}>
        <CardHeader
          // avatar={
          //   <Avatar aria-label="recipe" className={classes.avatar}>
          //     R
          //   </Avatar>
          // }
          // action={
          //   <IconButton aria-label="settings">
          //     <MoreVertIcon />
          //   </IconButton>
          // }
          // title="Shrimp and Chorizo Paella"
          // subheader="September 14, 2016"
          title={content}
          subheader={description}
        />
        {/* <CardMedia
          className={classes.media}
          image="/static/images/cards/paella.jpg"
          title="Paella dish"
        /> */}
        {/* <CardContent>
          <Typography variant="body2" color="textSecondary" component="p">
            This impressive paella is a perfect party dish and a fun meal to cook together with your
            guests. Add 1 cup of frozen peas along with the mussels, if you like.
          </Typography>
        </CardContent> */}
        <CardActions disableSpacing>
          <IconButton aria-label="install" ref={ch => this.addIconEles[cat] = ch} onClick={this.addDelHandleClick(cat, false)}>
            {cat !== 3 && cat !== 4?<AddIcon />:<PlayIcon />}
          </IconButton>
          {cat !== 4 && this.generateMenu(cat, false, selectedAdapterOps)}
          {cat === 4 && this.generateYAMLEditor(cat, false)}
          {cat !== 3 && <div className={classes.fileLabel}>
            <IconButton aria-label="delete" ref={ch => this.delIconEles[cat] = ch} className={classes.deleteRight} onClick={this.addDelHandleClick(cat, true)}>
              {/*<IconButton aria-label="Delete" color="primary" onClick={this.handleSubmit(key, true)}>*/}
              {/* <FontAwesomeIcon icon={faTrashAlt} transform="shrink-4" fixedWidth /> */}
              <DeleteIcon />
            </IconButton>
            {cat !== 4 && this.generateMenu(cat, true, selectedAdapterOps)}
            {cat === 4 && this.generateYAMLEditor(cat, true)}
          </div>}
          {/* <IconButton
            // className={clsx(classes.expand, {
            //   [classes.expandOpen]: expanded,
            // })}
            onClick={this.handleExpandClick}
            aria-expanded={expanded}
            aria-label="show more"
          >
            <ExpandMoreIcon />
          </IconButton> */}
        </CardActions>
      </Card>
    );
  }

  render() {
    const {classes, color, iconButtonClassName, avatarClassName, adapter, adapter_icon} = this.props;
    const {
      selectedOp,
      cmEditorValAdd,
      cmEditorValDel,
      namespace,

      selectionError,
      namespaceError,
      cmEditorValAddError,
      cmEditorValDelError,
    } = this.state;

    const filteredOps = [];
    if (adapter && adapter.ops && adapter.ops.length > 0){
      adapter.ops.forEach(({category}) => {
        if(typeof category === 'undefined'){
          category = 0;
        }
        if(filteredOps.indexOf(category) === -1){
          filteredOps.push(category);
        }
      });
      filteredOps.sort();
    }

    var self = this;
    // return (
    //   <NoSsr>
    //     <React.Fragment>
    //       <div className={classes.root}>
    //       <Grid container spacing={5}>
    //       <Grid item xs={12}>
    //         <TextField
    //           required
    //           id="namespace"
    //           name="namespace"
    //           label="Namespace"
    //           fullWidth
    //           value={namespace}
    //           error={namespaceError}
    //           margin="normal"
    //           variant="outlined"
    //           onChange={this.handleChange('namespace')}
    //         />
    //       </Grid>
    //       <Grid item xs={12}>
    //         <FormControl required error={selectionError} component="fieldset">
    //         <FormLabel component="legend">{`Play with ${adapter.name}`}</FormLabel>
    //         <FormGroup
    //         // aria-label={`Play with ${adapter.name}`}
    //         // name="query"
    //         // className={classes.group}
    //         // value={selectedOp}
    //         // onChange={this.handleChange('selectedOp')}
    //         >
    //         {Object.keys(adapter.ops).filter(({key}) => key !== 'custom').map(({key}, ind) => (
    //           <div>
    //             <IconButton aria-label="Apply" color="primary" onClick={this.handleSubmit(key, false)}>
    //             <FontAwesomeIcon icon={faArrowRight} transform="shrink-4" fixedWidth />
    //             </IconButton>
                
    //             <IconButton aria-label="Delete" color="primary" onClick={this.handleSubmit(key, true)}>
    //             <FontAwesomeIcon icon={faTrashAlt} transform="shrink-4" fixedWidth />
    //             </IconButton>

    //             {adapter.ops[ind].value}
    //           </div>
              
    //         ))}
    //        </FormGroup>
    //        </FormControl>
    //        </Grid>
    //        <Grid item xs={12}>
    //        <React.Fragment>
    //         <div className={classes.buttons}>
    //           <Button color="inherit" onClick={this.handleModalOpen} className={classes.custom}>
    //             <FontAwesomeIcon icon={faTerminal} fixedWidth className={classes.padRight} />
    //             Custom YAML
    //             <FontAwesomeIcon icon={faExternalLinkSquareAlt} fixedWidth className={classes.padLeft} />
    //           </Button>
    //         </div>
    //           <Dialog
    //             onClose={this.handleModalClose}
    //             aria-labelledby="adapter-dialog-title"
    //             open={this.state.customDialog}
    //             fullWidth={true}
    //             maxWidth={'md'}
    //           >
    //             <DialogTitle id="adapter-dialog-title" onClose={this.handleModalClose}>
    //               {adapter.name} Adapter - Custom YAML
    //             </DialogTitle>
    //             <Divider variant="fullWidth" light />
    //             <DialogContent>
    //             <Grid container spacing={5}>
    //             <Grid item xs={12}>
    //               <TextField
    //                 required
    //                 id="namespace"
    //                 name="namespace"
    //                 label="Namespace"
    //                 fullWidth
    //                 value={namespace}
    //                 error={namespaceError}
    //                 margin="normal"
    //                 variant="outlined"
    //                 onChange={this.handleChange('namespace')}
    //               />
    //             </Grid>
    //             <Grid item xs={12}>
    //               <CodeMirror
    //                   editorDidMount={editor => { this.cmEditor = editor }}
    //                   value={cmEditorVal}
    //                   options={{
    //                     mode: 'yaml',
    //                     theme: 'material',
    //                     lineNumbers: true,
    //                     lineWrapping: true,
    //                     gutters: ["CodeMirror-lint-markers"],
    //                     lint: true,
    //                     mode: "text/x-yaml"
    //                   }}
    //                   onBeforeChange={(editor, data, value) => {
    //                     this.setState({cmEditorVal: value});
    //                     if(value !== '' && this.cmEditor.state.lint.marked.length === 0) {
    //                       this.setState({ selectionError: false, cmEditorValError: false });  
    //                     }
    //                   }}
    //                   onChange={(editor, data, value) => {
    //                   }}
    //                 />
    //               </Grid>
    //               </Grid>
    //             </DialogContent>
    //             <Divider variant="fullWidth" light />
    //             <DialogActions>
    //               <IconButton aria-label="Apply" color="primary" onClick={this.handleSubmit('custom', false)}>
    //                 <FontAwesomeIcon icon={faArrowRight} transform="shrink-4" fixedWidth />
    //               </IconButton>
                  
    //               <IconButton aria-label="Delete" color="primary" onClick={this.handleSubmit('custom', false)}>
    //                 <FontAwesomeIcon icon={faTrashAlt} transform="shrink-4" fixedWidth />
    //               </IconButton>

    //             </DialogActions>
    //           </Dialog>
    //         </React.Fragment>
    //        </Grid>
    //       </Grid>
    //       </div>
    //     </React.Fragment>
    //   </NoSsr>
    // )

    if(this.props.adapCount > 1){
      return (
        <NoSsr>
          <React.Fragment>
            <div className={classes.root}>
              <Grid container spacing={5}>
                <Grid item xs={12} sm={6}>
                  <TextField
                    disabled
                    id="ap"
                    name="ap"
                    label="Adapter URL"
                    fullWidth
                    value={adapter.adapter_location}
                    margin="normal"
                    variant="outlined"
                    // onChange={this.handleChange('namespace')}
                    InputProps={{
                      startAdornment: (
                        <InputAdornment position="start">
                          {adapter_icon}
                        </InputAdornment>
                      ),
                    }}
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
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
                {filteredOps.map(val => (
                  <Grid item xs={12} sm={4}>
                    {this.generateCardForCategory(val)}
                  </Grid>
                ))}
                {/* <Grid item xs={12} sm={3}>
              {this.generateCardForCategory(1)}
            </Grid>
            <Grid item xs={12} sm={3}>
            {this.generateCardForCategory(2)}
            </Grid>
            <Grid item xs={12} sm={3}>
            {this.generateCardForCategory(3)}
            </Grid> */}
              </Grid>
            </div>
          </React.Fragment>
        </NoSsr>
      ) 
}
  
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.root}>
            <Grid container spacing={5}>
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
              {filteredOps.map(val => (
                <Grid item xs={12} sm={4}>
                  {this.generateCardForCategory(val)}
                </Grid>
              ))}
            </Grid>
          </div>
        </React.Fragment>
      </NoSsr>
    )
  
  }
}

MesheryAdapterPlayComponent.propTypes = {
  classes: PropTypes.object.isRequired,
  // index: PropTypes.number.isRequired,
  adapter: PropTypes.object.isRequired,
  // adapter_icon: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    updateProgress: bindActionCreators(updateProgress, dispatch),
  }
}

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps
)(withRouter(withSnackbar(MesheryAdapterPlayComponent))));
