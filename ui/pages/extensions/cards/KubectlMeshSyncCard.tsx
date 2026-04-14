import React from 'react';
import { Button, Grid2, Typography } from '@sistent/sistent';
import {
  UnifiedCardContainer,
  UnifiedDescription,
  UnifiedButtonContainer,
  ResponsiveImage,
} from '../shared';
import { GRID_SIZE } from '../constants';

const KubectlMeshSyncCard = () => {
  const handleLearnMore = (e) => {
    window.open('https://docs.meshery.io/extensions/kubectl-meshsync-snapshot', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="kubectl-meshsync-heading" variant="h5" component="div">
          Kubectl Plugin for MeshSync Snapshot
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/extensions/meshsync.svg"
            alt="MeshSync Logo"
            testId="kubectl-meshsync-image"
          />
          <div>
            Capture cluster state directly from kubectl with simplified networking and access
            requirements. Generate MeshSync snapshots for offline management and visualization in
            Meshery Server, without requiring full Meshery Operator deployment.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="kubectl-meshsync-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export default KubectlMeshSyncCard;
