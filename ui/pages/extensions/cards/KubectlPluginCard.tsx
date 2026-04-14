import React from 'react';
import { Button, Grid2, Typography } from '@sistent/sistent';
import {
  UnifiedCardContainer,
  UnifiedDescription,
  UnifiedButtonContainer,
  ResponsiveImage,
} from '../shared';
import { GRID_SIZE } from '../constants';

const KubectlPluginCard = () => {
  const handleLearnMore = (e) => {
    window.open('https://docs.meshery.io/extensions/kubectl-kanvas-snapshot', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="kubectl-plugin-heading" variant="h5" component="div">
          Kubectl Plugin for Kanvas Snapshot
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/extensions/kubectl.svg"
            alt="Kubectl Plugin Logo"
            testId="kubectl-plugin-image"
          />
          <div>
            Generate visual snapshots of your Kubernetes manifests directly from kubectl. Cluster
            configurations and workflows with Kanvas Snapshots. Receive snapshots via email or get
            instant terminal URL display.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="kubectl-plugin-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export default KubectlPluginCard;
