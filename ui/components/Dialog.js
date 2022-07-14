import {Box, Dialog, DialogTitle, DialogActions, Backdrop, Fade } from '@mui/material';
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";

const CustomDialogWrapper = styled(Dialog)(({ theme }) => ({
   display : "flex",
   alignItems : "center",
   justifyContent : "center", 
}))

const CustomBox = styled(Box)(({ theme }) => ({
   borderRadius : "5px", 
   textAlign : "center",
   paddingTop: theme.spacing(2),
   paddingButton: theme.spacing(2),
   paddingLeft: theme.spacing(3),
   paddingRight: theme.spacing(3),
   minWidth: theme.spacing(50),
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

const CustomDialog = ({
    open, Content, handleClose, container, Title, Buttons,
  }) =>  { const theme = useTheme();
    return (
      <CustomDialogWrapper
        open={open}
        onClose={handleClose}
        closeAfterTransition
        BackdropComponent={Backdrop}
        BackdropProps={{ timeout : 200, }}
        container={container}
      > 
      <DialogTitle>{Title}</DialogTitle>
      <CustomBox>
        <Fade in={open} sx={{ maxHeight : "90vh", overflow : "auto" }} >  
            {Content}
         </Fade>
      <DialogActions>{Buttons}</DialogActions>   
       </CustomBox>     
      </CustomDialogWrapper>
    );
  }

  export default CustomDialog;
  