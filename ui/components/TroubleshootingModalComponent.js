import * as React from 'react';
import Typography from '@material-ui/core/Typography';
import { Modal } from '@material-ui/core';
import InfoIcon from '@material-ui/icons/Info';
import Accordion from '@material-ui/core/Accordion';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import CloseIcon from '@material-ui/icons/Close';
import { makeStyles, Paper } from '@material-ui/core';
import { IconButton } from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  '@keyframes rotateCloseIcon': {
    from: {
      transform: 'rotate(0deg)',
    },
    to: {
      transform: 'rotate(360deg)',
    },
  },
  paper: {
    top: '50%',
    left: '50%',
    width: '47.5%',
    position: 'absolute',
    borderRadius: '10px',
    overflow: 'hidden',
    transform: 'translate(-50%,-50%)',
    backgroundColor: 'transparent',
    outline: 'none',
    [theme.breakpoints.down(1350)]: {
      width: '70%',
    },
  },
  headerContainer: {
    display: 'flex',
    alignItems: 'center',
  },
  headerText: {
    flexShrink: 0,
    margin: 10,
    fontSize: '1rem',
  },
  infoIcon: {
    color: theme.palette.type === 'dark' ? '#00B39F' : '#607d8b',
  },
  footerText: {
    color: '#EDEDED',
    fontSize: '.85rem',
    textDecoration: 'italics',
    fontFamily: 'Qanelas Soft, sans-serif',
  },
  modalHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingBottom: 10,
    padding: '0 .5rem',
    paddingTop: 10,
    backgroundColor: theme.palette.secondary.mainBackground,
  },
  modalFooter: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    paddingBottom: 10,
    padding: '0 .5rem',
    paddingTop: 10,
    backgroundColor: theme.palette.secondary.mainBackground,
  },
  modalContent: {
    backgroundColor: 'transparent',
  },
  details: {
    display: 'block',
  },
  modelHeader: {
    fontSize: '1rem',
    color: '#fff',
  },
  iconStyle: {
    color: '#fff',
  },
  iconContainer: {
    transition: 'all .3s',
    '&:hover': {
      backgroundColor: 'transparent !important',
      animation: '$rotateCloseIcon 1s',
    },
  },
  accordionContainer: {
    margin: '0 !important',
    borderBottom: '1px solid #ccc',
  },
  accordionSummary: {
    '&.Mui-expanded': {
      backgroundColor: theme.palette.type === 'dark' ? '#303030' : '#f1f1f1',
    },
  },
  accDetailHead: {
    color: theme.palette.type == 'dark' ? 'f1f1f1' : '#444',
    fontFamily: 'Qanelas Soft, sans-serif',
  },
  troubleshootListitem: {
    fontSize: '0.9rem',
    marginBottom: '1rem',
  },
  keyStyleContainer: {
    display: 'inline-block',
    padding: '0.1rem 0.5rem',
    background: theme.palette.secondary.elevatedComponent,
    // boxShadow: "0 1px 5px rgba(0,0,0,0.1)",
    margin: '0.3rem',
    borderRadius: '5px',
    boxShadow:
      'rgba(0, 0, 0, 0.17) 0px -0px 0px -5px inset, rgba(0, 0, 0, 0.15) 0px 0px 0px -3px inset, rgba(0, 0, 0, 0.1) 0px 4px 30px 0px inset, rgba(0, 0, 0, 0.06) 0px 2px 1px, rgba(0, 0, 0, 0.09) 0px 0px 2px, rgba(0, 0, 0, 0.09) 0px 8px 4px, rgba(0, 0, 0, 0.1) 0px 1px 0px, rgba(0, 0, 0, 0.1) 0px -2px 0px',
  },
  accDetailsContainer: {
    flexDirection: 'column',
  },
  troubleshootHelpLink: {
    color: theme.palette.type == 'dark' ? '#00B39F' : 'rgb(57, 102, 121)',
    fontWeight: 'bold',
    textDecoration: 'none',
  },
  contactHelpLink: {
    color: '#EDEDED',
    fontWeight: 'bold',
    textDecoration: 'none',
  },
  threadPara: {
    fontSize: '1.1rem',
    textAlign: 'center',
    marginTop: '0',
    fontStyle: 'italic',
  },
  outbtn: {
    margin: '0px 8px',
  },
}));

