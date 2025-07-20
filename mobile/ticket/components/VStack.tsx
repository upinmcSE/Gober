import { Stack, StackProps } from "@/components/Stack";
import React from "react";

interface VStackProps extends StackProps { }

export function VStack(props: VStackProps) {
  return (
    <Stack { ...props } direction="column">
      { props.children }
    </Stack>
  );
}