import { Chip } from '@material-ui/core';

export const ConnectionChip = ({ handlePing, title, icon }) => (
  <Chip label={title} onClick={() => handlePing()} icon={icon} variant="outlined" />
);
