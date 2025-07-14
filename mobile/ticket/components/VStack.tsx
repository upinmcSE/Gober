import React from "react";
import { Stack, StackProps } from "@/components/Stack";
import { ViewStyle } from "react-native";

interface VStackProps extends StackProps {
  style?: ViewStyle; // thêm prop style
}

export function VStack(props: VStackProps) {
  return (
    <Stack {...props} direction="column">
      {props.children}
    </Stack>
  );
}