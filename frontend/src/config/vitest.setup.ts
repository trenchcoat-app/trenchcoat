import "@/config/apiClient";
import { beforeAll, afterEach, afterAll, vi } from 'vitest'
import { server } from '@/mocks/node.js'

vi.mock("react-i18next", () => ({
    useTranslation: () => ({ t: (key: string) => key }),
}));
 
beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())