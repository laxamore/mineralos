import type { NextPage } from "next"
import { FormEvent, useState } from "react"
import Link from "next/link"

const Register: NextPage = () => {
    const [registerState, setRegisterState] = useState({
        registerState: false,
        registerSuccess: false,
        registerMsg: '',
    })

    const registerHandler = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        const formData = new FormData(e.currentTarget)
        const data = {
            username: formData.get("username"),
            email: formData.get("email"),
            registerToken: formData.get("registerToken"),
            password: formData.get("password")
        }

        const response = await fetch(`${process.env.API_ENDPOINT}/api/v1/register`, {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            body: JSON.stringify({
                username: data.username,
                email: data.email,
                password: data.password
            }),
            headers: {
                regToken: `${data.registerToken}`
            }
        });

        if (response.status === 200) {
            setRegisterState({
                registerState: true,
                registerSuccess: true,
                registerMsg: 'Register Success',
            })
        }
        else if (response.status === 409) {
            setRegisterState({
                registerState: true,
                registerSuccess: false,
                registerMsg: 'Register Fail, User Already Exist',
            })
        }
        else {
            setRegisterState({
                registerState: true,
                registerSuccess: false,
                registerMsg: 'Register Fail',
            })
        }
    }

    return (
        <div className="absolute top-0 left-0 w-full h-full flex flex-row justify-center items-center">
            <div className="w-96 bg-gray-900 rounded-xl text-white font-extrabold text-xl">
                <div className="flex flex-col w-full justify-center items-center">
                    <h1 className="py-6">MineralOS</h1>

                    <form action="submit" onSubmit={registerHandler} className="mb-4">
                        <div className="flex flex-col text-sm font-semibold">
                            {
                                registerState.registerState ?
                                    <p className={`text-center ${registerState.registerSuccess ? 'text-green-500' : 'text-red-500'}`}>
                                        {registerState.registerMsg}
                                    </p> :
                                    <></>
                            }
                            <label htmlFor="username" className="py-2">Username</label>
                            <input className="text-black pl-2" type="text" name="username" />

                            <label htmlFor="email" className="py-2">Email</label>
                            <input className="text-black pl-2" type="email" name="email" />

                            <label htmlFor="registerToken" className="py-2">Register Token</label>
                            <input className="text-black pl-2" type="text" name="registerToken" />

                            <label htmlFor="password" className="py-2">Password</label>
                            <input className="text-black pl-2" type="password" name="password" />
                            <button type="submit" className="bg-blue-800 py-2 mt-4">Register</button>
                            <Link href="/login">
                                <a className="text-center pt-4 text-blue-500 hover:text-blue-400 cursor-pointer">Back to Login.</a>
                            </Link>
                        </div>

                    </form>
                </div>
            </div>
        </div>
    )
}

export default Register