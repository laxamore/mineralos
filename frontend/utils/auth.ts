import Router from 'next/router'

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

export function checkAuth() {
    return new Promise(async (resolve) => {
        if (inMemoryToken.jwt_token != "") {
            resolve(inMemoryToken)
        }
        else {
            const response = await fetch(`${process.env.API_ENDPOINT}/api/v1/refreshToken`, {
                method: 'POST',
                mode: 'cors',
                cache: 'no-cache',
                credentials: 'include',
            })

            const responseJSON = await response.json()

            if (responseJSON.jwt_token) {
                inMemoryToken = responseJSON
                resolve(inMemoryToken)
            }
            else {
                Router.push("/login")
            }
        }
    })
}

export function withAuth(handler: Function) {
    if (inMemoryToken.jwt_token != "") {
        const isTokenExpired = (inMemoryToken.jwt_token_expiry - Math.floor(Date.now() / 1000) - 10) < 0 ? true : false

        if (isTokenExpired) {
            fetch(`${process.env.API_ENDPOINT}/api/v1/refreshToken`, {
                method: 'POST',
                mode: 'cors',
                cache: 'no-cache',
                credentials: 'include',
            }).then(response => {
                response.json().then(responseJSON => {
                    if (responseJSON.jwt_token) {
                        inMemoryToken = responseJSON
                        handler(inMemoryToken)
                    }
                    else {
                        Router.push("/login")
                    }
                })
            })
        } else {
            handler(inMemoryToken)
        }
    }
    else {
        Router.push("/login")
    }
}

export function clearAuthToken() {
    inMemoryToken = {
        jwt_token: "",
        jwt_token_expiry: 0
    }
}