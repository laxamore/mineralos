import { Dispatch, FC, FormEvent, SetStateAction } from "react"
import { jwtObject, withAuth } from "../../utils/auth"
import Modal from "./modal"

type Props = {
    setShowModal: Dispatch<SetStateAction<boolean>>
    createRigSuccessHandler?: Function
    createRigFailedHandler?: Function
}

const CreateRigModal: FC<Props> = ({ setShowModal, createRigSuccessHandler, createRigFailedHandler }) => {
    return <Modal setShowModal={setShowModal}>
        <form onSubmit={(e: FormEvent<HTMLFormElement>) => {
            e.preventDefault()
            const formData = new FormData(e.currentTarget)
            const data = {
                rig_name: formData.get("rig_name"),
            }

            if (data.rig_name != "") {
                withAuth(async (token: jwtObject) => {
                    const response = await fetch(`${process.env.API_ENDPOINT}/api/v1/newRig`, {
                        method: 'POST',
                        mode: "cors",
                        body: JSON.stringify({
                            rig_name: data.rig_name
                        }),
                        headers: {
                            Authorization: `Bearer ${token.jwt_token}`
                        }
                    })

                    if (response.status == 200) {
                        const responseJSON = await response.json()
                        if (typeof createRigSuccessHandler != 'undefined') {
                            createRigSuccessHandler(responseJSON)
                        }
                    }
                    else {
                        if (typeof createRigFailedHandler != 'undefined') {
                            createRigFailedHandler()
                        }
                    }
                })
            }
        }}>
            <div className="flex flex-col mx-4 mb-4">
                <label htmlFor="rig_name" className="text-white font-semibold py-2">RIG Name</label>
                <input type="text" name="rig_name" required />
                <button type="submit" className="bg-blue-500 hover:bg-blue-600 text-white py-1 mt-4 rounded-lg">Create RIG</button>
            </div>
        </form>
    </Modal>
}

export default CreateRigModal