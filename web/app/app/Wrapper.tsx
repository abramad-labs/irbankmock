"use client";

import { Provider } from "@/components/ui/provider";
import { Toaster } from "@/components/ui/toaster";
import { ReactNode } from "react";

export const Wrapper = (props: { children: ReactNode }) => {
    return (
        <Provider defaultTheme="light">
            {props.children}
            <Toaster />
        </Provider>
    );
};
