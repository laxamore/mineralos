const { parsed: localEnv } = require('dotenv').config({
  path: '../.env',
})

/** @type {import('next').NextConfig} */
module.exports = {
  reactStrictMode: true,
  env: localEnv
}
