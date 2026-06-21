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

  const profileUrl = getProfileUrl();

  useEffect(() => {
    if (!userLoaded && isGetUserSuccess) {
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
  const refURL = btoa(window.location.href);

  if (userData?.status == 'anonymous') {
    const url = `${providerCapabilities?.providerUrl}?anonymousUserID=${userData?.id}&source=${sourceURL}&ref=${refURL}`;

    return (
      <Link href={url}>
        <Button variant="contained" color="primary" data-testid="sign-in-button">
          Sign In
        </Button>
      </Link>
    );
  }

  const avatar = (
    <Avatar
      sx={{ height: 36, width: 36 }}
      src={isGetUserSuccess ? userData?.avatarUrl : null}
      imgProps={{ referrerPolicy: 'no-referrer' }}
    />
  );

  const profileLinkProps = profileUrl
    ? {
        component: Link as any,
        href: profileUrl,
        target: '_blank',
        rel: 'noopener noreferrer',
        'aria-label': 'Open user profile',
      }
    : {
        disabled: true,
      };

  return (
    <div>
      <NoSsr>
        <div data-testid="profile-button">
          {
            <IconButtonAvatar color={color} {...profileLinkProps}>
              {avatar}
            </IconButtonAvatar>
          }
        </div>
      </NoSsr>
    </div>
  );
};

const UserProvider = (props) => {
  return <User {...props} />;
};

export default UserProvider;
