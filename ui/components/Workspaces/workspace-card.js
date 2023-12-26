import React from 'react';
import { Button, Grid } from '@material-ui/core';
import { Delete, Edit } from '@material-ui/icons';

import OrgIcon from '../../assets/icons/OrgIcon';
import {
  BulkSelectCheckbox,
  CardTitle,
  CardWrapper,
  DateLabel,
  DescriptionLabel,
  EmptyDescription,
  OrganizationName,
  StyledIconButton,
} from './styles';
import FlipCard from '../Environments/flip-card';
import theme from '../../themes/app';

export const formattoLongDate = (date) => {
  return new Date(date).toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
};

/**
 * Renders a Workspace card component.
 *
 * @param {Object} props - The component props.
 * @param {Object} props.workspaceDetails - The details of the workspace.
 * @param {string} props.workspaceDetails.name - The name of the workspace.
 * @param {string} props.workspaceDetails.description - The description of the workspace.
 * @param {Function} props.onDelete - Function to delete the workspace.
 * @param {Function} props.onEdit - Function to edit the workspace.
 * @param {Function} props.onSelect - Function to select workspace for bulk actions.
 * @param {Array} props.selectedWorkspaces - Selected workspace list for delete.
 *
 */

const WorkspaceCard = ({ workspaceDetails, onDelete, onEdit, onSelect, selectedWorkspaces }) => {
  return (
    <>
      <FlipCard
        disableFlip={
          selectedWorkspaces?.filter((id) => id == workspaceDetails.id).length === 1 ? true : false
        }
        frontComponents={
          <CardWrapper
            elevation={2}
            style={{
              minHeight: '290px',
            }}
          >
            <Grid style={{ display: 'flex', flexDirection: 'row', paddingBottom: '5px' }}>
              <CardTitle variant="body2" onClick={(e) => e.stopPropagation()}>
                {workspaceDetails?.name}
              </CardTitle>
            </Grid>
            <Grid style={{ display: 'flex', alignItems: 'center', marginTop: '10px' }}>
              <StyledIconButton onClick={(e) => e.stopPropagation()}>
                <OrgIcon width="24" height="24" />
              </StyledIconButton>
              <OrganizationName variant="span" onClick={(e) => e.stopPropagation()}>
                {workspaceDetails?.owner}
              </OrganizationName>
            </Grid>
          </CardWrapper>
        }
        backComponents={
          <CardWrapper
            elevation={2}
            style={{
              minHeight: '290px',
              background: 'linear-gradient(180deg, #007366 0%, #000 100%)',
            }}
          >
            <Grid xs={12}>
              <Grid xs={12} style={{ display: 'flex', flexDirection: 'row' }}>
                <Grid xs={6} style={{ display: 'flex', alignItems: 'flex-start' }}>
                  <BulkSelectCheckbox onClick={(e) => e.stopPropagation()} onChange={onSelect} />
                  <CardTitle
                    style={{ color: theme.palette.secondary.white }}
                    variant="body2"
                    onClick={(e) => e.stopPropagation()}
                  >
                    {workspaceDetails?.name}
                  </CardTitle>
                </Grid>
                <Grid
                  xs={6}
                  style={{
                    display: 'flex',
                    alignItems: 'flex-start',
                    justifyContent: 'flex-end',
                  }}
                >
                  <Button
                    style={{
                      minWidth: 'fit-content',
                      '&.MuiButtonBase-root:hover': {
                        bgcolor: 'transparent',
                      },
                      padding: 0,
                    }}
                    onClick={onEdit}
                    disabled={
                      selectedWorkspaces?.filter((id) => id == workspaceDetails.id).length === 1
                        ? true
                        : false
                    }
                  >
                    <Edit style={{ color: 'white', margin: '0 2px' }} />
                  </Button>
                  <Button
                    style={{
                      minWidth: 'fit-content',
                      '&.MuiButtonBase-root:hover': {
                        bgcolor: 'transparent',
                      },
                      padding: 0,
                    }}
                    onClick={onDelete}
                    disabled={
                      selectedWorkspaces?.filter((id) => id == workspaceDetails.id).length === 1
                        ? true
                        : false
                    }
                  >
                    <Delete style={{ color: 'white', margin: '0 2px' }} />
                  </Button>
                </Grid>
              </Grid>
              <Grid>
                {workspaceDetails.description ? (
                  <DescriptionLabel
                    onClick={(e) => e.stopPropagation()}
                    style={{
                      color: theme.palette.secondary.white,
                      maxHeight: '105px',
                      padding: '10px',
                    }}
                  >
                    {workspaceDetails.description}
                  </DescriptionLabel>
                ) : (
                  <EmptyDescription
                    onClick={(e) => e.stopPropagation()}
                    style={{
                      color: 'rgba(255, 255, 255, .6)',
                      paddingTop: '20px',
                      paddingBottom: '20px',
                    }}
                  >
                    No description
                  </EmptyDescription>
                )}
              </Grid>
            </Grid>
            <Grid
              style={{
                display: 'flex',
                flexDirection: 'row',
                position: 'absolute',
                bottom: '20px',
                width: '100%',
                color: `${theme.palette.secondary.white}99`,
              }}
            >
              <Grid xs={6} style={{ textAlign: 'left' }}>
                <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
                  Updated At: {formattoLongDate(workspaceDetails?.updated_at)}
                </DateLabel>
              </Grid>
              <Grid xs={6} style={{ textAlign: 'left' }}>
                <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
                  Created At: {formattoLongDate(workspaceDetails?.created_at)}
                </DateLabel>
              </Grid>
            </Grid>
          </CardWrapper>
        }
      />
    </>
  );
};

export default WorkspaceCard;
