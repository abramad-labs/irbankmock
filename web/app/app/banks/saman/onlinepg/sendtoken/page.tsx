"use client";

import { fetcherWithError, ResponseError } from "@/lib/fetcher";
import { SamanPublicTokenInfoResponse, SuccessErrorPair } from "@/types/banks/saman/types";
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
    ProgressCircle,
    EmptyState,
} from "@chakra-ui/react";
import { intervalToDuration, parseISO } from "date-fns";
import { useSearchParams } from "next/navigation";
import { ChangeEventHandler, useEffect, useState } from "react";
import { MdOutlineMoneyOff } from "react-icons/md";
import useSWR from "swr";

const InvalidPayment = (props: { title: string; description: string }) => {
    return (
        <EmptyState.Root size="lg">
            <EmptyState.Content>
                <EmptyState.Indicator>
                    <MdOutlineMoneyOff />
                </EmptyState.Indicator>
                <VStack textAlign="center">
                    <EmptyState.Title>{props.title}</EmptyState.Title>
                    <EmptyState.Description>
                        {props.description}
                    </EmptyState.Description>
                </VStack>
            </EmptyState.Content>
        </EmptyState.Root>
    );
};

export default function Home() {
    const searchParams = useSearchParams()
    const token = searchParams.get("token") ?? ""
    const { data, error, isLoading, mutate } =
        useSWR<SamanPublicTokenInfoResponse, ResponseError<SuccessErrorPair>>(
            `/banks/saman/management/public/token/?token=${encodeURIComponent(
                token
            )}`,
            fetcherWithError
        );
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


 const [remaining, setRemaining] = useState("");

  useEffect(() => {
    if(!data?.expiresAt)
        return

    const target = parseISO(data?.expiresAt);

    function update() {
      const now = new Date();
      if (now >= target) {
        setRemaining("00:00");
        mutate()
        return;
      }

      const dur = intervalToDuration({ start: now, end: target });
      const mm = String((dur.hours || 0) * 60 + (dur.minutes || 0)).padStart(2, "0");
      const ss = String(dur.seconds || 0).padStart(2, "0");
      setRemaining(`${mm}:${ss}`);
    }

    update();
    const id = setInterval(update, 1000);
    return () => clearInterval(id);
  }, [data?.expiresAt]);

    if (token === "")
        return (
            <InvalidPayment
                title="Invalid Payment"
                description="invalid token was provided"
            />
        );

    if (error)
        return <InvalidPayment title="Invalid Payment" description={error.data?.error ?? "failure in fetching token data"} />;

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
                            <Button colorPalette="red" onClick={handleCancel}>
                                Cancel
                            </Button>
                            <Button colorPalette="green" onClick={handleSubmit}>
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
                    <Heading mb={5}>Payment Details</Heading>
                    <Heading size="md">Remaining</Heading>
                    <Text mb={2}>{remaining}</Text>
                    <Heading size="md">Terminal Name</Heading>
                    <Text mb={2}>{data?.terminalName}</Text>
                    <Heading size="md">Terminal ID</Heading>
                    <Text mb={2}>{data?.terminalId}</Text>
                    <Heading size="md">Website</Heading>
                    <Text mb={2}>{data?.website}</Text>
                    <Heading size="md">Amount</Heading>
                    <Text>{data?.amount} IRR</Text>

                </Box>
            </HStack>
        </Container>
    );
}
