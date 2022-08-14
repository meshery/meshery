import { Box,  IconButton,  makeStyles, TextField, Typography } from '@material-ui/core';
import { Add, Remove } from '@material-ui/icons';
import React, { useEffect } from "react"
const useStyles = makeStyles(() => ({
  innerRoot : {
    marginTop : "2px",
    display : "flex",
    border : "1px solid #ccc",
    width : "40%",
  },
  input : {
    alignItems : "center",
    fontSize : '16px',
  }


}))
const CustomUpDownWidget = (props) => {
  const name = props.label
  const classes = useStyles();
  const [value, setValue] = React.useState(props.value||props.schema.default||0);
  console.log("schema",props)
  useEffect(() => {
    props?.onChange(value)
  },[value])
  return (
    <Box style={{ marginTop : "4px" }}>
      <Typography>{name}</Typography>
      <Box className={classes.innerRoot}>
        <IconButton onClick={() => {
          setValue(value+1)
        }}><Add /></IconButton>
        <TextField
          key={props.id}
          value={value}
          variant="outlined"
          onChange={e => setValue(parseInt(e.target.value)||0)}
          type="tel"
          margin="none"
          size="medium"
          InputProps={{
            disableUnderline : true,
            className : classes.input,
            inputProps : {
              style : {
                border : 'none',
                textAlign : 'center'
              }
            }

          }}

        />
        <IconButton onClick={() => {
          setValue(value => Math.max(value - 1,0))
        }}><Remove /></IconButton>
      </Box>
    </Box>
  )
}
export default CustomUpDownWidget;