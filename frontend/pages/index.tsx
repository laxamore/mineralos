import type { NextPage } from 'next'
import Router from 'next/router'
import { useEffect } from 'react'

const Login: NextPage = () => {
  useEffect(() => {
    Router.push('/login')
  }, [])

  return (
    <div className="w-full h-full absolute top-0 left-0">
    </div>
  )
}

export default Login
