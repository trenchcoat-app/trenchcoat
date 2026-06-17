import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: '../openapi/bundled.yaml',
  output: './src/api',
  plugins: [
    '@hey-api/client-fetch',
    '@tanstack/react-query',
  ],
});
