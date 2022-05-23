import { useEffect, useMemo, useState } from "react"
import Navbar from "../../components/navbar"
import { Content } from "../../components/content"
import { ContentContext } from "../../utils/context"
import { checkAuth, getAuthPayload, jwtObject, withAuth } from "../../utils/auth"
import CreateWalletModal from "../../components/modals/createWalletModal"
import { GetServerSideProps, NextPage } from "next"
import { Column, useTable } from "react-table"
import Router from "next/router"

const isServer = () => typeof window === 'undefined';

type Props = {
    data?: any;
}

const Wallets: NextPage<Props> = ({ data }) => {
    const [isAuth, setIsAuth] = useState(false)
    const [privilege, setPrivilege] = useState('readOnly')
    const [showCreateWalletModal, setShowCreateWalletModal] = useState(false);
    const [walletsData, setWalletsData] = useState([]);

    interface walletsColumnInterface {
        wallet_name: string
        wallet_address: string,
        coin: string,
    }

    const columns = useMemo<Column<walletsColumnInterface>[]>(() => [
        {
            Header: 'Wallet Name',
            accessor: 'wallet_name',
        },
        {
            Header: 'Wallet Address',
            accessor: 'wallet_address',
        },
        {
            Header: 'Coin',
            accessor: 'coin'
        },
    ], []);

    useEffect(() => {
        checkAuth().then(auth => {
            if (auth) {
                setIsAuth(true)
                const authPayload = getAuthPayload()
                setPrivilege(authPayload.privilege)
            }
        })
        if (data.wallets) {
            setWalletsData(data.wallets)
        }
        else {
            setWalletsData([])
        }
    }, [data.wallets])

    const tableInstance = useTable({ columns: columns, data: walletsData })
    const {
        rows,
    } = tableInstance

    return <div>
        {
            isAuth ?
                <>
                    <Navbar />
                    <div className="flex flex-col justify-center items-center h-full w-full">
                        <ContentContext.Provider value={[setShowCreateWalletModal]}>
                            <Content showCreateButton={true} createButtonName={"Create New Wallet"} privilege={privilege} showRefreshButton={false}>
                                <ul className="mt-2">
                                    {
                                        rows.map((val: any) => {
                                            return <li className={`border border-gray-500 p-3 rounded-2xl cursor-pointer mt-2`}
                                                key={val.original._id}>
                                                <div className="flex flex-row items-center">
                                                    <div className="w-1/6">
                                                        <p className="text-lg font-bold">{val.original.wallet_name}</p>
                                                    </div>
                                                    <div className="w-1/6">
                                                        <p className="text-lg font-bold">{val.original.coin}</p>
                                                    </div>
                                                    <div className="w-4/6">
                                                        <p className="text-lg font-bold">{val.original.wallet_address}</p>
                                                    </div>
                                                    <div className="flex flex-grow flex-row-reverse ml-4">
                                                        <button className={`bg-red-500 rounded-full px-2 ${privilege === 'admin' || privilege === 'readAndWrite' ? 'hover:bg-red-600' : 'opacity-50 cursor-default'}`}
                                                            disabled={privilege === 'admin' || privilege === 'readAndWrite' ? false : true} onClick={(e) => {
                                                                e.preventDefault();
                                                                e.stopPropagation();

                                                                withAuth(async (token: jwtObject) => {
                                                                    const response = await fetch(`${process.env.API_ENDPOINT}/api/v1/deleteWallet`, {
                                                                        method: 'DELETE',
                                                                        mode: 'cors',
                                                                        body: JSON.stringify({
                                                                            wallet_id: val.original._id,
                                                                        }),
                                                                        headers: {
                                                                            Authorization: `Bearer ${token.jwt_token}`,
                                                                        },
                                                                    })

                                                                    if (response.status == 200) {
                                                                        setWalletsData(walletsData.filter((wallet: any) => wallet._id != val.original._id))
                                                                    }
                                                                })
                                                            }}>X</button>
                                                    </div>
                                                </div>
                                            </li>
                                        })
                                    }
                                </ul>
                            </Content>
                        </ContentContext.Provider>
                    </div>
                    {showCreateWalletModal ?
                        <CreateWalletModal setShowModal={setShowCreateWalletModal} createWalletSuccessHandler={(res: never) => {
                            setShowCreateWalletModal(false)
                            setWalletsData([...walletsData, res])
                        }} /> : null
                    }
                </>
                : null
        }
    </div>
}

export const getServerSideProps: GetServerSideProps = async (ctx: any) => {
    const cookie = ctx.req.headers.cookie
    const getWalletsData = () => {
        return new Promise(resolve => {
            withAuth(async (token: jwtObject) => {
                const response = await fetch(`${isServer() ? process.env.API_ENDPOINT_SSR : process.env.API_ENDPOINT}/api/v1/getWallets`, {
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
            }, isServer(), cookie)
        })
    }

    const data: any = await getWalletsData()
    return { props: { data } }
}

export default Wallets