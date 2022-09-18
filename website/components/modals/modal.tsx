import { Dispatch, FC, SetStateAction } from "react";

type Props = {
  setShowModal: Dispatch<SetStateAction<boolean>>;
  children?: React.ReactNode;
};

const Modal: FC<Props> = ({ children, setShowModal }) => {
  return (
    <div className="flex justify-center items-center fixed w-full h-full top-0 left-0 z-10">
      <span
        className="absolute w-full h-full top-0 left-0  bg-gray-900 opacity-50"
        onClick={() => setShowModal(false)}
      ></span>
      <div className="relative flex flex-col z-20 bg-gray-700 rounded-2xl">
        <div className="flex flex-row justify-end items-center">
          <span
            className="bg-red-500 hover:bg-red-600 px-2 text-sm rounded-full text-white font-semibold m-2 cursor-pointer"
            onClick={() => setShowModal(false)}
          >
            X
          </span>
        </div>
        {children}
      </div>
    </div>
  );
};

export default Modal;
