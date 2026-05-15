import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../css/icons.styles', () => ({
  iconMedium: {},
}));

vi.mock('@sistent/sistent', () => ({
  OutlinedVisibilityOnIcon: () => <svg data-testid="visibility-on" />,
  OutlinedVisibilityOffIcon: () => <svg data-testid="visibility-off" />,
  Button: ({ children, onClick, ...rest }: any) => (
    <button onClick={onClick} {...rest}>
      {children}
    </button>
  ),
}));

import CatalogFilter from './CatalogFilter';

describe('CatalogFilter', () => {
  it('renders nothing when hideCatalog is true', () => {
    render(
      <CatalogFilter
        catalogVisibility={true}
        handleCatalogVisibility={() => {}}
        hideCatalog={true}
        classes={{ btnText: 'btn-text' }}
      />,
    );
    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });

  it('renders the catalog visibility-on icon when catalog is visible', () => {
    render(
      <CatalogFilter
        catalogVisibility={true}
        handleCatalogVisibility={() => {}}
        classes={{ btnText: 'btn-text' }}
      />,
    );
    expect(screen.getByTestId('visibility-on')).toBeInTheDocument();
    expect(screen.queryByTestId('visibility-off')).not.toBeInTheDocument();
  });

  it('renders the visibility-off icon when catalog is hidden', () => {
    render(
      <CatalogFilter
        catalogVisibility={false}
        handleCatalogVisibility={() => {}}
        classes={{ btnText: 'btn-text' }}
      />,
    );
    expect(screen.getByTestId('visibility-off')).toBeInTheDocument();
  });

  it('invokes handleCatalogVisibility on click', async () => {
    const user = userEvent.setup();
    const toggle = vi.fn();
    render(
      <CatalogFilter
        catalogVisibility={true}
        handleCatalogVisibility={toggle}
        classes={{ btnText: 'btn-text' }}
      />,
    );
    await user.click(screen.getByRole('button'));
    expect(toggle).toHaveBeenCalledTimes(1);
  });

  it('uses the provided className for the button text', () => {
    const { container } = render(
      <CatalogFilter
        catalogVisibility={true}
        handleCatalogVisibility={() => {}}
        classes={{ btnText: 'btn-text-class' }}
      />,
    );
    expect(container.querySelector('.btn-text-class')).toBeInTheDocument();
  });
});
