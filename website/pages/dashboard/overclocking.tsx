import { useEffect, useState } from "react";
import Navbar from "../../components/navbar";
import { checkAuth, getAuthPayload } from "../../utils/auth";
import Router from "next/router";

const Overclokcing = () => {
  const [isAuth, setIsAuth] = useState(false);
  const [privilege, setPrivilege] = useState("readOnly");

  useEffect(() => {
    checkAuth().then((auth) => {
      if (auth) {
        setIsAuth(true);
        const authPayload = getAuthPayload();
        setPrivilege(authPayload.privilege);
      } else {
        Router.push("/login");
      }
    });
  }, []);

  return <div>{isAuth ? <Navbar /> : null}</div>;
};

export default Overclokcing;
