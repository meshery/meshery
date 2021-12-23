import React from "react";
import { IconButton, TextField, Typography } from "@material-ui/core";
import HelpOutlineIcon from "../HelpOutlineIcon";
import EnlargedTextTooltip from "../EnlargedTextTooltip";
import { formatString } from "../../helpers"

const CustomInputField = (props) => {
  const name = props?.name || props?.idSchema['$id']?.split('_')[1]
  const prettifiedName = formatString(name) || 'Input'
  return (
    <div key={props.id}>
      <Typography variant="body1" style={{ fontWeight : "bold" }}>{prettifiedName}
        {props.schema?.description && (
          <EnlargedTextTooltip title={props.schema?.description}>
            <IconButton component="span" size="small">
              <HelpOutlineIcon />
            </IconButton>
          </EnlargedTextTooltip>
        )}
      </Typography>
      <TextField variant="outlined" size="small" style={{ margin : '10px 0 ' }} autoFocus key={props.id} value={props.value} id={props.id} onChange={e => props?.onChange(e.target.value)} placeholder={`${prettifiedName}`}/>
    </div>
  )
}

const MemoizedCustomInputField = React.memo(CustomInputField)

export default MemoizedCustomInputField;
