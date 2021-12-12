import { ClickAwayListener, makeStyles, TextField } from "@material-ui/core"
import React, { useState } from "react"

const useStyles = makeStyles(() => ({
  root: {
    width: "100%",
    color: "#fff",
    position: "fixed",
    top: 80,
    // left: 0,
    backgroundColor: "#477E96",
    padding: "4px 50px",
    transform: 'translateX(-40px)',
  },
  span: {
    // marginLeft: 300,
    color: "#fff",
    fontStyle: "italic",
    "&:hover": {
      cursor: "pointer",
      textDecoration: "underline"
    }
  },
  input: {
    background: "transparent",
    border: "none",
    color: "#fff",
    textDecoration: "underline",
    '&:focus': {
      outline: "none",
      border: "none"
    }
  }
}))

export default function CustomBreadCrumb({ title, onBack }) {
  const classes = useStyles()
  const [editing, setEditing] = useState(false);
  const [name, setName] = useState(title);

  const handleInputChange = event => {
    setName(event.target.value)
  }

  return (
    <div className={classes.root}>
      {"> "}
      <span
        className={classes.span}
        onClick={onBack}
      >
        Patterns
      </span>
      {" > "}

      {editing
        ?
        <ClickAwayListener
          onClickAway={() => setEditing(false)}
        >
          <input
            className={classes.input}
            value={name}
            onChange={handleInputChange}
            autoFocus
          />
        </ClickAwayListener>
        : <span
          className={classes.span}
          onClick={() => setEditing(true)}
        >{name}</span>
      }
    </div>
  )
}