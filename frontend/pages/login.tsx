import type { NextPage } from 'next'
import { useRouter } from 'next/router'
import { FormEvent } from 'react'
import { login } from "../utils/auth"
import Link from "next/link"

const Login: NextPage = () => {
    const Router = useRouter()

    const loginHandler = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        const formData = new FormData(e.currentTarget)
        const data = {
            username: formData.get("username"),
            password: formData.get("password")
        }
        const response = await fetch('http://localhost:5000/api/v1/login', {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            body: JSON.stringify(data),
            credentials: 'include',
        })

        if (response.status === 200) {
            const responseJSON = await response.json()
            login(responseJSON)
        }
    }

    return (
        <div className="absolute top-0 left-0 w-full h-full flex flex-row justify-center items-center">
            <div className="w-96 bg-gray-900 rounded-xl text-white font-extrabold text-xl">
                <div className="flex flex-col w-full justify-center items-center">
                    <h1 className="py-6">MineralOS</h1>

                    <form action="submit" onSubmit={loginHandler}>
                        <div className="flex flex-col text-sm font-semibold">
                            <label htmlFor="username" className="py-2">Username</label>
                            <input className="text-black pl-2" type="text" name="username" />
                            <label htmlFor="password" className="py-2">Password</label>
                            <input className="text-black pl-2" type="password" name="password" />
                            <button type="submit" className="bg-blue-800 py-2 mt-4">Login</button>
                            <Link href="/register">
                                <a className="text-center py-4 text-blue-500 hover:text-blue-400 cursor-pointer">Register.</a>
                            </Link>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}

export default Login
