import Router from "next/router";
import { useEffect } from "react";

const Page404 = () => {
  useEffect(() => {
    Router.push("/404.html");
  }, []);

  return <></>;
};

export default Page404;
