import React from 'react';
import { Box } from '@sistent/sistent';
import { CardContainer, FrontSideDescription } from 'css/icons.styles';

type ChildrenProps = {
  children: React.ReactNode;
};

type UnifiedDescriptionProps = Omit<
  React.ComponentProps<typeof FrontSideDescription>,
  'children' | 'sx'
> & {
  children: React.ReactNode;
  hasIcon?: boolean;
};

export const UnifiedCardContainer = ({ children, sx = {} }: ChildrenProps & { sx?: object }) => (
  <CardContainer
    sx={{
      height: '100%',
      minHeight: { xs: '280px', sm: '260px', lg: '280px' },
      display: 'flex',
      flexDirection: 'column',
      position: 'relative',
      ...sx,
    }}
  >
    {children}
  </CardContainer>
);

export const UnifiedDescription = ({
  children,
  hasIcon = false,
  ...props
}: UnifiedDescriptionProps) => (
  <FrontSideDescription
    {...props}
    sx={{
      flex: 1,
      display: 'flex',
      flexDirection: { xs: hasIcon ? 'column' : 'row', sm: 'row' },
      alignItems: { xs: hasIcon ? 'center' : 'flex-start', sm: 'flex-start' },
      textAlign: { xs: hasIcon ? 'center' : 'left', sm: 'left' },
      gap: { xs: hasIcon ? '12px' : '8px', sm: '8px' },
      marginBottom: { xs: '50px', sm: '45px', lg: '40px' },
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

export const UnifiedButtonContainer = ({ children }: ChildrenProps) => (
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
