import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

const BASE_URL = __ENV.BASE_URL;

export const options = {
    stages: [
        // Warm-up
        { duration: '30s', target: 10 },
        // Light load
        { duration: '1m', target: 25 },
        // Normal load
        { duration: '2m', target: 50 },
        // Peak load
        { duration: '3m', target: 100 },
        // Stress phase
        { duration: '2m', target: 200 },
        // Ramp down
        { duration: '30s', target: 0 },
    ],

    thresholds: {
        http_req_failed: [
            'rate<0.01',// <1% errors
            //'rate<0.005',
            //'count<10',
        ],
        'http_req_failed{group:single_blog_post}': ['rate<0.02'],
        http_req_duration: [
            'p(50)<400',
            'p(95)<500',
            'p(99)<1000',
        ],
        errors: ['rate<0.01'],
    },

    summaryTrendStats: [
        'avg',
        'min',
        'med',
        'max',
        'p(90)',
        'p(95)',
        'p(99)',
    ],
};

function validateResponse(response, name) {
    const success = check(response, {
        [`${name}: status is 2xx`]: (r) =>
            r.status >= 200 && r.status < 300,

        [`${name}: response time < 500ms`]: (r) =>
            r.timings.duration < 500,
    });

    errorRate.add(!success);

    return success;
}

export function setup() {
    console.log(`Testing ${BASE_URL}`);

    /*
    const loginResponse = http.post(
      `${BASE_URL}/api/auth/login`,
      JSON.stringify({
        email: 'test@example.com',
        password: 'password',
      }),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    validateResponse(loginResponse, 'login');

    return {
      token: loginResponse.json('token'),
    };
    */

    return {};
}

export default function (data) {
    const headers = {
        'Content-Type': 'application/json',
        // Authorization: `Bearer ${data.token}`,
    };

    group('single_blog_post', () => {
        const id = Math.floor(Math.random() * 1000) + 1;
        const res = http.get(`${BASE_URL}/blog/posts/${id}`, { headers });
        validateResponse(res, 'blog-post-details');
        sleep(1);
    });

    group('single_product', () => {
        const id = Math.floor(Math.random() * 1000) + 1;
        const res = http.get(`${BASE_URL}/products/${id}`, { headers });
        validateResponse(res, 'product-details');
        sleep(1);
    });
}

export function teardown(data) {
    console.log('Load test completed');
}
