import React, {
  useEffect,
  useState
} from 'react';
import Cookies from 'universal-cookie';
import { styled } from "@mui/material/styles"
import { Typography, Grid, Button, IconButton } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';

const DivPaper = styled("div")(({ theme }) => ({
  position : "fixed",
  width : 450,
  backgroundColor : theme.palette.background.paper,
  border : "0px solid #000",
  boxShadow : theme.shadows[5],
  padding : theme.spacing(1, 2, 3, 4),
  right : 0,
  bottom : 0,
  borderRadius : 10,

  ["@media (max-width: 455px)"] : {
    width : "100%",
  },

  zIndex : 5,
}));

const ImgDesignerImg = styled("img")(() => ({
  height : 300,
  margin : "auto",
}));

const TypographyHeader = styled(Typography)(() => ({
  paddingBottom : "0.5rem",
  paddingTop : "0.6rem",
  position : "absolute",
  fontWeight : "bold",

  ["@media (max-width: 455px)"] : {
    fontSize : "1rem",
  },
}));

const TypographyCaption = styled(Typography)(() => ({
  lineHeight : "1.2",
  paddingBottom : "15px",
  fontSize : ".75rem",
  textAlign : "center",
}));

const DivImgWrapper = styled("div")(() => ({
  padding : "15px 10px 15px 0",
  display : "flex",
}));

const DivHeaderWrapper = styled("div")(() => ({
  marginBottom : 12,
}));

function MeshMapEarlyAccessCardPopup() {
  const [isOpen, setIsOpen] = useState(true);
  const cookies = new Cookies("registered");

  const handleOpen = () => {
    const timer = setTimeout(() => {
      setIsOpen(true);
    }, 10000);
    return () => clearTimeout(timer);
  };

  useEffect(() => {
    if (cookies.get("registered")) {
      setIsOpen(false);
    } else if (!cookies.get("registered")) {
      cookies.set("registered", "true", {
        path : "/",
      });
      handleOpen();
    }
  }, []);

  if (isOpen) {
    return <MeshMapEarlyAccessCard closeForm={() => setIsOpen(false)} />;
  } else {
    return <></>;
  }
}

export function MeshMapEarlyAccessCard({
  rootStyle = {},
  closeForm = () => {},
}) {

  const handleSignUp = (e) => {
    window.open("https://layer5.io/meshmap", "_blank");
    e.stopPropagation();
  };

  return (
    <DivPaper style={rootStyle}>
      <DivHeaderWrapper>
        <TypographyHeader variant="h6">
          Get early access to MeshMap!
        </TypographyHeader>

        <div
          style={{
            display : "flex",
            justifyContent : "flex-end",
            whiteSpace : "nowrap",
            position : "relative",
          }}
        >
          <IconButton
            key="close"
            aria-label="Close"
            color="inherit"
            onClick={closeForm}
          >
            <CloseIcon />
          </IconButton>
        </div>
      </DivHeaderWrapper>

      <DivImgWrapper>
        <ImgDesignerImg src="/static/img/designer.png" />
      </DivImgWrapper>
      <TypographyCaption variant="subtitle1">
        <i>
          Friends dont let friends GitOps alone. Visually design and collaborate
          in real-time with other MeshMap users.
        </i>
      </TypographyCaption>
      <div style={{ display : "flex", justifyContent : "flex-end" }}>
        <Grid item xs={3}>
          <Button
            fullWidth
            variant="contained"
            color="primary"
            onClick={(e) => handleSignUp(e)}
          >
            Sign up
          </Button>
        </Grid>
      </div>
    </DivPaper>
  );
}

export default MeshMapEarlyAccessCardPopup;

