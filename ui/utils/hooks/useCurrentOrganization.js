import { useLegacySelector } from 'lib/store';

export const useCurrentOrganization = () => {
  const org = useLegacySelector((state) => {
    try {
      if (!state) return undefined;
      if (typeof state.get === 'function') {
        return state.get('organization');
      }
      return state.organization;
    } catch (err) {
      return undefined;
    }
  });

  return org;
};
