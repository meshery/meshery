import {
  EventBus,
  type MesheryExtensionEvent,
  type MesheryExtensionEventBus,
} from '@sistent/sistent';

/**
 * Cross-boundary event bus for communication between the Meshery UI core
 * and Meshery extensions (which run in separate contexts).
 *
 * The set of events allowed on this bus is declared once, in sistent's
 * `mesheryExtensionContract` module, and shared by the host and every extension
 * bundle. Read that module for the event list and payload shapes; never widen
 * this bus by publishing a literal that is not part of the union.
 *
 * The type argument is load-bearing. `EventBus<T>` has no default for `T`, so a
 * bare `new EventBus()` collapses `T` to its constraint and `publish()` silently
 * accepts any `{ type: string }`. A rename on either side of the boundary then
 * compiles cleanly and the feature becomes a runtime no-op with no error.
 *
 * Do NOT use this for intra-UI communication -- use Redux dispatch or XState events instead.
 * This bus exists solely for the extension boundary.
 */
export const mesheryEventBus: MesheryExtensionEventBus = new EventBus<MesheryExtensionEvent>();
