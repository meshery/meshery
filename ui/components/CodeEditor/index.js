/* eslint-disable */
import React from "react";
import dynamic from "next/dynamic";
const CodeEditor = dynamic(import("./CodeEditor"), {
  ssr: false,
});

const CodeEditorNoSSR = ({ defaultValue = "", name }) => (
  <div>
    <CodeEditor defaultValue={defaultValue} name={name} />
  </div>
);

export default CodeEditorNoSSR;
