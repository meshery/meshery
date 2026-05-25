import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

// ---------------------------------------------------------------------------
// environments.ts is a thin façade over the schemas-package hooks. Each
// exported hook wraps a `useSchemas…Mutation` and normalises argument shape
// (stringifying pagination, accepting multiple field-name aliases for the
// organization ID, etc.). We mock every underlying hook and assert wiring.
// ---------------------------------------------------------------------------

const { triggers, schemasGetEnvs, schemasGetEnvConns, result } = vi.hoisted(() => ({
  result: { isLoading: false, data: { ok: true } },
  schemasGetEnvs: vi.fn(),
  schemasGetEnvConns: vi.fn(),
  triggers: {
    create: vi.fn(),
    update: vi.fn(),
    del: vi.fn(),
    add: vi.fn(),
    remove: vi.fn(),
  },
}));

vi.mock('@meshery/schemas/mesheryApi', () => ({
  useGetEnvironmentsQuery: (...args: unknown[]) => schemasGetEnvs(...args),
  useGetEnvironmentConnectionsQuery: (...args: unknown[]) => schemasGetEnvConns(...args),
  useCreateEnvironmentMutation: () => [triggers.create, result],
  useUpdateEnvironmentMutation: () => [triggers.update, result],
  useDeleteEnvironmentMutation: () => [triggers.del, result],
  useAddConnectionToEnvironmentMutation: () => [triggers.add, result],
  useRemoveConnectionFromEnvironmentMutation: () => [triggers.remove, result],
}));

import {
  useGetEnvironmentsQuery,
  useGetEnvironmentConnectionsQuery,
  useCreateEnvironmentMutation,
  useSaveEnvironmentMutation,
  useUpdateEnvironmentMutation,
  useDeleteEnvironmentMutation,
  useAddConnectionToEnvironmentMutation,
  useRemoveConnectionFromEnvironmentMutation,
} from '../environments';

