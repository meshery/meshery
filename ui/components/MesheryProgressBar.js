import { useEffect, useRef } from 'react';
import { connect } from 'react-redux';
import { LinearProgress } from '@material-ui/core';
import { enqueueSnackbar, closeSnackbar } from 'notistack';

const MesheryProgressBar = ({ showProgress }) => {
  const keyRef = useRef('');

  useEffect(() => {
    if (showProgress) {
      keyRef.current = enqueueSnackbar(
        <div style={{ width: 250 }}>
          <LinearProgress />
        </div>,
        { variant: 'default', persist: true },
      );
    } else {
      closeSnackbar(keyRef.current);
    }
  }, [showProgress]);

  return null;
};

const mapStateToProps = (state) => ({ showProgress: state.get('showProgress') });

export default connect(mapStateToProps)(MesheryProgressBar);
