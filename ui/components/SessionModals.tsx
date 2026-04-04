import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import {
  Modal,
  ModalBody,
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalFooter,
  Typography,
  CircularProgress,
} from '@sistent/sistent';
import { SessionState } from '@/store/slices/sessions';
import { recordActivity, triggerSessionExpired } from 'lib/sessionTimer';

const SessionExpiringModal = () => {
  const sessionStatus = useSelector(
    (state: { sessions: { status: SessionState } }) => state.sessions.status,
  );
  const [isRefreshing, setIsRefreshing] = useState(false);

  const handleStayLoggedIn = () => {
    setIsRefreshing(true);
    fetch('/api/user/prefs', { credentials: 'include' })
      .then(() => {
        recordActivity();
      })
      .catch(() => {
        triggerSessionExpired();
      })
      .finally(() => {
        setIsRefreshing(false);
      });
  };

  const handleLogout = () => {
    window.location.href = '/user/logout';
  };

  return (
    <Modal
      open={sessionStatus === 'expiring'}
      closeModal={handleStayLoggedIn}
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
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            gap: '1rem',
            width: '100%',
          }}
        >
          <ModalButtonSecondary onClick={handleLogout}>Log Out</ModalButtonSecondary>
          <ModalButtonPrimary onClick={handleStayLoggedIn} disabled={isRefreshing} autoFocus>
            {isRefreshing ? (
              <div style={{ display: 'flex', alignItems: 'center' }}>
                <CircularProgress size={20} style={{ marginRight: '0.5rem', color: '#fff' }} />
                <Typography variant="body1">Refreshing...</Typography>
              </div>
            ) : (
              'Stay Logged In'
            )}
          </ModalButtonPrimary>
        </div>
      </ModalFooter>
    </Modal>
  );
};

const SessionExpiredModal = () => {
  const sessionStatus = useSelector(
    (state: { sessions: { status: SessionState } }) => state.sessions.status,
  );
  const isExpired = sessionStatus === 'expired';

  const redirectTo =
    typeof window !== 'undefined' && window.location.host.endsWith('3000')
      ? '/user/login'
      : typeof window !== 'undefined'
        ? window.location.href
        : '/';

  const handleLogin = () => {
    window.location.href = redirectTo;
  };

  // Auto-redirect after 5 seconds
  useEffect(() => {
    if (!isExpired) return;

    const timer = setTimeout(() => {
      window.location.href = redirectTo;
    }, 5000);

    return () => clearTimeout(timer);
  }, [isExpired, redirectTo]);

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
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            width: '100%',
          }}
        >
          <ModalButtonPrimary onClick={handleLogin} autoFocus>
            Log In
          </ModalButtonPrimary>
        </div>
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
