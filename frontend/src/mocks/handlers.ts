import { HttpResponse } from 'msw'
import { api } from '@/mocks/http';
import type { Account, Session, SignInBody, SignInOkResponse } from '@/api'
 
export const handlers = [
  api.post('/api/v1/auth/sign-in', async ({ request }) => {
    const body : SignInBody = await request.json() as SignInBody;
    const account: Account = {
      id: 'acc-123',
      email: body.email,
      displayName: 'John Maverick',
    }
    // Temporarily added since OpenAPI spec change has not yet been merged
    const session: Session = {
      token: 'mock-jwt-token',
      expiresAt: new Date(
        Date.now() + 1000 * 60 * 60 * 24 // 24 hours
      ).toISOString(),
      accountId: account.id,
    }

    const response: SignInOkResponse = {
      session,
      account,
    }

    return HttpResponse.json(response)
  }),
]