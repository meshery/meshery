import React from 'react';
import { Box } from '@sistent/sistent';
import { CardContainer, FrontSideDescription } from '../../css/icons.styles';

export const UnifiedCardContainer = ({ children, sx = {} }) => {
  return (
    <CardContainer
      sx={{
        minHeight: {
          xs: '280px',
          sm: '260px',
          lg: '280px',
        },
        display: 'flex',
        flexDirection: 'column',
        position: 'relative',
        ...sx,
      }}
    >
      {children}
    </CardContainer>
  );
};

export const UnifiedDescription = ({ children, hasIcon = false, ...props }) => {
  return (
    <FrontSideDescription
      {...props}
      sx={{
        flex: 1,
        display: 'flex',
        flexDirection: {
          xs: hasIcon ? 'column' : 'row',
          sm: 'row',
        },
        alignItems: {
          xs: hasIcon ? 'center' : 'flex-start',
          sm: 'flex-start',
        },
        textAlign: {
          xs: hasIcon ? 'center' : 'left',
          sm: 'left',
        },
        gap: {
          xs: hasIcon ? '12px' : '8px',
          sm: '8px',
        },
        marginBottom: {
          xs: '50px',
          sm: '45px',
          lg: '40px',
        },
        overflow: 'hidden',
        wordWrap: 'break-word',
        hyphens: 'auto',
        '& > *:not(img)': {
          maxWidth: '100%',
          overflow: 'hidden',
          wordWrap: 'break-word',
          hyphens: 'auto',
        },
      }}
    >
      {children}
    </FrontSideDescription>
  );
};

export const UnifiedButtonContainer = ({ children }) => {
  return (
    <Box
      sx={{
        position: 'absolute',
        bottom: 12,
        left: 12,
        right: 12,
        textAlign: 'right',
      }}
    >
      {children}
    </Box>
  );
};

export const ResponsiveImage = ({ src, alt, testId }) => {
  return (
    <img
      style={{
        height: 'auto',
        width: 'auto',
        maxWidth: '140px',
        maxHeight: '85px',
        minWidth: '100px',
        minHeight: '60px',
        flexShrink: 0,
      }}
      data-testid={testId}
      src={src}
      alt={alt || ''}
    />
  );
};
