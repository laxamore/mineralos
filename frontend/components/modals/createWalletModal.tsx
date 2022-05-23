import { Dispatch, FC, FormEvent, SetStateAction } from "react";
import { jwtObject, withAuth } from "../../utils/auth";
import Modal from "./modal";

type Props = {
  setShowModal: Dispatch<SetStateAction<boolean>>;
  createWalletSuccessHandler?: Function;
  createWalletFailedHandler?: Function;
};

const CreateWalletModal: FC<Props> = ({
  setShowModal,
  createWalletSuccessHandler,
  createWalletFailedHandler,
}) => {
  return (
    <Modal setShowModal={setShowModal}>
      <form
        onSubmit={(e: FormEvent<HTMLFormElement>) => {
          e.preventDefault();
          const formData = new FormData(e.currentTarget);
          const data = {
            wallet_name: formData.get("wallet_name"),
            wallet_address: formData.get("wallet_address"),
            coin: formData.get("coin"),
          };

          if (
            data.wallet_name != "" &&
            data.wallet_address != "" &&
            data.coin != ""
          ) {
            withAuth(async (token: jwtObject) => {
              const response = await fetch(
                `${process.env.API_ENDPOINT}/api/v1/newWallet`,
                {
                  method: "POST",
                  mode: "cors",
                  body: JSON.stringify({
                    wallet_name: data.wallet_name,
                    wallet_address: data.wallet_address,
                    coin: data.coin,
                  }),
                  headers: {
                    Authorization: `Bearer ${token.jwt_token}`,
                  },
                }
              );

              if (response.status == 200) {
                const responseJSON = await response.json();
                if (typeof createWalletSuccessHandler != "undefined") {
                  createWalletSuccessHandler(responseJSON);
                }
              } else {
                if (typeof createWalletFailedHandler != "undefined") {
                  createWalletFailedHandler();
                }
              }
            });
          }
        }}
      >
        <div className="flex flex-col mx-4 mb-4 w-80">
          <label
            htmlFor="wallet_name"
            className="text-white font-semibold py-2"
          >
            Wallet Name
          </label>
          <input type="text" name="wallet_name" required />

          <label
            htmlFor="wallet_address"
            className="text-white font-semibold py-2"
          >
            Wallet Address
          </label>
          <input type="text" name="wallet_address" required />

          <label htmlFor="coin" className="text-white font-semibold py-2">
            Coin
          </label>
          <input type="text" name="coin" required />

          <button
            type="submit"
            className="bg-blue-500 hover:bg-blue-600 text-white py-1 mt-4 rounded-lg"
          >
            Create Wallet
          </button>
        </div>
      </form>
    </Modal>
  );
};

export default CreateWalletModal;
