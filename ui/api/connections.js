import { promisifiedDataFetch } from '../lib/data-fetch';

export const getConnectionStatusSummary = async () => {
  return await promisifiedDataFetch('/api/integrations/connections/status');
};
