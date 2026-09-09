import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Keys } from '@meshery/schemas/permissions';

/**
 * Access-gate coverage for the design configurator route.
 *
 * The configurator gated its individual save/update/delete controls but
 * rendered the model browser and the design JSON in full for any authenticated
 * session, including an organization member holding zero permission keys - the
 * one page missed by 69db7bd7de64 ("Enforce permission guards on unguarded
 * pages"). See meshery/meshery#21492.
 *
 * The deny path is what these tests exist to pin: an allow-only test would have
 * passed against the ungated page too. The allow path is asserted alongside it
 * so a regression to a blanket deny is caught as well.
 *
 * `DesignConfigurator` owns cytoscape, xstate and a live meshmodel fetch, so it
 * is stubbed here (as its own suite records, it is not unit-renderable); this
 * suite is scoped to the branch the page itself owns.
 */

const { useHasPermission } = vi.hoisted(() => ({ useHasPermission: vi.fn() }));

vi.mock('@sistent/sistent', () => ({ useHasPermission }));

vi.mock('@/components/general/MesheryPage', () => ({
  MesheryPage: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="meshery-page">{children}</div>
  ),
}));

vi.mock('@/components/designs/configurator/MeshModel', () => ({
  default: () => <div data-testid="design-configurator" />,
}));

vi.mock('@/components/general/error-404', () => ({
  default: ({ permissionKey }: { permissionKey: { id: string } }) => (
    <div data-testid="default-error">{permissionKey?.id}</div>
  ),
}));

import DesignConfiguratorPage from '../../../../pages/configuration/designs/configurator';

describe('DesignConfiguratorPage access gate', () => {
  beforeEach(() => {
    useHasPermission.mockReset();
  });

  it('replaces the configurator with the permission-denied page when the key is absent', () => {
    useHasPermission.mockReturnValue(false);

    render(<DesignConfiguratorPage />);

    expect(screen.queryByTestId('design-configurator')).not.toBeInTheDocument();
    expect(screen.getByTestId('default-error')).toHaveTextContent(
      Keys.CatalogManagementViewDesigns.id,
    );
  });

  it('renders the configurator when the key is held', () => {
    useHasPermission.mockReturnValue(true);

    render(<DesignConfiguratorPage />);

    expect(screen.getByTestId('design-configurator')).toBeInTheDocument();
    expect(screen.queryByTestId('default-error')).not.toBeInTheDocument();
  });

  it('gates on View Designs - the same key as the designs list and navigator entry', () => {
    useHasPermission.mockReturnValue(true);

    render(<DesignConfiguratorPage />);

    expect(useHasPermission).toHaveBeenCalledWith(Keys.CatalogManagementViewDesigns);
  });
});
