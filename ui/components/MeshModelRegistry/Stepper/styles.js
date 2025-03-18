import { Box, styled } from '@layer5/sistent';
import Sketch from '@uiw/react-color-sketch';

export const StyledSketchContainer = styled(Box)(({ theme, primaryColor }) => ({
  width: '100%',
  height: '25px',
  borderRadius: '4px',
  border: `5px solid ${theme.palette.border.normal}`,
  backgroundColor: primaryColor,
  cursor: 'pointer',
  marginTop: '0.5rem',
}));
export const StyledSketchPicker = styled(Sketch)(({ theme }) => ({
  backgroundColor: `${theme.palette.background.surfaces} !important`,
}));

export const StyledSketchWrapper = styled(Box)(() => ({
  position: 'absolute',
  zIndex: 2,
  bottom: '25px',
}));
