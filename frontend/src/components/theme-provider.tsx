"use client";

import * as React from "react";
import { ThemeProvider as NextThemesProvider } from "next-themes";
import { type ThemeProviderProps } from "next-themes/dist/types";
import { usePathname } from "next/navigation";

export function ThemeProvider({ children }: ThemeProviderProps) {
  const pathname = usePathname();
  return (
    <NextThemesProvider
      forcedTheme={
        pathname === "/login" || pathname === "/register" ? "light" : undefined
      }
      attribute="class"
      defaultTheme="dark"
      enableSystem
    >
      {children}
    </NextThemesProvider>
  );
}
