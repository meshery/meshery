import React, { useContext, useMemo } from 'react';
import {
  createUseRemoteComponent,
  getDependencies,
  createRequires,
} from '@paciolan/remote-component';
import MesheryPerformanceComponent from '../../performance';
import PatternServiceFormCore from '../../meshery-mesh-interface/PatternServiceFormCore';
import InfoModal from '../../shared/Modal/Information/InfoModal';
import _PromptComponent from '../../general/PromptComponent';
import { ProviderUiAccessControl } from '../../../utils/disabledComponents';
import { useNotification } from '../../../utils/hooks/useNotification';
import Modal from '../../shared/Modal/Modal';
import ExportDesignModal from '../../designs/export/ExportDesignModal';
import { FormatStructuredData } from '../../data-formatter';
import { useFilterK8sContexts } from '@/utils/hooks/useKubernetesHook';
import { useDynamicComponent } from '@/utils/context/dynamicContext';
import { ValidateDesign } from '../../designs/lifecycle/ValidateDesign';
import { DryRunDesign } from '../../designs/lifecycle/DryRun';
import { DeployStepper, UnDeployStepper } from '../../designs/lifecycle/DeployStepper';
import { designValidationMachine } from 'machines/validator/designValidator';
import CAN from '@/utils/can';
import { mesheryEventBus } from '@/utils/eventBus';
import { ThemeTogglerCore } from '@/theme/hooks';
import RJSFForm from '../../meshery-mesh-interface/PatternService/RJSF';
import { DynamicFullScreenLoader } from '../../shared/LoadingState/DynamicFullscreenLoader';
import Troubleshoot from '../../TroubleshootingComponent';
import TypingFilter from '@/components/shared/FormFields/typing-filter';
import CreateModelModal from '../../registry/CreateModelModal';
import ImportModelModal from '../../registry/ImportModelModal';
import { ViewInfoModal } from '../../workspaces/ViewInfoModal';
import { selectK8sConfig, selectSelectedK8sClusters } from '@/store/slices/mesheryUi';
import { useSelector } from 'react-redux';
import { store } from '../../../store';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import { RelationshipEvaluationTraceFormatter } from '../NotificationCenter/formatters/relationship_evaluation';
import { useRegistryModal } from '@/utils/hooks/useRegistryModal';
import { MESHERY_EXTENSION_CONTRACT_VERSION } from '@sistent/sistent';

const requires = createRequires(getDependencies);
const useRemoteComponent = createUseRemoteComponent({ requires });

type NavigatorExtensionProps = {
  url: string;
};

const extensionExposedMesheryStore = {
  selectedK8sClusters: {
    get: () => selectSelectedK8sClusters(store.getState()),
  },
  k8sConfig: {
    get: () => selectK8sConfig(store.getState()),
  },
};

/**
 * Inert stand-in for the removed `subscribeConfiguration` GraphQL subscription.
 * Matches the old call shape — `(onNext, variables) => ({ dispose })` — so an
 * extension bundle published against the previous host contract still mounts and
 * unmounts cleanly. It never emits: `onNext` is intentionally never invoked.
 */
const noopSubscription = () => ({ dispose: () => {} });

function PerformanceTestComponent(props: React.ComponentProps<typeof MesheryPerformanceComponent>) {
  return (
    <ProviderStoreWrapper>
      <MesheryPerformanceComponent {...props} />
    </ProviderStoreWrapper>
  );
}

/**
 * Everything the injected capability bag needs from React state or context. The
 * rest of the bag is module-scoped, so a caller only has to supply these.
 */
export type ExtensionInjectPropsDeps = {
  providerCapabilities: unknown;
  selectedK8sContexts: unknown;
  currentOrganization: unknown;
  openWorkspaceModal: unknown;
  openRegistryModal: unknown;
  setCurrentLoadedResourceInOrgWorkspaceSession: unknown;
};

/**
 * Builds the capability bag handed to an extension bundle as `injectProps`.
 *
 * Every key here is part of the host <-> extension contract declared in sistent's
 * `mesheryExtensionContract`: extensions read these off the injected bag, so
 * renaming or dropping one resolves to `undefined` at the extension's use site
 * rather than failing anywhere near here. Kept as a pure factory so a unit test
 * can assert the bag against the contract and catch such a rename before merge.
 */
