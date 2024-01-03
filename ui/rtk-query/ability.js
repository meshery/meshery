import { useState, useEffect } from 'react';
import { ability } from '../utils/can';
import { useLazyGetUserKeysQuery } from './userKeys';

export const useGetUserAbilities = (org, skip) => {
  const [data, setData] = useState(null);

  /**
   * RTk Lazy Query
   */
  const [getUserQuery] = useLazyGetUserKeysQuery();

  useEffect(() => {
    getUserQuery({ orgId: org.id }, { skip })
      .unwrap()
      .then((res) => {
        const abilities = res.keys?.map((key) => ({
          action: key.id,
          subject: key.function,
        }));

        setData({
          ...res,
          abilities: abilities,
        });
      })
      .catch((error) => {
        console.error('Error when fetching keys in useGetUserAbilities custom hook', error);
      });
  }, [org.id, getUserQuery, skip]);

  return data;
};

export const useGetCurrentAbilities = (org, skip) => {
  const res = useGetUserAbilities(org, skip)
  if(res?.abilities){
    ability.update(res.abilities);

    // set keys to session storage
    sessionStorage.setItem('keys', JSON.stringify(res.keys));
  }
  return res
}