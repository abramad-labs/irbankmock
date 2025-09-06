"use client";

import {
    cancelToken,
    failToken,
    submitToken,
} from "@/clients/banks/saman/saman";
import { toaster } from "@/components/ui/toaster";
import { fetcherWithError, ResponseError } from "@/lib/fetcher";
import {
    BankSepTokenFinalizeResponse,
    SamanPublicTokenInfoResponse,
    SuccessErrorPair,
} from "@/types/banks/saman/types";
import { CommonError } from "@/types/errors";
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
    Center,
    Grid,
    GridItem,
} from "@chakra-ui/react";
import { AxiosError } from "axios";
import { intervalToDuration, parseISO } from "date-fns";
import { useSearchParams } from "next/navigation";
import { ChangeEventHandler, Suspense, useEffect, useState } from "react";
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

const BankInputForm = ({
    setFinalizeResponse,
}: {
    setFinalizeResponse: (resp: BankSepTokenFinalizeResponse) => void;
}) => {
    const searchParams = useSearchParams();
    const token = searchParams.get("token") ?? "";
    const { data, error, isLoading, mutate } = useSWR<
        SamanPublicTokenInfoResponse,
        ResponseError<SuccessErrorPair>
    >(
        `/banks/saman/public/token/?token=${encodeURIComponent(token)}`,
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

    const [formLoading, setFormLoading] = useState(false);

    const handleChange: ChangeEventHandler<HTMLInputElement> = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    const handleSubmit = () => {
        setFormLoading(true);
        submitToken({
            token: token,
            captcha: formData.captcha,
            cardNumber: formData.cardNumber,
            cardPassword: formData.cardPassword,
            cvv: parseInt(formData.cvv, 10),
            expiryMonth: parseInt(formData.expiryMonth, 10),
            expiryYear: parseInt(formData.expiryYear, 10),
        })
            .then((resp) => {
                setFinalizeResponse(resp.data);
            })
            .catch((err: AxiosError<CommonError>) => {
                const error = err.response?.data?.error ?? err.message;
                toaster.create({
                    title: "Payment Error",
                    description: `Error finalizing payment: ${error}`,
                    type: "error",
                });
                setFormLoading(false);
            });
    };

    const handleFail = () => {
        setFormLoading(true);
        failToken({
            token: token,
        })
            .then((resp) => {
                setFinalizeResponse(resp.data);
            })
            .catch((err: AxiosError<CommonError>) => {
                const error = err.response?.data?.error ?? err.message;
                toaster.create({
                    title: "Payment Error",
                    description: `Error finalizing payment: ${error}`,
                    type: "error",
                });
                setFormLoading(false);
            });
    };

    const handleCancel = () => {
        setFormLoading(true);
        cancelToken({
            token: token,
        })
            .then((resp) => {
                setFinalizeResponse(resp.data);
            })
            .catch((err: AxiosError<CommonError>) => {
                const error = err.response?.data?.error ?? err.message;
                toaster.create({
                    title: "Payment Error",
                    description: `Error finalizing payment: ${error}`,
                    type: "error",
                });
                setFormLoading(false);
            });
    };

    const [remaining, setRemaining] = useState("");

    useEffect(() => {
        if (!data?.expiresAt) return;

        const target = parseISO(data?.expiresAt);

        function update() {
            const now = new Date();
            if (now >= target) {
                setRemaining("00:00");
                mutate();
                return;
            }

            const dur = intervalToDuration({ start: now, end: target });
            const mm = String(
                (dur.hours || 0) * 60 + (dur.minutes || 0)
            ).padStart(2, "0");
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
        return (
            <InvalidPayment
                title="Invalid Payment"
                description={
                    error.data?.error ?? "failure in fetching token data"
                }
            />
        );

    if (isLoading)
        return (
            <ProgressCircle.Root value={null} size="sm">
                <ProgressCircle.Circle>
                    <ProgressCircle.Track />
                    <ProgressCircle.Range />
                </ProgressCircle.Circle>
            </ProgressCircle.Root>
        );

    const disableSubmit =
        formData.cardNumber.length !== 16 || !/\d+/.test(formData.cardNumber);
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
                                disabled={formLoading}
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
                                type="number"
                                maxLength={4}
                                disabled={formLoading}
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
                                    disabled={formLoading}
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
                                    disabled={formLoading}
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
                                disabled={formLoading}
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
                                disabled={formLoading}
                                onChange={handleChange}
                                placeholder="Enter card password"
                            />
                        </Box>

                        <HStack gap="4" pt="2">
                            <Button
                                colorPalette="pink"
                                disabled={formLoading}
                                onClick={handleCancel}
                            >
                                Cancel
                            </Button>
                            <Button
                                colorPalette="red"
                                onClick={handleFail}
                                disabled={disableSubmit || formLoading}
                            >
                                Fail
                            </Button>
                            <Button
                                colorPalette="green"
                                onClick={handleSubmit}
                                disabled={disableSubmit || formLoading}
                            >
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
};

const RedirectPageView = ({
    tokenFinalizeResponse,
}: {
    tokenFinalizeResponse: BankSepTokenFinalizeResponse;
}) => {
    const cb = tokenFinalizeResponse.callbackData;
    return (
        <Container p={10}>
            <Heading mx="auto" mb={5}>
                Saman Bank Internet Payment Gateway (Test)
            </Heading>
            <Stack m="auto" height={500} maxWidth={900}>
                <Center
                    w="100%"
                    mx="auto"
                    p="4"
                    borderWidth="1px"
                    borderRadius="lg"
                    boxShadow="md"
                >
                    <Stack direction="column" gap="4" align="stretch">
                        <Heading>
                            You are being redirected to the merchant website...
                        </Heading>
                        <Center textAlign="center" gap="4" pt="2">
                            <form
                                method="post"
                                action={tokenFinalizeResponse?.redirectURL}
                            >
                                <input
                                    type="hidden"
                                    name="MID"
                                    value={cb.MID}
                                />
                                <input
                                    type="hidden"
                                    name="TerminalId"
                                    value={cb.terminalId}
                                />
                                <input
                                    type="hidden"
                                    name="AffectiveAmount"
                                    value={cb.affectiveAmount}
                                />
                                {/*
                                    The typo in `OrginalAmount` is intended and originally roots back
                                    to the older versions of the SEP document.
                                    They removed this key from the newer versions
                                    but we kept it due to compatibility reasons.
                                 */}
                                <input
                                    type="hidden"
                                    name="OrginalAmount"
                                    value={cb.amount}
                                />
                                <input
                                    type="hidden"
                                    name="Amount"
                                    value={cb.amount}
                                />
                                <input
                                    type="hidden"
                                    name="HashedCardNumber"
                                    value={cb.hashedCardNumber}
                                />
                                <input
                                    type="hidden"
                                    name="RefNum"
                                    value={cb.refNum}
                                />
                                <input
                                    type="hidden"
                                    name="ResNum"
                                    value={cb.resNum}
                                />
                                <input
                                    type="hidden"
                                    name="RRN"
                                    value={cb.rrn}
                                />
                                <input
                                    type="hidden"
                                    name="SecurePan"
                                    value={cb.securePan}
                                />
                                <input
                                    type="hidden"
                                    name="State"
                                    value={cb.state}
                                />
                                <input
                                    type="hidden"
                                    name="Status"
                                    value={cb.status}
                                />
                                <input
                                    type="hidden"
                                    name="Token"
                                    value={cb.token}
                                />
                                <input
                                    type="hidden"
                                    name="TraceNo"
                                    value={cb.traceNo}
                                />
                                <input
                                    type="hidden"
                                    name="Wage"
                                    value={cb.wage}
                                />
                                <Button type="submit" colorPalette="blue">
                                    Redirect Now
                                </Button>
                            </form>
                        </Center>
                    </Stack>
                </Center>
                <Box
                    mx="auto"
                    w="100%"
                    h="100%"
                    p="4"
                    borderWidth="1px"
                    borderRadius="lg"
                    boxShadow="md"
                >
                    <Heading mb={5}>Payment Details</Heading>
                    <Grid templateColumns="repeat(2, 1fr)">
                        <GridItem>
                            <Heading size="md">Terminal ID</Heading>
                            <Text mb={2}>{cb.terminalId}</Text>
                            <Heading size="md">Amount</Heading>
                            <Text mb={2}>{cb.amount}</Text>
                            <Heading size="md">Card</Heading>
                            <Text mb={2}>{cb.securePan}</Text>
                            <Heading size="md">Reference Num</Heading>
                            <Text mb={2}>{cb.refNum}</Text>
                            <Heading size="md">Reservation Num</Heading>
                            <Text>{cb.resNum}</Text>
                        </GridItem>
                        <GridItem>
                            <Heading size="md">RRN</Heading>
                            <Text mb={2}>{cb.rrn}</Text>
                            <Heading size="md">State</Heading>
                            <Text mb={2}>{cb.state}</Text>
                            <Heading size="md">Status</Heading>
                            <Text mb={2}>{cb.status}</Text>
                            <Heading size="md">Trace No</Heading>
                            <Text mb={2}>{cb.traceNo}</Text>
                            <Heading size="md">Wage</Heading>
                            <Text>{cb.wage} IRR</Text>
                        </GridItem>
                    </Grid>
                </Box>
            </Stack>
        </Container>
    );
};

export default function Home() {
    const [pageMode, setPageMode] = useState<"input" | "redirect">("input");
    const [tokenFinalizeResponse, setTokenFinalizeResponse] = useState<
        BankSepTokenFinalizeResponse | undefined
    >(undefined);

    const setFinalizeResponseCallback = (
        resp: BankSepTokenFinalizeResponse
    ) => {
        setTokenFinalizeResponse(resp);
    };

    useEffect(() => {
        if (tokenFinalizeResponse) {
            setPageMode("redirect");
        }
    }, [tokenFinalizeResponse]);

    return pageMode === "input" ? (
        <Suspense>
            <BankInputForm setFinalizeResponse={setFinalizeResponseCallback} />
        </Suspense>
    ) : tokenFinalizeResponse ? (
        <RedirectPageView tokenFinalizeResponse={tokenFinalizeResponse} />
    ) : (
        <InvalidPayment
            title="Unexpected Error"
            description="Server didn't respond with correct data"
        />
    );
}
