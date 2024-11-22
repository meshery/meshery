import { test as base, expect } from '@playwright/test';

export const test = base.extend({
  // Define an option and provide a default value.
  // We can later override it in the config.
  provider: ['None', { option: true }],
});

export { expect };
