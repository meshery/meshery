import { LinearProgress } from '@layer5/sistent';
import React, { useEffect, useRef } from 'react';
import { useSnackbar } from 'notistack';
import { useSelectorRtk } from '@/store/hooks';

const MesheryProgressBar = () => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const { showProgress } = useSelectorRtk((state) => state.ui);
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

export default MesheryProgressBar;
