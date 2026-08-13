import path from 'path';
import dotenv from 'dotenv';

// ui/.env has to be loaded before the captures below, and this module is
// imported directly by the setup projects and the specs as well as by
// playwright.config.js, so it owns the load. The path is resolved against this
// file rather than the working directory, a missing .env is normal, and real
// environment variables win - dotenv does not override them by default.
if (!process.env.CI) {
  dotenv.config({ path: path.resolve(__dirname, '../../.env') });
}

const MESHERY_SERVER_URL = process.env.MESHERY_SERVER_URL || 'http://localhost:9081';
const REMOTE_PROVIDER_URL = process.env.REMOTE_PROVIDER_URL || 'https://cloud.meshery.io'; // AUTO-GENERATED from install/providers.env - run `make providers-propagate`

const USER_EMAIL = process.env.REMOTE_PROVIDER_USER_EMAIL;
const USER_PASSWORD = process.env.REMOTE_PROVIDER_USER_PASSWORD;

const REMOTE_PROVIDER_USER = {
  email: USER_EMAIL,
  password: USER_PASSWORD,
};
const PROVIDER_SELECTION_URL = `${MESHERY_SERVER_URL}/provider`;
const PROVIDER_TOKEN = process.env.PROVIDER_TOKEN;

if (!PROVIDER_TOKEN) {
  if (!USER_EMAIL && !USER_PASSWORD) {
    console.warn(
      'No remote-provider credentials configured. Set PROVIDER_TOKEN, or both ' +
        'REMOTE_PROVIDER_USER_EMAIL and REMOTE_PROVIDER_USER_PASSWORD. Without one of those the ' +
        'chromium-meshery-provider project fails its setup; the Local provider project is unaffected.',
    );
  } else if (!USER_EMAIL || !USER_PASSWORD) {
    console.warn(
      `Incomplete remote-provider credentials: ${
        USER_EMAIL ? 'REMOTE_PROVIDER_USER_PASSWORD' : 'REMOTE_PROVIDER_USER_EMAIL'
      } is not set. ` +
        'Set both, or set PROVIDER_TOKEN instead. The chromium-meshery-provider project fails ' +
        'its setup without them; the Local provider project is unaffected.',
    );
  }
}

const AUTHFILELOCALPROVIDER = 'playwright/.auth/user-local-provider.json';

const AUTHFILEMESHERYPROVIDER = 'playwright/.auth/user-meshery-provider.json';

export const ENV = {
  MESHERY_SERVER_URL,
  PROVIDER_SELECTION_URL,
  REMOTE_PROVIDER_URL,
  REMOTE_PROVIDER_USER,
  PROVIDER_TOKEN,
  AUTHFILELOCALPROVIDER,
  AUTHFILEMESHERYPROVIDER,
};
