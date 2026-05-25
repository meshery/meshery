import { describe, expect, it, vi } from 'vitest';

vi.mock('./configuration/config', () => ({ ConfigurationTableConfig: vi.fn() }));
vi.mock('./network/config', () => ({ NetWorkTableConfig: vi.fn() }));
vi.mock('./security/config', () => ({ SecurityTypesConfig: vi.fn() }));
vi.mock('./storage/config', () => ({ StorageTableConfig: vi.fn() }));
vi.mock('./workloads/config', () => ({ WorkloadTableConfig: vi.fn() }));
vi.mock('./namespace/config', () => ({ NamespaceTableConfig: vi.fn() }));
vi.mock('./nodes/config', () => ({ NodeTableConfig: vi.fn() }));
vi.mock('./crds/config', () => ({ CustomResourceConfig: vi.fn() }));

import {
  ALL_VIEW,
  ResourceMenuConfig,
  ResourcesConfig,
  SINGLE_VIEW,
  generateDynamicURL,
  getAllCustomResourceDefinitionsKinds,
} from './config';

describe('resources/config constants', () => {
  it('exports the documented ALL_VIEW and SINGLE_VIEW labels', () => {
    expect(ALL_VIEW).toBe('all');
    expect(SINGLE_VIEW).toBe('single');
  });

  it('marks Node/Namespace categories as having no submenu', () => {
    expect(ResourcesConfig.Node.submenu).toBe(false);
    expect(ResourcesConfig.Namespace.submenu).toBe(false);
  });

  it('marks the rest of the categories as having a submenu', () => {
    expect(ResourcesConfig.Workload.submenu).toBe(true);
    expect(ResourcesConfig.Configuration.submenu).toBe(true);
    expect(ResourcesConfig.Network.submenu).toBe(true);
    expect(ResourcesConfig.Security.submenu).toBe(true);
    expect(ResourcesConfig.Storage.submenu).toBe(true);
    expect(ResourcesConfig.CRDS.submenu).toBe(true);
  });

  it('publishes the expected workload kinds in ResourceMenuConfig', () => {
    expect(ResourceMenuConfig.Workload).toContain('Pod');
    expect(ResourceMenuConfig.Workload).toContain('Deployment');
    expect(ResourceMenuConfig.Workload).toContain('CronJob');
  });
});

describe('generateDynamicURL', () => {
  it('returns the category for top-level entries like Node and Namespace', () => {
    expect(generateDynamicURL('Node')).toBe('?resourceCategory=Node&resource=');
    expect(generateDynamicURL('Namespace')).toBe('?resourceCategory=Namespace&resource=');
  });

  it('locates a specific kind within its category', () => {
    expect(generateDynamicURL('Pod')).toBe('?resourceCategory=Workload&resource=Pod');
    expect(generateDynamicURL('Ingress')).toBe('?resourceCategory=Network&resource=Ingress');
    expect(generateDynamicURL('Secret')).toBe('?resourceCategory=Configuration&resource=Secret');
  });

  it('falls back to the CRDS category for unknown kinds', () => {
    expect(generateDynamicURL('SomeNewKind')).toBe('?resourceCategory=CRDS&resource=SomeNewKind');
  });
});

describe('getAllCustomResourceDefinitionsKinds', () => {
  it('filters out kinds that are in the known resource menu', () => {
    const kinds = [
      { Kind: 'Pod' },
      { Kind: 'Service' },
      { Kind: 'NotebookSpec' },
      { Kind: 'Node' },
      { Kind: 'Namespace' },
      { Kind: 'AcmeRandom' },
    ];
    const crds = getAllCustomResourceDefinitionsKinds(kinds);
    expect(crds.map((k) => k.Kind)).toEqual(['NotebookSpec', 'AcmeRandom']);
  });

  it('returns [] when the input is null/undefined', () => {
    expect(getAllCustomResourceDefinitionsKinds(undefined)).toEqual([]);
    expect(getAllCustomResourceDefinitionsKinds(null)).toEqual([]);
  });

  it('returns [] when no kinds qualify as custom resources', () => {
    const kinds = [{ Kind: 'Pod' }, { Kind: 'Node' }];
    expect(getAllCustomResourceDefinitionsKinds(kinds)).toEqual([]);
  });
});
