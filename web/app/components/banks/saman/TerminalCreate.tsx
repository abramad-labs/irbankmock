import { createTerminal } from "@/app/clients/banks/saman/saman";
import { Button, Field, Group, GroupProps, Input } from "@chakra-ui/react";
import { useState } from "react";

export type TerminalCreateProps = {
    groupProps?: GroupProps;
    terminalCreateFinalized?: () => void;
};

export const TerminalCreate = (props: TerminalCreateProps) => {
    const [submitting, setSubmitting] = useState(false);
    const [terminalName, setTerminalName] = useState("");
    const onClick = () => {
        setSubmitting(true);
        if (!submitting) {
            createTerminal({ name: terminalName })
                .then((x) => {
                    setTerminalName("");
                })
                .catch((x) => {
                    alert(x)
                      
                })
                .finally(() => {
                    setSubmitting(false);
                    if(props.terminalCreateFinalized) {
                        props.terminalCreateFinalized()
                    }
                });
        }
    };

    return (
        <Group {...props.groupProps}>
            <Field.Root required>
                <Field.Label>Terminal Name</Field.Label>
                <Input
                    value={terminalName}
                    onChange={(v) => setTerminalName(v.target.value)}
                    required
                    placeholder="Something"
                />
            </Field.Root>
            <Button alignSelf="end" onClick={onClick} loading={submitting}>
                Create
            </Button>
        </Group>
    );
};
