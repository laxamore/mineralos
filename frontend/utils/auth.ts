import Router from 'next/router'
import Login from '../pages';

export interface jwtObject {
    jwt_token: string
    jwt_token_expiry: number
}

let inMemoryToken: jwtObject = {
    jwt_token: "",
    jwt_token_expiry: 0
};

export function login(jwt: { jwt_token: string, jwt_token_expiry: number }, noRedirect = false) {
    inMemoryToken = {
        jwt_token: jwt.jwt_token,
        jwt_token_expiry: jwt.jwt_token_expiry
    };

    if (!noRedirect) {
        Router.push('/dashboard')
    }
}

export function checkAuth(ssr?: boolean, rtoken_cookies?: string) {
    return new Promise(async (resolve) => {
        if ((inMemoryToken.jwt_token_expiry * 1000) < (Date.now() + 10000)) {
            let options: RequestInit = {
                method: 'POST',
                mode: 'cors',
                cache: 'no-cache',
                credentials: 'include',
            }

            if (ssr) {
                options.headers = {
                    Cookie: `rtoken=${rtoken_cookies}`
                }
            }

            const response = await fetch(`${ssr ? process.env.API_ENDPOINT_SSR : process.env.API_ENDPOINT}/api/v1/refreshToken`, options)
            const responseJSON = await response.json()

            if (responseJSON.jwt_token) {
                inMemoryToken = responseJSON
                resolve(inMemoryToken)
            }
            else if (!ssr) {
                Router.push("/login")
            }
            else {
                resolve(false);
            }
        }
        else {
            resolve(inMemoryToken);
        }
    })
}

export function withAuth(handler: Function, ssr?: boolean, rtoken_cookie?: string) {
    if ((inMemoryToken.jwt_token_expiry * 1000) < (Date.now() + 10000)) {
        if (!ssr) {
            checkAuth().then(() => {
                handler(inMemoryToken)
            })
        }
        else {
            checkAuth(ssr, rtoken_cookie).then(() => {
                handler(inMemoryToken)
            })
        }
    } else {
        handler(inMemoryToken)
    }
}

export function clearAuthToken() {
    inMemoryToken = {
        jwt_token: "",
        jwt_token_expiry: 0
    }
}

export function getAuthPayload() {
    const encodedPayload = inMemoryToken.jwt_token.split('.')[1]
    return JSON.parse(Buffer.from(encodedPayload, 'base64').toString())
}