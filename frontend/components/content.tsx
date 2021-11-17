import { FC, useContext, useEffect } from "react"
import { RefreshIcon, ClockIcon } from "@heroicons/react/solid"

import { ContentContext, RefreshContext } from "../utils/context"

type Params = {
    privilege: string
    showCreateButton?: boolean
    showRefreshButton: boolean
    showRefreshButtonTimeout: boolean
    children?: React.ReactNode
}


export const Content: FC<Params> = ({ children, privilege, showCreateButton, showRefreshButton, showRefreshButtonTimeout }) => {
    const [setShowCreateRigModal] = useContext(ContentContext);
    const [refreshTimeout, setRefreshTimeout] = useContext(RefreshContext)

    return <div className="justify-start flex flex-col w-3/4 m-24 bg-gray-700 text-white rounded-2xl">
        <div className="flex flex-row p-4 border-b-2 items-center border-blue-700">
            {
                showCreateButton ?
                    <button className={`p-2 bg-blue-600 rounded-lg ${privilege === 'admin' || privilege === 'readAndWrite' ? 'hover:bg-blue-700' : 'opacity-50 cursor-default'}`}
                        onClick={() => {
                            setShowCreateRigModal(true)
                        }}
                        disabled={privilege === 'admin' || privilege === 'readAndWrite' ? false : true}>
                        Create New Rig
                    </button> : null
            }
            <div className="flex flex-grow justify-end">
                {
                    showRefreshButton ?
                        <RefreshIcon className={`justify-end p-2 w-8 h-8 mx-2 bg-blue-600 rounded-lg hover:bg-blue-700 cursor-pointer select-none`}
                            onClick={() => {
                                setShowCreateRigModal(true)
                            }}>
                        </RefreshIcon> : null
                }

                {
                    showRefreshButtonTimeout ?
                        <ClockIcon className={`justify-end p-2 w-8 h-8 mx-2 bg-blue-600 rounded-lg hover:bg-blue-700 cursor-pointer select-none ${!refreshTimeout ? 'opacity-50' : ''}`}
                            onClick={() => {
                                window.localStorage.setItem('refreshState', `${!refreshTimeout}`)
                                setRefreshTimeout(!refreshTimeout)
                            }}>
                        </ClockIcon> : null
                }
            </div>
        </div>
        {children}
    </div>
}