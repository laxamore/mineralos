import type { GetServerSideProps, NextPage } from 'next'
import { useEffect, useMemo, useState } from 'react'
import { checkAuth, withAuth, jwtObject, getAuthPayload } from "../../utils/auth"

import Navbar from '../../components/navbar'
import CreateRigModal from '../../components/modals/createRigModal'
import { Accessor, Column, useTable } from 'react-table'

const isServer = () => typeof window === 'undefined';

type Props = {
    data?: any;
}

const Dashboard: NextPage<Props> = ({ data }) => {
    const [isAuth, setIsAuth] = useState(false)
    const [privilege, setPrivilege] = useState('readOnly')
    const [showCreateRigModal, setShowCreateRigModal] = useState(false);
    const [rigsData, setRigsData] = useState([]);

    const rigsDataTable = useMemo(() => rigsData, [rigsData])
    interface rigsColumnInterface {
        status: string
        rig_name: string,
        hashrate: string,
        units: string,
    }

    const columns = useMemo<Column<rigsColumnInterface>[]>(() => [
        {
            Header: 'status',
            accessor: 'status',
        },
        {
            Header: 'rig_name',
            accessor: 'rig_name',
        },
        {
            Header: 'hashrate',
            accessor: 'hashrate'
        },
        {
            Header: 'units',
            accessor: 'units',
        }
    ], []);

    useEffect(() => {
        checkAuth().then(auth => {
            if (auth) {
                setIsAuth(true)
                const authPayload = getAuthPayload()
                setPrivilege(authPayload.privilege)
            }
        })
        setRigsData(data.rigs)
    }, [])

    const tableInstance = useTable({ columns: columns, data: rigsDataTable })
    const {
        rows,
    } = tableInstance

    return (
        <div>
            {isAuth ?
                <>
                    <Navbar />
                    <div className="flex flex-col justify-center items-center h-full w-full">
                        <div className="justify-start flex flex-col w-3/4 m-24 bg-gray-700 text-white rounded-2xl">
                            <div className="flex flex-row p-4 border-b-2 border-blue-700">
                                <button className={`p-2 bg-blue-600 rounded-lg ${privilege === 'admin' || privilege === 'readAndWrite' ? 'hover:bg-blue-700' : 'opacity-50 cursor-default'}`}
                                    onClick={() => {
                                        console.log(rigsData)
                                        console.log(rows)
                                        setShowCreateRigModal(true)
                                    }}
                                    disabled={privilege === 'admin' || privilege === 'readAndWrite' ? false : true}>
                                    Create New Rig
                                </button>
                            </div>
                            <ul className="mt-2">
                                {
                                    rows.map((val: any) => {
                                        return <li className="border border-green-600 p-3 rounded-2xl cursor-pointer mt-2" key={val.original.rig_id}>
                                            <div className="flex flex-row items-center">
                                                <div className="w-1/6">
                                                    <p className="text-lg font-bold">{val.original.rig_name}</p>
                                                    <p className="text-sm">0 -/s</p>
                                                </div>
                                                <div className="rounded-lg py-2 mx-2 w-5/6 h-12 bg-gray-800">

                                                </div>
                                                <div className="flex justify-end ml-4">
                                                    <button className={`bg-red-500 rounded-full px-2 ${privilege === 'admin' || privilege === 'readAndWrite' ? 'hover:bg-red-600' : 'opacity-50 cursor-default'}`}
                                                        disabled={privilege === 'admin' || privilege === 'readAndWrite' ? false : true} onClick={() => {
                                                            withAuth(async (token: jwtObject) => {
                                                                const response = await fetch(`${process.env.API_ENDPOINT}/api/v1/deleteRig`, {
                                                                    method: 'DELETE',
                                                                    mode: 'cors',
                                                                    body: JSON.stringify({
                                                                        rig_id: val.original.rig_id,
                                                                    }),
                                                                    headers: {
                                                                        Authorization: `Bearer ${token.jwt_token}`,
                                                                    },
                                                                })

                                                                if (response.status == 200) {
                                                                    setRigsData(rigsData.filter((rig: any) => rig.rig_id != val.original.rig_id))
                                                                }
                                                            })
                                                        }}>X</button>
                                                </div>
                                            </div>
                                        </li>
                                    })
                                }
                            </ul>
                        </div>
                    </div>

                    {showCreateRigModal ?
                        <CreateRigModal setShowModal={setShowCreateRigModal} createRigSuccessHandler={(res: never) => {
                            setRigsData([...rigsData, res])
                        }} /> : null
                    }
                </>
                :
                <></>
            }
        </div >
    )
}


export const getServerSideProps: GetServerSideProps = async (ctx: any) => {
    const rtoken = ctx.req.cookies['rtoken']
    const getRigsData = () => {
        return new Promise(resolve => {
            withAuth(async (token: jwtObject) => {
                const response = await fetch(`${process.env.API_ENDPOINT}/api/v1/getRigs`, {
                    method: 'GET',
                    mode: 'cors',
                    headers: {
                        Authorization: `Bearer ${token.jwt_token}`
                    }
                })

                if (response.status == 200) {
                    resolve(await response.json())
                }
                else {
                    resolve({ status: response.status })
                }
            }, isServer(), rtoken)
        })
    }

    const data: any = await getRigsData()
    return { props: { data } }
}

export default Dashboard