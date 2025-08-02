/**
 * Reusable chip hover styles utility
 * Provides consistent hover behavior for chips with overflow tooltip functionality
 */

export const createChipHoverStyles = (theme) => ({
  '&:hover': {
    zIndex: 1000,
    '& .MuiChip-label': {
      overflow: 'visible',
      whiteSpace: 'normal',
      position: 'absolute',
      left: 0,
      top: 0,
      width: 'max-content',
      maxWidth: 'min(300px, 80vw)',
      backgroundColor: theme.palette.background.paper || '#fff',
      border: `1px solid ${theme.palette.divider || '#e0e0e0'}`,
      borderRadius: '4px',
      padding: '8px 12px',
      boxShadow: theme.shadows[4] || '0px 2px 8px rgba(0,0,0,0.15)',
      zIndex: 1001,
      // Improved viewport boundary containment
      transform: 'translateX(clamp(-100%, 0px, calc(100vw - 100% - 20px)))',
      // Modern word breaking
      overflowWrap: 'break-word',
      hyphens: 'auto',
      // Ensure proper positioning
      transition: 'all 0.3s ease',
    },
  },
});

/**
 * Enhanced chip hover styles with better positioning logic
 * Uses CSS clamp() for better viewport boundary handling
 */
export const createEnhancedChipHoverStyles = (theme) => ({
  overflow: 'visible',
  position: 'relative',
  '& .MuiChip-label': {
    overflow: 'hidden',
    whiteSpace: 'nowrap',
    textOverflow: 'ellipsis',
    transition: 'all 0.3s ease',
  },
  '&:hover': {
    zIndex: 1000,
    '& .MuiChip-label': {
      overflow: 'visible',
      whiteSpace: 'normal',
      position: 'absolute',
      left: 0,
      top: 0,
      width: 'max-content',
      maxWidth: 'min(300px, 80vw)',
      backgroundColor: theme.palette.background.paper || '#fff',
      border: `1px solid ${theme.palette.divider || '#e0e0e0'}`,
      borderRadius: '4px',
      padding: '8px 12px',
      boxShadow: theme.shadows[4] || '0px 2px 8px rgba(0,0,0,0.15)',
      zIndex: 1001,
      // Better viewport boundary containment using clamp()
      transform: 'translateX(clamp(-100%, 0px, calc(100vw - 100% - 20px)))',
      // Modern word breaking properties
      overflowWrap: 'break-word',
      hyphens: 'auto',
      transition: 'all 0.3s ease',
    },
  },
});
