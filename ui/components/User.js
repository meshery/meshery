import React, { useState, useRef } from 'react';
import { Avatar } from '@layer5/sistent';
import IconButton from '@material-ui/core/IconButton';
import NoSsr from '@material-ui/core/NoSsr';
import Link from 'next/link';
import { Provider, connect } from 'react-redux';
import { store } from '../store';
import { bindActionCreators } from 'redux';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import { updateUser } from '../lib/store';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { Button } from '@material-ui/core';

/**
 * Extension Point: Avatar behavior for User Modes
 * Insert custom logic here to handle Single User mode, Anonymous User mode, Multi User mode behavior.
 */
const User = (props) => {
  const [userLoaded, setUserLoaded] = useState(false);
  const [account, setAccount] = useState([]);
  const capabilitiesLoadedRef = useRef(false);
  const { notify } = useNotification();

  const {
    data: userData,
    isSuccess: isGetUserSuccess,
    isError: isGetUserError,
    error: getUserError,
  } = useGetLoggedInUserQuery();

  const { capabilitiesRegistry } = props;

  const getProfileUrl = () => {
    return (account || [])?.find((item) => item.title === 'Profile')?.href;
  };

  const goToProfile = () => {
    const profileUrl = getProfileUrl();
    if (profileUrl) {
      window.location = profileUrl;
      return;
    }
  };

  if (!userLoaded && isGetUserSuccess) {
    props.updateUser({ user: userData });
    setUserLoaded(true);
  } else if (isGetUserError) {
    notify({
      message: 'Error fetching user',
      event_type: EVENT_TYPES.ERROR,
      details: getUserError?.data,
    });
  }

  if (!capabilitiesLoadedRef.current && capabilitiesRegistry) {
    capabilitiesLoadedRef.current = true;
    setAccount(ExtensionPointSchemaValidator('account')(capabilitiesRegistry?.extensions?.account));
  }

  const { color, iconButtonClassName, avatarClassName } = props;

  const source = new URL('/api/user/token', window.location.origin);
  const sourceURL = btoa(source.toString());
  const refURL = btoa(window.location.href);

  if (userData?.status == 'anonymous') {
    const url = `${capabilitiesRegistry?.provider_url}?anonymousUserID=${userData?.id}&source=${sourceURL}&ref=${refURL}`;

    return (
      <Link href={url}>
        <Button variant="contained" color="primary">
          Sign In
        </Button>
      </Link>
    );
  }

  return (
    <div>
      <NoSsr>
        <div data-test="profile-button">
          <IconButton
            color={color}
            className={iconButtonClassName}
            aria-haspopup="true"
            onClick={goToProfile}
          >
            <Avatar
              className={avatarClassName}
              src={isGetUserSuccess ? userData?.avatar_url : null}
              imgProps={{ referrerPolicy: 'no-referrer' }}
            />
          </IconButton>
        </div>
      </NoSsr>
    </div>
  );
};

const UserProvider = (props) => {
  return (
    <Provider store={store}>
      <User {...props} />
    </Provider>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateUser: bindActionCreators(updateUser, dispatch),
});

const mapStateToProps = (state) => ({
  capabilitiesRegistry: state.get('capabilitiesRegistry'),
});

export default connect(mapStateToProps, mapDispatchToProps)(UserProvider);
