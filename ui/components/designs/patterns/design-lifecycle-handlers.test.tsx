import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const parseDesignFileMock = vi.fn((fileContents: string) => ({ parsed: fileContents }));

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<Record<string, unknown>>();
  return {
    ...actual,
    ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
  };
});

vi.mock('@/assets/icons', () => ({
  DoneAll: () => <svg data-testid="done-all" />,
}));

vi.mock('../../../public/static/img/UndeployIcon', () => ({
  default: () => <svg data-testid="undeploy-icon" />,
}));

vi.mock('@/assets/icons/DryRunIcon', () => ({
  default: () => <svg data-testid="dry-run-icon" />,
}));

vi.mock('@/assets/icons/CheckIcon', () => ({
  default: () => <svg data-testid="check-icon" />,
}));

vi.mock('../lifecycle/DeployStepper', () => ({
  DeployStepper: ({ design }: any) => (
    <div data-testid="deploy-stepper" data-design={JSON.stringify(design)} />
  ),
  UnDeployStepper: ({ design }: any) => (
    <div data-testid="undeploy-stepper" data-design={JSON.stringify(design)} />
  ),
}));

vi.mock('../lifecycle/DryRun', () => ({
  DryRunDesign: ({ design }: any) => (
    <div data-testid="dry-run-design" data-design={JSON.stringify(design)} />
  ),
}));

vi.mock('../lifecycle/ValidateDesign', () => ({
  ValidateDesign: ({ design }: any) => (
    <div data-testid="validate-design" data-design={JSON.stringify(design)} />
  ),
}));

vi.mock('../lifecycle/common', () => ({
  DEPLOYMENT_TYPE: { DEPLOY: 'deploy', UNDEPLOY: 'undeploy' },
}));

vi.mock('../../../utils/utils', () => ({
  parseDesignFile: (file: string) => parseDesignFileMock(file),
}));

import { buildDesignLifecycleHandlers } from './design-lifecycle-handlers';

const makeDeps = () => {
  const designLifecycleModal = {
    openModal: vi.fn(),
    closeModal: vi.fn(),
  };
  return {
    designLifecycleModal,
    designValidationActorRef: { send: vi.fn() },
    selectedK8sContexts: ['ctx-a'],
    handleDeploy: vi.fn(),
    handleUndeploy: vi.fn(),
    notify: vi.fn(),
  };
};

