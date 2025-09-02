"use client";

import { TerminalCreate } from "@/components/banks/saman/TerminalCreate";
import { TerminalList } from "@/components/banks/saman/TerminalList";
import {
    Button,
    Container,
    Field,
    Group,
    Heading,
    Input,
} from "@chakra-ui/react";

export default function Home() {
    
    return (
        <Container p={10}>
            <Heading mb={5}>
                Saman Bank Internet Payment Gateway Management
            </Heading>
            <Heading size="lg">Create Terminal</Heading>
            <TerminalCreate groupProps={{mb: 2}} />
            <Heading size="lg">Terminal List</Heading>
            <TerminalList />
        </Container>
    );
}
