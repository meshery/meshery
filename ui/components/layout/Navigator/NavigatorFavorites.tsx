import React from 'react';
import Link from 'next/link';
import { CustomTooltip, ListItemIcon, styled, useTheme, Zoom } from '@sistent/sistent';
import {
  MainListIcon,
  NavigatorLink,
  RootDiv,
  SecondaryDivider,
  SideBarListItem,
  SideBarText,
} from '../../general/style';
import PinIcon from '@/assets/icons/PinIcon';
import PinToggleButton from './PinToggleButton';

const FavoritesHeader = styled('div')(({ theme }) => ({
  padding: '8px 24px 4px 24px',
  fontSize: '0.75rem',
  fontWeight: 600,
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  color: theme.palette.background.constant.white,
  opacity: 0.6,
  userSelect: 'none',
}));

const FavoritesEmptyCollapsed = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  height: 36,
  color: theme.palette.background.constant.white,
  opacity: 0.35,
}));

const FavoritesEmptyExpanded = styled('div')(({ theme }) => ({
  padding: '8px 24px',
  fontSize: '0.8125rem',
  color: theme.palette.background.constant.white,
  opacity: 0.5,
  fontStyle: 'italic',
}));

const FavoriteItemWrapper = styled(RootDiv)(() => ({
  position: 'relative',
  '&:hover button, &:focus-within button': {
    opacity: 1,
  },
}));

const CollapsedFavListItemIcon = styled(ListItemIcon)(() => ({
  marginLeft: '20%',
  marginBottom: '0.4rem',
}));

export interface NavigatorItem {
  id: string;
  title: string;
  icon?: React.ReactNode;
  hovericon?: React.ReactNode;
  href?: string;
  link?: boolean;
  show?: boolean;
  permissionKey?: string;
  [key: string]: any;
}

export interface NavigatorFavoritesProps {
  resolvedPinnedItems: NavigatorItem[];
  isDrawerCollapsed: boolean;
  hoveredId: string | null;
  setHoveredId: (id: string | null) => void;
  currentPath: string;
  onTogglePin: (id: string) => void;
}

export const NavigatorFavorites: React.FC<NavigatorFavoritesProps> = ({
  resolvedPinnedItems,
  isDrawerCollapsed,
  hoveredId,
  setHoveredId,
  currentPath,
  onTogglePin,
}) => {
  const theme = useTheme();

  return (
    <div data-testid="favorites-section">
      {!isDrawerCollapsed && (
        <FavoritesHeader data-testid="favorites-header">Favorites</FavoritesHeader>
      )}
      {resolvedPinnedItems.length === 0 ? (
        isDrawerCollapsed ? (
          <CustomTooltip title="No pinned favorites" placement="right" TransitionComponent={Zoom}>
            <FavoritesEmptyCollapsed data-testid="favorites-empty-state-collapsed">
              <PinIcon
                pinned={false}
                width="16px"
                height="16px"
                fill={theme.palette.background.constant.white}
              />
            </FavoritesEmptyCollapsed>
          </CustomTooltip>
        ) : (
          <FavoritesEmptyExpanded data-testid="favorites-empty-state">
            No pinned favorites
          </FavoritesEmptyExpanded>
        )
      ) : (
        resolvedPinnedItems.map((item) => {
          const isActive = currentPath === item.href;
          const updatedFavIcon =
            React.isValidElement(item.icon) && typeof item.icon.type !== 'string'
              ? React.cloneElement(item.icon as React.ReactElement<any>, {
                  fill: isActive ? theme.palette.icon.brand : theme.palette.common.white,
                })
              : item.icon;

          return (
            <FavoriteItemWrapper key={`favorite-${item.id}`}>
              <SideBarListItem
                dense
                link={true}
                isActive={isActive}
                onMouseOver={() => (isDrawerCollapsed ? setHoveredId(`fav-${item.id}`) : null)}
                onMouseLeave={() => (isDrawerCollapsed ? setHoveredId(null) : null)}
                permissionKey={item.permissionKey}
                permissionAction="showShield"
                {...(item.href ? { component: Link, href: item.href } : {})}
              >
                <NavigatorLink data-testid={`favorite-${item.id}`}>
                  <CustomTooltip
                    title={item.title}
                    placement="right"
                    disableFocusListener={!isDrawerCollapsed}
                    disableHoverListener={true}
                    disableTouchListener={!isDrawerCollapsed}
                    TransitionComponent={Zoom}
                  >
                    {isDrawerCollapsed && hoveredId === `fav-${item.id}` ? (
                      <div>
                        <CustomTooltip
                          title={item.title}
                          placement="right"
                          TransitionComponent={Zoom}
                        >
                          <CollapsedFavListItemIcon>
                            {item.hovericon ?? item.icon}
                          </CollapsedFavListItemIcon>
                        </CustomTooltip>
                      </div>
                    ) : (
                      <MainListIcon>{updatedFavIcon}</MainListIcon>
                    )}
                  </CustomTooltip>
                  <SideBarText drawerCollapsed={isDrawerCollapsed}>{item.title}</SideBarText>
                </NavigatorLink>
              </SideBarListItem>
              {!isDrawerCollapsed && (
                <PinToggleButton
                  id={item.id}
                  title={item.title}
                  pinned={true}
                  onToggle={onTogglePin}
                  testId={`unpin-button-${item.id}`}
                />
              )}
            </FavoriteItemWrapper>
          );
        })
      )}
      <SecondaryDivider />
    </div>
  );
};

export default NavigatorFavorites;