export const buildExtensionInjectProps = ({
  providerCapabilities,
  selectedK8sContexts,
  currentOrganization,
  openWorkspaceModal,
  openRegistryModal,
  setCurrentLoadedResourceInOrgWorkspaceSession,
}: ExtensionInjectPropsDeps) => ({
  // Lets an extension bundle detect that it was built against a different
  // contract revision than the host it was loaded into, which no build-time
  // check can catch: bundles are published artifacts loaded by whatever host
  // version happens to be deployed.
  contractVersion: MESHERY_EXTENSION_CONTRACT_VERSION,
  PatternServiceFormCore,
  RelationshipEvaluationResponseFormatter: RelationshipEvaluationTraceFormatter,
  MesheryPerformanceComponent: PerformanceTestComponent,
  selectedK8sContexts,
  // Meshery Server no longer exposes any GraphQL subscription. `resolver` is
  // kept only for backward compatibility with already-published extension
  // bundles: they call `resolver.subscription.ConfigurationSubscription(cb, vars)`
  // and later `.dispose()` on the result. Handing them `undefined` would throw
  // at mount, so we hand them an inert subscription that never emits. New
  // extensions should read designs/filters from the REST API instead.
  resolver: {
    query: {},
    mutation: {},
    subscription: {
      ConfigurationSubscription: noopSubscription,
    },
  },
  InfoModal,
  ViewInfoModal,
  ExportModal: ExportDesignModal,
  GenericRJSFModal: Modal,
  _PromptComponent,
  providerCapabilities,
  ProviderUiAccessControlClass: ProviderUiAccessControl,
  // Backward-compatible alias for already-published remote extensions.
  CapabilitiesRegistryClass: ProviderUiAccessControl,
  TypingFilter,
  useNotificationHook: useNotification,
  StructuredDataFormatter: FormatStructuredData,
  CreateModelModal,
  ImportModelModal,
  ValidateDesign,
  DryRunDesign,
  DeployStepper,
  UnDeployStepper,
  designValidationMachine,
  mesheryEventBus,
  ThemeTogglerCore,
  RJSForm: RJSFForm,
  hooks: {
    CAN,
    useFilterK8sContexts,
    useDynamicComponent,
  },
  mesheryStore: extensionExposedMesheryStore,
  currentOrganization,
  openWorkspaceModal,
  openRegistryModal,
  SetCurrentLoadedResourceInOrgWorkspaceSession: setCurrentLoadedResourceInOrgWorkspaceSession,
});

/**
 * The `$$typeof` tags of the object-shaped values React accepts as a *component*
 * type. Deliberately an allow-list rather than a bare `'$$typeof' in component`
 * check: many React-internal values carry `$$typeof` without being components —
 * most notably an already-rendered element (`<Foo />`), which a bundle can export
 * by accident. Accepting those would let the very failure this guard exists to
 * diagnose through to React's opaque "Element type is invalid" error.
 */
const RENDERABLE_COMPONENT_TAGS: ReadonlySet<symbol> = new Set([
  Symbol.for('react.forward_ref'),
  Symbol.for('react.memo'),
  Symbol.for('react.lazy'),
]);

/**
 * Whether a value can legally be passed to `React.createElement` as an element
 * type. Function components (and classes) are plain functions; `forwardRef`,
 * `memo` and `lazy` results are objects tagged with `$$typeof`.
 */
function isRenderableComponent(
  component: unknown,
): component is React.ComponentType<{ injectProps: unknown }> {
  if (typeof component === 'function') {
    return true;
  }
  return (
    typeof component === 'object' &&
    component !== null &&
    RENDERABLE_COMPONENT_TAGS.has((component as { $$typeof?: symbol }).$$typeof as symbol)
  );
}

/**
 * Best-effort description of what a bundle exported instead of a component, so
 * the diagnostic names the actual value rather than just "not a component".
 */
