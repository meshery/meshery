import React, { useState } from 'react'
import PropTypes from 'prop-types'
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
} from "@sistent/sistent"
import {
  TwitterShareButton,
  LinkedinShareButton,
  FacebookShareButton,
  TwitterIcon,
  LinkedinIcon,
  FacebookIcon,
} from "react-share"
import ReplyIcon from "@mui/icons-material/Reply"
const StyledDialogTitle = styled(DialogTitle)({
  display: "flex",
  justifyContent: "space-between",
  alignItems: "center",
})

const ShareIconButton = styled(IconButton)({
  transform: "scaleX(-1)",
})

const SocialPopper = styled(Popper)({
  width: 500,
})

const SocialPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(1),
}))

const SocialIconWrapper = styled("span")(({ theme }) => ({
  margin: theme.spacing(0.4),
}))

function MesheryChartDialog(props) {
  const { open, title, handleClose, content, socialMessage } = props
  const [socialExpand, setSocialExpand] = useState(false)
  const [anchorEl, setAnchorEl] = useState(null)

  const handleSocialExpandClick = (e) => {
    setAnchorEl(e.currentTarget)
    e.stopPropagation()
    setSocialExpand((prevState) => !prevState)
  }
  return (
    <React.Fragment>
      <Dialog 
      fullWidth maxWidth="md" 
      open={open} 
      onClose={handleClose} 
      aria-labelledby="chart-dialog-title">
        <StyledDialogTitle id="chart-dialog-title" data-testid="chart-dialog-title">
        {title && title.length ? title : 'Comparison'}
          <ShareIconButton aria-label="Share" onClick={handleSocialExpandClick}>
            <ReplyIcon />
          </ShareIconButton>
        </StyledDialogTitle>
        <SocialPopper open={socialExpand} anchorEl={anchorEl} transition style={{ zIndex: "1301" }}>
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
                    <LinkedinShareButton
                      url={'https://meshery.io'}
                      summary={socialMessage}
                    >
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
  )
}

MesheryChartDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  handleClose: PropTypes.func.isRequired,
  content: PropTypes.node.isRequired,
  title: PropTypes.string,
  socialMessage: PropTypes.string,
}

export default MesheryChartDialog
