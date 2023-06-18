import React from "react";
import { IconButton, InputAdornment, TextField, useTheme } from "@material-ui/core";
import HelpOutlineIcon from "../../../../assets/icons/HelpOutlineIcon";
import { CustomTextTooltip } from "../CustomTextTooltip";
import ErrorOutlineIcon from "../../../../assets/icons/ErrorOutlineIcon";
import { ERROR_COLOR } from "../../../../constants/colors";
import { iconSmall } from "../../../../css/icons.styles";
import { getHyperLinkDiv } from "../helper";
import { makeStyles } from '@material-ui/styles';

const BaseInput = (props) => {
  const additional = props.schema?.__additional_property; // check if the field is additional
  const name = (additional ? "Value" : props.label) // || props.id?.split('_')[-1].trim()
  const focused = props.options?.focused // true for datetime-local
  const prettifiedName = name || 'Enter a value'
  const [isFocused, setIsFocused] = React.useState(false);
  const style = {
    display : "flex",
    alignItems : "center"
  }
  const styles = makeStyles((theme) => ({
    customInputLabel : {
      color : theme.palette.secondary.text,
      backgroundColor : theme.palette.background.default,
      padding : "0.2rem",
      height : "1rem",
      borderRadius : "3px",
    },
  }));

  const theme = useTheme();
  console.log(theme)
  const classes = styles();
  return (
    console.log(prettifiedName),
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
          InputLabelProps={{
            className : (prettifiedName === "name" || prettifiedName === "namespace" || isFocused) ? classes.customInputLabel : '',
          }}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
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
                    <ErrorOutlineIcon width="14px" height="14px" fill="#B32700" style={{ verticalAlign : "middle", ...iconSmall }} />
                  </IconButton>
                </CustomTextTooltip>
              )}
              {props.schema?.description && (
                <CustomTextTooltip backgroundColor="#3C494F" flag={props?.formContext?.overrideFlag} title={getHyperLinkDiv(props.schema?.description)} interactive={true}>
                  <IconButton component="span" size="small">
                    <HelpOutlineIcon width="14px" height="14px" fill={theme.palette.type === 'dark' ? "white" : "gray"} style={{ verticalAlign : "middle", ...iconSmall }} />
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