function describeExportKind(exported: unknown): string {
  if (exported && typeof exported === 'object') {
    const tag = (exported as { $$typeof?: unknown }).$$typeof;
    if (typeof tag === 'symbol') {
      return `a React internal value tagged ${String(tag)}`;
    }
  }
  return `a value of type "${typeof exported}"`;
}

/**
 * Explains why a loaded bundle yielded nothing renderable.
 *
 * The two causes have different remedies and must not share a message: a missing
 * default export is a bundler configuration problem, while a default export that
 * is present but is not a component is a defect in the extension's own source.
 * Telling an author to rebuild with `commonjs2` when they actually exported
 * `<Foo />` instead of `Foo` sends them after the wrong fix entirely.
 */
function describeUnrenderableExport(url: string, exported: unknown): string {
  if (exported === undefined || exported === null) {
    return (
      `The extension at ${url} loaded but exposed no CommonJS default export ` +
      '(module.exports.default), so there is no component to render. The bundle was ' +
      'most likely built without output.library.type = "commonjs2". Rebuild and ' +
      'republish the extension, then reload this page.'
    );
  }
  return (
    `The extension at ${url} exposed a CommonJS default export, but it is not a React ` +
    `component (received ${describeExportKind(exported)}). Export the component itself ` +
    'rather than a rendered element or other value, then rebuild and republish the extension.'
  );
}

function NavigatorExtensionError({ error }: { error: unknown }) {
  return (
    <div role="alert">
      <h2>Uh-oh!😔 Please pardon our mesh.</h2>
      <div
        style={{
          backgroundColor: '#1E2117',
          color: '#FFFFFF',
          padding: '.85rem',
          borderRadius: '.2rem',
        }}
      >
        <code>{String(error)}</code>
      </div>
      <div style={{ marginTop: '1rem' }}>
        <Troubleshoot showDesignerButton={false} />
      </div>
    </div>
  );
}

function NavigatorExtension({ url }: NavigatorExtensionProps) {
  const {
    providerCapabilities,
    selectedK8sContexts,
    organization: currentOrganization,
  } = useSelector((state) => state.ui);
  // `useRemoteComponent` fetches the bundle, evaluates it as a CommonJS module,
  // and returns `module.exports.default` as the component. The remote bundle MUST
  // therefore export the component as a CommonJS default
  // (`module.exports = { default: Component, __esModule: true }`). NOTE: the hook
  // does NOT throw when `.default` is missing — it returns `RemoteComponent ===
  // undefined` with `err === undefined`, so a mis-built bundle fails silently
  // (no loader error; React then throws "Element type is invalid … got: undefined").
  // If an extension renders as undefined with no error, inspect the bundle's export
  // shape, not this code. See meshery-extensions/docs/troubleshoot.md.
  const [loading, err, RemoteComponent] = useRemoteComponent(url);
  const { openModalWithDefault, onLoadResource } = useContext(WorkspaceModalContext);
  const registryModal = useRegistryModal();

  const injectProps = useMemo(
    () =>
      buildExtensionInjectProps({
        providerCapabilities,
        selectedK8sContexts,
        currentOrganization,
        openWorkspaceModal: openModalWithDefault,
        openRegistryModal: registryModal,
        setCurrentLoadedResourceInOrgWorkspaceSession: onLoadResource,
      }),
    [
      providerCapabilities,
      currentOrganization,
      onLoadResource,
      openModalWithDefault,
      registryModal,
      selectedK8sContexts,
    ],
  );

  if (err != null) {
    return <NavigatorExtensionError error={err} />;
  }

  const hasComponent = isRenderableComponent(RemoteComponent);

  // A bundle without a CommonJS default export leaves both `err` and
  // `RemoteComponent` undefined, so this is the last point at which the actual
  // cause is still known. Without it React reports only "Element type is
  // invalid ... got: undefined" from inside the render tree.
  if (!loading && !hasComponent) {
    return <NavigatorExtensionError error={describeUnrenderableExport(url, RemoteComponent)} />;
  }

  return (
    <DynamicFullScreenLoader isLoading={loading}>
      {hasComponent ? <RemoteComponent injectProps={injectProps} /> : null}
    </DynamicFullScreenLoader>
  );
}

export default NavigatorExtension;
