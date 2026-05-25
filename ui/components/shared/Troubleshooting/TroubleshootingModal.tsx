/**
 * TroubleshootingModal — Meshery Extensions troubleshooting guide dialog.
 *
 * Phase 5.b.2 migration target: composes the shared `Modal` primitive so the
 * dialog inherits the same header / chrome that the rest of the app uses,
 * dropping the bespoke `StyledPaper` overlay the legacy component painted on
 * top of Sistent's raw `Modal`.
 *
 * The accordion content (Stale Data, Missing Data, Additional Resources) and
 * the "Need help?" footer copy are preserved verbatim so the diagnostic
 * material the dialog surfaces is unchanged.
 */
import * as React from 'react';
import {
  Typography,
  Accordion,
  AccordionDetails,
  AccordionSummary,
  InfoIcon,
  LIGHT_TEAL,
} from '@sistent/sistent';
import { styled } from '@/theme';
import { ExpandMore } from '@/assets/icons';
import { Modal } from '@/components/shared/Modal';

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

const AccordionContainer = styled(Accordion)(({ theme }) => ({
  margin: '0 !important',
  borderBottom: `1px solid ${theme.palette.divider}`,
}));

const AccordionSummaryStyled = styled(AccordionSummary)(({ theme }) => ({
  '&.Mui-expanded': {
    backgroundColor: theme.palette.background.card,
  },
}));

const AccDetailHead = styled(Typography)(({ theme }) => ({
  color: theme.palette.mode === 'dark' ? theme.palette.grey[300] : theme.palette.grey[700],
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
  boxShadow: theme.shadows[3],
}));

const AccDetailsContainer = styled(AccordionDetails)({
  flexDirection: 'column',
});

const TroubleshootHelpLink = styled('a')(({ theme }) => ({
  color:
    theme.palette.mode === 'dark'
      ? theme.palette.background.brand?.default
      : theme.palette.info.main,
  fontWeight: 'bold',
  textDecoration: 'none',
}));

const ContactHelpLink = styled('a')(({ theme }) => ({
  color: theme.palette.common.white,
  fontWeight: 'bold',
  textDecoration: 'none',
}));

const FooterText = styled(Typography)(({ theme }) => ({
  color: theme.palette.common.white,
  fontSize: '.85rem',
}));

const Details = styled('div')({
  display: 'block',
});

export interface TroubleshootingModalProps {
  /** Whether the modal is currently visible. */
  open: boolean;
  /** Controlled setter that the modal calls with `false` to close. */
  setOpen: (open: boolean) => void;
  /** Optional header-level error message; used to auto-expand the missing-data panel. */
  viewHeaderErrorMessage?: string;
  /** Optional inline error detail (reserved; kept for prop-shape parity). */
  viewDataErrorMessage?: string;
}

const TroubleshootingModal: React.FC<TroubleshootingModalProps> = (props) => {
  const [expanded, setExpanded] = React.useState<string | false>(false);
  const [missingData, setMissingData] = React.useState(
    !!props.viewHeaderErrorMessage?.includes('data missing'),
  );

  const handleChange = (panel: string) => (_event: unknown, isExpanded: boolean) => {
    setExpanded(isExpanded ? panel : false);
    if (missingData) {
      setMissingData(false);
    }
  };

  const handleClose = () => props.setOpen(false);

  return (
    <Modal
      isOpen={props.open}
      onClose={handleClose}
      title="Extensions Troubleshooting Guide"
      size="md"
      aria-labelledby="modal-modal-title"
      aria-describedby="modal-modal-description"
      actions={
        <FooterText variant="body2">
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
          expandIcon={<ExpandMore />}
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
                <KeyStyleContainer>E</KeyStyleContainer> to force reload if you are getting stale
                copy of the component due to caching.
              </TroubleshootListitem>
              <TroubleshootListitem>
                Use Incognito or Private browsing mode. If you are still getting stale copy of the
                component, try opening a <code>incognito tab</code>
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
          expandIcon={<ExpandMore />}
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
          expandIcon={<ExpandMore />}
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
    </Modal>
  );
};

export default TroubleshootingModal;
