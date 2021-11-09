import type { NextPage } from 'next'
import Router from 'next/router'
import { useEffect } from 'react'

const Login: NextPage = () => {
  useEffect(() => {
    Router.push('/login')
  }, [])

  return (
    <div >
    </div>
  )
}

export default Login
