import React from "react";
import { styled, Typography, Box, Link, Grid, Divider } from "@sistent/sistent";
import { FavoriteIcon } from "@sistent/sistent";


const FooterContainer = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' 
    ? theme.palette.background.default 
    : theme.palette.background.paper,
  padding: theme.spacing(6, 2),
  borderTop: `1px solid ${theme.palette.divider}`,
  marginTop: 'auto',
}));

const FooterSection = styled(Box)(({ theme }) => ({
  marginBottom: theme.spacing(3),
}));

const SectionTitle = styled(Typography)(({ theme }) => ({
  fontWeight: 600,
  marginBottom: theme.spacing(2),
  color: theme.palette.text.primary,
  fontSize: '0.875rem',
  textTransform: 'uppercase',
  letterSpacing: '0.5px',
}));

const FooterLink = styled(Link)(({ theme }) => ({
  display: 'block',
  color: theme.palette.text.secondary,
  textDecoration: 'none',
  marginBottom: theme.spacing(1),
  fontSize: '0.875rem',
  transition: 'color 0.2s ease',
  '&:hover': {
    color: theme.palette.primary.main,
    textDecoration: 'none',
  },
}));

const CopyrightText = styled(Typography)(({ theme }) => ({
  color: theme.palette.text.disabled,
  fontSize: '0.75rem',
  textAlign: 'center',
}));

const StyledFavoriteIcon = styled(FavoriteIcon)(({ theme }) => ({
  display: "inline",
  verticalAlign: "middle",
  fontSize: "1rem",
  fill: theme.palette.error.main,
  margin: "0 4px",
}));


const footerSections = {
  project: {
    title: "Project",
    links: [
      { label: "About", href: "https://meshery.io" },
      { label: "Documentation", href: "https://docs.meshery.io" },
      { label: "GitHub", href: "https://github.com/meshery/meshery" },
      { label: "Roadmap", href: "https://github.com/meshery/meshery/blob/master/ROADMAP.md" },
      { label: "Contributing", href: "https://docs.meshery.io/project/contributing" },
    ],
  },
  gettingStarted: {
    title: "Getting Started",
    links: [
      { label: "Quick Start", href: "https://docs.meshery.io/installation/quick-start" },
      { label: "Tutorials", href: "https://docs.meshery.io/guides" },
      { label: "Playground", href: "https://play.meshery.io" },
      { label: "Catalog", href: "https://meshery.io/catalog" },
      { label: "Extensions", href: "https://meshery.io/extensions" },
    ],
  },
  community: {
    title: "Community",
    links: [
      { label: "Slack", href: "https://slack.meshery.io" },
      { label: "Discussion Forum", href: "https://discuss.layer5.io" },
      { label: "Community Calendar", href: "https://meshery.io/calendar" },
      { label: "Events", href: "https://layer5.io/community/events" },
      { label: "Newcomers", href: "https://layer5.io/community/newcomers" },
    ],
  },
  social: {
    title: "Social",
    links: [
      { label: "Twitter", href: "https://twitter.com/mesheryio" },
      { label: "LinkedIn", href: "https://www.linkedin.com/company/meshery" },
      { label: "YouTube", href: "https://www.youtube.com/channel/UCgXlqWDCg-9RP1eckf0s6KA" },
      { label: "Blog", href: "https://layer5.io/blog" },
      { label: "Newsletter", href: "https://layer5.io/subscribe" },
    ],
  },
};

export default function UnifiedFooter() {
  const currentYear = new Date().getFullYear();

  return (
    <FooterContainer component="footer">
      <Grid container spacing={4}>
        
        <Grid item xs={12} sm={6} md={3}>
          <FooterSection>
            <SectionTitle variant="h6">{footerSections.project.title}</SectionTitle>
            {footerSections.project.links.map((link, index) => (
              <FooterLink
                key={index}
                href={link.href}
                target="_blank"
                rel="noopener noreferrer"
              >
                {link.label}
              </FooterLink>
            ))}
          </FooterSection>
        </Grid>

        
        <Grid item xs={12} sm={6} md={3}>
          <FooterSection>
            <SectionTitle variant="h6">{footerSections.gettingStarted.title}</SectionTitle>
            {footerSections.gettingStarted.links.map((link, index) => (
              <FooterLink
                key={index}
                href={link.href}
                target="_blank"
                rel="noopener noreferrer"
              >
                {link.label}
              </FooterLink>
            ))}
          </FooterSection>
        </Grid>

        
        <Grid item xs={12} sm={6} md={3}>
          <FooterSection>
            <SectionTitle variant="h6">{footerSections.community.title}</SectionTitle>
            {footerSections.community.links.map((link, index) => (
              <FooterLink
                key={index}
                href={link.href}
                target="_blank"
                rel="noopener noreferrer"
              >
                {link.label}
              </FooterLink>
            ))}
          </FooterSection>
        </Grid>

        
        <Grid item xs={12} sm={6} md={3}>
          <FooterSection>
            <SectionTitle variant="h6">{footerSections.social.title}</SectionTitle>
            {footerSections.social.links.map((link, index) => (
              <FooterLink
                key={index}
                href={link.href}
                target="_blank"
                rel="noopener noreferrer"
              >
                {link.label}
              </FooterLink>
            ))}
          </FooterSection>
        </Grid>
      </Grid>

      
      <Divider sx={{ my: 3 }} />

      
      <Box sx={{ textAlign: 'center' }}>
        <Typography variant="body2" sx={{ mb: 1 }}>
          Built with <StyledFavoriteIcon /> by the Meshery Community
        </Typography>
        <CopyrightText>
          © {currentYear} Meshery. All rights reserved.
        </CopyrightText>
      </Box>
    </FooterContainer>
  );
}