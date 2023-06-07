import React from "react";
import PropTypes from  "prop-types";
import { Paper as MuiPaper, Box } from "@mui/material";

export const Paper = ({ elevation, ...rest }) => {
  return (
    <Box
      sx={{
        display: 'flex',
        flexWrap: 'wrap',
        '& > :not(style)': {
          m: 1,
          width: 128,
          height: 128,
        },
      }}
    >
      <MuiPaper elevation={elevation} {...rest} />
    </Box>
  )
}

Paper.propTypes = {
  label : PropTypes.number,
  variant : PropTypes.oneOf(['elevation', 'outlined'])
}

Paper.defaulProps = {
  label : 1,
  variant : "outlined",
}