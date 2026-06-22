import React, { useEffect, useState } from "react";
import {
  alpha,
  BLACK,
  Box,
  SAFFRON,
  styled,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  WarningIcon,
} from "@sistent/sistent";

const CARD_WIDTH_PX = 445;
const HEADER_HEIGHT_PX = 56;
const CONTENT_TOP_PX = 34;
const CONTENT_BOTTOM_PX = 30;
const CONTENT_SIDE_PX = 24;
const STRIP_HEIGHT_PX = 22;

const SessionExpiredContent = styled(Box)(() => ({
  width: "100%",
  overflowWrap: "break-word",
  textAlign: "center",
  padding: `${CONTENT_TOP_PX}px ${CONTENT_SIDE_PX}px ${CONTENT_BOTTOM_PX}px ${CONTENT_SIDE_PX}px`,
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
}));

const IconContainer = styled(Box)(() => ({
  width: "28px",
  height: "28px",
  display: "inline-flex",
  alignItems: "center",
  justifyContent: "center",
  position: "absolute",
  left: "16px",
  top: "50%",
  transform: "translateY(-50%)",
  color: "#000000",
}));

function AlertUnauthenticatedSession() {
  const [countDown, setCountDown] = useState(3);

  useEffect(() => {
    const timer = setTimeout(() => {
      if (countDown === 1) {
        // Propagate existing request parameters, if present.
        const existingQueryString = window.location.search;
        window.location = `/user/login${existingQueryString}`;
        return;
      }
      setCountDown((countDown) => countDown - 1);
    }, 1000);
    return () => clearTimeout(timer);
  }, [countDown]);

  return (
    <Dialog
      open
      disableEscapeKeyDown
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
      BackdropProps={{
        style: {
          backgroundColor: alpha(BLACK, 0.35),
        },
      }}
      PaperProps={{
        style: {
          width: "calc(100% - 32px)",
          maxWidth: `${CARD_WIDTH_PX}px`,
          borderBottom: `${STRIP_HEIGHT_PX}px solid ${SAFFRON}`,
          borderRadius: "6px",
          overflow: "hidden",
          boxShadow: `0 18px 42px ${alpha(BLACK, 0.35)}`,
        },
      }}
    >
      <DialogTitle
        id="alert-dialog-title"
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          position: "relative",
          textAlign: "center",
          padding: "0 16px",
          color: "#000000",
          backgroundColor: SAFFRON,
          fontSize: "18px",
          lineHeight: "26px",
          fontWeight: 500,
          minHeight: `${HEADER_HEIGHT_PX}px`,
        }}
      >
        <IconContainer>
          <WarningIcon color="#000000" width={26} height={26} />
        </IconContainer>
        Session Expired
      </DialogTitle>
      <DialogContent sx={{ padding: 0 }}>
        <SessionExpiredContent id="alert-dialog-description">
          <Typography sx={{ color: "text.default", fontSize: "15px", marginBottom: "12px" }}>
            User not authenticated
          </Typography>
          <Typography sx={{ color: "text.default", fontSize: "15px" }}>
            You will be redirected to Login page in {countDown}
          </Typography>
        </SessionExpiredContent>
      </DialogContent>
    </Dialog>
  );
}

export default AlertUnauthenticatedSession;
