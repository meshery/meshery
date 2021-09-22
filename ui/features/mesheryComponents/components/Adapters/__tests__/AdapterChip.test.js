/**
 * @jest-environment jsdom
 */

import React from "react";
import { AdapterChip } from "@/features/mesheryComponents/components/Adapters/AdapterChip";
import { render, screen, fireEvent } from "@/utils/test-utils";

// each adapter should have a logo

/** @type {import("@/features/mesheryComponents/mesheryComponentsSlice").AdapterType} */
const mockAdapter1 = {
  adapter_location: "mesherylocal.layer5.io:10000",
  isActive: false,
  name: "istio",
};

/** @type {import("@/features/mesheryComponents/mesheryComponentsSlice").AdapterType} */
const mockAdapter2 = {
  adapter_location: "mesherylocal.layer5.io:10000",
  isActive: true,
  name: "istio",
};

// TODO: Test for click event

describe("Adapter chip", () => {
  it("renders an inactive adapter chip properly with the correct label", () => {
    render(<AdapterChip adapter={mockAdapter1} />);
    screen.getByText("mesherylocal.layer5.io");
  });

  it("renders the correct tooltip text for inactive adapter", async () => {
    render(<AdapterChip adapter={mockAdapter1} />);
    const chip = screen.getByTestId("adapter-chip");
    fireEvent.mouseOver(chip);
    await screen.findByText("This adapter is inactive");
  });

  it("renders the correct tooltip text for active adapter", async () => {
    render(<AdapterChip adapter={mockAdapter2} />);
    const chip = screen.getByTestId("adapter-chip");
    fireEvent.mouseOver(chip);
    await screen.findByText("Active");
  });
});
