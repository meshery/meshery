import React from 'react';
import { MeshModelToolbar, CardStyle } from '@/assets/styles/general/tool.styles';
import {
  Button,
  AddCircleIcon as AddIcon,
  ExternalLinkIcon as LinkIcon,
  FileUploadIcon as UploadIcon,
  styled,
} from '@sistent/sistent';
import { iconSmall } from 'css/icons.styles';
import { MODELS, RELATIONSHIPS } from '../../constants/navigator';
const TabBarContainer = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  gap: '0.75rem',
  flexWrap: 'wrap',
}));

const TabCardCount = styled('span')(() => ({
  fontSize: '1rem',
  marginLeft: '4px',
}));

export const TabBar = ({
  openImportModal,
  openCreateModal,
  view,
  openRelationshipModal,
}: {
  openImportModal: () => void;
  openCreateModal: () => void;
  view: string;
  openRelationshipModal: () => void;
}) => {
  return (
    <MeshModelToolbar>
      <TabBarContainer>
        {view === MODELS && (
          <>
            <Button
              aria-label="Create Model"
              variant="contained"
              color="primary"
              onClick={openCreateModal}
              style={{ display: 'flex' }}
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
              style={{ display: 'flex' }}
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
            style={{ display: 'flex' }}
            disabled={false}
            startIcon={<LinkIcon />}
            data-testid="TabBar-Button-CreateRelationship"
          >
            Create Relationship
          </Button>
        )}
      </TabBarContainer>
    </MeshModelToolbar>
  );
};

export const TabCard = ({
  label,
  count,
  active,
  onClick,
}: {
  label: string;
  count: number;
  active: boolean;
  onClick: () => void;
}) => {
  return (
    <CardStyle isSelected={active} onClick={onClick}>
      <TabCardCount>{`(${count?.toLocaleString() || 0})`}</TabCardCount>
      {label}
    </CardStyle>
  );
};
