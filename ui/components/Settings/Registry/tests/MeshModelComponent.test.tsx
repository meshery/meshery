import { render, screen } from '@testing-library/react';
import MeshModelComponent from '../MeshModelComponent';
import * as canModule from '@/utils/can';

// Mock CAN utility
jest.mock('@/utils/can');

// Mock permission constants
jest.mock('@/utils/permission_constants', () => ({
  keys: {
    CREATE_MODEL: { action: 'create-model-action', subject: 'Create Model' },
    IMPORT_MODEL: { action: 'import-model-action', subject: 'Import Model' },
  },
}));

// Mock next/router
jest.mock('next/router', () => ({
  useRouter: jest.fn(() => ({
    query: {},
    push: jest.fn(),
  })),
}));

// Mock RTK-Query hooks
jest.mock('@/rtk-query/meshModel', () => ({
  useLazyGetMeshModelsQuery: jest.fn(() => [jest.fn(), { data: null }]),
  useLazyGetComponentsQuery: jest.fn(() => [jest.fn(), { data: null }]),
  useLazyGetRelationshipsQuery: jest.fn(() => [jest.fn(), { data: null }]),
  useLazyGetRegistrantsQuery: jest.fn(() => [jest.fn(), { data: null }]),
}));

// Mock NoSsr to render children immediately in tests
jest.mock('@sistent/sistent', () => ({
  ...jest.requireActual('@sistent/sistent'),
  NoSsr: ({ children }) => <>{children}</>,
}));

describe('MeshModelComponent permission gating', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('disables Create Model button when user lacks permission', () => {
    (canModule.CAN as jest.Mock).mockImplementation((action) => {
      if (action === 'create-model-action') return false;
      return true;
    });

    render(<MeshModelComponent />);
    expect(
      screen.getByRole('button', { name: /create model/i })
    ).toBeDisabled();
  });

  it('enables Create Model button when user has permission', () => {
    (canModule.CAN as jest.Mock).mockReturnValue(true);

    render(<MeshModelComponent />);
    expect(
      screen.getByRole('button', { name: /create model/i })
    ).not.toBeDisabled();
  });

  it('disables Import Model button when user lacks permission', () => {
    (canModule.CAN as jest.Mock).mockImplementation((action) => {
      if (action === 'import-model-action') return false;
      return true;
    });

    render(<MeshModelComponent />);
    expect(
      screen.getByRole('button', { name: /import model/i })
    ).toBeDisabled();
  });

  it('enables Import Model button when user has permission', () => {
    (canModule.CAN as jest.Mock).mockReturnValue(true);

    render(<MeshModelComponent />);
    expect(
      screen.getByRole('button', { name: /import model/i })
    ).not.toBeDisabled();
  });
});