describe('buildDesignLifecycleHandlers', () => {
  it('returns all modal and direct lifecycle handlers', () => {
    const handlers = buildDesignLifecycleHandlers(makeDeps());
    expect(typeof handlers.openDeployModal).toBe('function');
    expect(typeof handlers.openUndeployModal).toBe('function');
    expect(typeof handlers.openDryRunModal).toBe('function');
    expect(typeof handlers.openValidateModal).toBe('function');
    expect(typeof handlers.directDeploy).toBe('function');
    expect(typeof handlers.directUndeploy).toBe('function');
    expect(typeof handlers.directDryRun).toBe('function');
    expect(typeof handlers.directValidate).toBe('function');
  });

  it('openDeployModal stops propagation and opens a modal containing the DeployStepper', () => {
    const deps = makeDeps();
    const handlers = buildDesignLifecycleHandlers(deps);
    const stop = vi.fn();
    handlers.openDeployModal({ stopPropagation: stop } as any, 'yaml-bytes', 'My Design');

    expect(stop).toHaveBeenCalled();
    expect(deps.designLifecycleModal.openModal).toHaveBeenCalled();
    const args = deps.designLifecycleModal.openModal.mock.calls[0][0];
    expect(args.title).toBe('Deploy design "My Design"');
    expect(parseDesignFileMock).toHaveBeenCalledWith('yaml-bytes');

    render(args.reactNode);
    expect(screen.getByTestId('deploy-stepper')).toBeInTheDocument();
  });

  it('openUndeployModal opens a modal with an UnDeployStepper', () => {
    const deps = makeDeps();
    const handlers = buildDesignLifecycleHandlers(deps);
    handlers.openUndeployModal({ stopPropagation: vi.fn() } as any, 'yaml', 'X');

    const args = deps.designLifecycleModal.openModal.mock.calls[0][0];
    expect(args.title).toBe('Undeploy design "X"');
    render(args.reactNode);
    expect(screen.getByTestId('undeploy-stepper')).toBeInTheDocument();
  });

  it('openDryRunModal opens a modal with a DryRunDesign body', () => {
    const deps = makeDeps();
    const handlers = buildDesignLifecycleHandlers(deps);
    handlers.openDryRunModal({ stopPropagation: vi.fn() } as any, 'yaml', 'X');

    const args = deps.designLifecycleModal.openModal.mock.calls[0][0];
    expect(args.title).toBe('Dryrun design "X"');
    render(args.reactNode);
    expect(screen.getByTestId('dry-run-design')).toBeInTheDocument();
  });

  it('openValidateModal opens a modal with a ValidateDesign body', () => {
    const deps = makeDeps();
    const handlers = buildDesignLifecycleHandlers(deps);
    handlers.openValidateModal({ stopPropagation: vi.fn() } as any, 'yaml', 'X');

    const args = deps.designLifecycleModal.openModal.mock.calls[0][0];
    expect(args.title).toBe('Validate design "X"');
    render(args.reactNode);
    expect(screen.getByTestId('validate-design')).toBeInTheDocument();
  });

  it('directDeploy executes handleDeploy and notifies the user', async () => {
    const deps = makeDeps();
    const handlers = buildDesignLifecycleHandlers(deps);
    const stop = vi.fn();

    await handlers.directDeploy(
      { stopPropagation: stop } as any,
      'yaml-bytes',
      'My Design',
      'id-123',
    );

    expect(stop).toHaveBeenCalled();
    expect(deps.handleDeploy).toHaveBeenCalledWith({
      design: expect.objectContaining({ id: 'id-123', parsed: 'yaml-bytes' }),
      selectedK8sContexts: ['ctx-a'],
    });
    expect(deps.notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Deploying design "My Design"',
      }),
    );
  });

  it('directUndeploy executes handleUndeploy and notifies the user', async () => {
    const deps = makeDeps();
    const handlers = buildDesignLifecycleHandlers(deps);
    const stop = vi.fn();

    await handlers.directUndeploy(
      { stopPropagation: stop } as any,
      'yaml-bytes',
      'My Design',
      'id-123',
    );

    expect(stop).toHaveBeenCalled();
    expect(deps.handleUndeploy).toHaveBeenCalledWith({
      design: expect.objectContaining({ id: 'id-123', parsed: 'yaml-bytes' }),
      selectedK8sContexts: ['ctx-a'],
    });
    expect(deps.notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Undeploying design "My Design"',
      }),
    );
  });

  it('directDryRun dispatches dry run command to validation actor and notifies', () => {
    const deps = makeDeps();
    const handlers = buildDesignLifecycleHandlers(deps);
    const stop = vi.fn();

    handlers.directDryRun({ stopPropagation: stop } as any, 'yaml-bytes', 'My Design');

    expect(stop).toHaveBeenCalled();
    expect(deps.designValidationActorRef.send).toHaveBeenCalled();
    expect(deps.notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Running dry run for design "My Design"',
      }),
    );
  });

  it('directValidate dispatches validation schema command to validation actor and notifies', () => {
    const deps = makeDeps();
    const handlers = buildDesignLifecycleHandlers(deps);
    const stop = vi.fn();

    handlers.directValidate({ stopPropagation: stop } as any, 'yaml-bytes', 'My Design');

    expect(stop).toHaveBeenCalled();
    expect(deps.designValidationActorRef.send).toHaveBeenCalled();
    expect(deps.notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Validating design "My Design"',
      }),
    );
  });

  it('direct handlers emit error notifications and abort when design parsing fails', async () => {
    parseDesignFileMock.mockReturnValueOnce(null as any);
    const deps = makeDeps();
    const handlers = buildDesignLifecycleHandlers(deps);
    const expectedError = expect.objectContaining({
      message: 'Failed to parse design "Bad Design"',
      event_type: expect.objectContaining({ type: 'error' }),
    });

    await handlers.directDeploy(undefined, 'invalid-yaml', 'Bad Design', 'id-1');
    expect(deps.handleDeploy).not.toHaveBeenCalled();
    expect(deps.notify).toHaveBeenCalledTimes(1);
    expect(deps.notify).toHaveBeenCalledWith(expectedError);

    deps.notify.mockClear();
    parseDesignFileMock.mockReturnValueOnce(null as any);
    await handlers.directUndeploy(undefined, 'invalid-yaml', 'Bad Design', 'id-1');
    expect(deps.handleUndeploy).not.toHaveBeenCalled();
    expect(deps.notify).toHaveBeenCalledTimes(1);
    expect(deps.notify).toHaveBeenCalledWith(expectedError);

    deps.notify.mockClear();
    parseDesignFileMock.mockReturnValueOnce(null as any);
    handlers.directDryRun(undefined, 'invalid-yaml', 'Bad Design');
    expect(deps.designValidationActorRef.send).not.toHaveBeenCalled();
    expect(deps.notify).toHaveBeenCalledTimes(1);
    expect(deps.notify).toHaveBeenCalledWith(expectedError);

    deps.notify.mockClear();
    parseDesignFileMock.mockReturnValueOnce(null as any);
    handlers.directValidate(undefined, 'invalid-yaml', 'Bad Design');
    expect(deps.designValidationActorRef.send).not.toHaveBeenCalled();
    expect(deps.notify).toHaveBeenCalledTimes(1);
    expect(deps.notify).toHaveBeenCalledWith(expectedError);
  });
});
