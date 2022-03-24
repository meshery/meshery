import React from "react";
import { IconButton, InputAdornment, TextField } from "@material-ui/core";
import HelpOutlineIcon from "../HelpOutlineIcon";
import EnlargedTextTooltip from "../EnlargedTextTooltip";
import { formatString } from "../../helpers"

// The special input fields
export const omitTitleFields = ["name", "namespace"]

const CustomInputField = (props) => {
  const name = props?.name || props?.idSchema['$id']?.split('_')[1].trim()
  const prettifiedName = formatString(name) || 'Input'
  const style = {
    display : "flex",
    alignItems : "center",
    justifyContent : "space-between"
  }

  return (
    <div key={props.id} style={style}>
      <TextField
        InputLabelProps={{ style : { pointerEvents : "auto", padding : "2px" } }}
        variant="outlined"
        size="small"
        key={props.id}
        value={props.value}
        id={props.id}
        onChange={e => props?.onChange(e.target.value)}
        label={`${prettifiedName}`}
        InputProps={{ style : { padding : "2px 0px 3px 0px" },
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
