import { useState, useEffect } from 'react';
import { ability } from '../utils/can';
import { useLazyGetUserKeysQuery } from './userKeys';

/**
 * Custom hook to fetch and retrieve user abilities based on the provided organization.
 * Uses the useGetUserKeysQuery hook internally to fetch user keys.
 *
 * @param {Object} org - The organization object.
 * @param {string} org.id - The unique identifier of the organization.
 * @param {boolean} skip - Flag to determine whether to skip fetching.
 *
 * @returns {Object} An object containing the response from useGetUserKeysQuery and additional `data` property with mapped abilities.
 */
export const useGetCurrentAbilities = (org, skip) => {
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
          action: key.function,
          subject: key.id,
        }));

        if (res.keys.length > 0) {
          ability.update(res.keys);
          sessionStorage.setItem('keys', JSON.stringify(res.keys));
        }
        setData({
          ...res,
          abilities: abilities,
        });
      }).catch((error)=>{
        console.error("Error when fetching keys in useGetUserAbilities custom hook", error)
      });
  }, [org.id, getUserQuery, skip]);

  return data;
};
