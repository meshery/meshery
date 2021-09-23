import React from "react";
import { Modal, Backdrop, Fade } from '@mui/core'
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";

const CustomModal = styled(Modal)(({ theme }) => ({
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

export const GenericModal = ({
    open, Content, handleClose, container
  }) =>  { const theme = useTheme();
    return (
      <CustomModal
        sx={{theme}}
        open={open}
        onClose={handleClose}
        closeAfterTransition
        BackdropComponent={Backdrop}
        BackdropProps={{ timeout : 200, }}
        container={container}
      >
        <Fade in={open} style={{ maxHeight : "90vh", overflow : "auto" }} >{Content}</Fade>
      </CustomModal>
    );
  }

  export default CustomModal;
  