import { styled } from "@mui/material/styles";
import { Box } from "@mui/material";

export const ChartDiv = styled(Box)(() => ({
  padding: "1rem",
  borderRadius: "1rem",
  width: "50%",
  height: "23rem",
  marginRight: "1rem",
  display: "block",
}));

export const DesignCardUrl = styled("a")(() => ({
  textDecoration: "none"
}));

export const CopyButton = styled("div")(() => ({
  textDecoration: "none"
}));

export const DesignCard = styled("div")(() => ({
  width: "14rem",
  height: "17rem",
  border: "1px solid rgba(0, 0, 0, 0.15)",
  position: "relative",
  borderRadius: "1rem",
  textAlign: "center",
  transformStyle: "preserve-3d",
  boxShadow: "0 4px 8px 0 rgba(0,0,0,0.2)",
  marginBottom: "1.25rem",
  marginLeft: "10px",
  marginRight: "10px",
  display: "block",
  perspective: "1000px",
  transition: "all .9s",
  "&:hover": {
    cursor: "pointer",
    transform: "rotateY(180deg)"
  },
}));

export const DesignInnerCard = styled("div")(() => ({
  position: "relative",
  width: "100%",
  height: "100%",
  textAlign: "center",
  transition: "transform 0.6s",
  transformStyle: "preserve-3d",
  boxShadow: "0 4px 8px 0 rgba(0,0,0,0.2)",
  borderRadius: "0.9375rem"
}));

export const CardFront = styled("div")(() => ({
  boxShadow: "2px 2px 6px 0px #00d3a9",
  position: "absolute",
  width: "100%",
  height: "100%",
  WebkitBackfaceVisibility: "hidden",
  borderRadius: "0.9375rem",
  backfaceVisibility: "hidden",
  background:
    "linear-gradient(to left bottom, #f1f3f5, #f4f5f7, #f7f7f9, #ffffff, #ffffff, #ffffff, #ffffff, #ffffff, #ffffff, #f7f7f9, #f4f5f7, #f1f3f5);"
}));

export const CardBack = styled("div")(() => ({
  boxShadow: "2px 2px 6px 0px #00d3a9",
  position: "absolute",
  width: "100%",
  height: "100%",
  WebkitBackfaceVisibility: "hidden",
  borderRadius: "0.9375rem",
  backfaceVisibility: "hidden",
  color: "white",
  transform: "rotateY(180deg)",
  background: "linear-gradient(250deg, #477e96 0%, #00b39f 35%, rgb(60, 73, 79) 100%)"
}));

export const DesignType = styled("span")(() => ({
  position: "absolute",
  top: "0",
  right: "0",
  minWidth: "3rem",
  padding: "0 0.75rem",
  fontSize: "0.875rem",
  textTransform: "capitalize",
  color: "black",
  borderRadius: "0 1rem 0 2rem",
  background: "rgb(0, 211, 169)"
}));

export const DesignName = styled("h5")(() => ({
  textTransform: "capitalize",
  color: "rgba(0, 0, 0, 0.75)",
  fontSize: "1.125rem",
  marginTop: "4rem",
  marginBottom: "1rem",
  padding: "0rem 1.5rem",
  position: "relative",
  overflow: "hidden",
  whiteSpace: "nowrap",
  textOverflow: "ellipsis",
  textAlign: "center",
  "& :after": {
    content: "''",
    textAlign: "right",
    position: "absolute",
    bottom: "0",
    right: "0",
    width: "70%",
    background:
      "linear-gradient(to right, rgba(255, 255, 255, 0), rgba(255, 255, 255, 1) 50%)"
  }
}));

export const DesignId = styled("p")(() => ({
  color: "rgba(0, 0, 0, 0.45)",
  fontSize: "0.75rem",
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  gap: "0.5rem",
  position: "absolute",
  bottom: "15px",
  width: "100%"
}));

export const DesignVisibility = styled("div")(() => ({
  border: "1px solid gray",
  padding: "0.1rem 0.5rem",
  position: "absolute",
  top: 30,
  left: 20,
  fontSize: "0.7rem",
  textTransform: "capitalize",
  color: "rgba(26, 26, 26, .8)",
  borderRadius: "0.4rem"
}));

export const DesignDetailsDiv = styled("div")(() => ({
  height: "7.5rem",
  display: "flex",
  marginTop: "-1rem",
  flexDirection: "column",
  padding: "0rem 1.5rem",
  justifyContent: "start",
  alignItems: "start",
}));

export const DesignDetails = styled("p")(() => ({
  color: "white",
  fontSize: "0.9rem",
  fontWeight: "semi-bold",
  lineHeight: "1.75rem",
  textOverflow: "ellipsis",
  whiteSpace: "nowrap",
  margin: 0,
  width: "100%",
  overflow: "hidden",
  textAlign: "left"
}));