import { describe, expect, it, vi } from 'vitest';

// The step bodies pull in the full connections UI (and its icon assets); onNext
// only needs the payload helpers, so the rendered components are stubbed.
vi.mock('../ConnectionWizardStepContent', () => ({
  ConnectionKindSelectionStep: () => null,
  ConnectionReviewStep: () => null,
  CredentialAssociationStep: () => null,
  GenericConnectionDetailsStep: () => null,
}));
vi.mock('./useKindPermission', () => ({ useKindPermission: () => () => true }));

import { genericRegisterStep } from './genericSteps';
import type { GenericRecord, WizardContext, WizardData } from './types';

// Drives genericRegisterStep.onNext the way the wizard does: each attempt sees
// the data as of its start (a React render), and patches apply afterwards.
const setup = () => {
  let data = {
    kindConfig: {
      kind: 'prometheus',
      label: 'Prometheus',
      type: 'telemetry',
      subType: 'monitoring',
    },
    connectionFormData: { url: 'http://prometheus:9090' },
    credentialMode: 'new',
    selectedCredentialId: '',
    credentialName: 'prom-cred',
    credentialFormData: {},
    skipCredentialVerification: true,
    registrationId: null,
    registered: null,
    registrationResult: null,
    registrationError: null,
    postConfig: {},
  } as unknown as WizardData;

  const services = {
    notify: vi.fn(),
    registerConnection: vi.fn(async (_body: GenericRecord) => ({})),
    connectConnection: vi.fn(async (_body: GenericRecord) => ({})),
    cancelRegistration: vi.fn(async (_id: string) => {}),
    credentials: [],
  };

  const attempt = () => {
    const ctx = {
      mode: 'create',
      data,
      patch: (partial: Partial<WizardData>) => {
        data = { ...data, ...partial };
      },
      services,
    } as unknown as WizardContext;
    return genericRegisterStep.onNext!(ctx);
  };

  return {
    services,
    attempt,
    edit: (partial: Partial<WizardData>) => {
      data = { ...data, ...partial };
    },
    get data() {
      return data;
    },
  };
};

describe('genericRegisterStep.onNext', () => {
  it('registers then connects on the first attempt', async () => {
    const w = setup();

    await expect(w.attempt()).resolves.toBe(true);

    expect(w.services.registerConnection).toHaveBeenCalledTimes(1);
    expect(w.services.connectConnection).toHaveBeenCalledTimes(1);
    const registerBody = w.services.registerConnection.mock.calls[0][0];
    expect(registerBody.status).toBe('register');
    expect(w.services.connectConnection.mock.calls[0][0].id).toBe(registerBody.id);
    expect(w.services.cancelRegistration).not.toHaveBeenCalled();
  });

  it('retries only connect when the form is unchanged after a failed connect', async () => {
    const w = setup();
    w.services.connectConnection.mockRejectedValueOnce(new Error('connect failed'));

    await expect(w.attempt()).resolves.toBe(false);
    const registeredId = w.data.registered?.id;
    expect(registeredId).toBeTruthy();

    await expect(w.attempt()).resolves.toBe(true);

    expect(w.services.registerConnection).toHaveBeenCalledTimes(1);
    expect(w.services.connectConnection).toHaveBeenCalledTimes(2);
    expect(w.services.connectConnection.mock.calls[1][0].id).toBe(registeredId);
    expect(w.services.cancelRegistration).not.toHaveBeenCalled();
  });

  it('restarts the registration when the form was edited after a failed connect', async () => {
    const w = setup();
    w.services.connectConnection.mockRejectedValueOnce(new Error('connect failed'));

    await w.attempt();
    const staleId = w.data.registered!.id;

    w.edit({ connectionFormData: { url: 'http://prometheus.monitoring:9090' } });
    await expect(w.attempt()).resolves.toBe(true);

    expect(w.services.cancelRegistration).toHaveBeenCalledWith(staleId);
    expect(w.services.registerConnection).toHaveBeenCalledTimes(2);
    const reregister = w.services.registerConnection.mock.calls[1][0];
    expect(reregister.id).not.toBe(staleId);
    expect(reregister.metadata).toEqual({ url: 'http://prometheus.monitoring:9090' });
    expect(w.services.connectConnection.mock.calls[1][0].id).toBe(reregister.id);
    expect(w.data.registrationId).toBe(reregister.id);
  });

  it('keeps the stale registration when cancelling it fails, so a later retry can try again', async () => {
    const w = setup();
    w.services.connectConnection.mockRejectedValueOnce(new Error('connect failed'));
    await w.attempt();
    const staleId = w.data.registered!.id;

    w.edit({ connectionFormData: { url: 'http://other:9090' } });
    w.services.cancelRegistration.mockRejectedValueOnce(new Error('cancel failed'));

    await expect(w.attempt()).resolves.toBe(false);
    expect(w.services.registerConnection).toHaveBeenCalledTimes(1);
    expect(w.data.registered?.id).toBe(staleId);
    expect(w.data.registrationId).toBe(staleId);
  });
});
