import { useQuery } from '@tanstack/react-query';
import { pingOptions } from '@api/@tanstack/react-query.gen';

/**
 * Minimal component to demonstrate how to consume the auto-generated
 * Tanstack Query Hooks from the OpenAPI schema.
 */
export function PingComponent() {
    const { data, error, isLoading } = useQuery(pingOptions());

    if (isLoading) return <p>Ping Loading...</p>;
    if (error) return <p>Ping Error: {error.message}</p>;

    return <div>Ping API Message: {data?.message}</div>;
}

