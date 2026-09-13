import React from 'react';
import {
  MeshModelToolbar as StyledMeshModelToolbar,
  CardStyle,
} from '@/assets/styles/general/tool.styles';
import { MODELS, RELATIONSHIPS } from '../../constants/navigator';
import {
  Button,
  AddCircleIcon as AddIcon,
  ExternalLinkIcon as LinkIcon,
  FileUploadIcon as UploadIcon,
  styled,
} from '@sistent/sistent';
import { iconSmall } from 'css/icons.styles';

const ToolbarActionContainer = styled('div')({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  gap: '0.75rem',
  flexWrap: 'wrap',
});

const TabCountText = styled('span')({
  fontSize: '1rem',
  marginLeft: '4px',
});

export const TabBar = ({ openImportModal, openCreateModal, view, openRelationshipModal }: any) => (
  <StyledMeshModelToolbar>
    <ToolbarActionContainer>
      {view === MODELS && (
        <>
          <Button
            aria-label="Create Model"
            variant="contained"
            color="primary"
            onClick={openCreateModal}
            style={{ display: 'flex' }}
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
            style={{ display: 'flex' }}
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
          style={{ display: 'flex' }}
          startIcon={<LinkIcon />}
          data-testid="TabBar-Button-CreateRelationship"
        >
          Create Relationship
        </Button>
      )}
    </ToolbarActionContainer>
  </StyledMeshModelToolbar>
);

export const TabCard = ({ label, count, active, onClick }: any) => (
  <CardStyle isSelected={active} elevation={3} onClick={onClick}>
    <TabCountText>{`(${count?.toLocaleString() || 0})`}</TabCountText>
    {label}
  </CardStyle>
);
