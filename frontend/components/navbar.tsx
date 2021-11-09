import type { NextPage } from "next"
import styles from "./navbar.module.css"

const Navbar: NextPage = () => {
    return <div className="flex flex-row items-center w-full absolute top-0 left-0 h-14 bg-gray-700 text-gray-200">
        <span className="px-8 font-extrabold text-xl">MineralOS</span>
        <ul className={`flex flex-row justify-end w-full mx-16 ${styles.menu}`}>
            <li>Rigs</li>
            <li>Overclocking</li>
            <li>Wallets</li>
            <li>Settings</li>
        </ul>
    </div>
}

export default Navbar