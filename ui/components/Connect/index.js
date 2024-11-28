import { useRouter } from 'next/router';
import React, { useEffect, useState } from 'react';
import ConnectionStepper from './stepper';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateProgress } from '../../lib/store';
import { Provider } from 'react-redux';
import NoSsr from '@material-ui/core/NoSsr';
import { store } from '../../store';
import { ErrorBoundary } from '../General/ErrorBoundary';
// import Verify from "./verify";

const ConnectionWizard = ({ user }) => {
  const router = useRouter();
  const [connection, setConnection] = useState({
    connectionType: '', // github, kubernetes, etc
    connectionAction: '', // new, callback
    connectionId: '', // id of the connection, rn it is just the user uuid
  });

  useEffect(() => {
    const { connection } = router.query;
    if (connection) {
      const newConnection = {
        connectionType: connection[0],
        connectionAction: connection[1],
        connectionId: connection[2],
      };

      // Check if the new connection is different from the current connection
      if (
        newConnection.connectionType !== connection.connectionType ||
        newConnection.connectionAction !== connection.connectionAction ||
        newConnection.connectionId !== connection.connectionId
      ) {
        setConnection(newConnection);

        const shouldRedirect =
          newConnection.connectionAction === 'new' &&
          user &&
          (!newConnection.connectionType ||
            !newConnection.connectionAction ||
            !newConnection.connectionId);
        // profileData.id !== newConnection.connectionId
        // let queryString = '';
        // if (githubLogin) {
        //   queryString = `?githubLogin=${githubLogin}`;
        //   if (id) {
        //     queryString = `${queryString}&id=${id}`;
        //   }
        // }
        console.log('shouldRedirect', shouldRedirect);
      }
    }
  }, [router.query]);

  return (
    <>
      {connection.connectionAction === 'new' && (
        <ConnectionStepper connectionType={connection.connectionType} />
      )}
    </>
  );
};

const ConnectionWizardNew = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary
        FallbackComponent={() => null}
        onError={(e) => console.error('Error in connection wizard', e)}
      >
        <Provider store={store}>
          <ConnectionWizard {...props} />
        </Provider>
      </ErrorBoundary>
    </NoSsr>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  return {
    user: state.get('user')?.toObject(),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(ConnectionWizardNew);
