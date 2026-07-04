import "@/config/apiClient";
import { beforeAll, afterEach, afterAll, vi } from 'vitest'
import { server } from '@/mocks/node.js'

/**
 * Mock i18n to avoid test setup complexity and locale dependencies.
 * We only verify that correct translation keys are used.
 */
vi.mock("react-i18next", () => ({
    useTranslation: () => ({ t: (key: string) => key }),
}));
 
beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())