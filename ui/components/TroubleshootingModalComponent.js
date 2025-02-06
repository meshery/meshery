import * as React from 'react';
import { Modal } from '@mui/material';
import { keyframes } from '@mui/material/styles';

import {
  Typography,
  Accordion,
  AccordionDetails,
  AccordionSummary,
  styled,
  Paper,
  IconButton,
  InfoIcon,
  LIGHT_TEAL,
} from '@layer5/sistent';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import CloseIcon from '@mui/icons-material/Close';
import { UsesSistent } from './SistentWrapper';

const StyledPaper = styled(Paper)(({ theme }) => ({
  top: '50%',
  left: '50%',
  width: '47.5%',
  position: 'absolute',
  borderRadius: '10px',
  overflow: 'hidden',
  transform: 'translate(-50%, -50%)',
  backgroundColor: 'transparent',
  outline: 'none',
  [theme.breakpoints.down(1350)]: {
    width: '70%',
  },
}));

const HeaderContainer = styled('div')(() => ({
  display: 'flex',
  alignItems: 'center',
}));

const HeaderText = styled(Typography)(() => ({
  flexShrink: 0,
  margin: 10,
  fontSize: '1rem',
}));

const Info = styled(InfoIcon)(({ theme }) => ({
  color: theme.palette.type === 'dark' ? theme.palette.background.brand?.default : LIGHT_TEAL,
}));

const FooterText = styled(Typography)(() => ({
  color: 'white',
  fontSize: '.85rem',
  textDecoration: 'italics',
  fontFamily: 'Qanelas Soft, sans-serif',
}));

const ModalHeader = styled(Typography)(({ theme }) => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  paddingBottom: 10,
  padding: '0.5rem',
  paddingTop: 10,
  backgroundColor: theme.palette.background.brand.default,
}));

const ModalFooter = styled('div')(({ theme }) => ({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  paddingBottom: 10,
  padding: '0.5rem',
  paddingTop: 10,
  backgroundColor: theme.palette.background.brand.default,
}));

const rotateCloseIcon = keyframes`
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
`;

const ModalContent = styled('div')({
  backgroundColor: 'transparent',
});

const Details = styled('div')({
  display: 'block',
});

const ModelHeader = styled(Typography)({
  fontSize: '1rem',
  color: 'white',
});

const IconStyle = styled(CloseIcon)({
  color: 'white',
});

const IconContainer = styled(IconButton)({
  transition: 'all .3s',
  '&:hover': {
    backgroundColor: 'transparent !important',
    animation: `${rotateCloseIcon} 1s`,
  },
});

const AccordionContainer = styled(Accordion)({
  margin: '0 !important',
  borderBottom: '1px solid #ccc',
});

const AccordionSummaryStyled = styled(AccordionSummary)(({ theme }) => ({
  '&.Mui-expanded': {
    backgroundColor: theme.palette.background.card,
  },
}));

const AccDetailHead = styled(Typography)(({ theme }) => ({
  color: theme.palette.mode === 'dark' ? '#f1f1f1' : '#444',
  fontFamily: 'Qanelas Soft, sans-serif',
}));

const TroubleshootListitem = styled('li')({
  fontSize: '0.9rem',
  marginBottom: '1rem',
});

const KeyStyleContainer = styled('div')(({ theme }) => ({
  display: 'inline-block',
  padding: '0.1rem 0.5rem',
  background: theme.palette.background.tabs,
  margin: '0.3rem',
  borderRadius: '5px',
  boxShadow:
    'rgba(0, 0, 0, 0.17) 0px -0px 0px -5px inset, rgba(0, 0, 0, 0.15) 0px 0px 0px -3px inset, rgba(0, 0, 0, 0.1) 0px 4px 30px 0px inset, rgba(0, 0, 0, 0.06) 0px 2px 1px, rgba(0, 0, 0, 0.09) 0px 0px 2px, rgba(0, 0, 0, 0.09) 0px 8px 4px, rgba(0, 0, 0, 0.1) 0px 1px 0px, rgba(0, 0, 0, 0.1) 0px -2px 0px',
}));

const AccDetailsContainer = styled(AccordionDetails)({
  flexDirection: 'column',
});

const TroubleshootHelpLink = styled('a')(({ theme }) => ({
  color:
    theme.palette.mode === 'dark' ? theme.palette.background.brand?.default : 'rgb(57, 102, 121)',
  fontWeight: 'bold',
  textDecoration: 'none',
}));

const ContactHelpLink = styled('a')({
  color: 'white',
  fontWeight: 'bold',
  textDecoration: 'none',
});

