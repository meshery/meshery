import { useState, useEffect } from 'react';
import _ from 'lodash';
import ExtensionPointSchemaValidator from '../../../utils/ExtensionPointSchemaValidator';
import { EVENT_TYPES } from '../../../lib/event-types';
import { useNotification } from '../../../utils/hooks/useNotification';
import {
  useGetProviderCapabilitiesQuery,
  useGetUserPrefQuery,
  useUpdateUserPrefMutation,
  useUpdateUserPrefWithContextMutation,
} from '@/rtk-query/user';
import { useDispatch, useSelector } from 'react-redux';
import { toggleCatalogContent, updateProgress } from '@/store/slices/mesheryUi';

export const useUserPreferences = (
  initialAnonymousStats?: boolean,
  initialPerfResultStats?: boolean,
) => {
  const [anonymousStats, setAnonymousStats] = useState(initialAnonymousStats);
  const [perfResultStats, setPerfResultStats] = useState(initialPerfResultStats);
  const [userPrefs, setUserPrefs] = useState(ExtensionPointSchemaValidator('user_prefs')());
  const [providerType, setProviderType] = useState('');
  const [catalogContent, setCatalogContent] = useState(true);
  const [extensionPreferences, setExtensionPreferences] = useState({});
  const [capabilitiesLoaded, setCapabilitiesLoaded] = useState(false);
  const [providerInfo, setProviderInfo] = useState({});

  const dispatch = useDispatch();
  const { capabilitiesRegistry } = useSelector((state: any) => state.ui);

  const {
    data: userData,
    isSuccess: isUserDataFetched,
    isError: isUserDataError,
    error: userDataError,
  } = useGetUserPrefQuery();

  const { data: capabilitiesData, isSuccess: isCapabilitiesDataFetched } =
    useGetProviderCapabilitiesQuery();

  const [updateUserPref] = useUpdateUserPrefMutation();
  const [updateUserPrefWithContext] = useUpdateUserPrefWithContextMutation();

  const { notify } = useNotification();

  const handleCatalogContentToggle = () => {
    dispatch(toggleCatalogContent({ catalogVisibility: !catalogContent }));
    setCatalogContent(!catalogContent);
    handleCatalogPreference(!catalogContent);
  };

  const handleCatalogPreference = (catalogContent: boolean) => {
    let body = Object.assign({}, extensionPreferences);
    body['catalogContent'] = catalogContent;
    updateUserPref({ usersExtensionPreferences: body })
      .unwrap()
      .then(() => {
        notify({
          message: `Catalog Content was ${catalogContent ? 'enab' : 'disab'}led`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch(() => {
        handleError('There was an error sending your preference');
      });
  };

  const handleToggle = (name: string) => () => {
    if (name === 'anonymousUsageStats') {
      setAnonymousStats(!anonymousStats);
      handleChange(name, !anonymousStats);
    } else {
      setPerfResultStats(!perfResultStats);
      handleChange(name, !perfResultStats);
    }
  };

  const handleError = (msg: string) => {
    updateProgress({ showProgress: false });
    notify({ message: msg, event_type: EVENT_TYPES.ERROR });
  };

  const handleChange = (name: string, resultState: boolean) => {
    let val = resultState,
      msg;
    if (name === 'anonymousUsageStats') {
      msg = val
        ? 'Sending anonymous usage statistics was enabled'
        : 'Sending anonymous usage statistics was disabled';
    } else {
      msg = val
        ? 'Sending anonymous performance results was enabled'
        : 'Sending anonymous performance results was disabled';
    }

    const requestBody = JSON.stringify({
      anonymousUsageStats: name === 'anonymousUsageStats' ? val : anonymousStats,
      anonymousPerfResults: name === 'anonymousPerfResults' ? val : perfResultStats,
    });

    updateProgress({ showProgress: true });
    updateUserPrefWithContext({ body: requestBody })
      .unwrap()
      .then((result) => {
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          notify({ message: msg, event_type: val ? EVENT_TYPES.SUCCESS : EVENT_TYPES.INFO });
        }
      })
      .catch(() => {
        handleError('There was an error sending your preference');
      });
  };

  const handleUpdateUserPref = (key: string, value: unknown) => {
    if (!userData) return;
    const updates = _.set(_.cloneDeep(userData), key, value);
    updateUserPrefWithContext({ body: JSON.stringify(updates) });
  };

  useEffect(() => {
    if (capabilitiesRegistry && !capabilitiesLoaded) {
      setCapabilitiesLoaded(true);
      setUserPrefs(
        ExtensionPointSchemaValidator('user_prefs')(capabilitiesRegistry?.extensions?.user_prefs),
      );
      setProviderType(capabilitiesRegistry?.provider_type);
    }
  }, [capabilitiesRegistry]);

  useEffect(() => {
    if (isUserDataFetched && userData) {
      setExtensionPreferences(userData?.usersExtensionPreferences);
      setCatalogContent(userData?.usersExtensionPreferences?.catalogContent);
    } else if (isUserDataError) {
      console.error(userDataError);
    }
  }, [isUserDataFetched, userData]);

  useEffect(() => {
    if (isCapabilitiesDataFetched && capabilitiesData) {
      setProviderInfo(capabilitiesData);
    }
  }, [isCapabilitiesDataFetched, capabilitiesData]);

  return {
    anonymousStats,
    perfResultStats,
    catalogContent,
    userPrefs,
    providerType,
    providerInfo,
    userData,
    handleCatalogContentToggle,
    handleToggle,
    handleUpdateUserPref,
  };
};
