import { TextField } from "@material-ui/core";
import React ,{ useEffect } from "react";
export const CustomUpDownField = (props) => {
  const name = props.label
  const [value, setValue] = React.useState(props?.value||props?.schema?.default||"\n");
  console.log("schema",props)
  useEffect(() => {
    props?.onChange(value)
  },[value])
  return (
    <TextField
      id="standard-number"
      label={name}
      key={props.id}
      value={value}
      variant="standard"
      onChange={e => setValue(parseInt(e.target.value)||0)}
      type="number"
      margin="none"
      error={props.rawErrors?.length > 0}
      helperText={props.rawErrors?.length > 0 ? props.rawErrors[0] : ""}
      size="large"
    />
  )
}