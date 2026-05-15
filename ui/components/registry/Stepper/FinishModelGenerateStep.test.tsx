import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const importMeshModel = vi.fn();
let mutationReturn: any = [importMeshModel, { isLoading: false, error: null }];

vi.mock('lodash', () => ({
  capitalize: (s: string) => s.charAt(0).toUpperCase() + s.slice(1),
}));

vi.mock('@/components/designs/lifecycle/common', () => ({
  Loading: ({ message }: any) => <div data-testid="loading">{message}</div>,
}));

vi.mock('@/components/layout/NotificationCenter', () => {
  const ctx = (require('react') as typeof import('react')).createContext({
    operationsCenterActorRef: {
      on: (_event: any, _cb: any) => ({ unsubscribe: () => {} }),
    },
  });
  return { NotificationCenterContext: ctx };
});

vi.mock('machines/operationsCenter', () => ({
  OPERATION_CENTER_EVENTS: {
    EVENT_RECEIVED_FROM_SERVER: 'event-received-from-server',
  },
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useImportMeshModelMutation: () => mutationReturn,
}));

vi.mock('@/components/layout/NotificationCenter/formatters/model_registration', () => ({
  ModelImportedSection: ({ modelDetails }: any) => (
    <div data-testid="imported-section">{JSON.stringify(modelDetails || null)}</div>
  ),
  ModelImportMessages: ({ message }: any) => (
    <div data-testid="import-messages">{JSON.stringify(message || null)}</div>
  ),
}));

vi.mock('@/components/layout/NotificationCenter/formatters/error', () => ({
  ErrorMetadataFormatter: ({ metadata }: any) => (
    <div data-testid="error-metadata">{JSON.stringify(metadata || null)}</div>
  ),
}));

import FinishModelGenerateStep from './FinishModelGenerateStep';

describe('FinishModelGenerateStep', () => {
  beforeEach(() => {
    importMeshModel.mockReset();
    mutationReturn = [importMeshModel, { isLoading: false, error: null }];
  });

  it('renders a Loading message capitalized while loading', () => {
    mutationReturn = [importMeshModel, { isLoading: true, error: null }];

    render(<FinishModelGenerateStep requestBody={{}} generateType="generate" />);

    expect(screen.getByTestId('loading')).toHaveTextContent('Generateing model');
  });

  it('renders the error formatter when there is an error', () => {
    mutationReturn = [importMeshModel, { isLoading: false, error: { msg: 'boom' } }];

    render(<FinishModelGenerateStep requestBody={{}} generateType="generate" />);

    expect(screen.getByTestId('error-metadata')).toBeInTheDocument();
  });

  it('invokes importMeshModel on mount with the request body', () => {
    render(<FinishModelGenerateStep requestBody={{ foo: 'bar' }} generateType="generate" />);

    expect(importMeshModel).toHaveBeenCalledWith({ importBody: { foo: 'bar' } });
  });

  it('renders import messages and imported section in the success state', () => {
    render(<FinishModelGenerateStep requestBody={{}} generateType="generate" />);

    expect(screen.getByTestId('import-messages')).toBeInTheDocument();
    expect(screen.getByTestId('imported-section')).toBeInTheDocument();
  });
});
