import React from 'react';
import { Button, Grid2, Typography } from '@sistent/sistent';
import {
  UnifiedCardContainer,
  UnifiedDescription,
  UnifiedButtonContainer,
  ResponsiveImage,
} from '../shared';
import { GRID_SIZE } from '../constants';

const DesignEmbedCard = () => {
  const handleLearnMore = (e) => {
    window.open('https://docs.layer5.io/kanvas/designer/embedding-designs/', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography variant="h5" component="div">
          Meshery Design Embed
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/extensions/meshmap.svg"
            alt="Meshery Design Logo"
            testId="design-embed-icon"
          />
          <div>
            Meshery Design Embedding lets you export designs in an interactive format that
            seamlessly integrates with websites, blogs, and platforms using HTML, CSS, and
            JavaScript, making it easy to share with stakeholders.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="design-embed-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export default DesignEmbedCard;
