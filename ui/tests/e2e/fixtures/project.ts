import { test as base, expect } from '@playwright/test';
import fs from 'fs';
import os from 'os';
import yaml from 'js-yaml';

interface ClusterMetaData {
  name: string;
  kind: string;
}

interface TestOptions {
  provider: string;
  clusterMetaData: ClusterMetaData;
}

/* eslint-disable react-hooks/rules-of-hooks */
export const test = base.extend<TestOptions>({
  provider: ['None', { option: true }],
  clusterMetaData: async ({}, use) => {
    const kubeConfigPath = `${os.homedir()}/.kube/config`;
    const kubeConfigRaw = fs.readFileSync(kubeConfigPath, 'utf8');
    const kubeConfig = yaml.load(kubeConfigRaw) as Record<string, unknown>;
    const currentContextName = kubeConfig['current-context'] as string;
    const contexts = kubeConfig.contexts as Array<{ name: string; context: { cluster: string } }>;
    const context = contexts.find((ctx) => ctx.name === currentContextName);
    const clusterName = context?.context?.cluster;
    const clusters = kubeConfig.clusters as Array<{
      name: string;
      cluster: Record<string, string>;
    }>;
    const clusterEntry = clusters.find((c) => c.name === clusterName);
    const clusterMetaData: ClusterMetaData = {
      name: clusterName || 'kind-kind-cluster',
      kind: clusterEntry?.cluster?.['kind'] || 'Kubernetes',
    };
    await use(clusterMetaData);
  },
});

export { expect };
