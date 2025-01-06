import React from 'react';
import { FormControl, InputLabel, MenuItem, Select, styled } from '@material-ui/core';

const StyledFormControl = styled(FormControl)({
  minWidth: 200,
});

export const ResourceSelector = ({ kinds, selectedkind, onKindChange }) => {
  const handleChange = (event) => {
    onKindChange(event.target.value);
  };

  return (
    <StyledFormControl size="small">
      <InputLabel>Select Resource</InputLabel>
      <Select value={selectedkind} label="Select Cluster" onChange={handleChange}>
        <MenuItem value="all">All Resource</MenuItem>
        {kinds.map((kind) => (
          <MenuItem key={kind} value={kind}>
            {kind}
          </MenuItem>
        ))}
      </Select>
    </StyledFormControl>
  );
};
