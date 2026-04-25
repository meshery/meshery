import React from 'react';
import { Button, Grid2, Typography } from '@sistent/sistent';
import {
  UnifiedCardContainer,
  UnifiedDescription,
  UnifiedButtonContainer,
  ResponsiveImage,
} from '../shared';
import { GRID_SIZE } from '../constants';

const HelmKanvasCard = () => {
  const handleLearnMore = (e) => {
    window.open('https://docs.meshery.io/extensions/helm-kanvas-snapshot', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography variant="h5" component="div">
          Kanvas Snapshot Helm Plugin
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/extensions/helm_chart.svg"
            alt="Helm Chart Logo"
            testId="helm-kanvas-icon"
          />
          <div>
            The Kanvas Snapshot Helm Plugin allows you to generate a visual snapshot of your Helm
            charts directly from the command line. It simplifies the process of creating Meshery
            Snapshots, providing a visual representation of packaged Helm charts.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="helm-kanvas-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export default HelmKanvasCard;
