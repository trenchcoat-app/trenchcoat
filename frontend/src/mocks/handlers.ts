import { HttpResponse } from 'msw'
import { api } from '@/mocks/http';
import type { Account, SignInBody, SignInOkResponse } from '@/api'
 
export const handlers = [
  api.post('/api/v1/auth/sign-in', async ({ request }) => {
    const body : SignInBody = await request.json() as SignInBody;
    const account: Account = {
      id: 'acc-123',
      email: body.email,
      displayName: 'John Maverick',
    }

    const response: SignInOkResponse = {
      account,
    }

    return HttpResponse.json(response)
  }),
]