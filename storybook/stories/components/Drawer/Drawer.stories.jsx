import React from "react";
import { Drawer } from "./Drawer";

export default {
  title : "Components/Drawer",
  component : Drawer,
  tags: ['autodocs'],
};

export function TemporaryDrawerTop() {
  return (
    <Drawer direction="top" />
  )
}

export function TemporaryDrawerLeft() {
  return (
    <Drawer direction="left" />
  )
}

export function TemporaryDrawerBottom() {
  return (
    <Drawer direction="bottom" />
  )
}

export function TemporaryDrawerRight() {
  return (
    <Drawer direction="right" />
  )
}