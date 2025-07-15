"use client";

import {
    Box,
    Container,
    Heading,
    Input,
    Stack,
    VStack,
    Text,
    HStack,
    Button,
} from "@chakra-ui/react";
import { ChangeEventHandler, EventHandler, useState } from "react";

export default function Home() {
    const [formData, setFormData] = useState({
        cardNumber: "",
        cvv: "",
        expiryMonth: "",
        expiryYear: "",
        captcha: "",
        cardPassword: "",
    });

    const handleChange: ChangeEventHandler<HTMLInputElement> = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    const handleSubmit = () => {
        console.log("Submitting payment:", formData);
    };

    const handleCancel = () => {
        setFormData({
            cardNumber: "",
            cvv: "",
            expiryMonth: "",
            expiryYear: "",
            captcha: "",
            cardPassword: "",
        });
    };

    return (
        <Container p={10}>
            <Heading mx="auto" mb={5}>
                Saman Bank Internet Payment Gateway (Test)
            </Heading>
            <HStack m="auto" height={500} maxWidth={900}>
                <Box
                    w="100%"
                    mx="auto"
                    p="4"
                    borderWidth="1px"
                    borderRadius="lg"
                    boxShadow="md"
                >
                    <Stack direction="column" gap="4" align="stretch">
                        <Box>
                            <Text mb="1" fontWeight="bold">
                                Card Number
                            </Text>
                            <Input
                                name="cardNumber"
                                type="text"
                                maxLength={16}
                                value={formData.cardNumber}
                                onChange={handleChange}
                                placeholder="Enter 16-digit card number"
                            />
                        </Box>

                        <Box>
                            <Text mb="1" fontWeight="bold">
                                CVV
                            </Text>
                            <Input
                                name="cvv"
                                type="password"
                                maxLength={4}
                                value={formData.cvv}
                                onChange={handleChange}
                                placeholder="Enter CVV"
                            />
                        </Box>

                        <HStack gap="4">
                            <Box flex="1">
                                <Text mb="1" fontWeight="bold">
                                    Expiry Month
                                </Text>
                                <Input
                                    name="expiryMonth"
                                    type="number"
                                    min={1}
                                    max={12}
                                    value={formData.expiryMonth}
                                    onChange={handleChange}
                                    placeholder="MM"
                                />
                            </Box>

                            <Box flex="1">
                                <Text mb="1" fontWeight="bold">
                                    Expiry Year
                                </Text>
                                <Input
                                    name="expiryYear"
                                    type="number"
                                    min={2024}
                                    max={2100}
                                    value={formData.expiryYear}
                                    onChange={handleChange}
                                    placeholder="YYYY"
                                />
                            </Box>
                        </HStack>

                        <Box>
                            <Text mb="1" fontWeight="bold">
                                Captcha (ignored)
                            </Text>
                            <Input
                                name="captcha"
                                type="text"
                                value={formData.captcha}
                                onChange={handleChange}
                                placeholder="Enter captcha text"
                            />
                        </Box>

                        <Box>
                            <Text mb="1" fontWeight="bold">
                                Card Password
                            </Text>
                            <Input
                                name="cardPassword"
                                type="password"
                                value={formData.cardPassword}
                                onChange={handleChange}
                                placeholder="Enter card password"
                            />
                        </Box>

                        <HStack gap="4" pt="2">
                            <Button colorScheme="red" onClick={handleCancel}>
                                Cancel
                            </Button>
                            <Button colorScheme="green" onClick={handleSubmit}>
                                Submit Payment
                            </Button>
                        </HStack>
                    </Stack>
                </Box>
                <Box
                    maxW="400px"
                    mx="auto"
                    w="100%"
                    h="100%"
                    p="4"
                    borderWidth="1px"
                    borderRadius="lg"
                    boxShadow="md"
                >
                    <Heading>Payment Details</Heading>
                </Box>
            </HStack>
        </Container>
    );
}
