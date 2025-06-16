import { test as base, expect } from '@playwright/test';
import fs from 'fs';
import os from 'os';
import yaml from 'js-yaml';

export const test = base.extend({
  // Define an option and provide a default value.
  // We can later override it in the config.
  provider: ['None', { option: true }],
  // eslint-disable-next-line no-empty-pattern
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