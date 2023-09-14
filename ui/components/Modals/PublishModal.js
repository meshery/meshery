import React, { useState, useEffect } from "react";
import Modal from "../Modal";
import PublicIcon from "@material-ui/icons/Public";
import _ from "lodash";
import { getMeshModels } from "../../api/meshmodel";
import { modifyRJSFSchema } from "../../utils/utils";
import dataFetch from "../../lib/data-fetch";

// This modal is used in MeshMap also
export default function PublishModal(props) {
  const { open, title, handleClose, handleSubmit } = props;
  const [publishSchema, setPublishSchema] = useState({});

  useEffect(() => {
    dataFetch(
      "/api/schema/resource/publish",
      {
        method : "GET",
        credentials : "include",
      },
      async (result) => {
        try {
          const { models } = await getMeshModels();
          const modelNames = _.uniq(models?.map((model) => model.displayName));

          // Modify the schema using the utility function
          const modifiedSchema = modifyRJSFSchema(result.rjsfSchema, "properties.compatibility.items.enum", modelNames);

          setPublishSchema({ rjsfSchema : modifiedSchema, uiSchema : result.uiSchema });
        } catch (err) {
          console.error(err);
          setPublishSchema(result);
        }
      }
    );
  },[]);

  return (
    <Modal
      open={open}
      schema={publishSchema.rjsfSchema}
      uiSchema={publishSchema.uiSchema}
      title={title}
      handleClose={handleClose}
      handleSubmit={handleSubmit}
      submitBtnText="Submit for Approval"
      submitBtnIcon={<PublicIcon data-cy="import-button" />}
      showInfoIcon={{
        text : "Upon submitting your catalog item, an approval flow will be initiated.",
        link : "https://docs.meshery.io/concepts/catalog",
      }}
    />
  );
}