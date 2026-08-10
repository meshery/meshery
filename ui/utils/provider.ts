/**
 * Provider-type helpers.
 *
 * Meshery runs against either the built-in local provider or a remote provider
 * (Meshery Cloud). The provider's type is reported by the provider-capabilities
 * response as `providerType` ("local" | "remote"; see
 * server/models/providers.go). These helpers centralize the check so callers do
 * not hand-roll `providerCapabilities?.providerType === 'local'` inline.
 *
 * Note on the local provider: since the schemas v1beta3 account-consolidation
 * cutover, the local provider's user is a synthetic UUID (LocalProviderUserID)
 * and no longer carries the legacy `userId: "meshery"` string. Code that used to
 * detect the local provider via `user.userId === 'meshery'` must use
 * `isLocalProvider(providerCapabilities)` instead.
 *
 * Schema note: `providerType` is sourced from the provider-capabilities response,
 * which is NOT yet covered by `@meshery/schemas` - it is explicitly "pending the
 * provider-capabilities schema tracked separately in the identifier-uniformity
 * program" (see @meshery/schemas cloudApi). This helper therefore takes a
 * loosely-typed capabilities object as a stopgap; once that schema lands, the
 * generated `mesheryApi` type should back the argument and the local
 * `getProviderCapabilities` endpoint (ui/rtk-query/user.ts) should be retired.
 */
type ProviderType = 'local' | 'remote';

type ProviderCapabilitiesLike = { providerType?: ProviderType } | null | undefined;

export const isLocalProvider = (capabilities?: ProviderCapabilitiesLike): boolean =>
  capabilities?.providerType === 'local';

export const isRemoteProvider = (capabilities?: ProviderCapabilitiesLike): boolean =>
  capabilities?.providerType === 'remote';
