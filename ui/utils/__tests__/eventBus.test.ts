import { describe, expect, it } from 'vitest';
import { MESHERY_EXTENSION_EVENT, type MesheryExtensionEvent } from '@sistent/sistent';
import { mesheryEventBus } from '../eventBus';

const dispatchToStore: MesheryExtensionEvent = {
  type: MESHERY_EXTENSION_EVENT.DispatchToMesheryStore,
  data: { type: 'core/updateProgress', payload: { showProgress: true } },
};

const k8sContextsUpdated: MesheryExtensionEvent = {
  type: MESHERY_EXTENSION_EVENT.K8sContextsUpdated,
  data: { selectedK8sContexts: ['ctx-1'] },
};

describe('mesheryEventBus', () => {
  it('is a singleton EventBus instance with publish/on/onAny methods', () => {
    expect(mesheryEventBus).toBeDefined();
    expect(typeof mesheryEventBus.publish).toBe('function');
    expect(typeof mesheryEventBus.on).toBe('function');
    expect(typeof mesheryEventBus.onAny).toBe('function');
  });

  it('rejects event types outside the extension contract at compile time', () => {
    // @ts-expect-error - the whole point of typing the bus: an event literal the
    // contract does not declare must fail the build. At runtime it would be
    // published to nobody, which is indistinguishable from a working feature.
    mesheryEventBus.publish({ type: 'NOT_IN_THE_CONTRACT', data: {} });
  });

  it('delivers published events to subscribers of the matching type', async () => {
    const received: MesheryExtensionEvent[] = [];
    const subscription = mesheryEventBus
      .on(MESHERY_EXTENSION_EVENT.DispatchToMesheryStore)
      .subscribe((event) => {
        received.push(event);
      });

    mesheryEventBus.publish(dispatchToStore);
    mesheryEventBus.publish(k8sContextsUpdated);

    // Allow rxjs microtask flush
    await Promise.resolve();

    expect(received).toHaveLength(1);
    expect(received[0]).toEqual(dispatchToStore);
    subscription.unsubscribe();
  });

  it('delivers events to onAny subscribers regardless of type', async () => {
    const seenTypes: string[] = [];
    const subscription = mesheryEventBus.onAny().subscribe((event) => {
      seenTypes.push(event.type);
    });

    mesheryEventBus.publish(k8sContextsUpdated);
    mesheryEventBus.publish(dispatchToStore);

    await Promise.resolve();

    expect(seenTypes).toEqual(
      expect.arrayContaining([
        MESHERY_EXTENSION_EVENT.K8sContextsUpdated,
        MESHERY_EXTENSION_EVENT.DispatchToMesheryStore,
      ]),
    );
    subscription.unsubscribe();
  });

  it('stops delivering events to a subscriber after it unsubscribes', async () => {
    const received: MesheryExtensionEvent[] = [];
    const subscription = mesheryEventBus
      .on(MESHERY_EXTENSION_EVENT.K8sContextsUpdated)
      .subscribe((event) => {
        received.push(event);
      });

    mesheryEventBus.publish({
      type: MESHERY_EXTENSION_EVENT.K8sContextsUpdated,
      data: { selectedK8sContexts: ['round-1'] },
    });
    await Promise.resolve();
    expect(received).toHaveLength(1);

    subscription.unsubscribe();
    mesheryEventBus.publish({
      type: MESHERY_EXTENSION_EVENT.K8sContextsUpdated,
      data: { selectedK8sContexts: ['round-2'] },
    });
    await Promise.resolve();
    expect(received).toHaveLength(1);
  });
});
