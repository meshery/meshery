/**
 * Shared header-icon wrapper for design-domain modals.
 *
 * The Sistent modal header expects a small icon (~24px) rendered to the left
 * of the title. Both `ImportDesignModal` and `PublishDesignModal` (plus future
 * design modals) want the same square slot for the Pattern SVG, so the
 * styled wrapper lives here to keep the two modal files focused.
 */
import { FC, ReactNode } from 'react';
import { styled, useTheme } from '@/theme';
import Pattern from '../../public/static/img/drawer-icons/pattern_svg';

const DesignHeaderIconSlot = styled('span')({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  height: '24px',
  width: '24px',
  fontSize: '1.45rem',
});

interface DesignModalHeaderIconProps {
  /** Optional override; defaults to the Pattern SVG used by all design modals. */
  children?: ReactNode;
}

export const DesignModalHeaderIcon: FC<DesignModalHeaderIconProps> = ({ children }) => {
  const theme = useTheme();
  return (
    <DesignHeaderIconSlot>
      {children ?? <Pattern fill={theme.palette.common.white} />}
    </DesignHeaderIconSlot>
  );
};