const TroubleshootingModal = (props) => {
  const [expanded, setExpanded] = React.useState(false);
  const [missingData, setMissingData] = React.useState(
    !!props.viewHeaderErrorMessage?.includes('data missing'),
  );
  const handleChange = (panel) => (event, isExpanded) => {
    setExpanded(isExpanded ? panel : false);
    if (missingData) {
      setMissingData(false);
    }
  };

  const handleClose = () => props?.setOpen(false);

  return (
    <UsesSistent>
      <Modal
        open={props?.open}
        onClose={handleClose}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
      >
        <StyledPaper elevation={6} square="true">
          <ModalContent>
            <ModalHeader>
              <Typography variant="h5"></Typography>
              <ModelHeader variant="h5">Extensions Troubleshooting Guide</ModelHeader>
              <IconContainer data-cy="modal-close-btn" onClick={handleClose}>
                <IconStyle />
              </IconContainer>
            </ModalHeader>
            <AccordionContainer expanded={expanded === 'panel1'} onChange={handleChange('panel1')}>
              <AccordionSummaryStyled
                expandIcon={<ExpandMoreIcon />}
                aria-controls="panel1bh-content"
                id="panel1bh-header"
              >
                <HeaderContainer>
                  <Info />
                  <HeaderText variant="h6">Stale Data</HeaderText>
                </HeaderContainer>
              </AccordionSummaryStyled>
              <Details>
                <AccDetailHead>
                  <strong>Browser</strong>
                </AccDetailHead>
                <ul>
                  <Typography>
                    <TroubleshootListitem>
                      Run <KeyStyleContainer>CTRL</KeyStyleContainer> +
                      <KeyStyleContainer>SHIFT</KeyStyleContainer>+
                      <KeyStyleContainer>R</KeyStyleContainer> or
                      <KeyStyleContainer>CMD</KeyStyleContainer>+
                      <KeyStyleContainer>OPTION</KeyStyleContainer>+
                      <KeyStyleContainer>E</KeyStyleContainer> to force reload if you are getting
                      stale copy of the component due to caching.
                    </TroubleshootListitem>
                    <TroubleshootListitem>
                      Use Incognito or Private browsing mode. If you are still getting stale copy of
                      the component, try opening a <code>incognito tab</code>
                    </TroubleshootListitem>
                  </Typography>
                </ul>
              </Details>
            </AccordionContainer>
            <AccordionContainer
              expanded={expanded === 'panel2' || missingData}
              onChange={handleChange('panel2')}
            >
              <AccordionSummaryStyled
                expandIcon={<ExpandMoreIcon />}
                aria-controls="panel2bh-content"
                id="panel2bh-header"
              >
                <HeaderContainer>
                  <Info />
                  <HeaderText variant="h6">Missing Data</HeaderText>
                </HeaderContainer>
              </AccordionSummaryStyled>
              <AccDetailsContainer>
                <AccDetailHead>
                  <strong>Meshery Database</strong>
                </AccDetailHead>
                <ul>
                  <Typography>
                    <TroubleshootListitem>
                      Verify MeshSync data is being received. Run{' '}
                      <code>kubectl get svc -n meshery</code>. Docker Desktop: VPNkit commonly fails
                      to assign an IP address to Meshery Broker (MeshSync). Verify that the Meshery
                      Broker service has external IP address assigned.
                    </TroubleshootListitem>
                    <TroubleshootListitem>
                      Confirm that your machine&apos;s firewall isn&apos;t getting in the way.
                    </TroubleshootListitem>
                    <TroubleshootListitem>
                      Dump Meshery Database. Run <code>rm -rf ~/.meshery/config</code>.
                    </TroubleshootListitem>
                  </Typography>
                </ul>
              </AccDetailsContainer>
            </AccordionContainer>

            <AccordionContainer expanded={expanded === 'panel4'} onChange={handleChange('panel4')}>
              <AccordionSummaryStyled
                expandIcon={<ExpandMoreIcon />}
                aria-controls="panel4bh-content"
                id="panel4bh-header"
              >
                <HeaderContainer>
                  <Info />
                  <HeaderText variant="h6">Additional Resources</HeaderText>
                </HeaderContainer>
              </AccordionSummaryStyled>
              <AccDetailsContainer>
                <HeaderText variant="h6">
                  <strong>Troubleshooting Tips</strong>
                </HeaderText>
                <ul>
                  <Typography>
                    <TroubleshootListitem>
                      {' '}
                      <TroubleshootHelpLink
                        href="https://meshery.io/community#community-forums/t/what-are-some-troubleshooting-tips-for-meshmap"
                        target="_blank"
                        rel="noreferrer"
                      >
                        &quot;What are some troubleshooting tips for MeshMap?&quot;
                      </TroubleshootHelpLink>
                    </TroubleshootListitem>
                  </Typography>
                </ul>
              </AccDetailsContainer>
            </AccordionContainer>
            <ModalFooter>
              <FooterText variant="h6">
                Need help? Contact us via{' '}
                <ContactHelpLink
                  href="mailto:maintainers@meshery.io"
                  target="_blank"
                  rel="noreferrer"
                >
                  email
                </ContactHelpLink>{' '}
                or{' '}
                <ContactHelpLink href="https://meshery.io/community#community-forums" target="_blank" rel="noreferrer">
                  community forum
                </ContactHelpLink>
                .
              </FooterText>
            </ModalFooter>
          </ModalContent>
        </StyledPaper>
      </Modal>
    </UsesSistent>
  );
};

export default TroubleshootingModal;
