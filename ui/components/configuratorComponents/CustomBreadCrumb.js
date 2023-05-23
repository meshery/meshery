import { ClickAwayListener, withStyles } from "@material-ui/core"
import React, { useEffect, useState } from "react"

const styles = (theme) => ({
  designWrapper : {
    width : "100%",
    color : "#fff",
    position : "fixed",
    top : 80,
    // left: 0,
    backgroundColor : theme.palette.type === 'dark' ? "#222222" : "#477E96",
    zIndex : "1",
    marginLeft : "4px",
    padding : "4px 50px",
    transform : 'translateX(-40px)',
  },
  span : {
    // marginLeft: 300,
    color : "#fff",
    fontStyle : "italic",
    "&:hover" : {
      cursor : "pointer",
      textDecoration : "underline"
    }
  },
  input : {
    background : "transparent",
    border : "none",
    color : "#fff",
    textDecoration : "underline",
    '&:focus' : {
      outline : "none",
      border : "none"
    }
  }
});

function CustomBreadCrumb({ title, onBack, titleChangeHandler, classes }) {
  const [editing, setEditing] = useState(false);
  const [name, setName] = useState(title);

  useEffect(() => {
    const timer = setTimeout(() => {
      titleChangeHandler(name.trim())
    }, 400);

    return () => clearTimeout(timer)
  }, [name])

  useEffect(() => {
    setName(title)
  }, [title])

  const handleInputChange = event => {
    setName(event.target.value)
  }

  return (
    <div className={classes.designWrapper}>
      {"> "}
      <span
        className={classes.span}
        onClick={onBack}
      >
        Designs
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
        >{title}</span>
      }
    </div>
  )
}
export default withStyles(styles)(CustomBreadCrumb);