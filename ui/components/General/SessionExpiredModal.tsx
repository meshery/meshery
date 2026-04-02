import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import {
  Button,
  Typography,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
} from '@sistent/sistent';
import { selectIsSessionExpired, clearSessionExpired } from '../../store/slices/session';

/**
 * Modal shown when a 401 is detected via the auth middleware.
 * Gives users a chance to re-authenticate without losing client-side state.
 */
export default function SessionExpiredModal() {
  const isExpired = useSelector(selectIsSessionExpired);
  const dispatch = useDispatch();

  const handleReLogin = () => {
    dispatch(clearSessionExpired());
    window.location.href = '/user/login';
  };

  const handleDismiss = () => {
    dispatch(clearSessionExpired());
  };

  if (!isExpired) {
    return null;
  }

  return (
    <Dialog open={isExpired} data-testid="meshery-session-expired-modal">
      <DialogTitle>Session Expired</DialogTitle>
      <DialogContent>
        <Typography>
          Your session has expired. You can log in again to continue, or dismiss this to stay on the
          current page.
        </Typography>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleDismiss} data-testid="meshery-session-dismiss-btn">
          Dismiss
        </Button>
        <Button variant="contained" onClick={handleReLogin} data-testid="meshery-session-login-btn">
          Log In
        </Button>
      </DialogActions>
    </Dialog>
  );
}
