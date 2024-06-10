import http from 'k6/http';
import { check } from 'k6';

export let options = {
    vus: 30, // Number of virtual users
    duration: '30s', // Duration of the test
};

export default function () {
    let res = http.get('http://localhost:9000/api/data?id=1');
    check(res, {
        'is status 200': (r) => r.status === 200,
    });
}