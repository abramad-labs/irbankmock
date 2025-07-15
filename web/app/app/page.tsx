import {
    Button,
    Card,
    Center,
    Container,
    Grid,
    Group,
    Heading,
    Link as ChakraLink,
    Stack,
} from "@chakra-ui/react";
import Link from "next/link";

export default function Home() {
    return (
        <Container p={10}>
            <main>
                <Center>
                    <Stack>
                        <Heading>IR Bank Mock</Heading>
                        <Group>
                            <Card.Root>
                                <Card.Header>
                                    <Heading>Saman</Heading>
                                </Card.Header>
                                <Card.Body>
                                    Manage Saman Electronic Payment
                                </Card.Body>
                                <Card.Footer>
                                    <Button asChild>
                                        <ChakraLink asChild>
                                            <Link href="/banks/saman/manage">Manage</Link>
                                        </ChakraLink>
                                    </Button>
                                </Card.Footer>
                            </Card.Root>
                        </Group>
                    </Stack>
                </Center>
            </main>
            <footer></footer>
        </Container>
    );
}
