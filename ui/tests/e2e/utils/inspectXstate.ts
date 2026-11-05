import { Page, JSHandle } from '@playwright/test';
import { waitFor } from './waitFor';

interface WaitForEventParams {
  page: Page;
  eventType: string;
  timeout_after?: number;
}

export const waitForEvent = async ({
  page,
  eventType,
  timeout_after = 3000,
}: WaitForEventParams): Promise<JSHandle> => {
  await waitFor(
    () =>
      page.evaluate(() => {
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
      return new Promise((resolve, reject) => {
        const debuggerActor = (window as any)?.debuggingActorRef;

        const timeout = setTimeout(() => {
          reject(new Error(`Timeout after ${timeout_after} ms for event ${eventType}`));
        }, timeout_after);

        const in_buffer = debuggerActor
          .getSnapshot()
          .context.events_buffer.find((e: any) => e.data.eventType === eventType);

        if (in_buffer) {
          clearTimeout(timeout);
          debuggerActor.send({ type: 'FLUSH_BUFFER' });
          return resolve(in_buffer.data.incommingEvent);
        }

        const subscription = debuggerActor.on('LOGGED_EVENT', (event: any) => {
          if (event.data.eventType === eventType) {
            subscription.unsubscribe();
            clearTimeout(timeout);
            debuggerActor.send({ type: 'FLUSH_BUFFER' });
            resolve(event?.data?.incommingEvent);
          }
        });
      });
    },
    [eventType, timeout_after] as const,
  );

  return eventHandle;
};
