import React, { useState, useRef, useEffect } from 'react';
import { Avatar, Button } from '@sistent/sistent';
import NoSsr from '@mui/material/NoSsr';
import Link from 'next/link';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { IconButtonAvatar } from './Header.styles';
import { useDispatch, useSelector } from 'react-redux';
import { updateUser } from '@/store/slices/mesheryUi';
/**
 * Extension Point: Avatar behavior for User Modes
 * Insert custom logic here to handle Single User mode, Anonymous User mode, Multi User mode behavior.
 */
type UserProps = {
  color?:
    | 'default'
    | 'inherit'
    | 'primary'
    | 'secondary'
    | 'error'
    | 'info'
    | 'success'
    | 'warning';
  [key: string]: any;
};

type AccountItem = {
  title: string;
  href: string;
  [key: string]: any;
};

const User = (props: UserProps) => {
  const [userLoaded, setUserLoaded] = useState(false);
  const [account, setAccount] = useState<AccountItem[]>([]);
  const capabilitiesLoadedRef = useRef(false);
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const { capabilitiesRegistry } = useSelector((state: any) => state.ui);
  const {
    data: userData,
    isSuccess: isGetUserSuccess,
    isError: isGetUserError,
    error: getUserError,
  } = useGetLoggedInUserQuery(undefined, {});

  const getProfileUrl = () => {
    return (account || [])?.find((item) => item.title === 'Cloud Account')?.href;
  };

  const goToProfile = () => {
    const profileUrl = getProfileUrl();
    if (profileUrl) {
      window.location.href = profileUrl;
      return;
    }
  };

  useEffect(() => {
    if (!userLoaded && isGetUserSuccess) {
      dispatch(updateUser({ user: userData }));
      setUserLoaded(true);
    } else if (isGetUserError) {
      notify({
        message: 'Error fetching user',
        event_type: EVENT_TYPES.ERROR,
        details: (getUserError as any)?.data,
      });
    }
  }, [userData, isGetUserSuccess, isGetUserError]);

  useEffect(() => {
    if (!capabilitiesLoadedRef.current && capabilitiesRegistry) {
      capabilitiesLoadedRef.current = true;
      setAccount(
        ExtensionPointSchemaValidator('account')(capabilitiesRegistry?.extensions?.account),
      );
    }
  }, [capabilitiesRegistry]);

  const { color } = props;

  const source = new URL('/api/user/token', window.location.origin);
  const sourceURL = btoa(source.toString());
  const refURL = btoa(window.location.href);

  if (userData?.status == 'anonymous') {
    const url = `${capabilitiesRegistry?.provider_url}?anonymousUserID=${userData?.id}&source=${sourceURL}&ref=${refURL}`;

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
          <IconButtonAvatar
            {...(color ? { color } : {})}
            aria-haspopup="true"
            onClick={goToProfile}
          >
            <Avatar
              sx={{ height: 36, width: 36 }}
              src={isGetUserSuccess ? userData?.avatar_url : null}
              imgProps={{ referrerPolicy: 'no-referrer' }}
            />
          </IconButtonAvatar>
        </div>
      </NoSsr>
    </div>
  );
};

type UserProviderProps = {
  color?:
    | 'default'
    | 'inherit'
    | 'primary'
    | 'secondary'
    | 'error'
    | 'info'
    | 'success'
    | 'warning';
  [key: string]: any;
};

const UserProvider = (props: UserProviderProps) => {
  return <User {...props} />;
};

export default UserProvider;