const TroubleshootingModal = (props) => {
  const classes = useStyles();
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
    <Modal
      open={props?.open}
      onClose={handleClose}
      aria-labelledby="modal-modal-title"
      aria-describedby="modal-modal-description"
    >
      <Paper elevation={6} className={classes.paper} square="true">
        <div className={classes.modalContent}>
          <div className={classes.modalHeader}>
            <Typography variant="h5"></Typography>
            <Typography className={classes.modelHeader} variant="h5">
              Meshmap Troubleshooting Guide
            </Typography>
            <IconButton
              data-cy="modal-close-btn"
              className={classes.iconContainer}
              onClick={handleClose}
            >
              <CloseIcon className={classes.iconStyle} />
            </IconButton>
          </div>
          <Accordion
            className={classes.accordionContainer}
            expanded={expanded === 'panel1'}
            onChange={handleChange('panel1')}
          >
            <AccordionSummary
              expandIcon={<ExpandMoreIcon />}
              aria-controls="panel1bh-content"
              id="panel1bh-header"
              className={classes.accordionSummary}
            >
              <div className={classes.headerContainer}>
                <InfoIcon className={classes.infoIcon} />
                <Typography variant="h6" className={classes.headerText}>
                  Stale Data
                </Typography>
              </div>
            </AccordionSummary>
            <AccordionDetails className={classes.details}>
              <Typography className={classes.accDetailHead}>
                <strong>Browser</strong>
              </Typography>
              <ul>
                <Typography>
                  <li className={classes.troubleshootListitem}>
                    Run <div className={classes.keyStyleContainer}>CTRL</div> +
                    <div className={classes.keyStyleContainer}>SHIFT</div>+
                    <div className={classes.keyStyleContainer}>R</div> or
                    <div className={classes.keyStyleContainer}>CMD</div>+
                    <div className={classes.keyStyleContainer}>OPTION</div>+
                    <div className={classes.keyStyleContainer}>E</div> to force reload if you are
                    getting stale copy of the component due to caching.
                  </li>
                  <li className={classes.troubleshootListitem}>
                    Use Incognito or Private browsing mode. If you are still getting stale copy of
                    the component, try opening a <code>incognito tab</code>
                  </li>
                </Typography>
              </ul>
            </AccordionDetails>
          </Accordion>
          <Accordion
            className={classes.accordionContainer}
            expanded={expanded === 'panel2' || missingData}
            onChange={handleChange('panel2')}
          >
            <AccordionSummary
              expandIcon={<ExpandMoreIcon />}
              aria-controls="panel2bh-content"
              id="panel2bh-header"
              className={classes.accordionSummary}
            >
              <div className={classes.headerContainer}>
                <InfoIcon className={classes.infoIcon} />
                <Typography variant="h6" className={classes.headerText}>
                  Missing Data
                </Typography>
              </div>
            </AccordionSummary>
            <AccordionDetails className={classes.accDetailsContainer}>
              <Typography className={classes.accDetailHead}>
                <strong>Meshery Database</strong>
              </Typography>
              <ul>
                <Typography>
                  <li className={classes.troubleshootListitem}>
                    Verify MeshSync data is being received. Run{' '}
                    <code>kubectl get svc -n meshery</code>. Docker Desktop: VPNkit commonly fails
                    to assign an IP address to Meshery Broker (MeshSync). Verify that the Meshery
                    Broker service has external IP address assigned.
                  </li>
                  <li className={classes.troubleshootListitem}>
                    Confirm that your machine&apos;s firewall isn&apos;t getting in the way.
                  </li>
                  <li className={classes.troubleshootListitem}>
                    Dump Meshery Database. Run <code>rm -rf ~/.meshery/config</code>.
                  </li>
                </Typography>
              </ul>
            </AccordionDetails>
          </Accordion>
          <Accordion
            className={classes.accordionContainer}
            expanded={expanded === 'panel3'}
            onChange={handleChange('panel3')}
          >
            <AccordionSummary
              expandIcon={<ExpandMoreIcon />}
              aria-controls="panel3bh-content"
              id="panel3bh-header"
              className={classes.accordionSummary}
            >
              <div className={classes.headerContainer}>
                <InfoIcon className={classes.infoIcon} />
                <Typography variant="h6" className={classes.headerText}>
                  Incompatibility
                </Typography>
              </div>
            </AccordionSummary>
            <AccordionDetails className={classes.accDetailsContainer}>
              <Typography className={classes.accDetailHead}>
                <strong>GraphQL/Golang Plugin</strong>
              </Typography>
              <ul>
                <Typography>
                  <li className={classes.troubleshootListitem}>
                    <b>Building</b>: Ensure that <code>go.mod</code> in the{' '}
                    <code>meshery/meshery</code> repository is identical to the <code>go.mod</code>{' '}
                    in the <code>layer5labs/meshery-extensions</code> repository are identical.
                  </li>
                  <li className={classes.troubleshootListitem}>
                    <b>Loading</b>: Confirm that the <code>plugin</code> version offered by Meshery
                    Cloud (at{' '}
                    <a
                      className={classes.troubleshootHelpLink}
                      href="https://meshery.layer5.io/capabilities"
                    >
                      https://meshery.layer5.io/capabilities
                    </a>
                    ), from which MeshMap files are retreived, matches the <code>plugin</code>{' '}
                    version that Meshery Server is using as the filesystem reference in{' '}
                    <code>~/.meshery/provider/Meshery/vx.x.x/</code>.
                  </li>
                </Typography>
              </ul>
            </AccordionDetails>
          </Accordion>
          <Accordion
            className={classes.accordionContainer}
            expanded={expanded === 'panel4'}
            onChange={handleChange('panel4')}
          >
            <AccordionSummary
              expandIcon={<ExpandMoreIcon />}
              aria-controls="panel4bh-content"
              id="panel4bh-header"
              className={classes.accordionSummary}
            >
              <div className={classes.headerContainer}>
                <InfoIcon className={classes.infoIcon} />
                <Typography variant="h6" className={classes.headerText}>
                  Additional Resources
                </Typography>
              </div>
            </AccordionSummary>
            <AccordionDetails className={classes.accDetailsContainer}>
              <Typography variant="h6" className={classes.headerText}>
                <strong>Troubleshooting Tips</strong>
              </Typography>
              <ul>
                <Typography>
                  <li className={classes.troubleshootListitem}>
                    {' '}
                    <a
                      className={classes.troubleshootHelpLink}
                      href="http://discuss.meshery.io/t/what-are-some-troubleshooting-tips-for-meshmap"
                      target="_blank"
                      rel="noreferrer"
                    >
                      &quot;What are some troubleshooting tips for MeshMap?&quot;
                    </a>
                  </li>
                </Typography>
              </ul>
            </AccordionDetails>
          </Accordion>
          <div className={classes.modalFooter}>
            <Typography className={classes.footerText} variant="h6">
              Need help? Contact us via{' '}
              <a
                className={classes.contactHelpLink}
                href="mailto:meshmap@layer5.io"
                target="_blank"
                rel="noreferrer"
              >
                email
              </a>{' '}
              or{' '}
              <a
                className={classes.contactHelpLink}
                href="http://discuss.meshery.io"
                target="_blank"
                rel="noreferrer"
              >
                community forum
              </a>
              .
            </Typography>
          </div>
        </div>
      </Paper>
    </Modal>
  );
};

export default TroubleshootingModal;
