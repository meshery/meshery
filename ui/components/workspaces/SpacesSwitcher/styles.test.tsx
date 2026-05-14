import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (_Component: any, _opts?: any) => (_factory?: any) => {
    const Styled = ({ children, ...props }: any) => <div {...props}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    Drawer: ({ children }: any) => <div>{children}</div>,
    List: ({ children }: any) => <ul>{children}</ul>,
    ListItem: ({ children }: any) => <li>{children}</li>,
    ListItemText: ({ children }: any) => <span>{children}</span>,
    ListItemIcon: ({ children }: any) => <span>{children}</span>,
    DARK_BLUE_GRAY: '#243a44',
    Box: ({ children }: any) => <div>{children}</div>,
    Avatar: ({ children }: any) => <div>{children}</div>,
    Button: ({ children }: any) => <button>{children}</button>,
    Checkbox: () => <input type="checkbox" />,
    Typography: ({ children }: any) => <span>{children}</span>,
    Modal: ({ children }: any) => <div>{children}</div>,
  };
});

import * as styles from './styles';

describe('SpacesSwitcher/styles', () => {
  it('exports the expected styled components', () => {
    expect(styles.DrawerHeader).toBeDefined();
    expect(styles.StyledDrawer).toBeDefined();
    expect(styles.StyledMainContent).toBeDefined();
    expect(styles.DesignList).toBeDefined();
    expect(styles.LoadingContainer).toBeDefined();
    expect(styles.GhostContainer).toBeDefined();
    expect(styles.GhostImage).toBeDefined();
    expect(styles.GhostText).toBeDefined();
    expect(styles.StyledListItem).toBeDefined();
    expect(styles.StyledUserDetailsContainer).toBeDefined();
    expect(styles.StyledAvatarContainer).toBeDefined();
    expect(styles.StyledMainMenuComponent).toBeDefined();
    expect(styles.StyledListItemText).toBeDefined();
    expect(styles.StyledListIcon).toBeDefined();
    expect(styles.StyledUpdatedText).toBeDefined();
    expect(styles.StyledSmallAvatarContainer).toBeDefined();
    expect(styles.StyledSmallAvatar).toBeDefined();
    expect(styles.StyledResponsiveButton).toBeDefined();
    expect(styles.StyledMuiDoubleCheckbox).toBeDefined();
    expect(styles.StyledTypography).toBeDefined();
    expect(styles.StyledModal).toBeDefined();
  });

  it('renders the styled components without error', () => {
    const { container } = render(
      <styles.StyledMainContent>
        <styles.DesignList>
          <styles.StyledListItem>hello</styles.StyledListItem>
        </styles.DesignList>
      </styles.StyledMainContent>,
    );
    expect(container).toBeInTheDocument();
  });
});
