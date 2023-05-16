import { useEffect, useRef, useState } from "react";
import { EditorView } from "@codemirror/view";
import { basicSetup } from "codemirror";

export default {
    title: "Components/CodeMirror",
      tags: ['autodocs'],
};

function useCodeMirror() {
    const ref = useRef();
    const [view, setView] = useState();

    useEffect(() => {
        const view = new EditorView({
            extensions: [
                basicSetup,
            ],
            parent: ref.current,
        });

        setView(view);

        return () => {
            view.destroy();
            setView(undefined);
        };
    }, []);

    return { ref, view };
};

function onUpdate(onChange) {
    return EditorView.updateListener.of((viewUpdate) => {
        if (viewUpdate.docChanged) {
            const doc = viewUpdate.state.doc;
            const value = doc.toString();
            onChange(value, viewUpdate);
        }
    });
}

function useCodeEditor({ value, onChange }) {
    const { ref, view } = useCodeMirror([onUpdate(onChange)]);

    useEffect(() => {
        if (view) {
            const editorValue = view.state.doc.toString();

            if (value !== editorValue) {
                view.dispatch({
                    changes: {
                        from: 0,
                        to: editorValue.length,
                        insert: value || "",
                    },
                });
            }
        }
    }, [value, view]);

    return ref;
}

function CodeMirrorComponent({ value, onChange }) {
    const ref = useCodeEditor({ value, onChange });

    return <div ref={ref} />;
};

export function GenericEditor() {
    const [code, setCode] = useState("console.log");

    return (
        <CodeMirrorComponent
            value={code}
            onChange={(newCode) => {
                setCode(newCode);
            }}
        />
    )
}