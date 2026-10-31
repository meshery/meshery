import { useEffect, useState } from 'react';
import {
  Box,
  Button,
  CustomTooltip,
  SearchIcon,
  Typography,
  alpha,
  useTheme,
} from '@sistent/sistent';
import { useCommandPalette } from '@/utils/context/CommandPaletteContextProvider';

type NavigatorWithUserAgentData = Navigator & {
  userAgentData?: {
    platform?: string;
  };
};

export const CommandPaletteTrigger = () => {
  const { openPalette } = useCommandPalette();
  const theme = useTheme();
  const [shortcutLabel, setShortcutLabel] = useState('Ctrl K');

  useEffect(() => {
    try {
      const nav =
        typeof navigator !== 'undefined' ? (navigator as NavigatorWithUserAgentData) : undefined;
      const platform = (nav?.userAgentData?.platform || nav?.platform || '').toLowerCase();
      if (platform.includes('mac') || platform.includes('iphone') || platform.includes('ipad')) {
        setShortcutLabel('⌘K');
      } else {
        setShortcutLabel('Ctrl K');
      }
    } catch {
      setShortcutLabel('Ctrl K');
    }
  }, []);

  const tooltipTitle = shortcutLabel === '⌘K' ? 'Search commands (⌘K)' : 'Search commands (Ctrl K)';

  return (
    <CustomTooltip title={tooltipTitle}>
      <Button
        onClick={openPalette}
        data-testid="command-palette-trigger"
        aria-label="Search commands"
        aria-keyshortcuts="Control+K Meta+K"
        sx={{
          display: 'inline-flex',
          alignItems: 'center',
          gap: 1,
          px: { xs: 1, md: 1.5 },
          py: 0.5,
          minWidth: 'auto',
          borderRadius: 2,
          backgroundColor: alpha(theme.palette.common.white, 0.12),
          border: `1px solid ${alpha(theme.palette.common.white, 0.2)}`,
          color: theme.palette.common.white,
          textTransform: 'none',
          '&:hover': {
            backgroundColor: alpha(theme.palette.common.white, 0.2),
            borderColor: alpha(theme.palette.common.white, 0.35),
          },
        }}
      >
        <SearchIcon
          sx={{
            fontSize: '1.25rem',
            color: theme.palette.common.white,
          }}
        />
        <Box
          component="span"
          sx={{
            display: { xs: 'none', md: 'inline-flex' },
            alignItems: 'center',
            gap: 1,
          }}
        >
          <Typography
            variant="body2"
            sx={{
              color: alpha(theme.palette.common.white, 0.85),
              fontSize: '0.8125rem',
            }}
          >
            Search...
          </Typography>
          <Box
            component="kbd"
            aria-hidden="true"
            sx={{
              display: 'inline-flex',
              alignItems: 'center',
              px: 0.75,
              py: 0.125,
              fontSize: '0.75rem',
              fontWeight: 600,
              borderRadius: 1,
              backgroundColor: alpha(theme.palette.common.white, 0.15),
              color: theme.palette.common.white,
              border: `1px solid ${alpha(theme.palette.common.white, 0.25)}`,
              fontFamily: 'inherit',
            }}
          >
            {shortcutLabel}
          </Box>
        </Box>
      </Button>
    </CustomTooltip>
  );
};

export default CommandPaletteTrigger;
