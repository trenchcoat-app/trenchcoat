// Browser-side MSW worker (used for local dev mocking, not tests)
import { setupWorker } from 'msw/browser';
import { handlers } from "@/mocks/handlers";
 
export const worker = setupWorker(...handlers);