// Node-side MSW server used by vitest (see src/config/vitest.setup.ts)
import { setupServer } from 'msw/node'
import { handlers } from "@/mocks/handlers";
 
export const server = setupServer(...handlers)