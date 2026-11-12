import React, { useState, useRef, useEffect } from 'react';
import { Avatar, Button, NoSsr } from '@sistent/sistent';
import Link from 'next/link';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { IconButtonAvatar } from './layout/Header/Header.styles';
import { useDispatch, useSelector } from 'react-redux';
import { updateUser } from '@/store/slices/mesheryUi';
/**
 * Extension Point: Avatar behavior for User Modes
 * Insert custom logic here to handle Single User mode, Anonymous User mode, Multi User mode behavior.
 */
const User = (props) => {
  const [userLoaded, setUserLoaded] = useState(false);
  const [account, setAccount] = useState([]);
  const capabilitiesLoadedRef = useRef(false);
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const { providerCapabilities } = useSelector((state) => state.ui);
  const {
    data: userData,
    isSuccess: isGetUserSuccess,
    isError: isGetUserError,
    error: getUserError,
  } = useGetLoggedInUserQuery();

  const getProfileUrl = () => {
    return (account || [])?.find((item) => item.title === 'Cloud Account')?.href;
  };

  const goToProfile = () => {
    const profileUrl = getProfileUrl();
    if (profileUrl) {
      window.open(profileUrl, '_blank', 'noopener,noreferrer');
      return;
    }
    notify({
      message: 'Please log in to access this profile',
      event_type: EVENT_TYPES.WARNING,
    });
  };

  useEffect(() => {
    if (!userLoaded && isGetUserSuccess) {
      // userData is normalized by getLoggedInUser's transformResponse
      // (userId backfilled from id for the v1beta3 Cloud response).
      dispatch(updateUser({ user: userData }));
      setUserLoaded(true);
    } else if (isGetUserError) {
      notify({
        message: 'Error fetching user',
        event_type: EVENT_TYPES.ERROR,
        details: getUserError?.data,
      });
    }
  }, [userData, isGetUserSuccess, isGetUserError]);

  useEffect(() => {
    if (!capabilitiesLoadedRef.current && providerCapabilities) {
      capabilitiesLoadedRef.current = true;
      setAccount(
        ExtensionPointSchemaValidator('account')(providerCapabilities?.extensions?.account),
      );
    }
  }, [providerCapabilities]);

  const { color } = props;

  const source = new URL('/api/user/token', window.location.origin);
  const sourceURL = btoa(source.toString());
  // Relative ref keeps search and hash (including mode=design). An absolute
  // href is a same-origin URL the server used to reject, dropping the user
  // on /extension/meshmap with the query stripped.
  const refURL = btoa(
    `${window.location.pathname}${window.location.search}${window.location.hash}`,
  );

  if (userData?.status == 'anonymous') {
    // btoa emits standard base64, whose '+' decodes back to a space in a query
    // value. URLSearchParams percent-encodes every value, so the server reads
    // the ref and source it was sent.
    const params = new URLSearchParams({
      anonymousUserID: userData?.id ?? '',
      source: sourceURL,
      ref: refURL,
    });
    const url = `${providerCapabilities?.providerUrl}?${params.toString()}`;

    return (
      <Link href={url}>
        <Button variant="contained" color="primary" data-testid="sign-in-button">
          Sign In
        </Button>
      </Link>
    );
  }

  return (
    <div>
      <NoSsr>
        <div data-testid="profile-button">
          <IconButtonAvatar color={color} aria-haspopup="true" onClick={goToProfile}>
            <Avatar
              sx={{ height: 36, width: 36 }}
              src={isGetUserSuccess ? userData?.avatarUrl : null}
              imgProps={{ referrerPolicy: 'no-referrer' }}
            />
          </IconButtonAvatar>
        </div>
      </NoSsr>
    </div>
  );
};

const UserProvider = (props) => {
  return <User {...props} />;
};

export default UserProvider;
