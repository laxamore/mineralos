import { GetServerSideProps, NextPage } from "next"
import Router from "next/router";
import { useEffect, useState } from "react"
import { checkAuth, getAuthPayload, jwtObject, withAuth } from "../../../utils/auth"
import Navbar from '../../../components/navbar'
import { Content } from '../../../components/content'
import { ContentContext } from "../../../utils/context"
import { EyeIcon, EyeOffIcon } from "@heroicons/react/solid"


const isServer = () => typeof window === 'undefined';

type Params = {
    data: any
}

const Rigs: NextPage<Params> = ({ data }) => {
    const [isAuth, setIsAuth] = useState(false)
    const [privilege, setPrivilege] = useState('readOnly')
    const [showKey, setShowKey] = useState(false)

    useEffect(() => {
        if (data.status === 404) {
            Router.push('/404.html')
        }

        checkAuth().then(auth => {
            if (auth) {
                setIsAuth(true)
                const authPayload = getAuthPayload()
                setPrivilege(authPayload.privilege)
            }
        })
    }, [])

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
                                    <p>: {data.rig_id}</p>
                                </div>
                                <div className="flex flex-row items-center">
                                    <div className="flex flex-col">
                                        <div className="flex flex-row">
                                            <p className="w-32">RIG Key</p>
                                            <p className={`${!showKey ? 'tracking-tight' : ''}`}>: {showKey ? data.key : "∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗"}</p>
                                        </div>
                                        <div className="flex flex-row">
                                            <p className="w-32">RIG Public Key</p>
                                            <p className={`${!showKey ? 'tracking-tight' : ''}`}>: {showKey ? data.pubkey : "∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗∗"}</p>
                                        </div>
                                    </div>
                                    <div className="w-6 mx-4 cursor-pointer" onClick={() => setShowKey(!showKey)}>
                                        {
                                            showKey ? <EyeOffIcon /> : <EyeIcon />
                                        }
                                    </div>
                                </div>
                            </div>
                        </Content>
                    </ContentContext.Provider>
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
                        status: response.status
                    })
                }
                else {
                    resolve({ status: response.status })
                }
            }, isServer(), rtoken)
        })
    }

    const data: any = await getRigsData()
    return {
        props: { data }
    }
}

export default Rigs