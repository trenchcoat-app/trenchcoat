import { render as rtlRender } from "@testing-library/react";
import type { RenderOptions } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactElement, ReactNode } from "react";

// Replacement for React Testing Library's `render` that wraps components
// with QueryClientProvider. Use this for tests involving TanStack Query hooks.
export const renderWithQueryClient = ( ui: ReactElement, options?: Omit<RenderOptions, "wrapper"> ) => {
    const queryClient = new QueryClient({
        defaultOptions: {
            queries: { retry: false },
            mutations: { retry: false },
        },
    });

    const QueryClientWrapper = ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient}>
            {children}
        </QueryClientProvider>
    );

    return rtlRender(ui, {
        wrapper: QueryClientWrapper,
        ...options,
    });
};