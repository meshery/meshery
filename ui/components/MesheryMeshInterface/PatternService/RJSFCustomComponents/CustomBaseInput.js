import React from "react";
import { IconButton, InputAdornment, TextField, useTheme } from "@material-ui/core";
import HelpOutlineIcon from "../../../../assets/icons/HelpOutlineIcon";
import { CustomTextTooltip } from "../CustomTextTooltip";
import ErrorOutlineIcon from "../../../../assets/icons/ErrorOutlineIcon";
import { ERROR_COLOR } from "../../../../constants/colors";
import { iconSmall } from "../../../../css/icons.styles";
import { getHyperLinkDiv } from "../helper";

const BaseInput = (props) => {
  const additional = props.schema?.__additional_property; // check if the field is additional
  const name = (additional ? "Value" : props.label) // || props.id?.split('_')[-1].trim()
  const focused = props.options?.focused // true for datetime-local
  const prettifiedName = name || 'Enter a value'
  const style = {
    display : "flex",
    alignItems : "center",
  }

  const theme = useTheme();
  return (
    <>
      <div key={props.id} style={style}>
        <TextField
          variant={additional ? "standard" : "outlined"}
          size="small"
          focused={focused}
          type={props.options?.inputType}
          key={props.id}
          value={additional && props?.value === "New Value" ? "" : props?.value} // remove the default value i.e. New Value for additionalFields
          id={props.id}
          margin="dense"
          error={props.rawErrors?.length > 0}
          onChange={e => props?.onChange(e.target.value === "" ? props.options.emptyValue : e.target.value)}
          label={`${prettifiedName}`}
          InputProps={{
            style : { padding : "0px 0px 0px 0px", },
            endAdornment : (<InputAdornment position="start">
              {props.rawErrors?.length > 0 && (
                <CustomTextTooltip
                  backgroundColor={ERROR_COLOR}
                  flag={props?.formContext?.overrideFlag}
                  title={
                    <div>
                      {props.rawErrors?.map((error, index) => (
                        <div key={index}>{error}</div>
                      ))}
                    </div>

                  }
                  interactive={true}>
                  <IconButton component="span" size="small">
                    <ErrorOutlineIcon width="14px" height="14px" fill="#B32700" style={{ verticalAlign : "middle", ...iconSmall }}/>
                  </IconButton>
                </CustomTextTooltip>
              )}
              {props.schema?.description && (
                <CustomTextTooltip backgroundColor="#3C494F" flag={props?.formContext?.overrideFlag} title={getHyperLinkDiv(props.schema?.description)} interactive={true}>
                  <IconButton component="span" size="small">
                    <HelpOutlineIcon width="14px" height="14px"  fill={theme.palette.type === 'dark' ? "white" : "gray"}  style={{ verticalAlign : "middle", ...iconSmall }}/>
                  </IconButton>
                </CustomTextTooltip>
              )}
            </InputAdornment>),
          }} />
      </div>
    </>
  )
}

export default BaseInput;
