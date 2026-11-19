import React from 'react';
import { CustomTooltip, FavoriteIcon, IconButton, styled } from '@sistent/sistent';
import { RootDiv } from '../../general/style';

export interface PinToggleButtonProps {
  id: string;
  title: string;
  pinned: boolean;
  onToggle: (id: string) => void;
  testId?: string;
}

interface StyledPinButtonProps {
  pinned?: boolean;
}

export const NavRowWrapper = styled(RootDiv)(() => ({
  position: 'relative',
  '&:hover > [data-pin-toggle], &:focus-within > [data-pin-toggle]': {
    opacity: 1,
  },
}));

export const StyledPinButton = styled(IconButton, {
  shouldForwardProp: (prop) => prop !== 'pinned',
})<StyledPinButtonProps>(({ theme, pinned }) => ({
  position: 'absolute',
  right: 12,
  top: '50%',
  transform: 'translateY(-50%)',
  zIndex: 2,
  background: 'transparent',
  border: 'none',
  borderRadius: 4,
  cursor: 'pointer',
  padding: 4,
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  opacity: pinned ? 1 : 0,
  outline: 'none',
  transition: 'opacity 150ms ease-in-out',
  color: pinned ? theme.palette.icon.brand : theme.palette.background.constant.white,
  '& svg': {
    width: '16px',
    height: '16px',
  },
  '& svg path[fill="none"]': {
    stroke: 'none',
  },
  '&:hover': {
    opacity: 1,
    color: pinned ? theme.palette.icon.brand : theme.palette.common.white,
  },
  '&:focus-within': {
    opacity: 1,
  },
  '&:focus-visible': {
    opacity: 1,
    outline: `2px solid ${theme.palette.icon.brand}`,
    outlineOffset: 2,
    color: pinned ? theme.palette.icon.brand : theme.palette.common.white,
  },
  '@media (hover: none)': {
    opacity: 1,
  },
}));

export const PinToggleButton: React.FC<PinToggleButtonProps> = ({
  id,
  title,
  pinned,
  onToggle,
  testId,
}) => {
  const label = pinned ? `Unpin ${title} from Favorites` : `Pin ${title} to Favorites`;

  return (
    <CustomTooltip title={label}>
      <StyledPinButton
        data-pin-toggle="true"
        type="button"
        data-testid={testId || `pin-button-${id}`}
        aria-pressed={pinned}
        aria-label={label}
        onClick={(e) => {
          e.preventDefault();
          e.stopPropagation();
          onToggle(id);
        }}
        pinned={pinned}
      >
        {pinned ? (
          <FavoriteIcon width="16px" height="16px" fill="currentColor" />
        ) : (
          <FavoriteIcon
            width="16px"
            height="16px"
            fill="none"
            stroke="currentColor"
            strokeWidth={2}
          />
        )}
      </StyledPinButton>
    </CustomTooltip>
  );
};

export default PinToggleButton;
