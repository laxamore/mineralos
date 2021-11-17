const { parsed: mainEnv } = require('dotenv').config({
  path: '../.env',
})

/** @type {import('next').NextConfig} */
module.exports = {
  reactStrictMode: true,
  env: {
    ...mainEnv,
    API_ENDPOINT_SSR: process.env.DOCKER ? 'http://backend_api:5000' : mainEnv.API_ENDPOINT,
    PORT: process.env.PORT
  }
}
