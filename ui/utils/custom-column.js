import React, { useState } from "react";
import IconButton from "@mui/material/IconButton";
import Checkbox from "@mui/material/Checkbox";
import FormControlLabel from "@mui/material/FormControlLabel";
import Tooltip from "@mui/material/Tooltip";
import Popper from "@mui/material/Popper";
import ClickAwayListener from "@mui/material/ClickAwayListener";
import { Paper } from "@mui/material";
import { withStyles } from "@mui/styles";
import ColumnIcon from "../assets/icons/coulmn";


const CustomColumnVisibilityControl = ({ columns, customToolsProps }) => {

  const [anchorEl, setAnchorEl] = useState(null);

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

      <Popper
        open={Boolean(anchorEl)}
        anchorEl={anchorEl}
        placement="bottom-end"
        modifiers={[
          {
            name : "preventOverflow",
            enabled : true,
            options : {
              padding : 16
            }
          }
        ]}
      >
        <ClickAwayListener onClickAway={handleClose}>
          <Paper sx={{ padding : "1rem" }} >
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
      </Popper>
    </>
  );
};

export default withStyles() (CustomColumnVisibilityControl);
