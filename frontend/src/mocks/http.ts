import { http } from 'msw'

const API_URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

export const api = {
    get: (path: string, resolver: Parameters<typeof http.get>[1]) =>
        http.get(`${API_URL}${path}`, resolver),

    post: (path: string, resolver: Parameters<typeof http.post>[1]) =>
        http.post(`${API_URL}${path}`, resolver),

    put: (path: string, resolver: Parameters<typeof http.put>[1]) =>
        http.put(`${API_URL}${path}`, resolver),

    patch: (path: string, resolver: Parameters<typeof http.patch>[1]) =>
        http.patch(`${API_URL}${path}`, resolver),

    delete: (path: string, resolver: Parameters<typeof http.delete>[1]) =>
        http.delete(`${API_URL}${path}`, resolver),
}