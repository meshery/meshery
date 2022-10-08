import { FormHelperText, TextField } from "@material-ui/core";
import React from "react";
import ErrorIcon from "@material-ui/icons/Error";
const style = {
  display : "flex",
  alignItems : "center",
}
const CustomHelperText = (props) => {
  return (
    <div style={style} id={props.id}>
      <ErrorIcon style={{ color : "#B32700", marginRight : "0.1rem", height : "1rem" }} />
      <FormHelperText error style={{ color : "#B32700" }}>{props.errormsg}</FormHelperText>
    </div>

  )
}

const CustomUpDownField = (props) => {
  const name = props.label
  return (
    <>
      <TextField
        id={props.id}
        size="small"
        label={name}
        key={props.id}
        value={props?.value}
        variant="outlined"
        onChange={e => e.target.value < 0 ? props?.onChange(e.target.value = "") : props?.onChange(e.target.value)}
        type="number"
        error={props.rawErrors?.length > 0}
        style={{ marginTop : '0.3em' }}
        InputProps={{
          style : { padding : "0px 0px 0px 0px", backgroundColor : "rgba(255, 255, 255, 0.4)" },
        }}
      />
      <div style={{ display : "flex" }}>
        {props.rawErrors?.map((errormsg, i) => (
          (errormsg === "required property" ? null
            : <CustomHelperText key={i} errormsg={errormsg} />
          )
        ))}
      </div>
    </>

  )
}

export { CustomUpDownField, CustomHelperText }