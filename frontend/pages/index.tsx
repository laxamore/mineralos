import type { NextPage } from "next";
import Router from "next/router";
import { useEffect } from "react";
import { checkAuth } from "../utils/auth";

const Login: NextPage = () => {
  useEffect(() => {
    checkAuth().then((res) => {
      if (res) {
        Router.push("/dashboard");
      } else {
        Router.push("/login");
      }
    });
  }, []);

  return <div className="w-full h-full absolute top-0 left-0"></div>;
};

export default Login;
