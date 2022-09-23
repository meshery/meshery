import React from "react";
import { IconButton, InputAdornment, TextField } from "@material-ui/core";
import HelpOutlineIcon from "../HelpOutlineIcon";
import EnlargedTextTooltip from "../EnlargedTextTooltip";
import { CustomHelperText } from "./CustomUpDownWidget";
const CustomInputField = (props) => {
  const name = props.label // || props.id?.split('_')[-1].trim()
  const prettifiedName = name || 'Enter a value'
  const style = {
    display : "flex",
    alignItems : "center",
    justifyContent : "space-between",
    marginRight : "4px"
  }

  return (
    <div key={props.id} style={style}>
      <TextField
        variant="outlined"
        size="small"
        key={props.id}
        value={props.value}
        id={props.id}
        error={props.rawErrors?.length > 0}
        onChange={e => props?.onChange(e.target.value)}
        label={`${prettifiedName}`}
        InputProps={{ style : { padding : "2px 0px 3px 0px", backgroundColor : "rgba(255, 255, 255, 0.4)" },
          endAdornment : (<InputAdornment position="start">
            {props.schema?.description && (
              <EnlargedTextTooltip title={props.schema?.description}>
                <IconButton component="span" size="small">
                  <HelpOutlineIcon />
                </IconButton>
              </EnlargedTextTooltip>
            )}
          </InputAdornment>),  }}/>
      <div style={{ display : "flex" }}>
        {props.rawErrors?.map((errormsg, i) => (
          (errormsg === "is a required property" ? null
            : <CustomHelperText key={i} errormsg={errormsg} id={props.id} />
          )
        ))}
      </div>
    </div>
  )
}

export default CustomInputField;
