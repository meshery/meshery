import { AdaptersListContainer, ComponentsStatusContainer } from "@/features/mesheryComponents";
import { nanoid } from "@reduxjs/toolkit";
import React from "react";

export default function Home() {
  return (
    <div>
      <ComponentsStatusContainer render={({ components }) => <div>{components.meshsync.connectionStatus}</div>} />
      <AdaptersListContainer
        render={({ adapters }) => (
          <ul>
            {adapters.map((ad) => (
              <li key={nanoid()}>{ad.adapter_location}</li>
            ))}
          </ul>
        )}
      />
    </div>
  );
}
