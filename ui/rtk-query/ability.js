import { useLegacySelector } from 'lib/store';
import { ability } from '../utils/can';
import { useGetUserKeysQuery } from './userKeys';
import _ from 'lodash';
import CustomErrorMessage from '@/components/ErrorPage';
import LoadingScreen from '@/components/LoadingComponents/LoadingComponentServer';
import { randomLoadingMessage } from '@/components/LoadingComponents/loadingMessages';

export const useGetUserAbilities = (org, skip) => {
  const { data, ...res } = useGetUserKeysQuery(
    {
      orgId: org?.id,
    },
    {
      skip,
    },
  );

  const abilities =
    data?.keys?.map((key) => ({
      action: key.id,
      subject: _.lowerCase(key.function),
    })) || [];

  return {
    ...res,
    abilities,
  };
};

export const useGetCurrentAbilities = (org, setKeys) => {
  const shouldSkip = !org || !org.id;
  const res = useGetUserAbilities(org, shouldSkip);
  if (res?.abilities) {
    ability.update(res.abilities);
    setKeys({ keys: res.keys });
  }
  return res;
};

export const LoadSessionGuard = ({ children }) => {
  // this assumes that the organization is already loaded at the app mount time
  // otherwise, this will not work
  const org = useLegacySelector((state) => state.get('organization'));
  const { isLoading, error } = useGetCurrentAbilities(org, () => {});

  if (error) {
    return (
      <CustomErrorMessage
        message={error.message || 'An error occurred while fetching your organization permissions'}
      />
    );
  }

  return (
    <LoadingScreen isLoading={isLoading || !org?.id} message={randomLoadingMessage}>
      {children}
    </LoadingScreen>
  );
};
