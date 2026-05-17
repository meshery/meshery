import React, { useState } from 'react';
import AddKubeconfig from '@/assets/gifs/add-kubeconfig.gif';
import SwitchClusters from '@/assets/gifs/switch-clusters.gif';

import {
  Typography,
  Button,
  Box,
  Modal,
  ModalBody,
  ModalFooter,
  ModalButtonPrimary,
  styled,
  useTheme,
  Link,
} from '@sistent/sistent';
import { useRouter } from 'next/router';

const EmptyStateContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  padding: theme.spacing(8, 2),
  textAlign: 'center',
  backgroundColor: theme.palette.background.paper,
  borderRadius: theme.spacing(2),
  minHeight: '400px',
  border: `1px solid ${theme.palette.divider}`,
}));

const EmptyStateTitle = styled(Typography)(({ theme }) => ({
  fontSize: '1.5rem',
  fontWeight: 600,
  color: theme.palette.text.primary,
  marginBottom: theme.spacing(2),
}));

const EmptyStateDescription = styled(Typography)(({ theme }) => ({
  fontSize: '1rem',
  color: theme.palette.text.secondary,
  marginBottom: theme.spacing(4),
  maxWidth: '500px',
  lineHeight: 1.6,
}));

const ActionButton = styled(Button)(({ theme }) => ({
  marginTop: theme.spacing(2),
  padding: theme.spacing(1.5, 4),
  fontWeight: 600,
}));

const StepList = styled('ul')(({ theme }) => ({
  textAlign: 'left',
  margin: theme.spacing(2, 0),
  paddingLeft: theme.spacing(3),
  '& li': {
    marginBottom: theme.spacing(1),
    color: theme.palette.text.secondary,
  },
}));

const MeshSyncEmptyState = () => {
  const [infoModalOpen, setInfoModalOpen] = useState(false);
  const theme = useTheme();
  const router = useRouter();

  const handleOpenInfoModal = () => {
    setInfoModalOpen(true);
  };

  const handleCloseInfoModal = () => {
    setInfoModalOpen(false);
  };

  const handleGoToConnections = () => {
    router.push('/management/connections');
    setInfoModalOpen(false);
  };

  return (
    <>
      <EmptyStateContainer>
        <EmptyStateTitle>No Resources Discovered Yet</EmptyStateTitle>
        <EmptyStateDescription>
          MeshSync automatically discovers and synchronizes your Kubernetes infrastructure. Connect
          a cluster to start seeing your resources here in real-time.
        </EmptyStateDescription>
        <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap', justifyContent: 'center' }}>
          <ActionButton variant="outlined" onClick={handleOpenInfoModal}>
            Learn More
          </ActionButton>
        </Box>
      </EmptyStateContainer>

      <Modal
        open={infoModalOpen}
        closeModal={handleCloseInfoModal}
        title="Getting Started with MeshSync"
        headerIcon={
          <img
            src={
              theme.palette.mode === 'dark'
                ? '/static/img/meshsync-white.svg'
                : '/static/img/meshsync.svg'
            }
            alt="MeshSync"
            style={{ width: '24px', height: '24px' }}
          />
        }
        maxWidth="md"
      >
        <ModalBody>
          <Typography variant="body1" sx={{ mb: 3 }}>
            MeshSync is a custom Kubernetes controller that provides tiered discovery and continual
            synchronization with Meshery Server as to the state of your Kubernetes clusters.
          </Typography>

          <StepList>
            <li>
              <strong>Cluster Detection</strong> - Meshery automatically discovers your available
              clusters. You can also upload your kubeconfig in Connections.
            </li>
            <Box>
              <img src={AddKubeconfig?.src} style={{ maxWidth: '100%' }} />
            </Box>
            <li>
              <strong>Deploy Meshery Operator</strong> - Ensure your deployment mode is Operator and
              Operator is deployed in your cluster
            </li>
            <li>
              <strong>Switching Clusters</strong> - If you want to view resources from a different
              cluster, use the Cluster Switcher in the Nav bar to select it.
            </li>
            <Box>
              <img
                src={SwitchClusters?.src}
                style={{ maxWidth: '100%' }}
                alt="Add Kubernetes connections by importing your kubeconfig"
              />
            </Box>
          </StepList>

          <Typography variant="body2" sx={{ mt: 3, fontStyle: 'italic', color: 'text.secondary' }}>
            Once connected, MeshSync will automatically discover and sync infrastructure from your
            clusters. Resources will appear here in real-time as they are discovered.
          </Typography>

          <Box sx={{ mt: 3, p: 2, borderRadius: 1, display: 'flex', justifyContent: 'center' }}>
            <Typography variant="body2">
              Visit our{' '}
              <Link
                href="https://docs.meshery.io/concepts/architecture/meshsync"
                target="_blank"
                rel="noopener"
              >
                documentation
              </Link>{' '}
              to learn more about MeshSync.
            </Typography>
          </Box>
        </ModalBody>
        <ModalFooter
          helpText={
            'Connect a Kubernetes cluster to populate MeshSync data and start managing your infrastructure with Meshery.\n[Learn more](https://docs.meshery.io/installation/kubernetes) about Kubernetes installation.'
          }
        >
          <ModalButtonPrimary onClick={handleGoToConnections}>Add Cluster</ModalButtonPrimary>
        </ModalFooter>
      </Modal>
    </>
  );
};

export default MeshSyncEmptyState;
