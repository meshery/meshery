import { Typography, useTheme, Box, Card, styled } from '@sistent/sistent';
import React from 'react';
import FlipCard from './FlipCard';

interface ResourceCardProps {
  title: string;
  subtitle?: React.ReactNode;
  headerRight?: React.ReactNode;
  footer?: React.ReactNode;
  children?: React.ReactNode;
  // Back side props
  backTitle?: string;
  backSubtitle?: React.ReactNode;
  backHeaderRight?: React.ReactNode;
  backFooter?: React.ReactNode;
  backContent?: React.ReactNode;
  // Control props
  onClick?: (ev: React.MouseEvent) => void;
  onFlip?: () => void;
  isFlippable?: boolean;
}

const StyledCard = styled(Card)(({ theme }) => ({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'space-between',
  borderRadius: '8px',
  '&:hover': {
    boxShadow: theme.shadows[4],
  },
  padding: '16px',
  backgroundColor: theme.palette.background.paper,
}));

/**
 * Base ResourceCard primitive for Meshery UI.
 * Standardizes the layout for resources like Patterns, Filters, etc.
 * Supports flippable behavior for back-side content.
 */
const ResourceCard: React.FC<ResourceCardProps> = ({
  title,
  subtitle,
  headerRight,
  footer,
  children,
  backTitle,
  backSubtitle,
  backHeaderRight,
  backFooter,
  backContent,
  onClick,
  onFlip,
  isFlippable = false,
}) => {
  const theme = useTheme();

  const renderSide = (
    t: string,
    s: React.ReactNode,
    hr: React.ReactNode,
    f: React.ReactNode,
    c: React.ReactNode,
    isBack = false,
  ) => (
    <StyledCard onClick={!isBack ? onClick : undefined}>
      <Box sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
        <Box display="flex" justifyContent="space-between" alignItems="flex-start">
          <Typography
            variant="h6"
            sx={{
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              width: '100%',
            }}
          >
            {t}
          </Typography>
          {hr}
        </Box>
        {s && (
          <Box mb={1}>
            <Typography
              variant="caption"
              sx={{
                fontStyle: 'italic',
                color: theme.palette.text.secondary,
              }}
            >
              {s}
            </Typography>
          </Box>
        )}
        <Box mt={1} sx={{ flexGrow: 1 }}>
          {c}
        </Box>
      </Box>
      {f && (
        <Box
          mt={2}
          display="flex"
          justifyContent="flex-start"
          alignItems="center"
          gap="8px"
          flexWrap="wrap"
        >
          {f}
        </Box>
      )}
    </StyledCard>
  );

  const Front = renderSide(title, subtitle, headerRight, footer, children);
  const Back =
    isFlippable && (backContent || backTitle)
      ? renderSide(
          backTitle || title,
          backSubtitle || subtitle,
          backHeaderRight,
          backFooter,
          backContent,
          true,
        )
      : null;

  if (isFlippable && Back) {
    return <FlipCard frontComponents={Front} backComponents={Back} onFlip={onFlip} />;
  }

  return Front;
};

export default ResourceCard;
