import { connect } from 'react-redux';
import { LinearProgress } from '@material-ui/core';
import React, { useEffect, useRef } from 'react';
import { useSnackbar } from 'notistack';

const MesheryProgressBar = ({ showProgress }) => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const snackbarKey = useRef(null);

  useEffect(() => {
    if (showProgress) {
      snackbarKey.current = enqueueSnackbar(
        <div style={{ width: 250 }}>
          <LinearProgress />
        </div>,
        { variant: 'default', persist: true },
      );
    } else if (snackbarKey.current) {
      closeSnackbar(snackbarKey.current);
      snackbarKey.current = null;
    }
  }, [showProgress, enqueueSnackbar, closeSnackbar]);

  return null;
};

const mapStateToProps = (state) => ({ showProgress: state.get('showProgress') });

export default connect(mapStateToProps)(MesheryProgressBar);
