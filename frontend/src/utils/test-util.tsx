import { render as rtlRender } from "@testing-library/react";
import type { RenderOptions } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ReactElement, ReactNode } from "react";

const createTestQueryClient = () =>
    new QueryClient({
        defaultOptions: {
            queries: { retry: false },
            mutations: { retry: false },
        },
    });

// Wraps components with QueryClientProvider which is needed for tests on components that use TanStack Query.
const QueryClientWrapper = ({ children }: { children: ReactNode }) => {
    const queryClient = createTestQueryClient();
    return (
        <QueryClientProvider client={queryClient}>
            {children}
        </QueryClientProvider>
    );
};

const customRender = (ui: ReactElement, options?: Omit<RenderOptions, "wrapper">) =>
    rtlRender(ui, { wrapper: QueryClientWrapper, ...options });

export * from "@testing-library/react";
export { customRender as render };