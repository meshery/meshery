import { Component } from 'react';
import { withSnackbar } from 'notistack';
import { connect } from 'react-redux';
import { LinearProgress } from '@material-ui/core';

class MesheryProgressBar extends Component {
    key = '';

    shouldComponentUpdate(nextProps) {
      const { showProgress } = this.props;
      // if ((this.key !== '' && !showProgress) || (this.key === '' && showProgress)){
      //     return true;
      // }
      return showProgress !== nextProps.showProgress;
    }

    componentDidUpdate() {
      const { showProgress } = this.props;
      if (showProgress) {
        this.key = this.props.enqueueSnackbar(<div style={{ width : 250 }}><LinearProgress /></div>, { variant : 'default',
          persist : true, });
      } else {
        this.props.closeSnackbar(this.key);
      }
    }

    render() {
      return null;
    }
}

const mapStateToProps = (state) => ({ showProgress : state.get('showProgress') });

export default connect(
  mapStateToProps,
)(withSnackbar(MesheryProgressBar));
