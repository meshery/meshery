import { EventBus } from '@sistent/sistent';

/**
 * Cross-boundary event bus for communication between the Meshery UI core
 * and Meshery extensions (which run in separate contexts).
 *
 * Event types:
 * - DISPATCH_TO_MESHERY_STORE: Extensions -> Redux store dispatch (subscribed in store/index.ts)
 * - K8S_CONTEXTS_UPDATED: Redux store -> Extensions notification (published in mesheryUi.ts)
 *
 * Do NOT use this for intra-UI communication -- use Redux dispatch or XState events instead.
 * This bus exists solely for the extension boundary.
 */
export const mesheryEventBus = new EventBus();
