// Needed for fetching server-side code coverage. See: https://github.com/bahmutov/next-and-cypress-example

import { NextApiHandler } from 'next';

let handler: NextApiHandler | undefined;

if (process.env.NODE_ENV === 'development') {
  handler = require('@cypress/code-coverage/middleware/nextjs');
}

export default handler;
