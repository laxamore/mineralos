import { clearAuthToken, jwtObject, withAuth } from "../utils/auth"
import styles from "./navbar.module.css"
import Router from "next/router"
import { FC } from "react"
import Link from 'next/link'

const Navbar: FC = () => {
    const logoutHandler = () => {
        withAuth(async (token: jwtObject) => {
            const response = await fetch(`${process.env.API_ENDPOINT}/api/v1/logout`, {
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${token.jwt_token}`
                },
                credentials: "include"
            })

            const responseJSON = await response.json()
            if (responseJSON === "logout success") {
                clearAuthToken()
                Router.push("/login")
            }
        })
    }

    return <div className="sticky flex flex-row items-center w-full top-0 left-0 h-14 bg-gray-700 text-gray-200">
        <span className="px-8 font-extrabold text-xl">MineralOS</span>
        <ul className={`flex flex-row justify-end w-full mx-16 font-bold ${styles.menu}`}>
            <li>
                <Link href="/dashboard">
                    <a>Rigs</a>
                </Link>
            </li>
            <li>
                <Link href="/dashboard/overclocking">
                    <a>Overclokcing</a>
                </Link>
            </li>
            <li>
                <Link href="/dashboard/wallets">
                    <a>Wallets</a>
                </Link>
            </li>
            <li>
                <Link href="/dashboard/miners">
                    <a>Miners</a>
                </Link>
            </li>
            <li>
                <Link href="/dashboard/settings">
                    <a>Settings</a>
                </Link>
            </li>
        </ul>
        <span className="px-8 font-extrabold text-red-500 hover:text-red-600 cursor-pointer" onClick={logoutHandler}>Logout</span>
    </div>
}

export default Navbar