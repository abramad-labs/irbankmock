"use client";

import { TerminalCreate } from "@/components/banks/saman/TerminalCreate";
import { TerminalList } from "@/components/banks/saman/TerminalList";
import {
    Container,
    Heading,
} from "@chakra-ui/react";
import { v4 } from "uuid";
import { useState } from "react";

export default function Home() {
    const [refershKey, setRefreshKey] = useState<undefined | string>(undefined)
    const refreshList = () => {
        setRefreshKey(v4())
    }
    return (
        <Container p={10}>
            <Heading mb={5}>
                Saman Bank Internet Payment Gateway Management
            </Heading>
            <Heading size="lg">Create Terminal</Heading>
            <TerminalCreate groupProps={{mb: 2}} terminalCreateFinalized={refreshList} />
            <Heading size="lg">Terminal List</Heading>
            <TerminalList refreshKey={refershKey}/>
        </Container>
    );
}
