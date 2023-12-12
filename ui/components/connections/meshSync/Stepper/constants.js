import FinishFlagIcon from '../../../../assets/icons/FinishFlagIcon';
import { ConnectionDetails, CredentialDetails, Finish } from './StepperContent';
import ConnectionIcon from '../../../../assets/icons/Connection';
import CredentialIcon from '../../../../assets/icons/CredentialIcon';

export const registerConnectionSteps = ['Connection Details', 'Credential Details', 'Finish'];

export const registerConnectionIcons = {
  1: <ConnectionIcon height={26} width={26} />,
  2: <CredentialIcon height={32} width={32} />,
  3: <FinishFlagIcon width="2rem" />,
};

export const registerConnectionContent = {
  1: {
    component: <ConnectionDetails />,
    props: ['handleNext', 'sharedData', 'setSharedData'],
  },
  2: {
    component: <CredentialDetails />,
    props: ['handleNext', 'sharedData', 'setSharedData'],
  },
  3: {
    component: <Finish />,
    props: ['sharedData', 'setSharedData'],
  },
};

export const ConnectionStepperTips = [
  'GitOps is more fun with friends. Invite one today.',
  'Meshery works with multiple Kubernetes clusters. Connect each of yours by uploading your kubeconfig and letting Meshery do the rest.',
  'Import your Docker Compose, Helm Chart, and Kubernetes mainfests. Let Meshery visualize and deploy them.',
];

export const ConnectionDetailContent = {
  title: '',
  tips: 'Establish a connection by registering this resource.',
  btnText: 'Next',
  cancel: false,
};

export const CredentialDetailContent = {
  title: '',
  tips: '',
  btnText: 'Verify Connection',
  cancel: true,
};

export const FinishContent = {
  title: 'Connection Created',
  subtitle: 'Congratulations 🎉, you have registered a new connection.',
  btnText: 'Finish',
};
