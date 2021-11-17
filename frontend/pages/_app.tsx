import '../styles/globals.css'
import type { AppProps } from 'next/app'
import { RefreshContext } from "../utils/context"
import { useState } from 'react'

function MyApp({ Component, pageProps }: AppProps) {
  const [refreshTimeout, setRefreshTimeout] = useState(false)
  return <RefreshContext.Provider value={[refreshTimeout, setRefreshTimeout]}>
    <Component {...pageProps} />
  </RefreshContext.Provider>
}

export default MyApp
