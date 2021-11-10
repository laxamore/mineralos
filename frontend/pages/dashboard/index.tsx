import type { NextPage } from 'next'
import { useEffect, useState } from 'react'
import { checkAuth, withAuth, jwtObject, getAuthPayload } from "../../utils/auth"

import Navbar from '../../components/navbar'

const Dashboard: NextPage = () => {
    const [isAuth, setIsAuth] = useState(false)
    const [privilege, setPrivilege] = useState('readOnly')

    useEffect(() => {
        checkAuth().then(auth => {
            if (auth) {
                setIsAuth(true)
                const authPayload = getAuthPayload()
                setPrivilege(authPayload.privilege)
            }
        })
    }, [])

    return (
        <div>
            {isAuth ?
                <>
                    <Navbar />
                    <div className="flex flex-col justify-center items-center h-full w-full">
                        <div className="justify-start flex flex-col w-3/4 m-24 bg-gray-700 text-white rounded-2xl">
                            <div className="flex flex-row p-4 border-b-2 border-blue-700">
                                <button className={`p-2 bg-blue-600 rounded-lg ${privilege === 'admin' || privilege === 'readAndWrite' ? 'hover:bg-blue-700' : 'opacity-50 cursor-default'}`} onClick={() => {
                                    const authPayload = getAuthPayload()
                                    console.log(authPayload.privilege)
                                }} disabled={privilege === 'admin' || privilege === 'readAndWrite' ? false : true}>Create New Rig</button>
                            </div>
                            <ul className="mt-2">
                                <li className="border border-green-600 p-8 rounded-2xl cursor-pointer mt-2">
                                    <div className="flex flex-row items-center">
                                        <div className="w-1/6">
                                            <p className="text-lg font-bold">RIG_NAME</p>
                                            <p className="text-sm">0 -/s</p>
                                        </div>
                                        <div className="rounded-lg py-2 mx-2 w-5/6 h-12 bg-gray-800">

                                        </div>
                                        <div className="flex justify-end ml-4">
                                            <button className={`bg-red-500 rounded-full py-2 px-4 ${privilege === 'admin' || privilege === 'readAndWrite' ? 'hover:bg-red-600' : 'opacity-50 cursor-default'}`}
                                                disabled={privilege === 'admin' || privilege === 'readAndWrite' ? false : true}>X</button>
                                        </div>
                                    </div>
                                </li>
                                {/* <li className="border border-green-600 p-8 rounded-lg cursor-pointer mt-2">Rig 2</li>
                                <li className="border border-green-600 p-8 rounded-lg cursor-pointer mt-2">Rig 3</li>
                                <li className="border border-green-600 p-8 rounded-lg cursor-pointer mt-2">Rig 4</li>
                                <li className="border border-green-600 p-8 rounded-lg cursor-pointer mt-2">Rig 5</li>
                                <li className="border border-green-600 p-8 rounded-lg cursor-pointer mt-2">Rig 6</li> */}
                            </ul>
                        </div>
                    </div>
                </>
                :
                <></>
            }
        </div >
    )
}

export default Dashboard