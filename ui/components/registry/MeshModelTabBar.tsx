import React from 'react';
import {
  Button,
  AddCircleIcon as AddIcon,
  ExternalLinkIcon as LinkIcon,
  FileUploadIcon as UploadIcon,
} from '@sistent/sistent';
import { MeshModelToolbar } from '@/assets/styles/general/tool.styles';
import { styled } from '@/theme';
import { MODELS, RELATIONSHIPS } from '../../constants/navigator';
import { iconSmall } from 'css/icons.styles';

const ToolbarActions = styled('div')({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  gap: '0.75rem',
  flexWrap: 'wrap',
});

type MeshModelTabBarProps = {
  openImportModal: () => void;
  openCreateModal: () => void;
  openRelationshipModal: () => void;
  view: string;
};

/**
 * Toolbar actions rendered above the registry resource tree for creating and
 * importing models or relationships.
 */
export const MeshModelTabBar = ({
  openImportModal,
  openCreateModal,
  view,
  openRelationshipModal,
}: MeshModelTabBarProps) => {
  return (
    <MeshModelToolbar>
      <ToolbarActions>
        {view === MODELS && (
          <>
            <Button
              aria-label="Create Model"
              variant="contained"
              color="primary"
              onClick={openCreateModal}
              disabled={false} //TODO: Need to make key for this component
              startIcon={<AddIcon style={iconSmall} />}
              data-testid="TabBar-Button-CreateModel"
            >
              Create Model
            </Button>
            <Button
              aria-label="Import Model"
              variant="contained"
              color="primary"
              onClick={openImportModal}
              disabled={false} //TODO: Need to make key for this component
              startIcon={<UploadIcon />}
              data-testid="TabBar-Button-ImportModel"
            >
              Import Model
            </Button>
          </>
        )}

        {view === RELATIONSHIPS && (
          <Button
            aria-label="Create Relationship"
            variant="contained"
            color="primary"
            onClick={openRelationshipModal}
            disabled={false}
            startIcon={<LinkIcon />}
            data-testid="TabBar-Button-CreateRelationship"
          >
            Create Relationship
          </Button>
        )}
      </ToolbarActions>
    </MeshModelToolbar>
  );
};
