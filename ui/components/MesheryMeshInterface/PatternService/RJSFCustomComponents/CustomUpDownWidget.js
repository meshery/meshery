import { FormHelperText, TextField } from "@material-ui/core";
import React from "react";
import ErrorIcon from "@material-ui/icons/Error";
const style={
  display : "flex",
  alignItems : "center",
}
const CustomHelperText=(props) => {
  return (
    <div style={style} id={props.id}>
      <ErrorIcon style={{ color : "red",marginRight : "2px",height : "18px" }}/>
      <FormHelperText error>{props.errormsg}</FormHelperText>
    </div>

  )
}

const CustomUpDownField = (props) => {
  const name = props.label
  return (
    <>
      <TextField
        id="standard-number"
        label={name}
        key={props.id}
        value={props?.value}
        variant="standard"
        onChange={e => props?.onChange(e.target.value)}
        type="number"
        margin="none"
        error={props.rawErrors?.length > 0}
        size="large"
      />
      <div style={{ display : "flex" }}>
        {props.rawErrors?.map((errormsg, i) => (
          (errormsg==="is a required property"?null
            :<CustomHelperText key={i} errormsg={errormsg}/>
          )
        ))}
      </div>
    </>

  )
}

export { CustomUpDownField,CustomHelperText }