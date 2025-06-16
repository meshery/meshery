import React from 'react';
import { Typography, Container, Grid2, Box } from '@sistent/sistent';
import RelationshipBuilder from '../../components/RelationshipBuilder';
import Head from 'next/head';
import { useNotification } from '../../utils/hooks/useNotification';

function RelationshipBuilderPage() {
  const { notify } = useNotification();

  return (
    <div>
      <Head>
        <title>Relationship Builder - Meshery</title>
      </Head>
      <Container maxWidth="lg">
        <Box mb={4}>
          <Typography variant="h5">Create Relationship Definition</Typography>
          <Typography variant="body1" color="textSecondary">
            Define the interaction pattern between components in your system
          </Typography>
        </Box>

        <Grid2 container spacing={4}>
          <Grid2 item xs={12} md={12}>
            <RelationshipBuilder
              onSuccess={() => {
                notify({
                  message: 'Relationship successfully created',
                  variant: 'success',
                });
              }}
              onError={(error) => {
                notify({
                  message: `Error creating relationship: ${error}`,
                  variant: 'error',
                });
              }}
            />
          </Grid2>
        </Grid2>
      </Container>
    </div>
  );
}

export default RelationshipBuilderPage;
