import { createTerminal } from "@/app/clients/banks/saman/saman";
import { toaster } from "@/components/ui/toaster";
import { CommonError } from "@/types/errors";
import { Button, Field, Group, GroupProps, Input } from "@chakra-ui/react";
import { AxiosError } from "axios";
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
                .catch((err: AxiosError<CommonError>) => {
                    const error = err.response?.data?.error ?? err.message
                    toaster.create({
                        title: "Terminal",
                        description: `Error creating terminal: ${error}`,
                        type: "error",
                    });
                })
                .finally(() => {
                    setSubmitting(false);
                    if (props.terminalCreateFinalized) {
                        props.terminalCreateFinalized();
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
