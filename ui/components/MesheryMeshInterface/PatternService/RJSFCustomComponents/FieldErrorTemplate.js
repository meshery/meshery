import React from "react";
import ListItem from "@material-ui/core/ListItem";
import FormHelperText from "@material-ui/core/FormHelperText";
import List from "@material-ui/core/List";

export default function FieldErrorTemplate(props) {
  const { errors = [], idSchema } = props;
  if (errors.length === 0) {
    return null;
  }
  const id = `${idSchema.$id}__error`;

  return (
    <List dense={true} disablePadding={true}>
      {errors.map((error, i) => {
        return (error==="is a required property" ? null:
          <ListItem key={i} disableGutters={true}>
            <FormHelperText id={id}>{error}</FormHelperText>
          </ListItem>
        )
      })}
    </List>
  );
}