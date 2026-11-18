import { test as base, expect } from '@playwright/test';
import fs from 'fs';
import os from 'os';
import * as yaml from 'js-yaml';

/* eslint-disable react-hooks/rules-of-hooks */

export const test = base.extend({
  // Define an option and provide a default value.
  // We can later override it in the config.
  provider: ['Local', { option: true }],

  // The built-in Playwright `request` fixture creates a bare APIRequestContext
  // without browser Origin metadata. Local-provider CI pins Local, so unsafe
  // API calls made directly by tests must identify the local server origin,
  // just as same-origin browser requests do.
  request: async ({ playwright, baseURL }, use) => {
    const headers = baseURL ? { Origin: new URL(baseURL).origin } : undefined;
    const request = await playwright.request.newContext({ baseURL, extraHTTPHeaders: headers });
    await use(request);
    await request.dispose();
  },

  clusterMetaData: async ({}, use) => {
    const kubeConfigPath = `${os.homedir()}/.kube/config`;
    const kubeConfigRaw = fs.readFileSync(kubeConfigPath, 'utf8');
    const kubeConfig = yaml.load(kubeConfigRaw);

    const currentContextName = kubeConfig['current-context'];
    const context = kubeConfig.contexts.find((ctx) => ctx.name === currentContextName);
    const clusterName = context?.context?.cluster;
    const clusterEntry = kubeConfig.clusters.find((c) => c.name === clusterName);

    const clusterMetaData = {
      name: clusterName || 'kind-kind-cluster',
      kind: clusterEntry?.cluster?.['kind'] || 'Kubernetes',
    };

    await use(clusterMetaData);
  },
});

export { expect };
