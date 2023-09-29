import React from "react";
import { VerifyContainer, VerifyContent } from "../styles";
import { CircularProgress, Typography } from "@mui/material";
import { useRouter } from "next/router";
import useNotificationHandlers from "../../../utility/notification-handlers";
import axios from "axios";
import { GITHUB_APP_URL } from "../github-connect/constants";
import { VerifyGithubInstallation } from "../github-connect/github";

export default function ConnectionVerify() {
  const router = useRouter();
  const { handleError } = useNotificationHandlers();
  React.useEffect(() => {
    const { id, githubLogin } = router.query;
    VerifyGithubInstallation()
      .then(res => {
        router.push(`/connect/github/new?githubLogin=${githubLogin}&id=${id}`);
      })
      .catch(err => {
        handleError(err);
      });
  }, []);

  return (
    <VerifyContainer>
      <VerifyContent elevation={3}>
        <Typography variant="body1">Verifying your Github Installation......</Typography>
        <CircularProgress color="inherit" />
      </VerifyContent>
    </VerifyContainer>
  );
}
