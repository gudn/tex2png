import http from 'k6/http'
import { check } from 'k6'

function getBody() {
  return '$1+\\frac{1}{2}$'
}

export const options = {
  stages: [
    { duration: '5m', target: 400 },
  ],
}

export default function () {
  const url = 'http://localhost:3000/'
  const body = getBody()
  const res = http.post(url, body)
  check(res, { 'is status 200': r => r.status === 200 })
  check(res, {
    'is contentType is image': r => r.headers['Content-Type'] === 'image/png',
  })
}
