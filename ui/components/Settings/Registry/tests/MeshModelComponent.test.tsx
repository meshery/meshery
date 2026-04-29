/* eslint-env jest */
import { render, screen } from '@testing-library/react';
import MeshModelComponent from '../MeshModelComponent';
import * as canModule from '@/utils/can';

jest.mock('@/utils/can');
jest.mock('@/utils/permission_constants', () => ({
  keys: {
    CREATE_MODEL: { action: 'create-model-action', subject: 'Create Model' },
    IMPORT_MODEL: { action: 'import-model-action', subject: 'Import Model' },
  },
}));
jest.mock('next/router', () => ({
  useRouter: jest.fn(() => ({ query: {}, push: jest.fn() })),
}));
jest.mock('@/rtk-query/meshModel', () => ({
  useLazyGetMeshModelsQuery: jest.fn(() => [jest.fn(), { data: null }]),
  useLazyGetComponentsQuery: jest.fn(() => [jest.fn(), { data: null }]),
  useLazyGetRelationshipsQuery: jest.fn(() => [jest.fn(), { data: null }]),
  useLazyGetRegistrantsQuery: jest.fn(() => [jest.fn(), { data: null }]),
}));
jest.mock('@sistent/sistent', () => ({
  ...jest.requireActual('@sistent/sistent'),
  NoSsr: ({ children }) => <>{children}</>,
}));

describe('MeshModelComponent permission gating', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('disables Create Model button when user lacks permission', () => {
    (canModule.CAN as jest.Mock).mockImplementation((action) => action !== 'create-model-action');
    render(<MeshModelComponent />);
    expect(screen.getByRole('button', { name: /create model/i })).toBeDisabled();
  });

  it('enables Create Model button when user has permission', () => {
    (canModule.CAN as jest.Mock).mockReturnValue(true);
    render(<MeshModelComponent />);
    expect(screen.getByRole('button', { name: /create model/i })).not.toBeDisabled();
  });

  it('disables Import Model button when user lacks permission', () => {
    (canModule.CAN as jest.Mock).mockImplementation((action) => action !== 'import-model-action');
    render(<MeshModelComponent />);
    expect(screen.getByRole('button', { name: /import model/i })).toBeDisabled();
  });

  it('enables Import Model button when user has permission', () => {
    (canModule.CAN as jest.Mock).mockReturnValue(true);
    render(<MeshModelComponent />);
    expect(screen.getByRole('button', { name: /import model/i })).not.toBeDisabled();
  });
});
