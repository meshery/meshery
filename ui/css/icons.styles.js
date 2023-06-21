// icon styles, setting general height and width properties to solves scaling and consistency problems

export const iconSmall = {
  height : 20,
  width : 20
}

export const iconMedium = {
  height : 24,
  width : 24
}

export const iconLarge = {
  height : 32,
  width : 32
}

export const iconXLarge = {
  height : 40,
  width : 40
}

export const extensionStyles = (theme) => ({
  button : {
    borderRadius : 5,
    minWidth : 100,
    color : "#fff",
    "&:hover" : {
      boxShadow : "0px 2px 4px -1px rgb(0 0 0 / 20%), 0px 4px 5px 0px rgb(0 0 0 / 14%), 0px 1px 10px 0px rgb(0 0 0 / 12%)"
    },
  },
  card : {
    padding : theme.spacing(3),
    borderRadius : theme.spacing(1),
    transformStyle : "preserve-3d",
    boxShadow : "0 4px 8px 0 rgba(0,0,0,0.2)",
    backgroundColor : theme.palette.secondary.elevatedComponents,
    minHeight : "250px",
    position : "relative",
  },
  a : {
    textDecoration : "none",
  },
  img : {
    paddingRight : "1rem",
    height : "auto",
    width : "auto",
    maxWidth : "120px",
    maxHeight : "120px",
  },
  frontSideDescription : {
    paddingTop : "1rem",
    paddingBottom : "1rem",
    textAlign : "left",
    display : "flex",
    flexDirection : 'row',
  },
  link : {
    textDecoration : "none",
    color : theme.palette.secondary.link2,
  },
  switchBase : {
    color : '#647881',
    "&$checked" : { color : '#00b39f' },
    "&$checked + $track" : { backgroundColor : 'rgba(0,179,159,0.5)' },
  },
  track : { backgroundColor : 'rgba(100,120,129,0.5)', },
  checked : {},
});