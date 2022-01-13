import React from "react";
import { IconButton, InputAdornment, TextField } from "@material-ui/core";
import HelpOutlineIcon from "../HelpOutlineIcon";
import EnlargedTextTooltip from "../EnlargedTextTooltip";
import { formatString } from "../../helpers"

const CustomInputField = (props) => {
  const name = props?.name || props?.idSchema['$id']?.split('_')[1]
  const prettifiedName = formatString(name) || 'Input'
  const style= {
    display : "flex",
    alignItems : "center",
    justifyContent : "space-between"
  }
  return (
    <div key={props.id} style={style}>
      <TextField
        InputLabelProps={{ style : { pointerEvents : "auto", padding : "3px" } }}
        variant="outlined"
        size="small" autoFocus
        key={props.id}
        value={props.value}
        id={props.id}
        onChange={e => props?.onChange(e.target.value)}
        label={`${prettifiedName}`}
        InputProps={{ style : { padding : "4px 0px 5px 17px" },
          endAdornment : (<InputAdornment position="start">
            {props.schema?.description && (
              <EnlargedTextTooltip title={props.schema?.description}>
                <IconButton component="span" size="small">
                  <HelpOutlineIcon />
                </IconButton>
              </EnlargedTextTooltip>
            )}
          </InputAdornment>),  }}/>
    </div>
  )
}

const MemoizedCustomInputField = React.memo(CustomInputField)

export default MemoizedCustomInputField;
