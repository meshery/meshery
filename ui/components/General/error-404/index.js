import { Typography } from '@material-ui/core';
import { ErrorTypes } from '@/constants/common';
import { useTheme } from '@material-ui/core/styles';
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
  const theme = useTheme();

  return (
    <ErrorMain>
      <ErrorContainer>
        <img
          width="400px"
          height="300px"
          src={
            theme.palette.type === 'dark'
              ? '/static/img/meshery-logo/meshery-logo-white-text.png'
              : '/static/img/meshery-logo/meshery-logo-light-text.png'
          }
          alt="Meshery logo"
        />
        <ErrorComponent>
          <Typography variant="h4" component="h4" align="center" className="errormsg">
            {errorTitle
              ? errorTitle
              : "Oops! It seems like you don't have the necessary permissions to view this page."}
          </Typography>
          <Typography
            variant="body1"
            component="p"
            align="left"
            style={{ paddingLeft: '1.25rem', paddingTop: '1rem' }}
          >
            Possible reasons:
          </Typography>
          <ol style={{ textAlign: 'left', fontSize: '1rem' }}>
            <li>
              <strong>Insufficient Permissions:</strong> Your account may lack the required
              permissions to access this page. To resolve this, please reach out to your
              administrator and request the necessary access.
            </li>
            <li>
              <strong>Check Selected Organization:</strong> Ensure you are in the correct
              organization in which you have access to view this page. You can verify your
              organization settings <ErrorLink href="/user/preferences">here</ErrorLink>. If needed,
              switch to an organization where you have the required permissions.
            </li>
          </ol>
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
