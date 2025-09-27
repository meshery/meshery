import React from 'react';
import {
  Modal,
  ModalBody,
  ModalFooter,
  Typography,
  Box,
  Button
} from '@sistent/sistent';
import {
  ConnectionStateTransitionVisual,
  CONNECTION_STATE_DESCRIPTIONS,
} from './ConnectionStateHelpers';
import { CONNECTION_STATES } from '../../utils/Enum';

export const ConnectionStateTransitionDialog = ({
  open,
  onClose,
  onConfirm,
  currentState,
  newState,
  subtitle,
  isLoading = false,
}) => {
  const currentStateInfo = CONNECTION_STATE_DESCRIPTIONS[currentState];

  const newStateInfo = CONNECTION_STATE_DESCRIPTIONS[newState];

  const handleConfirm = () => {
    onConfirm();
  };

  const getTransitionType = () => {
    const stateOrder = [
      CONNECTION_STATES.DISCOVERED,
      CONNECTION_STATES.REGISTERED,
      CONNECTION_STATES.CONNECTED,
      CONNECTION_STATES.MAINTENANCE,
      CONNECTION_STATES.DISCONNECTED,
      CONNECTION_STATES.IGNORED,
      CONNECTION_STATES.DELETED,
    ];

    const currentIndex = stateOrder.indexOf(currentState);
    const newIndex = stateOrder.indexOf(newState);

    if (newState === CONNECTION_STATES.DELETED) return 'destructive';
    if (newState === CONNECTION_STATES.IGNORED) return 'warning';
    if (newIndex > currentIndex) return 'progressive';
    if (newIndex < currentIndex) return 'regressive';
    return 'neutral';
  };

  const transitionType = getTransitionType();

  if (!open) {
    return null;
  }

  return (
    <Modal
      open={open}
      onClose={onClose}
      maxWidth="sm"
      title="Transition Connection State"
      fullWidth
    >
      <ModalBody>
        <Box sx={{ mb: 3 }}>
          <Typography variant="subtitle2" sx={{ mb: 1 }}>
            State Transition: {currentStateInfo.title} → {newStateInfo.title}
          </Typography>
          <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
            <ConnectionStateTransitionVisual
              fromState={currentState}
              toState={newState}
              showLabels={true}
              size="medium"
              useTransitionConfig={true}
            />
          </Box>
        </Box>
        <Box sx={{ mb: 2 }}>
          <Typography variant="body2">
            {subtitle}
          </Typography>
        </Box>
      </ModalBody>

      <ModalFooter variant='filled' helpText="Learn more about the [lifecycle of connections and the behavior of state transitions](https://docs.meshery.io/concepts/logical/connections) in Meshery Docs.">
        <Button
          onClick={handleConfirm}
          variant="contained"
          color={transitionType === 'destructive' ? 'error' : 'primary'}
          disabled={isLoading}
        >
          {isLoading ? 'Processing...' : `Transition to ${newStateInfo.title}`}
        </Button>
      </ModalFooter>
    </Modal>
  );
};
