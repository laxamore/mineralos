import { GetServerSideProps, NextPage } from "next"
import Router from "next/router";
import { useEffect, useMemo, useState } from "react"
import { checkAuth, getAuthPayload, jwtObject, withAuth } from "../../../utils/auth"
import Navbar from '../../../components/navbar'
import { Content } from '../../../components/content'
import { ContentContext } from "../../../utils/context"
import { EyeIcon, EyeOffIcon } from "@heroicons/react/solid"
import { Column, Row, useTable } from "react-table";
import OverclockModal from '../../../components/modals/overclockModal'

const isServer = () => typeof window === 'undefined';

type Params = {
    data: any
}

const Rigs: NextPage<Params> = ({ data }) => {
    const [isAuth, setIsAuth] = useState(false)
    const [privilege, setPrivilege] = useState('readOnly')
    const [showKey, setShowKey] = useState(false)
    const [gpusData, setGpusData] = useState<any>([]);
    const [gpu, setGpu] = useState({})
    const [showOverclockModal, setShowOverclockModal] = useState(false)

    let amdIndex: number = -1, nvidiaIndex: number = -1;

    interface gpusColumnInterface {
        gpuname: string,
        hashrate: string,
        fs: string,
        cc: string,
        cv: string,
        mc: string,
        mv: string,
        pl: string
    }

    const columns = useMemo<Column<gpusColumnInterface>[]>(() => [
        {
            Header: 'GPU Name',
            accessor: 'gpuname',
        },
        {
            Header: '',
            accessor: 'hashrate',
        },
        {
            Header: 'FS',
            accessor: 'fs',
        },
        {
            Header: 'CC',
            accessor: 'cc',
        },
        {
            Header: 'CV',
            accessor: 'cv',
        },
        {
            Header: 'MC',
            accessor: 'mc',
        },
        {
            Header: 'MV',
            accessor: 'mv',
        },
        {
            Header: 'PL',
            accessor: 'pl',
        },
    ], []);

    useEffect(() => {
        if (data.responseStatus === 404) {
            Router.push('/404.html')
        }
        checkAuth().then(auth => {
            if (auth) {
                setIsAuth(true)
                const authPayload = getAuthPayload()
                setPrivilege(authPayload.privilege)
            }
        })
        if (data.status) {
            if (data.status.gpus) {
                setGpusData(data.status.gpus)
            }
        }
        else {
            setGpusData([])
        }
    }, [])

    const tableInstance = useTable({ columns: columns, data: gpusData })
    const {
        getTableProps,
        getTableBodyProps,
        headerGroups,
        rows,
        prepareRow,
    } = tableInstance

    return <div>
        {isAuth ?
            <>
                <Navbar />

                <div className="flex flex-col justify-center items-center h-full w-full">
                    <ContentContext.Provider value={[]}>
                        <Content showRefreshButtonTimeout={true} showRefreshButton={true} privilege={privilege}>
                            <div className="flex flex-col font-semibold px-8 py-4">
                                <div className="flex flex-row">
                                    <p className="w-32">RIG ID</p>
                                    <p className={`${!showKey ? 'tracking-tight' : ''}`}>: {showKey ? data.rig_id : "∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗"}</p>

                                    <div className="w-6 mx-4 cursor-pointer" onClick={() => setShowKey(!showKey)}>
                                        {
                                            showKey ? <EyeOffIcon /> : <EyeIcon />
                                        }
                                    </div>

                                    <div className="flex flex-grow flex-row-reverse">
                                        <div className="flex flex-row px-4">
                                            <p className="pr-2">AMD Drivers:</p>
                                            <p>{data.status ? data.status.drivers.amd : ''}</p>
                                        </div>
                                        <div className="flex flex-row px-4">
                                            <p className="pr-2">NVIDIA Drivers:</p>
                                            <p>{data.status ? data.status.drivers.nvidia : ''}</p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <table {...getTableProps()} className="my-4">
                                <thead>
                                    {headerGroups.map(headerGroup => (
                                        <tr {...headerGroup.getHeaderGroupProps()}>
                                            {headerGroup.headers.map(column => (
                                                <th
                                                    {...column.getHeaderProps()}
                                                    className={
                                                        `p-3
                                                        ${column.id == "gpuname" ? 'text-left w-9/12' : 'text-center'}
                                                        ${column.id == "hashrate" ? 'w-2/12' : ''}
                                                        `
                                                    }
                                                >
                                                    {column.render('Header')}
                                                </th>
                                            ))}
                                        </tr>
                                    ))}
                                </thead>
                                <tbody {...getTableBodyProps()}>
                                    {rows.map((row: any, index) => {
                                        let ocID: number;

                                        if (row.original.gpuvendor == "AMD") {
                                            amdIndex++
                                            ocID = amdIndex;
                                        }
                                        else if (row.original.gpuvendor == "NVIDIA") {
                                            nvidiaIndex++
                                            ocID = nvidiaIndex;
                                        }

                                        prepareRow(row)
                                        return (
                                            <tr {...row.getRowProps()} className="border-2 border-gray-500 cursor-pointer"
                                                onClick={() => {
                                                    setShowOverclockModal(true)
                                                    setGpu({
                                                        idx: index,
                                                        rig_id: data.rig_id,
                                                        gpuvendor: row.original.gpuvendor,
                                                        ocID: ocID,
                                                        gpuinfo: row.values,
                                                    })
                                                }}
                                            >
                                                {row.cells.map((cell: any) => {
                                                    return (
                                                        <td
                                                            {...cell.getCellProps()}
                                                            className={
                                                                `p-3 
                                                                ${cell.column.id == "gpuname" ? ` w-9/12 text-lg font-bold ${row.original.gpuvendor === "NVIDIA" ? 'text-green-400' : 'text-red-500'}` : 'text-base text-white text-center'}
                                                                ${cell.column.id == "hashrate" ? 'w-2/12' : ''}
                                                                `
                                                            }
                                                        >
                                                            {cell.value ? cell.value : "0"}
                                                            {
                                                                cell.column.id == "gpuname" ?
                                                                    <p className={`text-sm font-normal ${row.original.gpuvendor === "NVIDIA" ? 'text-green-600' : 'text-red-700'}`}>
                                                                        {row.original.memorysize}
                                                                    </p> :
                                                                    null
                                                            }
                                                        </td>
                                                    )
                                                })}
                                            </tr>
                                        )
                                    })}
                                </tbody>
                            </table>
                        </Content>
                    </ContentContext.Provider>

                    {
                        showOverclockModal ?
                            <OverclockModal setShowModal={setShowOverclockModal} gpu={gpu} overclockSuccessHandler={(res: any, idx: number) => {
                                setShowOverclockModal(false)
                                gpusData[idx] = { ...gpusData[idx], ...res.oc }
                                setGpusData([...gpusData])
                            }} /> : null
                    }
                </div>
            </>
            : null}
    </div>
}

