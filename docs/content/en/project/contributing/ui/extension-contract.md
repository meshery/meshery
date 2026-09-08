---
title: Contributing to Meshery UI - Extension Contract and Event Bus
description: How the Meshery UI core and remotely-loaded Meshery extensions communicate - the load pipeline, the shared runtime, the injected capability bag, the Meshery event bus, and the rules for changing any of them.
categories: [contributing]
aliases: [/project/contributing/contributing-ui-extension-contract]
---

<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisite Reading</strong></p>
  <ol>
    <li><a href="{{< ref "project/contributing/ui/ui" >}}">Contributing to Meshery UI</a></li>
    <li><a href="{{< ref "reference/extensibility/ui.md" >}}">Extensibility: UI</a></li>
    <li><a href="{{< ref "reference/extensibility/providers/index.md" >}}">Extensibility: Providers</a></li>
  </ol>
</div>

## Table of Contents

- [The boundary](#the-boundary)
- [Why there is a contract](#why-there-is-a-contract)
- [Where the contract lives](#where-the-contract-lives)
- [How an extension gets loaded](#how-an-extension-gets-loaded)
  - [The shared runtime](#the-shared-runtime)
  - [Export shape and its two failure modes](#export-shape-and-its-two-failure-modes)
- [The injected capability bag](#the-injected-capability-bag)
  - [What is in the bag](#what-is-in-the-bag)
  - [How the extension consumes it](#how-the-extension-consumes-it)
  - [Contract versioning](#contract-versioning)
  - [Backward compatibility](#backward-compatibility)
  - [Worked example: adding a capability](#worked-example-adding-a-capability)
- [The Meshery event bus](#the-meshery-event-bus)
  - [Delivery semantics](#delivery-semantics)
  - [Publishing](#publishing)
  - [Subscribing](#subscribing)
- [Event reference](#event-reference)
- [Worked example: adding an event](#worked-example-adding-an-event)
- [Changing the contract](#changing-the-contract)
- [Tests that gate the boundary](#tests-that-gate-the-boundary)
- [Troubleshooting](#troubleshooting)
- [File reference](#file-reference)

## The boundary

Meshery extensions — Kanvas being the one most contributors will meet — are **not** compiled
with Meshery UI. They are separately-built, separately-released JavaScript bundles that
Meshery Server proxies and Meshery UI fetches, evaluates, and mounts at runtime.

```
  meshery/meshery (host)                      meshery/meshery-extensions (bundle)
  ─────────────────────────                   ──────────────────────────────────
  NavigatorExtension.tsx                      App.tsx
      │                                            │
      │  injectProps  ────────────────────────►  destructured, normalized,
      │  (one untyped object)                    stashed in module-scope maps
      │                                            │
      │  ◄──────────  mesheryEventBus  ──────────►  │   (bi-directional, typed)
      │                                            │
      └──────────────┬─────────────────────────────┘
                     │
             @sistent/sistent
          mesheryExtensionContract
     (event literals, payload types, declared
      capabilities, contract version)
```

Everything crossing that line is either a key on `injectProps` or an event on
`mesheryEventBus`. Both are declared once, in sistent, and shared by the two repos.

## Why there is a contract

There is **no compile-time link across the boundary**:

- The host never imports the extension, so TypeScript cannot verify that the keys the host
  provides are the keys the extension reads.
- The extension never imports the host, so renaming a key here still compiles there.
- Bundle and host are released independently. A bundle built against one revision of the
  host may be loaded by an older or newer host, indefinitely.

A rename on either side therefore does not fail the build and does not throw. It resolves
to `undefined` at the extension's use site, and the feature becomes a **silent runtime
no-op**. This is not hypothetical — two shipped examples:

| Change | Symptom in production |
| --- | --- |
| `OPEN_DESIGN_IN_KANVAS` → `OPEN_DESIGN_IN_EXTENSION` | "Open in Kanvas" published an event nobody was subscribed to. No error; the click did nothing. |
| `capabilitiesRegistry` → `providerCapabilities` | The extension constructed its access-control object from `undefined`. No error; permission-gated UI silently misbehaved. |

The contract exists to convert that class of failure into a compile error or a failing test.

## Where the contract lives

The contract is declared **once**, in the `mesheryExtensionContract` module of
[`@sistent/sistent`](https://github.com/meshery/sistent). Read that module for the
authoritative definitions. The exports Meshery UI consumes:

| Export | Kind | Purpose |
| --- | --- | --- |
| `MESHERY_EXTENSION_EVENT` | const enum-like object | Every event literal allowed on the bus. **Derive every event `type` from this — never type the string.** |
| `MesheryExtensionEvent` | type | Discriminated union of `{ type, data }` across all contract events. |
| `MesheryExtensionEventBus` | type | Type of the shared bus singleton. |
| `EventBus<T>` | class | The bus implementation itself (RxJS-backed). |
| `MESHERY_EXTENSION_CONTRACT_VERSION` | const | Revision of the contract the host was built against; injected into every extension. |
| `reportInjectedCapabilities` | function | Given an `injectProps` bag, reports which declared capabilities are present and which are missing. |
| `isInjectedCapabilityReportSatisfied` | function | Whether a report has no missing capabilities. |
| `describeInjectedCapabilityReport` | function | Human-readable rendering of a report, used as test-failure output. |
| `createCanShow` | function | Builds the permission-gated `CanShow` component, and publishes missing-permission / missing-capability events onto the bus it is given. |

{{% alert color="warning" title="Sistent is the source of truth" %}}
Do not add a second declaration of an event literal or an injected key anywhere in
`meshery/meshery`. Hand-duplicated literals are the exact defect the contract exists to
prevent. If a capability or event is missing, add it to sistent first, then consume it here.
{{% /alert %}}

## How an extension gets loaded

The full path from provider capability to mounted component:

1. **The provider declares the extension.** `GET /api/provider/capabilities` returns an
   `extensions` object keyed by extension point — `navigator`, `account`, `userPrefs`,
   `collaborator`, plus `full_page` entries. Each entry carries an `href` (the route it
   claims) and a `component` (the bundle URI).
2. **The Navigator renders menu entries** for the `navigator` extension point.
   `ui/components/layout/Navigator/Navigator.tsx` runs the declared entries through
   `ExtensionPointSchemaValidator('navigator')` and renders them as menu items.
3. **The route resolves.** Navigating to a declared `href` lands on the catch-all page
   `ui/pages/extension/[...component].tsx`, which matches the current path against the
   provider's declared extensions and dispatches `extensionType` into Redux.
4. **`ExtensionSandbox` picks the bundle.**
   `ui/components/ExtensionSandbox.tsx` walks the extension tree for the entry whose `href`
   matches the current path, takes its `component` URI, and turns it into a fetchable URL
   with `createPathForRemoteComponent()` — which prefixes `/api/provider/extension`. The
   bundle is therefore served **through Meshery Server**, not fetched cross-origin from the
   provider.
5. **`NavigatorExtension` loads and mounts it.**
   `ui/components/layout/Navigator/NavigatorExtension.tsx` calls `useRemoteComponent(url)`,
   which fetches the bundle, evaluates it as a CommonJS module against the shared runtime,
   and returns `module.exports.default`. The host then renders it with exactly one prop:
   `injectProps`.

For non-navigator extension points, `ExtensionSandbox` renders the generic
`ui/components/general/RemoteComponent` instead, which forwards whatever props it was given
but mounts the bundle **without** the capability bag. Only the Navigator extension point
receives `injectProps`. That path also has none of the export-shape diagnostics described
below — a mis-built bundle mounted through it still fails with React's opaque error, so
debug it against the checklist in [Troubleshooting](#troubleshooting) rather than waiting
for a useful message.

### The shared runtime

A remote bundle must not ship its own copy of React — two React instances in one page
produce "invalid hook call" failures, and two MUI/emotion instances produce broken theming.
So the bundle is built with those packages as **externals**, and the host satisfies each
`require()` at evaluation time from `ui/remote-component.config.js`:

```javascript
module.exports = {
  resolve: {
    react: require('react'),
    'react-dom': require('react-dom'),
    '@mui/material': require('@mui/material'),
    '@emotion/react': require('@emotion/react'),
    // …plus xstate, redux, rjsf, cytoscape, dockview, xterm, dnd-kit, tippy, …
  },
};
```

Three consequences worth knowing:

- **A missing key is a hard failure.** If a bundle `require()`s something absent from
  `resolve`, remote-component throws *"… does not exist in dependencies"* and the extension
  does not mount. Adding a new third-party dependency to an extension therefore requires a
  matching entry here.
- **CSS subpath imports need stub entries.** Bundles do side-effect-only imports such as
  `require('@xterm/xterm/css/xterm.css')`. The actual stylesheets are injected globally in
  `ui/pages/_app.tsx`, so the config maps those subpaths to `{}` — a key must exist or the
  resolver throws, but the value is irrelevant.
- **Versions must stay in lockstep.** The host's installed version of a shared package is
  the one the extension gets, regardless of what the extension declared. A major-version
  bump on either side is a coordinated change.

### Export shape and its two failure modes

The bundle must export its component as a CommonJS default:

```javascript
module.exports = { default: Component, __esModule: true };
```

In practice that means building with `output.library.type = "commonjs2"`. A bundle built
any other way still loads — it just resolves to nothing useful. Two distinct failures
follow, and `NavigatorExtension` reports them with **different** messages on purpose:

| Failure | What the loader returns | Actual cause | Remedy the host reports |
| --- | --- | --- | --- |
| No CommonJS default export | `RemoteComponent === undefined`, `err === undefined` — the hook does **not** throw | Bundler configuration | Rebuild with `output.library.type = "commonjs2"` and republish |
| Default export is not a component (e.g. `<Foo />` instead of `Foo`) | An object carrying a `$$typeof` tag | Defect in the extension's own source | Export the component itself, not a rendered element |

Without those guards, both cases surface only as React's opaque *"Element type is invalid …
got: undefined"*, thrown from deep inside the render tree, far from the real problem.

The component check is deliberately an **allow-list** of `$$typeof` tags
(`react.forward_ref`, `react.memo`, `react.lazy`) rather than a bare `'$$typeof' in value`
test. Many React-internal values carry `$$typeof` without being component types — most
notably an already-rendered element, which is precisely the mistake the second row above
describes. A bare check would wave it straight through.

If an extension renders blank with nothing in the console, inspect the bundle's export
shape before anything else.

## The injected capability bag

`buildExtensionInjectProps` in `NavigatorExtension.tsx` builds the bag. It is deliberately a
**pure factory** — everything it needs from React state or context arrives as a parameter —
so a unit test can build the bag and assert it against the contract without rendering
anything. Keep it that way; the moment it reads a hook directly, the gate test gets harder
to write than the shortcut is worth.

The component memoizes the result on its inputs, so extensions are not re-rendered by an
unrelated store update.

### What is in the bag

| Group | Keys |
| --- | --- |
| Contract | `contractVersion` |
| Components | `PatternServiceFormCore`, `InfoModal`, `ViewInfoModal`, `ExportModal`, `GenericRJSFModal`, `RJSForm`, `TypingFilter`, `CreateModelModal`, `ImportModelModal`, `ValidateDesign`, `DryRunDesign`, `DeployStepper`, `UnDeployStepper`, `MesheryPerformanceComponent`, `StructuredDataFormatter`, `RelationshipEvaluationResponseFormatter`, `ThemeTogglerCore`, `_PromptComponent` |
| Hooks | `hooks.CAN`, `hooks.useHasPermission`, `hooks.useFilterK8sContexts`, `hooks.useDynamicComponent`, and `useNotificationHook` at top level |
| State | `providerCapabilities`, `selectedK8sContexts`, `currentOrganization`, `mesheryStore` |
| Callbacks | `openWorkspaceModal`, `openRegistryModal`, `SetCurrentLoadedResourceInOrgWorkspaceSession` |
| Machines | `designValidationMachine` |
| Bus | `mesheryEventBus` |
| Access control | `ProviderUiAccessControlClass` (and its legacy alias `CapabilitiesRegistryClass`) |
| Legacy | `resolver` |

Two of these deserve elaboration:

**`mesheryStore` is getter-based, not a snapshot.**

```javascript
const extensionExposedMesheryStore = {
  selectedK8sClusters: { get: () => selectSelectedK8sClusters(store.getState()) },
  k8sConfig: { get: () => selectK8sConfig(store.getState()) },
};
```

The extension calls `.get()` at the moment it needs the value, so it always reads live
store state. Handing over a plain value instead would freeze it at mount time. Add new
store exposure the same way.

**`ProviderUiAccessControlClass` is a class, not an instance.** The host passes the class;
the extension constructs it with the raw `/api/provider/capabilities` payload
(`new CapabilitiesRegistryClass(capabilitiesRegistry)`). `ProviderUiAccessControl`
(`ui/utils/disabledComponents.js`) answers the provider-driven questions —
`isNavigatorComponentEnabled`, `isHeaderComponentEnabled`, `isExtensionComponentEnabled` —
which in a restricted (playground) environment walk `restrictedAccess.allowedComponents` to
decide whether a given piece of UI may render at all.

### How the extension consumes it

On the extension side (`meshery-extensions/ui/src/App.tsx` and
`ui/src/globals/mesherySdk.ts`), the bag is destructured once at the top level and stashed
into module-scoped maps, so deeply-nested extension code can reach it without prop drilling:

```typescript
setMesheryFunctionality("mesheryEventBus", mesheryEventBus);
setMesheryFunctionality("CAN", normalizedCAN);
injectedReactComponents.set(GLOBAL_COMPONENTS.ValidateDesign, normalizedValidateDesign);
injectedMachines.set(GLOBAL_MACHINES.DesignValidatorMachine, normalizedDesignValidationMachine);
```

and read back through typed accessors — `getMesheryEventBus()`, `getCapabilitiesRegistry()`,
`createInjectedComponent("ExportDesignModal")`, and so on. `injectedMesheryFunctionalitiesMap`
in `mesherySdk.ts` is the extension's local typing of the same contract.

Every value passes through `normalizeInjectedValue` first. It repeatedly unwraps
`{ default, __esModule }` wrappers until it reaches the real value:

```typescript
while (hasDefaultExport(resolvedValue) && !seen.has(resolvedValue)) {
  resolvedValue = resolvedValue.default;
}
```

This exists because the ESM/CJS interop between host build, remote-component evaluation and
bundle build can double-wrap an injected component, so `Component` arrives as
`{ default: Component }`. Rendering that directly fails with the same opaque React error as
a bad export. **If you inject a value that is itself an object with a `default` key, be
aware it will be unwrapped** — the `seen` set only guards against cycles, not against
legitimately-shaped data.

### Contract versioning

`contractVersion: MESHERY_EXTENSION_CONTRACT_VERSION` is injected into every bundle. Nothing
at build time can catch a host/bundle mismatch, because bundles are published artifacts
loaded by whichever host happens to be deployed. The version gives the extension a runtime
signal it can use to degrade gracefully or warn, rather than silently misbehaving.

Bump it in sistent whenever a change is not backward-compatible for already-published
bundles.

### Backward compatibility

Published bundles are loaded by new hosts, so keys cannot simply be deleted:

- **`CapabilitiesRegistryClass`** is retained as an alias of `ProviderUiAccessControlClass`.
  The rename happened here; the alias keeps older bundles working. The unit test asserts the
  two keys are the same reference, so the alias cannot drift.
- **`resolver`** is retained even though Meshery Server no longer exposes any GraphQL
  subscription. Published bundles call
  `resolver.subscription.ConfigurationSubscription(cb, vars)` and later `.dispose()` on the
  returned handle; handing them `undefined` would throw at mount. The host supplies an inert
  subscription matching the old call shape — `(onNext, variables) => ({ dispose })` — that
  never invokes `onNext`. New extensions must read designs and filters from the REST API.

When you retire a capability, keep an alias or an inert stand-in until every published bundle
has been rebuilt, and leave a comment saying why the dead-looking key is there. Otherwise the
next contributor deletes it as cruft.

### Worked example: adding a capability

Say you want extensions to be able to open the environment-selection modal.

1. **Declare it in sistent.** Add `openEnvironmentModal` to the declared capability list in
   `mesheryExtensionContract` so `reportInjectedCapabilities` knows to look for it.
2. **Add the parameter to the factory** in `NavigatorExtension.tsx` — it comes from a hook or
   context, so it must be a parameter, not a direct hook call:

   ```typescript
   export type ExtensionInjectPropsDeps = {
     // …existing
     openEnvironmentModal: unknown;
   };

   export const buildExtensionInjectProps = ({ /* … */ openEnvironmentModal }) => ({
     // …existing
     openEnvironmentModal,
   });
   ```

3. **Wire it at the call site**, inside the existing `useMemo`, and add it to the dependency
   array — a missed dependency here hands the extension a stale callback.
4. **Consume it in the extension**: add it to `AppProps['injectProps']` and
   `injectedMesheryFunctionalitiesMap`, then `setMesheryFunctionality("openEnvironmentModal", …)`.
5. **Run `npm test` in `ui/`.** The capability-report assertion in
   `NavigatorExtension.test.tsx` now covers the new key; if you forgot step 2 or 3 it fails
   here rather than in production.

## The Meshery event bus

`ui/utils/eventBus.ts` exports the singleton:

```typescript
import { EventBus, type MesheryExtensionEvent, type MesheryExtensionEventBus } from '@sistent/sistent';

export const mesheryEventBus: MesheryExtensionEventBus = new EventBus<MesheryExtensionEvent>();
```

{{% alert color="warning" title="The type argument is load-bearing" %}}
`EventBus<T>` has **no default** for `T`. A bare `new EventBus()` collapses `T` to its
constraint, and `publish()` then silently accepts any `{ type: string }` — disabling
publish-site checking entirely, which is the single check standing between a typo and a
dead feature. Keep this bus typed as `EventBus<MesheryExtensionEvent>`.
{{% /alert %}}

{{% alert color="info" title="This bus is only for the extension boundary" %}}
Do **not** use `mesheryEventBus` for intra-UI communication. Within Meshery UI, use Redux
dispatch or XState events. Every literal on this bus is a cross-repo API surface; adding one
for a purely internal purpose burdens the extension repo with a contract it will never use.
{{% /alert %}}

### Delivery semantics

The implementation is a thin wrapper over an RxJS `Subject`:

```typescript
publish(event) { this.eventSubject.next(event); }
on(type)       { return this.eventObservable.pipe(filter((e) => e.type === type)); }
onAny()        { return this.eventObservable; }
```

That has three practical consequences:

- **No replay, no buffering.** It is a plain `Subject`, not a `ReplaySubject` or
  `BehaviorSubject`. An event published while nobody is subscribed is dropped permanently —
  a late subscriber never sees it. This is the reason the navigation publishers first check
  `isExtensionOpen()`: publishing to an unmounted extension would be a no-op, so they fall
  back to a router push instead.
- **Delivery is synchronous.** `publish()` runs every matching subscriber before it returns.
  Do not rely on a publish being deferred, and be aware that an exception thrown by one
  subscriber propagates back into the publisher.
- **`on()` filters at runtime but is typed as the full union.** See the narrowing note under
  [Subscribing](#subscribing).

Subscriptions are RxJS subscriptions. Hold the handle and `unsubscribe()` on teardown —
this bus is a module-scoped singleton that outlives every component, so a leaked
subscription leaks for the lifetime of the page.

| Method | Behavior |
| --- | --- |
| `publish(event)` | Emits a `MesheryExtensionEvent`. Type-checked against the contract union. Synchronous. |
| `on(type)` | Observable of events of that type only. |
| `onAny()` | Observable of every event, regardless of type. Useful for debugging; avoid in feature code. |

### Publishing

Always derive the literal from the enum:

```javascript
import { MESHERY_EXTENSION_EVENT } from '@sistent/sistent';
import { mesheryEventBus } from '@/utils/eventBus';

mesheryEventBus.publish({
  type: MESHERY_EXTENSION_EVENT.OpenDesignInExtension,
  data: { designId, designName },
});
```

The host's navigation helpers in `ui/utils/utils.tsx` all follow the same shape — publish if
the extension is mounted, otherwise navigate:

```javascript
export const openDesignInExtension = (designId, designName, router) => {
  if (isExtensionOpen()) {
    mesheryEventBus.publish({
      type: MESHERY_EXTENSION_EVENT.OpenDesignInExtension,
      data: { designId, designName },
    });
    return;
  }
  router.push(`/extension/meshmap?mode=design&type=design&id=${designId}`);
};
```

Preserve that branch when adding similar actions. Given the no-replay semantics above,
publishing without the guard means the action is silently lost whenever the extension is not
already mounted.

### Subscribing

`ui/store/index.ts` wires the store to the bus and shows the one non-obvious detail:

```typescript
mesheryEventBus.on(MESHERY_EXTENSION_EVENT.DispatchToMesheryStore).subscribe((event) => {
  // `EventBus.on` filters at runtime but is typed as the full event union, so the
  // discriminant has to be re-checked here for `event.data` to narrow to a
  // dispatchable action rather than the union of every event payload.
  if (event.type !== MESHERY_EXTENSION_EVENT.DispatchToMesheryStore) return;
  store.dispatch(event.data);
});
```

The redundant-looking `if` is required for narrowing — `on()` is declared as
`Observable<T>`, not `Observable<Extract<T, { type: E }>>`. Do not delete it.

On the extension side the same bus is bridged into XState with `fromCallback`, so a state
machine can receive host events as ordinary machine events:

```typescript
export const MesheryEventsSubscriberActor = (events: MesheryEventTypes["type"][]) =>
  fromCallback(({ sendBack }) => {
    const subscribers = events.map((eventType) =>
      getMesheryEventBus().on(eventType).subscribe(sendBack)
    );
    return () => subscribers.forEach((s) => s.unsubscribe());
  });
```

Note the teardown function — that is the `unsubscribe()` discipline described above, and the
pattern to copy for any long-lived subscriber.

## Event reference

`MESHERY_EXTENSION_EVENT` in sistent is authoritative. The events Meshery UI publishes or
subscribes to today:

| Enum member | Wire literal | Direction | Payload | Site |
| --- | --- | --- | --- | --- |
| `DispatchToMesheryStore` | `DISPATCH_TO_MESHERY_STORE` | extension → host | `{ type, payload? }` — a Redux action | consumed in `ui/store/index.ts` |
| `K8sContextsUpdated` | `K8S_CONTEXTS_UPDATED` | host → extension | `{ selectedK8sContexts }` | `setK8sContexts` thunk, `ui/store/slices/mesheryUi.ts` |
| `OpenDesignInExtension` | `OPEN_DESIGN_IN_EXTENSION` | host → extension | `{ designId, designName }` | `openDesignInExtension`, `ui/utils/utils.tsx` |
| `OpenViewInExtension` | `OPEN_VIEW_IN_EXTENSION` | host → extension | `{ viewId, viewName }` | `openViewInExtension`, `ui/utils/utils.tsx` |
| `OpenViewScopedToDesign` | `OPEN_VIEW_SCOPED_TO_DESIGN` | host → extension | `{ designId, designName }` | `openViewScopedToDesignInOperator`, `ui/utils/utils.tsx` |
| `MergeDesign` | `MERGE_DESIGN` | host → extension | `{ id, name }` | `mergeDesignWithCurrent`, `ui/utils/utils.tsx` |

Plus the access-control events that shared sistent components emit onto whichever bus they
are given — in Meshery UI that wiring is `createCanShow(getProviderUiAccessControl, CAN, () => mesheryEventBus)`
in `ui/utils/can.ts`:

| Wire literal | Direction | Meaning |
| --- | --- | --- |
| `MISSING_PERMISSION` | host → extension | A `CanShow`-gated element was denied by the user's abilities. |
| `MISSING_CAPABILITY` | host → extension | A `CanShow`-gated element was denied by the provider's declared capabilities. |
| `FeatureRequiresUserAccount` | either | A feature was attempted that requires a signed-in account. |

Kanvas subscribes to those three to drive its anonymous-user flow — prompting sign-in or an
upgrade instead of rendering a dead control.

{{% alert color="info" title="DispatchToMesheryStore is the widest surface" %}}
This event lets an extension dispatch **arbitrary Redux actions into the host store**. That
makes host action-type strings and payload shapes an implicit part of the contract: renaming
a slice action can break an extension even though no contract file changed. Prefer a
purpose-built event, or an injected callback, over widening what extensions dispatch.
{{% /alert %}}

## Worked example: adding an event

Say the host needs to tell the extension that the active workspace changed.

1. **Declare it in sistent** — add the literal to `MESHERY_EXTENSION_EVENT` and the payload
   member to the `MesheryExtensionEvent` union:

   ```typescript
   WorkspaceChanged: 'WORKSPACE_CHANGED',
   // …
   type WorkspaceChangedEvent = {
     type: typeof MESHERY_EXTENSION_EVENT.WorkspaceChanged;
     data: { workspaceId: string; workspaceName: string };
   };
   ```

2. **Publish from the host**, deriving the literal from the enum, and guarding with
   `isExtensionOpen()` if the event is a navigation action:

   ```javascript
   mesheryEventBus.publish({
     type: MESHERY_EXTENSION_EVENT.WorkspaceChanged,
     data: { workspaceId, workspaceName },
   });
   ```

3. **Assert the publish site** in the owning slice/util test, the way
   `ui/store/slices/__tests__/mesheryUi.test.ts` asserts `K8sContextsUpdated`.
4. **Subscribe in the extension** — add the literal to `MesheryEventTypes` and to the
   relevant `MesheryEventsSubscriberActor` list (or a direct `.on(...).subscribe(...)` with
   teardown).
5. **Bump `MESHERY_EXTENSION_CONTRACT_VERSION`** if the change is not backward-compatible.
6. **Update the [event reference](#event-reference) table on this page.**

Skipping step 6 is how the table goes stale, and a stale table is worse than none — it is
what contributors read instead of the sistent module.

## Changing the contract

The order matters. Sistent first, then host, then extension:

1. **Change sistent.** Add the event literal and payload, or the declared capability, in
   `mesheryExtensionContract`.
2. **Bump `MESHERY_EXTENSION_CONTRACT_VERSION`** for anything not backward-compatible with
   already-published bundles.
3. **Update the host** — publish/subscribe via the enum member; add the capability to
   `buildExtensionInjectProps` and its call site's `useMemo` dependencies.
4. **Update `meshery/meshery-extensions`** to read the new key or handle the new event.
5. **Keep an alias or inert stand-in** for anything renamed or removed until published
   bundles are rebuilt.
6. **Run `npm test` in `ui/`** and update this page.

### Working against an unreleased sistent

Contract changes almost always mean an unreleased `@sistent/sistent`. During development,
`ui/node_modules/@sistent/sistent` is typically overwritten in place with a locally-built
dist. A green local run then proves nothing about CI, which installs the published package.

Re-verify with `npm ci` before declaring a contract change done. A local build usually keeps
the published version string, so **matching versions are not evidence that the installed
contents are the published ones**. A version *mismatch* against `ui/package.json` is a useful
tell that this has happened; a match proves nothing.

### Forbidden

- MUST NOT type an event literal by hand — derive it from `MESHERY_EXTENSION_EVENT`.
- MUST NOT construct the bus as a bare `new EventBus()`; it must be
  `new EventBus<MesheryExtensionEvent>()`.
- MUST NOT use `mesheryEventBus` for communication that stays inside Meshery UI.
- MUST NOT redeclare contract events or injected keys locally in `meshery/meshery`.
- MUST NOT drop an injected key without an alias or inert stand-in while older bundles are
  still in the wild.
- MUST NOT call hooks directly inside `buildExtensionInjectProps` — it must stay a pure
  factory so the gate test can call it.
- MUST NOT add a third-party dependency to an extension without a matching entry in
  `ui/remote-component.config.js`.

## Tests that gate the boundary

Run with `npm test` in `ui/` (Vitest).

| File | What it guarantees |
| --- | --- |
| `ui/components/layout/Navigator/NavigatorExtension.test.tsx` | The bag built by `buildExtensionInjectProps` satisfies **every** capability the contract declares, via `reportInjectedCapabilities` / `isInjectedCapabilityReportSatisfied`, with `describeInjectedCapabilityReport` as the failure message. **This is the gate that catches a capability rename before merge.** Also covers `contractVersion` being injected, the `CapabilitiesRegistryClass` alias holding, the getter-based `mesheryStore`, the missing-default-export and not-a-component diagnostics (including that only the former mentions `commonjs2`), `memo`/`forwardRef`-wrapped exports, and that a missing export is not reported while still loading. |
| `ui/utils/__tests__/eventBus.test.ts` | The bus is a singleton exposing `publish`/`on`/`onAny`; `on` filters by type while `onAny` does not; unsubscribing stops delivery; and — via `@ts-expect-error` — an event literal outside the contract **fails the build**. |
| `ui/store/slices/__tests__/mesheryUi.test.ts` | The `setK8sContexts` thunk publishes `K8sContextsUpdated` with the selected contexts. |

If you add an event, add a publish-site assertion alongside these. An event nobody asserts on
is indistinguishable from an event nobody receives.

## Troubleshooting

| Symptom | Likely cause | Where to look |
| --- | --- | --- |
| Extension area is blank, no console error | Bundle has no CommonJS default export | The host's own error panel now reports this; rebuild with `output.library.type = "commonjs2"` |
| React: *"Element type is invalid … got: undefined"* | Same as above, from an older host, or a double-wrapped injected component | Bundle export shape; `normalizeInjectedValue` on the extension side |
| *"… does not exist in dependencies"* | Bundle `require()`s a package absent from the shared runtime | `ui/remote-component.config.js` |
| "Invalid hook call" inside the extension | Two React copies — the bundle shipped its own instead of externalizing it | Extension build config; `resolve.react` in the shared runtime |
| A host action does nothing in the extension | Event literal renamed on one side, or published while the extension was unmounted | `MESHERY_EXTENSION_EVENT`; the `isExtensionOpen()` guard |
| An injected value is `undefined` in the extension | Capability renamed or dropped, or missing from the `useMemo` deps | `buildExtensionInjectProps`; run the gate test |
| A feature works locally, fails in CI | `node_modules/@sistent/sistent` overwritten with a local dist | Re-run `npm ci` and retest |
| Permission-gated UI misbehaves | `ProviderUiAccessControl` constructed from the wrong payload | `ui/utils/disabledComponents.js`; `capabilitiesRegistry` vs `providerCapabilities` |

## File reference

| Path | Role |
| --- | --- |
| `ui/utils/eventBus.ts` | The singleton bus. |
| `ui/components/layout/Navigator/NavigatorExtension.tsx` | Loads the bundle; builds and injects the capability bag; reports load failures. |
| `ui/components/ExtensionSandbox.tsx` | Resolves which bundle URI serves the current path, for each extension point. |
| `ui/pages/extension/[...component].tsx` | Catch-all route for full-page extensions. |
| `ui/remote-component.config.js` | Shared runtime handed to remote bundles. |
| `ui/utils/utils.tsx` | Host-side navigation publishers and `isExtensionOpen()`. |
| `ui/store/index.ts` | Subscribes the Redux store to `DispatchToMesheryStore`. |
| `ui/store/slices/mesheryUi.ts` | Publishes `K8sContextsUpdated`. |
| `ui/utils/can.ts` | Wires `createCanShow` to the bus for permission/capability events. |
| `ui/utils/disabledComponents.js` | `ProviderUiAccessControl` — provider-driven UI gating. |
| `@sistent/sistent` → `mesheryExtensionContract` | The contract itself. Authoritative. |
| `meshery-extensions/ui/src/App.tsx` | Extension-side entry point; destructures `injectProps`. |
| `meshery-extensions/ui/src/globals/mesherySdk.ts` | Extension-side accessors and local typing of the contract. |
