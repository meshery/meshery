import React from 'react';
import { Button, Grid2, Typography } from '@sistent/sistent';
import {
  UnifiedCardContainer,
  UnifiedDescription,
  UnifiedButtonContainer,
  ResponsiveImage,
} from '../shared';
import { GRID_SIZE } from '../constants';

const CloudExtensionCard = () => {
  const handleLearnMore = (e) => {
    window.open('https://meshery.io/extensions/layer5-cloud', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="layer5-cloud-heading" variant="h5" component="div">
          Layer5 Cloud
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/extensions/layer5.svg"
            alt="Layer5 Cloud Logo"
            testId="layer5-cloud-image"
          />
          <div>
            A comprehensive platform offering identity and collaboration services, private catalogs,
            GitOps, and multi-Meshery management. Leverage its extensible authorization framework
            and organizational hierarchy for streamlined cloud infrastructure management.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="layer5-cloud-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export default CloudExtensionCard;
