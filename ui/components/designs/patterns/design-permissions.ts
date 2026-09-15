type UserLike = { id?: string } | null | undefined;
type DesignLike = { userId?: string } | null | undefined;

/**
 * canEditDesign is the single source of truth for the "can edit this design"
 * gate, shared by the card/grid view (MesheryPatternCard) and the table view
 * (MesheryPatterns) so the two can't drift apart.
 *
 * Model: a user may edit a design if they hold the global CatalogManagementEditDesign
 * permission, OR they own the design. The truthy `user?.id` guard stops an
 * unauthenticated user (no id) from matching an owner-less design via
 * `undefined === undefined`.
 *
 * Scope: this gates VISIBILITY of the edit affordance. Click-through is
 * additionally gated on CAN(Keys.CatalogManagementEditDesign) at each call site (`disabled={!CAN}`),
 * so an owner without the permission sees the affordance but it stays disabled.
 */
export const canEditDesign = (
  user: UserLike,
  design: DesignLike,
  hasEditPermission: boolean,
): boolean => hasEditPermission || (!!user?.id && user.id === design?.userId);
