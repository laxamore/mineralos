import { useEffect, useState } from "react";
import Navbar from "../../components/navbar";
import { checkAuth, getAuthPayload } from "../../utils/auth";

const Overclokcing = () => {
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

    return <div>
        {
            isAuth ?
                <Navbar />
                : null
        }
    </div>
}

export default Overclokcing;