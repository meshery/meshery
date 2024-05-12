import { Typography } from '@material-ui/core';
import { ErrorTypes } from '@/constants/common';
import { useTheme } from '@material-ui/core/styles';
// import InstallMeshery, { MesheryAction } from "../../dashboard/install-meshery-card";
import Socials from './socials';
import {
  ErrorComponent,
  ErrorContainer,
  ErrorContentContainer,
  // ErrorLink,
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
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
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
        <Typography variant="h4" component="h4" align="center" className="errormsg">
          {errorTitle
            ? errorTitle
            : "Oops! It seems like you don't have the necessary permissions to view this page."}
        </Typography>
        {errorType === ErrorTypes.UNKNOWN ? (
          <UnknownServerSideError errorContent={errorContent} />
        ) : null}
      </div>
      <ErrorContainer>
        <ErrorComponent></ErrorComponent>
      </ErrorContainer>
      <Socials />
    </ErrorMain>
  );
};

export default DefaultError;
