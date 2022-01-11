import React from "react";
import { IconButton, TextField, Typography } from "@material-ui/core";
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

  if (omitTitleFields.includes(name)) {
    console.log("Fuck that's cool", prettifiedName)
    return <div key={props.id} style={style}>
      <TextField label={prettifiedName} variant="outlined" size="small" autoFocus key={props.id} value={props.value} id={props.id} onChange={e => props?.onChange(e.target.value)} placeholder={`${prettifiedName}`} />
    </div>
  }

  return (
    <div key={props.id} style={style}>
      <Typography variant="body1" style={{ fontWeight : "bold" }}>{prettifiedName}
        {props.schema?.description && (
          <EnlargedTextTooltip title={props.schema?.description}>
            <IconButton component="span" size="small">
              <HelpOutlineIcon />
            </IconButton>
          </EnlargedTextTooltip>
        )}
      </Typography>
      <TextField variant="outlined" size="small" autoFocus key={props.id} value={props.value} id={props.id} onChange={e => props?.onChange(e.target.value)} placeholder={`${prettifiedName}`} />
    </div>
  )
}

const MemoizedCustomInputField = React.memo(CustomInputField)

export default MemoizedCustomInputField;
