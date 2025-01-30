import React from 'react';
import {
  Typography,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Slide,
  CloseIcon,
  styled,
} from '@layer5/sistent';

const DialogTitleStyled = styled(DialogTitle)(({ theme }) => ({
  margin: 0,
  padding: theme.spacing(2),
  textAlign: 'center',
  backgroundColor: '#252E31',
  color: theme.palette.background.constant.white,
}));

const CloseButton = styled(IconButton)(({ theme }) => ({
  position: 'absolute',
  color: theme.palette.background.constant.white,
  right: theme.spacing(1),
  top: theme.spacing(1),
  transform: 'rotate(-90deg)',
  '&:hover': {
    transform: 'rotate(90deg)',
    transition: 'all .3s ease-in',
  },
}));

const ImgWrapper = styled('div')({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  padding: '2rem',
});

const InsideImgWrapper = styled('img')(({ theme }) => ({
  padding: '0rem 0.5rem',
  content:
    theme.palette.mode === 'dark'
      ? "url('/static/img/meshery-logo-text.svg')"
      : "url('/static/img/meshery-logo-light-text.svg')",
}));

const InsideImgWrapperLogo = styled('img')({
  padding: '0rem 0.5rem',
});

const InstallButton = styled(Button)({
  marginBottom: '1rem',
});

const Transition = React.forwardRef(function Transition(props, ref) {
  return <Slide direction="up" ref={ref} {...props} />;
});

export default function PlaygroundMeshDeploy(props) {
  const handlePage = (e) => {
    window.open('https://meshery.io/#getting-started', '_blank');
    e.stopPropagation();
  };

  return (
    <div>
      <Dialog
        aria-labelledby="customized-dialog-title"
        open={props.isOpen}
        onClose={props.closeForm}
        TransitionComponent={Transition}
      >
        <DialogTitleStyled id="customized-dialog-title">
          <Typography variant="h6">The Cloud Native Playground</Typography>
          <CloseButton aria-label="close" onClick={props.closeForm}>
            <CloseIcon />
          </CloseButton>
        </DialogTitleStyled>
        <DialogContent>
          <ImgWrapper>
            <InsideImgWrapperLogo width="20%" height="20%" src="/static/img/meshery-logo.png" />
            <InsideImgWrapper width="50%" height="50%" />
          </ImgWrapper>
          <Typography gutterBottom>
            Meshery Playground gives you hands-on experience with designing cloud native systems -
            from your browser - using every CNCF project. Choose a{' '}
            <a href="https://layer5.io/learn/learning-paths" style={{ color: '#00b39f' }}>
              Learning Path
            </a>{' '}
            and work through labs as you visually and collaboratively learn-by-doing without having
            to install a single thing.
          </Typography>
          <Typography gutterBottom>
            To ensure that Meshery Playground remains a clean sandbox for all to use, many of
            Meshery&apos;s features are disabled. For full access to all of Meshery&apos;s features,
            deploy your own instance of Meshery.
          </Typography>
        </DialogContent>
        <DialogActions sx={{ display: 'flex', justifyContent: 'center' }}>
          <InstallButton size="large" variant="contained" color="primary" onClick={handlePage}>
            Install Meshery
          </InstallButton>
        </DialogActions>
      </Dialog>
    </div>
  );
}
