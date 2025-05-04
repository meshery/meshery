import React from 'react';
import { FormatDryRunResponse } from '../../DesignLifeCycle/DryRun';
import { formatDryRunResponse } from 'machines/validator/designValidator';

export const DryRunResponse = ({ response }) => {
  return <FormatDryRunResponse dryRunErrors={formatDryRunResponse(response)} />;
};
