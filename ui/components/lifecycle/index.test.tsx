import { describe, expect, it, vi } from 'vitest';

vi.mock('../workspaces', () => ({
  default: () => null,
}));

vi.mock('../environments', () => ({
  default: () => null,
}));

import { WorkspacesComponent, EnvironmentComponent } from './index';

describe('lifecycle/index', () => {
  it('re-exports WorkspacesComponent and EnvironmentComponent', () => {
    expect(WorkspacesComponent).toBeDefined();
    expect(EnvironmentComponent).toBeDefined();
  });
});
