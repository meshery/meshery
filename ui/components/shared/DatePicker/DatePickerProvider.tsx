import React from 'react';
import { AdapterMoment } from '@mui/x-date-pickers/AdapterMoment';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';

interface MesheryDatePickerProviderProps {
  children: React.ReactNode;
}

const MesheryDatePickerProvider = ({ children }: MesheryDatePickerProviderProps) => {
  return <LocalizationProvider dateAdapter={AdapterMoment}>{children}</LocalizationProvider>;
};

export default MesheryDatePickerProvider;
