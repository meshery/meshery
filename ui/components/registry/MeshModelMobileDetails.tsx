import React, { useMemo } from 'react';
import { BottomSheet, ThemeProvider, createTheme } from '@sistent/sistent';
import { useTheme } from '@/theme';
import MeshModelDetails from './MeshModelDetails';

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
  const theme = useTheme();
  const bottomSheetTheme = useMemo(
    () =>
      createTheme(theme, {
        zIndex: {
          modal: theme.zIndex.tooltip + 1,
        },
      }),
    [theme],
  );

  const open = Boolean(
    showDetailsData?.data &&
    Object.keys(showDetailsData.data).length > 0 &&
    showDetailsData.type !== 'none',
  );

  return (
    <ThemeProvider theme={bottomSheetTheme}>
      <BottomSheet
        open={open}
        onClose={() => setShowDetailsData({ type: '', data: {} })}
        title={showDetailsData.data?.displayName || showDetailsData.type}
        closeButtonAriaLabel="Close details"
      >
        <MeshModelDetails view={view} showDetailsData={showDetailsData} />
      </BottomSheet>
    </ThemeProvider>
  );
};

export default MeshModelMobileDetails;