describe('useGetEnvironmentsQuery wrapper', () => {
  beforeEach(() => {
    schemasGetEnvs.mockReset();
    schemasGetEnvs.mockReturnValue({ data: { environments: [] } });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('stringifies page and pagesize and passes other params verbatim', () => {
    useGetEnvironmentsQuery({
      page: 2,
      pagesize: 50,
      search: 'staging',
      order: 'desc',
      orgId: 'o',
    });
    expect(schemasGetEnvs).toHaveBeenCalledWith(
      {
        search: 'staging',
        order: 'desc',
        page: '2',
        pagesize: '50',
        orgId: 'o',
      },
      undefined,
    );
  });

  it('handles missing pagination fields', () => {
    useGetEnvironmentsQuery({ search: 'q' });
    expect(schemasGetEnvs).toHaveBeenCalledWith(
      expect.objectContaining({ search: 'q', page: undefined, pagesize: undefined }),
      undefined,
    );
  });

  it('returns the underlying hook result', () => {
    const expected = { data: { environments: [{ id: '1' }] } };
    schemasGetEnvs.mockReturnValue(expected);
    expect(useGetEnvironmentsQuery({})).toBe(expected);
  });
});

describe('useGetEnvironmentConnectionsQuery wrapper', () => {
  beforeEach(() => {
    schemasGetEnvConns.mockReset();
    schemasGetEnvConns.mockReturnValue({ data: { connections: [] } });
  });

  it('forwards environmentId, search, order, filter and stringified pagination', () => {
    useGetEnvironmentConnectionsQuery({
      environmentId: 'env-1',
      search: 's',
      order: 'asc',
      page: 1,
      pagesize: 10,
      filter: 'kubernetes',
    });

    expect(schemasGetEnvConns).toHaveBeenCalledWith(
      {
        environmentId: 'env-1',
        search: 's',
        order: 'asc',
        page: '1',
        pagesize: '10',
        filter: 'kubernetes',
      },
      undefined,
    );
  });
});

describe('useCreateEnvironmentMutation wrapper', () => {
  beforeEach(() => {
    Object.values(triggers).forEach((t) => t.mockReset());
  });

  it('returns [trigger, result] tuple', () => {
    const [trigger, res] = useCreateEnvironmentMutation();
    expect(typeof trigger).toBe('function');
    expect(res).toEqual(result);
  });

  it('reshapes payload to {body: {name,description,organizationId}}', () => {
    const [trigger] = useCreateEnvironmentMutation();
    trigger({
      environmentPayload: {
        name: 'prod',
        description: 'production',
        organizationId: 'org-1',
      },
    });
    expect(triggers.create).toHaveBeenCalledWith({
      body: {
        name: 'prod',
        description: 'production',
        organizationId: 'org-1',
      },
    });
  });

  it('accepts organization_id as a fallback for organizationId', () => {
    const [trigger] = useCreateEnvironmentMutation();
    trigger({
      environmentPayload: { name: 'staging', description: 'd', organization_id: 'org-2' },
    });
    expect(triggers.create.mock.calls[0][0].body.organizationId).toBe('org-2');
  });

  it('accepts OrganizationID (PascalCase) as a final fallback', () => {
    const [trigger] = useCreateEnvironmentMutation();
    trigger({
      environmentPayload: { name: 'dev', description: 'd', OrganizationID: 'org-3' },
    });
    expect(triggers.create.mock.calls[0][0].body.organizationId).toBe('org-3');
  });
});

describe('useSaveEnvironmentMutation wrapper', () => {
  beforeEach(() => {
    Object.values(triggers).forEach((t) => t.mockReset());
  });

  it('forwards the raw body unchanged', () => {
    const [trigger] = useSaveEnvironmentMutation();
    const body = { name: 'foo', custom_field: 'bar' };
    trigger({ body });
    expect(triggers.create).toHaveBeenCalledWith({ body });
  });
});

describe('useUpdateEnvironmentMutation wrapper', () => {
  beforeEach(() => {
    Object.values(triggers).forEach((t) => t.mockReset());
  });

  it('forwards environmentId plus reshaped body', () => {
    const [trigger] = useUpdateEnvironmentMutation();
    trigger({
      environmentId: 'env-7',
      environmentPayload: { name: 'rename', description: 'desc', organizationId: 'org' },
    });
    expect(triggers.update).toHaveBeenCalledWith({
      environmentId: 'env-7',
      body: { name: 'rename', description: 'desc', organizationId: 'org' },
    });
  });
});

describe('useDeleteEnvironmentMutation wrapper', () => {
  beforeEach(() => {
    Object.values(triggers).forEach((t) => t.mockReset());
  });

  it('passes the environmentId through to the underlying mutation', () => {
    const [trigger] = useDeleteEnvironmentMutation();
    trigger({ environmentId: 'env-9' });
    expect(triggers.del).toHaveBeenCalledWith({ environmentId: 'env-9' });
  });
});

describe('useAddConnectionToEnvironmentMutation wrapper', () => {
  beforeEach(() => {
    Object.values(triggers).forEach((t) => t.mockReset());
  });

  it('passes environmentId/connectionId pair to the underlying mutation', () => {
    const [trigger] = useAddConnectionToEnvironmentMutation();
    trigger({ environmentId: 'env-1', connectionId: 'conn-1' });
    expect(triggers.add).toHaveBeenCalledWith({
      environmentId: 'env-1',
      connectionId: 'conn-1',
    });
  });
});

describe('useRemoveConnectionFromEnvironmentMutation wrapper', () => {
  beforeEach(() => {
    Object.values(triggers).forEach((t) => t.mockReset());
  });

  it('passes environmentId/connectionId pair to the underlying mutation', () => {
    const [trigger] = useRemoveConnectionFromEnvironmentMutation();
    trigger({ environmentId: 'env-1', connectionId: 'conn-2' });
    expect(triggers.remove).toHaveBeenCalledWith({
      environmentId: 'env-1',
      connectionId: 'conn-2',
    });
  });
});
