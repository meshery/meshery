import { withSnackbar } from 'notistack';
import { useEffect } from 'react';
import { connect, useSelector } from 'react-redux';
import { LinearProgress } from '@material-ui/core';

let key = '';

const MesheryProgressBar = ({ enqueueSnackbar, closeSnackbar }) => {
  const showProgress = useSelector((state) => state.get('showProgress'));

  useEffect(() => {
    if (showProgress) {
      key = enqueueSnackbar(
        <div style={{ width: 250 }}>
          <LinearProgress />
        </div>,
        { variant: 'default', persist: true },
      );
    } else if (key) {
      closeSnackbar(key);
      key = '';
    }
  }, [showProgress, enqueueSnackbar, closeSnackbar]);

  return null;
};

export default connect()(withSnackbar(MesheryProgressBar));
