import type { NextPage } from 'next'
import { useEffect, useState } from 'react'
import { checkAuth, withAuth, jwtObject } from "../../utils/auth"

import Navbar from '../../components/navbar'

const Dashboard: NextPage = () => {
    const [isAuth, setIsAuth] = useState(false)

    useEffect(() => {
        checkAuth().then(auth => {
            if (auth) {
                setIsAuth(true)
            }
        })
    }, [])

    return (
        <div>
            {isAuth ?
                <>
                    <Navbar />
                    <h1>Hello World!!</h1>
                    <button onClick={() => {
                        withAuth((token: jwtObject) => {
                            fetch("http://localhost:5000/api/v1/hello", {
                                method: 'GET',
                                mode: 'cors',
                                headers: {
                                    'Accept': 'application/json',
                                    'Content-Type': 'application/json',
                                    'Authorization': `Bearer ${token.jwt_token}`,
                                }
                            }).then(response => {
                                response.json().then(responseJSON => console.log(responseJSON))
                            })
                        })
                    }}>nani</button>
                </>
                :
                <></>
            }
        </div >
    )
}

export default Dashboard