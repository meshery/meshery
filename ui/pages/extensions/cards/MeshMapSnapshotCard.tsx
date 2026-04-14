import React from 'react';
import { Button, Grid2, Typography } from '@sistent/sistent';
import {
  UnifiedCardContainer,
  UnifiedDescription,
  UnifiedButtonContainer,
  ResponsiveImage,
} from '../shared';
import { GRID_SIZE } from '../constants';

const MeshMapSnapshotCard = ({ githubActionEnabled = false }) => {
  const handleEnable = (e) => {
    window.open('https://cloud.layer5.io/connect/github/new/', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="kanvas-snapshot-heading" variant="h5" component="div">
          GitHub Action: Kanvas Snapshot
        </Typography>

        <UnifiedDescription data-testid="kanvas-snapshot-description" variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/extensions/meshmap-snapshot-logo.svg"
            alt="Kanvas Snapshot Logo"
            testId="kanvas-snapshot-image"
          />
          <div>
            Connect Kanvas to your GitHub repo and see changes pull request-to-pull request. Get
            snapshots of your infrastructure directly in your PRs.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="kanvas-snapshot-enable-btn"
            disabled={githubActionEnabled}
            onClick={(e) => handleEnable(e)}
          >
            {githubActionEnabled ? 'Remove' : 'Enable'}
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export default MeshMapSnapshotCard;
