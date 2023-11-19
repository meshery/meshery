import { NOTIFICATIONCOLORS } from '../../themes';
import AlertIcon from '../../assets/icons/AlertIcon';
import ErrorIcon from '../../assets/icons/ErrorIcon.js';
import { Colors } from '../../themes/app';
import ReadIcon from '../../assets/icons/ReadIcon';
import Ajv from 'ajv';
import _ from 'lodash';
import InfoIcon from '../../assets/icons/InfoIcon';

export const SEVERITY = {
  INFO: 'informational',
  ERROR: 'error',
  WARNING: 'warning',
  SUCCESS: 'success',
};

// This class is required to add to any svg or button that opens notification center
// To prevent the clickaway listner from blocking it
export const NOTIFICATION_CENTER_TOGGLE_CLASS = 'toggle-notification-center';

export const SEVERITY_TO_NOTIFICATION_TYPE_MAPPING = {
  [SEVERITY.INFO]: 'info',
  [SEVERITY.ERROR]: 'error',
  [SEVERITY.WARNING]: 'warning',
  [SEVERITY.SUCCESS]: 'success',
};

export const STATUS = {
  READ: 'read',
  UNREAD: 'unread',
};

export const STATUS_STYLE = {
  [STATUS.READ]: {
    icon: ReadIcon,
    color: Colors.charcoal,
    darkColor: '#BCC7CC',
  },
};

export const SEVERITY_STYLE = {
  [SEVERITY.INFO]: {
    icon: InfoIcon,
    color: NOTIFICATIONCOLORS.INFO,
    darkColor: NOTIFICATIONCOLORS.INFO,
  },
  [SEVERITY.ERROR]: {
    icon: ErrorIcon,
    color: NOTIFICATIONCOLORS.ERROR,
    darlColor: NOTIFICATIONCOLORS.ERROR_DARK,
  },
  [SEVERITY.WARNING]: {
    icon: AlertIcon,
    color: NOTIFICATIONCOLORS.WARNING,
    darkColor: NOTIFICATIONCOLORS.WARNING,
  },
  [SEVERITY.SUCCESS]: {
    icon: InfoIcon,
    color: NOTIFICATIONCOLORS.SUCCESS,
    darkColor: NOTIFICATIONCOLORS.SUCCESS,
  },
};

//TODO: This should be generated from OPENAPI schema
const EVENT_SCHEMA = {
  type: 'object',
  properties: {
    id: { type: 'string' },
    description: {
      type: 'string',
      default: '',
    },
    severity: {
      type: 'string',
      enum: Object.values(SEVERITY),
      default: SEVERITY.INFO,
    },
    status: {
      type: 'string',
      enum: Object.values(STATUS),
      default: STATUS.UNREAD,
    },
    created_at: { type: 'string' },
    updated_at: { type: 'string' },
    user_id: { type: 'string' },
    system_id: { type: 'string' },
    operation_id: { type: 'string' },
    action: { type: 'string' },
    category: { type: 'string' },
    metadata: {
      type: ['object', 'null'],
    },
  },
  required: [
    'id',
    'severity',
    'status',
    'created_at',
    'updated_at',
    'user_id',
    'system_id',
    'action',
  ],
};

// Validate event against EVENT_SCHEMA and return [isValid,validatedEvent]
export const validateEvent = (event) => {
  const eventCopy = _.cloneDeep(event) || {};
  eventCopy.status = eventCopy.status.trim() || STATUS.UNREAD;
  eventCopy.severity = eventCopy.severity.trim() || SEVERITY.INFO;
  const ajv = new Ajv({
    useDefaults: true,
  });
  const validate = ajv.compile(EVENT_SCHEMA);
  const valid = validate(eventCopy);
  return [valid, eventCopy];
};

// return validated events (adds default values if not present)
export const validateEvents = (events) => {
  return events
    .map((event) => {
      const [isValid, validatedEvent] = validateEvent(event);
      return isValid ? validatedEvent : null;
    })
    .filter((event) => event);
};

const EVENT_METADATA_SCHEMA = {
  type: 'object',
  properties: {
    error: {
      type: 'object',
      properties: {
        Code: { type: 'string' },
        LongDescription: { type: 'array', items: { type: 'string' }, default: [] },
        ProbableCause: { type: 'array', items: { type: 'string' }, default: [] },
        Severity: { type: 'number', default: 1 },
        ShortDescription: { type: 'array', items: { type: 'string' }, default: [] },
        SuggestedRemediation: { type: 'array', items: { type: 'string' }, default: [] },
      },
      required: [
        'Code',
        'LongDescription',
        'ProbableCause',
        'Severity',
        'ShortDescription',
        'SuggestedRemediation',
      ],
    },
  },
  required: ['error'],
};

export const validateEventMetadata = (metadata) => {
  const metadataCopy = _.cloneDeep(metadata) || {};
  const ajv = new Ajv();
  const validate = ajv.compile(EVENT_METADATA_SCHEMA);
  const valid = validate(metadataCopy);
  return [valid, metadataCopy];
};
