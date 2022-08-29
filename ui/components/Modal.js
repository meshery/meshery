import { Modal, Box, Backdrop, Fade } from '@mui/material'
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";

const CustomModalWrapper = styled(Modal)(({ theme }) => ({
   display : "flex",
   alignItems : "center",
   justifyContent : "center", 
}))

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

const CustomModal = ({
    open, Content, handleClose, container
  }) =>  { const theme = useTheme();
    return (
      <CustomModalWrapper
        open={open}
        onClose={handleClose}
        closeAfterTransition
        BackdropProps={{ timeout : 200, }}
        container={container}
      >
        {Content}
      </CustomModalWrapper>
    );
  }

  export default CustomModal;
  