/* eslint-disable */
import { Page, JSHandle } from '@playwright/test';
import { waitFor } from './waitFor';

interface WaitForEventParams {
  page: Page;
  eventType: string;
  timeout_after?: number;
}

/**
 * Waits for a specific event from the XState debugger.
 */
export const waitForEvent = async ({
  page,
  eventType,
  timeout_after = 3000,
}: WaitForEventParams): Promise<JSHandle> => {
  await waitFor(
    () =>
      page.evaluate(() => {
        // We cast window to any to access the debugger
        const debuggerActor = (window as any)?.debuggingActorRef;
        if (!debuggerActor) {
          return null;
        }
        return debuggerActor?.getSnapshot?.()?.status === 'active';
      }),
    timeout_after,
  );

  const eventHandle = await page.evaluateHandle(
    async ([eventType, timeout_after]) => {
      const debuggerActor = (window as any)?.debuggingActorRef;

      const timeout = setTimeout(() => {
        throw new Error(`Timeout after ${timeout_after} ms for event ${eventType}`);
      }, timeout_after);

      return new Promise((resolve) => {
        const in_buffer = debuggerActor
          .getSnapshot()
          // Using any for the event buffer to avoid complex typing
          .context.events_buffer.find((e: any) => e.data.eventType === eventType);

        if (in_buffer) {
          debuggerActor.send({ type: 'FLUSH_BUFFER' });
          return resolve(in_buffer.data.incommingEvent);
        }

        const subscription = debuggerActor.on('LOGGED_EVENT', (event: any) => {
          if (event.data.eventType === eventType) {
            subscription.unsubscribe();
            debuggerActor.send({ type: 'FLUSH_BUFFER' });
            clearTimeout(timeout);
            resolve(event?.data?.incommingEvent);
          }
        });
      });
    },
    [eventType, timeout_after] as const,
  );

  return eventHandle;
};
