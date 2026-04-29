import React from 'react';
import { render, screen } from '@testing-library/react';
import MeshModelComponent from '../MeshModelComponent';
import * as canModule from '@/utils/can';

// Mock the required modules
jest.mock('@/utils/can');
jest.mock('@/utils/permission_constants', () => ({
  keys: {
    CREATE_MODEL: { action: 'create-model-action', subject: 'Create Model' },
    IMPORT_MODEL: { action: 'import-model-action', subject: 'Import Model' },
  },
}));

jest.mock('next/router', () => ({
  useRouter: jest.fn(() => ({
    query: {},
    push: jest.fn(),
    asPath: '/',
    pathname: '/',
  })),
}));

jest.mock('@/rtk-query/meshModel', () => ({
  useLazyGetMeshModelsQuery: jest.fn(() => [jest.fn(), { data: null, isLoading: false, isFetching: false }]),
  useLazyGetComponentsQuery: jest.fn(() => [jest.fn(), { data: null, isLoading: false, isFetching: false }]),
  useLazyGetRelationshipsQuery: jest.fn(() => [jest.fn(), { data: null, isLoading: false, isFetching: false }]),
  useLazyGetRegistrantsQuery: jest.fn(() => [jest.fn(), { data: null, isLoading: false, isFetching: false }]),
}));

jest.mock('@sistent/sistent', () => {
  const actual = jest.requireActual('@sistent/sistent');
  return {
    ...actual,
    NoSsr: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  };
});

describe('MeshModelComponent permission gating', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('disables Create Model button when user lacks create_model permission', () => {
    const mockCAN = jest.spyOn(canModule, 'default' as any).mockImplementation(
      (action, subject) => {
        if (action === 'create-model-action' && subject === 'Create Model') return false;
        return true;
      },
    );

    render(<MeshModelComponent />);
    const createButton = screen.getByRole('button', { name: /create model/i });
    expect(createButton).toBeDisabled();

    mockCAN.mockRestore();
  });

  it('enables Create Model button when user has create_model permission', () => {
    const mockCAN = jest.spyOn(canModule, 'default' as any).mockReturnValue(true);

    render(<MeshModelComponent />);
    const createButton = screen.getByRole('button', { name: /create model/i });
    expect(createButton).not.toBeDisabled();

    mockCAN.mockRestore();
  });

  it('disables Import Model button when user lacks import_model permission', () => {
    const mockCAN = jest.spyOn(canModule, 'default' as any).mockImplementation(
      (action, subject) => {
        if (action === 'import-model-action' && subject === 'Import Model') return false;
        return true;
      },
    );

    render(<MeshModelComponent />);
    const importButton = screen.getByRole('button', { name: /import model/i });
    expect(importButton).toBeDisabled();

    mockCAN.mockRestore();
  });

  it('enables Import Model button when user has import_model permission', () => {
    const mockCAN = jest.spyOn(canModule, 'default' as any).mockReturnValue(true);

    render(<MeshModelComponent />);
    const importButton = screen.getByRole('button', { name: /import model/i });
    expect(importButton).not.toBeDisabled();

    mockCAN.mockRestore();
  });

  it('enables both buttons when user has all permissions', () => {
    const mockCAN = jest.spyOn(canModule, 'default' as any).mockReturnValue(true);

    render(<MeshModelComponent />);
    const createButton = screen.getByRole('button', { name: /create model/i });
    const importButton = screen.getByRole('button', { name: /import model/i });

    expect(createButton).not.toBeDisabled();
    expect(importButton).not.toBeDisabled();

    mockCAN.mockRestore();
  });

  it('disables both buttons when user lacks both permissions', () => {
    const mockCAN = jest.spyOn(canModule, 'default' as any).mockReturnValue(false);

    render(<MeshModelComponent />);
    const createButton = screen.getByRole('button', { name: /create model/i });
    const importButton = screen.getByRole('button', { name: /import model/i });

    expect(createButton).toBeDisabled();
    expect(importButton).toBeDisabled();

    mockCAN.mockRestore();
  });
});
