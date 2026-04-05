import React, { useEffect, useRef, useState } from 'react';
import { useSelector } from 'react-redux';
import {
  Box,
  Modal,
  ModalBody,
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalFooter,
  Typography,
  CircularProgress,
  styled,
} from '@sistent/sistent';
import { ExpirationState } from '@/store/slices/sessions';
import { recordActivity, triggerSessionExpired } from 'lib/sessionTimer';

const selectSessionStatus = (state: { sessions: { status: ExpirationState } }) =>
  state.sessions.status;

function getLoginRedirectUrl() {
  return '/user/login?redirect=' + encodeURIComponent(window.location.pathname);
}

const FooterActions = styled(Box)({
  display: 'flex',
  justifyContent: 'center',
  gap: '1rem',
  width: '100%',
});

const SessionExpiringModal = () => {
  const sessionStatus = useSelector(selectSessionStatus);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const isMountedRef = useRef(true);

  useEffect(() => {
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  const handleStayLoggedIn = () => {
    setIsRefreshing(true);
    fetch('/api/user/prefs', { credentials: 'include' })
      .then((res) => {
        if (res.ok) {
          recordActivity();
        } else {
          triggerSessionExpired();
        }
      })
      .catch(() => {
        triggerSessionExpired();
      })
      .finally(() => {
        if (isMountedRef.current) {
          setIsRefreshing(false);
        }
      });
  };

  const handleLogout = () => {
    window.location.href = '/user/logout';
  };

  return (
    <Modal
      open={sessionStatus === 'expiring'}
      closeModal={() => recordActivity()}
      title="Session Expiring Soon"
      maxWidth="xs"
      disableEscapeKeyDown
      disableBackdropClick
      aria-describedby="session-expiring-description"
    >
      <ModalBody>
        <Typography variant="body1" id="session-expiring-description">
          Your session will expire in a few minutes due to inactivity. Would you like to stay logged
          in?
        </Typography>
      </ModalBody>
      <ModalFooter variant="filled">
        <FooterActions>
          <ModalButtonSecondary onClick={handleLogout}>Log Out</ModalButtonSecondary>
          <ModalButtonPrimary onClick={handleStayLoggedIn} disabled={isRefreshing} autoFocus>
            {isRefreshing ? (
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <CircularProgress size={20} sx={{ mr: '0.5rem', color: 'common.white' }} />
                <Typography variant="body1">Refreshing...</Typography>
              </Box>
            ) : (
              'Stay Logged In'
            )}
          </ModalButtonPrimary>
        </FooterActions>
      </ModalFooter>
    </Modal>
  );
};

const SessionExpiredModal = () => {
  const sessionStatus = useSelector(selectSessionStatus);
  const isExpired = sessionStatus === 'expired';

  const handleLogin = () => {
    window.location.href = getLoginRedirectUrl();
  };

  useEffect(() => {
    if (!isExpired) return;

    const timer = setTimeout(() => {
      window.location.href = getLoginRedirectUrl();
    }, 5000);

    return () => clearTimeout(timer);
  }, [isExpired]);

  return (
    <Modal
      open={isExpired}
      closeModal={handleLogin}
      title="Session Expired"
      maxWidth="xs"
      disableEscapeKeyDown
      disableBackdropClick
      aria-describedby="session-expired-description"
    >
      <ModalBody>
        <Typography variant="body1" id="session-expired-description">
          Your session has expired. You will be redirected to log in.
        </Typography>
      </ModalBody>
      <ModalFooter variant="filled">
        <FooterActions>
          <ModalButtonPrimary onClick={handleLogin} autoFocus>
            Log In
          </ModalButtonPrimary>
        </FooterActions>
      </ModalFooter>
    </Modal>
  );
};

const SessionModals = () => {
  return (
    <>
      <SessionExpiringModal />
      <SessionExpiredModal />
    </>
  );
};

export default SessionModals;
