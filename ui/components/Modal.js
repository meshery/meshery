import { Modal, Backdrop, Fade } from "@mui/core";
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";

const CustomModalWrapper = styled(Modal)(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
}));

/**
 *
 * @param {{
 *  open?: boolean,
 *  Content?: JSX.Element,
 *  handleClose?: (event: {}, reason: "backdropClick" | "escapeKeyDown") => void,
 *  container?: React.ReactInstance | (() => React.ReactInstance)
 * }} props
 * @returns
 */

const CustomModal = ({ open, Content, handleClose, container }) => {
  const theme = useTheme();
  return (
    <CustomModalWrapper
      open={open}
      onClose={handleClose}
      closeAfterTransition
      BackdropComponent={Backdrop}
      BackdropProps={{ timeout: 200 }}
      container={container}     
    >
      <Fade in={open} sx={{ maxHeight: "90vh", overflow: "auto" }}>
        {Content}
      </Fade>
    </CustomModalWrapper>
  );
};

export default CustomModal;
