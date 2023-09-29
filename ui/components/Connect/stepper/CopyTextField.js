import React, { useRef } from "react";
import {
  Button,
  TextField,
  Typography,
  withStyles,
} from "@material-ui/core";
import { useTheme } from "@material-ui/core/styles"

const CopyTextField = ({ content, token }) => {
  const textFieldRef = useRef(null);
  const theme = useTheme()

  const copyToClipboard = () => {
    textFieldRef.current.select();
    document.execCommand("copy");
    window.getSelection().removeAllRanges();
  };

  return (
    <>
      <Typography
        style={{
          fontSize: "0.75rem",
          fontWeight: "bold",
          textTransform: "uppercase",
        }}
      >
        {content.title}
      </Typography>
      <TextField
        inputRef={textFieldRef}
        value={token ? token : content.value}
        fullWidth
        style={{
          marginBottom: "1.5rem",
        }}
        InputProps={{
          style: {
            padding: "0.25rem 0",
            height: "2rem",
            borderRadius: "4px",
            fontSize: "12px",
          },
          readOnly: true,
          endAdornment: (
            <Button
              onClick={copyToClipboard}
              style={{ color: theme.palette.charcoal }}
            >
              Copy
            </Button>
          ),
        }}
      />
    </>
  );
};

export default withStyles(styles)(CopyTextField);