export const getServerSideProps: GetServerSideProps = async (ctx) => {
    const { rig_id } = ctx.query
    const rtoken = ctx.req.cookies['rtoken']

    const getRigsData = () => {
        return new Promise(resolve => {
            withAuth(async (token: jwtObject) => {
                const response = await fetch(`${isServer() ? process.env.API_ENDPOINT_SSR : process.env.API_ENDPOINT}/api/v1/getRig/${rig_id}`, {
                    method: 'GET',
                    mode: 'cors',
                    headers: {
                        Authorization: `Bearer ${token.jwt_token}`
                    }
                })

                if (response.status == 200) {
                    const responseJSON = await response.json()
                    resolve({
                        ...responseJSON,
                        responseStatus: response.status
                    })
                }
                else {
                    resolve({ responseStatus: response.status })
                }
            }, isServer(), rtoken)
        })
    }

    const data: any = await getRigsData()
    let amdIdx = -1, nvidiaIdx = -1;

    data.status.gpus.map((ele: any, idx: number) => {
        let ocinfo = {
            fs: 0,
            cc: 0,
            cv: 0,
            mc: 0,
            mv: 0,
            pl: 0
        }

        let ocID = 0;
        if (ele.gpuvendor == "AMD") {
            amdIdx++
            ocID = amdIdx;
        } else if (ele.gpuvendor == "NVIDIA") {
            nvidiaIdx++
            ocID = nvidiaIdx;
        }

        if (typeof data.oc != "undefined") {
            if (typeof data.oc[ele.gpuvendor] !== "undefined") {
                if (typeof data.oc[ele.gpuvendor][ocID] !== "undefined") {
                    ocinfo = data.oc[ele.gpuvendor][ocID]
                }
            }
        }

        data.status.gpus[idx] = { ...data.status.gpus[idx], ...ocinfo }
    })
    return {
        props: { data }
    }
}

export default Rigs