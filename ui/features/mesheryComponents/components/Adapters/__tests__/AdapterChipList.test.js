/**
 * @jest-environment jsdom
 */

import React from "react";
import { render, screen } from "@/utils/test-utils";
import { AdaptersChipList } from "@/features/mesheryComponents";

/** @type {import("@/features/mesheryComponents/mesheryComponentsSlice").AdaptersListType} */
const mockAdaptersList = [
  {
    adapter_location: "mesherylocal.layer5.io:10000",
    isActive: false,
    name: "istio",
  },
  {
    adapter_location: "mesherylocal.layer5.io:10001",
    isActive: true,
    name: "istio",
  },
  {
    adapter_location: "mesherylocal.layer5.io:10002",
    isActive: true,
    name: "istio",
  },
  {
    adapter_location: "mesherylocal.layer5.io:10003",
    isActive: false,
    name: "istio",
  },
];

describe("Adapter chip list", () => {
  it("renders all the adapter chips", () => {
    render(<AdaptersChipList adapters={mockAdaptersList} />);
    const chips = screen.getAllByTestId("adapter-chip");
    expect(chips.length).toBe(4);
  });
});
