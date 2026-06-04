import { describe, expect, it } from 'vitest';
import { mesheryEventBus } from '../eventBus';

describe('mesheryEventBus', () => {
  it('is a singleton EventBus instance with publish/on/onAny methods', () => {
    expect(mesheryEventBus).toBeDefined();
    expect(typeof mesheryEventBus.publish).toBe('function');
    expect(typeof mesheryEventBus.on).toBe('function');
    expect(typeof mesheryEventBus.onAny).toBe('function');
  });

  it('delivers published events to subscribers of the matching type', async () => {
    const received: Array<{ type: string; data: { count: number } }> = [];
    const subscription = mesheryEventBus
      .on('DISPATCH_TO_MESHERY_STORE')
      .subscribe((event: { type: string; data: { count: number } }) => {
        received.push(event);
      });

    mesheryEventBus.publish({ type: 'DISPATCH_TO_MESHERY_STORE', data: { count: 1 } });
    mesheryEventBus.publish({ type: 'K8S_CONTEXTS_UPDATED', data: { count: 99 } });

    // Allow rxjs microtask flush
    await Promise.resolve();

    expect(received).toHaveLength(1);
    expect(received[0]).toEqual({ type: 'DISPATCH_TO_MESHERY_STORE', data: { count: 1 } });
    subscription.unsubscribe();
  });

  it('delivers events to onAny subscribers regardless of type', async () => {
    const seenTypes: string[] = [];
    const subscription = mesheryEventBus.onAny().subscribe((event: { type: string }) => {
      seenTypes.push(event.type);
    });

    mesheryEventBus.publish({ type: 'K8S_CONTEXTS_UPDATED', data: {} });
    mesheryEventBus.publish({ type: 'DISPATCH_TO_MESHERY_STORE', data: {} });

    await Promise.resolve();

    expect(seenTypes).toEqual(
      expect.arrayContaining(['K8S_CONTEXTS_UPDATED', 'DISPATCH_TO_MESHERY_STORE']),
    );
    subscription.unsubscribe();
  });

  it('stops delivering events to a subscriber after it unsubscribes', async () => {
    const received: unknown[] = [];
    const subscription = mesheryEventBus.on('K8S_CONTEXTS_UPDATED').subscribe((event) => {
      received.push(event);
    });

    mesheryEventBus.publish({ type: 'K8S_CONTEXTS_UPDATED', data: { round: 1 } });
    await Promise.resolve();
    expect(received).toHaveLength(1);

    subscription.unsubscribe();
    mesheryEventBus.publish({ type: 'K8S_CONTEXTS_UPDATED', data: { round: 2 } });
    await Promise.resolve();
    expect(received).toHaveLength(1);
  });
});
