import { fetcher } from "@/lib/fetcher";
import { SamanTerminalsResponse } from "@/types/banks/saman/types";
import {
    Alert,
    Code,
    List,
    ProgressCircle,
    Stack,
    Table,
} from "@chakra-ui/react";
import { useEffect } from "react";
import useSWR from "swr";

export type TerminalListProps = {
    refreshKey?: string;
};

export const TerminalList = (props: TerminalListProps) => {
    const { data, error, isLoading, mutate } = useSWR<SamanTerminalsResponse>(
        "/banks/saman/management/terminal",
        fetcher
    );

    useEffect(() => {
        if (props.refreshKey) {
            mutate();
        }
    }, [props.refreshKey]);

    if (error) return <div>failed loading terminals: {error}</div>;

    if (isLoading)
        return (
            <ProgressCircle.Root value={null} size="sm">
                <ProgressCircle.Circle>
                    <ProgressCircle.Track />
                    <ProgressCircle.Range />
                </ProgressCircle.Circle>
            </ProgressCircle.Root>
        );

    return (
        <Stack>
            <Alert.Root status="info">
                <Alert.Indicator />
                <Alert.Title>
                    Use these paths to send request to terminal:
                    <List.Root>
                        <List.Item>
                            Payment Gateway: &#9;<Code>{data?.endpoints?.paymentGateway}</Code>
                        </List.Item>
                        <List.Item>
                            Payment Token: &#9;<Code>{data?.endpoints?.paymentToken}</Code>
                        </List.Item>
                        <List.Item>
                            Receipt: <Code>{data?.endpoints?.receipt}</Code>
                        </List.Item>
                        <List.Item>
                            Verify Transaction: <Code>{data?.endpoints?.verifyTransaction}</Code>
                        </List.Item>
                        <List.Item>
                            Reverse Transaction: <Code>{data?.endpoints?.reverseTransaction}</Code>
                        </List.Item>
                    </List.Root>
                </Alert.Title>
            </Alert.Root>
            <Table.Root size="sm">
                <Table.Header>
                    <Table.Row>
                        <Table.ColumnHeader>ID</Table.ColumnHeader>
                        <Table.ColumnHeader>Name</Table.ColumnHeader>
                        <Table.ColumnHeader>Username</Table.ColumnHeader>
                        <Table.ColumnHeader>Password</Table.ColumnHeader>
                        <Table.ColumnHeader textAlign="end">
                            Action
                        </Table.ColumnHeader>
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    {data?.terminals.map((item) => (
                        <Table.Row key={item.id}>
                            <Table.Cell>{item.id}</Table.Cell>
                            <Table.Cell>{item.name}</Table.Cell>
                            <Table.Cell>{item.username}</Table.Cell>
                            <Table.Cell>{item.password}</Table.Cell>
                            <Table.Cell textAlign="end"></Table.Cell>
                        </Table.Row>
                    ))}
                </Table.Body>
            </Table.Root>
        </Stack>
    );
};
