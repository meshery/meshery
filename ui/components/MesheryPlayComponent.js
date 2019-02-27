import IconButton from '@material-ui/core/IconButton';
import Avatar from '@material-ui/core/Avatar';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux'
import { updateUser } from '../lib/store';

import MenuList from '@material-ui/core/MenuList';
import Grow from '@material-ui/core/Grow';
import MenuItem from '@material-ui/core/MenuItem';
import Popper from '@material-ui/core/Popper';
import Paper from '@material-ui/core/Paper';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import NoSsr from '@material-ui/core/NoSsr';
import dataFetch from '../lib/data-fetch';
import {Controlled as CodeMirror} from 'react-codemirror2'
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/material.css';
import 'codemirror/addon/lint/lint.css';
// import 'codemirror/mode/yaml/yaml';

// import dynamic from 'next/dynamic'
// dynamic(() => import('codemirror/mode/yaml/yaml'), {
//   ssr: false
// })
// import 'codemirror/mode/yaml/yaml';
// import 'codemirror/addon/lint/lint';
// import 'codemirror/addon/lint/yaml-lint';
// import 'js-yaml';
if (typeof window !== 'undefined') { 
  require('codemirror/mode/yaml/yaml'); 
  require('codemirror/addon/lint/lint');
  require('codemirror/addon/lint/yaml-lint');
  if (typeof window.jsyaml === 'undefined'){
    window.jsyaml = require('js-yaml');
  }
}


class MesheryPlayComponent extends React.Component {

  state = {
    user: null,
    open: false,
    cmEditorVal: '',
  }

  handleToggle = () => {
    this.setState(state => ({ open: !state.open }));
  };

  handleClose = event => {
    if (this.anchorEl.contains(event.target)) {
      return;
    }
    this.setState({ open: false });
  };

  render() {
    const {color, iconButtonClassName, avatarClassName, ...other} = this.props;
    let avatar_url, user_id;
    if (this.state.user !== null){
      avatar_url = this.state.user.avatar_url;
      user_id = this.state.user.user_id;
    }
    const { cmEditorVal } = this.state;
    return (
      <NoSsr>
      <CodeMirror
          value={cmEditorVal}
          options={{
            mode: 'yaml',
            theme: 'material',
            lineNumbers: true,
            lineWrapping: true,
            gutters: ["CodeMirror-lint-markers"],
            lint: true,
            mode: "text/x-yaml"
          }}
          onBeforeChange={(editor, data, value) => {
            this.setState({cmEditorVal: value});
          }}
          onChange={(editor, data, value) => {
          }}
        />
      </NoSsr>
    )
  }
}

// MesheryPlayComponent.propTypes = {
//   classes: PropTypes.object.isRequired,
// };

const mapDispatchToProps = dispatch => {
    return {
        // updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch)
    }
}
const mapStateToProps = state => {
    // console.log("header - mapping state to props. . . new title: "+ state.get("page").get("title"));
    // console.log("state: " + JSON.stringify(state));
    // const k8sconfig = state.get("k8sConfig");
    let newprops = {};
    // if (typeof k8sconfig !== 'undefined'){
    //   newprops = { 
    //     inClusterConfig: k8sconfig.get('inClusterConfig'),
    //     // k8sfile: '', 
    //     contextName: k8sconfig.get('contextName'), 
    //     meshLocationURL: k8sconfig.get('meshLocationURL'), 

    //     reconfigureCluster: k8sconfig.get('reconfigureCluster'),
    //   }
    // }
    return newprops;
}

export default connect(
    mapStateToProps,
    mapDispatchToProps
  )(MesheryPlayComponent);
