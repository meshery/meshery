import React, { useState } from "react";
import IconButton from "@mui/material/IconButton";
import Checkbox from "@mui/material/Checkbox";
import FormControlLabel from "@mui/material/FormControlLabel";
import Tooltip from "@mui/material/Tooltip";
import Popper from "@mui/material/Popper";
import ClickAwayListener from "@mui/material/ClickAwayListener";
import { Paper } from "@mui/material";
import { withStyles } from "@material-ui/core/styles";
import ColumnIcon from "../assets/icons/coulmn";
import Slide from "@mui/material/Slide";
import Box from "@mui/material/Box";


const styles = theme => ({
  epaper : {
    background : theme.palette.secondary.headerColor,
    padding : "1rem"
  } });

const CustomColumnVisibilityControl = ({ columns, customToolsProps, classes }) => {

  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);

  const handleOpen = event => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleColumnVisibilityChange = (columnName, isVisible) => {
    customToolsProps.setColumnVisibility(prevState => ({
      ...prevState,
      [columnName] : isVisible
    }));
  };



  return (
    <>
      <Tooltip title="View Columns">
        <IconButton
          onClick={handleOpen}
          sx={{
            "&:hover" : {
              borderRadius : "4px"
            },
          }}
        >
          < ColumnIcon  />
        </IconButton>
      </Tooltip>

      <Box>
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
            <Slide {...TransitionProps} direction="down" timeout={350} in={open} mountOnEnter unmountOnExit>
              <Box>
                <ClickAwayListener onClickAway={handleClose}>
                  <Paper  className={classes.epaper}>
                    <div style={{ display : "flex", flexDirection : "column" }}>
                      {columns.map(col => (
                        <FormControlLabel
                          key={col.name}
                          control={
                            <Checkbox
                              checked={customToolsProps.columnVisibility[col.name]}
                              onChange={e =>
                                handleColumnVisibilityChange(col.name, e.target.checked)
                              }
                            />
                          }
                          label={col.label}
                        />
                      ))}
                    </div>
                  </Paper>
                </ClickAwayListener>
              </Box>
            </Slide>
          )}
        </Popper>
      </Box>
    </>
  );
};

export default withStyles(styles) (CustomColumnVisibilityControl);
