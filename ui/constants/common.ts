export const ErrorTypes = {
  MESHERY_DEPLOYMENT_INCOMPATIBLE: 'MESHERY_DEPLOYMENT_INCOMPATIBLE',
  PAGE_NOT_FOUND: 'PAGE_NOT_FOUND',
  UNKNOWN: 'UNKNOWN',
} as const;

export type ErrorType = (typeof ErrorTypes)[keyof typeof ErrorTypes];

export const KUBERNETES = 'kubernetes';
export const CUSTOM_RESOURCE_DEFINITION = 'customresourcedefinition';

export const FALLBACK_MESHERY_IMAGE_PATH = '/static/img/meshery-logo/meshery-logo.svg';
export const FALLBACK_KUBERNETES_IMAGE_PATH = '/static/img/kubernetes.svg';
