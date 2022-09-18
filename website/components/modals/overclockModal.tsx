import { Dispatch, FC, FormEvent, SetStateAction } from "react";
import { jwtObject, withAuth } from "../../utils/auth";
import Modal from "./modal";

type Props = {
  setShowModal: Dispatch<SetStateAction<boolean>>;
  overclockSuccessHandler?: Function;
  overclockFailHandler?: Function;
  gpu: any;
};

const OverclockModal: FC<Props> = ({
  setShowModal,
  overclockSuccessHandler,
  overclockFailHandler,
  gpu,
}) => {
  return (
    <Modal setShowModal={setShowModal}>
      <form
        onSubmit={(e: FormEvent<HTMLFormElement>) => {
          e.preventDefault();
          const formData = new FormData(e.currentTarget);
          const data = {
            cc: formData.get("cc") ? parseInt(`${formData.get("cc")}`) : 0,
            cv: formData.get("cv") ? parseInt(`${formData.get("cv")}`) : 0,
            mc: formData.get("mc") ? parseInt(`${formData.get("mc")}`) : 0,
            mv: formData.get("mv") ? parseInt(`${formData.get("mv")}`) : 0,
            fs: formData.get("fs") ? parseInt(`${formData.get("fs")}`) : 0,
            pl: formData.get("pl") ? parseInt(`${formData.get("pl")}`) : 0,
          };

          if (data.pl < 0) data.pl = 0;
          else if (data.pl > 100) data.pl = 100;

          if (data.fs < 0) data.fs = 0;
          else if (data.fs > 100) data.fs = 100;

          withAuth(async (token: jwtObject) => {
            const response = await fetch(
              `${process.env.API_ENDPOINT}/api/v1/updateOC`,
              {
                method: "PUT",
                mode: "cors",
                body: JSON.stringify({
                  rig_id: gpu.rig_id,
                  vendor: gpu.gpuvendor,
                  id: gpu.ocID,
                  ...data,
                }),
                headers: {
                  Authorization: `Bearer ${token.jwt_token}`,
                },
              }
            );

            if (response.status == 200) {
              const responseJSON = await response.json();
              if (typeof overclockSuccessHandler != "undefined") {
                overclockSuccessHandler(responseJSON, gpu.idx);
              }
            } else {
              if (typeof overclockFailHandler != "undefined") {
                overclockFailHandler();
              }
            }
          });
        }}
      >
        {gpu.gpuvendor == "AMD" ? (
          <div className="flex flex-col mx-4 mb-4">
            <label
              htmlFor="gpu_name"
              className="text-white text-xl font-semibold pb-6"
            >
              {gpu.gpuinfo.gpuname}
            </label>

            <label htmlFor="cc" className="text-white font-semibold py-2">
              Core Clock
            </label>
            <input
              type="number"
              name="cc"
              required
              className="px-2"
              defaultValue={gpu.gpuinfo.cc ? gpu.gpuinfo.cc : 0}
            />

            <label htmlFor="cv" className="text-white font-semibold py-2">
              Core Volt
            </label>
            <input
              type="text"
              name="cv"
              required
              className="px-2"
              defaultValue={gpu.gpuinfo.cv ? gpu.gpuinfo.cv : 0}
            />

            <label htmlFor="mc" className="text-white font-semibold py-2">
              Mem Clock
            </label>
            <input
              type="text"
              name="mc"
              required
              className="px-2"
              defaultValue={gpu.gpuinfo.mc ? gpu.gpuinfo.mc : 0}
            />

            <label htmlFor="mv" className="text-white font-semibold py-2">
              Mem Volt
            </label>
            <input
              type="text"
              name="mv"
              required
              className="px-2"
              defaultValue={gpu.gpuinfo.mv ? gpu.gpuinfo.mv : 0}
            />

            <label htmlFor="fs" className="text-white font-semibold py-2">
              Fan Speed
            </label>
            <input
              type="text"
              name="fs"
              required
              className="px-2"
              defaultValue={gpu.gpuinfo.fs ? gpu.gpuinfo.fs : 0}
            />

            <label htmlFor="pl" className="text-white font-semibold py-2">
              Power Limit
            </label>
            <input
              type="text"
              name="pl"
              required
              className="px-2"
              defaultValue={gpu.gpuinfo.pl ? gpu.gpuinfo.pl : 0}
            />

            <button
              type="submit"
              className="bg-blue-500 hover:bg-blue-600 text-white py-1 mt-6 rounded-lg"
            >
              Update Overclock
            </button>
          </div>
        ) : (
          <div className="flex flex-col mx-4 mb-4">
            <label
              htmlFor="gpu_name"
              className="text-white text-xl font-semibold pb-6"
            >
              {gpu.gpuinfo.gpuname}
            </label>

            <label htmlFor="cc" className="text-white font-semibold py-2">
              Core Clock
            </label>
            <input
              type="number"
              name="cc"
              required
              className="px-2"
              defaultValue={gpu.gpuinfo.cc ? gpu.gpuinfo.cc : 0}
            />

            <label htmlFor="mc" className="text-white font-semibold py-2">
              Mem Clock
            </label>
            <input
              type="text"
              name="mc"
              required
              className="px-2"
              defaultValue={gpu.gpuinfo.mc ? gpu.gpuinfo.mc : 0}
            />

            <label htmlFor="fs" className="text-white font-semibold py-2">
              Fan Speed
            </label>
            <input
              type="text"
              name="fs"
              required
              className="px-2"
              defaultValue={gpu.gpuinfo.fs ? gpu.gpuinfo.fs : 0}
            />

            <label htmlFor="pl" className="text-white font-semibold py-2">
              Power Limit
            </label>
            <input
              type="text"
              name="pl"
              required
              className="px-2"
              defaultValue={gpu.gpuinfo.pl ? gpu.gpuinfo.pl : 0}
            />

            <button
              type="submit"
              className="bg-blue-500 hover:bg-blue-600 text-white py-1 mt-6 rounded-lg"
            >
              Update Overclock
            </button>
          </div>
        )}
      </form>
    </Modal>
  );
};

export default OverclockModal;
