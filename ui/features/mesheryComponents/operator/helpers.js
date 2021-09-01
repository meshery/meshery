import fetchMesheryOperatorStatus from "../graphql/queries/OperatorStatusQuery";
import subscribeOperatorStatusEvents from "../graphql/subscriptions/OperatorStatusEvents";
import changeOperatorStatusMutation from "../graphql/mutations/OperatorStatus";

export const fetchMesheryComponentsStatus = () =>
  new Promise((res, rej) => {
    fetchMesheryOperatorStatus().subscribe({
      next: (result) => {
        res(result);
      },
      error: (err) => rej(err),
    });
  });

/**
 *
 * @param {(res) => void} cb callback that dispatches `updateConnectoinStatus` action
 * @returns {Promise}
 */
export const initialiseOperatorStatusEventsSubscriptions = (cb) =>
  // have to rethink this implementation, in this implementation, the error case might not be handled properly

  new Promise((res, rej) => {
    subscribeOperatorStatusEvents(cb, () =>
      rej("Subscription not initialised succesfully")
    );
    res("Initialised succesfully");
  });

/**
 *
 * @typedef  OperatorStateChangeResponse
 * @property {Object} response
 * @property {Object} errors
 */

/**
 *
 * @param {"ENABLED" | "DISABLED"} desiredState
 * @returns {Promise<OperatorStateChangeResponse>}
 */

export const changeOperatorStatus = (desiredState) =>
  new Promise((res, rej) => {
    changeOperatorStatusMutation(
      (response, errors) =>
        res({
          response,
          errors,
        }),
      (err) => rej(err),
      { status: desiredState }
    );
  });
