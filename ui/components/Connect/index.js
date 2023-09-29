import { useRouter } from "next/router";
import React, { useEffect, useState } from "react";
import ConnectionStepper from "./stepper";
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateProgress } from '../../lib/store';
// import Verify from "./verify";

const ConnectionWizardNew = ({user}) => {
  const router = useRouter();
  const [connection, setConnection] = useState({
    connectionType: "", // github, kubernetes, etc
    connectionAction: "", // new, callback
    connectionId: "" // id of the connection, rn it is just the user uuid
  });
  console.log("user from connection wizard", user)

  useEffect(() => {
    const { connection, id, githubLogin } = router.query;
    console.log("connection", connection);
    if (connection) {
      const newConnection = {
        connectionType: connection[0],
        connectionAction: connection[1],
        connectionId: connection[2]
      };

      // Check if the new connection is different from the current connection
      if (
        newConnection.connectionType !== connection.connectionType ||
        newConnection.connectionAction !== connection.connectionAction ||
        newConnection.connectionId !== connection.connectionId
      ) {
        setConnection(newConnection);

        const shouldRedirect =
          newConnection.connectionAction === "new" &&
          user &&
          (!newConnection.connectionType ||
            !newConnection.connectionAction ||
            !newConnection.connectionId);
// profileData.id !== newConnection.connectionId
        let queryString = "";
        if (githubLogin) {
          queryString = `?githubLogin=${githubLogin}`;
          if (id) {
            queryString = `${queryString}&id=${id}`;
          }
        }
        console.log("shouldRedirect", shouldRedirect);
        // if (shouldRedirect && newConnection.connectionType !== "helm") {
        //   const redirectUrl = `${GITHUB_NEW_CONNECTION}/${profileData.id}${queryString}`;
        //   router.replace(redirectUrl);
        // }
        // if (shouldRedirect) {
        //   const redirectUrl = `${GITHUB_NEW_CONNECTION}/${newConnection.connectionType}/${newConnection.connectionAction}/${newConnection.connectionId}${queryString}`;
        //   router.replace(redirectUrl);
        // }
      }
    }
  }, [router.query]);

  return (
    <>
      {connection.connectionAction === "new" && (
        <ConnectionStepper connectionType={connection.connectionType} />
      )}
      {/* {connection.connectionAction === "verify" && <Verify />} TODO: enable if required */}
    </>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  return {
    user: state.get('user')?.toObject()
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(ConnectionWizardNew);
