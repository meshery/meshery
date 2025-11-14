import React, { useState } from 'react';
import PropTypes from 'prop-types';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
  Paper,
  ClickAwayListener,
  Fade,
  Popper,
  styled,
} from '@sistent/sistent';
import {
  TwitterShareButton,
  LinkedinShareButton,
  FacebookShareButton,
  TwitterIcon,
  LinkedinIcon,
  FacebookIcon,
} from 'react-share';
import ReplyIcon from '@mui/icons-material/Reply';
import {
  fortioResultToJsChartData,
  makeChart,
  makeOverlayChart,
  makeMultiChart,
} from '../lib/chartjs-formatter';

const StyledDialogTitle = styled(DialogTitle)({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
});

const TitleContainer = styled('div')({
  display: 'flex',
  alignItems: 'center',
  flex: 1,
});

const ShareIconButton = styled(IconButton)({
  transform: 'scaleX(-1)',
});

const SocialPopper = styled(Popper)(({ theme }) => ({
  maxWidth: theme.spacing(30),
  zIndex: theme.zIndex.modal + 1,
}));

const SocialPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(1),
}));

const SocialIconWrapper = styled('span')(({ theme }) => ({
  margin: theme.spacing(0.4),
}));

function MesheryChartDialog(props) {
  const { open, title, handleClose, content, data, rawdata } = props;
  const [socialExpand, setSocialExpand] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);
  const [socialMessage, setSocialMessage] = useState('');

  const getSocialMessageForPerformanceTest = (rps, percentile) => {
    return `I achieved ${rps.trim()} RPS running my service at a P99.9 of ${percentile} ms using @mesheryio with @smp_spec! Find out how fast your service is with`;
  };

  const handleSocialExpandClick = (e) => {
    if (!data || !data.length) return;

    // Compute chart data similar to MesheryChart component
    let chartData;
    if (data.length === 1) {
      chartData = makeChart(fortioResultToJsChartData(rawdata || [], data[0]));
    } else if (data.length === 2) {
      chartData = makeOverlayChart(
        fortioResultToJsChartData(rawdata || [], data[0]),
        fortioResultToJsChartData(rawdata || [], data[1]),
      );
    } else if (data.length > 2) {
      chartData = makeMultiChart(rawdata || [], data);
    }

    if (chartData && chartData.options && chartData.options.metadata) {
      const rps = chartData.options.metadata.qps?.display?.value?.split(' ')[1] || '';
      const percentile = chartData.percentiles?.[4]?.Value || '';
      if (rps && percentile) {
        setSocialMessage(getSocialMessageForPerformanceTest(rps, percentile));
      }
    }

    setAnchorEl(e.currentTarget);
    e.stopPropagation();
    setSocialExpand((prevState) => !prevState);
  };

  const dialogTitle = title && title.length ? title : 'Comparison';

  return (
    <React.Fragment>
      <Dialog
        fullWidth
        maxWidth="md"
        open={open}
        onClose={handleClose}
        aria-labelledby="chart-dialog-title"
      >
        <StyledDialogTitle id="chart-dialog-title" data-testid="chart-dialog-title">
          <TitleContainer>{dialogTitle}</TitleContainer>
          {data && data.length > 0 && (
            <ShareIconButton
              aria-label="Share"
              onClick={handleSocialExpandClick}
              data-testid="share-button"
            >
              <ReplyIcon />
            </ShareIconButton>
          )}
        </StyledDialogTitle>
        <SocialPopper open={socialExpand} anchorEl={anchorEl} transition placement="bottom-end">
          {({ TransitionProps }) => (
            <ClickAwayListener onClickAway={() => setSocialExpand(false)}>
              <Fade {...TransitionProps} timeout={350}>
                <SocialPaper>
                  <SocialIconWrapper>
                    <TwitterShareButton
                      url={'https://meshery.io'}
                      title={socialMessage}
                      hashtags={['opensource']}
                    >
                      <TwitterIcon size={32} />
                    </TwitterShareButton>
                  </SocialIconWrapper>
                  <SocialIconWrapper>
                    <LinkedinShareButton url={'https://meshery.io'} summary={socialMessage}>
                      <LinkedinIcon size={32} />
                    </LinkedinShareButton>
                  </SocialIconWrapper>
                  <SocialIconWrapper>
                    <FacebookShareButton
                      url={'https://meshery.io'}
                      quote={socialMessage}
                      hashtag={'#opensource'}
                    >
                      <FacebookIcon size={32} />
                    </FacebookShareButton>
                  </SocialIconWrapper>
                </SocialPaper>
              </Fade>
            </ClickAwayListener>
          )}
        </SocialPopper>
        <DialogContent>
          <DialogContentText>{content}</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} color="primary">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </React.Fragment>
  );
}

MesheryChartDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  handleClose: PropTypes.func.isRequired,
  content: PropTypes.node.isRequired,
  title: PropTypes.string,
  data: PropTypes.array,
  rawdata: PropTypes.array,
};

export default MesheryChartDialog;
