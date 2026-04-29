import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import * as canModule from '@/utils/can';

vi.mock('@/utils/can');
vi.mock('@/utils/permission_constants', () => ({
  keys: {
    CREATE_MODEL: { action: 'create-model-action', subject: 'Create Model' },
    IMPORT_MODEL: { action: 'import-model-action', subject: 'Import Model' },
  },
}));
vi.mock('next/router', () => ({
  useRouter: vi.fn(() => ({ query: {}, push: vi.fn() })),
}));
vi.mock('@/rtk-query/meshModel', () => ({
  useLazyGetMeshModelsQuery: vi.fn(() => [vi.fn(), { data: null }]),
  useLazyGetComponentsQuery: vi.fn(() => [vi.fn(), { data: null }]),
  useLazyGetRelationshipsQuery: vi.fn(() => [vi.fn(), { data: null }]),
  useLazyGetRegistrantsQuery: vi.fn(() => [vi.fn(), { data: null }]),
}));
vi.mock('@sistent/sistent', async () => {
  const actual = await vi.importActual('@sistent/sistent');
  return { ...actual, NoSsr: ({ children }) => <>{children}</> };
});

import MeshModelComponent from '../MeshModelComponent';

describe('MeshModelComponent permission gating', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('disables Create Model button when user lacks permission', () => {
    (canModule.CAN as ReturnType<typeof vi.fn>).mockImplementation((action) => action !== 'create-model-action');
    render(<MeshModelComponent />);
    expect(screen.getByRole('button', { name: /create model/i })).toBeDisabled();
  });

  it('enables Create Model button when user has permission', () => {
    (canModule.CAN as ReturnType<typeof vi.fn>).mockReturnValue(true);
    render(<MeshModelComponent />);
    expect(screen.getByRole('button', { name: /create model/i })).not.toBeDisabled();
  });

  it('disables Import Model button when user lacks permission', () => {
    (canModule.CAN as ReturnType<typeof vi.fn>).mockImplementation((action) => action !== 'import-model-action');
    render(<MeshModelComponent />);
    expect(screen.getByRole('button', { name: /import model/i })).toBeDisabled();
  });

  it('enables Import Model button when user has permission', () => {
    (canModule.CAN as ReturnType<typeof vi.fn>).mockReturnValue(true);
    render(<MeshModelComponent />);
    expect(screen.getByRole('button', { name: /import model/i })).not.toBeDisabled();
  });
});
