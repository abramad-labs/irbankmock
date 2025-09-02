import { fetcher } from "@/lib/fetcher";
import { SamanTerminalsResponse } from "@/types/banks/saman/types";
import { ProgressCircle, Table } from "@chakra-ui/react";
import useSWR from "swr";

export const TerminalList = (props: {}) => {
    const { data, error, isLoading } = useSWR<SamanTerminalsResponse>(
        "/banks/saman/management/terminal",
        fetcher
    );

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
    );
};
