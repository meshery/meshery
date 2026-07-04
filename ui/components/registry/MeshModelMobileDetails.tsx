import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  IconButton,
  CloseIcon,
  alpha,
  styled,
} from '@sistent/sistent';
import MeshModelDetails from './MeshModelDetails';

const TitleSpan = styled('span')(() => ({
  fontWeight: 600,
  fontSize: '1rem',
  flex: 1,
  minWidth: 0,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
}));

const ContentDiv = styled('div')(() => ({
  width: '100%',
}));

const MeshModelMobileDetails = ({
  showDetailsData,
  setShowDetailsData,
  view,
}: {
  showDetailsData: {
    type: string;
    data: any;
  };
  setShowDetailsData: React.Dispatch<
    React.SetStateAction<{
      type: string;
      data: any;
    }>
  >;
  view: string;
}) => {
  return (
    <Dialog
      open={Boolean(
        showDetailsData?.data &&
        Object.keys(showDetailsData.data).length > 0 &&
        showDetailsData.type !== 'none',
      )}
      onClose={() => setShowDetailsData({ type: '', data: {} })}
      fullWidth
      maxWidth="sm"
      sx={{
        zIndex: 1600,
        '& .MuiDialog-container': {
          alignItems: 'flex-end',
        },
        '& .MuiDialog-paper': {
          margin: 0,
          width: '100%',
          maxWidth: '100%',
          borderRadius: '12px 12px 0 0',
          maxHeight: '80vh',
        },
      }}
    >
      <DialogTitle
        sx={(theme) => ({
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          padding: '0.75rem 1rem',
          borderBottom: `1px solid ${alpha(theme.palette.divider, 1)}`,
        })}
      >
        <TitleSpan>{showDetailsData.data?.displayName || showDetailsData.type}</TitleSpan>
        <IconButton
          aria-label="Close details"
          onClick={() => setShowDetailsData({ type: '', data: {} })}
          size="small"
          sx={{ flexShrink: 0 }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <DialogContent sx={{ padding: '1rem', overflowY: 'auto' }}>
        <ContentDiv>
          <MeshModelDetails view={view} showDetailsData={showDetailsData} />
        </ContentDiv>
      </DialogContent>
    </Dialog>
  );
};

export default MeshModelMobileDetails;
