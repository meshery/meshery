import React, { useState, useRef } from 'react';
import { Typography, Grid2, Button, Paper } from '@sistent/sistent';
import RJSFWrapper from '../MesheryMeshInterface/PatternService/RJSF_wrapper';
import { RelationshipDefinitionV1Alpha3OpenApiSchema } from '@meshery/schemas';

/**
 *  Relationship Builder component using React JSON Schema Form
 * @param {Object} props - The component props
 * @param {Function} props.onSuccess - Callback function when relationship is successfully created
 * @param {Function} props.onError - Callback function when there's an error creating the relationship
 */
const RelationshipBuilder = ({ onSuccess, onError }) => {
  const formRef = useRef();
  const [formData, setFormData] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  // -> Import the Relationship definition schema from Meshery schemas
  const RelationshipDefinitionV1Alpha3Schema =
    RelationshipDefinitionV1Alpha3OpenApiSchema.components.schemas.RelationshipDefinition;

  const handleChange = (data) => {
    setFormData(data.formData);
  };

  const handleSubmit = async () => {
    if (formRef.current) {
      // Trigger validation of the form
      formRef.current.validateForm();
      // Check if there are any errors after validation
      const formState = formRef.current?.state;
      if (formState) {
        setIsSubmitting(true);

        try {
          console.log('Valid Relationship Definition:', formData);
          onSuccess && onSuccess(formData);
          setFormData({});
        } catch (error) {
          console.error('Error creating relationship:', error);
          onError && onError(error.message || 'Failed to create relationship');
        } finally {
          setIsSubmitting(false);
        }
      } else {
        onError && onError('Please fix validation errors before submitting');
      }
    }
  };
  return (
    <Paper elevation={1} style={{ padding: '20px', marginTop: '20px' }}>
      <Grid2 container spacing={3}>
        <Grid2 item xs={12}>
          <Typography variant="h5" gutterBottom>
            Relationship Builder
          </Typography>
          <Typography variant="body2" gutterBottom>
            Create a new relationship definition to define how components interact with each other.
          </Typography>
        </Grid2>
        <Grid2 item xs={12}>
          <RJSFWrapper
            jsonSchema={RelationshipDefinitionV1Alpha3Schema}
            formData={formData}
            onChange={handleChange}
            formRef={formRef}
            liveValidate={false}
          />
        </Grid2>
        <Grid2 item xs={12}>
          <Button
            variant="contained"
            color="primary"
            onClick={handleSubmit}
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Creating...' : 'Create Relationship'}
          </Button>
        </Grid2>
      </Grid2>
    </Paper>
  );
};

export default RelationshipBuilder;
