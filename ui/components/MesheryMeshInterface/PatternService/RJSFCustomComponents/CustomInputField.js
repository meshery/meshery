import React from "react";
import { IconButton, InputAdornment, TextField } from "@material-ui/core";
import HelpOutlineIcon from "../HelpOutlineIcon";
import EnlargedTextTooltip from "../EnlargedTextTooltip";

const CustomInputField = (props) => {
  const name = props.label // || props.id?.split('_')[-1].trim()
  const prettifiedName = name || 'Enter a value'
  const style = {
    display : "flex",
    alignItems : "center",
  }

  return (
    <>
      <div key={props.id} style={style}>
        <TextField
          variant="outlined"
          size="small"
          key={props.id}
          value={props?.value}
          id={props.id}
          error={props.rawErrors?.length > 0}
          onChange={e => props?.onChange(e.target.value=== "" ? props.options.emptyValue : e.target.value)}
          label={`${prettifiedName}`}
          InputProps={{
            style : { padding : "0px 0px 0px 0px", backgroundColor : "rgba(255, 255, 255, 0.4)" },
            endAdornment : (<InputAdornment position="start">
              {props.schema?.description && (
                <EnlargedTextTooltip title={props.schema?.description}>
                  <IconButton component="span" size="small">
                    <HelpOutlineIcon />
                  </IconButton>
                </EnlargedTextTooltip>
              )}
            </InputAdornment>),
          }} />
      </div>
    </>
  )
}

export default CustomInputField;
