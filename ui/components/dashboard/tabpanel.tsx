import React from 'react';
import { Typography } from '@sistent/sistent';

type TabPanelProps = {
  children?: React.ReactNode;
  value: number | string;
  index: number | string;
};

export function TabPanel({ children, value, index, ...other }: TabPanelProps) {
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Typography component="div" style={{ paddingTop: 2 }}>
          {children}
        </Typography>
      )}
    </div>
  );
}
