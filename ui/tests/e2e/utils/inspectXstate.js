// Usage: await inspectXstate(page, machine, options);
import { waitFor } from './waitFor';

/**
 * Waits for a specific event from the XState debugger.
 * @param {Object} params - Parameters object.
 * @param {import('playwright').Page} params.page - Playwright page instance.
 * @param {string} params.eventType - The type of event to wait for.
 * @param {number} [params.timeout_after=3000] - Timeout duration in milliseconds.
 * @returns {Promise<import('playwright').JSHandle>} - Promise resolving to the handle of the received event.
 */
export const waitForEvent = async ({ page, eventType, timeout_after = 3000 }) => {
  await waitFor(
    () =>
      page.evaluate(() => {
        const debuggerActor = window?.debuggingActorRef;
        if (!debuggerActor) {
          return null;
        }
        return debuggerActor?.getSnapshot?.()?.status == 'active';
      }),
    timeout_after,
  );

  const eventHandle = await page.evaluateHandle(
    async ([eventType, timeout_after]) => {
      const debuggerActor = window?.debuggingActorRef;
      const timeout = setTimeout(() => {
        throw new Error(`Timeout after ${timeout_after} ms for event ${eventType}`);
      }, timeout_after);

      return new Promise((resolve) => {
        const in_buffer = debuggerActor
          .getSnapshot()
          .context.events_buffer.find((e) => e.data.eventType == eventType);

        if (in_buffer) {
          debuggerActor.send({ type: 'FLUSH_BUFFER' });
          return resolve(in_buffer.data.incommingEvent);
        }

        const subscription = debuggerActor.on('LOGGED_EVENT', (event) => {
          // console.log("Inspect", event.data.eventType);
          if (event.data.eventType == eventType) {
            subscription.unsubscribe();
            debuggerActor.send({ type: 'FLUSH_BUFFER' });
            clearTimeout(timeout);
            resolve(event?.data?.incommingEvent);
          }
        });
      });
    },
    [eventType, timeout_after],
  );

  return eventHandle;
};
