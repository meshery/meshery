import React from 'react';
import { CatalogIcon, Grid2, Switch, Typography, Box, useTheme } from '@sistent/sistent';
import { UnifiedCardContainer, UnifiedDescription } from '../shared';
import { GRID_SIZE } from '../constants';

const CatalogCard = ({ catalogContent, handleToggle }) => {
  const theme = useTheme();

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="catalog-section-heading" variant="h5" component="div">
          {'Meshery Catalog'}
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <CatalogIcon
            data-testid="catalog-toggle-switch"
            style={{
              height: '85px',
              width: '140px',
              flexShrink: 0,
            }}
          />

          <div
            style={{
              display: 'inline',
              position: 'relative',
            }}
          >
            Enable access to the cloud native catalog, supporting design patterns, WebAssembly
            filters (<span style={{ fontStyle: 'italic' }}>soon</span>), and OPA policies (
            <span style={{ fontStyle: 'italic' }}>soon</span>). Import any catalog item and
            customize.
          </div>
        </UnifiedDescription>

        <Box
          sx={{
            position: 'absolute',
            bottom: 12,
            left: 12,
            right: 12,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <Typography variant="subtitle2" style={{ fontStyle: 'italic' }}>
            Explore the{' '}
            <a
              href="https://meshery.io/catalog"
              target="_blank"
              data-testid="catalog-link"
              rel="noreferrer"
              style={{
                textDecoration: 'none',
                color: theme.palette.text.brand,
              }}
            >
              Meshery Catalog
            </a>
          </Typography>

          <Switch
            checked={catalogContent}
            onChange={handleToggle}
            name="OperatorSwitch"
            color="primary"
          />
        </Box>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export default CatalogCard;
