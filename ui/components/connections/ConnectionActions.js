import React from 'react';
import { Box, Button, Typography, Divider } from '@sistent/sistent';
import AddIcon from '@mui/icons-material/Add';
import MesherySettingsEnvButtons from '../MesherySettingsEnvButtons';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import useTestIDsGenerator from '../hooks/useTestIDs';

const ConnectionActions = ({ onOpenRegistrationModal }) => {
  const testIDs = useTestIDsGenerator('connection');

  // const handleOpenKubernetesWizard = () => {
  //   if (onOpenConnectionWizard) {
  //     onOpenConnectionWizard('kubernetes');
  //   }
  // };

  const handleOpenRegistrationModal = () => {
    if (onOpenRegistrationModal) {
      onOpenRegistrationModal();
    }
  };
  // this file is just to make buttons for register connection and add cluster
  return (
    <Box sx={{ display: 'flex', flexDirection: 'row', gap: 1, width: '100%' }}>
      {/* New Connection Wizard Button */}

      {/* {CAN(keys.ADD_CLUSTER.action, keys.ADD_CLUSTER.subject) && (
        <Button
          variant="contained"
          onClick={handleOpenKubernetesWizard}
          startIcon={<AutoAwesomeIcon />}
          sx={{
            width: '100%',
            borderRadius: 1.25,
            padding: '8px 16px',
            background: 'linear-gradient(45deg, #00B39F 30%, #00D3A9 90%)',
            '&:hover': {
              background: 'linear-gradient(45deg, #009688 30%, #00B39F 90%)',
            },
          }}
          data-testid={testIDs('openConnectionWizard')}
        >
          <Typography
            variant="button"
            sx={{
              fontWeight: 600,
              textTransform: 'none',
            }}
          >
            Connection Wizard
          </Typography>
        </Button>
      )} */}

      {/* Register Connection Modal Button */}
      {CAN(keys.ADD_CLUSTER.action, keys.ADD_CLUSTER.subject) && (
        <Button
          variant="outlined"
          onClick={handleOpenRegistrationModal}
          startIcon={<AddIcon />}
          sx={{
            width: '100%',
            borderRadius: 1.25,
            padding: '8px 16px',
          }}
          data-testid={testIDs('openRegistrationModal')}
        >
          <Typography
            variant="button"
            sx={{
              fontWeight: 500,
              textTransform: 'none',
            }}
          >
            Register Connection
          </Typography>
        </Button>
      )}

      {/* Divider with "OR" */}
      <Box sx={{ display: 'flex', alignItems: 'center', my: 1 }}>
        <Divider sx={{ flex: 1 }} />
        <Typography variant="body2" color="textSecondary" sx={{ mx: 2, fontSize: '0.75rem' }}>
          OR
        </Typography>
        <Divider sx={{ flex: 1 }} />
      </Box>

      {/* Legacy Upload Button */}
      <MesherySettingsEnvButtons />
    </Box>
  );
};

export default ConnectionActions;
