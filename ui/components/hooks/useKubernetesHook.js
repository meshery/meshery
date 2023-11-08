import { useDispatch } from 'react-redux';
import { updateProgress } from '../../lib/store';
import { useNotification } from '../../utils/hooks/useNotification';
import { errorHandlerGenerator, successHandlerGenerator } from '../ConnectionWizard/helpers/common';
import { pingKubernetes } from '../ConnectionWizard/helpers/kubernetesHelpers';

export default function useKubernetesHook() {
  const { notify } = useNotification();
  const dispatch = useDispatch();
  const ping = (name, server, connectionID) => {
    dispatch(updateProgress({ showProgress: true }));
    pingKubernetes(
      successHandlerGenerator(notify, `Kubernetes context ${name} at ${server} pinged`, () =>
        dispatch(updateProgress({ showProgress: false })),
      ),
      errorHandlerGenerator(notify, `Kubernetes context ${name} at ${server} not reachable`, () =>
        dispatch(updateProgress({ showProgress: false })),
      ),
      connectionID,
    );
  };

  return ping;
}
