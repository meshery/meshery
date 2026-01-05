import React from 'react';
import { Box, Typography, Chip, Divider } from '@sistent/sistent';
import { CONNECTION_STATES } from '../../utils/Enum';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { ConnectedStateChip, DeletedStateChip, DisconnectedStateChip, DiscoveredStateChip, IgnoredStateChip, NotFoundStateChip, RegisteredStateChip } from './ConnectionChip';

// State descriptions mapping
export const CONNECTION_STATE_DESCRIPTIONS = {
  [CONNECTION_STATES.DISCOVERED]: {
    title: 'Discovered',
  },
  [CONNECTION_STATES.REGISTERED]: {
    title: 'Registered',
  },
  [CONNECTION_STATES.CONNECTED]: {
    title: 'Connected',
  },
  [CONNECTION_STATES.IGNORED]: {
    title: 'Ignored',
  },
  [CONNECTION_STATES.MAINTENANCE]: {
    title: 'Maintenance',
  },
  [CONNECTION_STATES.DISCONNECTED]: {
    title: 'Disconnected',
  },
  [CONNECTION_STATES.DELETED]: {
    title: 'Deleted',
  },
  [CONNECTION_STATES.NOTFOUND]: {
    title: 'Not Found',
  },
};

export const CONNECTION_STATE_TRANSITIONS = {

  [`${CONNECTION_STATES.CONNECTED}_${CONNECTION_STATES.DISCONNECTED}`]: {
    arrowText: 'Removes Operator, Stays registered',
    doubleArrow: true,
  },
  [`${CONNECTION_STATES.DISCONNECTED}_${CONNECTION_STATES.CONNECTED}`]: {
    arrowText: 'Deploy Operator',
    doubleArrow: true,
  },
  [`${CONNECTION_STATES.CONNECTED}_${CONNECTION_STATES.IGNORED}`]: {
    arrowText: 'Unplanned Maintenance',
  },
  [`${CONNECTION_STATES.CONNECTED}_${CONNECTION_STATES.DELETED}`]: {
    arrowText: 'Undeploy Operator and unregister',
    doubleArrow: false,
  },
  [`${CONNECTION_STATES.DISCONNECTED}_${CONNECTION_STATES.DELETED}`]: {
    arrowText: 'Unregister',
  },
};

export const getTransitionConfig = (fromState, toState) => {
  const key = `${fromState}_${toState}`;
  return CONNECTION_STATE_TRANSITIONS[key] || {
    arrowText: null,
    leftText: null,
    rightText: null,
    doubleArrow: false,
  };
};

export const ConnectionStateTransitionVisual = ({ 
  fromState, 
  toState, 
  showLabels = true,
  size = 'small',
  arrowText = undefined, 
  leftText = undefined,
  rightText = undefined,
  doubleArrow = undefined,
  compact = false,
  useTransitionConfig = true
}) => {
  const transitionConfig = useTransitionConfig ? getTransitionConfig(fromState, toState) : {};

  const finalArrowText = arrowText !== undefined ? arrowText : transitionConfig.arrowText;
  const finalLeftText = leftText !== undefined ? leftText : transitionConfig.leftText;
  const finalRightText = rightText !== undefined ? rightText : transitionConfig.rightText;
  const finalDoubleArrow = doubleArrow !== undefined ? doubleArrow : transitionConfig.doubleArrow;

  const getStateChip = (state) => {
    const chipProps = { 
      label: showLabels ? CONNECTION_STATE_DESCRIPTIONS[state]?.title || state : '', 
      size,
      variant: 'filled'
    };

    switch (state) {
      case CONNECTION_STATES.CONNECTED:
        return (
          <ConnectedStateChip value={CONNECTION_STATE_DESCRIPTIONS[state]?.title} {...chipProps} />
        );
      case CONNECTION_STATES.REGISTERED:
        return (
          <RegisteredStateChip value={CONNECTION_STATE_DESCRIPTIONS[state]?.title} {...chipProps} />
        );
      case CONNECTION_STATES.DISCOVERED:
        return (
          <DiscoveredStateChip value={CONNECTION_STATE_DESCRIPTIONS[state]?.title} {...chipProps} />
        );
      case CONNECTION_STATES.IGNORED:
        return <IgnoredStateChip value={CONNECTION_STATE_DESCRIPTIONS[state]?.title} {...chipProps} />;
      case CONNECTION_STATES.DELETED:
        return <DeletedStateChip value={CONNECTION_STATE_DESCRIPTIONS[state]?.title} {...chipProps} />;
      case CONNECTION_STATES.DISCONNECTED:
        return (
          <DisconnectedStateChip
            value={CONNECTION_STATE_DESCRIPTIONS[state]?.title}
            {...chipProps}
          />
        );
      case CONNECTION_STATES.NOTFOUND:
        return (
          <NotFoundStateChip value={CONNECTION_STATE_DESCRIPTIONS[state]?.title} {...chipProps} />
        );
      default:
        return <Chip {...chipProps} />;
    }
  };

  const ArrowComponent = () => (
    <Box
      sx={{
        position: 'relative',
        display: 'flex',
        alignItems: 'center',
        mx: compact ? 0.5 : 1,
      }}
    >
      {finalArrowText && (
        <Typography
          variant="caption"
          sx={{
            position: 'absolute',
            top: -25,
            left: '50%',
            transform: 'translateX(-50%)',
            whiteSpace: 'nowrap',
            fontSize: '0.7rem',
          }}
        >
          {finalArrowText}
        </Typography>
      )}
      {finalDoubleArrow ? (
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
          <ArrowForwardIcon
            sx={{ color: 'text.secondary', fontSize: 'large'}}
          />
          <ArrowBackIcon sx={{ color: 'text.secondary',fontSize: 'large', mt: -0.2 }} />
        </Box>
      ) : (
        <ArrowForwardIcon sx={{ color: 'text.secondary', fontSize: 'large' }} />
      )}
    </Box>
  );

  return (
    <Box sx={{ 
      display: 'flex', 
      alignItems: 'center', 
      gap: compact ? 0.5 : 1,
      flexWrap: 'wrap',
      justifyContent: 'center',
      position: 'relative',
      minHeight: finalArrowText ? 40 : 'auto'
    }}>
      {/* Left text */}
      {finalLeftText && (
        <Typography 
          variant="caption" 
          sx={{ 
            fontSize: '0.7rem',
            color: 'text.secondary',
            fontWeight: 500,
            mr: 0.5
          }}
        >
          {finalLeftText}
        </Typography>
      )}
      
      {/* From State Chip */}
      {getStateChip(fromState)}
      
      {/* Arrow with optional text */}
      <ArrowComponent />
      
      {/* To State Chip */}
      {getStateChip(toState)}
      
      {/* Right text */}
      {finalRightText && (
        <Typography 
          variant="caption" 
          sx={{ 
            fontSize: '0.7rem',
            color: 'text.secondary',
            fontWeight: 500,
            ml: 0.5
          }}
        >
          {finalRightText}
        </Typography>
      )}
    </Box>
  );
};
