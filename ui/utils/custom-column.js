import React, { useState, useRef } from "react";
import IconButton from "@mui/material/IconButton";
import Checkbox from "@mui/material/Checkbox";
import FormControlLabel from "@mui/material/FormControlLabel";
import Tooltip from "@mui/material/Tooltip";
import Popper from "@mui/material/Popper";
import ClickAwayListener from "@mui/material/ClickAwayListener";
import { withStyles } from "@material-ui/core/styles";
import ColumnIcon from "../assets/icons/coulmn";
import Slide from "@mui/material/Slide";
import Box from "@mui/material/Box";
import { Card } from "@material-ui/core";

const styles = (theme) => ({
  epaper : {
    "&: .MuiPaper-root" : {
      background : theme.palette.secondary.headerColor,
      color : theme.palette.secondary.textMain,
    },
    padding : "1rem",
    boxShadow :
      "0px 4px 0px -2px rgb(0 179 159 / 10%), 0px 4px 0px -2px rgb(0 179 159 / 10%), 0px 4px 0px -2px rgb(0 179 159 / 10%)",
  },
});

const CustomColumnVisibilityControl = ({ columns, customToolsProps, classes }) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);
  const containerRef = useRef(null);

  const handleOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleColumnVisibilityChange = (columnName, isVisible) => {
    customToolsProps.setColumnVisibility((prevState) => ({
      ...prevState,
      [columnName] : isVisible,
    }));
  };

  // const paperStyle = {
  //   background : theme.palette.secondary.link,
  //   color: theme.palette.secondary.textMain,
  //   fontWeight: "bold",
  // };

  return (
    <div>
      <Tooltip title="View Columns">
        <IconButton
          onClick={handleOpen}
          sx={{
            "&:hover" : {
              borderRadius : "4px",
            },
          }}
        >
          <ColumnIcon />
        </IconButton>
      </Tooltip>

      <Box sx={{ overflow : "hidden" }} ref={containerRef}>
        <Popper
          open={Boolean(anchorEl)}
          anchorEl={anchorEl}
          placement="bottom-start"
          anchorOrigin={{
            vertical : "bottom",
            horizontal : "center",
          }}
          transformOrigin={{
            vertical : "top",
            horizontal : "center",
          }}
          transition
        >
          {({ TransitionProps }) => (
            <Slide
              {...TransitionProps}
              direction="down"
              in={open}
              timeout={400}
              mountOnEnter
              unmountOnExit
              container={containerRef.current}
            >
              <Box>
                <ClickAwayListener onClickAway={handleClose}>
                  <Card className={classes.epaper} elevation={2}>
                    <div style={{ display : "flex", flexDirection : "column" }}>
                      {columns.map((col) => (
                        <FormControlLabel
                          key={col.name}
                          control={
                            <Checkbox
                              checked={customToolsProps.columnVisibility[col.name]}
                              onChange={(e) => handleColumnVisibilityChange(col.name, e.target.checked)}
                              sx={{
                                "&.Mui-checked" : {
                                  color : "#00B39F",
                                }
                              }}
                            />
                          }
                          label={col.label}
                        />
                      ))}
                    </div>

                  </Card>
                </ClickAwayListener>
              </Box>
            </Slide>
          )}
        </Popper>
      </Box>
    </div>
  );
};

export default withStyles(styles)(CustomColumnVisibilityControl);
