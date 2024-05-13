import { ErrorTypes } from '@/constants/common';
import { useTheme } from '@material-ui/core/styles';
import Socials from './socials';
import {
  ErrorSection,
  ErrorSectionContainer,
  ErrorContainer,
  ErrorContentContainer,
  ErrorMain,
  ErrorSectionContent,
  StyledButton,
} from './styles';
import { Typography, Divider, InfoCircleIcon } from '@layer5/sistent';
import OrgSwitcher from './OrgSwitcher';
import RequestForm from './RequestForm';
import CurrentSessionInfo from './CurrentSession';
import { Tooltip } from '@mui/material';

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
  // const { roles } = useGetUserRolesQuery({});
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
        <ErrorSectionContainer>
          <ErrorSection>
            <Typography variant="h5" component="h5" align="center" fontWeight={600}>
              YOUR CURRENT SESSION
            </Typography>
            <CurrentSessionInfo />
          </ErrorSection>
          <Divider
            orientation="vertical"
            flexItem
            sx={{
              display: {
                xs: 'none',
                lg: 'block',
              },
            }}
          />
          <ErrorSection>
            <Typography variant="h5" component="h5" align="center" fontWeight={600}>
              YOUR OPTIONS
            </Typography>
            <ErrorSectionContent>
              <OrgSwitcher />
              <Divider />
              <RequestForm />
            </ErrorSectionContent>
          </ErrorSection>
        </ErrorSectionContainer>
        <Tooltip title="This is a tooltip">
          <InfoCircleIcon
            height={32}
            width={32}
            style={{
              alignSelf: 'flex-end',
              marginInline: '2rem',
            }}
          />
        </Tooltip>
      </ErrorContainer>
      <StyledButton variant="contained">Return to Dashboard</StyledButton>
      <Socials />
    </ErrorMain>
  );
};

export default DefaultError;
