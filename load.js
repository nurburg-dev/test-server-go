import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 10 }, // Ramp up to 10 VUs in 10 seconds
    { duration: '30s', target: 10 }, // Stay at 10 VUs for 30 seconds
    { duration: '10s', target: 0 },  // Ramp down to 0 VUs in 10 seconds
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.01'],   // Error rate should be less than 1%
  },
};

export default function () {
  const res = http.get('http://test-service:9000/');
  check(res, {
    'is status 200': (r) => r.status === 200,
    'response time is acceptable': (r) => r.timings.duration < 500,
  });
  sleep(1); // Add a 1-second pause between requests
}
