"use client";

import { Provider } from "@/components/ui/provider";
import { ReactNode } from "react";


export const Wrapper = (props: { children: ReactNode }) => {
    return (
       <Provider defaultTheme="light">{props.children}</Provider>
    );
};
