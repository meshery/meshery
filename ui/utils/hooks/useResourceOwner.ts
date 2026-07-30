import { useGetProviderCapabilitiesQuery, useGetUserByIdQuery } from '@/rtk-query/user';
import { isRemoteProvider } from '@/utils/provider';

export type OwnerProfile = {
  id?: string;
  firstName?: string;
  lastName?: string;
  avatarUrl?: string;
};

type UseResourceOwnerResult = {
  /** The resolved owner profile, or undefined while unresolvable. */
  owner?: OwnerProfile;
  /**
   * Whether the owner has a Meshery Cloud profile page. The built-in provider's
   * user is a synthetic, local-only identity, so linking its avatar to
   * `cloud.meshery.io/user/<id>` would land on a 404. This is true only once the
   * provider is known to be remote: while the provider-capabilities query is
   * still in flight the type is unresolved, and treating that as "has a cloud
   * profile" would flash a dead, 404-bound link on cold load, so an unresolved
   * provider reports no cloud profile.
   */
  hasCloudProfile: boolean;
};

/**
 * Resolves the owner profile of a content resource (design, filter).
 *
 * Producers that follow the schemas design contract join the owner's profile
 * onto the resource as `user`. That profile is preferred and the by-id lookup is
 * skipped, so no redundant request is issued and the owner resolves even on
 * providers with no per-id user endpoint. Producers that emit only an owner id
 * (filters, and the v1beta1 workspace design page) fall back to the lookup.
 *
 * This is the single source of truth for that precedence, shared by the design
 * card, the filter card and the shared Info modal so the three cannot drift.
 */
export const useResourceOwner = (
  ownerId?: string,
  embeddedOwner?: OwnerProfile | null,
): UseResourceOwnerResult => {
  const { data: fetchedOwner } = useGetUserByIdQuery(ownerId, {
    skip: Boolean(embeddedOwner),
  });
  const { data: providerCapabilities } = useGetProviderCapabilitiesQuery();

  return {
    owner: embeddedOwner ?? fetchedOwner,
    hasCloudProfile: isRemoteProvider(providerCapabilities),
  };
};
