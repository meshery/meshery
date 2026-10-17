import React from 'react';
import { styled } from '@sistent/sistent';
import PinIcon from '@/assets/icons/PinIcon';

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

export const StyledPinButton = styled('button', {
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
  return (
    <StyledPinButton
      type="button"
      data-testid={testId || `pin-button-${id}`}
      aria-pressed={pinned}
      aria-label={pinned ? `Unpin ${title} from Favorites` : `Pin ${title} to Favorites`}
      onClick={(e) => {
        e.preventDefault();
        e.stopPropagation();
        onToggle(id);
      }}
      pinned={pinned}
    >
      <PinIcon pinned={pinned} width="16px" height="16px" fill="currentColor" />
    </StyledPinButton>
  );
};

export default PinToggleButton;
