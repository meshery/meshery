import React, { useEffect, useState } from 'react';
import {
  alpha,
  BLACK,
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  SAFFRON,
  Typography,
  useTheme,
  WarningIcon,
} from '@sistent/sistent';

const SESSION_WARNING_EVENT = 'meshery:session-warning';
const SESSION_WARNING_DISMISS_EVENT = 'meshery:session-warning-dismiss';
const SESSION_EXPIRED_EVENT = 'meshery:session-expired';
const CARD_WIDTH_PX = 445;
const HEADER_HEIGHT_PX = 56;
const CONTENT_TOP_PX = 34;
const CONTENT_BOTTOM_PX = 30;
const CONTENT_SIDE_PX = 24;
const STRIP_HEIGHT_PX = 22;

export default function SessionAuthModals() {
  const theme = useTheme();
  const [showWarning, setShowWarning] = useState(false);
  const [showExpired, setShowExpired] = useState(false);
  const [redirectIn, setRedirectIn] = useState(3);

  useEffect(() => {
    if (typeof window === 'undefined') return;

    const openWarning = () => setShowWarning(true);
    const closeWarning = () => setShowWarning(false);
    const openExpired = () => {
      setShowWarning(false);
      setRedirectIn(3);
      setShowExpired(true);
    };

    window.addEventListener(SESSION_WARNING_EVENT, openWarning);
    window.addEventListener(SESSION_WARNING_DISMISS_EVENT, closeWarning);
    window.addEventListener(SESSION_EXPIRED_EVENT, openExpired);

    return () => {
      window.removeEventListener(SESSION_WARNING_EVENT, openWarning);
      window.removeEventListener(SESSION_WARNING_DISMISS_EVENT, closeWarning);
      window.removeEventListener(SESSION_EXPIRED_EVENT, openExpired);
    };
  }, []);

  const handleStayLoggedIn = () => {
    fetch('/api/user/prefs', { credentials: 'include' }).catch(() => {});
    setShowWarning(false);
  };

  const handleLogOut = () => {
    window.location.href = '/user/logout';
  };

  const handleLogIn = () => {
    const redirectTo = window.location.host.endsWith('3000') ? '/user/login' : window.location.href;
    window.location.href = redirectTo;
  };

  useEffect(() => {
    if (!showExpired) return;

    if (redirectIn <= 0) {
      handleLogIn();
      return;
    }

    const timeout = setTimeout(() => setRedirectIn((prev) => prev - 1), 1000);
    return () => clearTimeout(timeout);
  }, [showExpired, redirectIn]);

  return (
    <>
      <Dialog open={showWarning} data-testid="meshery-session-warning-modal">
        <DialogTitle>Session is about to expire</DialogTitle>
        <DialogContent>
          <Typography>
            Your session will expire in a few minutes due to inactivity. Would you like to stay
            logged in?
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleLogOut} data-testid="meshery-session-warning-logout-btn">
            Log Out
          </Button>
          <Button
            variant="contained"
            onClick={handleStayLoggedIn}
            data-testid="meshery-session-warning-stay-btn"
          >
            Stay Logged In
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog
        open={showExpired}
        data-testid="meshery-session-expired-modal"
        disableEscapeKeyDown
        BackdropProps={{
          style: {
            backgroundColor: alpha(BLACK, 0.35),
          },
        }}
        PaperProps={{
          style: {
            width: '100%',
            maxWidth: `${CARD_WIDTH_PX}px`,
            borderRadius: '10px',
            overflow: 'hidden',
            boxShadow: `0 18px 42px ${alpha(BLACK, 0.35)}`,
            borderBottom: `${STRIP_HEIGHT_PX}px solid ${SAFFRON}`,
          },
        }}
      >
        <DialogTitle
          style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            position: 'relative',
            background: SAFFRON,
            color: theme.palette.text.default,
            fontWeight: 500,
            fontSize: '18px',
            lineHeight: '26px',
            minHeight: `${HEADER_HEIGHT_PX}px`,
            padding: '0 16px',
          }}
        >
          <Box
            sx={{
              width: '24px',
              height: '24px',
              display: 'inline-flex',
              alignItems: 'center',
              justifyContent: 'center',
              position: 'absolute',
              left: '16px',
              top: '50%',
              transform: 'translateY(-50%)',
            }}
          >
            <WarningIcon color={BLACK} width={24} height={24} />
          </Box>
          Session Expired
        </DialogTitle>
        <DialogContent sx={{ padding: 0 }}>
          <Box
            sx={{
              backgroundColor: theme.palette.background.default,
              padding: `${CONTENT_TOP_PX}px ${CONTENT_SIDE_PX}px ${CONTENT_BOTTOM_PX}px ${CONTENT_SIDE_PX}px`,
              textAlign: 'center',
            }}
          >
            <Typography
              sx={{ color: theme.palette.text.default, marginBottom: '12px', fontSize: '15px' }}
            >
              User not authenticated
            </Typography>
            <Typography sx={{ color: theme.palette.text.default, fontSize: '15px' }}>
              You will be redirected to Login page in {redirectIn}
            </Typography>
          </Box>
        </DialogContent>
      </Dialog>
    </>
  );
}
