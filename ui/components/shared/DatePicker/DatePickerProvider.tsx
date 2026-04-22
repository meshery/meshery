import React from 'react';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterMoment } from '@mui/x-date-pickers/AdapterMoment';

interface MesheryDatePickerProviderProps {
  children: React.ReactNode;
}

const MesheryDatePickerProvider = ({ children }: MesheryDatePickerProviderProps) => {
  return <LocalizationProvider dateAdapter={AdapterMoment}>{children}</LocalizationProvider>;
};

export default MesheryDatePickerProvider;