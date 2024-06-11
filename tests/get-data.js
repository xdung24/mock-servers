import http from 'k6/http';
import { check } from 'k6';

// Test configuration
export const options = {
    maxRedirects: 1 ,
    thresholds: {
        'http_req_duration{status:0}': ['max>=0'],
        'http_req_duration{status:200}': ['max>=0'],
        'http_req_duration{status:400}': ['max>=0'],
        'http_req_duration{status:500}': ['max>=0'],
        'http_req_duration{status:502}': ['max>=0'],
        'http_req_duration{method:POST}': ['max>=0'],
    },  
    summaryTrendStats: ['min', 'med', 'avg', 'p(90)', 'p(95)', 'max', 'count'],
    discardResponseBodies: true,
    // Ramp the number of virtual users up and down
    stages: [
      { duration: "30s", target: 10 },
      { duration: "30s", target: 20 },
      { duration: "20s", target: 0 },
    ],
  };
  
export default function () {
    let res = http.get('http://192.168.1.14:9000/api/data?id=1');
    check(res, {
        'is status 200': (r) => r.status === 200,
    });
}