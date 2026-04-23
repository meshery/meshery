import * as React from 'react';
import {
  Typography,
  Accordion,
  AccordionDetails,
  AccordionSummary,
  styled,
  InfoIcon,
  LIGHT_TEAL,
  useTheme,
} from '@sistent/sistent';
import { ExpandMore as ExpandMoreIcon, Information } from '../icons';

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
  color: theme.palette.mode === 'dark' ? theme.palette.background.brand?.default : LIGHT_TEAL,
}));

const FooterText = styled(Typography)(({ theme }) => ({
  color: theme.palette.common.white,
  fontSize: '.85rem',
  fontStyle: 'italic',
}));

const AccordionContainer = styled(Accordion)(({ theme }) => ({
  margin: '0 !important',
  boxShadow: 'none',
  borderBottom: `1px solid ${theme.palette.divider}`,
  '&:before': {
    display: 'none',
  },
});

const AccordionSummaryStyled = styled(AccordionSummary)(({ theme }) => ({
  '&.Mui-expanded': {
    backgroundColor: theme.palette.background.card,
  },
}));

const AccDetailHead = styled(Typography)(({ theme }) => ({
  color: theme.palette.text.primary,
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
    theme.palette.mode === 'dark' ? theme.palette.background.brand?.default : theme.palette.primary.main,
  fontWeight: 'bold',
  textDecoration: 'none',
}));

const ContactHelpLink = styled('a')(({ theme }) => ({
  color: theme.palette.common.white,
  fontWeight: 'bold',
  textDecoration: 'none',
}));

const TroubleshootingModal = (props) => {
  const { open, setOpen, viewHeaderErrorMessage } = props;
  const [expanded, setExpanded] = React.useState<string | false>(false);
  const [missingData, setMissingData] = React.useState(
    !!viewHeaderErrorMessage?.includes('data missing'),
  );
  const theme = useTheme();

  const handleChange = (panel: string) => (event: React.SyntheticEvent, isExpanded: boolean) => {
    setExpanded(isExpanded ? panel : false);
    if (missingData) {
      setMissingData(false);
    }
  };

  const handleClose = () => setOpen(false);

  return (
    <Modal
      open={open}
      closeModal={handleClose}
      title="Extensions Troubleshooting Guide"
      headerIcon={<Information fill={theme.palette.common.white} height="2rem" width="2rem" />}
      maxWidth="md"
      footer={
        <FooterText>
          Need help? Contact us via{' '}
          <ContactHelpLink href="mailto:maintainers@meshery.io" target="_blank" rel="noreferrer">
            email
          </ContactHelpLink>{' '}
          or{' '}
          <ContactHelpLink
            href="https://meshery.io/community#community-forums"
            target="_blank"
            rel="noreferrer"
          >
            community forum
          </ContactHelpLink>
          .
        </FooterText>
      }
    >
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
        <AccDetailsContainer>
          <AccDetailHead>
            <strong>Browser</strong>
          </AccDetailHead>
          <ul>
            <Typography component="div">
              <TroubleshootListitem>
                Run <KeyStyleContainer>CTRL</KeyStyleContainer> +
                <KeyStyleContainer>SHIFT</KeyStyleContainer>+
                <KeyStyleContainer>R</KeyStyleContainer> or
                <KeyStyleContainer>CMD</KeyStyleContainer>+
                <KeyStyleContainer>OPTION</KeyStyleContainer>+
                <KeyStyleContainer>E</KeyStyleContainer> to force reload if you are getting stale
                copy of the component due to caching.
              </TroubleshootListitem>
              <TroubleshootListitem>
                Use Incognito or Private browsing mode. If you are still getting stale copy of the
                component, try opening a <code>incognito tab</code>
              </TroubleshootListitem>
            </Typography>
          </ul>
        </AccDetailsContainer>
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
            <Typography component="div">
              <TroubleshootListitem>
                Verify MeshSync data is being received. Run <code>kubectl get svc -n meshery</code>.
                Docker Desktop: VPNkit commonly fails to assign an IP address to Meshery Broker
                (MeshSync). Verify that the Meshery Broker service has external IP address assigned.
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
            <Typography component="div">
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
    </Modal>
  );
};

export default TroubleshootingModal;
