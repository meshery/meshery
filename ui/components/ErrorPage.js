import React, { useState, useEffect } from 'react';
import {
  NoSsr,
  Typography,
  Link
} from "@mui/material";
import { styled } from "@mui/material/styles"

const DivRootClass = styled("div")(() => ({
  backgroundColor : "#fff",
  padding : "2rem",
  textAlign : "center",
  borderRadius : 4,
  height : "100%",
}));

const DivErrorSection = styled("div")(() => ({}));

const DivMessage = styled("div")(() => ({
  fontSize : "3rem",
  lineHeight : "2rem",
  marginBottom : "2rem",
}));

const DivErrMessage = styled("div")(() => ({
  fontWeight : "400",
  fontSize : "1.5rem",
  color : "gray",
  fontStyle : "italic",
  marginTop : "2.5rem",
}));

const ImgMesh = styled("img")(() => ({
  display : "block",
  margin : "auto",
  marginTop : "3.125rem",
  maxWidth : "50%",
  height : "45%",
}));

const PHelpMessage = styled("p")(() => ({
  marginTop : "5rem",
  color : "rgba(0, 0, 0, 0.87)",
}));

const customMessages = [
  "Oh, no. Please pardon our meshy app.",
  "Oops. Please excuse the mesh.",
  "Things tend to get a bit meshy around here.",
  "Please pardon our mesh.",
  "How did this mesh happen?",
  "Well, isn't this a mesh?",
  "Yikes. Things are a mesh here.",
];

function CustomErrorMessage() {
  const [customMessage, setCustomMessage] = useState(customMessages[0]);

  useEffect(() => {
    setCustomMessage(
      customMessages[Math.floor(Math.random() * customMessages.length)]
    );
  }, []);

  return (
    <NoSsr>
      <DivRootClass>
        <DivErrorSection>
          <Typography variant="h1">
            <DivMessage>{customMessage}</DivMessage>
          </Typography>
          <Typography variant="h5">
            <DivErrMessage>Page does not exist.</DivErrMessage>
          </Typography>
        </DivErrorSection>
        <ImgMesh src="/static/img/service-mesh.svg" alt="service meshed" />
        <Typography variant="body1">
          <PHelpMessage>
            Start a conversation at Layer5 community{" "}
            <Link
              underline="none"
              href="https://discuss.layer5.io/c/meshery/5"
              target="_blank"
            >
              discussion forum
            </Link>
            .
          </PHelpMessage>
        </Typography>
      </DivRootClass>
    </NoSsr>
  );
}

export default CustomErrorMessage;

