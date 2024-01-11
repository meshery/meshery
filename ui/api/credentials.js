import { promisifiedDataFetch } from 'lib/data-fetch';

export const getCredentialByID = async (credentialID) => {
  return await promisifiedDataFetch(`/api/integrations/credentials/${credentialID}`);
};
