const MESHERY_SERVER_URL = process.env.MESHERY_SERVER_URL || 'http://localhost:9081';
const REMOTE_PROVIDER_URL = process.env.REMOTE_PROVIDER_URL || 'http://localhost:9876';
const REMOTE_PROVIDER_USER = {
  email: process.env.REMOTE_PROVIDER_USER_EMAIL || 'test-admin@layer5.io',
  password: process.env.REMOTE_PROVIDER_USER_PASSWORD || 'test-admin',
};
const PROVIDER_SELECTION_URL = `${MESHERY_SERVER_URL}/provider`;
const PROVIDER_TOKEN = process.env.PROVIDER_TOKEN;

export const ENV = {
  MESHERY_SERVER_URL,
  PROVIDER_SELECTION_URL,
  REMOTE_PROVIDER_URL,
  REMOTE_PROVIDER_USER,
  PROVIDER_TOKEN,
};
