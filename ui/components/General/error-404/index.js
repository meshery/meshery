import { Typography } from '@material-ui/core';
import { ErrorTypes } from '@/constants/common';
// import InstallMeshery, { MesheryAction } from "../../dashboard/install-meshery-card";
import Socials from './socials';
import {
  ErrorComponent,
  ErrorContainer,
  ErrorContentContainer,
  ErrorLink,
  ErrorMain,
} from './styles';

//TODO: Add component for meshery version compatiblity error
// const MesheryVersionCompatiblity = () => {
//   return (
//     <div>
//       <Typography variant="p" component="p" align="center">
//         <InstallMeshery action={MesheryAction.UPGRADE.KEY} />
//       </Typography>
//     </div>
//   );
// };

const UnknownServerSideError = (props) => {
  const { errorContent } = props;
  return (
    <div>
      <ErrorContentContainer>
        <Typography variant="p" component="p" align="center">
          {errorContent}
        </Typography>
      </ErrorContentContainer>
    </div>
  );
};

const DefaultError = (props) => {
  const { errorTitle, errorContent, errorType } = props;

  return (
    <ErrorMain>
      <ErrorContainer>
        <img
          width="400px"
          height="300px"
          src="/static/img/meshery-logo/meshery-logo-light-text.png"
          alt="Meshery logo"
        />
        <ErrorComponent>
          <Typography variant="h4" component="h4" align="center" className="errormsg">
            {errorTitle
              ? errorTitle
              : 'You are not authorized to view this page. Contact your administrator'}
          </Typography>
          {errorType === ErrorTypes.UNKNOWN ? (
            <UnknownServerSideError errorContent={errorContent} />
          ) : null}
          <div style={{ marginTop: '3rem' }}>
            <div>
              <Typography variant="p" component="p" align="center">
                Navigate to <ErrorLink href="/">Dashboard</ErrorLink>
              </Typography>
            </div>
            <div style={{ marginTop: '0.8rem' }}>
              <Typography variant="p" component="p" align="center">
                For help, please inquire on the
                <ErrorLink href="https://discuss.layer5.io"> discussion forum</ErrorLink> or the{' '}
                <ErrorLink href="https://slack.layer5.io"> Slack workspace</ErrorLink>.
              </Typography>
            </div>
          </div>
        </ErrorComponent>
        <Socials />
      </ErrorContainer>
    </ErrorMain>
  );
};

export default DefaultError;
