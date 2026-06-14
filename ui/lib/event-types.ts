export const NOTIFICATION_STATUS = {
  VIEWED: 'viewed',
  NEW: 'new',
};

export const EVENT_TYPES = {
  SUCCESS: {
    type: 'success',
  },
  DEFAULT: {
    type: 'default',
  },
  INFO: {
    type: 'info',
  },
  WARNING: {
    type: 'warning',
  },
  ERROR: {
    type: 'error',
  },
};

export const SERVER_EVENT_TYPES = {
  0: EVENT_TYPES.SUCCESS,
  1: EVENT_TYPES.WARNING,
  2: EVENT_TYPES.ERROR,
};
