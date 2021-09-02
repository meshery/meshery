import { AdaptersChipList, AdaptersListContainer, ComponentsStatusContainer } from "@/features/mesheryComponents";
import React from "react";

export default function Home() {
  return (
    <div>
      <ComponentsStatusContainer render={({ components }) => <div>{components.meshsync.connectionStatus}</div>} />
      <AdaptersListContainer
        render={({ adapters, loading }) => <AdaptersChipList adapters={adapters} loading={loading} />}
      />
    </div>
  );
}
