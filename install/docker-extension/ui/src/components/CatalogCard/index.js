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
} from "./styles";
import Tooltip from "@mui/material/Tooltip";
import ServiceMeshIcon from "../../assets/service-mesh";
import WhiteDesignIcon from "../../assets/design-white";
import DesignIcon from "../../assets/design";
import CopyIcon from "../../assets/copy-button";

function CatalogCard({ pattern, patternType, catalog }) {
  const [copied, setCopied] = useState(false);

  const handleCopyClick = event => {
    event.preventDefault();

    navigator.clipboard.writeText(`MESHERY${pattern.id.split("-")[2]}`);
    setCopied(true);

    setTimeout(() => {
      setCopied(false);
    }, 2000);
  };

  return pattern?.visibility == "published" ? (
    <a href={`https://meshery.layer5.io/catalog/${pattern.id}`} target="_blank" rel="noreferrer">
      <DesignCard>
        <DesignInnerCard>
          <CardFront>
            <DesignVisibility>{pattern.visibility}</DesignVisibility>
            <DesignType>{patternType}</DesignType>
            <DesignName>{pattern.name}</DesignName>
            <ServiceMeshIcon width={120} height={100} />
            <DesignId>
              <DesignIcon width={20} height={20} />
              {` MESHERY${pattern.id.split("-")[2]}`}
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
    </a>
  ) : (
    <DesignCard>
      <DesignInnerCard>
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
            {catalog && <DesignDetails> Created By : {pattern.first_name}</DesignDetails>}
            <DesignDetails>Created At : {pattern.created_at.slice(0, 10)}</DesignDetails>
            <DesignDetails>Updated At : {pattern.updated_at.slice(0, 10)}</DesignDetails>
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
  );
}

export default CatalogCard;
