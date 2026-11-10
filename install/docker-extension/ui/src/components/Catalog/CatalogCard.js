import React, { useState } from "react";
import {
  CardBack,
  CardFront,
  DesignCard,
  DesignInnerCard,
  DesignVisibility,
  DesignType,
  DesignName,
  DesignId,
  DesignDetailsDiv,
  DesignDetails,
  CopyButton
} from "./style";
import Tooltip from "@mui/material/Tooltip";
import ServiceMeshIcon from "../../assets/service-mesh";
import WhiteDesignIcon from "../../assets/design-white";
import DesignIcon from "../../assets/design";
import CopyIcon from "../../assets/copy-button";
import { mesheryCloudUrl } from '../utils/constants';

function CatalogCard({ pattern, patternType, catalog }) {
  const [copied, setCopied] = useState(false);

  const handleCopyClick = async event => {
    event.preventDefault();

    try {
      await navigator.clipboard.writeText(
        `${mesheryCloudUrl}/catalog/${pattern.id}`
      );
      setCopied(true);
      setTimeout(() => {
        setCopied(false);
      }, 2000);
    } catch (error) {
      // Handle clipboard write error
      console.error("Error copying to clipboard:", error);
    }
  };

  const OpenDesignInCatalog = (designId) => {
    if(!designId) return
    window.ddClient.host.openExternal(
                      `${mesheryCloudUrl}/catalog/content/design/${designId}`,
                    )
  }

  return (
      <DesignCard>
        <DesignInnerCard onClick={() => OpenDesignInCatalog(pattern?.id)}>
          <CardFront>
            <DesignVisibility>{pattern.visibility}</DesignVisibility>
            <DesignType>{patternType}</DesignType>
            <DesignName>{pattern.name}</DesignName>
            <ServiceMeshIcon width={120} height={100} />
            <DesignId>
              <DesignIcon width={20} height={20} />
              {`MESHERY${pattern.id.split("-")[2]}`}
            </DesignId>
          </CardFront>
          <CardBack>
            <DesignVisibility
              style={{
                color: "white",
                border: `1px solid #00d3a9`
              }}
            >
              {pattern.visibility}
            </DesignVisibility>
            <DesignType>{patternType}</DesignType>
            <DesignName style={{ color: "white" }}>{pattern.name}</DesignName>
            <DesignDetailsDiv>
              {catalog && (
                <DesignDetails> Created By : {pattern.first_name}</DesignDetails>
              )}
              <DesignDetails>
                Created At : {pattern.created_at.slice(0, 10)}
              </DesignDetails>
              <DesignDetails>
                Updated At : {pattern.updated_at.slice(0, 10)}
              </DesignDetails>
              <DesignDetails>Version : v1</DesignDetails>
            </DesignDetailsDiv>
            <DesignId style={{ color: "white" }}>
              <WhiteDesignIcon width={20} height={20} />
              {` MESHERY${pattern.id.split("-")[2]}`}

              <Tooltip title={copied ? "Copied" : "Copy"}>
                <CopyButton onClick={handleCopyClick}>
                  <CopyIcon width={20} height={20} />
                </CopyButton>
              </Tooltip>
            </DesignId>
          </CardBack>
        </DesignInnerCard>
      </DesignCard>
  )
}

export default CatalogCard;
