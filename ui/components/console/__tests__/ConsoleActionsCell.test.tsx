import React from 'react';
import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import ConsoleActionsCell, { consoleAddressFor } from '../ConsoleActionsCell';

const k8sConfig = [{ kubernetesServerId: 'cluster-1', connectionId: 'conn-1', name: 'local' }];

const row = (kind: string, name = 'meshery-meshsync') => ({
  id: `${kind}/${name}`,
  kind,
  cluster_id: 'cluster-1',
  metadata: { name, namespace: 'meshery' },
});

describe('consoleAddressFor', () => {
  it('addresses a workload row with the workload kind, not always pod', () => {
    // Hard-coding `resource: 'pod'` sent the deployment's name to pods/<name>,
    // which is why every workload row answered "Console target not found".
    expect(consoleAddressFor(row('Deployment'), k8sConfig)).toEqual({
      connectionId: 'conn-1',
      target: { resource: 'deployment', namespace: 'meshery', name: 'meshery-meshsync' },
    });
  });

  it('addresses a pod row as a pod', () => {
    expect(consoleAddressFor(row('Pod', 'meshery-server-7d9f'), k8sConfig)?.target.resource).toBe(
      'pod',
    );
  });

  it('names no address for a resource that runs no containers', () => {
    expect(consoleAddressFor(row('ConfigMap', 'settings'), k8sConfig)).toBeNull();
  });

  it('names no address when the row has no name or no connection', () => {
    expect(
      consoleAddressFor({ kind: 'Deployment', cluster_id: 'cluster-1', metadata: {} }, k8sConfig),
    ).toBeNull();
    expect(
      consoleAddressFor({ kind: 'Deployment', metadata: { name: 'x' } }, k8sConfig),
    ).toBeNull();
  });
});

describe('ConsoleActionsCell', () => {
  it('offers the console buttons for a workload row', () => {
    render(<ConsoleActionsCell resource={row('Deployment')} k8sConfig={k8sConfig} />);

    screen.getByRole('button', { name: 'View logs' });
    screen.getByRole('button', { name: 'Open a terminal' });
  });

  it('offers nothing for a resource that hosts no console', () => {
    render(<ConsoleActionsCell resource={row('ConfigMap', 'settings')} k8sConfig={k8sConfig} />);

    expect(screen.queryByRole('button', { name: 'View logs' })).toBeNull();
    expect(screen.queryByRole('button', { name: 'Open a terminal' })).toBeNull();
  });
});
