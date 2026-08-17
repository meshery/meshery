import React, { useMemo, useState } from 'react';
import {
  Typography,
  Button,
  useTheme,
  Box,
  Chip,
  IconButton,
  CustomTooltip,
  CheckIcon,
  CopyIcon,
  WarningIcon,
  alpha,
} from '@sistent/sistent';
import { useGetSystemVersionQuery } from '@/rtk-query/user';

export interface MesheryVersionCompatibilityNoticeProps {
  currentVersion?: string;
  requiredVersion?: string;
  componentName?: string;
  upgradeCommand?: string;
}

/**
 * MesheryVersionCompatibilityNotice Component
 *
 * Renders a rich, theme-aware notice when a version mismatch or compatibility
 * issue occurs between Meshery components (Server, UI, CLI, Adapters).
 */
const MesheryVersionCompatibilityNotice: React.FC<MesheryVersionCompatibilityNoticeProps> = ({
  currentVersion = 'v0.7.0',
  requiredVersion = 'v0.7.1+',
  componentName = 'Meshery Server',
  upgradeCommand = 'mesheryctl system update',
}) => {
  const theme = useTheme();
  const [copied, setCopied] = useState(false);
  const { data: systemVersionData } = useGetSystemVersionQuery();

  const resolvedCurrentVersion = useMemo(() => {
    if (systemVersionData?.build) {
      return systemVersionData.build;
    }

    return currentVersion;
  }, [currentVersion, systemVersionData?.build]);

  const handleCopy = async () => {
    if (typeof navigator === 'undefined' || !navigator.clipboard) {
      return;
    }

    try {
      await navigator.clipboard.writeText(upgradeCommand);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      setCopied(false);
    }
  };

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        padding: theme.spacing(3),
        borderRadius: '12px',
        backgroundColor: alpha(
          theme.palette.common.white,
          theme.palette.mode === 'dark' ? 0.05 : 0.02,
        ),
        border: `1px solid ${
          theme.palette.mode === 'dark'
            ? alpha(theme.palette.common.white, 0.1)
            : alpha(theme.palette.common.black, 0.08)
        }`,
        maxWidth: '560px',
        margin: '1.5rem auto',
        textAlign: 'center',
      }}
    >
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          gap: 1,
          marginBottom: theme.spacing(1.5),
        }}
      >
        <WarningIcon
          sx={{
            color: theme.palette.warning?.main,
            fontSize: '2rem',
          }}
        />
        <Typography
          variant="h5"
          component="h5"
          fontWeight={600}
          sx={{ color: theme.palette.text.default }}
        >
          Version Compatibility Notice
        </Typography>
      </Box>

      <Typography variant="body1" component="p" sx={{ color: theme.palette.text.secondary, mb: 2 }}>
        The currently running <strong>{componentName}</strong> version may not be fully compatible
        with this interface.
      </Typography>

      <Box
        sx={{
          display: 'flex',
          gap: 1.5,
          marginBottom: theme.spacing(2.5),
          alignItems: 'center',
          flexWrap: 'wrap',
          justifyContent: 'center',
        }}
      >
        <Chip
          label={`Current: ${resolvedCurrentVersion}`}
          color="warning"
          variant="outlined"
          size="medium"
        />
        <Chip
          label={`Required: ${requiredVersion}`}
          color="success"
          variant="outlined"
          size="medium"
        />
      </Box>

      <Box
        sx={{
          width: '100%',
          backgroundColor:
            theme.palette.mode === 'dark' ? theme.palette.grey[900] : theme.palette.grey[100],
          borderRadius: '8px',
          padding: '8px 14px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          marginBottom: theme.spacing(2.5),
          fontFamily: 'monospace',
          border: `1px solid ${theme.palette.divider}`,
        }}
      >
        <Typography
          variant="body2"
          component="code"
          sx={{
            fontFamily: 'monospace',
            color: theme.palette.text.primary,
          }}
        >
          {upgradeCommand}
        </Typography>
        <CustomTooltip title={copied ? 'Copied!' : 'Copy command'}>
          <IconButton size="small" onClick={handleCopy} aria-label="Copy upgrade command">
            {copied ? (
              <CheckIcon style={{ fontSize: 18, color: theme.palette.success.main }} />
            ) : (
              <CopyIcon style={{ fontSize: 18 }} />
            )}
          </IconButton>
        </CustomTooltip>
      </Box>

      <Box sx={{ display: 'flex', gap: 2 }}>
        <Button
          variant="contained"
          color="primary"
          href="https://docs.meshery.io/guides/mesheryctl#upgrading-meshery"
          target="_blank"
          rel="noopener noreferrer"
        >
          Upgrade Guide
        </Button>
        <Button
          variant="outlined"
          color="secondary"
          href="https://meshery.io/community#community-forums"
          target="_blank"
          rel="noopener noreferrer"
        >
          Discussion Forum
        </Button>
      </Box>
    </Box>
  );
};

export default MesheryVersionCompatibilityNotice;
